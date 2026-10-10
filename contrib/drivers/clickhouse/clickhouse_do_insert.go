// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// DoInsert inserts or updates data for given table.
// The list parameter must contain at least one record, which was previously validated.
//
// The records are grouped by their fields and inserted on `link`, and the result reports the
// number of records inserted.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	// Save and Replace insert the records, which a ReplacingMergeTree table deduplicates, but
	// inserting them cannot skip the conflicting records as InsertIgnore does.
	if option.InsertOption == gdb.InsertOptionIgnore {
		return nil, errUnsupportedInsertIgnore
	}
	if link == nil {
		if link, err = d.MasterLink(); err != nil {
			return nil, err
		}
	}
	var (
		groupKeys []string
		groups    = make(map[string]gdb.List)
	)
	for _, record := range list {
		keys := make([]string, 0, len(record))
		for key := range record {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		groupKey := strings.Join(keys, ",")
		if _, ok := groups[groupKey]; !ok {
			groupKeys = append(groupKeys, groupKey)
		}
		groups[groupKey] = append(groups[groupKey], record)
	}
	for _, groupKey := range groupKeys {
		err = d.insertRecords(ctx, link, table, strings.Split(groupKey, ","), groups[groupKey], option.BatchCount)
		if err != nil {
			return nil, err
		}
	}
	return &Result{rowsAffected: int64(len(list))}, nil
}

// insertRecords inserts `records`, which all have the fields `keys`, into `table` on `link`.
//
// The records are sent in one batch of the underlying driver, which binds the values of every
// column type as they are. Binding the values into a statement instead loses some of them, like
// a time outside the range of DateTime or binary data, and the core expands a slice argument
// into several arguments. A record with a gdb.Raw value is inserted through DoExec instead, with
// the raw value written into the statement, at most `batchCount` records a statement.
func (d *Driver) insertRecords(
	ctx context.Context, link gdb.Link, table string, keys []string, records gdb.List, batchCount int,
) error {
	var (
		charL, charR = d.GetChars()
		insertSql    = fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES ",
			d.QuotePrefixTableName(table), charL+strings.Join(keys, charR+","+charL)+charR,
		)
		batch   = make(gdb.List, 0, len(records))
		holders = make([]string, 0)
		params  = make([]any, 0)
	)
	if batchCount <= 0 {
		batchCount = len(records)
	}
	flush := func() error {
		if len(holders) == 0 {
			return nil
		}
		_, err := d.DoExec(ctx, link, insertSql+strings.Join(holders, ","), params...)
		holders, params = holders[:0], params[:0]
		return err
	}
	for _, record := range records {
		if !hasRawValue(record) {
			batch = append(batch, record)
			continue
		}
		values := make([]string, 0, len(keys))
		for _, key := range keys {
			switch value := record[key].(type) {
			case gdb.Raw:
				values = append(values, string(value))
			case *gdb.Raw:
				values = append(values, string(*value))
			default:
				values = append(values, "?")
				params = append(params, value)
			}
		}
		holders = append(holders, "("+strings.Join(values, ",")+")")
		if len(holders) >= batchCount {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	if err := flush(); err != nil {
		return err
	}
	if len(batch) == 0 {
		return nil
	}
	beginner, ok := link.(interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	})
	if !ok {
		return gerror.NewCodef(gcode.CodeNotSupported, `link of type "%T" cannot start an insert batch`, link)
	}
	tx, err := beginner.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	batchSql := insertSql + "(" + strings.TrimSuffix(strings.Repeat("?,", len(keys)), ",") + ")"
	if err = d.insertBatch(ctx, &batchLink{tx: tx}, batchSql, keys, batch); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// insertBatch prepares `insertSql` on `link` and executes it for each of `records`.
func (d *Driver) insertBatch(ctx context.Context, link gdb.Link, insertSql string, keys []string, records gdb.List) error {
	stmt, err := d.DoPrepare(ctx, link, insertSql)
	if err != nil {
		return err
	}
	defer func() { _ = stmt.Close() }()
	for _, record := range records {
		values := make([]any, 0, len(keys))
		for _, key := range keys {
			values = append(values, record[key])
		}
		if _, err = stmt.ExecContext(ctx, values...); err != nil {
			return err
		}
	}
	return nil
}

// hasRawValue reports whether `record` has a gdb.Raw value.
func hasRawValue(record gdb.Map) bool {
	for _, value := range record {
		switch value.(type) {
		case gdb.Raw, *gdb.Raw:
			return true
		}
	}
	return false
}

// batchLink is the gdb.Link of the transaction that an insert batch of the underlying driver
// is sent in.
type batchLink struct {
	tx *sql.Tx
}

// QueryContext executes a query on the transaction of the batch.
func (l *batchLink) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return l.tx.QueryContext(ctx, query, args...)
}

// ExecContext executes a statement on the transaction of the batch.
func (l *batchLink) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return l.tx.ExecContext(ctx, query, args...)
}

// PrepareContext prepares the insert statement of the batch.
func (l *batchLink) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return l.tx.PrepareContext(ctx, query)
}

// IsOnMaster reports true, as inserts run on the master.
func (l *batchLink) IsOnMaster() bool {
	return true
}

// IsTransaction reports true, as the batch is sent in a transaction of the underlying driver.
func (l *batchLink) IsTransaction() bool {
	return true
}

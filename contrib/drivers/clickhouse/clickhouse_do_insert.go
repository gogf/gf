// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
)

// DoInsert inserts or updates data for given table.
// The list parameter must contain at least one record, which was previously validated.
//
// The records are inserted on `link` in their order, each run of records with the same fields
// and either all or none with a gdb.Raw value in batches of at most `option.BatchCount` records,
// and the result reports the number of records inserted. ClickHouse has no transactions, so the
// batches inserted before a batch that fails are kept.
//
// Save and Replace insert the records, which a ReplacingMergeTree table deduplicates, but
// inserting them cannot skip the conflicting records, so InsertIgnore fails as unsupported.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	if option.InsertOption == gdb.InsertOptionIgnore {
		return nil, errUnsupportedInsertIgnore
	}
	fields, err := d.tableFields(ctx, table)
	if err != nil {
		return nil, err
	}
	if list, err = d.convertDecimalValues(ctx, fields, list); err != nil {
		return nil, err
	}
	if link == nil {
		if link, err = d.MasterLink(); err != nil {
			return nil, err
		}
	}
	for start := 0; start < len(list); {
		var (
			keys = recordKeys(list[start])
			raw  = hasRawValue(list[start])
			end  = start + 1
		)
		for end < len(list) && hasRawValue(list[end]) == raw && hasKeys(list[end], keys) {
			end++
		}
		if raw {
			err = d.insertRawRecords(ctx, link, table, keys, list[start:end], option.BatchCount, fields)
		} else {
			err = d.insertRecords(ctx, link, table, keys, list[start:end], option.BatchCount)
		}
		if err != nil {
			return nil, err
		}
		start = end
	}
	return &Result{rowsAffected: int64(len(list))}, nil
}

// insertSql returns the statement inserting the fields `keys` into `table`, up to its VALUES.
func (d *Driver) insertSql(table string, keys []string) string {
	charL, charR := d.GetChars()
	return fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES ",
		d.QuotePrefixTableName(table), charL+strings.Join(keys, charR+","+charL)+charR,
	)
}

// insertRawRecords inserts `records`, which all have the fields `keys` and a gdb.Raw value,
// into `table`, whose fields are `fields`, on `link` through DoExec, with the raw values written
// into the statement and the other values as rawValueSql returns them, at most `batchCount`
// records a statement.
func (d *Driver) insertRawRecords(
	ctx context.Context, link gdb.Link, table string, keys []string, records gdb.List, batchCount int,
	fields map[string]*gdb.TableField,
) error {
	var (
		insertSql = d.insertSql(table, keys)
		holders   = make([]string, 0)
		params    = make([]any, 0)
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
		values := make([]string, 0, len(keys))
		for _, key := range keys {
			switch value := record[key].(type) {
			case gdb.Raw:
				values = append(values, string(value))
			case *gdb.Raw:
				values = append(values, string(*value))
			default:
				var fieldType string
				if field := fields[key]; field != nil {
					fieldType = field.Type
				}
				valueSql, valueParams, err := rawValueSql(fieldType, value)
				if err != nil {
					return err
				}
				values = append(values, valueSql)
				params = append(params, valueParams...)
			}
		}
		holders = append(holders, "("+strings.Join(values, ",")+")")
		if len(holders) >= batchCount {
			if err := flush(); err != nil {
				return err
			}
		}
	}
	return flush()
}

// rawValueSql returns the SQL of `value` for a column of type `fieldType` in the statement of a
// record with a gdb.Raw value, and the arguments of its placeholders.
//
// The underlying driver binds a time to the second within the range of DateTime and binary data
// as an array, and the core expands a slice argument into several arguments. So a time is written
// as a literal exact to the nanosecond, or as its date in its own time zone for a Date or Date32
// column as the underlying driver writes it, binary data is bound as a string, and an Array, Map
// or Tuple value fails with an error of code gcode.CodeNotSupported.
func rawValueSql(fieldType string, value any) (string, []any, error) {
	if t, ok := timeOf(value); ok {
		switch fieldTypeFamily(fieldType) {
		case "date":
			return "toDate('" + t.Format(time.DateOnly) + "')", nil, nil
		case "date32":
			return "toDate32('" + t.Format(time.DateOnly) + "')", nil, nil
		}
		timeSql, err := timeSql(t)
		return timeSql, nil, err
	}
	switch v := value.(type) {
	case []byte:
		return "?", []any{string(v)}, nil
	case fmt.Stringer:
	default:
		switch reflect.Indirect(reflect.ValueOf(value)).Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return "", nil, gerror.NewCodef(
				gcode.CodeNotSupported, `value of type "%T" cannot be inserted with a gdb.Raw value`, value,
			)
		}
	}
	return "?", []any{value}, nil
}

// timeOf returns the time of `value` if it is a time.Time, a gtime.Time or a non-nil pointer to
// one of them.
func timeOf(value any) (time.Time, bool) {
	switch v := value.(type) {
	case time.Time:
		return v, true
	case *time.Time:
		if v != nil {
			return *v, true
		}
	case gtime.Time:
		return v.Time, true
	case *gtime.Time:
		if v != nil {
			return v.Time, true
		}
	}
	return time.Time{}, false
}

// timeSql returns the SQL of the time `t` as a DateTime64 of its instant exact to the nanosecond,
// which a column of any time zone converts without loss. It returns an error of code
// gcode.CodeInvalidParameter if `t` has nanoseconds and is out of the range of DateTime64(9).
func timeSql(t time.Time) (string, error) {
	switch {
	case t.Nanosecond()%int(time.Millisecond) == 0:
		return "fromUnixTimestamp64Milli(" + strconv.FormatInt(t.UnixMilli(), 10) + ")", nil
	case t.Nanosecond()%int(time.Microsecond) == 0:
		return "fromUnixTimestamp64Micro(" + strconv.FormatInt(t.UnixMicro(), 10) + ")", nil
	case t.Before(time.Unix(0, math.MinInt64)) || t.After(time.Unix(0, math.MaxInt64)):
		return "", gerror.NewCodef(gcode.CodeInvalidParameter, `time %s is out of the range of DateTime64(9)`, t)
	default:
		return "fromUnixTimestamp64Nano(" + strconv.FormatInt(t.UnixNano(), 10) + ")", nil
	}
}

// insertRecords inserts `records`, which all have the fields `keys`, into `table` on `link`, in
// batches of at most `batchCount` records.
//
// The records are sent in batches of the underlying driver, which binds the values of every
// column type as they are. Binding the values into a statement instead loses some of them, like
// a time outside the range of DateTime or binary data, and the core expands a slice argument
// into several arguments. Each batch is the INSERT statement of its records for the core, as
// sendBatch describes.
//
// In dry run, which starting a batch does not honor, each batch goes through DoExec as its
// INSERT statement instead, which the core logs and catches without executing it. Otherwise,
// under CatchSQL and ToSQL, the INSERT statement of each batch, formatted with its values as the
// core formats the statements it catches, is caught, and the batch is sent under CatchSQL only,
// as ToSQL must not execute it.
func (d *Driver) insertRecords(
	ctx context.Context, link gdb.Link, table string, keys []string, records gdb.List, batchCount int,
) error {
	if batchCount <= 0 {
		batchCount = len(records)
	}
	var (
		holder   = "(" + strings.TrimSuffix(strings.Repeat("?,", len(keys)), ",") + ")"
		batchSql = d.insertSql(table, keys) + holder
		manager  = catchSQLManager(ctx)
	)
	for start := 0; start < len(records); start += batchCount {
		var (
			batch     = records[start:min(start+batchCount, len(records))]
			statement = batchSql + strings.Repeat(","+holder, len(batch)-1)
			values    []any
		)
		for _, record := range batch {
			for _, key := range keys {
				values = append(values, record[key])
			}
		}
		if d.GetDryRun() {
			if _, err := d.DoExec(ctx, link, statement, values...); err != nil {
				return err
			}
			continue
		}
		if manager != nil {
			manager.SQLArray.Append(gdb.FormatSqlWithArgs(statement, values))
			if !manager.DoCommit {
				continue
			}
		}
		if err := d.sendBatch(ctx, link, batchSql, statement, len(keys), values); err != nil {
			return err
		}
	}
	return nil
}

// ctxKeyCatchSQL is the context key of the gdb.CatchSQLManager that the core passes the statements
// of CatchSQL and ToSQL to.
const ctxKeyCatchSQL gctx.StrKey = "CtxKeyCatchSQL"

// catchSQLManager returns the gdb.CatchSQLManager of the CatchSQL or ToSQL that `ctx` runs in, or
// nil if it runs in neither.
func catchSQLManager(ctx context.Context) *gdb.CatchSQLManager {
	manager, _ := ctx.Value(ctxKeyCatchSQL).(*gdb.CatchSQLManager)
	return manager
}

// sendBatch sends the records whose `width` values each are `values` in a batch of the underlying
// driver started on `link`, which prepares `insertSql` and appends the records.
//
// The batch is committed through DoCommit as the INSERT statement `statement` with `values`, so it
// is logged and traced once, as that statement, and fails with an error of code
// gcode.CodeDbOperationError carrying the statement, as a statement does; it fails so before its
// commit too, and on a panic of the underlying driver, which the core recovers from as well. The
// ExecTimeout of the configuration limits the whole batch, from its start to its commit, and the
// PrepareTimeout does not apply separately.
func (d *Driver) sendBatch(
	ctx context.Context, link gdb.Link, insertSql, statement string, width int, values []any,
) (err error) {
	beginner, ok := link.(interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	})
	if !ok {
		return gerror.NewCodef(gcode.CodeNotSupported, `link of type "%T" cannot start an insert batch`, link)
	}
	if timeout := d.GetConfig().ExecTimeout; timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	tx, err := beginner.BeginTx(ctx, nil)
	if err != nil {
		return gerror.WrapCode(gcode.CodeDbOperationError, err, gdb.FormatSqlWithArgs(statement, values))
	}
	defer func() {
		exception := recover()
		if exception == nil {
			return
		}
		if v, ok := exception.(error); ok && gerror.HasStack(v) {
			err = v
		} else {
			err = gerror.WrapCode(
				gcode.CodeDbOperationError,
				gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception),
				gdb.FormatSqlWithArgs(statement, values),
			)
		}
		if rollbackErr := tx.Rollback(); rollbackErr != nil && !errors.Is(rollbackErr, sql.ErrTxDone) {
			err = errors.Join(err, rollbackErr)
		}
	}()
	if err = appendBatch(ctx, tx, insertSql, width, values); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.Join(err, rollbackErr)
		}
		return gerror.WrapCode(gcode.CodeDbOperationError, err, gdb.FormatSqlWithArgs(statement, values))
	}
	_, err = d.DoCommit(ctx, gdb.DoCommitInput{
		Tx:            tx,
		Sql:           statement,
		Args:          values,
		Type:          gdb.SqlTypeTXCommit,
		IsTransaction: true,
	})
	return err
}

// appendBatch prepares `insertSql` on `tx`, which starts the batch of the underlying driver, and
// appends the records whose `width` values each are `values` to it.
func appendBatch(ctx context.Context, tx *sql.Tx, insertSql string, width int, values []any) (err error) {
	stmt, err := tx.PrepareContext(ctx, insertSql)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()
	for start := 0; start < len(values); start += width {
		if _, err = stmt.ExecContext(ctx, values[start:start+width]...); err != nil {
			return err
		}
	}
	return nil
}

// recordKeys returns the fields of `record` in ascending order.
func recordKeys(record gdb.Map) []string {
	keys := make([]string, 0, len(record))
	for key := range record {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// hasKeys reports whether `record` has exactly the fields `keys`.
func hasKeys(record gdb.Map, keys []string) bool {
	if len(record) != len(keys) {
		return false
	}
	for _, key := range keys {
		if _, ok := record[key]; !ok {
			return false
		}
	}
	return true
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

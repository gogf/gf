// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	internalIdentityFieldInCtx gctx.StrKey = "identity_field"
	fieldExtraIdentity         string      = "IDENTITY"
	// mergeSavePoint is the savepoint that the MERGE statements of several records roll back to in
	// the transaction of the caller if one of them fails.
	mergeSavePoint string = "gf_merge"
)

// DoInsert inserts or updates data for given table.
// The list parameter must contain at least one record, which was previously validated.
// The table is written to without the alias that `table` may give it, which SQL Server rejects in
// INSERT and MERGE.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	table = tableWithoutAlias(table)
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		// MSSQL does not support REPLACE INTO syntax, use SAVE instead.
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionIgnore:
		// MSSQL does not support INSERT IGNORE syntax, use MERGE instead.
		return d.doInsertIgnore(ctx, link, table, list, option)

	default:
		identityField, err := d.getIdentityField(ctx, table)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, internalIdentityFieldInCtx, identityField)
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
}

// getIdentityField returns the name of the IDENTITY column of `table`, or an empty string if the
// table has none.
func (d *Driver) getIdentityField(ctx context.Context, table string) (string, error) {
	tableFields, err := d.GetDB().TableFields(ctx, table)
	if err != nil {
		return "", err
	}
	for _, field := range tableFields {
		if field.Extra == fieldExtraIdentity {
			return field.Name, nil
		}
	}
	return "", nil
}

// doSave support upsert for MSSQL
func (d *Driver) doSave(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, true)
}

// doInsertIgnore implements INSERT IGNORE operation using MERGE statement for MSSQL database.
// It only inserts records when there's no conflict on primary/unique keys.
func (d *Driver) doInsertIgnore(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, false)
}

// doMergeInsert implements MERGE-based insert operations for MSSQL database, executing one MERGE
// statement for each record of `list`, all of them in one transaction, or back to a savepoint of
// the transaction of `link` or `ctx` if there is one, so that a failing record writes none of them.
// When withUpdate is true, it performs upsert (insert or update).
// When withUpdate is false, it performs insert ignore (insert only when no conflict).
func (d *Driver) doMergeInsert(
	ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption, withUpdate bool,
) (result sql.Result, err error) {
	if len(list) > 1 && !d.GetDryRun() {
		if (link == nil || !link.IsTransaction()) && gdb.TXFromCtx(ctx, d.GetGroup()) == nil {
			err = d.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
				result, err = d.doMergeRecords(ctx, link, table, list, option, withUpdate)
				return err
			})
			return result, err
		}
		if _, err = d.DoExec(ctx, link, "SAVE TRANSACTION "+mergeSavePoint); err != nil {
			return nil, err
		}
		if result, err = d.doMergeRecords(ctx, link, table, list, option, withUpdate); err != nil {
			if _, e := d.DoExec(ctx, link, "ROLLBACK TRANSACTION "+mergeSavePoint); e != nil {
				err = errors.Join(err, e)
			}
		}
		return result, err
	}
	return d.doMergeRecords(ctx, link, table, list, option, withUpdate)
}

// doMergeRecords executes one MERGE statement for each record of `list`, as doMergeInsert
// describes, and returns the sum of their affected rows.
func (d *Driver) doMergeRecords(
	ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption, withUpdate bool,
) (result sql.Result, err error) {
	// If OnConflict is not specified, automatically get the primary key of the table
	conflictKeys := option.OnConflict
	if len(conflictKeys) == 0 {
		primaryKeys, err := d.Core.GetPrimaryKeys(ctx, table)
		if err != nil {
			return nil, gerror.WrapCode(
				gcode.CodeInternalError,
				err,
				`failed to get primary keys for table`,
			)
		}
		foundPrimaryKey := false
		for _, primaryKey := range primaryKeys {
			for dataKey := range list[0] {
				if strings.EqualFold(dataKey, primaryKey) {
					foundPrimaryKey = true
					break
				}
			}
			if foundPrimaryKey {
				break
			}
		}
		if !foundPrimaryKey {
			return nil, gerror.NewCodef(
				gcode.CodeMissingParameter,
				`Replace/Save/InsertIgnore operation requires conflict detection: `+
					`either specify OnConflict() columns or ensure table '%s' has a primary key in the data`,
				table,
			)
		}
		// TODO consider composite primary keys.
		conflictKeys = primaryKeys
	}

	var (
		charL, charR       = d.GetChars()
		conflictKeySet     = gset.NewStrSet(false)
		quotedConflictKeys = make([]string, len(conflictKeys))
		mergeResult        = new(Result)
	)

	// conflictKeys slice type conv to set type
	for index, conflictKey := range conflictKeys {
		conflictKeySet.Add(gstr.ToUpper(conflictKey))
		quotedConflictKeys[index] = charL + conflictKey + charR
	}

	for _, one := range list {
		var (
			oneLen = len(one)
			keys   = make([]string, 0, oneLen)

			// queryHolders:	Handle data with Holder that need to be merged
			// queryValues:		Handle data that need to be merged
			// insertKeys:		Handle valid keys that need to be inserted
			// insertValues:	Handle values that need to be inserted
			// updateValues:	Handle values that need to be updated (only when withUpdate=true)
			queryHolders = make([]string, oneLen)
			queryValues  = make([]any, 0, oneLen)
			insertKeys   = make([]string, oneLen)
			insertValues = make([]string, oneLen)
			updateValues []string
		)

		for key := range one {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for index, key := range keys {
			keyWithChar := charL + key + charR
			if s, ok := one[key].(gdb.Raw); ok {
				queryHolders[index] = gconv.String(s)
			} else {
				queryHolders[index] = "?"
				queryValues = append(queryValues, one[key])
			}
			insertKeys[index] = keyWithChar
			insertValues[index] = "T2." + keyWithChar
		}
		// Build updateValues only when withUpdate is true
		if withUpdate {
			updateValues = d.formatMergeUpdateValues(keys, conflictKeySet, option)
		}

		sqlStr := parseSqlForMerge(table, queryHolders, insertKeys, insertValues, updateValues, quotedConflictKeys)
		r, err := d.DoExec(ctx, link, sqlStr, queryValues...)
		if err != nil {
			return r, err
		}
		n, err := r.RowsAffected()
		if err != nil {
			return r, err
		}
		mergeResult.rowsAffected += n
	}
	return mergeResult, nil
}

// formatMergeUpdateValues returns the assignments of the MERGE UPDATE SET clause for a record
// with the given `keys`. It follows OnDuplicate/OnDuplicateEx of `option` if specified, or else
// updates every key except conflict keys, and except soft created fields for Save. A conflict key
// assigned from its own source column is left out, as it cannot change the matched row.
func (d *Driver) formatMergeUpdateValues(
	keys []string, conflictKeySet *gset.StrSet, option gdb.DoInsertOption,
) (updateValues []string) {
	charL, charR := d.GetChars()
	if option.OnDuplicateStr != "" {
		return []string{option.OnDuplicateStr}
	}
	if len(option.OnDuplicateMap) > 0 {
		updateKeys := make([]string, 0, len(option.OnDuplicateMap))
		for key := range option.OnDuplicateMap {
			updateKeys = append(updateKeys, key)
		}
		sort.Strings(updateKeys)
		for _, key := range updateKeys {
			keyWithChar := charL + key + charR
			switch value := option.OnDuplicateMap[key].(type) {
			case gdb.Raw, *gdb.Raw:
				updateValues = append(updateValues, fmt.Sprintf(`T1.%s = %s`, keyWithChar, gconv.String(value)))

			case gdb.Counter, *gdb.Counter:
				var counter gdb.Counter
				switch v := value.(type) {
				case gdb.Counter:
					counter = v
				case *gdb.Counter:
					counter = *v
				}
				operator, columnVal := "+", counter.Value
				if columnVal < 0 {
					operator, columnVal = "-", -columnVal
				}
				updateValues = append(updateValues, fmt.Sprintf(
					`T1.%s = T1.%s%s%s`,
					keyWithChar, charL+counter.Field+charR, operator, gconv.String(columnVal),
				))

			default:
				column := gconv.String(value)
				if conflictKeySet.Contains(gstr.ToUpper(key)) && strings.EqualFold(key, column) {
					continue
				}
				updateValues = append(updateValues, fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, charL+column+charR))
			}
		}
		return updateValues
	}
	for _, key := range keys {
		if conflictKeySet.Contains(gstr.ToUpper(key)) {
			continue
		}
		if option.InsertOption == gdb.InsertOptionSave && d.IsSoftCreatedFieldName(key) {
			continue
		}
		keyWithChar := charL + key + charR
		updateValues = append(updateValues, fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, keyWithChar))
	}
	return updateValues
}

// parseSqlForMerge generates MERGE statement for MSSQL database.
// When updateValues is empty, it only inserts (INSERT IGNORE behavior).
// When updateValues is provided, it performs upsert (INSERT or UPDATE).
// Examples:
// - INSERT IGNORE: MERGE INTO table WITH (HOLDLOCK) AS T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...)
// - UPSERT: MERGE INTO table WITH (HOLDLOCK) AS T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...) WHEN MATCHED THEN UPDATE SET ...
func parseSqlForMerge(table string,
	queryHolders, insertKeys, insertValues, updateValues, duplicateKey []string,
) (sqlStr string) {
	var (
		queryHolderStr  = strings.Join(queryHolders, ",")
		insertKeyStr    = strings.Join(insertKeys, ",")
		insertValueStr  = strings.Join(insertValues, ",")
		duplicateKeyStr string
	)

	// Build ON condition
	for index, keys := range duplicateKey {
		if index != 0 {
			duplicateKeyStr += " AND "
		}
		duplicateKeyStr += fmt.Sprintf("T1.%s = T2.%s", keys, keys)
	}

	// Build SQL based on whether UPDATE is needed
	pattern := gstr.Trim(
		`MERGE INTO %s WITH (HOLDLOCK) AS T1 USING (VALUES(%s)) T2 (%s) ON (%s) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s)`,
	)
	if len(updateValues) > 0 {
		// Upsert: INSERT or UPDATE
		pattern += ` WHEN MATCHED THEN UPDATE SET %s`
		return fmt.Sprintf(
			pattern+";",
			table,
			queryHolderStr,
			insertKeyStr,
			duplicateKeyStr,
			insertKeyStr,
			insertValueStr,
			strings.Join(updateValues, ","),
		)
	}
	// Insert Ignore: INSERT only
	return fmt.Sprintf(pattern+";", table, queryHolderStr, insertKeyStr, duplicateKeyStr, insertKeyStr, insertValueStr)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"database/sql"
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
	internalPrimaryKeyInCtx gctx.StrKey = "primary_key_field"
)

// DoInsert inserts or updates data for given table.
// The list parameter must contain at least one record, which was previously validated.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		// Oracle does not support REPLACE INTO syntax, use SAVE instead.
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionIgnore:
		// Oracle does not support INSERT IGNORE syntax, use MERGE instead.
		return d.doInsertIgnore(ctx, link, table, list, option)

	case gdb.InsertOptionDefault:
		// For default insert, set primary key field in context to support LastInsertId.
		// Only set it when the primary key is not provided in the data, for performance reason.
		tableFields, err := d.GetCore().GetDB().TableFields(ctx, table)
		if err == nil && len(list) > 0 {
			for _, field := range tableFields {
				if strings.EqualFold(field.Key, "pri") {
					// Check if primary key is provided in the data.
					pkProvided := false
					for key := range list[0] {
						if strings.EqualFold(key, field.Name) {
							pkProvided = true
							break
						}
					}
					// Only use RETURNING when primary key is not provided, for performance reason.
					if !pkProvided {
						pkField := *field
						ctx = context.WithValue(ctx, internalPrimaryKeyInCtx, pkField)
					}
					break
				}
			}
		}

	default:
	}
	var (
		batchResult  = new(gdb.SqlResult)
		charL, charR = d.GetChars()
	)
	// Format "INSERT...INTO..." statement.
	// Note: Use standard INSERT INTO syntax instead of INSERT ALL to ensure triggers fire
	for _, item := range list {
		var (
			keys   = make([]string, 0, len(item))
			values = make([]string, 0, len(item))
			params = make([]any, 0, len(item))
		)
		for k := range item {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if s, ok := item[k].(gdb.Raw); ok {
				values = append(values, gconv.String(s))
			} else {
				values = append(values, "?")
				params = append(params, item[k])
			}
		}

		// Execute individual INSERT for each record to trigger row-level triggers
		r, err := d.DoExec(ctx, link, fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES(%s)",
			table, charL+strings.Join(keys, charR+","+charL)+charR, strings.Join(values, ","),
		), params...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected += n
		}
	}
	return batchResult, nil
}

// doSave support upsert for Oracle
func (d *Driver) doSave(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, true)
}

// doInsertIgnore implements INSERT IGNORE operation using MERGE statement for Oracle database.
// It only inserts records when there's no conflict on primary/unique keys.
func (d *Driver) doInsertIgnore(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, false)
}

// doMergeInsert implements MERGE-based insert operations for Oracle database.
// When withUpdate is true, it performs upsert (insert or update).
// When withUpdate is false, it performs insert ignore (insert only when no conflict).
func (d *Driver) doMergeInsert(
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
		charL, charR   = d.GetChars()
		conflictKeySet = gset.NewStrSet(false)
		batchResult    = new(gdb.SqlResult)
	)

	// conflictKeys slice type conv to set type
	for _, conflictKey := range conflictKeys {
		conflictKeySet.Add(gstr.ToUpper(conflictKey))
	}

	for _, one := range list {
		var (
			oneLen = len(one)
			keys   = make([]string, 0, oneLen)

			// queryHolders:	Handle data with Holder that need to be upsert
			// queryValues:		Handle data that need to be upsert
			// insertKeys:		Handle valid keys that need to be inserted
			// insertValues:	Handle values that need to be inserted
			// updateValues:	Handle values that need to be updated
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
				queryHolders[index] = fmt.Sprintf("%s AS %s", gconv.String(s), keyWithChar)
			} else {
				queryHolders[index] = fmt.Sprintf("? AS %s", keyWithChar)
				queryValues = append(queryValues, one[key])
			}
			insertKeys[index] = keyWithChar
			insertValues[index] = fmt.Sprintf("T2.%s", keyWithChar)
		}
		// Build updateValues only when withUpdate is true
		if withUpdate {
			updateValues = d.formatMergeUpdateValues(keys, conflictKeySet, option)
		}

		sqlStr := parseSqlForMerge(table, queryHolders, insertKeys, insertValues, updateValues, conflictKeys)
		r, err := d.DoExec(ctx, link, sqlStr, queryValues...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected += n
		}
	}
	return batchResult, nil
}

// formatMergeUpdateValues returns the assignments of the MERGE UPDATE SET clause for a record
// with the given `keys`. It follows OnDuplicate/OnDuplicateEx of `option` if specified, or else
// updates every key except conflict keys and, unless the whole record is replaced, soft created
// fields. A conflict key assigned from its own source column is left out, as it cannot change
// the matched row.
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
	// Filter conflict keys, and soft created fields unless the whole record is replaced.
	for _, key := range keys {
		if conflictKeySet.Contains(gstr.ToUpper(key)) {
			continue
		}
		if option.InsertOption != gdb.InsertOptionReplace && d.Core.IsSoftCreatedFieldName(key) {
			continue
		}
		keyWithChar := charL + key + charR
		updateValues = append(updateValues, fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, keyWithChar))
	}
	return updateValues
}

// parseSqlForMerge generates MERGE statement for Oracle database.
// When updateValues is empty, it only inserts (INSERT IGNORE behavior).
// When updateValues is provided, it performs upsert (INSERT or UPDATE).
// Examples:
// - INSERT IGNORE: MERGE INTO table T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...)
// - UPSERT: MERGE INTO table T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...) WHEN MATCHED THEN UPDATE SET ...
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
		`MERGE INTO %s T1 USING (SELECT %s FROM DUAL) T2 ON (%s) WHEN ` +
			`NOT MATCHED THEN INSERT(%s) VALUES (%s)`,
	)
	if len(updateValues) > 0 {
		// Upsert: INSERT or UPDATE
		pattern += gstr.Trim(` WHEN MATCHED THEN UPDATE SET %s`)
		return fmt.Sprintf(
			pattern, table, queryHolderStr, duplicateKeyStr, insertKeyStr, insertValueStr,
			strings.Join(updateValues, ","),
		)
	}
	// Insert Ignore: INSERT only
	return fmt.Sprintf(pattern, table, queryHolderStr, duplicateKeyStr, insertKeyStr, insertValueStr)
}

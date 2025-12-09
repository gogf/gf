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
	"strings"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
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

	default:
	}
	var (
		keys   []string
		values []string
		params []any
	)
	// Retrieve the table fields and length.
	var (
		listLength  = len(list)
		valueHolder = make([]string, 0)
	)
	for k := range list[0] {
		keys = append(keys, k)
		valueHolder = append(valueHolder, "?")
	}
	var (
		batchResult    = new(gdb.SqlResult)
		charL, charR   = d.GetChars()
		keyStr         = charL + strings.Join(keys, charL+","+charR) + charR
		valueHolderStr = strings.Join(valueHolder, ",")
	)
	// Format "INSERT...INTO..." statement.
	intoStrArray := make([]string, 0)
	for i := 0; i < len(list); i++ {
		for _, k := range keys {
			if s, ok := list[i][k].(gdb.Raw); ok {
				params = append(params, gconv.String(s))
			} else {
				params = append(params, list[i][k])
			}
		}
		values = append(values, valueHolderStr)
		intoStrArray = append(
			intoStrArray,
			fmt.Sprintf(
				"INTO %s(%s) VALUES(%s)",
				table, keyStr, valueHolderStr,
			),
		)
		if len(intoStrArray) == option.BatchCount || (i == listLength-1 && len(valueHolder) > 0) {
			r, err := d.DoExec(ctx, link, fmt.Sprintf(
				"INSERT ALL %s SELECT * FROM DUAL",
				strings.Join(intoStrArray, " "),
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
			params = params[:0]
			intoStrArray = intoStrArray[:0]
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
		conflictKeys = primaryKeys
	}

	var (
		one            = list[0]
		oneLen         = len(one)
		charL, charR   = d.GetChars()
		conflictKeySet = gset.New(false)

		// queryHolders:	Handle data with Holder that need to be upsert
		// queryValues:		Handle data that need to be upsert
		// insertKeys:		Handle valid keys that need to be inserted
		// insertValues:	Handle values that need to be inserted
		// updateValues:	Handle values that need to be updated
		queryHolders = make([]string, oneLen)
		queryValues  = make([]any, oneLen)
		insertKeys   = make([]string, oneLen)
		insertValues = make([]string, oneLen)
		updateValues []string
	)

	// conflictKeys slice type conv to set type
	for _, conflictKey := range conflictKeys {
		conflictKeySet.Add(gstr.ToUpper(conflictKey))
	}

	index := 0
	for key, value := range one {
		keyWithChar := charL + key + charR
		queryHolders[index] = fmt.Sprintf("? AS %s", keyWithChar)
		queryValues[index] = value
		insertKeys[index] = keyWithChar
		insertValues[index] = fmt.Sprintf("T2.%s", keyWithChar)

		// Build updateValues only when withUpdate is true
		// Filter conflict keys and soft created fields from updateValues
		if withUpdate && !(conflictKeySet.Contains(key) || d.Core.IsSoftCreatedFieldName(key)) {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, keyWithChar),
			)
		}
		index++
	}

	var (
		batchResult = new(gdb.SqlResult)
		sqlStr      = parseSqlForMerge(table, queryHolders, insertKeys, insertValues, updateValues, conflictKeys)
	)
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
	return batchResult, nil
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
	pattern := gstr.Trim(`MERGE INTO %s T1 USING (SELECT %s FROM DUAL) T2 ON (%s) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s)`)
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

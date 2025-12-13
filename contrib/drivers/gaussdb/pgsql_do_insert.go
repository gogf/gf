// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb

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
	ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		// Treat Replace as Save operation
		return d.doSave(ctx, link, table, list, option)

	// GaussDB does not support InsertIgnore with ON CONFLICT, use MERGE instead
	case gdb.InsertOptionIgnore:
		return d.doInsertIgnore(ctx, link, table, list, option)

	case gdb.InsertOptionDefault:
		// Get table fields to retrieve the primary key TableField object (not just the name)
		// because DoExec needs the `TableField.Type` to determine if LastInsertId is supported.
		tableFields, err := d.GetCore().GetDB().TableFields(ctx, table)
		if err == nil {
			for _, field := range tableFields {
				if strings.EqualFold(field.Key, "pri") {
					pkField := *field
					ctx = context.WithValue(ctx, internalPrimaryKeyInCtx, pkField)
					break
				}
			}
		}

	default:
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}

// doSave implements upsert operation using MERGE statement for GaussDB.
func (d *Driver) doSave(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, true)
}

// doInsertIgnore implements INSERT IGNORE operation using MERGE statement for GaussDB.
// It only inserts records when there's no conflict on primary/unique keys.
func (d *Driver) doInsertIgnore(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, false)
}

// doUpdateThenInsert handles upsert when conflict keys need to be updated.
// GaussDB MERGE cannot update columns in ON clause, so we use UPDATE + INSERT instead.
func (d *Driver) doUpdateThenInsert(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	charL, charR := d.GetChars()
	var (
		batchResult   = new(gdb.SqlResult)
		totalAffected int64
	)

	for _, data := range list {
		// Build UPDATE statement
		var (
			updateFields []string
			updateValues []any
			whereFields  []string
			whereValues  []any
			valueIndex   = 1
		)

		// Process OnDuplicateMap to build UPDATE SET clause
		for updateKey, updateValue := range option.OnDuplicateMap {
			keyWithChar := charL + updateKey + charR
			switch v := updateValue.(type) {
			case gdb.Raw, *gdb.Raw:
				rawStr := fmt.Sprintf("%v", v)
				rawStr = strings.ReplaceAll(rawStr, "EXCLUDED.", "")
				rawStr = strings.ReplaceAll(rawStr, "EXCLUDED ", "")
				updateFields = append(updateFields, fmt.Sprintf("%s = %s", keyWithChar, rawStr))
			case gdb.Counter, *gdb.Counter:
				var counter gdb.Counter
				if c, ok := v.(gdb.Counter); ok {
					counter = c
				} else if c, ok := v.(*gdb.Counter); ok {
					counter = *c
				}
				operator := "+"
				columnVal := counter.Value
				if columnVal < 0 {
					operator = "-"
					columnVal = -columnVal
				}
				fieldWithChar := charL + counter.Field + charR
				// For UPDATE statement, use the data value instead of referencing another column
				if dataValue, ok := data[counter.Field]; ok {
					updateFields = append(updateFields, fmt.Sprintf("%s = $%d %s %v", keyWithChar, valueIndex, operator, columnVal))
					updateValues = append(updateValues, dataValue)
					valueIndex++
				} else {
					updateFields = append(updateFields, fmt.Sprintf("%s = %s %s %v", keyWithChar, fieldWithChar, operator, columnVal))
				}
			default:
				// Map value to another field name or use the value from data
				valueStr := gconv.String(updateValue)
				if dataValue, ok := data[valueStr]; ok {
					updateFields = append(updateFields, fmt.Sprintf("%s = $%d", keyWithChar, valueIndex))
					updateValues = append(updateValues, dataValue)
					valueIndex++
				} else {
					updateFields = append(updateFields, fmt.Sprintf("%s = $%d", keyWithChar, valueIndex))
					updateValues = append(updateValues, updateValue)
					valueIndex++
				}
			}
		}

		// Build WHERE clause using OnConflict keys
		for _, conflictKey := range option.OnConflict {
			if dataValue, ok := data[conflictKey]; ok {
				keyWithChar := charL + conflictKey + charR
				whereFields = append(whereFields, fmt.Sprintf("%s = $%d", keyWithChar, valueIndex))
				whereValues = append(whereValues, dataValue)
				valueIndex++
			}
		}

		if len(updateFields) > 0 && len(whereFields) > 0 {
			updateSQL := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
				table,
				strings.Join(updateFields, ", "),
				strings.Join(whereFields, " AND "),
			)
			updateResult, updateErr := d.DoExec(ctx, link, updateSQL, append(updateValues, whereValues...)...)
			if updateErr != nil {
				return nil, updateErr
			}

			affected, _ := updateResult.RowsAffected()
			if affected > 0 {
				// UPDATE successful
				totalAffected += affected
				continue
			}
		}

		// If UPDATE affected 0 rows, do INSERT
		var (
			insertKeys    []string
			insertHolders []string
			insertValues  []any
			insertIndex   = 1
		)
		for key, value := range data {
			keyWithChar := charL + key + charR
			insertKeys = append(insertKeys, keyWithChar)
			insertHolders = append(insertHolders, fmt.Sprintf("$%d", insertIndex))
			insertValues = append(insertValues, value)
			insertIndex++
		}

		insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			table,
			strings.Join(insertKeys, ", "),
			strings.Join(insertHolders, ", "),
		)
		insertResult, insertErr := d.DoExec(ctx, link, insertSQL, insertValues...)
		if insertErr != nil {
			// Ignore duplicate key errors (race condition: another transaction inserted between our UPDATE and INSERT)
			if strings.Contains(insertErr.Error(), "duplicate key") ||
				strings.Contains(insertErr.Error(), "unique constraint") {
				continue
			}
			return nil, insertErr
		}

		affected, _ := insertResult.RowsAffected()
		totalAffected += affected
	}

	batchResult.Result = &gdb.SqlResult{}
	batchResult.Affected = totalAffected
	return batchResult, nil
}

// doMergeInsert implements MERGE-based insert operations for GaussDB.
// When withUpdate is true, it performs upsert (insert or update).
// When withUpdate is false, it performs insert ignore (insert only when no conflict).
func (d *Driver) doMergeInsert(
	ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption, withUpdate bool,
) (result sql.Result, err error) {
	// Check if OnDuplicateMap contains conflict keys
	// GaussDB MERGE statement cannot update columns used in ON clause
	// If user wants to update conflict keys, we need to use a different approach
	if withUpdate && len(option.OnDuplicateMap) > 0 && len(option.OnConflict) > 0 {
		conflictKeySet := gset.NewStrSetFrom(option.OnConflict)
		hasConflictKeyUpdate := false
		for updateKey := range option.OnDuplicateMap {
			if conflictKeySet.Contains(strings.ToLower(updateKey)) ||
				conflictKeySet.Contains(strings.ToUpper(updateKey)) ||
				conflictKeySet.Contains(updateKey) {
				hasConflictKeyUpdate = true
				break
			}
		}
		if hasConflictKeyUpdate {
			// Use UPDATE + INSERT approach when conflict keys need to be updated
			return d.doUpdateThenInsert(ctx, link, table, list, option)
		}
	}

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
			// For InsertIgnore without primary key, try normal insert and ignore duplicate errors
			// For Save/Replace, primary key is required
			if !withUpdate {
				result, err := d.Core.DoInsert(ctx, link, table, list, option)
				if err != nil {
					// Ignore duplicate key errors for InsertIgnore
					if strings.Contains(err.Error(), "duplicate key") ||
						strings.Contains(err.Error(), "unique constraint") {
						return result, nil
					}
					return result, err
				}
				return result, nil
			}
			return nil, gerror.NewCodef(
				gcode.CodeMissingParameter,
				`Replace/Save operation requires conflict detection: `+
					`either specify OnConflict() columns or ensure table '%s' has a primary key in the data`,
				table,
			)
		}
		// TODO consider composite primary keys.
		conflictKeys = primaryKeys
	}

	var (
		one            = list[0]
		oneLen         = len(one)
		charL, charR   = d.GetChars()
		conflictKeySet = gset.New(false)

		// queryHolders:	Handle data with Holder that need to be merged
		// queryValues:		Handle data that need to be merged
		// insertKeys:		Handle valid keys that need to be inserted
		// insertValues:	Handle values that need to be inserted
		// updateValues:	Handle values that need to be updated (only when withUpdate=true)
		queryHolders = make([]string, oneLen)
		queryValues  = make([]any, oneLen)
		insertKeys   = make([]string, oneLen)
		insertValues = make([]string, oneLen)
		updateValues []string
	)

	// conflictKeys slice type conv to set type
	for _, conflictKey := range conflictKeys {
		conflictKeySet.Add(strings.ToUpper(conflictKey))
	}

	index := 0
	for key, value := range one {
		keyWithChar := charL + key + charR
		queryHolders[index] = fmt.Sprintf("$%d AS %s", index+1, keyWithChar)
		queryValues[index] = value
		insertKeys[index] = keyWithChar
		insertValues[index] = fmt.Sprintf("T2.%s", keyWithChar)
		index++
	}

	// Build updateValues only when withUpdate is true
	if withUpdate {
		// Check if OnDuplicateStr or OnDuplicateMap is specified for custom update logic
		if option.OnDuplicateStr != "" {
			// Parse OnDuplicateStr (e.g., "field1,field2" or "field1, field2")
			fields := gstr.SplitAndTrim(option.OnDuplicateStr, ",")
			for _, field := range fields {
				fieldWithChar := charL + field + charR
				updateValues = append(
					updateValues,
					fmt.Sprintf(`T1.%s = T2.%s`, fieldWithChar, fieldWithChar),
				)
			}
		} else if len(option.OnDuplicateMap) > 0 {
			// Use OnDuplicateMap for custom update mapping
			for updateKey, updateValue := range option.OnDuplicateMap {
				// Skip conflict keys - they cannot be updated in MERGE
				if conflictKeySet.Contains(strings.ToUpper(updateKey)) {
					continue
				}
				keyWithChar := charL + updateKey + charR
				switch v := updateValue.(type) {
				case gdb.Raw, *gdb.Raw:
					// Raw SQL expression
					// Replace EXCLUDED (PostgreSQL ON CONFLICT syntax) with T2 (MERGE syntax)
					rawStr := fmt.Sprintf("%v", v)
					rawStr = strings.ReplaceAll(rawStr, "EXCLUDED.", "T2.")
					rawStr = strings.ReplaceAll(rawStr, "EXCLUDED ", "T2 ")
					updateValues = append(
						updateValues,
						fmt.Sprintf(`T1.%s = %s`, keyWithChar, rawStr),
					)
				case gdb.Counter, *gdb.Counter:
					// Counter operation
					var counter gdb.Counter
					if c, ok := v.(gdb.Counter); ok {
						counter = c
					} else if c, ok := v.(*gdb.Counter); ok {
						counter = *c
					}
					operator := "+"
					columnVal := counter.Value
					if columnVal < 0 {
						operator = "-"
						columnVal = -columnVal
					}
					fieldWithChar := charL + counter.Field + charR
					updateValues = append(
						updateValues,
						fmt.Sprintf(`T1.%s = T2.%s %s %v`, keyWithChar, fieldWithChar, operator, columnVal),
					)
				default:
					// Map value to another field name
					valueStr := gconv.String(updateValue)
					valueWithChar := charL + valueStr + charR
					updateValues = append(
						updateValues,
						fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, valueWithChar),
					)
				}
			}
		} else {
			// Default: update all fields except conflict keys and soft created fields
			for key := range one {
				if conflictKeySet.Contains(strings.ToUpper(key)) || d.Core.IsSoftCreatedFieldName(key) {
					continue
				}
				keyWithChar := charL + key + charR
				updateValues = append(
					updateValues,
					fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, keyWithChar),
				)
			}
		}
	}

	var (
		batchResult = new(gdb.SqlResult)
		sqlStr      string
	)

	// For InsertIgnore (withUpdate=false), we need to check if record exists first
	if !withUpdate {
		// Build WHERE clause to check if record exists
		var whereConditions []string
		var checkValues []any
		checkIndex := 1
		for _, key := range conflictKeys {
			if value, ok := one[key]; ok {
				keyWithChar := charL + key + charR
				whereConditions = append(whereConditions, fmt.Sprintf("%s = $%d", keyWithChar, checkIndex))
				checkValues = append(checkValues, value)
				checkIndex++
			}
		}
		whereClause := strings.Join(whereConditions, " AND ")

		// Check if record exists
		checkSQL := fmt.Sprintf("SELECT 1 FROM %s WHERE %s LIMIT 1", table, whereClause)
		checkResult, checkErr := d.DoQuery(ctx, link, checkSQL, checkValues...)
		if checkErr != nil {
			return nil, checkErr
		}

		// If record exists, return result with 0 affected rows
		if len(checkResult) > 0 {
			batchResult.Result = &gdb.SqlResult{}
			batchResult.Affected = 0
			return batchResult, nil
		}

		// Record doesn't exist, proceed with insert
		// For InsertIgnore, we just do a simple INSERT (no MERGE needed since we checked it doesn't exist)
		var insertSQL strings.Builder
		insertSQL.WriteString(fmt.Sprintf("INSERT INTO %s (", table))
		insertSQL.WriteString(strings.Join(insertKeys, ","))
		insertSQL.WriteString(") VALUES (")
		for i := range insertKeys {
			if i > 0 {
				insertSQL.WriteString(",")
			}
			insertSQL.WriteString(fmt.Sprintf("$%d", i+1))
		}
		insertSQL.WriteString(")")

		r, err := d.DoExec(ctx, link, insertSQL.String(), queryValues...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected = n
		}
		return batchResult, nil
	}

	// For Save/Replace (withUpdate=true), use MERGE
	sqlStr = parseSqlForMerge(table, queryHolders, insertKeys, insertValues, updateValues, conflictKeys, charL, charR)
	r, err := d.DoExec(ctx, link, sqlStr, queryValues...)
	if err != nil {
		return r, err
	}
	// GaussDB's MERGE statement may not return correct RowsAffected
	// We manually set it to 1 since MERGE always affects exactly one row
	if n, err := r.RowsAffected(); err != nil {
		return r, err
	} else {
		batchResult.Result = r
		// If RowsAffected returns 0, manually set to 1 for MERGE operations
		if n == 0 {
			batchResult.Affected = 1
		} else {
			batchResult.Affected += n
		}
	}
	return batchResult, nil
}

// parseSqlForMerge generates MERGE statement for GaussDB.
// When updateValues is empty, it only inserts (INSERT IGNORE behavior).
// When updateValues is provided, it performs upsert (INSERT or UPDATE).
// Examples:
// - INSERT IGNORE: MERGE INTO table T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...)
// - UPSERT: MERGE INTO table T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...) WHEN MATCHED THEN UPDATE SET ...
func parseSqlForMerge(table string,
	queryHolders, insertKeys, insertValues, updateValues, duplicateKey []string, charL, charR string,
) (sqlStr string) {
	var (
		intoStr   = fmt.Sprintf("MERGE INTO %s AS T1", table)
		usingStr  = fmt.Sprintf("USING (SELECT %s) AS T2", strings.Join(queryHolders, ","))
		onStr     string
		insertStr = fmt.Sprintf(
			"WHEN NOT MATCHED THEN INSERT (%s) VALUES (%s)",
			strings.Join(insertKeys, ","),
			strings.Join(insertValues, ","),
		)
		updateStr string
	)

	// Build ON condition
	var onConditions []string
	for _, key := range duplicateKey {
		keyWithChar := charL + key + charR
		onConditions = append(onConditions, fmt.Sprintf("T1.%s = T2.%s", keyWithChar, keyWithChar))
	}
	onStr = "ON (" + strings.Join(onConditions, " AND ") + ")"

	// Build UPDATE clause only when updateValues is provided
	if len(updateValues) > 0 {
		updateStr = fmt.Sprintf(" WHEN MATCHED THEN UPDATE SET %s", strings.Join(updateValues, ","))
	}

	sqlStr = fmt.Sprintf("%s %s %s %s%s", intoStr, usingStr, onStr, insertStr, updateStr)
	return
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	createdFiledNames = []string{"created_at", "create_at"}
)

const (
	autoIncrementName    = "auto_increment"
	mssqlOutPutKey       = "OUTPUT"
	mssqlInsertedObjName = "INSERTED"
	mssqlAffectFd        = " 1 as AffectCount"
	affectCountFieldName = "AffectCount"
	mssqlPrimaryKeyName  = "PRI"
	fdId                 = "ID"
)

// DoInsert inserts or updates data for given table. rewrite db.core.DoInsert
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by mssql driver`,
		)
		//default:
		//	return d.Core.DoInsert(ctx, link, table, list, option)
	}
	var (
		keys           []string      // Field names.
		values         []string      // Value holder string array, like: (?,?,?)
		params         []interface{} // Values that will be committed to underlying database driver.
		onDuplicateStr string        // onDuplicateStr is used in "ON DUPLICATE KEY UPDATE" statement.
	)
	// Group the list by fields. Different fields to different list.
	// It here uses ListMap to keep sequence for data inserting.
	var keyListMap = gmap.NewListMap()
	for _, item := range list {
		var (
			tmpKeys              = make([]string, 0)
			tmpKeysInSequenceStr string
		)
		for k := range item {
			tmpKeys = append(tmpKeys, k)
		}
		keys, err = d.fieldsToSequence(ctx, table, tmpKeys)
		if err != nil {
			return nil, err
		}
		tmpKeysInSequenceStr = gstr.Join(keys, ",")

		if !keyListMap.Contains(tmpKeysInSequenceStr) {
			keyListMap.Set(tmpKeysInSequenceStr, make(gdb.List, 0))
		}
		tmpKeysInSequenceList := keyListMap.Get(tmpKeysInSequenceStr).(gdb.List)
		tmpKeysInSequenceList = append(tmpKeysInSequenceList, item)
		keyListMap.Set(tmpKeysInSequenceStr, tmpKeysInSequenceList)
	}
	if keyListMap.Size() > 1 {
		var (
			tmpResult    sql.Result
			sqlResult    gdb.SqlResult
			rowsAffected int64
		)
		keyListMap.Iterator(func(key, value interface{}) bool {
			tmpResult, err = d.DoInsert(ctx, link, table, value.(gdb.List), option)
			if err != nil {
				return false
			}
			rowsAffected, err = tmpResult.RowsAffected()
			if err != nil {
				return false
			}
			sqlResult.Result = tmpResult
			sqlResult.Affected += rowsAffected
			return true
		})
		return &sqlResult, nil
	}
	// Prepare the batch result pointer.
	var (
		charL, charR = d.GetDB().GetChars()
		batchResult  = new(gdb.SqlResult)
		keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
		operation    = gdb.GetInsertOperationByOption(option.InsertOption)
	)
	if option.InsertOption == gdb.InsertOptionSave {
		onDuplicateStr = d.formatOnDuplicate(keys, option)
	}
	var (
		listLength  = len(list)
		valueHolder = make([]string, 0)
	)
	for i := 0; i < listLength; i++ {
		values = values[:0]
		// Note that the map type is unordered,
		// so it should use slice+key to retrieve the value.
		for _, k := range keys {
			if s, ok := list[i][k].(gdb.Raw); ok {
				values = append(values, gconv.String(s))
			} else {
				values = append(values, "?")
				params = append(params, list[i][k])
			}
		}
		valueHolder = append(valueHolder, "("+gstr.Join(values, ",")+")")
		// Batch package checks: It meets the batch number, or it is the last element.
		if len(valueHolder) == option.BatchCount || (i == listLength-1 && len(valueHolder) > 0) {
			var (
				//stdSqlResult sql.Result
				stdSqlResult gdb.Result
				//affectedRows int64
				retResult interface{}
			)

			stdSqlResult, err = d.GetDB().DoQuery(ctx, link, fmt.Sprintf(
				"%s INTO %s(%s) %s VALUES%s %s ",
				operation, d.QuotePrefixTableName(table), keysStr,
				d.GetInsertOutputSql(ctx, table),
				gstr.Join(valueHolder, ","),
				onDuplicateStr,
			), params...)
			if err != nil {
				retResult = &InsertResult{lastInsertId: 0, rowsAffected: 0, err: err}
				return retResult.(sql.Result), err
			}
			var (
				aCount int64 // affect count
				lId    int64 // last insert id
			)
			if len(stdSqlResult) == 0 {
				err = gerror.WrapCode(gcode.CodeDbOperationError, gerror.New("affectcount is zero"), `sql.Result.RowsAffected failed`)
				retResult = &InsertResult{lastInsertId: 0, rowsAffected: 0, err: err}
				return retResult.(sql.Result), err
			}
			// get affect count
			aCount = stdSqlResult[0].GMap().GetVar(affectCountFieldName).Int64()
			// get last_insert_id
			lId = stdSqlResult[0].GMap().GetVar(fdId).Int64()

			retResult = &InsertResult{lastInsertId: lId, rowsAffected: aCount}

			batchResult.Result = retResult.(sql.Result)
			batchResult.Affected += aCount

			params = params[:0]
			valueHolder = valueHolder[:0]
		}
	}
	return batchResult, nil
}

// InsertResult instance of sql.Result
type InsertResult struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (r *InsertResult) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

func (r *InsertResult) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

// GetInsertOutputSql  gen get last_insert_id code
func (m *Driver) GetInsertOutputSql(ctx context.Context, table string) string {
	fds, errFd := m.GetDB().TableFields(ctx, table)
	if errFd != nil {
		return ""
	}
	extraSqlAry := make([]string, 0)
	extraSqlAry = append(extraSqlAry, fmt.Sprintf("%s %s", mssqlOutPutKey, mssqlAffectFd))
	incrNo := 0
	if len(fds) > 0 {
		for _, fd := range fds {
			// has primary key and is auto-incement
			if fd.Extra == autoIncrementName && fd.Key == mssqlPrimaryKeyName && !fd.Null {
				incrNoStr := ""
				if incrNo == 0 { //fixed first field named id, convenient to get
					incrNoStr = fmt.Sprintf(" as %s", fdId)
				}

				extraSqlAry = append(extraSqlAry, fmt.Sprintf("%s.%s%s", mssqlInsertedObjName, fd.Name, incrNoStr))
				incrNo++
			}
			//fmt.Printf("null:%t name:%s key:%s k:%s \n", fd.Null, fd.Name, fd.Key, k)
		}
	}
	//fmt.Println(extraSqlAry)
	return strings.Join(extraSqlAry, ",")
	//";select ID = convert(bigint, SCOPE_IDENTITY()), AffectCount = @@ROWCOUNT;"
}

func (d *Driver) fieldsToSequence(ctx context.Context, table string, fields []string) ([]string, error) {
	var (
		fieldSet               = gset.NewStrSetFrom(fields)
		fieldsResultInSequence = make([]string, 0)
		tableFields, err       = d.GetDB().TableFields(ctx, table)
	)
	if err != nil {
		return nil, err
	}
	// Sort the fields in order.
	var fieldsOfTableInSequence = make([]string, len(tableFields))
	for _, field := range tableFields {
		fieldsOfTableInSequence[field.Index] = field.Name
	}
	// Sort the input fields.
	for _, fieldName := range fieldsOfTableInSequence {
		if fieldSet.Contains(fieldName) {
			fieldsResultInSequence = append(fieldsResultInSequence, fieldName)
		}
	}
	return fieldsResultInSequence, nil
}

func (d *Driver) formatOnDuplicate(columns []string, option gdb.DoInsertOption) string {
	var onDuplicateStr string
	if option.OnDuplicateStr != "" {
		onDuplicateStr = option.OnDuplicateStr
	} else if len(option.OnDuplicateMap) > 0 {
		for k, v := range option.OnDuplicateMap {
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			switch v.(type) {
			case gdb.Raw, *gdb.Raw:
				onDuplicateStr += fmt.Sprintf(
					"%s=%s",
					d.QuoteWord(k),
					v,
				)
			default:
				onDuplicateStr += fmt.Sprintf(
					"%s=VALUES(%s)",
					d.QuoteWord(k),
					d.QuoteWord(gconv.String(v)),
				)
			}
		}
	} else {
		for _, column := range columns {
			// If it's SAVE operation, do not automatically update the creating time.
			if d.isSoftCreatedFieldName(column) {
				continue
			}
			if len(onDuplicateStr) > 0 {
				onDuplicateStr += ","
			}
			onDuplicateStr += fmt.Sprintf(
				"%s=VALUES(%s)",
				d.QuoteWord(column),
				d.QuoteWord(column),
			)
		}
	}
	return fmt.Sprintf("ON DUPLICATE KEY UPDATE %s", onDuplicateStr)
}

func (d *Driver) isSoftCreatedFieldName(fieldName string) bool {
	if fieldName == "" {
		return false
	}
	if config := d.GetDB().GetConfig(); config.CreatedAt != "" {
		if equalFoldWithoutChars(fieldName, config.CreatedAt) {
			return true
		}
		return gstr.InArray(append([]string{config.CreatedAt}, createdFiledNames...), fieldName)
	}
	for _, v := range createdFiledNames {
		if equalFoldWithoutChars(fieldName, v) {
			return true
		}
	}
	return false
}

// doSave support upsert for SQL server
func (d *Driver) doSave(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	if len(option.OnConflict) == 0 {
		return nil, gerror.NewCode(
			gcode.CodeMissingParameter, `Please specify conflict columns`,
		)
	}

	if len(list) == 0 {
		return nil, gerror.NewCode(
			gcode.CodeInvalidRequest, `Save operation list is empty by mssql driver`,
		)
	}

	var (
		one          = list[0]
		oneLen       = len(one)
		charL, charR = d.GetChars()

		conflictKeys   = option.OnConflict
		conflictKeySet = gset.New(false)

		// queryHolders:	Handle data with Holder that need to be upsert
		// queryValues:		Handle data that need to be upsert
		// insertKeys:		Handle valid keys that need to be inserted
		// insertValues:	Handle values that need to be inserted
		// updateValues:	Handle values that need to be updated
		queryHolders = make([]string, oneLen)
		queryValues  = make([]interface{}, oneLen)
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
		queryHolders[index] = "?"
		queryValues[index] = value
		insertKeys[index] = charL + key + charR
		insertValues[index] = "T2." + charL + key + charR

		// filter conflict keys in updateValues.
		// And the key is not a soft created field.
		if !(conflictKeySet.Contains(key) || d.Core.IsSoftCreatedFieldName(key)) {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, charL+key+charR, charL+key+charR),
			)
		}
		index++
	}

	batchResult := new(gdb.SqlResult)
	sqlStr := parseSqlForUpsert(table, queryHolders, insertKeys, insertValues, updateValues, conflictKeys)
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

// parseSqlForUpsert
// MERGE INTO {{table}} T1
// USING ( VALUES( {{queryHolders}}) T2 ({{insertKeyStr}})
// ON (T1.{{duplicateKey}} = T2.{{duplicateKey}} AND ...)
// WHEN NOT MATCHED THEN
// INSERT {{insertKeys}} VALUES {{insertValues}}
// WHEN MATCHED THEN
// UPDATE SET {{updateValues}}
func parseSqlForUpsert(table string,
	queryHolders, insertKeys, insertValues, updateValues, duplicateKey []string,
) (sqlStr string) {
	var (
		queryHolderStr  = strings.Join(queryHolders, ",")
		insertKeyStr    = strings.Join(insertKeys, ",")
		insertValueStr  = strings.Join(insertValues, ",")
		updateValueStr  = strings.Join(updateValues, ",")
		duplicateKeyStr string
		pattern         = gstr.Trim(`MERGE INTO %s T1 USING (VALUES(%s)) T2 (%s) ON (%s) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s) WHEN MATCHED THEN UPDATE SET %s;`)
	)

	for index, keys := range duplicateKey {
		if index != 0 {
			duplicateKeyStr += " AND "
		}
		duplicateTmp := fmt.Sprintf("T1.%s = T2.%s", keys, keys)
		duplicateKeyStr += duplicateTmp
	}

	return fmt.Sprintf(pattern,
		table,
		queryHolderStr,
		insertKeyStr,
		duplicateKeyStr,
		insertKeyStr,
		insertValueStr,
		updateValueStr,
	)
}

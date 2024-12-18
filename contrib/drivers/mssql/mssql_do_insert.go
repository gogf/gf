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
	"strings"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
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
func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption, ext ...interface{}) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by mssql driver`,
		)
	default:
		outPutStr := d.GetInsertOutputSql(ctx, table)
		var insertHandler gdb.InsertHandler
		insertHandler = func(db gdb.DB, ctx context.Context, link gdb.Link, sqlStr string, args ...interface{}) (sql.Result, error) {
			var (
				stdSqlResult gdb.Result
				retResult    interface{}
			)
			stdSqlResult, err = d.GetDB().DoQuery(ctx, link, sqlStr, args...)
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
			return retResult.(sql.Result), nil
		}
		return d.Core.DoInsert(ctx, link, table, list, option, insertHandler, outPutStr)
	}
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

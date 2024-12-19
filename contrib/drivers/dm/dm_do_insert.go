// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

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

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		// TODO:: Should be Supported
		return nil, gerror.NewCode(
			gcode.CodeNotSupported, `Replace operation is not supported by dm driver`,
		)
	}

	return d.Core.DoInsert(ctx, link, table, list, option)
}

// doSave support upsert for dm
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
			gcode.CodeInvalidRequest, `Save operation list is empty by oracle driver`,
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
		keyWithChar := charL + key + charR
		queryHolders[index] = fmt.Sprintf("? AS %s", keyWithChar)
		queryValues[index] = value
		insertKeys[index] = keyWithChar
		insertValues[index] = fmt.Sprintf("T2.%s", keyWithChar)

		// filter conflict keys in updateValues.
		// And the key is not a soft created field.
		if !(conflictKeySet.Contains(key) || d.Core.IsSoftCreatedFieldName(key)) {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, keyWithChar),
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
// USING ( SELECT {{queryHolders}} FROM DUAL T2
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
		pattern         = gstr.Trim(`MERGE INTO %s T1 USING (SELECT %s FROM DUAL) T2 ON (%s) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s) WHEN MATCHED THEN UPDATE SET %s;`)
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
		duplicateKeyStr,
		insertKeyStr,
		insertValueStr,
		updateValueStr,
	)
}

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
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"strings"
)

// MERGE INTO test USING dual ON ( ID = 1 )
// WHEN MATCHED THEN
// UPDATE
// SET PASSPORT = '1',
// PASSWORD = '2',
// NICKNAME = '3',
// CREATE_TIME = '3',
// SALARY = 4
// WHEN NOT MATCHED THEN
// INSERT (ID, PASSPORT, PASSWORD, NICKNAME, CREATE_TIME, SALARY )
// VALUES
// (1, 'a', 'b', 'c', 'd', 100);
func (d *Driver) doSave(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	if len(option.OnConflict) == 0 {
		return nil, gerror.New("Please specify conflict columns")
	}

	if len(list) == 0 {
		return nil, gerror.NewCode(
			gcode.CodeInvalidRequest, `Save operation list is empty by oracle driver`,
		)
	}

	var (
		charL, charR = d.GetChars()
		//onConflict   = option.OnConflict
		keys                   []string
		insertKeys             []string
		insertValues           []string
		updateValues           []string
		queryValues            []string
		valueCharL, valueCharR = "\"", "\""
		one                    = list[0]
	)
	// Column names need to be aligned in the syntax
	for k := range list[0] {
		keys = append(keys, k)
	}

	for key, value := range one {
		insertKeys = append(insertKeys, charL+key+charR)
		insertValues = append(insertValues, "T2."+charL+key+charR)
		// TODO
		if key != "id" {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, charL+key+charR, charL+key+charR),
			)
		}

		saveValue := gconv.String(value)
		queryValues = append(
			queryValues,
			fmt.Sprintf(
				valueCharL+"%s"+valueCharR+" AS "+charL+"%s"+charR,
				saveValue, key,
			),
		)
	}
	// insertKeys:   Handle valid keys that need to be inserted and updated
	// insertValues: Handle values that need to be inserted
	// updateValues: Handle values that need to be updated
	// queryValues:  Handle only one insert with column name
	//insertKeys, insertValues, updateValues, queryValues := parseValue(list[0])
	//
	// unionValues: Handling values that need to be inserted and updated
	unionValues := parseUnion(list[1:], keys)
	gutil.Dump(unionValues)
	//
	//batchResult := new(gdb.SqlResult)
	//sqlStr := parseSql(
	//	insertKeys, insertValues, updateValues, queryValues, []string{"231321321"}, table, "id",
	//)
	//fmt.Println(sqlStr)
	//
	//fmt.Println(sqlStr)
	//return nil, err
	//r, err := d.DoExec(ctx, link, sqlStr)
	//if err != nil {
	//	return r, err
	//}
	//if n, err := r.RowsAffected(); err != nil {
	//	return r, err
	//} else {
	//	batchResult.Result = r
	//	batchResult.Affected += n
	//}
	//return batchResult, nil
	return nil, err
}

func parseUnion(list gdb.List, keys []string) (unionValues []string) {
	for _, mapper := range list {
		var saveValue []string
		var (
			valueCharL = "\""
			valueCharR = "\""
		)
		for _, column := range keys {
			if mapper[column] == nil {
				continue
			}
			saveValue = append(saveValue,
				fmt.Sprintf(
					valueCharL+"%s"+valueCharR,
					gconv.String(mapper[column]),
				))
		}
		unionValues = append(
			unionValues,
			fmt.Sprintf(`UNION ALL SELECT %s FROM DUAL`, strings.Join(saveValue, ",")),
		)
	}
	return
}

func parseSql(
	insertKeys, insertValues, updateValues, queryValues, unionValues []string, table, duplicateKey string,
) (sqlStr string) {
	var (
		queryValueStr  = strings.Join(queryValues, ",")
		unionValueStr  = strings.Join(unionValues, " ")
		insertKeyStr   = strings.Join(insertKeys, ",")
		insertValueStr = strings.Join(insertValues, ",")
		updateValueStr = strings.Join(updateValues, ",")
		pattern        = gstr.Trim(`
MERGE INTO %s T1 USING (SELECT %s FROM DUAL %s) T2 ON %s 
WHEN NOT MATCHED 
THEN 
INSERT(%s) VALUES (%s) 
WHEN MATCHED 
THEN 
UPDATE SET %s; 
COMMIT;
`)
	)
	return fmt.Sprintf(
		pattern,
		table, queryValueStr, unionValueStr,
		fmt.Sprintf("(T1.%s = T2.%s)", duplicateKey, duplicateKey),
		insertKeyStr, insertValueStr, updateValueStr,
	)
}

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
	_ "gitee.com/chunanyong/dm"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// DoInsert inserts or updates data forF given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		// TODO:: Should be Supported
		return nil, gerror.NewCode(
			gcode.CodeNotSupported, `Replace operation is not supported by dm driver`,
		)

	case gdb.InsertOptionSave:
		// This syntax currently only supports design tables whose primary key is ID.
		listLength := len(list)
		if listLength == 0 {
			return nil, gerror.NewCode(
				gcode.CodeInvalidRequest, `Save operation list is empty by dm driver`,
			)
		}
		var (
			keysSort     []string
			charL, charR = d.GetChars()
		)
		// Column names need to be aligned in the syntax
		for k := range list[0] {
			keysSort = append(keysSort, k)
		}
		var char = struct {
			charL        string
			charR        string
			valueCharL   string
			valueCharR   string
			duplicateKey string
			keys         []string
		}{
			charL:      charL,
			charR:      charR,
			valueCharL: "'",
			valueCharR: "'",
			// TODO:: Need to dynamically set the primary key of the table
			duplicateKey: "ID",
			keys:         keysSort,
		}

		// insertKeys:   Handle valid keys that need to be inserted and updated
		// insertValues: Handle values that need to be inserted
		// updateValues: Handle values that need to be updated
		// queryValues:  Handle only one insert with column name
		insertKeys, insertValues, updateValues, queryValues := parseValue(list[0], char)
		// unionValues: Handling values that need to be inserted and updated
		unionValues := parseUnion(list[1:], char)

		batchResult := new(gdb.SqlResult)
		// parseSql():
		// MERGE INTO {{table}} T1
		// USING ( SELECT {{queryValues}} FROM DUAL
		// {{unionValues}} ) T2
		// ON (T1.{{duplicateKey}} = T2.{{duplicateKey}})
		// WHEN NOT MATCHED THEN
		// INSERT {{insertKeys}} VALUES {{insertValues}}
		// WHEN MATCHED THEN
		// UPDATE SET {{updateValues}}
		sqlStr := parseSql(
			insertKeys, insertValues, updateValues, queryValues, unionValues, table, char.duplicateKey,
		)
		r, err := d.DoExec(ctx, link, sqlStr)
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
	return d.Core.DoInsert(ctx, link, table, list, option)
}

func parseValue(listOne gdb.Map, char struct {
	charL        string
	charR        string
	valueCharL   string
	valueCharR   string
	duplicateKey string
	keys         []string
}) (insertKeys []string, insertValues []string, updateValues []string, queryValues []string) {
	for _, column := range char.keys {
		if listOne[column] == nil {
			// remove unassigned struct object
			continue
		}
		insertKeys = append(insertKeys, char.charL+column+char.charR)
		insertValues = append(insertValues, "T2."+char.charL+column+char.charR)
		if column != char.duplicateKey {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, char.charL+column+char.charR, char.charL+column+char.charR),
			)
		}

		saveValue := gconv.String(listOne[column])
		queryValues = append(
			queryValues,
			fmt.Sprintf(
				char.valueCharL+"%s"+char.valueCharR+" AS "+char.charL+"%s"+char.charR,
				saveValue, column,
			),
		)
	}
	return
}

func parseUnion(list gdb.List, char struct {
	charL        string
	charR        string
	valueCharL   string
	valueCharR   string
	duplicateKey string
	keys         []string
}) (unionValues []string) {
	for _, mapper := range list {
		var saveValue []string
		for _, column := range char.keys {
			if mapper[column] == nil {
				continue
			}
			// va := reflect.ValueOf(mapper[column])
			// ty := reflect.TypeOf(mapper[column])
			// switch ty.Kind() {
			// case reflect.String:
			// 	saveValue = append(saveValue, char.valueCharL+va.String()+char.valueCharR)

			// case reflect.Int:
			// 	saveValue = append(saveValue, strconv.FormatInt(va.Int(), 10))

			// case reflect.Int64:
			// 	saveValue = append(saveValue, strconv.FormatInt(va.Int(), 10))

			// default:
			// 	// The fish has no chance getting here.
			// 	// Nothing to do.
			// }
			saveValue = append(saveValue,
				fmt.Sprintf(
					char.valueCharL+"%s"+char.valueCharR,
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

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/mattn/go-oci8"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

// Package oracle implements gdb.Driver, which supports operations for Oracle.
package oracle

import (
	_ "github.com/mattn/go-oci8"

	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Driver is the driver for oracle database.
type Driver struct {
	*gdb.Core
}

var (
	// tableFieldsMap caches the table information retrieved from database.
	tableFieldsMap = gmap.New(true)
)

func init() {
	if err := gdb.Register(`oracle`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for Oracle.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for oracle.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

// Open creates and returns an underlying sql.DB object for oracle.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "oci8"
	)
	if config.Link != "" {
		source = config.Link
	} else {
		source = fmt.Sprintf(
			"%s/%s@%s:%s/%s",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}

	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

// FilteredLink retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *Driver) FilteredLink() string {
	linkInfo := d.GetConfig().Link
	if linkInfo == "" {
		return ""
	}
	s, _ := gregex.ReplaceString(
		`(.+?)\s*/\s*(.+)\s*@\s*(.+)\s*:\s*(\d+)\s*/\s*(.+)`,
		`$1/xxx@$3:$4/$5`,
		linkInfo,
	)
	return s
}

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	defer func() {
		newSql, newArgs, err = d.Core.DoFilter(ctx, link, newSql, newArgs)
	}()

	var index int
	// Convert placeholder char '?' to string ":vx".
	newSql, _ = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf(":v%d", index)
	})
	newSql, _ = gregex.ReplaceString("\"", "", newSql)
	// Handle string datetime argument.
	for i, v := range args {
		if reflect.TypeOf(v).Kind() == reflect.String {
			valueStr := gconv.String(v)
			if gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, valueStr) {
				// args[i] = fmt.Sprintf(`TO_DATE('%s','yyyy-MM-dd HH:MI:SS')`, valueStr)
				args[i], _ = time.ParseInLocation("2006-01-02 15:04:05", valueStr, time.Local)
			}
		}
	}
	newSql = d.parseSql(newSql)
	newArgs = args
	return
}

// parseSql does some replacement of the sql before commits it to underlying driver,
// for support of oracle server.
func (d *Driver) parseSql(sql string) string {
	var (
		patten      = `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,{0,1}\s*(\d*))`
		allMatch, _ = gregex.MatchAllString(patten, sql)
	)
	if len(allMatch) == 0 {
		return sql
	}
	var (
		index   = 0
		keyword = strings.ToUpper(strings.TrimSpace(allMatch[index][0]))
	)
	index++
	switch keyword {
	case "SELECT":
		if len(allMatch) < 2 || strings.HasPrefix(allMatch[index][0], "LIMIT") == false {
			break
		}
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}
		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
		if len(queryExpr) != 4 ||
			strings.EqualFold(queryExpr[1], "SELECT") == false ||
			strings.EqualFold(queryExpr[3], "LIMIT") == false {
			break
		}
		first, limit := 0, 0
		for i := 1; i < len(allMatch[index]); i++ {
			if len(strings.TrimSpace(allMatch[index][i])) == 0 {
				continue
			}

			if strings.HasPrefix(allMatch[index][i], "LIMIT") {
				if allMatch[index][i+2] != "" {
					first, _ = strconv.Atoi(allMatch[index][i+1])
					limit, _ = strconv.Atoi(allMatch[index][i+2])
				} else {
					limit, _ = strconv.Atoi(allMatch[index][i+1])
				}
				break
			}
		}
		sql = fmt.Sprintf(
			"SELECT * FROM "+
				"(SELECT GFORM.*, ROWNUM ROWNUM_ FROM (%s %s) GFORM WHERE ROWNUM <= %d)"+
				" WHERE ROWNUM_ >= %d",
			queryExpr[1], queryExpr[2], limit, first,
		)
	}
	return sql
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
// Note that it ignores the parameter `schema` in oracle database, as it is not necessary.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	result, err = d.DoGetAll(ctx, nil, "SELECT TABLE_NAME FROM USER_TABLES ORDER BY TABLE_NAME")
	if err != nil {
		return
	}
	for _, m := range result {
		for _, v := range m {
			tables = append(tables, v.String())
		}
	}
	return
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
//
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			"function TableFields supports only single table operations",
		)
	}
	useSchema := d.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`oracle_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result       gdb.Result
				link         gdb.Link
				structureSql = fmt.Sprintf(`
SELECT 
	COLUMN_NAME AS FIELD, 
	CASE DATA_TYPE  
	WHEN 'NUMBER' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
	WHEN 'FLOAT' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
	ELSE DATA_TYPE||'('||DATA_LENGTH||')' END AS TYPE  
FROM USER_TAB_COLUMNS WHERE TABLE_NAME = '%s' ORDER BY COLUMN_ID`,
					strings.ToUpper(table),
				)
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			structureSql, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(structureSql))
			result, err = d.DoGetAll(ctx, link, structureSql)
			if err != nil {
				return nil
			}
			fields = make(map[string]*gdb.TableField)
			for i, m := range result {
				fields[strings.ToLower(m["FIELD"].String())] = &gdb.TableField{
					Index: i,
					Name:  strings.ToLower(m["FIELD"].String()),
					Type:  strings.ToLower(m["TYPE"].String()),
				}
			}
			return fields
		},
	)
	if v != nil {
		fields = v.(map[string]*gdb.TableField)
	}
	return
}

// DoInsert inserts or updates data for given table.
// This function is usually used for custom interface definition, you do not need call it manually.
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `option` values are as follows:
// 0: insert:  just insert, if there's unique/primary key in the data, it returns error;
// 1: replace: if there's unique/primary key in the data, it deletes it from table and inserts a new one;
// 2: save:    if there's unique/primary key in the data, it updates it or else inserts a new one;
// 3: ignore:  if there's unique/primary key in the data, it ignores the inserting;
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Save operation is not supported by mssql driver`)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Replace operation is not supported by mssql driver`)
	}

	var (
		keys   []string
		values []string
		params []interface{}
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
	intoStr := make([]string, 0)
	for i := 0; i < len(list); i++ {
		for _, k := range keys {
			params = append(params, list[i][k])
		}
		values = append(values, valueHolderStr)
		intoStr = append(intoStr, fmt.Sprintf("INTO %s(%s) VALUES(%s)", table, keyStr, valueHolderStr))
		if len(intoStr) == option.BatchCount || (i == listLength-1 && len(valueHolder) > 0) {
			r, err := d.DoExec(ctx, link, fmt.Sprintf(
				"INSERT ALL %s SELECT * FROM DUAL",
				strings.Join(intoStr, " "),
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
			intoStr = intoStr[:0]
		}
	}
	return batchResult, nil
}

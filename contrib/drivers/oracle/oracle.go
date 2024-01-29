// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package oracle implements gdb.Driver, which supports operations for database Oracle.
//
// Note:
// 1. It does not support Save/Replace features.
// 2. It does not support LastInsertId.
package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	gora "github.com/sijms/go-ora/v2"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// Driver is the driver for oracle database.
type Driver struct {
	*gdb.Core
}

const (
	quoteChar = `"`
)

var (
	tablesSqlTmp         = `SELECT TABLE_NAME FROM USER_TABLES ORDER BY TABLE_NAME`
	newSqlReplacementTmp = `
SELECT * FROM (
	SELECT GFORM.*, ROWNUM ROWNUM_ FROM (%s %s) GFORM WHERE ROWNUM <= %d
) 
	WHERE ROWNUM_ > %d
`
	tableFieldsSqlTmp = `
SELECT 
    COLUMN_NAME AS FIELD, 
    CASE   
    WHEN (DATA_TYPE='NUMBER' AND NVL(DATA_SCALE,0)=0) THEN 'INT'||'('||DATA_PRECISION||','||DATA_SCALE||')'
    WHEN (DATA_TYPE='NUMBER' AND NVL(DATA_SCALE,0)>0) THEN 'FLOAT'||'('||DATA_PRECISION||','||DATA_SCALE||')'
    WHEN DATA_TYPE='FLOAT' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
    ELSE DATA_TYPE||'('||DATA_LENGTH||')' END AS TYPE,NULLABLE
FROM USER_TAB_COLUMNS WHERE TABLE_NAME = '%s' ORDER BY COLUMN_ID
`
)

func init() {
	var err error
	tableFieldsSqlTmp = formatSqlTmp(tableFieldsSqlTmp)
	newSqlReplacementTmp = formatSqlTmp(newSqlReplacementTmp)
	if err = gdb.Register(`oracle`, New()); err != nil {
		panic(err)
	}
}

// formatSqlTmp formats sql template string into one line.
func formatSqlTmp(sqlTmp string) string {
	var err error
	// format sql template string.
	sqlTmp, err = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(sqlTmp))
	if err != nil {
		panic(err)
	}
	sqlTmp, err = gregex.ReplaceString(`\s{2,}`, " ", gstr.Trim(sqlTmp))
	if err != nil {
		panic(err)
	}
	return sqlTmp
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
		underlyingDriverName = "oracle"
	)

	options := map[string]string{
		"CONNECTION TIMEOUT": "60",
		"PREFETCH_ROWS":      "25",
	}

	if config.Debug {
		options["TRACE FILE"] = "oracle_trace.log"
	}
	// [username:[password]@]host[:port][/service_name][?param1=value1&...&paramN=valueN]
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`@(.+?)/([\w\.\-]+)+`, "@$1/"+config.Name, source)
		}
	} else {
		if config.Extra != "" {
			var extraMap map[string]interface{}
			if extraMap, err = gstr.Parse(config.Extra); err != nil {
				return nil, err
			}
			for k, v := range extraMap {
				options[k] = gconv.String(v)
			}
		}
		source = gora.BuildUrl(
			config.Host, gconv.Int(config.Port), config.Name, config.User, config.Pass, options,
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

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return quoteChar, quoteChar
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	var index int
	newArgs = args
	// Convert placeholder char '?' to string ":vx".
	newSql, err = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf(":v%d", index)
	})
	if err != nil {
		return
	}
	newSql, err = gregex.ReplaceString("\"", "", newSql)
	if err != nil {
		return
	}
	newSql, err = d.parseSql(newSql)
	if err != nil {
		return
	}
	return d.Core.DoFilter(ctx, link, newSql, newArgs)
}

// parseSql does some replacement of the sql before commits it to underlying driver,
// for support of oracle server.
func (d *Driver) parseSql(toBeCommittedSql string) (string, error) {
	var (
		err       error
		operation = gstr.StrTillEx(toBeCommittedSql, " ")
		keyword   = strings.ToUpper(gstr.Trim(operation))
	)
	switch keyword {
	case "SELECT":
		toBeCommittedSql, err = d.handleSelectSqlReplacement(toBeCommittedSql)
		if err != nil {
			return "", err
		}
	}
	return toBeCommittedSql, nil
}

func (d *Driver) handleSelectSqlReplacement(toBeCommittedSql string) (newSql string, err error) {
	var (
		match  [][]string
		patten = `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,{0,1}\s*(\d*))`
	)
	match, err = gregex.MatchAllString(patten, toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(match) == 0 {
		return toBeCommittedSql, nil
	}
	var index = 1
	if len(match) < 2 || strings.HasPrefix(match[index][0], "LIMIT") == false {
		return toBeCommittedSql, nil
	}
	// only handle `SELECT ... LIMIT ...` statement.
	queryExpr, err := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", toBeCommittedSql)
	if err != nil {
		return "", err
	}
	if len(queryExpr) == 0 {
		return toBeCommittedSql, nil
	}
	if len(queryExpr) != 4 ||
		strings.EqualFold(queryExpr[1], "SELECT") == false ||
		strings.EqualFold(queryExpr[3], "LIMIT") == false {
		return toBeCommittedSql, nil
	}
	page, limit := 0, 0
	for i := 1; i < len(match[index]); i++ {
		if len(strings.TrimSpace(match[index][i])) == 0 {
			continue
		}
		if strings.HasPrefix(match[index][i], "LIMIT") {
			if match[index][i+2] != "" {
				page, err = strconv.Atoi(match[index][i+1])
				if err != nil {
					return "", err
				}
				limit, err = strconv.Atoi(match[index][i+2])
				if err != nil {
					return "", err
				}
				if page <= 0 {
					page = 1
				}
				limit = (page/limit + 1) * limit
				page, err = strconv.Atoi(match[index][i+1])
				if err != nil {
					return "", err
				}
			} else {
				limit, err = strconv.Atoi(match[index][i+1])
				if err != nil {
					return "", err
				}
			}
			break
		}
	}
	var newReplacedSql = fmt.Sprintf(
		newSqlReplacementTmp,
		queryExpr[1], queryExpr[2], limit, page,
	)
	return newReplacedSql, nil
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
// Note that it ignores the parameter `schema` in oracle database, as it is not necessary.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	// DO NOT use `usedSchema` as parameter for function `SlaveLink`.
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, tablesSqlTmp)
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
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result       gdb.Result
		link         gdb.Link
		usedSchema   = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		structureSql = fmt.Sprintf(tableFieldsSqlTmp, strings.ToUpper(table))
	)
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	result, err = d.DoSelect(ctx, link, structureSql)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		isNull := false
		if m["NULLABLE"].String() == "Y" {
			isNull = true
		}

		fields[m["FIELD"].String()] = &gdb.TableField{
			Index: i,
			Name:  m["FIELD"].String(),
			Type:  m["TYPE"].String(),
			Null:  isNull,
		}
	}
	return fields, nil
}

// DoInsert inserts or updates data for given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Save operation is not supported by oracle driver`,
		)

	case gdb.InsertOptionReplace:
		return nil, gerror.NewCode(
			gcode.CodeNotSupported,
			`Replace operation is not supported by oracle driver`,
		)
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

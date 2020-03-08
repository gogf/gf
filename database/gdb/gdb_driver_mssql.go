// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/denisenkom/go-mssqldb"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"
	"strconv"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

// DriverMssql is the driver for SQL server database.
type DriverMssql struct {
	*Core
}

// New creates and returns a database object for SQL server.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverMssql) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverMssql{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for mssql.
func (d *DriverMssql) Open(config *ConfigNode) (*sql.DB, error) {
	source := ""
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("sqlserver", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// GetChars returns the security char for this type of database.
func (d *DriverMssql) GetChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// HandleSqlBeforeExec deals with the sql string before commits it to underlying sql driver.
func (d *DriverMssql) HandleSqlBeforeExec(query string) string {
	var index int
	// Convert place holder char '?' to string "@px".
	str, _ := gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})
	str, _ = gregex.ReplaceString("\"", "", str)
	return d.parseSql(str)
}

func (d *DriverMssql) parseSql(sql string) string {
	// SELECT * FROM USER WHERE ID=1 LIMIT 1
	if m, _ := gregex.MatchString(`^SELECT(.+)LIMIT 1$`, sql); len(m) > 1 {
		return fmt.Sprintf(`SELECT TOP 1 %s`, m[1])
	}
	// SELECT * FROM USER WHERE AGE>18 ORDER BY ID DESC LIMIT 100, 200
	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if gregex.IsMatchString(patten, sql) == false {
		return sql
	}
	res, err := gregex.MatchAllString(patten, sql)
	if err != nil {
		return ""
	}
	index := 0
	keyword := strings.TrimSpace(res[index][0])
	keyword = strings.ToUpper(keyword)
	index++
	switch keyword {
	case "SELECT":
		// 不含LIMIT关键字则不处理
		if len(res) < 2 ||
			(strings.HasPrefix(res[index][0], "LIMIT") == false &&
				strings.HasPrefix(res[index][0], "limit") == false) {
			break
		}
		// 不含LIMIT则不处理
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}
		// 判断SQL中是否含有order by
		selectStr := ""
		orderStr := ""
		haveOrder := gregex.IsMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)
		if haveOrder {
			// 取order by 前面的字符串
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)
			if len(queryExpr) != 4 ||
				strings.EqualFold(queryExpr[1], "SELECT") == false ||
				strings.EqualFold(queryExpr[3], "ORDER BY") == false {
				break
			}
			selectStr = queryExpr[2]

			// 取order by表达式的值
			orderExpr, _ := gregex.MatchString("((?i)ORDER BY)(.+)((?i)LIMIT)", sql)
			if len(orderExpr) != 4 ||
				strings.EqualFold(orderExpr[1], "ORDER BY") == false ||
				strings.EqualFold(orderExpr[3], "LIMIT") == false {
				break
			}
			orderStr = orderExpr[2]
		} else {
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
			if len(queryExpr) != 4 ||
				strings.EqualFold(queryExpr[1], "SELECT") == false ||
				strings.EqualFold(queryExpr[3], "LIMIT") == false {
				break
			}
			selectStr = queryExpr[2]
		}

		// 取limit后面的取值范围
		first, limit := 0, 0
		for i := 1; i < len(res[index]); i++ {
			if len(strings.TrimSpace(res[index][i])) == 0 {
				continue
			}

			if strings.HasPrefix(res[index][i], "LIMIT") ||
				strings.HasPrefix(res[index][i], "limit") {
				first, _ = strconv.Atoi(res[index][i+1])
				limit, _ = strconv.Atoi(res[index][i+2])
				break
			}
		}

		if haveOrder {
			sql = fmt.Sprintf(
				"SELECT * FROM "+
					"(SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROWNUMBER_, %s ) as TMP_ "+
					"WHERE TMP_.ROWNUMBER_ > %d AND TMP_.ROWNUMBER_ <= %d",
				orderStr, selectStr, first, limit,
			)
		} else {
			if first == 0 {
				first = limit
			} else {
				first = limit - first
			}
			sql = fmt.Sprintf(
				"SELECT * FROM (SELECT TOP %d * FROM (SELECT TOP %d %s) as TMP1_ ) as TMP2_ ",
				first, limit, selectStr,
			)
		}
	default:
	}
	return sql
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *DriverMssql) Tables(schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.DB.GetSlave(schema...)
	if err != nil {
		return nil, err
	}

	result, err = d.DB.DoGetAll(link, `SELECT NAME FROM SYSOBJECTS WHERE XTYPE='U' AND STATUS >= 0 ORDER BY NAME`)
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

// TableFields retrieves and returns the fields information of specified table of current schema.
func (d *DriverMssql) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	table = gstr.Trim(table)
	if gstr.Contains(table, " ") {
		panic("function TableFields supports only single table operations")
	}
	checkSchema := d.DB.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v := d.DB.GetCache().GetOrSetFunc(
		fmt.Sprintf(`mssql_table_fields_%s_%s`, table, checkSchema), func() interface{} {
			var result Result
			var link *sql.DB
			link, err = d.DB.GetSlave(checkSchema)
			if err != nil {
				return nil
			}
			result, err = d.DB.DoGetAll(link, fmt.Sprintf(`
			SELECT c.name as FIELD, CASE t.name 
				WHEN 'numeric' THEN t.name + '(' + convert(varchar(20),c.xprec) + ',' + convert(varchar(20),c.xscale) + ')' 
				WHEN 'char' THEN t.name + '(' + convert(varchar(20),c.length)+ ')'
				WHEN 'varchar' THEN t.name + '(' + convert(varchar(20),c.length)+ ')'
				ELSE t.name + '(' + convert(varchar(20),c.length)+ ')' END as TYPE
			FROM systypes t,syscolumns c WHERE t.xtype=c.xtype 
			AND c.id = (SELECT id FROM sysobjects WHERE name='%s') 
			ORDER BY c.colid`, strings.ToUpper(table)))
			if err != nil {
				return nil
			}
			fields = make(map[string]*TableField)
			for i, m := range result {
				fields[strings.ToLower(m["FIELD"].String())] = &TableField{
					Index: i,
					Name:  strings.ToLower(m["FIELD"].String()),
					Type:  strings.ToLower(m["TYPE"].String()),
				}
			}
			return fields
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}

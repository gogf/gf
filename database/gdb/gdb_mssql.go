// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/lib/pq"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

package gdb

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

type dbMssql struct {
	*dbBase
}

func (db *dbMssql) Open(config *ConfigNode) (*sql.DB, error) {
	source := ""
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}
	if db, err := sql.Open("sqlserver", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

func (db *dbMssql) getChars() (charLeft string, charRight string) {
	return "\"", "\""
}

func (db *dbMssql) handleSqlBeforeExec(query string) string {
	index := 0
	str, _ := gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})

	str, _ = gregex.ReplaceString("\"", "", str)

	return db.parseSql(str)
}

func (db *dbMssql) parseSql(sql string) string {
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
		haveOrderby := gregex.IsMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)
		if haveOrderby {
			// 取order by 前面的字符串
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)
			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "ORDER BY") == false {
				break
			}
			selectStr = queryExpr[2]

			// 取order by表达式的值
			orderExpr, _ := gregex.MatchString("((?i)ORDER BY)(.+)((?i)LIMIT)", sql)
			if len(orderExpr) != 4 || strings.EqualFold(orderExpr[1], "ORDER BY") == false || strings.EqualFold(orderExpr[3], "LIMIT") == false {
				break
			}
			orderStr = orderExpr[2]
		} else {
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "LIMIT") == false {
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

			if strings.HasPrefix(res[index][i], "LIMIT") || strings.HasPrefix(res[index][i], "limit") {
				first, _ = strconv.Atoi(res[index][i+1])
				limit, _ = strconv.Atoi(res[index][i+2])
				break
			}
		}

		if haveOrderby {
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

// 返回当前数据库所有的数据表名称
// TODO
func (db *dbMssql) Tables() (tables []string, err error) {
	return
}

// 获得指定表表的数据结构，构造成map哈希表返回，其中键名为表字段名称，键值暂无用途(默认为字段数据类型).
func (db *dbMssql) TableFields(table string) (fields map[string]*TableField, err error) {
	v := db.cache.GetOrSetFunc("mssql_table_fields_"+table, func() interface{} {
		var result Result
		result, err = db.GetAll(fmt.Sprintf(`
		SELECT c.name as FIELD, CASE t.name 
			WHEN 'numeric' THEN t.name + '(' + convert(varchar(20),c.xprec) + ',' + convert(varchar(20),c.xscale) + ')' 
			WHEN 'char' THEN t.name + '(' + convert(varchar(20),c.length)+ ')'
			WHEN 'varchar' THEN t.name + '(' + convert(varchar(20),c.length)+ ')'
			ELSE t.name + '(' + convert(varchar(20),c.length)+ ')' END as TYPE
		FROM systypes t,syscolumns c WHERE t.xtype=c.xtype AND c.id = (SELECT id FROM sysobjects WHERE name='%s') ORDER BY c.colid`, strings.ToUpper(table)))
		if err != nil {
			return nil
		}
		fields = make(map[string]*TableField)
		for i, m := range result {
			// SQLServer返回的field为大写的需要转为小写的
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

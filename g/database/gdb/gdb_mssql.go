// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
/*
@author wenzi1<liyz23@qq.com>
@date 20181109
说明：
    1.需要导入sqlserver驱动： github.com/denisenkom/go-mssqldb
    2.不支持save/replace方法
    3.不支持LastInsertId方法
*/
package gdb

import (
	"database/sql"
	"fmt"
	"gitee.com/johng/gf/g/util/gregex"
	"strconv"
	"strings"
)

// 数据库链接对象
type dbMssql struct {
	*dbBase
}

// 创建SQL操作对象
func (db *dbMssql) Open(config *ConfigNode) (*sql.DB, error) {
	source := ""
	if config.Linkinfo != "" {
		source = config.Linkinfo
	} else {
		source = fmt.Sprintf("user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name)
	}
	if db, err := sql.Open("sqlserver", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// 获得关键字操作符
func (db *dbMssql) getChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// 在执行sql之前对sql进行进一步处理
func (db *dbMssql) handleSqlBeforeExec(query string) string {
	index := 0
	str, _ := gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})

	str, _ = gregex.ReplaceString("\"", "", str)

	return db.parseSql(str)
}

//将MYSQL的SQL语法转换为MSSQL的语法
//1.由于mssql不支持limit写法所以需要对mysql中的limit用法做转换
func (db *dbMssql) parseSql(sql string) string {
	//下面的正则表达式匹配出SELECT和INSERT的关键字后分别做不同的处理，如有LIMIT则将LIMIT的关键字也匹配出
	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if gregex.IsMatchString(patten, sql) == false {
		fmt.Println("not matched..")
		return sql
	}

	res, err := gregex.MatchAllString(patten, sql)
	if err != nil {
		fmt.Println("MatchString error.", err)
		return ""
	}

	index := 0
	keyword := strings.TrimSpace(res[index][0])
	keyword = strings.ToUpper(keyword)

	index++
	switch keyword {
	case "SELECT":
		//不含LIMIT关键字则不处理
		if len(res) < 2 || (strings.HasPrefix(res[index][0], "LIMIT") == false && strings.HasPrefix(res[index][0], "limit") == false) {
			break
		}

		//不含LIMIT则不处理
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}

		//判断SQL中是否含有order by
		selectStr := ""
		orderbyStr := ""
		haveOrderby := gregex.IsMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)
		if haveOrderby {
			//取order by 前面的字符串
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)ORDER BY)", sql)

			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "ORDER BY") == false {
				break
			}
			selectStr = queryExpr[2]

			//取order by表达式的值
			orderbyExpr, _ := gregex.MatchString("((?i)ORDER BY)(.+)((?i)LIMIT)", sql)
			if len(orderbyExpr) != 4 || strings.EqualFold(orderbyExpr[1], "ORDER BY") == false || strings.EqualFold(orderbyExpr[3], "LIMIT") == false {
				break
			}
			orderbyStr = orderbyExpr[2]
		} else {
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "LIMIT") == false {
				break
			}
			selectStr = queryExpr[2]
		}

		//取limit后面的取值范围
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
			sql = fmt.Sprintf("SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROWNUMBER_, %s   ) as TMP_ WHERE TMP_.ROWNUMBER_ > %d AND TMP_.ROWNUMBER_ <= %d", orderbyStr, selectStr, first, limit)
		} else {
			if first == 0 {
				first = limit
			} else {
				first = limit - first
			}
			sql = fmt.Sprintf("SELECT * FROM (SELECT TOP %d * FROM (SELECT TOP %d %s) as TMP1_ ) as TMP2_ ", first, limit, selectStr)
		}
	default:
	}
	return sql
}

// 获得指定表表的数据结构，构造成map哈希表返回，其中键名为表字段名称，键值暂无用途(默认为字段数据类型).
func (db *dbMssql) getTableFields(table string) (fields map[string]string, err error) {
	// 缓存不存在时会查询数据表结构，缓存后不过期，直至程序重启(重新部署)
	v := db.cache.GetOrSetFunc("table_fields_"+table, func() interface{} {
		result := (Result)(nil)
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
		fields = make(map[string]string)
		for _, m := range result {
			fields[strings.ToLower(m["FIELD"].String())] = strings.ToLower(m["TYPE"].String()) //sqlserver返回的field为大写的需要转为小写的
		}
		return fields
	}, 0)
	if err == nil {
		fields = v.(map[string]string)
	}
	return
}

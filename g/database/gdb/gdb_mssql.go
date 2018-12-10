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


var linkMssql = &dbmssql{}

// 数据库链接对象
type dbmssql struct {
	Db
}

// 创建SQL操作对象
func (db *dbmssql) Open(c *ConfigNode) (*sql.DB, error) {
	var source string
	if c.Linkinfo != "" {
		source = c.Linkinfo
	} else {
		source = fmt.Sprintf("user id=%s;password=%s;server=%s;port=%s;database=%s;encrypt=disable", c.User, c.Pass, c.Host, c.Port, c.Name)
	}
	if db, err := sql.Open("sqlserver", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// 获得关键字操作符 - 左
func (db *dbmssql) getQuoteCharLeft() string {
	return "\""
}

// 获得关键字操作符 - 右
func (db *dbmssql) getQuoteCharRight() string {
	return "\""
}

// 在执行sql之前对sql进行进一步处理
func (db *dbmssql) handleSqlBeforeExec(q *string) *string {
	index := 0
	str, _ := gregex.ReplaceStringFunc("\\?", *q, func(s string) string {
		index++
		return fmt.Sprintf("@p%d", index)
	})

	str, _ = gregex.ReplaceString("\"", "", str)

	return db.parseSql(&str)
}

//将MYSQL的SQL语法转换为MSSQL的语法
//1.由于mssql不支持limit写法所以需要对mysql中的limit用法做转换
func (db *dbmssql) parseSql(sql *string) *string {
	//下面的正则表达式匹配出SELECT和INSERT的关键字后分别做不同的处理，如有LIMIT则将LIMIT的关键字也匹配出
	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if gregex.IsMatchString(patten, *sql) == false {
		fmt.Println("not matched..")
		return sql
	}

	res, err := gregex.MatchAllString(patten, *sql)
	if err != nil {
		fmt.Println("MatchString error.", err)
		return nil
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
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", *sql) == false {
			break
		} 

		//判断SQL中是否含有order by
		selectStr := ""
		orderbyStr := ""
		haveOrderby := gregex.IsMatchString("((?i)SELECT)(.+)((?i)ORDER BY)", *sql)
		if haveOrderby {
			//取order by 前面的字符串
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)ORDER BY)", *sql)
			
			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "ORDER BY") == false{
				break
			}
			selectStr = queryExpr[2]

			//取order by表达式的值
			orderbyExpr, _ := gregex.MatchString("((?i)ORDER BY)(.+)((?i)LIMIT)", *sql)
			if len(orderbyExpr) != 4 || strings.EqualFold(orderbyExpr[1], "ORDER BY") == false || strings.EqualFold(orderbyExpr[3], "LIMIT") == false{
				break
			}
			orderbyStr = orderbyExpr[2]
		} else {
			queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", *sql)
			if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "LIMIT") == false{
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
			*sql = fmt.Sprintf("SELECT * FROM (SELECT ROW_NUMBER() OVER (ORDER BY %s) as ROWNUMBER_, %s   ) as TMP_ WHERE TMP_.ROWNUMBER_ > %d AND TMP_.ROWNUMBER_ <= %d", orderbyStr, selectStr, first, limit)
		} else {
			if first == 0 {
				first = limit
			} else {
				first = limit - first
			}
			*sql = fmt.Sprintf("SELECT * FROM (SELECT TOP %d * FROM (SELECT TOP %d %s) as TMP1_ ) as TMP2_ ", first, limit, selectStr)
		}
	default:
	}
	return sql
}
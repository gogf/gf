// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
/*
@author wenzi1<liyz23@qq.com>
@date 20181026
说明：
    1.需要导入oracle驱动： github.com/mattn/go-oci8
    2.不支持save/replace方法，可以调用这2个方法估计会报错，还没测试过,(应该是可以通过oracle的merge来实现这2个功能的，还没仔细研究)
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
type dbOracle struct {
	*dbBase
}

// 创建SQL操作对象
func (db *dbOracle) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.Linkinfo != "" {
		source = config.Linkinfo
	} else {
		source = fmt.Sprintf("%s/%s@%s", config.User, config.Pass, config.Name)
	}
	if db, err := sql.Open("oci8", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// 获得关键字操作符
func (db *dbOracle) getChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// 在执行sql之前对sql进行进一步处理
func (db *dbOracle) handleSqlBeforeExec(query string) string {
	index := 0
	str, _ := gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf(":%d", index)
	})

	str, _ = gregex.ReplaceString("\"", "", str)

	return db.parseSql(str)
}

//由于ORACLE中对LIMIT和批量插入的语法与MYSQL不一致，所以这里需要对LIMIT和批量插入做语法上的转换
func (db *dbOracle) parseSql(sql string) string {
	//下面的正则表达式匹配出SELECT和INSERT的关键字后分别做不同的处理，如有LIMIT则将LIMIT的关键字也匹配出
	patten := `^\s*(?i)(SELECT)|(INSERT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
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

		//取limit前面的字符串
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}

		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
		if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "LIMIT") == false {
			break
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

		//也可以使用between,据说这种写法的性能会比between好点,里层SQL中的ROWNUM_ >= limit可以缩小查询后的数据集规模
		sql = fmt.Sprintf("SELECT * FROM (SELECT GFORM.*, ROWNUM ROWNUM_ FROM (%s %s) GFORM WHERE ROWNUM <= %d) WHERE ROWNUM_ >= %d", queryExpr[1], queryExpr[2], limit, first)
	case "INSERT":
		//获取VALUE的值，匹配所有带括号的值,会将INSERT INTO后的值匹配到，所以下面的判断语句会判断数组长度是否小于3
		valueExpr, err := gregex.MatchAllString(`(\s*\(([^\(\)]*)\))`, sql)
		if err != nil {
			return sql
		}

		//判断VALUE后的值是否有多个，只有在批量插入的时候才需要做转换，如只有1个VALUE则不需要做转换
		if len(valueExpr) < 3 {
			break
		}

		//获取INTO后面的值
		tableExpr, err := gregex.MatchString(`(?i)\s*(INTO\s+\w+\(([^\(\)]*)\))`, sql)
		if err != nil {
			return sql
		}
		tableExpr[0] = strings.TrimSpace(tableExpr[0])

		sql = "INSERT ALL"
		for i := 1; i < len(valueExpr); i++ {
			sql += fmt.Sprintf(" %s VALUES%s", tableExpr[0], strings.TrimSpace(valueExpr[i][0]))
		}
		sql += " SELECT 1 FROM DUAL"

	default:
	}
	return sql
}

// 获得指定表表的数据结构，构造成map哈希表返回，其中键名为表字段名称，键值暂无用途(默认为字段数据类型).
func (db *dbOracle) getTableFields(table string) (fields map[string]string, err error) {
	// 缓存不存在时会查询数据表结构，缓存后不过期，直至程序重启(重新部署)
	v := db.cache.GetOrSetFunc("table_fields_"+table, func() interface{} {
		result := (Result)(nil)
		result, err = db.GetAll(fmt.Sprintf(`
		SELECT COLUMN_NAME AS FIELD, CASE DATA_TYPE 
		    WHEN 'NUMBER' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
			WHEN 'FLOAT' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
			ELSE DATA_TYPE||'('||DATA_LENGTH||')' END AS TYPE  
		FROM USER_TAB_COLUMNS WHERE TABLE_NAME = '%s' ORDER BY COLUMN_ID`, strings.ToUpper(table)))
		if err != nil {
			return nil
		}

		fields = make(map[string]string)
		for _, m := range result {
			fields[strings.ToLower(m["FIELD"].String())] = strings.ToLower(m["TYPE"].String()) //ORACLE返回的值默认都是大写的，需要转为小写
		}
		return fields
	}, 0)
	if err == nil {
		fields = v.(map[string]string)
	}
	return
}

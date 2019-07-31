// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/g/text/gregex"
	"strings"
)

// PostgreSQL的适配.
//
// 使用时需要import:
//
// _ "github.com/gogf/gf/third/github.com/lib/pq"
//
// @todo 需要完善replace和save的操作覆盖

// 数据库链接对象
type dbPgsql struct {
	*dbBase
}

// 创建SQL操作对象，内部采用了lazy link处理
func (db *dbPgsql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.User, config.Pass, config.Host, config.Port, config.Name)
	}
	if db, err := sql.Open("postgres", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// 动态切换数据库
func (db *dbPgsql) setSchema(sqlDb *sql.DB, schema string) error {
	_, err := sqlDb.Exec("SET search_path TO " + schema)
	return err
}

// 获得关键字操作符
func (db *dbPgsql) getChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// 在执行sql之前对sql进行进一步处理
func (db *dbPgsql) handleSqlBeforeExec(query string) string {
	index := 0
	query, _ = gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf("$%d", index)
	})
	// 分页语法替换
	query, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $1 OFFSET $2`, query)
	return query
}

func (db *dbPgsql) getTableFields(table string) (fields map[string]string, err error) {
	// 缓存不存在时会查询数据表结构，缓存后不过期，直至程序重启(重新部署)
	table, _ = gregex.ReplaceString("\"", "", table)
	v := db.cache.GetOrSetFunc("pgsql_table_fields_"+table, func() interface{} {
		result := (Result)(nil)
		result, err = db.GetAll(fmt.Sprintf(`
		SELECT a.attname AS field, t.typname AS type FROM pg_class c, pg_attribute a 
        LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,pg_type t
        WHERE c.relname = '%s' and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid ORDER BY a.attnum`, strings.ToLower(table)))
		if err != nil {
			return nil
		}

		fields = make(map[string]string)
		for _, m := range result {
			fields[m["field"].String()] = m["type"].String()
		}
		return fields
	}, 0)
	if err == nil {
		fields = v.(map[string]string)
	}
	return
}

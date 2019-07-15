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
		source = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", config.User, config.Pass, config.Host, config.Port, config.Name)
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

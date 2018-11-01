// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// @author wxkj<wxscz@qq.com>

package gdb

import (
	"database/sql"
)

// 使用时需要import:
// _ "gitee.com/johng/gf/third/github.com/mattn/go-sqlite3"

// Sqlite接口对象
// @author wxkj<wxscz@qq.com>
var linkSqlite = &dbsqlite{}


// 数据库链接对象
type dbsqlite struct {
	Db
}

func (db *dbsqlite) Open(c *ConfigNode) (*sql.DB, error) {
	var source string
	if c.Linkinfo != "" {
		source = c.Linkinfo
	} else {
		source = c.Name
	}
	if db, err := sql.Open("sqlite3", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// 获得关键字操作符 - 左
func (db *dbsqlite) getQuoteCharLeft() string {
	return "`"
}

// 获得关键字操作符 - 右
func (db *dbsqlite) getQuoteCharRight() string {
	return "`"
}

// 在执行sql之前对sql进行进一步处理
// @todo 需要增加对Save方法的支持，可使用正则来实现替换，
// @todo 将ON DUPLICATE KEY UPDATE触发器修改为两条SQL语句(INSERT OR IGNORE & UPDATE)
func (db *dbsqlite) handleSqlBeforeExec(q *string) *string {

	return q
}
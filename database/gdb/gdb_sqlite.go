// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/mattn/go-sqlite3"
// 2. It does not support Save/Replace features.

package gdb

import (
	"database/sql"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"
)

type dbSqlite struct {
	*dbBase
}

// Open creates and returns a underlying sql.DB object for sqlite.
func (db *dbSqlite) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = config.Name
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("sqlite3", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// getChars returns the security char for this type of database.
func (db *dbSqlite) getChars() (charLeft string, charRight string) {
	return "`", "`"
}

// Tables retrieves and returns the tables of current schema.
// TODO
func (db *dbSqlite) Tables(schema ...string) (tables []string, err error) {
	return
}

// TableFields retrieves and returns the fields information of specified table of current schema.
// TODO
func (db *dbSqlite) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	table = gstr.Trim(table)
	if gstr.Contains(table, " ") {
		panic("function TableFields supports only single table operations")
	}
	return
}

// handleSqlBeforeExec deals with the sql string before commits it to underlying sql driver.
// @todo 需要增加对Save方法的支持，可使用正则来实现替换，
// @todo 将ON DUPLICATE KEY UPDATE触发器修改为两条SQL语句(INSERT OR IGNORE & UPDATE)
func (db *dbSqlite) handleSqlBeforeExec(sql string) string {
	return sql
}

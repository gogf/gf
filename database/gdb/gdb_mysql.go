// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"

	_ "github.com/gf-third/mysql"
)

// 数据库链接对象
type dbMysql struct {
	*dbBase
}

// 创建SQL操作对象，内部采用了lazy link处理
func (db *dbMysql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=true&loc=Local",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset)
	}
	if db, err := sql.Open("gf-mysql", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// getChars returns the quote chars for field.
func (db *dbMysql) getChars() (charLeft string, charRight string) {
	return "`", "`"
}

// handleSqlBeforeExec handles the sql before posts it to database.
func (db *dbMysql) handleSqlBeforeExec(query string) string {
	return query
}

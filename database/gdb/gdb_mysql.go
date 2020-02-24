// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/internal/intlog"

	_ "github.com/gf-third/mysql"
)

type dbMysql struct {
	*dbBase
}

// Open creates and returns a underlying database connection with given configuration.
func (db *dbMysql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true&parseTime=true&loc=Local",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
	}
	intlog.Printf("Open: %s", source)
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
func (db *dbMysql) handleSqlBeforeExec(sql string) string {
	return sql
}

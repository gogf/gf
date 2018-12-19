// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gdb

import (
    "fmt"
    "database/sql"
)

// 数据库链接对象
type dbMysql struct {
    *dbBase
}

// 创建SQL操作对象，内部采用了lazy link处理
func (db *dbMysql) Open (config *ConfigNode) (*sql.DB, error) {
    var source string
    if config.Linkinfo != "" {
        source = config.Linkinfo
    } else {
        source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true",
            config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset)
    }
    if db, err := sql.Open("mysql", source); err == nil {
        return db, nil
    } else {
        return nil, err
    }
}

// 获得关键字操作符
func (db *dbMysql) getChars () (charLeft string, charRight string) {
    return "`", "`"
}

// 在执行sql之前对sql进行进一步处理
func (db *dbMysql) handleSqlBeforeExec(query string) string {
    return query
}
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

// MySQL接口对象
var linkMysql = &dbmysql{}


// 数据库链接对象
type dbmysql struct {
    Db
}

// 创建SQL操作对象，内部采用了lazy link处理
func (db *dbmysql) Open (c *ConfigNode) (*sql.DB, error) {
    var source string
    if c.Linkinfo != "" {
        source = c.Linkinfo
    } else {
        source = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pass, c.Host, c.Port, c.Name)
    }
    if db, err := sql.Open("mysql", source); err == nil {
        return db, nil
    } else {
        return nil, err
    }
}

// 获得关键字操作符 - 左
func (db *dbmysql) getQuoteCharLeft () string {
    return "`"
}

// 获得关键字操作符 - 右
func (db *dbmysql) getQuoteCharRight () string {
    return "`"
}

// 在执行sql之前对sql进行进一步处理
func (db *dbmysql) handleSqlBeforeExec(q *string) *string {
    return q
}
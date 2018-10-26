// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.


package gdb

import (
    "fmt"
    "regexp"
    "database/sql"
)

// PostgreSQL的适配.
// 使用时需要import:
// _ "gitee.com/johng/gf/third/github.com/lib/pq"
// @todo 需要完善replace和save的操作覆盖

// PostgreSQL接口对象
var linkPgsql = &dbpgsql{}


// 数据库链接对象
type dbpgsql struct {
    Db
}

// 创建SQL操作对象，内部采用了lazy link处理
func (db *dbpgsql) Open (c *ConfigNode) (*sql.DB, error) {
    var source string
    if c.Linkinfo != "" {
        source = c.Linkinfo
    } else {
        source = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", c.User, c.Pass, c.Host, c.Port, c.Name)
    }
    if db, err := sql.Open("postgres", source); err == nil {
        return db, nil
    } else {
        return nil, err
    }
}

// 获得关键字操作符 - 左
func (db *dbpgsql) getQuoteCharLeft () string {
    return "\""
}

// 获得关键字操作符 - 右
func (db *dbpgsql) getQuoteCharRight () string {
    return "\""
}

// 在执行sql之前对sql进行进一步处理
func (db *dbpgsql) handleSqlBeforeExec(q *string) *string {
    reg   := regexp.MustCompile("\\?")
    index := 0
    str   := reg.ReplaceAllStringFunc(*q, func (s string) string {
        index ++
        return fmt.Sprintf("$%d", index)
    })
    return &str
}
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

// 数据库链接对象
type dbPgsql struct {
    *dbBase
}

// 创建SQL操作对象，内部采用了lazy link处理
func (db *dbPgsql) Open (config *ConfigNode) (*sql.DB, error) {
    var source string
    if config.Linkinfo != "" {
        source = config.Linkinfo
    } else {
        source = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", config.User, config.Pass, config.Host, config.Port, config.Name)
    }
    if db, err := sql.Open("postgres", source); err == nil {
        return db, nil
    } else {
        return nil, err
    }
}

// 获得关键字操作符
func (db *dbPgsql) getChars () (charLeft string, charRight string) {
    return "\"", "\""
}

// 在执行sql之前对sql进行进一步处理
func (db *dbPgsql) handleSqlBeforeExec(query string) string {
    reg   := regexp.MustCompile("\\?")
    index := 0
    str   := reg.ReplaceAllStringFunc(query, func (s string) string {
        index ++
        return fmt.Sprintf("$%d", index)
    })
    return str
}
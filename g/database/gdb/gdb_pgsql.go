// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "database/sql"
    "fmt"
    "regexp"
    "gitee.com/johng/gf/g/os/glog"
)

// postgresql的适配
// @todo 需要完善replace和save的操作覆盖

// 数据库链接对象
type pgsqlLink struct {
    dbLink
}

// 创建SQL操作对象，内部采用了lazy link处理
func (l *pgsqlLink) Open (c *ConfigNode) (*sql.DB, error) {
    var dbsource string
    if c.Linkinfo != "" {
        dbsource = c.Linkinfo
    } else {
        dbsource = fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s", c.User, c.Pass, c.Host, c.Port, c.Name)
    }
    db, err := sql.Open("postgres", dbsource)
    if err != nil {
        glog.Fatal(err)
    }
    return db, err
}

// 获得关键字操作符 - 左
func (l *pgsqlLink) getQuoteCharLeft () string {
    return "\""
}

// 获得关键字操作符 - 右
func (l *pgsqlLink) getQuoteCharRight () string {
    return "\""
}

// 在执行sql之前对sql进行进一步处理
func (l *pgsqlLink) handleSqlBeforeExec(q *string) *string {
    reg   := regexp.MustCompile("\\?")
    index := 0
    str   := reg.ReplaceAllStringFunc(*q, func (s string) string {
        index ++
        return fmt.Sprintf("$%d", index)
    })
    return &str
}


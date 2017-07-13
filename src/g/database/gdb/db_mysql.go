package gdb

import (
    "database/sql"
    "fmt"
    "log"
)

// 数据库链接对象
type mysqlLink struct {
    dbLink
}

// 创建SQL操作对象，内部采用了lazy link处理
func (l *mysqlLink) Open (c *ConfigNode) (*sql.DB, error) {
    var dbsource string
    if c.Linkinfo != "" {
        dbsource = c.Linkinfo
    } else {
        dbsource = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pass, c.Host, c.Port, c.Name)
    }
    db, err := sql.Open("mysql", dbsource)
    if err != nil {
        log.Fatal(err)
    }
    return db, err
}

// 获得关键字操作符 - 左
func (l *mysqlLink) getQuoteCharLeft () string {
    return "`"
}

// 获得关键字操作符 - 右
func (l *mysqlLink) getQuoteCharRight () string {
    return "`"
}
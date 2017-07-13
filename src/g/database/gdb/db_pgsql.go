package gdb

import (
    "database/sql"
    "fmt"
    "log"
)

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
        log.Fatal(err)
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

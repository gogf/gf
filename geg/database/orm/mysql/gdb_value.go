package main

import (
    "gitee.com/johng/gf/g/database/gdb"
    "fmt"
)

func main() {
    gdb.AddDefaultConfigNode(gdb.ConfigNode {
        Host    : "127.0.0.1",
        Port    : "3306",
        User    : "root",
        Pass    : "123456",
        Name    : "test",
        Type    : "mysql",
        Role    : "master",
        Charset : "utf8",
    })
    db, err := gdb.New()
    if err != nil {
        panic(err)
    }
    // 开启调试模式，以便于记录所有执行的SQL
    db.SetDebug(true)

    r, _ := db.Table("user").Where("uid=?", 1).One()
    if r != nil {
        fmt.Println(r["uid"].Int())
        fmt.Println(r["name"].String())
    }
}
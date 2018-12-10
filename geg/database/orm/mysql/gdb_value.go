package main

import (
    "gitee.com/johng/gf/g/database/gdb"
    "fmt"
    "gitee.com/johng/gf/g/encoding/gparser"
)

func main() {
    gdb.AddDefaultConfigNode(gdb.ConfigNode {
        Host    : "192.168.1.11",
        Port    : "3306",
        User    : "root",
        Pass    : "8692651",
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

    r, _ := db.Table("user").All()
    if r != nil {
        fmt.Println(r.ToList())
        b, e := gparser.VarToJson(r.ToList())
        fmt.Println(e)
        fmt.Println(string(b))
    }
}
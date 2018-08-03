package main

import (
    "gitee.com/johng/gf/g/database/gdb"
    "fmt"
    "gitee.com/johng/gf/g/encoding/gparser"
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

    one, err := db.Table("user").Where("uid=?", 1).One()
    if err != nil {
        panic(err)
    }

    // 使用内置方法转换为json/xml
    fmt.Println(one.ToJson())
    fmt.Println(one.ToXml())

    // 自定义方法方法转换为json/xml
    jsonContent, _ := gparser.VarToJson(one.ToMap())
    fmt.Println(jsonContent)
    xmlContent, _ := gparser.VarToXml(one.ToMap())
    fmt.Println(xmlContent)
}
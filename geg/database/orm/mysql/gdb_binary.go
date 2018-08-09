package main

import (
    "gitee.com/johng/gf/g/database/gdb"
    "fmt"
    "gitee.com/johng/gf/g/crypto/gaes"
    "gitee.com/johng/gf/g"
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

    key  := "0123456789123456"


    name := "john"
    encryptedName, err := gaes.Encrypt([]byte(name), []byte(key))
    if err != nil {
        fmt.Println(err)
    }

    // 写入
    r, err := db.Table("user").Data(g.Map{
        "uid"  : 1,
        "name" : encryptedName,
    }).Save()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(r.RowsAffected())

    // 查询
    one, err := db.Table("user").Where("name=?", encryptedName).One()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(one.ToMap())
}
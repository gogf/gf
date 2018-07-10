package main

import (
    "fmt"
    "gitee.com/johng/gf/g/database/gdb"
)

var db *gdb.Db

// 初始化配置及创建数据库
func init () {
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
    db, _ = gdb.New()
}


func main() {
    db.SetDebug(true)
    // 执行3条SQL查询
    for i := 1; i <= 3; i++ {
        db.Table("user").Where("uid=?", i).One()
    }
    // 构造一条错误查询
    db.Table("user").Where("no_such_field=?", "just_test").One()

    for k, v := range db.GetQueriedSqls() {
        fmt.Println(k, ":")
        fmt.Println("Sql  :", v.Sql)
        fmt.Println("Args :", v.Args)
        fmt.Println("Error:", v.Error)
        fmt.Println("Cost :", v.Cost)
        fmt.Println("Func :", v.Func)
    }
}
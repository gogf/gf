package main

//import (
//    _ "github.com/gogf/gf/third/github.com/mattn/go-sqlite3"
//    "github.com/gogf/gf/g/database/gdb"
//    "github.com/gogf/gf/g"
//    "fmt"
//)
//
//func main() {
//    gdb.SetConfig(gdb.Config{
//        "default": gdb.ConfigGroup{
//            gdb.ConfigNode{
//                Name: "/tmp/my.db",
//                Type: "sqlite",
//            },
//        },
//    })
//    db := g.Database()
//    if db == nil {
//        panic("db create failed")
//    }
//    defer db.Close()
//
//    // 创建表
//    sql := `CREATE TABLE user (
//        uid  INT PRIMARY KEY NOT NULL,
//        name VARCHAR(30) NOT NULL
//    );`
//    if _, err := db.Exec(sql); err != nil {
//        fmt.Println(err)
//    }
//
//    // 写入数据
//    result, err := db.Table("user").Data(g.Map{"uid" : 1, "name" : "john"}).Save()
//    if err == nil {
//        fmt.Println(result.RowsAffected())
//    } else {
//        fmt.Println(err)
//    }
//
//    // 删除表
//    sql  = `DROP TABLE user;`
//    if _, err := db.Exec(sql); err != nil {
//        fmt.Println(err)
//    }
//}

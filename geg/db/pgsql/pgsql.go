package main

import (
    "fmt"
    "time"
    "strconv"
    "gf/g/database/gdb"
)

// 本文件用于gf框架的postgresql数据库操作示例，不作为单元测试使用

var db gdb.Link

func init () {
    gdb.AddDefaultConfigNode(gdb.ConfigNode {
        Host : "127.0.0.1",
        Port : 5432,
        User : "postgres",
        Pass : "123456",
        Name : "test",
        Type : "pgsql",
    })
    db, _ = gdb.Instance()
}



// 创建测试数据库
func create() {
    fmt.Println("create:")
    _, err := db.Exec("CREATE SCHEMA IF NOT EXISTS \"test\"")
    if (err != nil) {
        fmt.Println(err)
    }

    s := `CREATE TABLE IF NOT EXISTS "user" (
            uid  int  PRIMARY KEY,
            name TEXT NOT NULL
        )
    `
    _, err = db.Exec(s)
    if (err != nil) {
        fmt.Println(err)
    }

    s = `
        CREATE TABLE IF NOT EXISTS user_detail (
            uid  int  PRIMARY KEY,
            site TEXT NOT NULL
        )
    `
    _, err = db.Exec(s)
    if (err != nil) {
        fmt.Println(err)
    }
    fmt.Println()
}

// 数据写入
func insert() {
    fmt.Println("insert:")
    r, err := db.Insert("user", &gdb.Map {
        "uid" : 1,
        "name": "john",
    })
    if (err == nil) {
        uid, err2 := r.LastInsertId()
        if err2 == nil {
            r, err = db.Insert("user_detail", &gdb.Map {
                "uid"  : string(uid),
                "site" : "http://johng.cn",
            })
            if err == nil {
                fmt.Printf("uid: %d\n", uid)
            } else {
                fmt.Println(err)
            }
        } else {
            fmt.Println(err2)
        }
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}


// 基本sql查询
func query() {
    fmt.Println("query:")
    list, err := db.GetAll("select * from \"user\"")
    if err == nil {
        fmt.Println(list)
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// replace into
func replace() {
    fmt.Println("replace:")
    r, err := db.Save("user", &gdb.Map {
        "uid": "1",
        "name": "john",
    })
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 数据保存
func save() {
    fmt.Println("save:")
    r, err := db.Save("user", &gdb.Map {
        "uid"  : "1",
        "name" : "john",
    })
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 批量写入
func batchInsert() {
    fmt.Println("batchInsert:")
    err := db.BatchInsert("user", &gdb.List {
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
    }, 10)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println()
}

// 数据更新
func update1() {
    fmt.Println("update1:")
    r, err := db.Update("user", &gdb.Map {"name": "john1"}, "uid=?", 1)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 数据更新
func update2() {
    fmt.Println("update2:")
    r, err := db.Update("user", "name='john2'", "uid=1")
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 数据更新
func update3() {
    fmt.Println("update3:")
    r, err := db.Update("user", "name=?", "uid=?", "john2", 1)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}


// 链式查询操作
func linkopSelect() {
    fmt.Println("linkopSelect:")
    r, err := db.Table("user u").
        LeftJoin("user_detail ud", "u.uid=ud.uid").
        Fields("u.*, ud.site").
        Condition("u.uid > ?", 1).
        Limit(0, 2).Select()
    if (err == nil) {
        fmt.Println(r)
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 错误操作
func linkopUpdate1() {
    fmt.Println("linkopUpdate1:")
    r, err := db.Table("henghe_setting").Update()
    if (err == nil) {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 通过Map指针方式传参方式
func linkopUpdate2() {
    fmt.Println("linkopUpdate2:")
    r, err := db.Table("user").Data(&gdb.Map{"name" : "john2"}).Condition("name=?", "john").Update()
    if (err == nil) {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 通过字符串方式传参
func linkopUpdate3() {
    fmt.Println("linkopUpdate3:")
    r, err := db.Table("user").Data("name='john3'").Condition("name=?", "john2").Update()
    if (err == nil) {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 主从io复用测试，在mysql中使用 show full processlist 查看链接信息
func keepPing() {
    fmt.Println("keepPing:")
    for {
        fmt.Println("ping...")
        db.PingMaster()
        db.PingSlave()
        time.Sleep(1*time.Second)
    }
}

// 数据库单例测试，在mysql中使用 show full processlist 查看链接信息
func instance() {
    fmt.Println("instance:")
    db1, _ := gdb.Instance()
    db2, _ := gdb.Instance()
    db3, _ := gdb.Instance()
    for {
        fmt.Println("ping...")
        db1.PingMaster()
        db1.PingSlave()
        db2.PingMaster()
        db2.PingSlave()
        db3.PingMaster()
        db3.PingSlave()
        time.Sleep(1*time.Second)
    }
}


func main() {
    create()
    create()
    insert()
    query()
    replace()
    save()
    batchInsert()
    update1()
    update2()
    update3()
    linkopSelect()
    linkopUpdate1()
    linkopUpdate2()
    linkopUpdate3()
}
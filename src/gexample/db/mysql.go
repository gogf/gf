package main

import (
    "fmt"
    "time"
    "strconv"
    "g/db/gdb"
)
var db = gdb.New(gdb.Config{
    Host : "127.0.0.1",
    Port : "3306",
    User : "root",
    Pass : "123456",
    Name : "test",
})

// 创建测试数据库
func create() {
    fmt.Println("create:")
    _, err := db.Exec("CREATE DATABASE IF NOT EXISTS test")
    if (err != nil) {
        fmt.Println(err)
    }

    s := `
        CREATE TABLE IF NOT EXISTS user (
            uid  INT(10) UNSIGNED AUTO_INCREMENT,
            name VARCHAR(45),
            PRIMARY KEY (uid)
        )
        ENGINE = InnoDB
        DEFAULT CHARACTER SET = utf8
    `
    _, err = db.Exec(s)
    if (err != nil) {
        fmt.Println(err)
    }

    s = `
        CREATE TABLE IF NOT EXISTS user_detail (
            uid   INT(10) UNSIGNED AUTO_INCREMENT,
            site  VARCHAR(255),
            PRIMARY KEY (uid)
        )
        ENGINE = InnoDB
        DEFAULT CHARACTER SET = utf8
    `
    _, err = db.Exec(s)
    if (err != nil) {
        fmt.Println(err)
    }
}

// 数据写入
func insert() {
    fmt.Println("insert1:")
    r, err := db.Insert("user", &gdb.DataMap {
        "name": "john",
    })
    if (err == nil) {
        uid, err2 := r.LastInsertId()
        if err2 == nil {
            r, err = db.Insert("user_detail", &gdb.DataMap {
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
}


// 基本sql查询
func query() {
    list, err := db.GetAll("select * from user limit 2")
    if err == nil {
        fmt.Println(list)
    } else {
        fmt.Println(err)
    }
}

// replace into
func replace() {
    r, err := db.Save("user", &gdb.DataMap {
        "uid": "1",
        "name": "john",
    })
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}

// 数据保存
func save() {
    r, err := db.Save("user", &gdb.DataMap {
        "uid"  : "1",
        "name" : "john",
    })
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}

// 批量写入
func batchInsert() {
    err := db.BatchInsert("user", &gdb.DataList {
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
        {"name": "john_" + strconv.FormatInt(time.Now().UnixNano(), 10)},
    }, 10)
    if err != nil {
        fmt.Println(err)
    }
}

// 数据更新
func update1() {
    r, err := db.Update("user", &gdb.DataMap {"name": "john1"}, "uid=?", 1)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}

// 数据更新
func update2() {
    r, err := db.Update("user", "name='john2'", "uid=1")
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}

// 数据更新
func update3() {
    r, err := db.Update("user", "name=?", "uid=?", "john2", 1)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}


// 链式查询操作
func linkopSelect() {
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
}

// 错误操作
func linkopUpdate1() {
    r, err := db.Table("henghe_setting").Update()
    if (err == nil) {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}

// 通过DataMap指针方式传参方式
func linkopUpdate2() {
    r, err := db.Table("user").Data(&gdb.DataMap{"name" : "john2"}).Condition("name=?", "john").Update()
    if (err == nil) {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}

// 通过字符串方式传参
func linkopUpdate3() {
    r, err := db.Table("user").Data("name='john3'").Condition("name=?", "john2").Update()
    if (err == nil) {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
}


func main() {
    create()
    create()
    //insert()
    //query()
    //replace()
    //save()
    //batchInsert()
    //update1()
    //update2()
    //update3()
    //linkopSelect()
    linkopUpdate1()
    //linkopUpdate2()
    //linkopUpdate3()

}
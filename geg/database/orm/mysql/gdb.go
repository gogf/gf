package main

import (
    "fmt"
    "time"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g"
)

// 本文件用于gf框架的mysql数据库操作示例，不作为单元测试使用

var db gdb.DB

// 初始化配置及创建数据库
func init () {
    gdb.AddDefaultConfigNode(gdb.ConfigNode {
       Host    : "127.0.0.1",
       Port    : "3306",
       User    : "root",
       Pass    : "12345678",
       Name    : "test",
       Type    : "mysql",
       Role    : "master",
       Charset : "utf8",
    })
    db, _ = gdb.New()

    //gins.Config().SetPath("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame")
    //db = g.Database()

    //gdb.SetConfig(gdb.ConfigNode {
    //    Host : "127.0.0.1",
    //    Port : 3306,
    //    User : "root",
    //    Pass : "123456",
    //    Name : "test",
    //    Type : "mysql",
    //})
    //db, _ = gdb.Instance()

    //gdb.SetConfig(gdb.Config {
    //    "default" : gdb.ConfigGroup {
    //        gdb.ConfigNode {
    //            Host     : "127.0.0.1",
    //            Port     : "3306",
    //            User     : "root",
    //            Pass     : "123456",
    //            Name     : "test",
    //            Type     : "mysql",
    //            Role     : "master",
    //            Priority : 100,
    //        },
    //        gdb.ConfigNode {
    //            Host     : "127.0.0.2",
    //            Port     : "3306",
    //            User     : "root",
    //            Pass     : "123456",
    //            Name     : "test",
    //            Type     : "mysql",
    //            Role     : "master",
    //            Priority : 100,
    //        },
    //        gdb.ConfigNode {
    //            Host     : "127.0.0.3",
    //            Port     : "3306",
    //            User     : "root",
    //            Pass     : "123456",
    //            Name     : "test",
    //            Type     : "mysql",
    //            Role     : "master",
    //            Priority : 100,
    //        },
    //        gdb.ConfigNode {
    //            Host     : "127.0.0.4",
    //            Port     : "3306",
    //            User     : "root",
    //            Pass     : "123456",
    //            Name     : "test",
    //            Type     : "mysql",
    //            Role     : "master",
    //            Priority : 100,
    //        },
    //    },
    //})
    //db, _ = gdb.Instance()
}



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
    fmt.Println()
}

// 数据写入
func insert() {
    fmt.Println("insert:")
    r, err := db.Insert("user", gdb.Map {
        "name": "john",
    })
    if err == nil {
        uid, err2 := r.LastInsertId()
        if err2 == nil {
            r, err = db.Insert("user_detail", gdb.Map {
                "uid"  : uid,
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
    list, err := db.GetAll("select * from user limit 2")
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
    r, err := db.Save("user", gdb.Map {
        "uid"  :  1,
        "name" : "john",
    })
    if err == nil {
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
    r, err := db.Save("user", gdb.Map {
        "uid"  : 1,
        "name" : "john",
    })
    if err == nil {
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
    _, err := db.BatchInsert("user", gdb.List {
        {"name": "john_1"},
        {"name": "john_2"},
        {"name": "john_3"},
        {"name": "john_4"},
    }, 10)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println()
}

// 数据更新
func update1() {
    fmt.Println("update1:")
    r, err := db.Update("user", gdb.Map {"name": "john1"}, "uid=?", 1)
    if err == nil {
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
    r, err := db.Update("user", gdb.Map{"name" : "john6"}, "uid=?", 1)
    if err == nil {
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
    if err == nil {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式查询操作1
func linkopSelect1() {
    fmt.Println("linkopSelect1:")
    r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*, ud.site").Where("u.uid > ?", 1).Limit(0, 2).Select()
    if err == nil {
        fmt.Println(r)
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式查询操作2
func linkopSelect2() {
    fmt.Println("linkopSelect2:")
    r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("u.*,ud.site").Where("u.uid=?", 1).One()
    if err == nil {
        fmt.Println(r)
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式查询操作3
func linkopSelect3() {
    fmt.Println("linkopSelect3:")
    r, err := db.Table("user u").LeftJoin("user_detail ud", "u.uid=ud.uid").Fields("ud.site").Where("u.uid=?", 1).Value()
    if err == nil {
        fmt.Println(r.String())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式查询数量1
func linkopCount1() {
    fmt.Println("linkopCount1:")
    r, err := db.Table("user u").Fields("uid").LeftJoin("user_detail ud", "u.uid=ud.uid").Where("u.uid=?", 1).Count()
    if err == nil {
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
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 通过Map指针方式传参方式
func linkopUpdate2() {
    fmt.Println("linkopUpdate2:")
    r, err := db.Table("user").Data(gdb.Map{"name" : "john2"}).Where("name=?", "john_1").Update()
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 通过字符串方式传参
func linkopUpdate3() {
    fmt.Println("linkopUpdate3:")
    r, err := db.Table("user").Data("name='john3'").Where("name=?", "john2").Update()
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// Where条件使用Map
func linkopUpdate4() {
    fmt.Println("linkopUpdate4:")
    r, err := db.Table("user").Data(gdb.Map{"name" : "john11111"}).Where(g.Map{"uid" : 1}).Update()
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式批量写入
func linkopBatchInsert1() {
    fmt.Println("linkopBatchInsert1:")
    r, err := db.Table("user").Data(gdb.List{
        {"name": "john_1"},
        {"name": "john_2"},
        {"name": "john_3"},
        {"name": "john_4"},
    }).Insert()
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式批量写入，指定每批次写入的条数
func linkopBatchInsert2() {
    fmt.Println("linkopBatchInsert2:")
    r, err := db.Table("user").Data(gdb.List{
        {"name": "john_1"},
        {"name": "john_2"},
        {"name": "john_3"},
        {"name": "john_4"},
    }).Batch(2).Insert()
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 链式批量保存
func linkopBatchSave() {
    fmt.Println("linkopBatchSave:")
    r, err := db.Table("user").Data(gdb.List{
        {"uid":1, "name": "john_1"},
        {"uid":2, "name": "john_2"},
        {"uid":3, "name": "john_3"},
        {"uid":4, "name": "john_4"},
    }).Save()
    if err == nil {
        fmt.Println(r.RowsAffected())
    } else {
        fmt.Println(err)
    }
    fmt.Println()
}

// 事务操作示例1
func transaction1() {
    fmt.Println("transaction1:")
    if tx, err := db.Begin(); err == nil {
        r, err := tx.Save("user", gdb.Map{
            "uid"  :  1,
            "name" : "john",
        })
        tx.Rollback()
        fmt.Println(r, err)
    }
    fmt.Println()
}

// 事务操作示例2
func transaction2() {
    fmt.Println("transaction2:")
    if tx, err := db.Begin(); err == nil {
        r, err := tx.Table("user").Data(gdb.Map{"uid":1, "name": "john_1"}).Save()
        tx.Commit()
        fmt.Println(r, err)
    }
    fmt.Println()
}

// 主从io复用测试，在mysql中使用 show full processlist 查看链接信息
func keepPing() {
    fmt.Println("keepPing:")
    for {
        fmt.Println("ping...")
        err := db.PingMaster()
        if err != nil {
            fmt.Println(err)
            return
        }
        err  = db.PingSlave()
        if err != nil {
            fmt.Println(err)
            return
        }
        time.Sleep(1*time.Second)
    }
}

// like语句查询
func likeQuery() {
    fmt.Println("likeQuery:")
    if r, err := db.Table("user").Where("name like ?", "%john%").Select(); err == nil {
        fmt.Println(r)
    } else {
        fmt.Println(err)
    }
}


// mapToStruct
func mapToStruct() {
    type User struct {
        Uid  int
        Name string
    }
    fmt.Println("mapToStruct:")
    if r, err := db.Table("user").Where("uid=?", 1).One(); err == nil {
        u := User{}
        if err := r.ToStruct(&u); err == nil {
            fmt.Println(r)
            fmt.Println(u)
        } else {
            fmt.Println(err)
        }
    } else {
        fmt.Println(err)
    }
}

// getQueriedSqls
func getQueriedSqls() {
    for k, v := range db.GetQueriedSqls() {
        fmt.Println(k, ":")
        fmt.Println("Sql  :", v.Sql)
        fmt.Println("Args :", v.Args)
        fmt.Println("Error:", v.Error)
        fmt.Println("Func :", v.Func)
    }
}

func main() {
    //db.SetDebug(true)
    //r, err := db.Table("test").Where("id=1").One()
    //fmt.Println(r["datetime"])
    //fmt.Println(r["datetime"].Time().Date())
    //fmt.Println(err)
    //create()
    //create()
    //insert()
    //query()
    //replace()
    //save()
    //batchInsert()
    //update1()
    //update2()
    //update3()
    linkopSelect1()
    //linkopSelect2()
    //linkopSelect3()
    //linkopCount1()
    //linkopUpdate1()
    //linkopUpdate2()
    //linkopUpdate3()
    //linkopUpdate4()
    //
    //transaction1()
    //transaction2()
    //
    //keepPing()
    //likeQuery()
    //mapToStruct()
    //getQueriedSqls()
}
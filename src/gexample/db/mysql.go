package main

import (
    "fmt"
    "g"
    "time"
    "strconv"
)
var db = g.Db.New(g.GDbConfig{
    Host : "127.0.0.1",
    Port : "3306",
    User : "root",
    Pass : "8692651",
    Name : "hhzl_gdg",
})

func query1() {
    _, err := db.Query("insert into henghe_setting(k,v) values('k', 'v')")
    fmt.Println(err)
}

func exec() {
    r, err := db.Exec("insert into henghe_setting(k,v) values('k', 'v')")
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func insert1() {
    r, err := db.Insert("henghe_setting", map[string]string {
        "k": strconv.FormatInt(time.Now().UnixNano(), 10),
        "v": "v",
    })
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func insert2() {
    r, err := db.Insert("henghe_setting", map[string]string {
        "k": "k",
        "v": "v",
    })
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func replace() {
    r, err := db.Save("henghe_setting", map[string]string {
        "k": "k",
        "v": "v3",
    })
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func save() {
    r, err := db.Save("henghe_setting", map[string]string {
        "k": "k",
        "v": "v4",
    })
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func batchInsert() {
    err := db.BatchInsert("henghe_setting", []map[string]string {
        {
            "k": strconv.FormatInt(time.Now().UnixNano(), 10),
            "v": "v",
        },
        {
            "k": strconv.FormatInt(time.Now().UnixNano(), 10),
            "v": "v",
        },
        {
            "k": strconv.FormatInt(time.Now().UnixNano(), 10),
            "v": "v",
        },
    }, 10)
    fmt.Println(err)
}

func update1() {
    r, err := db.Update("henghe_setting", map[string]string {"v": "v3"}, "k=?", "1499322119238299405")
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func update2() {
    r, err := db.Update("henghe_setting", "`v`='v4'", "`k`='1499322119238299405'")
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func update3() {
    r, err := db.Update("henghe_setting", "v=?", "k=?", "v5", "1499322119238299405")
    fmt.Println(err)
    if (err == nil) {
        fmt.Println(r.LastInsertId())
        fmt.Println(r.RowsAffected())
    }
}

func main() {
    query1()
    exec()
    //insert1()
    //insert2()
    //replace()
    //save()
    //batchInsert()
    //update1()
    //update2()
    //update3()
}
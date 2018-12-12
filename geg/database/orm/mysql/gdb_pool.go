package main

import (
    "gitee.com/johng/gf/g/database/gdb"
    "time"
)

func main() {
    gdb.AddDefaultConfigNode(gdb.ConfigNode {
        Host             : "127.0.0.1",
        Port             : "3306",
        User             : "root",
        Pass             : "12345678",
        Name             : "test",
        Type             : "mysql",
        Role             : "master",
        Charset          : "utf8",
        MaxIdleConnCount : 10,
        MaxOpenConnCount : 10,
        MaxConnLifetime  : 10,
    })
    db, err := gdb.New()
    if err != nil {
        panic(err)
    }
    // 开启调试模式，以便于记录所有执行的SQL
    db.SetDebug(true)

    for {
        for i := 0; i < 10; i++ {
            go db.Table("user").All()
        }
        time.Sleep(time.Second)
    }

}
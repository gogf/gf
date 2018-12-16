package main

import (
    "gitee.com/johng/gf/g"
    "time"
)

func main() {
    db := g.DB()
    db.SetMaxIdleConns(10)
    db.SetMaxOpenConns(10)
    db.SetConnMaxLifetime(10)

    // 开启调试模式，以便于记录所有执行的SQL
    db.SetDebug(true)

    for {
        for i := 0; i < 10; i++ {
            go db.Table("user").All()
        }
        time.Sleep(time.Second)
    }

}
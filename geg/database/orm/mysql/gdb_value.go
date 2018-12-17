package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    db := g.DB()
    // 开启调试模式，以便于记录所有执行的SQL
    db.SetDebug(true)

    r, _ := db.Table("test").Where("id IN (?)", []interface{}{1, 2}).All()
    if r != nil {
        fmt.Println(r.ToList())
    }
}
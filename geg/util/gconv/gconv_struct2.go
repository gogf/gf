package main

import (
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g"
    "fmt"
)

// 演示slice类型属性的赋值
func main() {
    type User struct {
        Scores []int
    }

    user   := new(User)
    scores := []int{99, 100, 60, 140}

    err := gconv.MapToStruct(g.Map{
        "Scores" : scores,
    }, user)
    if err != nil {
        fmt.Println(err)
    } else {
        g.Dump(user)
    }
}
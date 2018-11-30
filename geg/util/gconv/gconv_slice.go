package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
)

// struct转slice
func main() {
    type User struct {
        Uid  int
        Name string
    }
    // 对象
    fmt.Println(gconv.Interfaces(User{
        Uid      : 1,
        Name     : "john",
    }))
    // 指针
    fmt.Println(gconv.Interfaces(&User{
        Uid      : 1,
        Name     : "john",
    }))
}
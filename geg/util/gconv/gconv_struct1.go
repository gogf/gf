package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gconv"
)

type User struct {
    Uid   int
    Name  string
    Pass1 string `gconv:"password1"`
    Pass2 string `gconv:"password2"`
}

func main() {
    user    := (*User)(nil)

    // 使用默认映射规则绑定属性值到对象
    user     = new(User)
    params1 := g.Map{
        "uid"   : 1,
        "Name"  : "john",
        "PASS1" : "123",
        "PASS2" : "456",
    }
    if err := gconv.Struct(params1, user); err == nil {
        fmt.Println(user)
    }

    // 使用struct tag映射绑定属性值到对象
    user     = new(User)
    params2 := g.Map {
        "uid"       : 2,
        "name"      : "smith",
        "password1" : "111",
        "password2" : "222",
    }
    if err := gconv.Struct(params2, user); err == nil {
        fmt.Println(user)
    }
}
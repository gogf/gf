package main

import (
    "fmt"
)

// 自定义的struct
type User struct {
    Uid  int
    Name string
}

// map[uid:1 name:john3 email: type:1]
// {1 john3}
func main() {
    if r, err := db.Table("user").Where("uid=?", 1).One(); err == nil {
        u := User{}
        if err := r.ToStruct(&u); err == nil {
            fmt.Println(u)
        } else {
            fmt.Println(err)
        }
    } else {
        fmt.Println(err)
    }
}
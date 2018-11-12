package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gvalid"
)


// same校验
func main() {
    type User struct {
        Password        string `gvalid:"password@password"`
        ConfirmPassword string `gvalid:"confirm_password@password|same:password#|密码与确认密码不一致"`
    }

    user := &User{
        Password        : "123456",
        ConfirmPassword : "",
    }

    g.Dump(gvalid.CheckStruct(user, nil))
}

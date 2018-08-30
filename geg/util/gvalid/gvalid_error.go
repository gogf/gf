package main

import (
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/util/gvalid"
)


// same校验
func main() {
    type User struct {
        Password        string `gvalid:"password@password"`
        ConfiemPassword string `gvalid:"confirm_password@password|same:password#|密码与确认密码不一致"`
    }

    user := &User{
        Password        : "123456",
        ConfiemPassword : "",
    }

    err := gvalid.CheckStruct(user, nil)
    gutil.Dump(err)
    gutil.Dump(err.String())
    gutil.Dump(err.Strings())
    gutil.Dump(err.FirstItem())
    gutil.Dump(err.FirstRule())
    gutil.Dump(err.FirstString())
}

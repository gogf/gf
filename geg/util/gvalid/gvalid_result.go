package main

import (
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/util/gvalid"
)

func main() {
    type User struct {
        Name  string `gvalid:"name     @required|length:6,30#请输入用户名称|用户名称长度非法"`
        Pass1 string `gvalid:"password1@required|password3"`
        Pass2 string `gvalid:"password2@required|password3|same:password1#||两次密码不一致，请重新输入"`
    }

    user := &User{
        Name : "john",
        Pass1: "Abc123!@#",
        Pass2: "123",
    }

    err := gvalid.CheckStruct(user, nil)
    gutil.Dump(err)
    gutil.Dump(err.String())
    gutil.Dump(err.FirstString())
}

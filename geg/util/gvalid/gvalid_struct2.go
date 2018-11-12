package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gvalid"
)


// string默认值校验
func main() {
    type User struct {
        Uid string `gvalid:"uid@integer"`
    }

    user := &User{}

    g.Dump(gvalid.CheckStruct(user, nil))
}

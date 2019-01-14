// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
)


func Test_Struct_Basic(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Uid      int
            Name     string
            Site_Url string
            NickName string
            Pass1    string `gconv:"password1"`
            Pass2    string `gconv:"password2"`
        }
        user    := (*User)(nil)
        // 使用默认映射规则绑定属性值到对象
        user     = new(User)
        params1 := g.Map{
            "uid"       : 1,
            "Name"      : "john",
            "siteurl"   : "https://gfer.me",
            "nick_name" : "johng",
            "PASS1"     : "123",
            "PASS2"     : "456",
        }
        if err := gconv.Struct(params1, user); err != nil {
            gtest.Error(err)
        }
        gtest.Assert(user, &User{
            Uid      : 1,
            Name     : "john",
            Site_Url : "https://gfer.me",
            NickName : "johng",
            Pass1    : "123",
            Pass2    : "456",
        })

        // 使用struct tag映射绑定属性值到对象
        user     = new(User)
        params2 := g.Map {
            "uid"       : 2,
            "name"      : "smith",
            "site-url"  : "https://gfer.me",
            "nick name" : "johng",
            "password1" : "111",
            "password2" : "222",
        }
        if err := gconv.Struct(params2, user); err != nil {
            gtest.Error(err)
        }
        gtest.Assert(user, &User{
            Uid      : 2,
            Name     : "smith",
            Site_Url : "https://gfer.me",
            NickName : "johng",
            Pass1    : "111",
            Pass2    : "222",
        })
    })
}

// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gconv"
	"testing"
)


func Test_Map_Basic(t *testing.T) {
    gtest.Case(t, func() {
        m1 := map[string]string{
            "k" : "v",
        }
        m2 := map[int]string{
            3 : "v",
        }
        m3 := map[float64]float32{
           1.22 : 3.1,
        }
        gtest.Assert(gconv.Map(m1), g.Map{
            "k" : "v",
        })
        gtest.Assert(gconv.Map(m2), g.Map{
            "3" : "v",
        })
        gtest.Assert(gconv.Map(m3), g.Map{
            "1.22" : "3.1",
        })
    })
}

func Test_Map_StructWithGconvTag(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Uid      int
            Name     string
            SiteUrl  string `gconv:"-"`
            NickName string `gconv:"nickname, omitempty"`
            Pass1    string `gconv:"password1"`
            Pass2    string `gconv:"password2"`
        }
        user1 := User{
            Uid      : 100,
            Name     : "john",
            SiteUrl  : "https://goframe.org",
            Pass1    : "123",
            Pass2    : "456",
        }
        user2 := &user1
        map1  := gconv.Map(user1)
        map2  := gconv.Map(user2)
        gtest.Assert(map1["Uid"],       100)
        gtest.Assert(map1["Name"],      "john")
        gtest.Assert(map1["SiteUrl"],   nil)
        gtest.Assert(map1["NickName"],  nil)
        gtest.Assert(map1["nickname"],  nil)
        gtest.Assert(map1["password1"], "123")
        gtest.Assert(map1["password2"], "456")

        gtest.Assert(map2["Uid"],       100)
        gtest.Assert(map2["Name"],      "john")
        gtest.Assert(map2["SiteUrl"],   nil)
        gtest.Assert(map2["NickName"],  nil)
        gtest.Assert(map2["nickname"],  nil)
        gtest.Assert(map2["password1"], "123")
        gtest.Assert(map2["password2"], "456")
    })
}

func Test_Map_StructWithJsonTag(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Uid      int
            Name     string
            SiteUrl  string `json:"-"`
            NickName string `json:"nickname, omitempty"`
            Pass1    string `json:"password1"`
            Pass2    string `json:"password2"`
        }
        user1 := User{
            Uid      : 100,
            Name     : "john",
            SiteUrl  : "https://goframe.org",
            Pass1    : "123",
            Pass2    : "456",
        }
        user2 := &user1
        map1  := gconv.Map(user1)
        map2  := gconv.Map(user2)
        gtest.Assert(map1["Uid"],       100)
        gtest.Assert(map1["Name"],      "john")
        gtest.Assert(map1["SiteUrl"],   nil)
        gtest.Assert(map1["NickName"],  nil)
        gtest.Assert(map1["nickname"],  nil)
        gtest.Assert(map1["password1"], "123")
        gtest.Assert(map1["password2"], "456")

        gtest.Assert(map2["Uid"],       100)
        gtest.Assert(map2["Name"],      "john")
        gtest.Assert(map2["SiteUrl"],   nil)
        gtest.Assert(map2["NickName"],  nil)
        gtest.Assert(map2["nickname"],  nil)
        gtest.Assert(map2["password1"], "123")
        gtest.Assert(map2["password2"], "456")
    })
}

func Test_Map_PrivateAttribute(t *testing.T) {
    type User struct {
        Id   int
        name string
    }
    gtest.Case(t, func() {
        user := &User{1, "john"}
        gtest.Assert(gconv.Map(user), g.Map{"Id" : 1})
    })
}

func Test_Map_StructInherit(t *testing.T) {
	gtest.Case(t, func() {
		type Ids struct {
			Id         int    `json:"id"`
			Uid        int    `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport   string `json:"passport"`
			Password   string `json:"password"`
			Nickname   string `json:"nickname"`
		}
		user := new(User)
		user.Id         = 100
		user.Nickname   = "john"
		user.CreateTime = "2019"
		m := gconv.MapDeep(user)
		gtest.Assert(m["id"],          user.Id)
		gtest.Assert(m["nickname"],    user.Nickname)
		gtest.Assert(m["create_time"], user.CreateTime)
	})
}
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

func Test_Struct_Basic1(t *testing.T) {
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
            "siteurl"   : "https://goframe.org",
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
            Site_Url : "https://goframe.org",
            NickName : "johng",
            Pass1    : "123",
            Pass2    : "456",
        })

        // 使用struct tag映射绑定属性值到对象
        user     = new(User)
        params2 := g.Map {
            "uid"       : 2,
            "name"      : "smith",
            "site-url"  : "https://goframe.org",
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
            Site_Url : "https://goframe.org",
            NickName : "johng",
            Pass1    : "111",
            Pass2    : "222",
        })
    })
}

// 使用默认映射规则绑定属性值到对象
func Test_Struct_Basic2(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Uid     int
            Name    string
            SiteUrl string
            Pass1   string
            Pass2   string

        }
        user    := new(User)
        params  := g.Map {
            "uid"      : 1,
            "Name"     : "john",
            "site_url" : "https://goframe.org",
            "PASS1"    : "123",
            "PASS2"    : "456",
        }
        if err := gconv.Struct(params, user); err != nil {
            gtest.Error(err)
        }
        gtest.Assert(user, &User{
            Uid      : 1,
            Name     : "john",
            SiteUrl  : "https://goframe.org",
            Pass1    : "123",
            Pass2    : "456",
        })
    })
}

// slice类型属性的赋值
func Test_Struct_Attr_Slice(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Scores []int
        }
        scores := []interface{}{99, 100, 60, 140}
        user   := new(User)
        if err := gconv.Struct(g.Map{"Scores" : scores}, user); err != nil {
            gtest.Error(err)
        } else {
            gtest.Assert(user, &User{
                Scores : []int{99, 100, 60, 140},
            })
        }
    })
}

// 属性为struct对象
func Test_Struct_Attr_Struct(t *testing.T) {
    gtest.Case(t, func() {
        type Score struct {
            Name   string
            Result int
        }
        type User struct {
            Scores Score
        }

        user   := new(User)
        scores := map[string]interface{}{
            "Scores" : map[string]interface{}{
                "Name"   : "john",
                "Result" : 100,
            },
        }

        // 嵌套struct转换
        if err := gconv.Struct(scores, user); err != nil {
            gtest.Error(err)
        } else {
            gtest.Assert(user, &User{
                Scores : Score {
                    Name   : "john",
                    Result : 100,
                },
            })
        }
    })
}

// 属性为struct对象指针
func Test_Struct_Attr_Struct_Ptr(t *testing.T) {
    gtest.Case(t, func() {
        type Score struct {
            Name   string
            Result int
        }
        type User struct {
            Scores *Score
        }

        user   := new(User)
        scores := map[string]interface{}{
            "Scores" : map[string]interface{}{
                "Name"   : "john",
                "Result" : 100,
            },
        }

        // 嵌套struct转换
        if err := gconv.Struct(scores, user); err != nil {
            gtest.Error(err)
        } else {
            gtest.Assert(user.Scores, &Score {
                Name   : "john",
                Result : 100,
            })
        }
    })
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice1(t *testing.T) {
    gtest.Case(t, func() {
        type Score struct {
            Name   string
            Result int
        }
        type User struct {
            Scores []Score
        }

        user   := new(User)
        scores := map[string]interface{}{
            "Scores" : map[string]interface{}{
                "Name"   : "john",
                "Result" : 100,
            },
        }

        // 嵌套struct转换，属性为slice类型，数值为map类型
        if err := gconv.Struct(scores, user); err != nil {
            gtest.Error(err)
        } else {
            gtest.Assert(user.Scores, []Score {
                {
                    Name   : "john",
                    Result : 100,
                },
            })
        }
    })
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice2(t *testing.T) {
    gtest.Case(t, func() {
        type Score struct {
            Name   string
            Result int
        }
        type User struct {
            Scores []Score
        }

        user   := new(User)
        scores := map[string]interface{}{
            "Scores" : []interface{}{
                map[string]interface{}{
                    "Name"   : "john",
                    "Result" : 100,
                },
                map[string]interface{}{
                    "Name"   : "smith",
                    "Result" : 60,
                },
            },
        }

        // 嵌套struct转换，属性为slice类型，数值为slice map类型
        if err := gconv.Struct(scores, user); err != nil {
            gtest.Error(err)
        } else {
            gtest.Assert(user.Scores, []Score {
                {
                    Name   : "john",
                    Result : 100,
                },
                {
                    Name   : "smith",
                    Result : 60,
                },
            })
        }
    })
}

// 属性为struct对象slice ptr
func Test_Struct_Attr_Struct_Slice_Ptr(t *testing.T) {
    gtest.Case(t, func() {
        type Score struct {
            Name   string
            Result int
        }
        type User struct {
            Scores []*Score
        }

        user   := new(User)
        scores := map[string]interface{}{
            "Scores" : []interface{}{
                map[string]interface{}{
                    "Name"   : "john",
                    "Result" : 100,
                },
                map[string]interface{}{
                    "Name"   : "smith",
                    "Result" : 60,
                },
            },
        }

        // 嵌套struct转换，属性为slice类型，数值为slice map类型
        if err := gconv.Struct(scores, user); err != nil {
            gtest.Error(err)
        } else {
            gtest.Assert(len(user.Scores), 2)
            gtest.Assert(user.Scores[0], &Score {
                Name   : "john",
                Result : 100,
            })
            gtest.Assert(user.Scores[1], &Score {
                Name   : "smith",
                Result : 60,
            })
        }
    })
}

func Test_Struct_PrivateAttribute(t *testing.T) {
    type User struct {
        Id   int
        name string
    }
    gtest.Case(t, func() {
        user := new(User)
        err  := gconv.Struct(g.Map{"id" : 1, "name" : "john"}, user)
        gtest.Assert(err,       nil)
        gtest.Assert(user.Id,   1)
        gtest.Assert(user.name, "")
    })
}

func Test_Struct_Deep(t *testing.T) {
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
		data := g.Map{
			"id"          : 100,
			"uid"         : 101,
			"passport"    : "t1",
			"password"    : "123456",
			"nickname"    : "T1",
			"create_time" : "2019",
		}
		user := new(User)
		gconv.StructDeep(data, user)
		gtest.Assert(user.Id,         100)
		gtest.Assert(user.Uid,        101)
		gtest.Assert(user.Nickname,   "T1")
		gtest.Assert(user.CreateTime, "2019")
	})
}
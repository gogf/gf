// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
    "fmt"
    "github.com/gogf/gf/g/frame/gins"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)

func Test_SetGet(t *testing.T) {
    gtest.Case(t, func() {
        gins.Set("test-user", 1)
        gtest.Assert(gins.Get("test-user"),   1)
        gtest.Assert(gins.Get("none-exists"), nil)
    })
    gtest.Case(t, func() {
        gtest.Assert(gins.GetOrSet("test-1", 1), 1)
        gtest.Assert(gins.Get("test-1"), 1)
    })
    gtest.Case(t, func() {
        gtest.Assert(gins.GetOrSetFunc("test-2", func() interface{} {
            return 2
        }), 2)
        gtest.Assert(gins.Get("test-2"), 2)
    })
    gtest.Case(t, func() {
        gtest.Assert(gins.GetOrSetFuncLock("test-3", func() interface{} {
            return 3
        }), 3)
        gtest.Assert(gins.Get("test-3"), 3)
    })
    gtest.Case(t, func() {
        gtest.Assert(gins.SetIfNotExist("test-4", 4), true)
        gtest.Assert(gins.Get("test-4"), 4)
        gtest.Assert(gins.SetIfNotExist("test-4", 5), false)
        gtest.Assert(gins.Get("test-4"), 4)
    })
}

func Test_View(t *testing.T) {
    gtest.Case(t, func() {
        gtest.AssertNE(gins.View(), nil)
        b, e := gins.View().ParseContent(`{{"1540822968" | date "Y-m-d H:i:s"}}`, nil)
        gtest.Assert(e,         nil)
        gtest.Assert(string(b), "2018-10-29 22:22:48")
    })
    gtest.Case(t, func() {
        tpl := "t.tpl"
        err := gfile.PutContents(tpl, `{{"1540822968" | date "Y-m-d H:i:s"}}`)
        gtest.Assert(err, nil)
        defer gfile.Remove(tpl)

        b, e := gins.View().Parse("t.tpl", nil)
        gtest.Assert(e,         nil)
        gtest.Assert(string(b), "2018-10-29 22:22:48")
    })
    gtest.Case(t, func() {
        path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.Nanosecond())
        tpl  := fmt.Sprintf(`%s/%s`, path, "t.tpl")
        err  := gfile.PutContents(tpl, `{{"1540822968" | date "Y-m-d H:i:s"}}`)
        gtest.Assert(err, nil)
        defer gfile.Remove(tpl)
        err = gins.View().AddPath(path)
        gtest.Assert(err, nil)

        b, e := gins.View().Parse("t.tpl", nil)
        gtest.Assert(e,         nil)
        gtest.Assert(string(b), "2018-10-29 22:22:48")
    })
}

func Test_Config(t *testing.T) {
    config := `
# 模板引擎目录
viewpath = "/home/www/templates/"
test = "v=1"
# MySQL数据库配置
[database]
    [[database.default]]
        host     = "127.0.0.1"
        port     = "3306"
        user     = "root"
        pass     = ""
        name     = "test"
        type     = "mysql"
        role     = "master"
        charset  = "utf8"
        priority = "1"
    [[database.default]]
        host     = "127.0.0.1"
        port     = "3306"
        user     = "root"
        pass     = "8692651"
        name     = "test"
        type     = "mysql"
        role     = "master"
        charset  = "utf8"
        priority = "1"
# Redis数据库配置
[redis]
    disk  = "127.0.0.1:6379,0"
    cache = "127.0.0.1:6379,1"
`
    gtest.Case(t, func() {
        gtest.AssertNE(gins.Config(), nil)
    })
    gtest.Case(t, func() {
        path := "config.toml"
        err  := gfile.PutContents(path, config)
        gtest.Assert(err, nil)
        defer gfile.Remove(path)

        //fmt.Println(os.Getwd())
        //fmt.Println(gfile.Pwd())
        //fmt.Println(gfile.ScanDir(".", "*"))

        gtest.Assert(gins.Config().Get("test"), "v=1")
        gtest.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
        gtest.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
    })
    gtest.Case(t, func() {
        path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.Nanosecond())
        file := fmt.Sprintf(`%s/%s`, path, "config.toml")
        err  := gfile.PutContents(file, config)
        gtest.Assert(err, nil)
        defer gfile.Remove(file)
        err = gins.Config().AddPath(path)
        gtest.Assert(err, nil)

        gtest.Assert(gins.Config().Get("test"), "v=1")
        gtest.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
        gtest.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
    })
}
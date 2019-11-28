// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"testing"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_Database(t *testing.T) {
	config := `
# 模板引擎目录
viewpath = "/home/www/templates/"
test = "v=2"
# MySQL数据库配置
[database]
    [[database.default]]
        host     = "127.0.0.1"
        port     = "3306"
        user     = "root"
        pass     = "12345678"
        name     = "test"
        type     = "mysql"
        role     = "master"
		weight   = "1"
        charset  = "utf8"
    [[database.test]]
        host     = "127.0.0.1"
        port     = "3306"
        user     = "root"
        pass     = "12345678"
        name     = "test"
        type     = "mysql"
        role     = "master"
		weight   = "1"
        charset  = "utf8"
# Redis数据库配置
[redis]
    default = "127.0.0.1:6379,0"
    cache   = "127.0.0.1:6379,1"
`
	path := "config.toml"
	err := gfile.PutContents(path, config)
	gtest.Assert(err, nil)
	defer gfile.Remove(path)
	defer Config().Clear()

	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	gtest.Case(t, func() {
		//fmt.Println("gins Test_Database", Config().Get("test"))

		dbDefault := Database()
		dbTest := Database("test")
		gtest.AssertNE(dbDefault, nil)
		gtest.AssertNE(dbTest, nil)

		gtest.Assert(dbDefault.PingMaster(), nil)
		gtest.Assert(dbDefault.PingSlave(), nil)
		gtest.Assert(dbTest.PingMaster(), nil)
		gtest.Assert(dbTest.PingSlave(), nil)
	})
}

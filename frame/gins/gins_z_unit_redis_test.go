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

func Test_Redis(t *testing.T) {
	config := `
# 模板引擎目录
viewpath = "/home/www/templates/"
test = "v=3"
# MySQL数据库配置
[database]
    [[database.default]]
        host     = "127.0.0.1"
        port     = "3306"
        user     = "root"
        pass     = ""
        # pass     = "12345678"
        name     = "test"
        type     = "mysql"
        role     = "master"
        charset  = "utf8"
        priority = "1"
    [[database.test]]
        host     = "127.0.0.1"
        port     = "3306"
        user     = "root"
        pass     = ""
        # pass     = "12345678"
        name     = "test"
        type     = "mysql"
        role     = "master"
        charset  = "utf8"
        priority = "1"
# Redis数据库配置
[redis]
    default = "127.0.0.1:6379,7"
    cache   = "127.0.0.1:6379,8"
    disk    = "127.0.0.1:6379,9,?maxIdle=1&maxActive=10&idleTimeout=10&maxConnLifetime=10"
`
	path := "config.toml"
	err := gfile.PutContents(path, config)
	gtest.Assert(err, nil)
	defer gfile.Remove(path)
	defer Config().Clear()

	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	gtest.Case(t, func() {
		//fmt.Println("gins Test_Redis", Config().Get("test"))

		redisDefault := Redis()
		redisCache := Redis("cache")
		redisDisk := Redis("disk")
		gtest.AssertNE(redisDefault, nil)
		gtest.AssertNE(redisCache, nil)
		gtest.AssertNE(redisDisk, nil)

		r, err := redisDefault.Do("PING")
		gtest.Assert(err, nil)
		gtest.Assert(r, "PONG")

		r, err = redisCache.Do("PING")
		gtest.Assert(err, nil)
		gtest.Assert(r, "PONG")

		_, err = redisDisk.Do("SET", "k", "v")
		gtest.Assert(err, nil)
		r, err = redisDisk.Do("GET", "k")
		gtest.Assert(err, nil)
		gtest.Assert(r, []byte("v"))
	})
}

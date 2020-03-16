// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/os/gtime"
	"testing"
	"time"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/test/gtest"
)

func Test_Redis(t *testing.T) {
	redisContent := gfile.GetContents(gfile.Join(gdebug.TestDataPath(), "redis", "config.toml"))

	var err error
	dirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
	err = gfile.Mkdir(dirPath)
	gtest.Assert(err, nil)
	defer gfile.Remove(dirPath)

	name := "config.toml"
	err = gfile.PutContents(gfile.Join(dirPath, name), redisContent)
	gtest.Assert(err, nil)

	err = gins.Config().AddPath(dirPath)
	gtest.Assert(err, nil)

	defer gins.Config().Clear()

	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	gtest.Case(t, func() {
		//fmt.Println("gins Test_Redis", Config().Get("test"))

		redisDefault := gins.Redis()
		redisCache := gins.Redis("cache")
		redisDisk := gins.Redis("disk")
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

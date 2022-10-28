// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/gins"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_GINS_Redis(t *testing.T) {
	redisContent := gfile.GetContents(
		gtest.DataPath("redis", "config.toml"),
	)

	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.Temp(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.AssertNil(err)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), redisContent)
		t.AssertNil(err)

		err = gins.Config().GetAdapter().(*gcfg.AdapterFile).AddPath(dirPath)
		t.AssertNil(err)

		defer gins.Config().GetAdapter().(*gcfg.AdapterFile).Clear()

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		// fmt.Println("gins Test_Redis", Config().Get("test"))

		var (
			redisDefault = gins.Redis()
			redisCache   = gins.Redis("cache")
			redisDisk    = gins.Redis("disk")
		)
		t.AssertNE(redisDefault, nil)
		t.AssertNE(redisCache, nil)
		t.AssertNE(redisDisk, nil)

		r, err := redisDefault.Do(ctx, "PING")
		t.AssertNil(err)
		t.Assert(r, "PONG")

		r, err = redisCache.Do(ctx, "PING")
		t.AssertNil(err)
		t.Assert(r, "PONG")

		_, err = redisDisk.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)
		r, err = redisDisk.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))
	})
}

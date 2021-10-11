// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/gins"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gtime"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Redis(t *testing.T) {
	redisContent := gfile.GetContents(
		gdebug.TestDataPath("redis", "config.toml"),
	)

	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), redisContent)
		t.Assert(err, nil)

		err = gins.Config().GetAdapter().(*gcfg.AdapterFile).AddPath(dirPath)
		t.Assert(err, nil)

		defer gins.Config().GetAdapter().(*gcfg.AdapterFile).Clear()

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		//fmt.Println("gins Test_Redis", Config().Get("test"))

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
		t.Assert(err, nil)
		t.Assert(r, "PONG")

		_, err = redisDisk.Do(ctx, "SET", "k", "v")
		t.Assert(err, nil)
		r, err = redisDisk.Do(ctx, "GET", "k")
		t.Assert(err, nil)
		t.Assert(r, []byte("v"))
	})
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/frame/gins"
	"testing"
	"time"

	"github.com/gogf/gf/os/gcfg"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

var (
	configContent = gfile.GetContents(
		gfile.Join(gdebug.TestDataPath(), "config", "config.toml"),
	)
)

func Test_Config1(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertNE(configContent, "")
	})
	gtest.Case(t, func() {
		gtest.AssertNE(gins.Config(), nil)
	})
}

func Test_Config2(t *testing.T) {
	// relative path
	gtest.Case(t, func() {
		var err error
		dirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
		gtest.Assert(err, nil)

		err = gins.Config().AddPath(dirPath)
		gtest.Assert(err, nil)

		defer gins.Config().Clear()

		gtest.Assert(gins.Config().Get("test"), "v=1")
		gtest.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
		gtest.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
	})
	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	// relative path, config folder
	gtest.Case(t, func() {
		var err error
		dirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config/config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
		gtest.Assert(err, nil)

		err = gins.Config().AddPath(dirPath)
		gtest.Assert(err, nil)

		defer gins.Config().Clear()

		gtest.Assert(gins.Config().Get("test"), "v=1")
		gtest.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
		gtest.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)
	})
}

func Test_Config3(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Case(t, func() {
			var err error
			dirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
			err = gfile.Mkdir(dirPath)
			gtest.Assert(err, nil)
			defer gfile.Remove(dirPath)

			name := "test.toml"
			err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
			gtest.Assert(err, nil)

			err = gins.Config("test").AddPath(dirPath)
			gtest.Assert(err, nil)

			defer gins.Config("test").Clear()
			gins.Config("test").SetFileName("test.toml")

			gtest.Assert(gins.Config("test").Get("test"), "v=1")
			gtest.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
			gtest.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
		})
		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		gtest.Case(t, func() {
			var err error
			dirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
			err = gfile.Mkdir(dirPath)
			gtest.Assert(err, nil)
			defer gfile.Remove(dirPath)

			name := "config/test.toml"
			err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
			gtest.Assert(err, nil)

			err = gins.Config("test").AddPath(dirPath)
			gtest.Assert(err, nil)

			defer gins.Config("test").Clear()
			gins.Config("test").SetFileName("test.toml")

			gtest.Assert(gins.Config("test").Get("test"), "v=1")
			gtest.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
			gtest.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
		})
		// for gfsnotify callbacks to refresh cache of config file for next unit testing case.
		time.Sleep(500 * time.Millisecond)
	})
}

func Test_Config4(t *testing.T) {
	gtest.Case(t, func() {
		// absolute path
		gtest.Case(t, func() {
			path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.TimestampNano())
			file := fmt.Sprintf(`%s/%s`, path, "config.toml")
			err := gfile.PutContents(file, configContent)
			gtest.Assert(err, nil)
			defer gfile.Remove(file)
			defer gins.Config().Clear()

			gtest.Assert(gins.Config().AddPath(path), nil)
			gtest.Assert(gins.Config().Get("test"), "v=1")
			gtest.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
			gtest.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
		})
		time.Sleep(500 * time.Millisecond)

		gtest.Case(t, func() {
			path := fmt.Sprintf(`%s/%d/config`, gfile.TempDir(), gtime.TimestampNano())
			file := fmt.Sprintf(`%s/%s`, path, "config.toml")
			err := gfile.PutContents(file, configContent)
			gtest.Assert(err, nil)
			defer gfile.Remove(file)
			defer gins.Config().Clear()
			gtest.Assert(gins.Config().AddPath(path), nil)
			gtest.Assert(gins.Config().Get("test"), "v=1")
			gtest.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
			gtest.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
		})
		time.Sleep(500 * time.Millisecond)

		gtest.Case(t, func() {
			path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.TimestampNano())
			file := fmt.Sprintf(`%s/%s`, path, "test.toml")
			err := gfile.PutContents(file, configContent)
			gtest.Assert(err, nil)
			defer gfile.Remove(file)
			defer gins.Config("test").Clear()
			gins.Config("test").SetFileName("test.toml")
			gtest.Assert(gins.Config("test").AddPath(path), nil)
			gtest.Assert(gins.Config("test").Get("test"), "v=1")
			gtest.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
			gtest.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
		})
		time.Sleep(500 * time.Millisecond)

		gtest.Case(t, func() {
			path := fmt.Sprintf(`%s/%d/config`, gfile.TempDir(), gtime.TimestampNano())
			file := fmt.Sprintf(`%s/%s`, path, "test.toml")
			err := gfile.PutContents(file, configContent)
			gtest.Assert(err, nil)
			defer gfile.Remove(file)
			defer gins.Config().Clear()
			gins.Config("test").SetFileName("test.toml")
			gtest.Assert(gins.Config("test").AddPath(path), nil)
			gtest.Assert(gins.Config("test").Get("test"), "v=1")
			gtest.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
			gtest.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
		})
	})
}
func Test_Basic2(t *testing.T) {
	config := `log-path = "logs"`
	gtest.Case(t, func() {
		path := gcfg.DEFAULT_CONFIG_FILE
		err := gfile.PutContents(path, config)
		gtest.Assert(err, nil)
		defer func() {
			_ = gfile.Remove(path)
		}()

		gtest.Assert(gins.Config().Get("log-path"), "logs")
	})
}

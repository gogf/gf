// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gins_test

import (
	"fmt"
	"github.com/jin502437344/gf/debug/gdebug"
	"github.com/jin502437344/gf/frame/gins"
	"testing"
	"time"

	"github.com/jin502437344/gf/os/gcfg"

	"github.com/jin502437344/gf/os/gfile"
	"github.com/jin502437344/gf/os/gtime"
	"github.com/jin502437344/gf/test/gtest"
)

var (
	configContent = gfile.GetContents(
		gdebug.TestDataPath("config", "config.toml"),
	)
)

func Test_Config1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(configContent, "")
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(gins.Config(), nil)
	})
}

func Test_Config2(t *testing.T) {
	// relative path
	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
		t.Assert(err, nil)

		err = gins.Config().AddPath(dirPath)
		t.Assert(err, nil)

		defer gins.Config().Clear()

		t.Assert(gins.Config().Get("test"), "v=1")
		t.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
	})
	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	// relative path, config folder
	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config/config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
		t.Assert(err, nil)

		err = gins.Config().AddPath(dirPath)
		t.Assert(err, nil)

		defer gins.Config().Clear()

		t.Assert(gins.Config().Get("test"), "v=1")
		t.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)
	})
}

func Test_Config3(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "test.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
		t.Assert(err, nil)

		err = gins.Config("test").AddPath(dirPath)
		t.Assert(err, nil)

		defer gins.Config("test").Clear()
		gins.Config("test").SetFileName("test.toml")

		t.Assert(gins.Config("test").Get("test"), "v=1")
		t.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
	})
	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config/test.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), configContent)
		t.Assert(err, nil)

		err = gins.Config("test").AddPath(dirPath)
		t.Assert(err, nil)

		defer gins.Config("test").Clear()
		gins.Config("test").SetFileName("test.toml")

		t.Assert(gins.Config("test").Get("test"), "v=1")
		t.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
	})
	// for gfsnotify callbacks to refresh cache of config file for next unit testing case.
	time.Sleep(500 * time.Millisecond)
}

func Test_Config4(t *testing.T) {
	// absolute path
	gtest.C(t, func(t *gtest.T) {
		path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.TimestampNano())
		file := fmt.Sprintf(`%s/%s`, path, "config.toml")
		err := gfile.PutContents(file, configContent)
		t.Assert(err, nil)
		defer gfile.Remove(file)
		defer gins.Config().Clear()

		t.Assert(gins.Config().AddPath(path), nil)
		t.Assert(gins.Config().Get("test"), "v=1")
		t.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
	})
	time.Sleep(500 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		path := fmt.Sprintf(`%s/%d/config`, gfile.TempDir(), gtime.TimestampNano())
		file := fmt.Sprintf(`%s/%s`, path, "config.toml")
		err := gfile.PutContents(file, configContent)
		t.Assert(err, nil)
		defer gfile.Remove(file)
		defer gins.Config().Clear()
		t.Assert(gins.Config().AddPath(path), nil)
		t.Assert(gins.Config().Get("test"), "v=1")
		t.Assert(gins.Config().Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config().Get("redis.disk"), "127.0.0.1:6379,0")
	})
	time.Sleep(500 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		path := fmt.Sprintf(`%s/%d`, gfile.TempDir(), gtime.TimestampNano())
		file := fmt.Sprintf(`%s/%s`, path, "test.toml")
		err := gfile.PutContents(file, configContent)
		t.Assert(err, nil)
		defer gfile.Remove(file)
		defer gins.Config("test").Clear()
		gins.Config("test").SetFileName("test.toml")
		t.Assert(gins.Config("test").AddPath(path), nil)
		t.Assert(gins.Config("test").Get("test"), "v=1")
		t.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
	})
	time.Sleep(500 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		path := fmt.Sprintf(`%s/%d/config`, gfile.TempDir(), gtime.TimestampNano())
		file := fmt.Sprintf(`%s/%s`, path, "test.toml")
		err := gfile.PutContents(file, configContent)
		t.Assert(err, nil)
		defer gfile.Remove(file)
		defer gins.Config().Clear()
		gins.Config("test").SetFileName("test.toml")
		t.Assert(gins.Config("test").AddPath(path), nil)
		t.Assert(gins.Config("test").Get("test"), "v=1")
		t.Assert(gins.Config("test").Get("database.default.1.host"), "127.0.0.1")
		t.Assert(gins.Config("test").Get("redis.disk"), "127.0.0.1:6379,0")
	})
}
func Test_Basic2(t *testing.T) {
	config := `log-path = "logs"`
	gtest.C(t, func(t *gtest.T) {
		path := gcfg.DEFAULT_CONFIG_FILE
		err := gfile.PutContents(path, config)
		t.Assert(err, nil)
		defer func() {
			_ = gfile.Remove(path)
		}()

		t.Assert(gins.Config().Get("log-path"), "logs")
	})
}

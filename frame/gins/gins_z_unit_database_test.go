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

func Test_Database(t *testing.T) {
	databaseContent := gfile.GetContents(gfile.Join(gdebug.TestDataPath(), "database", "config.toml"))

	var err error
	dirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
	err = gfile.Mkdir(dirPath)
	gtest.Assert(err, nil)
	defer gfile.Remove(dirPath)

	name := "config.toml"
	err = gfile.PutContents(gfile.Join(dirPath, name), databaseContent)
	gtest.Assert(err, nil)

	err = gins.Config().AddPath(dirPath)
	gtest.Assert(err, nil)

	defer gins.Config().Clear()

	// for gfsnotify callbacks to refresh cache of config file
	time.Sleep(500 * time.Millisecond)

	gtest.Case(t, func() {
		//fmt.Println("gins Test_Database", Config().Get("test"))

		dbDefault := gins.Database()
		dbTest := gins.Database("test")
		gtest.AssertNE(dbDefault, nil)
		gtest.AssertNE(dbTest, nil)

		gtest.Assert(dbDefault.PingMaster(), nil)
		gtest.Assert(dbDefault.PingSlave(), nil)
		gtest.Assert(dbTest.PingMaster(), nil)
		gtest.Assert(dbTest.PingSlave(), nil)
	})
}

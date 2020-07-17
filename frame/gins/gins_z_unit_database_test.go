// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gins_test

import (
	"github.com/jin502437344/gf/debug/gdebug"
	"github.com/jin502437344/gf/frame/gins"
	"github.com/jin502437344/gf/os/gtime"
	"testing"
	"time"

	"github.com/jin502437344/gf/os/gfile"
	"github.com/jin502437344/gf/test/gtest"
)

func Test_Database(t *testing.T) {
	databaseContent := gfile.GetContents(
		gdebug.TestDataPath("database", "config.toml"),
	)
	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.Assert(err, nil)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), databaseContent)
		t.Assert(err, nil)

		err = gins.Config().AddPath(dirPath)
		t.Assert(err, nil)

		defer gins.Config().Clear()

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		//fmt.Println("gins Test_Database", Config().Get("test"))

		dbDefault := gins.Database()
		dbTest := gins.Database("test")
		t.AssertNE(dbDefault, nil)
		t.AssertNE(dbTest, nil)

		t.Assert(dbDefault.PingMaster(), nil)
		t.Assert(dbDefault.PingSlave(), nil)
		t.Assert(dbTest.PingMaster(), nil)
		t.Assert(dbTest.PingSlave(), nil)
	})
}

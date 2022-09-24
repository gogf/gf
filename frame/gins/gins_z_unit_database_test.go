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

func Test_Database(t *testing.T) {
	databaseContent := gfile.GetContents(
		gtest.DataPath("database", "config.toml"),
	)
	gtest.C(t, func(t *gtest.T) {
		var err error
		dirPath := gfile.Temp(gtime.TimestampNanoStr())
		err = gfile.Mkdir(dirPath)
		t.AssertNil(err)
		defer gfile.Remove(dirPath)

		name := "config.toml"
		err = gfile.PutContents(gfile.Join(dirPath, name), databaseContent)
		t.AssertNil(err)

		err = gins.Config().GetAdapter().(*gcfg.AdapterFile).AddPath(dirPath)
		t.AssertNil(err)

		defer gins.Config().GetAdapter().(*gcfg.AdapterFile).Clear()

		// for gfsnotify callbacks to refresh cache of config file
		time.Sleep(500 * time.Millisecond)

		// fmt.Println("gins Test_Database", Config().Get("test"))
		var (
			db        = gins.Database()
			dbDefault = gins.Database("default")
		)
		t.AssertNE(db, nil)
		t.AssertNE(dbDefault, nil)

		t.Assert(db.PingMaster(), nil)
		t.Assert(db.PingSlave(), nil)
		t.Assert(dbDefault.PingMaster(), nil)
		t.Assert(dbDefault.PingSlave(), nil)
	})
}

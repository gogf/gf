// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

var (
	ctx = context.TODO()
)

func Test_Rotate_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		p := gfile.Temp(gtime.TimestampNanoStr())
		err := l.SetConfigWithMap(g.Map{
			"Path":                 p,
			"File":                 "access.log",
			"StdoutPrint":          false,
			"RotateSize":           10,
			"RotateBackupLimit":    2,
			"RotateBackupExpire":   5 * time.Second,
			"RotateBackupCompress": 9,
			"RotateCheckInterval":  time.Second, // For unit testing only.
		})
		t.AssertNil(err)
		defer gfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			fmt.Println(ctx, "logging content index:", i)
			l.Print(ctx, s)
		}

		time.Sleep(time.Second * 3)

		files, err := gfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 2)

		content := gfile.GetContents(gfile.Join(p, "access.log"))
		t.Assert(gstr.Count(content, s), 1)

		time.Sleep(time.Second * 5)
		files, err = gfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
}

func Test_Rotate_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		p := gfile.Temp(gtime.TimestampNanoStr())
		err := l.SetConfigWithMap(g.Map{
			"Path":                 p,
			"File":                 "access.log",
			"StdoutPrint":          false,
			"RotateExpire":         time.Second,
			"RotateBackupLimit":    2,
			"RotateBackupExpire":   5 * time.Second,
			"RotateBackupCompress": 9,
			"RotateCheckInterval":  time.Second, // For unit testing only.
		})
		t.AssertNil(err)
		defer gfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			l.Print(ctx, s)
		}

		files, err := gfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 0)

		t.Assert(gstr.Count(gfile.GetContents(gfile.Join(p, "access.log")), s), 10)

		time.Sleep(time.Second * 3)

		files, err = gfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 1)

		t.Assert(gstr.Count(gfile.GetContents(gfile.Join(p, "access.log")), s), 0)

		time.Sleep(time.Second * 5)
		files, err = gfile.ScanDirFile(p, "*.gz")
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
}

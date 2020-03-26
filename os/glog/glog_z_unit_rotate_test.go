// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/text/gstr"
	"testing"
	"time"
)

func Test_Rotate_Size(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		p := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err := l.SetConfigWithMap(g.Map{
			"Path":               p,
			"File":               "access.log",
			"StdoutPrint":        false,
			"RotateSize":         10,
			"RotateBackLimit":    2,
			"RotateBackExpire":   5 * time.Second,
			"RotateBackCompress": 9,
			"RotateInterval":     time.Second, // For unit testing only.
		})
		t.Assert(err, nil)
		defer gfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			l.Print(s)
		}

		time.Sleep(time.Second * 3)

		files, err := gfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 2)

		content := gfile.GetContents(gfile.Join(p, "access.log"))
		t.Assert(gstr.Count(content, s), 1)

		time.Sleep(time.Second * 5)
		files, err = gfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 0)
	})
}

func Test_Rotate_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		p := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err := l.SetConfigWithMap(g.Map{
			"Path":               p,
			"File":               "access.log",
			"StdoutPrint":        false,
			"RotateExpire":       time.Second,
			"RotateBackLimit":    2,
			"RotateBackExpire":   5 * time.Second,
			"RotateBackCompress": 9,
			"RotateInterval":     time.Second, // For unit testing only.
		})
		t.Assert(err, nil)
		defer gfile.Remove(p)

		s := "1234567890abcdefg"
		for i := 0; i < 10; i++ {
			l.Print(s)
		}

		files, err := gfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 0)

		t.Assert(gstr.Count(gfile.GetContents(gfile.Join(p, "access.log")), s), 10)

		time.Sleep(time.Second * 3)

		files, err = gfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 1)

		t.Assert(gstr.Count(gfile.GetContents(gfile.Join(p, "access.log")), s), 0)

		time.Sleep(time.Second * 5)
		files, err = gfile.ScanDirFile(p, "*.gz")
		t.Assert(err, nil)
		t.Assert(len(files), 0)
	})
}

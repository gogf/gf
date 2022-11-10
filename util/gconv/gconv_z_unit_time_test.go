// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Time(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Duration(""), time.Duration(int64(0)))
		t.AssertEQ(gconv.GTime(""), gtime.New())
		t.AssertEQ(gconv.GTime(nil), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		s := "2011-10-10 01:02:03.456"
		t.AssertEQ(gconv.GTime(s), gtime.NewFromStr(s))
		t.AssertEQ(gconv.Time(nil), time.Time{})
		t.AssertEQ(gconv.Time(s), gtime.NewFromStr(s).Time)
		t.AssertEQ(gconv.Duration(100), 100*time.Nanosecond)
	})
	gtest.C(t, func(t *gtest.T) {
		s := "01:02:03.456"
		t.AssertEQ(gconv.GTime(s).Hour(), 1)
		t.AssertEQ(gconv.GTime(s).Minute(), 2)
		t.AssertEQ(gconv.GTime(s).Second(), 3)
		t.AssertEQ(gconv.GTime(s), gtime.NewFromStr(s))
		t.AssertEQ(gconv.Time(s), gtime.NewFromStr(s).Time)
	})
	gtest.C(t, func(t *gtest.T) {
		s := "0000-01-01 01:02:03"
		t.AssertEQ(gconv.GTime(s).Year(), 0)
		t.AssertEQ(gconv.GTime(s).Month(), 1)
		t.AssertEQ(gconv.GTime(s).Day(), 1)
		t.AssertEQ(gconv.GTime(s).Hour(), 1)
		t.AssertEQ(gconv.GTime(s).Minute(), 2)
		t.AssertEQ(gconv.GTime(s).Second(), 3)
		t.AssertEQ(gconv.GTime(s), gtime.NewFromStr(s))
		t.AssertEQ(gconv.Time(s), gtime.NewFromStr(s).Time)
	})
	gtest.C(t, func(t *gtest.T) {
		t1 := gtime.NewFromStr("2021-05-21 05:04:51.206547+00")
		t2 := gconv.GTime(gvar.New(t1))
		t3 := gvar.New(t1).GTime()
		t.AssertEQ(t1, t2)
		t.AssertEQ(t1.Local(), t2.Local())
		t.AssertEQ(t1, t3)
		t.AssertEQ(t1.Local(), t3.Local())
	})
}

func Test_Time_Slice_Attribute(t *testing.T) {
	type SelectReq struct {
		Arr []*gtime.Time
		One *gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		var s *SelectReq
		err := gconv.Struct(g.Map{
			"arr": g.Slice{"2021-01-12 12:34:56", "2021-01-12 12:34:57"},
			"one": "2021-01-12 12:34:58",
		}, &s)
		t.AssertNil(err)
		t.Assert(s.One, "2021-01-12 12:34:58")
		t.Assert(s.Arr[0], "2021-01-12 12:34:56")
		t.Assert(s.Arr[1], "2021-01-12 12:34:57")
	})
}

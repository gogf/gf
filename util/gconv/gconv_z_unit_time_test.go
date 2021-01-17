// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/frame/g"
	"testing"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func Test_Time(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "2011-10-10 01:02:03.456"
		t.AssertEQ(gconv.GTime(s), gtime.NewFromStr(s))
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
		t.Assert(err, nil)
		t.Assert(s.One, "2021-01-12 12:34:58")
		t.Assert(s.Arr[0], "2021-01-12 12:34:56")
		t.Assert(s.Arr[1], "2021-01-12 12:34:57")
	})
}

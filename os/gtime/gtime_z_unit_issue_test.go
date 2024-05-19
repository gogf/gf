// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// https://github.com/gogf/gf/issues/1681
func Test_Issue1681(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gtime.New("2022-03-08T03:01:14-07:00").Local().Time, gtime.New("2022-03-08T10:01:14Z").Local().Time)
		t.Assert(gtime.New("2022-03-08T03:01:14-08:00").Local().Time, gtime.New("2022-03-08T11:01:14Z").Local().Time)
		t.Assert(gtime.New("2022-03-08T03:01:14-09:00").Local().Time, gtime.New("2022-03-08T12:01:14Z").Local().Time)
		t.Assert(gtime.New("2022-03-08T03:01:14+08:00").Local().Time, gtime.New("2022-03-07T19:01:14Z").Local().Time)
	})
}

// https://github.com/gogf/gf/issues/2803
func Test_Issue2803(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		newTime := gtime.New("2023-07-26").LayoutTo("2006-01")
		t.Assert(newTime.Year(), 2023)
		t.Assert(newTime.Month(), 7)
		t.Assert(newTime.Day(), 1)
		t.Assert(newTime.Hour(), 0)
		t.Assert(newTime.Minute(), 0)
		t.Assert(newTime.Second(), 0)
	})
}

// https://github.com/gogf/gf/issues/3558
func Test_Issue3558(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeStr := "1880-10-24T00:00:00+08:05"
		gfTime := gtime.NewFromStr(timeStr)
		t.Assert(gfTime.Year(), 1880)
		t.Assert(gfTime.Month(), 10)
		t.Assert(gfTime.Day(), 24)
		t.Assert(gfTime.Hour(), 0)
		t.Assert(gfTime.Minute(), 0)
		t.Assert(gfTime.Second(), 0)

		stdTime, err := time.Parse(time.RFC3339, timeStr)
		t.AssertNil(err)
		stdTimeFormat := stdTime.Format("2006-01-02 15:04:05")
		gfTimeFormat := gfTime.Format("Y-m-d H:i:s")
		t.Assert(gfTimeFormat, stdTimeFormat)
	})
	gtest.C(t, func(t *gtest.T) {
		timeStr := "1880-10-24T00:00:00-08:05"
		gfTime := gtime.NewFromStr(timeStr)
		t.Assert(gfTime.Year(), 1880)
		t.Assert(gfTime.Month(), 10)
		t.Assert(gfTime.Day(), 24)
		t.Assert(gfTime.Hour(), 0)
		t.Assert(gfTime.Minute(), 0)
		t.Assert(gfTime.Second(), 0)
		stdTime, err := time.Parse(time.RFC3339, timeStr)
		t.AssertNil(err)
		stdTimeFormat := stdTime.Format("2006-01-02 15:04:05")
		gfTimeFormat := gfTime.Format("Y-m-d H:i:s")
		t.Assert(gfTimeFormat, stdTimeFormat)
	})
}

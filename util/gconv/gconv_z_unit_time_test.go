// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	timeStrTests  = "2024-04-22 12:00:00.123456789+00:00:00"
	timeTimeTests = time.Date(
		2024, 4, 22, 12, 0, 0, 123456789, time.UTC,
	)
)

func TestTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Time(nil), time.Time{})
		t.AssertEQ(gconv.Time(timeTimeTests), timeTimeTests)
	})
}

func TestDuration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Duration(nil), time.Duration(0))
		t.AssertEQ(gconv.Duration(timeTimeTests), time.Duration(0))
		t.AssertEQ(gconv.Duration("1m"), time.Minute)
		t.AssertEQ(gconv.Duration(time.Hour), time.Hour)
	})
}

func TestGtime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.GTime(""), gtime.New())
		t.AssertEQ(gconv.GTime(nil), nil)

		t.AssertEQ(gconv.GTime(gtime.New(timeStrTests)), gtime.New(timeStrTests))
		t.AssertEQ(gconv.GTime(timeTimeTests).Year(), 2024)
		t.AssertEQ(gconv.GTime(timeTimeTests).Month(), 4)
		t.AssertEQ(gconv.GTime(timeTimeTests).Day(), 22)
		t.AssertEQ(gconv.GTime(timeTimeTests).Hour(), 12)
		t.AssertEQ(gconv.GTime(timeTimeTests).Minute(), 0)
		t.AssertEQ(gconv.GTime(timeTimeTests).Second(), 0)
		t.AssertEQ(gconv.GTime(timeTimeTests).Nanosecond(), 123456789)
		t.AssertEQ(gconv.GTime(timeTimeTests).String(), "2024-04-22 12:00:00")
	})
}

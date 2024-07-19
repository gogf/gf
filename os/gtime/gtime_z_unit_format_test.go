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

func Test_Format(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp, err := gtime.StrToTime("2006-01-11 15:04:05", "Y-m-d H:i:s")
		timeTemp.ToZone("Asia/Shanghai")
		if err != nil {
			t.Error("test fail")
		}
		t.Assert(timeTemp.Format("\\T\\i\\m\\e中文Y-m-j G:i:s.u\\"), "Time中文2006-01-11 15:04:05.000")

		t.Assert(timeTemp.Format("d D j l"), "11 Wed 11 Wednesday")

		t.Assert(timeTemp.Format("F m M n"), "January 01 Jan 1")

		t.Assert(timeTemp.Format("Y y"), "2006 06")

		t.Assert(timeTemp.Format("a A g G h H i s u .u"), "pm PM 3 15 03 15 04 05 000 .000")

		t.Assert(timeTemp.Format("O P T"), "+0800 +08:00 CST")

		t.Assert(timeTemp.Format("r"), "Wed, 11 Jan 06 15:04 CST")

		t.Assert(timeTemp.Format("c"), "2006-01-11T15:04:05+08:00")

		//补零
		timeTemp1, err := gtime.StrToTime("2006-01-02 03:04:05", "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		t.Assert(timeTemp1.Format("Y-m-d h:i:s"), "2006-01-02 03:04:05")
		//不补零
		timeTemp2, err := gtime.StrToTime("2006-01-02 03:04:05", "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		t.Assert(timeTemp2.Format("Y-n-j G:i:s"), "2006-1-2 3:04:05")

		t.Assert(timeTemp2.Format("U"), "1136142245")

		// 测试数字型的星期
		times := []map[string]string{
			{"k": "2019-04-22", "f": "w", "r": "1"},
			{"k": "2019-04-23", "f": "w", "r": "2"},
			{"k": "2019-04-24", "f": "w", "r": "3"},
			{"k": "2019-04-25", "f": "w", "r": "4"},
			{"k": "2019-04-26", "f": "dw", "r": "265"},
			{"k": "2019-04-27", "f": "w", "r": "6"},
			{"k": "2019-03-10", "f": "w", "r": "0"},
			{"k": "2019-03-10", "f": "Y-m-d 星期:w", "r": "2019-03-10 星期:0"},
			{"k": "2019-04-25", "f": "N", "r": "4"},
			{"k": "2019-03-10", "f": "N", "r": "7"},
			{"k": "2019-03-01", "f": "S", "r": "st"},
			{"k": "2019-03-02", "f": "S", "r": "nd"},
			{"k": "2019-03-03", "f": "S", "r": "rd"},
			{"k": "2019-03-05", "f": "S", "r": "th"},

			{"k": "2019-01-01", "f": "第z天", "r": "第0天"},
			{"k": "2019-01-05", "f": "第z天", "r": "第4天"},
			{"k": "2020-05-05", "f": "第z天", "r": "第125天"},
			{"k": "2020-12-31", "f": "第z天", "r": "第365天"}, //润年
			{"k": "2020-02-12", "f": "第z天", "r": "第42天"},  //润年
			{"k": "2019-02-12", "f": "有t天", "r": "有28天"},
			{"k": "2020-02-12", "f": "20.2有t天", "r": "20.2有29天"},
			{"k": "2019-03-12", "f": "19.3有t天", "r": "19.3有31天"},
			{"k": "2019-11-12", "f": "19.11有t天", "r": "19.11有30天"},
			{"k": "2019-01-01", "f": "第W周", "r": "第1周"},
			{"k": "2017-01-01", "f": "第W周", "r": "第52周"},         //星期7
			{"k": "2002-01-01", "f": "第W周为星期2", "r": "第1周为星期2"},  //星期2
			{"k": "2016-01-01", "f": "第W周为星期5", "r": "第53周为星期5"}, //星期5
			{"k": "2014-01-01", "f": "第W周为星期3", "r": "第1周为星期3"},  //星期3
			{"k": "2015-01-01", "f": "第W周为星期4", "r": "第1周为星期4"},  //星期4
		}

		for _, v := range times {
			t1, err1 := gtime.StrToTime(v["k"], "Y-m-d")
			t.Assert(err1, nil)
			t.Assert(t1.Format(v["f"]), v["r"])
		}

	})
	gtest.C(t, func(t *gtest.T) {
		var ti *gtime.Time = nil
		t.Assert(ti.Format("Y-m-d h:i:s"), "")
		t.Assert(ti.FormatNew("Y-m-d h:i:s"), nil)
		t.Assert(ti.FormatTo("Y-m-d h:i:s"), nil)
		t.Assert(ti.Layout("2006-01-02 15:04:05"), "")
		t.Assert(ti.LayoutNew("2006-01-02 15:04:05"), nil)
		t.Assert(ti.LayoutTo("2006-01-02 15:04:05"), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		timeOnly, err := gtime.StrToTime("15:04:05")
		t.Assert(err, nil)

		t.Assert(timeOnly.Format("Y-m-d h:i:s"), "0001-01-01 03:04:05")
		t.Assert(timeOnly.FormatNew("Y-m-d h:i:s"), "0001-01-01 03:04:05")
		t.Assert(timeOnly.Format("Y-m-d H:i:s"), "0001-01-01 15:04:05")
		t.Assert(timeOnly.FormatNew("Y-m-d H:i:s"), "0001-01-01 15:04:05")
		t.Assert(timeOnly.Layout("2006-01-02 15:04:05"), "0001-01-01 15:04:05")
		t.Assert(timeOnly.LayoutNew("2006-01-02 15:04:05"), "0001-01-01 15:04:05")

		t.Assert(timeOnly.Clone().FormatTo("Y-m-d h:i:s"), "0001-01-01 03:04:05")
		t.Assert(timeOnly.Clone().FormatTo("Y-m-d H:i:s"), "0001-01-01 15:04:05")
		t.Assert(timeOnly.Clone().LayoutTo("2006-01-02 15:04:05"), "0001-01-01 15:04:05")

		d := timeOnly.UTC().Truncate(gtime.D)
		t.Assert(d.IsZero(), true)

		// std time.IsZero: zero time instant, January 1, year 1, 00:00:00 UTC
		zero := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		t.Assert(d.Time.Equal(zero), true)
	})
}

func Test_Format_ZeroString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 0000-00-00 00:00:00 is not a valid time
		timeTemp, err := gtime.StrToTime("0000-00-00 00:00:00")
		t.AssertNE(err, nil)
		t.AssertEQ(timeTemp, nil)
		t.Assert(timeTemp.String(), "")
	})

	const zeroTimeStr = "0001-01-01 00:00:00"

	gtest.C(t, func(t *gtest.T) {
		// std time.IsZero: zero time instant, January 1, year 1, 00:00:00 UTC
		zero := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		t.AssertEQ(zero.IsZero(), true)
		t.AssertEQ(zero.Format("2006-01-02 15:04:05"), zeroTimeStr)

		year0 := time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC)
		t.AssertEQ(year0.IsZero(), false)
		t.AssertEQ(year0.Format("2006-01-02 15:04:05"), "0000-01-01 00:00:00")
	})

	gtest.C(t, func(t *gtest.T) {
		// set timezone to UTC to get zero time instant
		err := gtime.SetTimeZone(time.UTC.String())
		t.Assert(err, nil)

		zeroDateTime, err := gtime.StrToTime(zeroTimeStr)
		t.AssertEQ(err, nil)
		t.Assert(zeroDateTime.IsZero(), true)
		t.Assert(zeroDateTime.String(), "")

		zeroTimeOnly, err := gtime.StrToTime("00:00:00")
		t.AssertEQ(err, nil)
		t.Assert(zeroTimeOnly.IsZero(), true)
		t.Assert(zeroTimeOnly.String(), "")
	})
}

func Test_FormatTo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.FormatTo("Y-m-01 00:00:01"), timeTemp.Time.Format("2006-01")+"-01 00:00:01")
	})
}

func Test_Layout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.Layout("2006-01-02 15:04:05"), timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
}

func Test_LayoutTo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.LayoutTo("2006-01-02 00:00:00"), timeTemp.Time.Format("2006-01-02 00:00:00"))
	})
}

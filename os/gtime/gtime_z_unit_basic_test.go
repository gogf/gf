// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_SetTimeZone(t *testing.T) {
	gtest.Case(t, func() {
		gtime.SetTimeZone("Asia/Shanghai")
		gtest.Assert(time.Local.String(), "Asia/Shanghai")
	})
}

func Test_Nanosecond(t *testing.T) {
	gtest.Case(t, func() {
		nanos := gtime.TimestampNano()
		timeTemp := time.Unix(0, nanos)
		gtest.Assert(nanos, timeTemp.UnixNano())
	})
}

func Test_Microsecond(t *testing.T) {
	gtest.Case(t, func() {
		micros := gtime.TimestampMicro()
		timeTemp := time.Unix(0, micros*1e3)
		gtest.Assert(micros, timeTemp.UnixNano()/1e3)
	})
}

func Test_Millisecond(t *testing.T) {
	gtest.Case(t, func() {
		millis := gtime.TimestampMilli()
		timeTemp := time.Unix(0, millis*1e6)
		gtest.Assert(millis, timeTemp.UnixNano()/1e6)
	})
}

func Test_Second(t *testing.T) {
	gtest.Case(t, func() {
		s := gtime.Timestamp()
		timeTemp := time.Unix(s, 0)
		gtest.Assert(s, timeTemp.Unix())
	})
}

func Test_Date(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gtime.Date(), time.Now().Format("2006-01-02"))
	})
}

func Test_Datetime(t *testing.T) {
	gtest.Case(t, func() {
		datetime := gtime.Datetime()
		timeTemp, err := gtime.StrToTime(datetime, "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		gtest.Assert(datetime, timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
}

func Test_ISO8601(t *testing.T) {
	gtest.Case(t, func() {
		iso8601 := gtime.ISO8601()
		gtest.Assert(iso8601, gtime.Now().Format("c"))
	})
}

func Test_RFC822(t *testing.T) {
	gtest.Case(t, func() {
		rfc822 := gtime.RFC822()
		gtest.Assert(rfc822, gtime.Now().Format("r"))
	})
}

func Test_StrToTime(t *testing.T) {
	gtest.Case(t, func() {
		//正常日期列表
		//正则的原因，日期"06.01.02"，"2006.01"，"2006..01"无法覆盖gtime.go的百分百
		var testDatetimes = []string{
			"2006-01-02 15:04:05",
			"2006/01/02 15:04:05",
			"2006.01.02 15:04:05.000",
			"2006.01.02 - 15:04:05",
			"2006.01.02 15:04:05 +0800 CST",
			"2006-01-02T20:05:06+05:01:01",
			"2006-01-02T14:03:04Z01:01:01",
			"2006-01-02T15:04:05Z",
			"02-jan-2006 15:04:05",
			"02/jan/2006 15:04:05",
			"02.jan.2006 15:04:05",
			"02.jan.2006:15:04:05",
		}

		for _, item := range testDatetimes {
			timeTemp, err := gtime.StrToTime(item)
			if err != nil {
				t.Error("test fail")
			}
			gtest.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")
		}

		//正常日期列表，时间00:00:00
		var testDates = []string{
			"2006.01.02",
			"2006.01.02 00:00",
			"2006.01.02 00:00:00.000",
		}

		for _, item := range testDates {
			timeTemp, err := gtime.StrToTime(item)
			if err != nil {
				t.Error("test fail")
			}
			gtest.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 00:00:00")
		}

		//测试格式化formatToStdLayout
		var testDateFormats = []string{
			"Y-m-d H:i:s",
			"\\T\\i\\m\\e Y-m-d H:i:s",
			"Y-m-d H:i:s\\",
			"Y-m-j G:i:s.u",
			"Y-m-j G:i:su",
		}

		var testDateFormatsResult = []string{
			"2007-01-02 15:04:05",
			"Time 2007-01-02 15:04:05",
			"2007-01-02 15:04:05",
			"2007-01-02 15:04:05.000",
			"2007-01-02 15:04:05.000",
		}

		for index, item := range testDateFormats {
			timeTemp, err := gtime.StrToTime(testDateFormatsResult[index], item)
			if err != nil {
				t.Error("test fail")
			}
			gtest.Assert(timeTemp.Time.Format("2006-01-02 15:04:05.000"), "2007-01-02 15:04:05.000")
		}

		//异常日期列表
		var testDatesFail = []string{
			"2006.01",
			"06..02",
			"20060102",
		}

		for _, item := range testDatesFail {
			_, err := gtime.StrToTime(item)
			if err == nil {
				t.Error("test fail")
			}
		}

		//test err
		_, err := gtime.StrToTime("2006-01-02 15:04:05", "aabbccdd")
		if err == nil {
			t.Error("test fail")
		}
	})
}

func Test_ConvertZone(t *testing.T) {
	gtest.Case(t, func() {
		//现行时间
		nowUTC := time.Now().UTC()
		testZone := "America/Los_Angeles"

		//转换为洛杉矶时间
		t1, err := gtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, "")
		if err != nil {
			t.Error("test fail")
		}

		//使用洛杉矶时区解析上面转换后的时间
		laStr := t1.Time.Format("2006-01-02 15:04:05")
		loc, err := time.LoadLocation(testZone)
		t2, err := time.ParseInLocation("2006-01-02 15:04:05", laStr, loc)

		//判断是否与现行时间匹配
		gtest.Assert(t2.UTC().Unix(), nowUTC.Unix())

	})

	//test err
	gtest.Case(t, func() {
		//现行时间
		nowUTC := time.Now().UTC()
		//t.Log(nowUTC.Unix())
		testZone := "errZone"

		//错误时间输入
		_, err := gtime.ConvertZone(nowUTC.Format("06..02 15:04:05"), testZone, "")
		if err == nil {
			t.Error("test fail")
		}
		//错误时区输入
		_, err = gtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, "")
		if err == nil {
			t.Error("test fail")
		}
		//错误时区输入
		_, err = gtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, testZone)
		if err == nil {
			t.Error("test fail")
		}
	})
}

func Test_StrToTimeFormat(t *testing.T) {
	gtest.Case(t, func() {

	})
}

func Test_ParseTimeFromContent(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.ParseTimeFromContent("我是中文2006-01-02 15:04:05我也是中文", "Y-m-d H:i:s")
		gtest.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.ParseTimeFromContent("我是中文2006-01-02 15:04:05我也是中文")
		gtest.Assert(timeTemp1.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		timeTemp2 := gtime.ParseTimeFromContent("我是中文02.jan.2006 15:04:05我也是中文")
		gtest.Assert(timeTemp2.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		//test err
		timeTempErr := gtime.ParseTimeFromContent("我是中文", "Y-m-d H:i:s")
		if timeTempErr != nil {
			t.Error("test fail")
		}
	})
}

func Test_FuncCost(t *testing.T) {
	gtest.Case(t, func() {
		gtime.FuncCost(func() {

		})
	})
}

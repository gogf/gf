// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_SetTimeZone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gtime.SetTimeZone("Asia/Shanghai"), nil)
		t.AssertNE(gtime.SetTimeZone("testzone"), nil)
		// t.Assert(time.Local.String(), "Asia/Shanghai")
	})
}

func Test_TimestampStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(len(gtime.TimestampMilliStr()), 0)
		t.AssertGT(len(gtime.TimestampMicroStr()), 0)
		t.AssertGT(len(gtime.TimestampNanoStr()), 0)
	})
}

func Test_Nanosecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		nanos := gtime.TimestampNano()
		timeTemp := time.Unix(0, nanos)
		t.Assert(nanos, timeTemp.UnixNano())
	})
}

func Test_Microsecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		micros := gtime.TimestampMicro()
		timeTemp := time.Unix(0, micros*1e3)
		t.Assert(micros, timeTemp.UnixNano()/1e3)
	})
}

func Test_Millisecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		millis := gtime.TimestampMilli()
		timeTemp := time.Unix(0, millis*1e6)
		t.Assert(millis, timeTemp.UnixNano()/1e6)
	})
}

func Test_Second(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gtime.Timestamp()
		timeTemp := time.Unix(s, 0)
		t.Assert(s, timeTemp.Unix())
	})
}

func Test_Date(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gtime.Date(), time.Now().Format("2006-01-02"))
	})
}

func Test_Datetime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		datetime := gtime.Datetime()
		timeTemp, err := gtime.StrToTime(datetime, "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		t.Assert(datetime, timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp, err := gtime.StrToTime("")
		t.Assert(err, nil)
		t.AssertLT(timeTemp.Unix(), 0)
		timeTemp, err = gtime.StrToTime("2006-01")
		t.AssertNE(err, nil)
		t.Assert(timeTemp, nil)
	})
}

func Test_ISO8601(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		iso8601 := gtime.ISO8601()
		t.Assert(iso8601, gtime.Now().Format("c"))
	})
}

func Test_RFC822(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		rfc822 := gtime.RFC822()
		t.Assert(rfc822, gtime.Now().Format("r"))
	})
}

func Test_StrToTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Correct datetime string.
		var testDateTimes = []string{
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

		for _, item := range testDateTimes {
			timeTemp, err := gtime.StrToTime(item)
			t.AssertNil(err)
			t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")
		}

		// Correct date string,.
		var testDates = []string{
			"2006.01.02",
			"2006.01.02 00:00",
			"2006.01.02 00:00:00.000",
		}

		for _, item := range testDates {
			timeTemp, err := gtime.StrToTime(item)
			t.AssertNil(err)
			t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 00:00:00")
		}

		// Correct time string.
		var testTimes = g.MapStrStr{
			"16:12:01":     "15:04:05",
			"16:12:01.789": "15:04:05.000",
		}

		for k, v := range testTimes {
			time1, err := gtime.StrToTime(k)
			t.AssertNil(err)
			time2, err := time.ParseInLocation(v, k, time.Local)
			t.AssertNil(err)
			t.Assert(time1.Time, time2)
		}

		// formatToStdLayout
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
			t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05.000"), "2007-01-02 15:04:05.000")
		}

		// 异常日期列表
		var testDatesFail = []string{
			"2006.01",
			"06..02",
		}

		for _, item := range testDatesFail {
			_, err := gtime.StrToTime(item)
			if err == nil {
				t.Error("test fail")
			}
		}

		// test err
		_, err := gtime.StrToTime("2006-01-02 15:04:05", "aabbccdd")
		if err == nil {
			t.Error("test fail")
		}
	})
}

func Test_ConvertZone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 现行时间
		nowUTC := time.Now().UTC()
		testZone := "America/Los_Angeles"

		// 转换为洛杉矶时间
		t1, err := gtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, "")
		if err != nil {
			t.Error("test fail")
		}

		// 使用洛杉矶时区解析上面转换后的时间
		laStr := t1.Time.Format("2006-01-02 15:04:05")
		loc, err := time.LoadLocation(testZone)
		t2, err := time.ParseInLocation("2006-01-02 15:04:05", laStr, loc)

		// 判断是否与现行时间匹配
		t.Assert(t2.UTC().Unix(), nowUTC.Unix())

	})

	// test err
	gtest.C(t, func(t *gtest.T) {
		// 现行时间
		nowUTC := time.Now().UTC()
		// t.Log(nowUTC.Unix())
		testZone := "errZone"

		// 错误时间输入
		_, err := gtime.ConvertZone(nowUTC.Format("06..02 15:04:05"), testZone, "")
		if err == nil {
			t.Error("test fail")
		}
		// 错误时区输入
		_, err = gtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, "")
		if err == nil {
			t.Error("test fail")
		}
		// 错误时区输入
		_, err = gtime.ConvertZone(nowUTC.Format("2006-01-02 15:04:05"), testZone, testZone)
		if err == nil {
			t.Error("test fail")
		}
	})
}

func Test_ParseDuration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		d, err := gtime.ParseDuration("1d")
		t.AssertNil(err)
		t.Assert(d.String(), "24h0m0s")
	})
	gtest.C(t, func(t *gtest.T) {
		d, err := gtime.ParseDuration("1d2h3m")
		t.AssertNil(err)
		t.Assert(d.String(), "26h3m0s")
	})
	gtest.C(t, func(t *gtest.T) {
		d, err := gtime.ParseDuration("-1d2h3m")
		t.AssertNil(err)
		t.Assert(d.String(), "-26h3m0s")
	})
	gtest.C(t, func(t *gtest.T) {
		d, err := gtime.ParseDuration("3m")
		t.AssertNil(err)
		t.Assert(d.String(), "3m0s")
	})
	// error
	gtest.C(t, func(t *gtest.T) {
		d, err := gtime.ParseDuration("-1dd2h3m")
		t.AssertNE(err, nil)
		t.Assert(d.String(), "0s")
	})
}

func Test_ParseTimeFromContent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.ParseTimeFromContent("我是中文2006-01-02 15:04:05我也是中文", "Y-m-d H:i:s")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.ParseTimeFromContent("我是中文2006-01-02 15:04:05我也是中文")
		t.Assert(timeTemp1.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		timeTemp2 := gtime.ParseTimeFromContent("我是中文02.jan.2006 15:04:05我也是中文")
		t.Assert(timeTemp2.Time.Format("2006-01-02 15:04:05"), "2006-01-02 15:04:05")

		// test err
		timeTempErr := gtime.ParseTimeFromContent("我是中文", "Y-m-d H:i:s")
		if timeTempErr != nil {
			t.Error("test fail")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		timeStr := "2021-1-27 9:10:24"
		t.Assert(gtime.ParseTimeFromContent(timeStr, "Y-n-d g:i:s").String(), "2021-01-27 09:10:24")
	})
}

func Test_FuncCost(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gtime.FuncCost(func() {

		})
	})
}

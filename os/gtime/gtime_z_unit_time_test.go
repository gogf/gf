// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtime_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_New(t *testing.T) {
	// time.Time
	gtest.C(t, func(t *gtest.T) {
		timeNow := time.Now()
		timeTemp := gtime.New(timeNow)
		t.Assert(timeTemp.Time.UnixNano(), timeNow.UnixNano())

		timeTemp1 := gtime.New()
		t.Assert(timeTemp1.Time, time.Time{})
	})
	// string
	gtest.C(t, func(t *gtest.T) {
		timeNow := gtime.Now()
		timeTemp := gtime.New(timeNow.String())
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), timeNow.Time.Format("2006-01-02 15:04:05"))
	})
	gtest.C(t, func(t *gtest.T) {
		timeNow := gtime.Now()
		timeTemp := gtime.New(timeNow.TimestampMicroStr())
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), timeNow.Time.Format("2006-01-02 15:04:05"))
	})
	// int64
	gtest.C(t, func(t *gtest.T) {
		timeNow := gtime.Now()
		timeTemp := gtime.New(timeNow.TimestampMicro())
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), timeNow.Time.Format("2006-01-02 15:04:05"))
	})
	// short datetime.
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.New("2021-2-9 08:01:21")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2021-02-09 08:01:21")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2021-02-09 08:01:21")

		timeTemp = gtime.New("2021-02-09 08:01:21", []byte("Y-m-d H:i:s"))
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2021-02-09 08:01:21")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2021-02-09 08:01:21")

		timeTemp = gtime.New([]byte("2021-02-09 08:01:21"))
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2021-02-09 08:01:21")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2021-02-09 08:01:21")

		timeTemp = gtime.New([]byte("2021-02-09 08:01:21"), "Y-m-d H:i:s")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2021-02-09 08:01:21")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2021-02-09 08:01:21")

		timeTemp = gtime.New([]byte("2021-02-09 08:01:21"), []byte("Y-m-d H:i:s"))
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2021-02-09 08:01:21")
		t.Assert(timeTemp.Time.Format("2006-01-02 15:04:05"), "2021-02-09 08:01:21")
	})
	//
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gtime.New(gtime.Time{}), nil)
		t.Assert(gtime.New(&gtime.Time{}), nil)
	})

	// unconventional
	gtest.C(t, func(t *gtest.T) {

		var testUnconventionalDates = []string{
			"2006-01.02",
			"2006.01-02",
		}

		for _, item := range testUnconventionalDates {
			timeTemp := gtime.New(item)
			t.Assert(timeTemp.TimestampMilli(), 0)
			t.Assert(timeTemp.TimestampMilliStr(), "")
			t.Assert(timeTemp.String(), "")
		}
	})
}

func Test_Nil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var t1 *gtime.Time
		t.Assert(t1.String(), "")
	})
	gtest.C(t, func(t *gtest.T) {
		var t1 gtime.Time
		t.Assert(t1.String(), "")
	})
}

func Test_NewFromStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2006-01-02 15:04:05")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.NewFromStr("2006.0102")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t1 := gtime.NewFromStr("2006-01-02 15:04:05")
		t.Assert(t1.String(), "2006-01-02 15:04:05")
		t.Assert(fmt.Sprintf("%s", t1), "2006-01-02 15:04:05")

		t2 := *t1
		t.Assert(t2.String(), "2006-01-02 15:04:05")
		t.Assert(fmt.Sprintf("{%s}", t2.String()), "{2006-01-02 15:04:05}")
	})
}

func Test_NewFromStrFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStrFormat("2006-01-02 15:04:05", "Y-m-d H:i:s")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.NewFromStrFormat("2006-01-02 15:04:05", "aabbcc")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})

	gtest.C(t, func(t *gtest.T) {
		t1 := gtime.NewFromStrFormat("2019/2/1", "Y/n/j")
		t.Assert(t1.Format("Y-m-d"), "2019-02-01")

		t2 := gtime.NewFromStrFormat("2019/10/12", "Y/n/j")
		t.Assert(t2.Format("Y-m-d"), "2019-10-12")
	})
}

func Test_NewFromStrLayout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStrLayout("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.NewFromStrLayout("2006-01-02 15:04:05", "aabbcc")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_NewFromTimeStamp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromTimeStamp(1554459846000)
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2019-04-05 18:24:06")
		timeTemp1 := gtime.NewFromTimeStamp(0)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "0001-01-01 00:00:00")
		timeTemp2 := gtime.NewFromTimeStamp(155445984)
		t.Assert(timeTemp2.Format("Y-m-d H:i:s"), "1974-12-05 11:26:24")
	})
}

func Test_Time_Second(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.Second(), timeTemp.Time.Second())
	})
}

func Test_Time_IsZero(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ti *gtime.Time = nil
		t.Assert(ti.IsZero(), true)
	})
}

func Test_Time_AddStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gt := gtime.New("2018-08-08 08:08:08")
		gt1, err := gt.AddStr("10T")
		t.Assert(gt1, nil)
		t.AssertNE(err, nil)
	})
}

func Test_Time_Equal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var t1 *gtime.Time = nil
		var t2 = gtime.New()
		t.Assert(t1.Equal(t2), false)
		t.Assert(t1.Equal(t1), true)
		t.Assert(t2.Equal(t1), false)
	})
}

func Test_Time_After(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var t1 *gtime.Time = nil
		var t2 = gtime.New()
		t.Assert(t1.After(t2), false)
		t.Assert(t2.After(t1), true)
	})
}

func Test_Time_Sub(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var t1 *gtime.Time = nil
		var t2 = gtime.New()
		t.Assert(t1.Sub(t2), time.Duration(0))
		t.Assert(t2.Sub(t1), time.Duration(0))
	})
}

func Test_Time_Nanosecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.Nanosecond(), timeTemp.Time.Nanosecond())
	})
}

func Test_Time_Microsecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.Microsecond(), timeTemp.Time.Nanosecond()/1e3)
	})
}

func Test_Time_Millisecond(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.Millisecond(), timeTemp.Time.Nanosecond()/1e6)
	})
}

func Test_Time_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		t.Assert(timeTemp.String(), timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
}

func Test_Time_ISO8601(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		now := gtime.Now()
		t.Assert(now.ISO8601(), now.Format("c"))
	})
}

func Test_Time_RFC822(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		now := gtime.Now()
		t.Assert(now.RFC822(), now.Format("r"))
	})
}

func Test_Clone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Clone()
		t.Assert(timeTemp.Time.Unix(), timeTemp1.Time.Unix())
	})
}

func Test_ToTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		t.Assert(timeTemp.Time.UnixNano(), timeTemp1.UnixNano())
	})
}

func Test_Add(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2006-01-02 15:04:05")
		timeTemp = timeTemp.Add(time.Second)
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:06")
	})
}

func Test_ToZone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp, _ = timeTemp.ToZone("America/Los_Angeles")
		t.Assert(timeTemp.Time.Location().String(), "America/Los_Angeles")

		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			t.Error("test fail")
		}
		timeTemp = timeTemp.ToLocation(loc)
		t.Assert(timeTemp.Time.Location().String(), "Asia/Shanghai")

		timeTemp1, _ := timeTemp.ToZone("errZone")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_AddDate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2006-01-02 15:04:05")
		timeTemp = timeTemp.AddDate(1, 2, 3)
		t.Assert(timeTemp.Format("Y-m-d H:i:s"), "2007-03-05 15:04:05")
	})
}

func Test_UTC(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.UTC()
		t.Assert(timeTemp.UnixNano(), timeTemp1.UTC().UnixNano())
	})
}

func Test_Local(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.Local()
		t.Assert(timeTemp.UnixNano(), timeTemp1.Local().UnixNano())
	})
}

func Test_Round(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp = timeTemp.Round(time.Hour)
		t.Assert(timeTemp.UnixNano(), timeTemp1.Round(time.Hour).UnixNano())
	})
}

func Test_Truncate(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp = timeTemp.Truncate(time.Hour)
		t.Assert(timeTemp.UnixNano(), timeTemp1.Truncate(time.Hour).UnixNano())
	})
}

func Test_StartOfMinute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfMinute()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 18:24:00")
	})
}

func Test_EndOfMinute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfMinute()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 18:24:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfMinute(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 18:24:59.999")
	})
}

func Test_StartOfHour(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfHour()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 18:00:00")
	})
}

func Test_EndOfHour(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfHour()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 18:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfHour(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 18:59:59.999")
	})
}

func Test_StartOfDay(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfDay()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-12 00:00:00")
	})
}

func Test_EndOfDay(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfDay()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 23:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfDay(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 23:59:59.999")
	})
}

func Test_StartOfWeek(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfWeek()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-06 00:00:00")
	})
}

func Test_EndOfWeek(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfWeek()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 23:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfWeek(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-12 23:59:59.999")
	})
}

func Test_StartOfMonth(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.StartOfMonth()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-12-01 00:00:00")
	})
}

func Test_EndOfMonth(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfMonth()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-12 18:24:06")
		timeTemp1 := timeTemp.EndOfMonth(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.999")
	})
}

func Test_StartOfQuarter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.StartOfQuarter()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-10-01 00:00:00")
	})
}

func Test_EndOfQuarter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfQuarter()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfQuarter(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.999")
	})
}

func Test_StartOfHalf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.StartOfHalf()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-07-01 00:00:00")
	})
}

func Test_EndOfHalf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfHalf()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfHalf(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.999")
	})
}

func Test_StartOfYear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.StartOfYear()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s"), "2020-01-01 00:00:00")
	})
}

func Test_EndOfYear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfYear()
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.000")
	})
	gtest.C(t, func(t *gtest.T) {
		timeTemp := gtime.NewFromStr("2020-12-06 18:24:06")
		timeTemp1 := timeTemp.EndOfYear(true)
		t.Assert(timeTemp1.Format("Y-m-d H:i:s.u"), "2020-12-31 23:59:59.999")
	})
}

func Test_OnlyTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		obj := gtime.NewFromStr("18:24:06")
		t.Assert(obj.String(), "18:24:06")
	})
}

func Test_DeepCopy(t *testing.T) {
	type User struct {
		Id          int
		CreatedTime *gtime.Time
	}
	gtest.C(t, func(t *gtest.T) {
		u1 := &User{
			Id:          1,
			CreatedTime: gtime.New("2022-03-08T03:01:14+08:00"),
		}
		u2 := gutil.Copy(u1).(*User)
		t.Assert(u1, u2)
	})
	// nil attribute.
	gtest.C(t, func(t *gtest.T) {
		u1 := &User{}
		u2 := gutil.Copy(u1).(*User)
		t.Assert(u1, u2)
	})
	gtest.C(t, func(t *gtest.T) {
		var t1 *gtime.Time = nil
		t.Assert(t1.DeepCopy(), nil)
	})
}

func Test_UnmarshalJSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var t1 gtime.Time
		t.AssertNE(json.Unmarshal([]byte("{}"), &t1), nil)
	})
}

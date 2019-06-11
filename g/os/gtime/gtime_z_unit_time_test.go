package gtime_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_New(t *testing.T) {
	gtest.Case(t, func() {
		timeNow := time.Now()
		timeTemp := gtime.New(timeNow)
		gtest.Assert(timeTemp.Time.UnixNano(), timeNow.UnixNano())

		timeTemp1 := gtime.New()
		gtest.Assert(timeTemp1.Time, time.Time{})
	})
}

func Test_NewFromStr(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.NewFromStr("2006-01-02 15:04:05")
		gtest.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.NewFromStr("20060102")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_NewFromStrFormat(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.NewFromStrFormat("2006-01-02 15:04:05", "Y-m-d H:i:s")
		gtest.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.NewFromStrFormat("2006-01-02 15:04:05", "aabbcc")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_NewFromStrLayout(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.NewFromStrLayout("2006-01-02 15:04:05", "2006-01-02 15:04:05")
		gtest.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:05")

		timeTemp1 := gtime.NewFromStrLayout("2006-01-02 15:04:05", "aabbcc")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_NewFromTimeStamp(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.NewFromTimeStamp(1554459846000)
		gtest.Assert(timeTemp.Format("Y-m-d H:i:s"), "2019-04-05 18:24:06")
		timeTemp1 := gtime.NewFromTimeStamp(0)
		gtest.Assert(timeTemp1.Format("Y-m-d H:i:s"), "0001-01-01 00:00:00")
	})
}

func Test_tSecond(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		gtest.Assert(timeTemp.Second(), timeTemp.Time.Unix())
	})
}

func Test_tNanosecond(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		gtest.Assert(timeTemp.Nanosecond(), timeTemp.Time.UnixNano())
	})
}

func Test_tMicrosecond(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		gtest.Assert(timeTemp.Microsecond(), timeTemp.Time.UnixNano()/1e3)
	})
}

func Test_tMillisecond(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		gtest.Assert(timeTemp.Millisecond(), timeTemp.Time.UnixNano()/1e6)
	})
}

func Test_String(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		gtest.Assert(timeTemp.String(), timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
}

func Test_Clone(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Clone()
		gtest.Assert(timeTemp.Time.Unix(), timeTemp1.Time.Unix())
	})
}

func Test_ToTime(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		gtest.Assert(timeTemp.ToTime().UnixNano(), timeTemp1.UnixNano())
	})
}

func Test_Add(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.NewFromStr("2006-01-02 15:04:05")
		timeTemp.Add(time.Second)
		gtest.Assert(timeTemp.Format("Y-m-d H:i:s"), "2006-01-02 15:04:06")
	})
}

func Test_ToZone(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		//
		timeTemp.ToZone("America/Los_Angeles")
		gtest.Assert(timeTemp.Time.Location().String(), "America/Los_Angeles")

		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			t.Error("test fail")
		}
		timeTemp.ToLocation(loc)
		gtest.Assert(timeTemp.Time.Location().String(), "Asia/Shanghai")

		timeTemp1, _ := timeTemp.ToZone("errZone")
		if timeTemp1 != nil {
			t.Error("test fail")
		}
	})
}

func Test_AddDate(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.NewFromStr("2006-01-02 15:04:05")
		timeTemp.AddDate(1, 2, 3)
		gtest.Assert(timeTemp.Format("Y-m-d H:i:s"), "2007-03-05 15:04:05")
	})
}

func Test_UTC(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.UTC()
		gtest.Assert(timeTemp.UnixNano(), timeTemp1.UTC().UnixNano())
	})
}

func Test_Local(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.Local()
		gtest.Assert(timeTemp.UnixNano(), timeTemp1.Local().UnixNano())
	})
}

func Test_Round(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.Round(time.Hour)
		gtest.Assert(timeTemp.UnixNano(), timeTemp1.Round(time.Hour).UnixNano())
	})
}

func Test_Truncate(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		timeTemp1 := timeTemp.Time
		timeTemp.Truncate(time.Hour)
		gtest.Assert(timeTemp.UnixNano(), timeTemp1.Truncate(time.Hour).UnixNano())
	})
}

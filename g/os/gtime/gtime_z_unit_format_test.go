package gtime_test

import (
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func Test_Format(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp, err := gtime.StrToTime("2006-01-11 15:04:05", "Y-m-d H:i:s")
		timeTemp.ToZone("Asia/Shanghai")
		if err != nil {
			t.Error("test fail")
		}
		gtest.Assert(timeTemp.Format("\\T\\i\\m\\e中文Y-m-j G:i:s.u\\"), "Time中文2006-01-11 15:04:05.000")

		gtest.Assert(timeTemp.Format("d D j l"), "11 Wed 11 Wednesday")

		gtest.Assert(timeTemp.Format("F m M n"), "January 01 Jan 1")

		gtest.Assert(timeTemp.Format("Y y"), "2006 06")

		gtest.Assert(timeTemp.Format("a A g G h H i s u .u"), "pm PM 3 15 03 15 04 05 000 .000")

		gtest.Assert(timeTemp.Format("O P T"), "+0800 +08:00 CST")

		gtest.Assert(timeTemp.Format("r"), "Wed, 11 Jan 06 15:04 CST")

		gtest.Assert(timeTemp.Format("c"), "2006-01-11T15:04:05+08:00")

		//补零
		timeTemp1, err := gtime.StrToTime("2006-01-02 03:04:05", "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		gtest.Assert(timeTemp1.Format("Y-m-d h:i:s"), "2006-01-02 03:04:05")
		//不补零
		timeTemp2, err := gtime.StrToTime("2006-01-02 03:04:05", "Y-m-d H:i:s")
		if err != nil {
			t.Error("test fail")
		}
		gtest.Assert(timeTemp2.Format("Y-n-j G:i:s"), "2006-1-2 3:04:05")

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
			{"k": "2019-03-05", "f": "S", "r": "th"},
			{"k": "2019-01-01", "f": "第z天", "r": "第0天"},
			{"k": "2019-01-05", "f": "第z天", "r": "第4天"},
			{"k": "2020-05-05", "f": "第z天", "r": "第125天"},
			{"k": "2020-12-31", "f": "第z天", "r": "第365天"}, // 润年
			{"k": "2020-02-12", "f": "第z天", "r": "第42天"}, // 润年
			// 测试一个日期在当年是第多少周
			// 2019-01-01为星期2
			// 2017-01-01为星期日
			{"k": "2019-06-04", "f": "第W周", "r": "第23周"}, // 06-04为星期2
			{"k": "2019-05-20", "f": "第W周", "r": "第21周"}, // 05-20为星期1
			{"k": "2017-07-26", "f": "第W周", "r": "第30周"}, // 05-20为星期1
			{"k": "2017-01-02", "f": "第W周", "r": "第2周"},  // 1号是星期天为最后一天，2号为第二周的周一，每周从周一开始,但php算得有问题。
			{"k": "2015-01-02", "f": "第W周", "r": "第1周"},  //1号星期4

		}

		for _, v := range times {
			t1, err1 := gtime.StrToTime(v["k"], "Y-m-d")
			gtest.Assert(err1, nil)
			gtest.Assert(t1.Format(v["f"]), v["r"])
		}







	})
}

func Test_Layout(t *testing.T) {
	gtest.Case(t, func() {
		timeTemp := gtime.Now()
		gtest.Assert(timeTemp.Layout("2006-01-02 15:04:05"), timeTemp.Time.Format("2006-01-02 15:04:05"))
	})
}

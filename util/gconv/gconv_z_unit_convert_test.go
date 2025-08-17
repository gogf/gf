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

func TestConvert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Convert(true, "bool"), true)
		t.AssertEQ(gconv.Convert(false, "bool"), false)

		t.Assert(gconv.Convert(int(0), "int"), int(0))
		t.Assert(gconv.Convert(int8(0), "int8"), int8(0))
		t.Assert(gconv.Convert(int16(0), "int16"), int16(0))
		t.Assert(gconv.Convert(int32(0), "int32"), int32(0))
		t.Assert(gconv.Convert(int64(1), "int64"), int64(1))

		t.Assert(gconv.Convert(uint(0), "uint"), uint(0))
		t.Assert(gconv.Convert(uint8(0), "uint8"), uint8(0))
		t.Assert(gconv.Convert(uint16(0), "uint16"), uint16(0))
		t.Assert(gconv.Convert(uint32(0), "uint32"), uint32(0))
		t.Assert(gconv.Convert(uint64(0), "uint64"), uint64(0))

		t.Assert(gconv.Convert(float32(0), "float32"), float32(0))
		t.Assert(gconv.Convert(float64(0), "float64"), float64(0))

		t.AssertEQ(gconv.Convert([]int{1, 2}, "[]int"), []int{1, 2})
		t.AssertEQ(gconv.Convert([]int8{1, 2}, "[]int8}"), []int8{1, 2})
		t.AssertEQ(gconv.Convert([]int16{1, 2}, "[]int16"), []int16{1, 2})
		t.AssertEQ(gconv.Convert([]int32{1, 2}, "[]int32"), []int32{1, 2})
		t.AssertEQ(gconv.Convert([]int64{1, 2}, "[]int64"), []int64{1, 2})
		t.AssertEQ(gconv.Convert([]uint{1, 2}, "[]uint"), []uint{1, 2})
		t.AssertEQ(gconv.Convert([]uint8{1, 2}, "[]uint8}"), []uint8{1, 2})
		t.AssertEQ(gconv.Convert([]uint16{1, 2}, "[]uint16"), []uint16{1, 2})
		t.AssertEQ(gconv.Convert([]uint32{1, 2}, "[]uint32"), []uint32{1, 2})
		t.AssertEQ(gconv.Convert([]uint64{1, 2}, "[]uint64"), []uint64{1, 2})
		t.AssertEQ(gconv.Convert([]float32{1, 2}, "[]float32"), []float32{1, 2})
		t.AssertEQ(gconv.Convert([]float64{1, 2}, "[]float64"), []float64{1, 2})
		t.AssertEQ(gconv.Convert([]string{"1", "2"}, "[]string"), []string{"1", "2"})
		t.AssertEQ(gconv.Convert([]byte{}, "[]byte"), []uint8{})

		var anyTest interface{} = nil
		t.AssertEQ(gconv.Convert(anyTest, "string"), "")
		t.AssertEQ(gconv.Convert("1", "string"), "1")

		t.AssertEQ(gconv.Convert("1989-01-02", "Time", "Y-m-d"),
			gconv.Time("1989-01-02", "Y-m-d"))

		t.AssertEQ(gconv.Convert(1989, "Time"),
			gconv.Time("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(gconv.Convert(1989, "gtime.Time"),
			*gconv.GTime("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(gconv.Convert(1989, "*gtime.Time"),
			gconv.GTime(1989))
		t.AssertEQ(gconv.Convert(1989, "Duration"),
			time.Duration(int64(1989)))
		t.AssertEQ(gconv.Convert("1989", "Duration"),
			time.Duration(int64(1989)))
		t.AssertEQ(gconv.Convert("1989", ""),
			"1989")

		// TODO gconv.Convert(gtime.Now(), "gtime.Time", 1) = {{0001-01-01 00:00:00 +0000 UTC}}
		t.AssertEQ(gconv.Convert(gtime.Now(), "gtime.Time", 1), *gtime.New())
		t.AssertEQ(gconv.Convert(gtime.Now(), "*gtime.Time", 1), gtime.New())
		t.AssertEQ(gconv.Convert(gtime.Now(), "GTime", 1), *gtime.New())

		var boolValue bool = true
		t.Assert(gconv.Convert(boolValue, "*bool"), true)
		t.Assert(gconv.Convert(&boolValue, "*bool"), true)

		var intNum int = 1
		t.Assert(gconv.Convert(intNum, "*int"), int(1))
		t.Assert(gconv.Convert(&intNum, "*int"), int(1))
		var int8Num int8 = 1
		t.Assert(gconv.Convert(int8Num, "*int8"), int(1))
		t.Assert(gconv.Convert(&int8Num, "*int8"), int(1))
		var int16Num int16 = 1
		t.Assert(gconv.Convert(int16Num, "*int16"), int(1))
		t.Assert(gconv.Convert(&int16Num, "*int16"), int(1))
		var int32Num int32 = 1
		t.Assert(gconv.Convert(int32Num, "*int32"), int(1))
		t.Assert(gconv.Convert(&int32Num, "*int32"), int(1))
		var int64Num int64 = 1
		t.Assert(gconv.Convert(int64Num, "*int64"), int(1))
		t.Assert(gconv.Convert(&int64Num, "*int64"), int(1))

		var uintNum uint = 1
		t.Assert(gconv.Convert(&uintNum, "*uint"), int(1))
		var uint8Num uint8 = 1
		t.Assert(gconv.Convert(uint8Num, "*uint8"), int(1))
		t.Assert(gconv.Convert(&uint8Num, "*uint8"), int(1))
		var uint16Num uint16 = 1
		t.Assert(gconv.Convert(uint16Num, "*uint16"), int(1))
		t.Assert(gconv.Convert(&uint16Num, "*uint16"), int(1))
		var uint32Num uint32 = 1
		t.Assert(gconv.Convert(uint32Num, "*uint32"), int(1))
		t.Assert(gconv.Convert(&uint32Num, "*uint32"), int(1))
		var uint64Num uint64 = 1
		t.Assert(gconv.Convert(uint64Num, "*uint64"), int(1))
		t.Assert(gconv.Convert(&uint64Num, "*uint64"), int(1))

		var float32Num float32 = 1.1
		t.Assert(gconv.Convert(float32Num, "*float32"), float32(1.1))
		t.Assert(gconv.Convert(&float32Num, "*float32"), float32(1.1))

		var float64Num float64 = 1.1
		t.Assert(gconv.Convert(float64Num, "*float64"), float64(1.1))
		t.Assert(gconv.Convert(&float64Num, "*float64"), float64(1.1))

		var stringValue string = "1"
		t.Assert(gconv.Convert(stringValue, "*string"), "1")
		t.Assert(gconv.Convert(&stringValue, "*string"), "1")

		var durationValue time.Duration = 1989
		var expectDurationValue = time.Duration(int64(1989))
		t.AssertEQ(gconv.Convert(&durationValue, "*time.Duration"),
			&expectDurationValue)
		t.AssertEQ(gconv.Convert(durationValue, "*time.Duration"),
			&expectDurationValue)

		var mapStrInt = map[string]int{"k1": 1}
		var mapStrStr = map[string]string{"k1": "1"}
		var mapStrAny = map[string]interface{}{"k1": 1}
		t.AssertEQ(gconv.Convert(mapStrInt, "map[string]string"), mapStrStr)
		t.AssertEQ(gconv.Convert(mapStrInt, "map[string]interface{}"), mapStrAny)
	})
}

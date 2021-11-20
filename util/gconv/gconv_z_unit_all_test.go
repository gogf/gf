// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type iString interface {
	String() string
}
type S struct {
}

func (s S) String() string {
	return "22222"
}

type iError interface {
	Error() string
}
type S1 struct {
}

func (s1 S1) Error() string {
	return "22222"
}

func Test_Bool_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.AssertEQ(gconv.Bool(any), false)
		t.AssertEQ(gconv.Bool(false), false)
		t.AssertEQ(gconv.Bool(nil), false)
		t.AssertEQ(gconv.Bool(0), false)
		t.AssertEQ(gconv.Bool("0"), false)
		t.AssertEQ(gconv.Bool(""), false)
		t.AssertEQ(gconv.Bool("false"), false)
		t.AssertEQ(gconv.Bool("off"), false)
		t.AssertEQ(gconv.Bool([]byte{}), false)
		t.AssertEQ(gconv.Bool([]string{}), false)
		t.AssertEQ(gconv.Bool([2]int{1, 2}), true)
		t.AssertEQ(gconv.Bool([]interface{}{}), false)
		t.AssertEQ(gconv.Bool([]map[int]int{}), false)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Bool(countryCapitalMap), true)

		t.AssertEQ(gconv.Bool("1"), true)
		t.AssertEQ(gconv.Bool("on"), true)
		t.AssertEQ(gconv.Bool(1), true)
		t.AssertEQ(gconv.Bool(123.456), true)
		t.AssertEQ(gconv.Bool(boolStruct{}), true)
		t.AssertEQ(gconv.Bool(&boolStruct{}), true)
	})
}

func Test_Int_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.AssertEQ(gconv.Int(any), 0)
		t.AssertEQ(gconv.Int(false), 0)
		t.AssertEQ(gconv.Int(nil), 0)
		t.Assert(gconv.Int(nil), 0)
		t.AssertEQ(gconv.Int(0), 0)
		t.AssertEQ(gconv.Int("0"), 0)
		t.AssertEQ(gconv.Int(""), 0)
		t.AssertEQ(gconv.Int("false"), 0)
		t.AssertEQ(gconv.Int("off"), 0)
		t.AssertEQ(gconv.Int([]byte{}), 0)
		t.AssertEQ(gconv.Int([]string{}), 0)
		t.AssertEQ(gconv.Int([2]int{1, 2}), 0)
		t.AssertEQ(gconv.Int([]interface{}{}), 0)
		t.AssertEQ(gconv.Int([]map[int]int{}), 0)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Int(countryCapitalMap), 0)

		t.AssertEQ(gconv.Int("1"), 1)
		t.AssertEQ(gconv.Int("on"), 0)
		t.AssertEQ(gconv.Int(1), 1)
		t.AssertEQ(gconv.Int(123.456), 123)
		t.AssertEQ(gconv.Int(boolStruct{}), 0)
		t.AssertEQ(gconv.Int(&boolStruct{}), 0)
	})
}

func Test_Int8_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Int8(any), int8(0))
		t.AssertEQ(gconv.Int8(false), int8(0))
		t.AssertEQ(gconv.Int8(nil), int8(0))
		t.AssertEQ(gconv.Int8(0), int8(0))
		t.AssertEQ(gconv.Int8("0"), int8(0))
		t.AssertEQ(gconv.Int8(""), int8(0))
		t.AssertEQ(gconv.Int8("false"), int8(0))
		t.AssertEQ(gconv.Int8("off"), int8(0))
		t.AssertEQ(gconv.Int8([]byte{}), int8(0))
		t.AssertEQ(gconv.Int8([]string{}), int8(0))
		t.AssertEQ(gconv.Int8([2]int{1, 2}), int8(0))
		t.AssertEQ(gconv.Int8([]interface{}{}), int8(0))
		t.AssertEQ(gconv.Int8([]map[int]int{}), int8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Int8(countryCapitalMap), int8(0))

		t.AssertEQ(gconv.Int8("1"), int8(1))
		t.AssertEQ(gconv.Int8("on"), int8(0))
		t.AssertEQ(gconv.Int8(int8(1)), int8(1))
		t.AssertEQ(gconv.Int8(123.456), int8(123))
		t.AssertEQ(gconv.Int8(boolStruct{}), int8(0))
		t.AssertEQ(gconv.Int8(&boolStruct{}), int8(0))
	})
}

func Test_Int16_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Int16(any), int16(0))
		t.AssertEQ(gconv.Int16(false), int16(0))
		t.AssertEQ(gconv.Int16(nil), int16(0))
		t.AssertEQ(gconv.Int16(0), int16(0))
		t.AssertEQ(gconv.Int16("0"), int16(0))
		t.AssertEQ(gconv.Int16(""), int16(0))
		t.AssertEQ(gconv.Int16("false"), int16(0))
		t.AssertEQ(gconv.Int16("off"), int16(0))
		t.AssertEQ(gconv.Int16([]byte{}), int16(0))
		t.AssertEQ(gconv.Int16([]string{}), int16(0))
		t.AssertEQ(gconv.Int16([2]int{1, 2}), int16(0))
		t.AssertEQ(gconv.Int16([]interface{}{}), int16(0))
		t.AssertEQ(gconv.Int16([]map[int]int{}), int16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Int16(countryCapitalMap), int16(0))

		t.AssertEQ(gconv.Int16("1"), int16(1))
		t.AssertEQ(gconv.Int16("on"), int16(0))
		t.AssertEQ(gconv.Int16(int16(1)), int16(1))
		t.AssertEQ(gconv.Int16(123.456), int16(123))
		t.AssertEQ(gconv.Int16(boolStruct{}), int16(0))
		t.AssertEQ(gconv.Int16(&boolStruct{}), int16(0))
	})
}

func Test_Int32_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Int32(any), int32(0))
		t.AssertEQ(gconv.Int32(false), int32(0))
		t.AssertEQ(gconv.Int32(nil), int32(0))
		t.AssertEQ(gconv.Int32(0), int32(0))
		t.AssertEQ(gconv.Int32("0"), int32(0))
		t.AssertEQ(gconv.Int32(""), int32(0))
		t.AssertEQ(gconv.Int32("false"), int32(0))
		t.AssertEQ(gconv.Int32("off"), int32(0))
		t.AssertEQ(gconv.Int32([]byte{}), int32(0))
		t.AssertEQ(gconv.Int32([]string{}), int32(0))
		t.AssertEQ(gconv.Int32([2]int{1, 2}), int32(0))
		t.AssertEQ(gconv.Int32([]interface{}{}), int32(0))
		t.AssertEQ(gconv.Int32([]map[int]int{}), int32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Int32(countryCapitalMap), int32(0))

		t.AssertEQ(gconv.Int32("1"), int32(1))
		t.AssertEQ(gconv.Int32("on"), int32(0))
		t.AssertEQ(gconv.Int32(int32(1)), int32(1))
		t.AssertEQ(gconv.Int32(123.456), int32(123))
		t.AssertEQ(gconv.Int32(boolStruct{}), int32(0))
		t.AssertEQ(gconv.Int32(&boolStruct{}), int32(0))
	})
}

func Test_Int64_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.AssertEQ(gconv.Int64("0x00e"), int64(14))
		t.Assert(gconv.Int64("022"), int64(18))

		t.Assert(gconv.Int64(any), int64(0))
		t.Assert(gconv.Int64(true), 1)
		t.Assert(gconv.Int64("1"), int64(1))
		t.Assert(gconv.Int64("0"), int64(0))
		t.Assert(gconv.Int64("X"), int64(0))
		t.Assert(gconv.Int64("x"), int64(0))
		t.Assert(gconv.Int64(int64(1)), int64(1))
		t.Assert(gconv.Int64(int(0)), int64(0))
		t.Assert(gconv.Int64(int8(0)), int64(0))
		t.Assert(gconv.Int64(int16(0)), int64(0))
		t.Assert(gconv.Int64(int32(0)), int64(0))
		t.Assert(gconv.Int64(uint64(0)), int64(0))
		t.Assert(gconv.Int64(uint32(0)), int64(0))
		t.Assert(gconv.Int64(uint16(0)), int64(0))
		t.Assert(gconv.Int64(uint8(0)), int64(0))
		t.Assert(gconv.Int64(uint(0)), int64(0))
		t.Assert(gconv.Int64(float32(0)), int64(0))

		t.AssertEQ(gconv.Int64(false), int64(0))
		t.AssertEQ(gconv.Int64(nil), int64(0))
		t.AssertEQ(gconv.Int64(0), int64(0))
		t.AssertEQ(gconv.Int64("0"), int64(0))
		t.AssertEQ(gconv.Int64(""), int64(0))
		t.AssertEQ(gconv.Int64("false"), int64(0))
		t.AssertEQ(gconv.Int64("off"), int64(0))
		t.AssertEQ(gconv.Int64([]byte{}), int64(0))
		t.AssertEQ(gconv.Int64([]string{}), int64(0))
		t.AssertEQ(gconv.Int64([2]int{1, 2}), int64(0))
		t.AssertEQ(gconv.Int64([]interface{}{}), int64(0))
		t.AssertEQ(gconv.Int64([]map[int]int{}), int64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Int64(countryCapitalMap), int64(0))

		t.AssertEQ(gconv.Int64("1"), int64(1))
		t.AssertEQ(gconv.Int64("on"), int64(0))
		t.AssertEQ(gconv.Int64(int64(1)), int64(1))
		t.AssertEQ(gconv.Int64(123.456), int64(123))
		t.AssertEQ(gconv.Int64(boolStruct{}), int64(0))
		t.AssertEQ(gconv.Int64(&boolStruct{}), int64(0))
	})
}

func Test_Uint_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.AssertEQ(gconv.Uint(any), uint(0))
		t.AssertEQ(gconv.Uint(false), uint(0))
		t.AssertEQ(gconv.Uint(nil), uint(0))
		t.Assert(gconv.Uint(nil), uint(0))
		t.AssertEQ(gconv.Uint(uint(0)), uint(0))
		t.AssertEQ(gconv.Uint("0"), uint(0))
		t.AssertEQ(gconv.Uint(""), uint(0))
		t.AssertEQ(gconv.Uint("false"), uint(0))
		t.AssertEQ(gconv.Uint("off"), uint(0))
		t.AssertEQ(gconv.Uint([]byte{}), uint(0))
		t.AssertEQ(gconv.Uint([]string{}), uint(0))
		t.AssertEQ(gconv.Uint([2]int{1, 2}), uint(0))
		t.AssertEQ(gconv.Uint([]interface{}{}), uint(0))
		t.AssertEQ(gconv.Uint([]map[int]int{}), uint(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Uint(countryCapitalMap), uint(0))

		t.AssertEQ(gconv.Uint("1"), uint(1))
		t.AssertEQ(gconv.Uint("on"), uint(0))
		t.AssertEQ(gconv.Uint(1), uint(1))
		t.AssertEQ(gconv.Uint(123.456), uint(123))
		t.AssertEQ(gconv.Uint(boolStruct{}), uint(0))
		t.AssertEQ(gconv.Uint(&boolStruct{}), uint(0))
	})
}

func Test_Uint8_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Uint8(any), uint8(0))
		t.AssertEQ(gconv.Uint8(uint8(1)), uint8(1))
		t.AssertEQ(gconv.Uint8(false), uint8(0))
		t.AssertEQ(gconv.Uint8(nil), uint8(0))
		t.AssertEQ(gconv.Uint8(0), uint8(0))
		t.AssertEQ(gconv.Uint8("0"), uint8(0))
		t.AssertEQ(gconv.Uint8(""), uint8(0))
		t.AssertEQ(gconv.Uint8("false"), uint8(0))
		t.AssertEQ(gconv.Uint8("off"), uint8(0))
		t.AssertEQ(gconv.Uint8([]byte{}), uint8(0))
		t.AssertEQ(gconv.Uint8([]string{}), uint8(0))
		t.AssertEQ(gconv.Uint8([2]int{1, 2}), uint8(0))
		t.AssertEQ(gconv.Uint8([]interface{}{}), uint8(0))
		t.AssertEQ(gconv.Uint8([]map[int]int{}), uint8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Uint8(countryCapitalMap), uint8(0))

		t.AssertEQ(gconv.Uint8("1"), uint8(1))
		t.AssertEQ(gconv.Uint8("on"), uint8(0))
		t.AssertEQ(gconv.Uint8(int8(1)), uint8(1))
		t.AssertEQ(gconv.Uint8(123.456), uint8(123))
		t.AssertEQ(gconv.Uint8(boolStruct{}), uint8(0))
		t.AssertEQ(gconv.Uint8(&boolStruct{}), uint8(0))
	})
}

func Test_Uint16_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Uint16(any), uint16(0))
		t.AssertEQ(gconv.Uint16(uint16(1)), uint16(1))
		t.AssertEQ(gconv.Uint16(false), uint16(0))
		t.AssertEQ(gconv.Uint16(nil), uint16(0))
		t.AssertEQ(gconv.Uint16(0), uint16(0))
		t.AssertEQ(gconv.Uint16("0"), uint16(0))
		t.AssertEQ(gconv.Uint16(""), uint16(0))
		t.AssertEQ(gconv.Uint16("false"), uint16(0))
		t.AssertEQ(gconv.Uint16("off"), uint16(0))
		t.AssertEQ(gconv.Uint16([]byte{}), uint16(0))
		t.AssertEQ(gconv.Uint16([]string{}), uint16(0))
		t.AssertEQ(gconv.Uint16([2]int{1, 2}), uint16(0))
		t.AssertEQ(gconv.Uint16([]interface{}{}), uint16(0))
		t.AssertEQ(gconv.Uint16([]map[int]int{}), uint16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Uint16(countryCapitalMap), uint16(0))

		t.AssertEQ(gconv.Uint16("1"), uint16(1))
		t.AssertEQ(gconv.Uint16("on"), uint16(0))
		t.AssertEQ(gconv.Uint16(int16(1)), uint16(1))
		t.AssertEQ(gconv.Uint16(123.456), uint16(123))
		t.AssertEQ(gconv.Uint16(boolStruct{}), uint16(0))
		t.AssertEQ(gconv.Uint16(&boolStruct{}), uint16(0))
	})
}

func Test_Uint32_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Uint32(any), uint32(0))
		t.AssertEQ(gconv.Uint32(uint32(1)), uint32(1))
		t.AssertEQ(gconv.Uint32(false), uint32(0))
		t.AssertEQ(gconv.Uint32(nil), uint32(0))
		t.AssertEQ(gconv.Uint32(0), uint32(0))
		t.AssertEQ(gconv.Uint32("0"), uint32(0))
		t.AssertEQ(gconv.Uint32(""), uint32(0))
		t.AssertEQ(gconv.Uint32("false"), uint32(0))
		t.AssertEQ(gconv.Uint32("off"), uint32(0))
		t.AssertEQ(gconv.Uint32([]byte{}), uint32(0))
		t.AssertEQ(gconv.Uint32([]string{}), uint32(0))
		t.AssertEQ(gconv.Uint32([2]int{1, 2}), uint32(0))
		t.AssertEQ(gconv.Uint32([]interface{}{}), uint32(0))
		t.AssertEQ(gconv.Uint32([]map[int]int{}), uint32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Uint32(countryCapitalMap), uint32(0))

		t.AssertEQ(gconv.Uint32("1"), uint32(1))
		t.AssertEQ(gconv.Uint32("on"), uint32(0))
		t.AssertEQ(gconv.Uint32(int32(1)), uint32(1))
		t.AssertEQ(gconv.Uint32(123.456), uint32(123))
		t.AssertEQ(gconv.Uint32(boolStruct{}), uint32(0))
		t.AssertEQ(gconv.Uint32(&boolStruct{}), uint32(0))
	})
}

func Test_Uint64_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.AssertEQ(gconv.Uint64("0x00e"), uint64(14))
		t.Assert(gconv.Uint64("022"), uint64(18))

		t.AssertEQ(gconv.Uint64(any), uint64(0))
		t.AssertEQ(gconv.Uint64(true), uint64(1))
		t.Assert(gconv.Uint64("1"), int64(1))
		t.Assert(gconv.Uint64("0"), uint64(0))
		t.Assert(gconv.Uint64("X"), uint64(0))
		t.Assert(gconv.Uint64("x"), uint64(0))
		t.Assert(gconv.Uint64(int64(1)), uint64(1))
		t.Assert(gconv.Uint64(int(0)), uint64(0))
		t.Assert(gconv.Uint64(int8(0)), uint64(0))
		t.Assert(gconv.Uint64(int16(0)), uint64(0))
		t.Assert(gconv.Uint64(int32(0)), uint64(0))
		t.Assert(gconv.Uint64(uint64(0)), uint64(0))
		t.Assert(gconv.Uint64(uint32(0)), uint64(0))
		t.Assert(gconv.Uint64(uint16(0)), uint64(0))
		t.Assert(gconv.Uint64(uint8(0)), uint64(0))
		t.Assert(gconv.Uint64(uint(0)), uint64(0))
		t.Assert(gconv.Uint64(float32(0)), uint64(0))

		t.AssertEQ(gconv.Uint64(false), uint64(0))
		t.AssertEQ(gconv.Uint64(nil), uint64(0))
		t.AssertEQ(gconv.Uint64(0), uint64(0))
		t.AssertEQ(gconv.Uint64("0"), uint64(0))
		t.AssertEQ(gconv.Uint64(""), uint64(0))
		t.AssertEQ(gconv.Uint64("false"), uint64(0))
		t.AssertEQ(gconv.Uint64("off"), uint64(0))
		t.AssertEQ(gconv.Uint64([]byte{}), uint64(0))
		t.AssertEQ(gconv.Uint64([]string{}), uint64(0))
		t.AssertEQ(gconv.Uint64([2]int{1, 2}), uint64(0))
		t.AssertEQ(gconv.Uint64([]interface{}{}), uint64(0))
		t.AssertEQ(gconv.Uint64([]map[int]int{}), uint64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Uint64(countryCapitalMap), uint64(0))

		t.AssertEQ(gconv.Uint64("1"), uint64(1))
		t.AssertEQ(gconv.Uint64("on"), uint64(0))
		t.AssertEQ(gconv.Uint64(int64(1)), uint64(1))
		t.AssertEQ(gconv.Uint64(123.456), uint64(123))
		t.AssertEQ(gconv.Uint64(boolStruct{}), uint64(0))
		t.AssertEQ(gconv.Uint64(&boolStruct{}), uint64(0))
	})
}

func Test_Float32_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Float32(any), float32(0))
		t.AssertEQ(gconv.Float32(false), float32(0))
		t.AssertEQ(gconv.Float32(nil), float32(0))
		t.AssertEQ(gconv.Float32(0), float32(0))
		t.AssertEQ(gconv.Float32("0"), float32(0))
		t.AssertEQ(gconv.Float32(""), float32(0))
		t.AssertEQ(gconv.Float32("false"), float32(0))
		t.AssertEQ(gconv.Float32("off"), float32(0))
		t.AssertEQ(gconv.Float32([]byte{}), float32(0))
		t.AssertEQ(gconv.Float32([]string{}), float32(0))
		t.AssertEQ(gconv.Float32([2]int{1, 2}), float32(0))
		t.AssertEQ(gconv.Float32([]interface{}{}), float32(0))
		t.AssertEQ(gconv.Float32([]map[int]int{}), float32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Float32(countryCapitalMap), float32(0))

		t.AssertEQ(gconv.Float32("1"), float32(1))
		t.AssertEQ(gconv.Float32("on"), float32(0))
		t.AssertEQ(gconv.Float32(float32(1)), float32(1))
		t.AssertEQ(gconv.Float32(123.456), float32(123.456))
		t.AssertEQ(gconv.Float32(boolStruct{}), float32(0))
		t.AssertEQ(gconv.Float32(&boolStruct{}), float32(0))
	})
}

func Test_Float64_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.Assert(gconv.Float64(any), float64(0))
		t.AssertEQ(gconv.Float64(false), float64(0))
		t.AssertEQ(gconv.Float64(nil), float64(0))
		t.AssertEQ(gconv.Float64(0), float64(0))
		t.AssertEQ(gconv.Float64("0"), float64(0))
		t.AssertEQ(gconv.Float64(""), float64(0))
		t.AssertEQ(gconv.Float64("false"), float64(0))
		t.AssertEQ(gconv.Float64("off"), float64(0))
		t.AssertEQ(gconv.Float64([]byte{}), float64(0))
		t.AssertEQ(gconv.Float64([]string{}), float64(0))
		t.AssertEQ(gconv.Float64([2]int{1, 2}), float64(0))
		t.AssertEQ(gconv.Float64([]interface{}{}), float64(0))
		t.AssertEQ(gconv.Float64([]map[int]int{}), float64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.Float64(countryCapitalMap), float64(0))

		t.AssertEQ(gconv.Float64("1"), float64(1))
		t.AssertEQ(gconv.Float64("on"), float64(0))
		t.AssertEQ(gconv.Float64(float64(1)), float64(1))
		t.AssertEQ(gconv.Float64(123.456), float64(123.456))
		t.AssertEQ(gconv.Float64(boolStruct{}), float64(0))
		t.AssertEQ(gconv.Float64(&boolStruct{}), float64(0))
	})
}

func Test_String_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s []rune
		t.AssertEQ(gconv.String(s), "")
		var any interface{} = nil
		t.AssertEQ(gconv.String(any), "")
		t.AssertEQ(gconv.String("1"), "1")
		t.AssertEQ(gconv.String("0"), string("0"))
		t.Assert(gconv.String("X"), string("X"))
		t.Assert(gconv.String("x"), string("x"))
		t.Assert(gconv.String(int64(1)), uint64(1))
		t.Assert(gconv.String(int(0)), string("0"))
		t.Assert(gconv.String(int8(0)), string("0"))
		t.Assert(gconv.String(int16(0)), string("0"))
		t.Assert(gconv.String(int32(0)), string("0"))
		t.Assert(gconv.String(uint64(0)), string("0"))
		t.Assert(gconv.String(uint32(0)), string("0"))
		t.Assert(gconv.String(uint16(0)), string("0"))
		t.Assert(gconv.String(uint8(0)), string("0"))
		t.Assert(gconv.String(uint(0)), string("0"))
		t.Assert(gconv.String(float32(0)), string("0"))
		t.AssertEQ(gconv.String(true), "true")
		t.AssertEQ(gconv.String(false), "false")
		t.AssertEQ(gconv.String(nil), "")
		t.AssertEQ(gconv.String(0), string("0"))
		t.AssertEQ(gconv.String("0"), string("0"))
		t.AssertEQ(gconv.String(""), "")
		t.AssertEQ(gconv.String("false"), "false")
		t.AssertEQ(gconv.String("off"), string("off"))
		t.AssertEQ(gconv.String([]byte{}), "")
		t.AssertEQ(gconv.String([]string{}), "[]")
		t.AssertEQ(gconv.String([2]int{1, 2}), "[1,2]")
		t.AssertEQ(gconv.String([]interface{}{}), "[]")
		t.AssertEQ(gconv.String(map[int]int{}), "{}")

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap["France"] = "巴黎"
		countryCapitalMap["Italy"] = "罗马"
		countryCapitalMap["Japan"] = "东京"
		countryCapitalMap["India "] = "新德里"
		t.AssertEQ(gconv.String(countryCapitalMap), `{"France":"巴黎","India ":"新德里","Italy":"罗马","Japan":"东京"}`)
		t.AssertEQ(gconv.String(int64(1)), "1")
		t.AssertEQ(gconv.String(123.456), "123.456")
		t.AssertEQ(gconv.String(boolStruct{}), "{}")
		t.AssertEQ(gconv.String(&boolStruct{}), "{}")

		var info iString
		info = new(S)
		t.AssertEQ(gconv.String(info), "22222")
		var errInfo iError
		errInfo = new(S1)
		t.AssertEQ(gconv.String(errInfo), "22222")
	})
}

func Test_Runes_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Runes("www"), []int32{119, 119, 119})
		var s []rune
		t.AssertEQ(gconv.Runes(s), nil)
	})
}

func Test_Rune_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Rune("www"), int32(0))
		t.AssertEQ(gconv.Rune(int32(0)), int32(0))
		var s []rune
		t.AssertEQ(gconv.Rune(s), int32(0))
	})
}

func Test_Bytes_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Bytes(nil), nil)
		t.AssertEQ(gconv.Bytes(int32(0)), []uint8{0, 0, 0, 0})
		t.AssertEQ(gconv.Bytes("s"), []uint8{115})
		t.AssertEQ(gconv.Bytes([]byte("s")), []uint8{115})
	})
}

func Test_Byte_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Byte(uint8(0)), uint8(0))
		t.AssertEQ(gconv.Byte("s"), uint8(0))
		t.AssertEQ(gconv.Byte([]byte("s")), uint8(115))
	})
}

func Test_Convert_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var any interface{} = nil
		t.AssertEQ(gconv.Convert(any, "string"), "")
		t.AssertEQ(gconv.Convert("1", "string"), "1")
		t.Assert(gconv.Convert(int64(1), "int64"), int64(1))
		t.Assert(gconv.Convert(int(0), "int"), int(0))
		t.Assert(gconv.Convert(int8(0), "int8"), int8(0))
		t.Assert(gconv.Convert(int16(0), "int16"), int16(0))
		t.Assert(gconv.Convert(int32(0), "int32"), int32(0))
		t.Assert(gconv.Convert(uint64(0), "uint64"), uint64(0))
		t.Assert(gconv.Convert(uint32(0), "uint32"), uint32(0))
		t.Assert(gconv.Convert(uint16(0), "uint16"), uint16(0))
		t.Assert(gconv.Convert(uint8(0), "uint8"), uint8(0))
		t.Assert(gconv.Convert(uint(0), "uint"), uint(0))
		t.Assert(gconv.Convert(float32(0), "float32"), float32(0))
		t.Assert(gconv.Convert(float64(0), "float64"), float64(0))
		t.AssertEQ(gconv.Convert(true, "bool"), true)
		t.AssertEQ(gconv.Convert([]byte{}, "[]byte"), []uint8{})
		t.AssertEQ(gconv.Convert([]string{}, "[]string"), []string{})
		t.AssertEQ(gconv.Convert([2]int{1, 2}, "[]int"), []int{1, 2})
		t.AssertEQ(gconv.Convert("1989-01-02", "Time", "Y-m-d"), gconv.Time("1989-01-02", "Y-m-d"))
		t.AssertEQ(gconv.Convert(1989, "Time"), gconv.Time("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(gconv.Convert(gtime.Now(), "gtime.Time", 1), *gtime.New())
		t.AssertEQ(gconv.Convert(1989, "gtime.Time"), *gconv.GTime("1970-01-01 08:33:09 +0800 CST"))
		t.AssertEQ(gconv.Convert(gtime.Now(), "*gtime.Time", 1), gtime.New())
		t.AssertEQ(gconv.Convert(gtime.Now(), "GTime", 1), *gtime.New())
		t.AssertEQ(gconv.Convert(1989, "*gtime.Time"), gconv.GTime(1989))
		t.AssertEQ(gconv.Convert(1989, "Duration"), time.Duration(int64(1989)))
		t.AssertEQ(gconv.Convert("1989", "Duration"), time.Duration(int64(1989)))
		t.AssertEQ(gconv.Convert("1989", ""), "1989")
	})
}

func Test_Slice_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := 123.456
		t.AssertEQ(gconv.Ints(value), []int{123})
		t.AssertEQ(gconv.Ints(nil), nil)
		t.AssertEQ(gconv.Ints([]string{"1", "2"}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]int{}), []int{})
		t.AssertEQ(gconv.Ints([]int8{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]int16{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]int32{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]int64{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]uint{1}), []int{1})
		t.AssertEQ(gconv.Ints([]uint8{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]uint16{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]uint32{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]uint64{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]bool{true}), []int{1})
		t.AssertEQ(gconv.Ints([]float32{1, 2}), []int{1, 2})
		t.AssertEQ(gconv.Ints([]float64{1, 2}), []int{1, 2})
		var inter []interface{} = make([]interface{}, 2)
		t.AssertEQ(gconv.Ints(inter), []int{0, 0})

		t.AssertEQ(gconv.Strings(value), []string{"123.456"})
		t.AssertEQ(gconv.Strings(nil), nil)
		t.AssertEQ(gconv.Strings([]string{"1", "2"}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]int{1}), []string{"1"})
		t.AssertEQ(gconv.Strings([]int8{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]int16{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]int32{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]int64{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]uint{1}), []string{"1"})
		t.AssertEQ(gconv.Strings([]uint8{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]uint16{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]uint32{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]uint64{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]bool{true}), []string{"true"})
		t.AssertEQ(gconv.Strings([]float32{1, 2}), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([]float64{1, 2}), []string{"1", "2"})
		var strer = make([]interface{}, 2)
		t.AssertEQ(gconv.Strings(strer), []string{"", ""})

		t.AssertEQ(gconv.Floats(value), []float64{123.456})
		t.AssertEQ(gconv.Floats(nil), nil)
		t.AssertEQ(gconv.Floats([]string{"1", "2"}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]int{1}), []float64{1})
		t.AssertEQ(gconv.Floats([]int8{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]int16{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]int32{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]int64{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]uint{1}), []float64{1})
		t.AssertEQ(gconv.Floats([]uint8{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]uint16{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]uint32{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]uint64{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]bool{true}), []float64{0})
		t.AssertEQ(gconv.Floats([]float32{1, 2}), []float64{1, 2})
		t.AssertEQ(gconv.Floats([]float64{1, 2}), []float64{1, 2})
		var floer = make([]interface{}, 2)
		t.AssertEQ(gconv.Floats(floer), []float64{0, 0})

		t.AssertEQ(gconv.Interfaces(value), []interface{}{123.456})
		t.AssertEQ(gconv.Interfaces(nil), nil)
		t.AssertEQ(gconv.Interfaces([]interface{}{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]string{"1"}), []interface{}{"1"})
		t.AssertEQ(gconv.Interfaces([]int{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]int8{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]int16{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]int32{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]int64{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]uint{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]uint8{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]uint16{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]uint32{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]uint64{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]bool{true}), []interface{}{true})
		t.AssertEQ(gconv.Interfaces([]float32{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([]float64{1}), []interface{}{1})
		t.AssertEQ(gconv.Interfaces([1]int{1}), []interface{}{1})

		type interSlice []int
		slices := interSlice{1}
		t.AssertEQ(gconv.Interfaces(slices), []interface{}{1})

		t.AssertEQ(gconv.Maps(nil), nil)
		t.AssertEQ(gconv.Maps([]map[string]interface{}{{"a": "1"}}), []map[string]interface{}{{"a": "1"}})
		t.AssertEQ(gconv.Maps(1223), []map[string]interface{}{nil})
		t.AssertEQ(gconv.Maps([]int{}), nil)
	})
}

// 私有属性不会进行转换
func Test_Slice_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int           `json:"id"`
		name string        `json:"name"`
		Ad   []interface{} `json:"ad"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{1, "john", []interface{}{2}}
		array := gconv.Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
		t.Assert(array[0].(*User).Ad, []interface{}{2})
	})
}

func Test_Map_Basic_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := map[string]string{
			"k": "v",
		}
		m2 := map[int]string{
			3: "v",
		}
		m3 := map[float64]float32{
			1.22: 3.1,
		}
		t.Assert(gconv.Map(m1), g.Map{
			"k": "v",
		})
		t.Assert(gconv.Map(m2), g.Map{
			"3": "v",
		})
		t.Assert(gconv.Map(m3), g.Map{
			"1.22": "3.1",
		})
		t.AssertEQ(gconv.Map(nil), nil)
		t.AssertEQ(gconv.Map(map[string]interface{}{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gconv.Map(map[int]interface{}{1: 1}), map[string]interface{}{"1": 1})
		t.AssertEQ(gconv.Map(map[uint]interface{}{1: 1}), map[string]interface{}{"1": 1})
		t.AssertEQ(gconv.Map(map[uint]string{1: "1"}), map[string]interface{}{"1": "1"})

		t.AssertEQ(gconv.Map(map[interface{}]interface{}{"a": 1}), map[interface{}]interface{}{"a": 1})
		t.AssertEQ(gconv.Map(map[interface{}]string{"a": "1"}), map[interface{}]string{"a": "1"})
		t.AssertEQ(gconv.Map(map[interface{}]int{"a": 1}), map[interface{}]int{"a": 1})
		t.AssertEQ(gconv.Map(map[interface{}]uint{"a": 1}), map[interface{}]uint{"a": 1})
		t.AssertEQ(gconv.Map(map[interface{}]float32{"a": 1}), map[interface{}]float32{"a": 1})
		t.AssertEQ(gconv.Map(map[interface{}]float64{"a": 1}), map[interface{}]float64{"a": 1})

		t.AssertEQ(gconv.Map(map[string]bool{"a": true}), map[string]interface{}{"a": true})
		t.AssertEQ(gconv.Map(map[string]int{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gconv.Map(map[string]uint{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gconv.Map(map[string]float32{"a": 1}), map[string]interface{}{"a": 1})
		t.AssertEQ(gconv.Map(map[string]float64{"a": 1}), map[string]interface{}{"a": 1})

	})
}

func Test_Map_StructWithGconvTag_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string   `gconv:"-"`
			NickName string   `gconv:"nickname,omitempty"`
			Pass1    string   `gconv:"password1"`
			Pass2    string   `gconv:"password2"`
			Ss       []string `gconv:"ss"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
			Ss:      []string{"sss", "2222"},
		}
		user2 := &user1
		map1 := gconv.Map(user1)
		map2 := gconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")
		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithJsonTag_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string   `json:"-"`
			NickName string   `json:"nickname, omitempty"`
			Pass1    string   `json:"password1,newpassword"`
			Pass2    string   `json:"password2"`
			Ss       []string `json:"omitempty"`
			ssb, ssa string
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
			Ss:      []string{"sss", "2222"},
			ssb:     "11",
			ssa:     "222",
		}
		user3 := User{
			Uid:      100,
			Name:     "john",
			NickName: "SSS",
			SiteUrl:  "https://goframe.org",
			Pass1:    "123",
			Pass2:    "456",
			Ss:       []string{"sss", "2222"},
			ssb:      "11",
			ssa:      "222",
		}
		user2 := &user1
		_ = gconv.Map(user1, "Ss")
		map1 := gconv.Map(user1, "json", "json2")
		map2 := gconv.Map(user2)
		map3 := gconv.Map(user3)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")
		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
		t.Assert(map3["NickName"], nil)
	})
}

func Test_Map_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{1, "john"}
		t.Assert(gconv.Map(user), g.Map{"Id": 1})
	})
}

func Test_Map_StructInherit_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string  `json:"passport"`
			Password string  `json:"password"`
			Nickname string  `json:"nickname"`
			S        *string `json:"nickname2"`
		}

		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		var s = "s"
		user.S = &s

		m := gconv.MapDeep(user)
		t.Assert(m["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], user.CreateTime)
		t.Assert(m["nickname2"], user.S)
	})
}

func Test_Struct_Basic1_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   int
			Result string
		}

		type Score2 struct {
			Name   int
			Result string
		}

		type User struct {
			Uid      int
			Name     string
			Site_Url string
			NickName string
			Pass1    string `gconv:"password1"`
			Pass2    string `gconv:"password2"`
			As       *Score
			Ass      Score
			Assb     []interface{}
		}
		// 使用默认映射规则绑定属性值到对象
		user := new(User)
		params1 := g.Map{
			"uid":       1,
			"Name":      "john",
			"siteurl":   "https://goframe.org",
			"nick_name": "johng",
			"PASS1":     "123",
			"PASS2":     "456",
			"As":        g.Map{"Name": 1, "Result": "22222"},
			"Ass":       &Score{11, "11"},
			"Assb":      []string{"wwww"},
		}
		_ = gconv.Struct(nil, user)
		_ = gconv.Struct(params1, nil)
		_ = gconv.Struct([]interface{}{nil}, user)
		_ = gconv.Struct(user, []interface{}{nil})

		var a = []interface{}{nil}
		ab := &a
		_ = gconv.Struct(params1, *ab)
		var pi *int = nil
		_ = gconv.Struct(params1, pi)

		_ = gconv.Struct(params1, user)
		_ = gconv.Struct(params1, user, map[string]string{"uid": "Names"})
		_ = gconv.Struct(params1, user, map[string]string{"uid": "as"})

		// 使用struct tag映射绑定属性值到对象
		user = new(User)
		params2 := g.Map{
			"uid":       2,
			"name":      "smith",
			"site-url":  "https://goframe.org",
			"nick name": "johng",
			"password1": "111",
			"password2": "222",
		}
		if err := gconv.Struct(params2, user); err != nil {
			gtest.Error(err)
		}
		t.Assert(user, &User{
			Uid:      2,
			Name:     "smith",
			Site_Url: "https://goframe.org",
			NickName: "johng",
			Pass1:    "111",
			Pass2:    "222",
		})
	})
}

// 使用默认映射规则绑定属性值到对象
func Test_Struct_Basic2_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid     int
			Name    string
			SiteUrl string
			Pass1   string
			Pass2   string
		}
		user := new(User)
		params := g.Map{
			"uid":      1,
			"Name":     "john",
			"site_url": "https://goframe.org",
			"PASS1":    "123",
			"PASS2":    "456",
		}
		if err := gconv.Struct(params, user); err != nil {
			gtest.Error(err)
		}
		t.Assert(user, &User{
			Uid:     1,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		})
	})
}

// 带有指针的基础类型属性
func Test_Struct_Basic3_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid  int
			Name *string
		}
		user := new(User)
		params := g.Map{
			"uid":  1,
			"Name": "john",
		}
		if err := gconv.Struct(params, user); err != nil {
			gtest.Error(err)
		}
		t.Assert(user.Uid, 1)
		t.Assert(*user.Name, "john")
	})
}

// slice类型属性的赋值
func Test_Struct_Attr_Slice_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []int
		}
		scores := []interface{}{99, 100, 60, 140}
		user := new(User)
		if err := gconv.Struct(g.Map{"Scores": scores}, user); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(user, &User{
				Scores: []int{99, 100, 60, 140},
			})
		}
	})
}

// 属性为struct对象
func Test_Struct_Attr_Struct_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换
		if err := gconv.Struct(scores, user); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(user, &User{
				Scores: Score{
					Name:   "john",
					Result: 100,
				},
			})
		}
	})
}

// 属性为struct对象指针
func Test_Struct_Attr_Struct_Ptr_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores *Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换
		if err := gconv.Struct(scores, user); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(user.Scores, &Score{
				Name:   "john",
				Result: 100,
			})
		}
	})
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice1_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": map[string]interface{}{
				"Name":   "john",
				"Result": 100,
			},
		}

		// 嵌套struct转换，属性为slice类型，数值为map类型
		if err := gconv.Struct(scores, user); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(user.Scores, []Score{
				{
					Name:   "john",
					Result: 100,
				},
			})
		}
	})
}

// 属性为struct对象slice
func Test_Struct_Attr_Struct_Slice2_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": []interface{}{
				map[string]interface{}{
					"Name":   "john",
					"Result": 100,
				},
				map[string]interface{}{
					"Name":   "smith",
					"Result": 60,
				},
			},
		}

		// 嵌套struct转换，属性为slice类型，数值为slice map类型
		if err := gconv.Struct(scores, user); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(user.Scores, []Score{
				{
					Name:   "john",
					Result: 100,
				},
				{
					Name:   "smith",
					Result: 60,
				},
			})
		}
	})
}

// 属性为struct对象slice ptr
func Test_Struct_Attr_Struct_Slice_Ptr_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Score struct {
			Name   string
			Result int
		}
		type User struct {
			Scores []*Score
		}

		user := new(User)
		scores := map[string]interface{}{
			"Scores": []interface{}{
				map[string]interface{}{
					"Name":   "john",
					"Result": 100,
				},
				map[string]interface{}{
					"Name":   "smith",
					"Result": 60,
				},
			},
		}

		// 嵌套struct转换，属性为slice类型，数值为slice map类型
		if err := gconv.Struct(scores, user); err != nil {
			gtest.Error(err)
		} else {
			t.Assert(len(user.Scores), 2)
			t.Assert(user.Scores[0], &Score{
				Name:   "john",
				Result: 100,
			})
			t.Assert(user.Scores[1], &Score{
				Name:   "smith",
				Result: 60,
			})
		}
	})
}

func Test_Struct_PrivateAttribute_All(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := gconv.Struct(g.Map{"id": 1, "name": "john"}, user)
		t.Assert(err, nil)
		t.Assert(user.Id, 1)
		t.Assert(user.name, "")
	})
}

func Test_Struct_Embedded_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		data := g.Map{
			"id":          100,
			"uid":         101,
			"passport":    "t1",
			"password":    "123456",
			"nickname":    "T1",
			"create_time": "2019",
		}
		user := new(User)
		gconv.Struct(data, user)
		t.Assert(user.Id, 100)
		t.Assert(user.Uid, 101)
		t.Assert(user.Nickname, "T1")
		t.Assert(user.CreateTime, "2019")
	})
}

func Test_Struct_Time_All(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			CreateTime time.Time
		}
		now := time.Now()
		user := new(User)
		gconv.Struct(g.Map{
			"create_time": now,
		}, user)
		t.Assert(user.CreateTime.UTC().String(), now.UTC().String())
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			CreateTime *time.Time
		}
		now := time.Now()
		user := new(User)
		gconv.Struct(g.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.UTC().String(), now.UTC().String())
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			CreateTime *gtime.Time
		}
		now := time.Now()
		user := new(User)
		gconv.Struct(g.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			CreateTime gtime.Time
		}
		now := time.Now()
		user := new(User)
		gconv.Struct(g.Map{
			"create_time": &now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			CreateTime gtime.Time
		}
		now := time.Now()
		user := new(User)
		gconv.Struct(g.Map{
			"create_time": now,
		}, user)
		t.Assert(user.CreateTime.Time.UTC().String(), now.UTC().String())
	})
}

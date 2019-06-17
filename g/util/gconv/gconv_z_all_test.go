package gconv_test

import (
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gconv"
	"testing"
	"time"
)

type apiString interface {
	String() string
}
type S struct {
}

func (s S) String() string {
	return "22222"
}

type apiError interface {
	Error() string
}
type S1 struct {
}

func (s1 S1) Error() string {
	return "22222"
}

func Test_Bool_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.AssertEQ(gconv.Bool(i), false)
		gtest.AssertEQ(gconv.Bool(false), false)
		gtest.AssertEQ(gconv.Bool(nil), false)
		gtest.AssertEQ(gconv.Bool(0), false)
		gtest.AssertEQ(gconv.Bool("0"), false)
		gtest.AssertEQ(gconv.Bool(""), false)
		gtest.AssertEQ(gconv.Bool("false"), false)
		gtest.AssertEQ(gconv.Bool("off"), false)
		gtest.AssertEQ(gconv.Bool([]byte{}), false)
		gtest.AssertEQ(gconv.Bool([]string{}), false)
		gtest.AssertEQ(gconv.Bool([2]int{1, 2}), true)
		gtest.AssertEQ(gconv.Bool([]interface{}{}), false)
		gtest.AssertEQ(gconv.Bool([]map[int]int{}), false)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Bool(countryCapitalMap), true)

		gtest.AssertEQ(gconv.Bool("1"), true)
		gtest.AssertEQ(gconv.Bool("on"), true)
		gtest.AssertEQ(gconv.Bool(1), true)
		gtest.AssertEQ(gconv.Bool(123.456), true)
		gtest.AssertEQ(gconv.Bool(boolStruct{}), true)
		gtest.AssertEQ(gconv.Bool(&boolStruct{}), true)
	})
}

func Test_Int_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.AssertEQ(gconv.Int(i), 0)
		gtest.AssertEQ(gconv.Int(false), 0)
		gtest.AssertEQ(gconv.Int(nil), 0)
		gtest.Assert(gconv.Int(nil), 0)
		gtest.AssertEQ(gconv.Int(0), 0)
		gtest.AssertEQ(gconv.Int("0"), 0)
		gtest.AssertEQ(gconv.Int(""), 0)
		gtest.AssertEQ(gconv.Int("false"), 0)
		gtest.AssertEQ(gconv.Int("off"), 0)
		gtest.AssertEQ(gconv.Int([]byte{}), 0)
		gtest.AssertEQ(gconv.Int([]string{}), 0)
		gtest.AssertEQ(gconv.Int([2]int{1, 2}), 0)
		gtest.AssertEQ(gconv.Int([]interface{}{}), 0)
		gtest.AssertEQ(gconv.Int([]map[int]int{}), 0)

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Int(countryCapitalMap), 0)

		gtest.AssertEQ(gconv.Int("1"), 1)
		gtest.AssertEQ(gconv.Int("on"), 0)
		gtest.AssertEQ(gconv.Int(1), 1)
		gtest.AssertEQ(gconv.Int(123.456), 123)
		gtest.AssertEQ(gconv.Int(boolStruct{}), 0)
		gtest.AssertEQ(gconv.Int(&boolStruct{}), 0)
	})
}

func Test_Int8_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Int8(i), int8(0))
		gtest.AssertEQ(gconv.Int8(false), int8(0))
		gtest.AssertEQ(gconv.Int8(nil), int8(0))
		gtest.AssertEQ(gconv.Int8(0), int8(0))
		gtest.AssertEQ(gconv.Int8("0"), int8(0))
		gtest.AssertEQ(gconv.Int8(""), int8(0))
		gtest.AssertEQ(gconv.Int8("false"), int8(0))
		gtest.AssertEQ(gconv.Int8("off"), int8(0))
		gtest.AssertEQ(gconv.Int8([]byte{}), int8(0))
		gtest.AssertEQ(gconv.Int8([]string{}), int8(0))
		gtest.AssertEQ(gconv.Int8([2]int{1, 2}), int8(0))
		gtest.AssertEQ(gconv.Int8([]interface{}{}), int8(0))
		gtest.AssertEQ(gconv.Int8([]map[int]int{}), int8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Int8(countryCapitalMap), int8(0))

		gtest.AssertEQ(gconv.Int8("1"), int8(1))
		gtest.AssertEQ(gconv.Int8("on"), int8(0))
		gtest.AssertEQ(gconv.Int8(int8(1)), int8(1))
		gtest.AssertEQ(gconv.Int8(123.456), int8(123))
		gtest.AssertEQ(gconv.Int8(boolStruct{}), int8(0))
		gtest.AssertEQ(gconv.Int8(&boolStruct{}), int8(0))
	})
}

func Test_Int16_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Int16(i), int16(0))
		gtest.AssertEQ(gconv.Int16(false), int16(0))
		gtest.AssertEQ(gconv.Int16(nil), int16(0))
		gtest.AssertEQ(gconv.Int16(0), int16(0))
		gtest.AssertEQ(gconv.Int16("0"), int16(0))
		gtest.AssertEQ(gconv.Int16(""), int16(0))
		gtest.AssertEQ(gconv.Int16("false"), int16(0))
		gtest.AssertEQ(gconv.Int16("off"), int16(0))
		gtest.AssertEQ(gconv.Int16([]byte{}), int16(0))
		gtest.AssertEQ(gconv.Int16([]string{}), int16(0))
		gtest.AssertEQ(gconv.Int16([2]int{1, 2}), int16(0))
		gtest.AssertEQ(gconv.Int16([]interface{}{}), int16(0))
		gtest.AssertEQ(gconv.Int16([]map[int]int{}), int16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Int16(countryCapitalMap), int16(0))

		gtest.AssertEQ(gconv.Int16("1"), int16(1))
		gtest.AssertEQ(gconv.Int16("on"), int16(0))
		gtest.AssertEQ(gconv.Int16(int16(1)), int16(1))
		gtest.AssertEQ(gconv.Int16(123.456), int16(123))
		gtest.AssertEQ(gconv.Int16(boolStruct{}), int16(0))
		gtest.AssertEQ(gconv.Int16(&boolStruct{}), int16(0))
	})
}

func Test_Int32_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Int32(i), int32(0))
		gtest.AssertEQ(gconv.Int32(false), int32(0))
		gtest.AssertEQ(gconv.Int32(nil), int32(0))
		gtest.AssertEQ(gconv.Int32(0), int32(0))
		gtest.AssertEQ(gconv.Int32("0"), int32(0))
		gtest.AssertEQ(gconv.Int32(""), int32(0))
		gtest.AssertEQ(gconv.Int32("false"), int32(0))
		gtest.AssertEQ(gconv.Int32("off"), int32(0))
		gtest.AssertEQ(gconv.Int32([]byte{}), int32(0))
		gtest.AssertEQ(gconv.Int32([]string{}), int32(0))
		gtest.AssertEQ(gconv.Int32([2]int{1, 2}), int32(0))
		gtest.AssertEQ(gconv.Int32([]interface{}{}), int32(0))
		gtest.AssertEQ(gconv.Int32([]map[int]int{}), int32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Int32(countryCapitalMap), int32(0))

		gtest.AssertEQ(gconv.Int32("1"), int32(1))
		gtest.AssertEQ(gconv.Int32("on"), int32(0))
		gtest.AssertEQ(gconv.Int32(int32(1)), int32(1))
		gtest.AssertEQ(gconv.Int32(123.456), int32(123))
		gtest.AssertEQ(gconv.Int32(boolStruct{}), int32(0))
		gtest.AssertEQ(gconv.Int32(&boolStruct{}), int32(0))
	})
}

func Test_Int64_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.AssertEQ(gconv.Int64("0x00e"), int64(14))
		gtest.Assert(gconv.Int64("022"), int64(18))

		gtest.Assert(gconv.Int64(i), int64(0))
		gtest.Assert(gconv.Int64(true), 1)
		gtest.Assert(gconv.Int64("1"), int64(1))
		gtest.Assert(gconv.Int64("0"), int64(0))
		gtest.Assert(gconv.Int64("X"), int64(0))
		gtest.Assert(gconv.Int64("x"), int64(0))
		gtest.Assert(gconv.Int64(int64(1)), int64(1))
		gtest.Assert(gconv.Int64(int(0)), int64(0))
		gtest.Assert(gconv.Int64(int8(0)), int64(0))
		gtest.Assert(gconv.Int64(int16(0)), int64(0))
		gtest.Assert(gconv.Int64(int32(0)), int64(0))
		gtest.Assert(gconv.Int64(uint64(0)), int64(0))
		gtest.Assert(gconv.Int64(uint32(0)), int64(0))
		gtest.Assert(gconv.Int64(uint16(0)), int64(0))
		gtest.Assert(gconv.Int64(uint8(0)), int64(0))
		gtest.Assert(gconv.Int64(uint(0)), int64(0))
		gtest.Assert(gconv.Int64(float32(0)), int64(0))

		gtest.AssertEQ(gconv.Int64(false), int64(0))
		gtest.AssertEQ(gconv.Int64(nil), int64(0))
		gtest.AssertEQ(gconv.Int64(0), int64(0))
		gtest.AssertEQ(gconv.Int64("0"), int64(0))
		gtest.AssertEQ(gconv.Int64(""), int64(0))
		gtest.AssertEQ(gconv.Int64("false"), int64(0))
		gtest.AssertEQ(gconv.Int64("off"), int64(0))
		gtest.AssertEQ(gconv.Int64([]byte{}), int64(0))
		gtest.AssertEQ(gconv.Int64([]string{}), int64(0))
		gtest.AssertEQ(gconv.Int64([2]int{1, 2}), int64(0))
		gtest.AssertEQ(gconv.Int64([]interface{}{}), int64(0))
		gtest.AssertEQ(gconv.Int64([]map[int]int{}), int64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Int64(countryCapitalMap), int64(0))

		gtest.AssertEQ(gconv.Int64("1"), int64(1))
		gtest.AssertEQ(gconv.Int64("on"), int64(0))
		gtest.AssertEQ(gconv.Int64(int64(1)), int64(1))
		gtest.AssertEQ(gconv.Int64(123.456), int64(123))
		gtest.AssertEQ(gconv.Int64(boolStruct{}), int64(0))
		gtest.AssertEQ(gconv.Int64(&boolStruct{}), int64(0))
	})
}

func Test_Uint_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.AssertEQ(gconv.Uint(i), uint(0))
		gtest.AssertEQ(gconv.Uint(false), uint(0))
		gtest.AssertEQ(gconv.Uint(nil), uint(0))
		gtest.Assert(gconv.Uint(nil), uint(0))
		gtest.AssertEQ(gconv.Uint(uint(0)), uint(0))
		gtest.AssertEQ(gconv.Uint("0"), uint(0))
		gtest.AssertEQ(gconv.Uint(""), uint(0))
		gtest.AssertEQ(gconv.Uint("false"), uint(0))
		gtest.AssertEQ(gconv.Uint("off"), uint(0))
		gtest.AssertEQ(gconv.Uint([]byte{}), uint(0))
		gtest.AssertEQ(gconv.Uint([]string{}), uint(0))
		gtest.AssertEQ(gconv.Uint([2]int{1, 2}), uint(0))
		gtest.AssertEQ(gconv.Uint([]interface{}{}), uint(0))
		gtest.AssertEQ(gconv.Uint([]map[int]int{}), uint(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Uint(countryCapitalMap), uint(0))

		gtest.AssertEQ(gconv.Uint("1"), uint(1))
		gtest.AssertEQ(gconv.Uint("on"), uint(0))
		gtest.AssertEQ(gconv.Uint(1), uint(1))
		gtest.AssertEQ(gconv.Uint(123.456), uint(123))
		gtest.AssertEQ(gconv.Uint(boolStruct{}), uint(0))
		gtest.AssertEQ(gconv.Uint(&boolStruct{}), uint(0))
	})
}

func Test_Uint8_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Uint8(i), uint8(0))
		gtest.AssertEQ(gconv.Uint8(uint8(1)), uint8(1))
		gtest.AssertEQ(gconv.Uint8(false), uint8(0))
		gtest.AssertEQ(gconv.Uint8(nil), uint8(0))
		gtest.AssertEQ(gconv.Uint8(0), uint8(0))
		gtest.AssertEQ(gconv.Uint8("0"), uint8(0))
		gtest.AssertEQ(gconv.Uint8(""), uint8(0))
		gtest.AssertEQ(gconv.Uint8("false"), uint8(0))
		gtest.AssertEQ(gconv.Uint8("off"), uint8(0))
		gtest.AssertEQ(gconv.Uint8([]byte{}), uint8(0))
		gtest.AssertEQ(gconv.Uint8([]string{}), uint8(0))
		gtest.AssertEQ(gconv.Uint8([2]int{1, 2}), uint8(0))
		gtest.AssertEQ(gconv.Uint8([]interface{}{}), uint8(0))
		gtest.AssertEQ(gconv.Uint8([]map[int]int{}), uint8(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Uint8(countryCapitalMap), uint8(0))

		gtest.AssertEQ(gconv.Uint8("1"), uint8(1))
		gtest.AssertEQ(gconv.Uint8("on"), uint8(0))
		gtest.AssertEQ(gconv.Uint8(int8(1)), uint8(1))
		gtest.AssertEQ(gconv.Uint8(123.456), uint8(123))
		gtest.AssertEQ(gconv.Uint8(boolStruct{}), uint8(0))
		gtest.AssertEQ(gconv.Uint8(&boolStruct{}), uint8(0))
	})
}

func Test_Uint16_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Uint16(i), uint16(0))
		gtest.AssertEQ(gconv.Uint16(uint16(1)), uint16(1))
		gtest.AssertEQ(gconv.Uint16(false), uint16(0))
		gtest.AssertEQ(gconv.Uint16(nil), uint16(0))
		gtest.AssertEQ(gconv.Uint16(0), uint16(0))
		gtest.AssertEQ(gconv.Uint16("0"), uint16(0))
		gtest.AssertEQ(gconv.Uint16(""), uint16(0))
		gtest.AssertEQ(gconv.Uint16("false"), uint16(0))
		gtest.AssertEQ(gconv.Uint16("off"), uint16(0))
		gtest.AssertEQ(gconv.Uint16([]byte{}), uint16(0))
		gtest.AssertEQ(gconv.Uint16([]string{}), uint16(0))
		gtest.AssertEQ(gconv.Uint16([2]int{1, 2}), uint16(0))
		gtest.AssertEQ(gconv.Uint16([]interface{}{}), uint16(0))
		gtest.AssertEQ(gconv.Uint16([]map[int]int{}), uint16(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Uint16(countryCapitalMap), uint16(0))

		gtest.AssertEQ(gconv.Uint16("1"), uint16(1))
		gtest.AssertEQ(gconv.Uint16("on"), uint16(0))
		gtest.AssertEQ(gconv.Uint16(int16(1)), uint16(1))
		gtest.AssertEQ(gconv.Uint16(123.456), uint16(123))
		gtest.AssertEQ(gconv.Uint16(boolStruct{}), uint16(0))
		gtest.AssertEQ(gconv.Uint16(&boolStruct{}), uint16(0))
	})
}

func Test_Uint32_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Uint32(i), uint32(0))
		gtest.AssertEQ(gconv.Uint32(uint32(1)), uint32(1))
		gtest.AssertEQ(gconv.Uint32(false), uint32(0))
		gtest.AssertEQ(gconv.Uint32(nil), uint32(0))
		gtest.AssertEQ(gconv.Uint32(0), uint32(0))
		gtest.AssertEQ(gconv.Uint32("0"), uint32(0))
		gtest.AssertEQ(gconv.Uint32(""), uint32(0))
		gtest.AssertEQ(gconv.Uint32("false"), uint32(0))
		gtest.AssertEQ(gconv.Uint32("off"), uint32(0))
		gtest.AssertEQ(gconv.Uint32([]byte{}), uint32(0))
		gtest.AssertEQ(gconv.Uint32([]string{}), uint32(0))
		gtest.AssertEQ(gconv.Uint32([2]int{1, 2}), uint32(0))
		gtest.AssertEQ(gconv.Uint32([]interface{}{}), uint32(0))
		gtest.AssertEQ(gconv.Uint32([]map[int]int{}), uint32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Uint32(countryCapitalMap), uint32(0))

		gtest.AssertEQ(gconv.Uint32("1"), uint32(1))
		gtest.AssertEQ(gconv.Uint32("on"), uint32(0))
		gtest.AssertEQ(gconv.Uint32(int32(1)), uint32(1))
		gtest.AssertEQ(gconv.Uint32(123.456), uint32(123))
		gtest.AssertEQ(gconv.Uint32(boolStruct{}), uint32(0))
		gtest.AssertEQ(gconv.Uint32(&boolStruct{}), uint32(0))
	})
}

func Test_Uint64_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.AssertEQ(gconv.Uint64("0x00e"), uint64(14))
		gtest.Assert(gconv.Uint64("022"), uint64(18))

		gtest.AssertEQ(gconv.Uint64(i), uint64(0))
		gtest.AssertEQ(gconv.Uint64(true), uint64(1))
		gtest.Assert(gconv.Uint64("1"), int64(1))
		gtest.Assert(gconv.Uint64("0"), uint64(0))
		gtest.Assert(gconv.Uint64("X"), uint64(0))
		gtest.Assert(gconv.Uint64("x"), uint64(0))
		gtest.Assert(gconv.Uint64(int64(1)), uint64(1))
		gtest.Assert(gconv.Uint64(int(0)), uint64(0))
		gtest.Assert(gconv.Uint64(int8(0)), uint64(0))
		gtest.Assert(gconv.Uint64(int16(0)), uint64(0))
		gtest.Assert(gconv.Uint64(int32(0)), uint64(0))
		gtest.Assert(gconv.Uint64(uint64(0)), uint64(0))
		gtest.Assert(gconv.Uint64(uint32(0)), uint64(0))
		gtest.Assert(gconv.Uint64(uint16(0)), uint64(0))
		gtest.Assert(gconv.Uint64(uint8(0)), uint64(0))
		gtest.Assert(gconv.Uint64(uint(0)), uint64(0))
		gtest.Assert(gconv.Uint64(float32(0)), uint64(0))

		gtest.AssertEQ(gconv.Uint64(false), uint64(0))
		gtest.AssertEQ(gconv.Uint64(nil), uint64(0))
		gtest.AssertEQ(gconv.Uint64(0), uint64(0))
		gtest.AssertEQ(gconv.Uint64("0"), uint64(0))
		gtest.AssertEQ(gconv.Uint64(""), uint64(0))
		gtest.AssertEQ(gconv.Uint64("false"), uint64(0))
		gtest.AssertEQ(gconv.Uint64("off"), uint64(0))
		gtest.AssertEQ(gconv.Uint64([]byte{}), uint64(0))
		gtest.AssertEQ(gconv.Uint64([]string{}), uint64(0))
		gtest.AssertEQ(gconv.Uint64([2]int{1, 2}), uint64(0))
		gtest.AssertEQ(gconv.Uint64([]interface{}{}), uint64(0))
		gtest.AssertEQ(gconv.Uint64([]map[int]int{}), uint64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Uint64(countryCapitalMap), uint64(0))

		gtest.AssertEQ(gconv.Uint64("1"), uint64(1))
		gtest.AssertEQ(gconv.Uint64("on"), uint64(0))
		gtest.AssertEQ(gconv.Uint64(int64(1)), uint64(1))
		gtest.AssertEQ(gconv.Uint64(123.456), uint64(123))
		gtest.AssertEQ(gconv.Uint64(boolStruct{}), uint64(0))
		gtest.AssertEQ(gconv.Uint64(&boolStruct{}), uint64(0))
	})
}

func Test_Float32_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Float32(i), float32(0))
		gtest.AssertEQ(gconv.Float32(false), float32(0))
		gtest.AssertEQ(gconv.Float32(nil), float32(0))
		gtest.AssertEQ(gconv.Float32(0), float32(0))
		gtest.AssertEQ(gconv.Float32("0"), float32(0))
		gtest.AssertEQ(gconv.Float32(""), float32(0))
		gtest.AssertEQ(gconv.Float32("false"), float32(0))
		gtest.AssertEQ(gconv.Float32("off"), float32(0))
		gtest.AssertEQ(gconv.Float32([]byte{}), float32(0))
		gtest.AssertEQ(gconv.Float32([]string{}), float32(0))
		gtest.AssertEQ(gconv.Float32([2]int{1, 2}), float32(0))
		gtest.AssertEQ(gconv.Float32([]interface{}{}), float32(0))
		gtest.AssertEQ(gconv.Float32([]map[int]int{}), float32(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Float32(countryCapitalMap), float32(0))

		gtest.AssertEQ(gconv.Float32("1"), float32(1))
		gtest.AssertEQ(gconv.Float32("on"), float32(0))
		gtest.AssertEQ(gconv.Float32(float32(1)), float32(1))
		gtest.AssertEQ(gconv.Float32(123.456), float32(123.456))
		gtest.AssertEQ(gconv.Float32(boolStruct{}), float32(0))
		gtest.AssertEQ(gconv.Float32(&boolStruct{}), float32(0))
	})
}

func Test_Float64_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.Assert(gconv.Float64(i), float64(0))
		gtest.AssertEQ(gconv.Float64(false), float64(0))
		gtest.AssertEQ(gconv.Float64(nil), float64(0))
		gtest.AssertEQ(gconv.Float64(0), float64(0))
		gtest.AssertEQ(gconv.Float64("0"), float64(0))
		gtest.AssertEQ(gconv.Float64(""), float64(0))
		gtest.AssertEQ(gconv.Float64("false"), float64(0))
		gtest.AssertEQ(gconv.Float64("off"), float64(0))
		gtest.AssertEQ(gconv.Float64([]byte{}), float64(0))
		gtest.AssertEQ(gconv.Float64([]string{}), float64(0))
		gtest.AssertEQ(gconv.Float64([2]int{1, 2}), float64(0))
		gtest.AssertEQ(gconv.Float64([]interface{}{}), float64(0))
		gtest.AssertEQ(gconv.Float64([]map[int]int{}), float64(0))

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.Float64(countryCapitalMap), float64(0))

		gtest.AssertEQ(gconv.Float64("1"), float64(1))
		gtest.AssertEQ(gconv.Float64("on"), float64(0))
		gtest.AssertEQ(gconv.Float64(float64(1)), float64(1))
		gtest.AssertEQ(gconv.Float64(123.456), float64(123.456))
		gtest.AssertEQ(gconv.Float64(boolStruct{}), float64(0))
		gtest.AssertEQ(gconv.Float64(&boolStruct{}), float64(0))
	})
}

func Test_String_All(t *testing.T) {
	gtest.Case(t, func() {
		var s []rune
		gtest.AssertEQ(gconv.String(s), "")
		var i interface{} = nil
		gtest.AssertEQ(gconv.String(i), "")
		gtest.AssertEQ(gconv.String("1"), "1")
		gtest.AssertEQ(gconv.String("0"), string("0"))
		gtest.Assert(gconv.String("X"), string("X"))
		gtest.Assert(gconv.String("x"), string("x"))
		gtest.Assert(gconv.String(int64(1)), uint64(1))
		gtest.Assert(gconv.String(int(0)), string("0"))
		gtest.Assert(gconv.String(int8(0)), string("0"))
		gtest.Assert(gconv.String(int16(0)), string("0"))
		gtest.Assert(gconv.String(int32(0)), string("0"))
		gtest.Assert(gconv.String(uint64(0)), string("0"))
		gtest.Assert(gconv.String(uint32(0)), string("0"))
		gtest.Assert(gconv.String(uint16(0)), string("0"))
		gtest.Assert(gconv.String(uint8(0)), string("0"))
		gtest.Assert(gconv.String(uint(0)), string("0"))
		gtest.Assert(gconv.String(float32(0)), string("0"))
		gtest.AssertEQ(gconv.String(true), "true")
		gtest.AssertEQ(gconv.String(false), "false")
		gtest.AssertEQ(gconv.String(nil), "")
		gtest.AssertEQ(gconv.String(0), string("0"))
		gtest.AssertEQ(gconv.String("0"), string("0"))
		gtest.AssertEQ(gconv.String(""), "")
		gtest.AssertEQ(gconv.String("false"), "false")
		gtest.AssertEQ(gconv.String("off"), string("off"))
		gtest.AssertEQ(gconv.String([]byte{}), "")
		gtest.AssertEQ(gconv.String([]string{}), "[]")
		gtest.AssertEQ(gconv.String([2]int{1, 2}), "[1,2]")
		gtest.AssertEQ(gconv.String([]interface{}{}), "[]")
		gtest.AssertEQ(gconv.String(map[int]int{}), "{}")

		var countryCapitalMap = make(map[string]string)
		/* map插入key - value对,各个国家对应的首都 */
		countryCapitalMap [ "France" ] = "巴黎"
		countryCapitalMap [ "Italy" ] = "罗马"
		countryCapitalMap [ "Japan" ] = "东京"
		countryCapitalMap [ "India " ] = "新德里"
		gtest.AssertEQ(gconv.String(countryCapitalMap), `{"France":"巴黎","India ":"新德里","Italy":"罗马","Japan":"东京"}`)
		gtest.AssertEQ(gconv.String(int64(1)), "1")
		gtest.AssertEQ(gconv.String(123.456), "123.456")
		gtest.AssertEQ(gconv.String(boolStruct{}), "{}")
		gtest.AssertEQ(gconv.String(&boolStruct{}), "{}")

		var info apiString
		info = new(S)
		gtest.AssertEQ(gconv.String(info), "22222")
		var errinfo apiError
		errinfo = new(S1)
		gtest.AssertEQ(gconv.String(errinfo), "22222")
	})
}

func Test_Runes_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertEQ(gconv.Runes("www"), []int32{119, 119, 119})
		var s []rune
		gtest.AssertEQ(gconv.Runes(s), nil)
	})
}

func Test_Rune_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertEQ(gconv.Rune("www"), int32(0))
		gtest.AssertEQ(gconv.Rune(int32(0)), int32(0))
		var s []rune
		gtest.AssertEQ(gconv.Rune(s), int32(0))
	})
}

func Test_Bytes_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertEQ(gconv.Bytes(nil), nil)
		gtest.AssertEQ(gconv.Bytes(int32(0)), []uint8{0, 0, 0, 0})
		gtest.AssertEQ(gconv.Bytes("s"), []uint8{115})
		gtest.AssertEQ(gconv.Bytes([]byte("s")), []uint8{115})
	})
}

func Test_Byte_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertEQ(gconv.Byte(uint8(0)), uint8(0))
		gtest.AssertEQ(gconv.Byte("s"), uint8(0))
		gtest.AssertEQ(gconv.Byte([]byte("s")), uint8(0))
	})
}

func Test_Convert_All(t *testing.T) {
	gtest.Case(t, func() {
		var i interface{} = nil
		gtest.AssertEQ(gconv.Convert(i, "string"), "")
		gtest.AssertEQ(gconv.Convert("1", "string"), "1")
		gtest.Assert(gconv.Convert(int64(1), "int64"), int64(1))
		gtest.Assert(gconv.Convert(int(0), "int"), int(0))
		gtest.Assert(gconv.Convert(int8(0), "int8"), int8(0))
		gtest.Assert(gconv.Convert(int16(0), "int16"), int16(0))
		gtest.Assert(gconv.Convert(int32(0), "int32"), int32(0))
		gtest.Assert(gconv.Convert(uint64(0), "uint64"), uint64(0))
		gtest.Assert(gconv.Convert(uint32(0), "uint32"), uint32(0))
		gtest.Assert(gconv.Convert(uint16(0), "uint16"), uint16(0))
		gtest.Assert(gconv.Convert(uint8(0), "uint8"), uint8(0))
		gtest.Assert(gconv.Convert(uint(0), "uint"), uint(0))
		gtest.Assert(gconv.Convert(float32(0), "float32"), float32(0))
		gtest.Assert(gconv.Convert(float64(0), "float64"), float64(0))
		gtest.AssertEQ(gconv.Convert(true, "bool"), true)
		gtest.AssertEQ(gconv.Convert([]byte{}, "[]byte"), []uint8{})
		gtest.AssertEQ(gconv.Convert([]string{}, "[]string"), []string{})
		gtest.AssertEQ(gconv.Convert([2]int{1, 2}, "[]int"), []int{0})
		gtest.AssertEQ(gconv.Convert(0, "Time", 0), gconv.Time("0000-01-01 00:00:00 +0800 CST"))
		gtest.AssertEQ(gconv.Convert(1989, "Time"), gconv.Time("2033-01-11 04:00:00 +0800 CST"))
		gtest.AssertEQ(gconv.Convert(gtime.Now(), "gtime.Time", 1), nil)
		gtest.AssertEQ(gconv.Convert(1989, "gtime.Time"), gtime.Time{gconv.Time("2033-01-11 04:00:00 +0800 CST")})
		gtest.AssertEQ(gconv.Convert(gtime.Now(), "*gtime.Time", 1), nil)
		gtest.AssertEQ(gconv.Convert(gtime.Now(), "GTime", 1), nil)
		gtest.AssertEQ(gconv.Convert(1989, "*gtime.Time"), gconv.GTime(1989))
		gtest.AssertEQ(gconv.Convert(1989, "Duration"), time.Duration(int64(1989)))
		gtest.AssertEQ(gconv.Convert("1989", "Duration"), time.Duration(int64(1989)))
		gtest.AssertEQ(gconv.Convert("1989", ""), "1989")
	})
}

func Test_Time_All(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertEQ(gconv.Duration(""), time.Duration(int64(0)))
		gtest.AssertEQ(gconv.GTime(""), gtime.New())
	})
}


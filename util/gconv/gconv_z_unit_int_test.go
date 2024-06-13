// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var intTests = []struct {
	value    interface{}
	expect   int
	expect8  int8
	expect16 int16
	expect32 int32
	expect64 int64
}{
	{true, 1, 1, 1, 1, 1},
	{false, 0, 0, 0, 0, 0},

	{int(0), 0, 0, 0, 0, 0},
	{int(123), 123, 123, 123, 123, 123},
	{int8(123), 123, 123, 123, 123, 123},
	{int16(123), 123, 123, 123, 123, 123},
	{int32(123), 123, 123, 123, 123, 123},
	{int64(123), 123, 123, 123, 123, 123},

	{uint(0), 0, 0, 0, 0, 0},
	{uint(123), 123, 123, 123, 123, 123},
	{uint8(123), 123, 123, 123, 123, 123},
	{uint16(123), 123, 123, 123, 123, 123},
	{uint32(123), 123, 123, 123, 123, 123},
	{uint64(123), 123, 123, 123, 123, 123},

	{uintptr(0), 0, 0, 0, 0, 0},
	{uintptr(123), 123, 123, 123, 123, 123},

	{rune(0), 0, 0, 0, 0, 0},
	{rune(49), 49, 49, 49, 49, 49},

	{float32(123), 123, 123, 123, 123, 123},
	{float64(123.456), 123, 123, 123, 123, 123},

	{[]byte(""), 0, 0, 0, 0, 0},

	{"", 0, 0, 0, 0, 0},
	{"0", 0, 0, 0, 0, 0},
	{"1", 1, 1, 1, 1, 1},
	{"+1", 1, 1, 1, 1, 1},
	{"-1", -1, -1, -1, -1, -1},
	{"0xA", 10, 10, 10, 10, 10},
	{"-0xA", -10, -10, -10, -10, -10},
	{"0XA", 10, 10, 10, 10, 10},
	{"-0XA", -10, -10, -10, -10, -10},
	{"123.456", 123, 123, 123, 123, 123},
	{"true", 0, 0, 0, 0, 0},
	{"false", 0, 0, 0, 0, 0},
	{"on", 0, 0, 0, 0, 0},
	{"off", 0, 0, 0, 0, 0},
	{"NaN", 0, 0, 0, 0, 0},

	{complex(1, 2), 0, 0, 0, 0, 0},
	{complex(123.456, 789.123), 0, 0, 0, 0, 0},

	{[3]int{1, 2, 3}, 0, 0, 0, 0, 0},
	{[]int{1, 2, 3}, 0, 0, 0, 0, 0},

	{map[int]int{1: 1}, 0, 0, 0, 0, 0},
	{map[string]string{"Earth": "大西洋"}, 0, 0, 0, 0, 0},

	{struct{}{}, 0, 0, 0, 0, 0},
	{nil, 0, 0, 0, 0, 0},

	{gvar.New(123), 123, 123, 123, 123, 123},
	{gvar.New(123.456), 123, 123, 123, 123, 123},
}

func TestInt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			t.AssertEQ(gconv.Int(test.value), test.expect)
		}
	})
}

func TestInt8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			t.AssertEQ(gconv.Int8(test.value), test.expect8)
		}
	})
}

func TestInt16(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			t.AssertEQ(gconv.Int16(test.value), test.expect16)
		}
	})
}

func TestInt32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			t.AssertEQ(gconv.Int32(test.value), test.expect32)
		}
	})
}

func TestInt64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			t.AssertEQ(gconv.Int64(test.value), test.expect64)
		}
	})
}

func TestInts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			if test.value == nil {
				t.AssertNil(gconv.Ints(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				ints      = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []int{
					test.expect, test.expect,
				}
			)
			ints = reflect.Append(ints, reflect.ValueOf(test.value))
			ints = reflect.Append(ints, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Ints(ints.Interface()), expects)
			t.AssertEQ(gconv.SliceInt(ints.Interface()), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		// string
		t.AssertEQ(gconv.Ints(""), []int{})
		t.AssertEQ(gconv.Ints("123"), []int{123})

		// []int8 json
		t.AssertEQ(gconv.Ints([]uint8(`{"Name":"Earth"}`)),
			[]int{123, 34, 78, 97, 109, 101, 34, 58, 34, 69, 97, 114, 116, 104, 34, 125})

		// []interface
		t.AssertEQ(gconv.Ints([]interface{}{1, 2, 3}), []int{1, 2, 3})

		// gvar.Var
		t.AssertEQ(gconv.Ints(
			gvar.New([]int{1, 2, 3}),
		), []int{1, 2, 3})

		// array
		t.AssertEQ(gconv.Ints("[1, 2]"), []int{1, 2})
	})
}

func TestInt32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			if test.value == nil {
				t.AssertNil(gconv.Int32s(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				int32s    = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []int32{
					test.expect32, test.expect32,
				}
			)
			int32s = reflect.Append(int32s, reflect.ValueOf(test.value))
			int32s = reflect.Append(int32s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Int32s(int32s.Interface()), expects)
			t.AssertEQ(gconv.SliceInt32(int32s.Interface()), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		// string
		t.AssertEQ(gconv.Int32s(""), []int32{})
		t.AssertEQ(gconv.Int32s("123"), []int32{123})

		// []int8 json
		t.AssertEQ(gconv.Int32s([]uint8(`{"Name":"Earth"}"`)),
			[]int32{123, 34, 78, 97, 109, 101, 34, 58, 34, 69, 97, 114, 116, 104, 34, 125, 34})

		// []interface
		t.AssertEQ(gconv.Int32s([]interface{}{1, 2, 3}), []int32{1, 2, 3})

		// gvar.Var
		t.AssertEQ(gconv.Int32s(
			gvar.New([]int32{1, 2, 3}),
		), []int32{1, 2, 3})
	})
}

func TestInt64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range intTests {
			if test.value == nil {
				t.AssertNil(gconv.Int64s(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				int64s    = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []int64{
					test.expect64, test.expect64,
				}
			)
			int64s = reflect.Append(int64s, reflect.ValueOf(test.value))
			int64s = reflect.Append(int64s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Int64s(int64s.Interface()), expects)
			t.AssertEQ(gconv.SliceInt64(int64s.Interface()), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		// string
		t.AssertEQ(gconv.Int64s(""), []int64{})
		t.AssertEQ(gconv.Int64s("123"), []int64{123})

		// []int8 json
		t.AssertEQ(gconv.Int64s([]uint8(`{"Name":"Earth"}"`)),
			[]int64{123, 34, 78, 97, 109, 101, 34, 58, 34, 69, 97, 114, 116, 104, 34, 125, 34})

		// []interface
		t.AssertEQ(gconv.Int64s([]interface{}{1, 2, 3}), []int64{1, 2, 3})

		// gvar.Var
		t.AssertEQ(gconv.Int64s(
			gvar.New([]int64{1, 2, 3}),
		), []int64{1, 2, 3})
	})
}

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

var uintTests = []struct {
	value    interface{}
	expect   uint
	expect8  uint8
	expect16 uint16
	expect32 uint32
	expect64 uint64
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
	{"0xA", 10, 10, 10, 10, 10},
	{"0XA", 10, 10, 10, 10, 10},
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
	{map[string]string{"Earth": "珠穆朗玛峰"}, 0, 0, 0, 0, 0},

	{struct{}{}, 0, 0, 0, 0, 0},
	//{make(chan interface{}), 0, 0, 0, 0, 0},
	//{func() {}, 0, 0, 0, 0, 0},
	{nil, 0, 0, 0, 0, 0},

	{gvar.New(123), 123, 123, 123, 123, 123},
	{gvar.New(123.456), 123, 123, 123, 123, 123},
}

func TestUint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			t.AssertEQ(gconv.Uint(test.value), test.expect)
		}
	})
}

func TestUint8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			t.AssertEQ(gconv.Uint8(test.value), test.expect8)
		}
	})
}

func TestUint16(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			t.AssertEQ(gconv.Uint16(test.value), test.expect16)
		}
	})
}

func TestUint32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			t.AssertEQ(gconv.Uint32(test.value), test.expect32)
		}
	})
}

func TestUint64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			t.AssertEQ(gconv.Uint64(test.value), test.expect64)
		}
	})
}

func TestUints(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			if test.value == nil {
				t.AssertEQ(gconv.Uints(test.value), nil)
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				uints     = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []uint{
					test.expect, test.expect,
				}
			)
			uints = reflect.Append(uints, reflect.ValueOf(test.value))
			uints = reflect.Append(uints, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Uints(uints.Interface()), expects)
			t.AssertEQ(gconv.SliceUint(uints.Interface()), expects)
		}
	})
}

func TestUint32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			if test.value == nil {
				t.AssertEQ(gconv.Uint32s(test.value), nil)
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				uint32s   = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []uint32{
					test.expect32, test.expect32,
				}
			)
			uint32s = reflect.Append(uint32s, reflect.ValueOf(test.value))
			uint32s = reflect.Append(uint32s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Uint32s(uint32s.Interface()), expects)
			t.AssertEQ(gconv.SliceUint32(uint32s.Interface()), expects)
		}
	})
}

func TestUint64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range uintTests {
			if test.value == nil {
				t.AssertEQ(gconv.Uint64s(test.value), nil)
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				uint64s   = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []uint64{
					test.expect64, test.expect64,
				}
			)
			uint64s = reflect.Append(uint64s, reflect.ValueOf(test.value))
			uint64s = reflect.Append(uint64s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Uint64s(uint64s.Interface()), expects)
			t.AssertEQ(gconv.SliceUint64(uint64s.Interface()), expects)
		}
	})
}

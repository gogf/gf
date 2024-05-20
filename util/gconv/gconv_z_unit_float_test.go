// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var floatTests = []struct {
	value    interface{}
	expect32 float32
	expect64 float64
}{
	{true, 0, 0},
	{false, 0, 0},

	{int(0), 0, 0},
	{int(123), 123, 123},
	{int8(123), 123, 123},
	{int16(123), 123, 123},
	{int32(123), 123, 123},
	{int64(123), 123, 123},

	{uint(0), 0, 0},
	{uint(123), 123, 123},
	{uint8(123), 123, 123},
	{uint16(123), 123, 123},
	{uint32(123), 123, 123},
	{uint64(123), 123, 123},

	{uintptr(0), 0, 0},
	{uintptr(123), 123, 123},

	{rune(0), 0, 0},
	{rune(49), 49, 49},

	{float32(123), 123, 123},
	{float64(123.456), 123.456, 123.456},

	{[]byte(""), 0, 0},

	{"0", 0, 0},
	{"", 0, 0},
	{"1", 1, 1},
	{"123.456", 123.456, 123.456},
	{"true", 0, 0},
	{"false", 0, 0},
	{"on", 0, 0},
	{"off", 0, 0},
	{"NaN", float32(math.NaN()), math.NaN()},

	{complex(1, 2), 0, 0},
	{complex(123.456, 789.123), 0, 0},

	{[3]int{1, 2, 3}, 0, 0},
	{[]int{1, 2, 3}, 0, 0},

	{map[int]int{1: 1}, 0, 0},
	{map[string]string{"Earth": "太平洋"}, 0, 0},

	{struct{}{}, 0, 0},
	{nil, 0, 0},

	{gvar.New(123), 123, 123},
	{gvar.New(123.456), 123.456, 123.456},
}

func TestFloat32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range floatTests {
			t.AssertEQ(gconv.Float32(test.value), test.expect32)
		}
	})
}

func TestFloat64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range floatTests {
			t.AssertEQ(gconv.Float64(test.value), test.expect64)
		}
	})
}

func TestFloat32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range floatTests {
			if test.value == nil {
				t.AssertNil(gconv.Float32s(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				float32s  = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []float32{
					test.expect32, test.expect32,
				}
			)
			float32s = reflect.Append(float32s, reflect.ValueOf(test.value))
			float32s = reflect.Append(float32s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Float32s(float32s.Interface()), expects)
			t.AssertEQ(gconv.SliceFloat32(float32s.Interface()), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		// string
		t.AssertEQ(gconv.Float32s(""), []float32{})
		t.AssertEQ(gconv.Float32s("123"), []float32{123})

		// []int8 json
		t.AssertEQ(gconv.Float32s([]uint8(`{"Name":"Earth"}"`)),
			[]float32{123, 34, 78, 97, 109, 101, 34, 58, 34, 69, 97, 114, 116, 104, 34, 125, 34})

		// []interface
		t.AssertEQ(gconv.Float32s([]interface{}{1, 2, 3}), []float32{1, 2, 3})

		// gvar.Var
		t.AssertEQ(gconv.Float32s(
			gvar.New([]float32{1, 2, 3}),
		), []float32{1, 2, 3})
	})
}

func TestFloat64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range floatTests {
			if test.value == nil {
				t.AssertNil(gconv.Float64s(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				float64s  = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []float64{
					test.expect64, test.expect64,
				}
			)
			float64s = reflect.Append(float64s, reflect.ValueOf(test.value))
			float64s = reflect.Append(float64s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Float64s(float64s.Interface()), expects)
			t.AssertEQ(gconv.SliceFloat64(float64s.Interface()), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		// string
		t.AssertEQ(gconv.Float64s(""), []float64{})
		t.AssertEQ(gconv.Float64s("123"), []float64{123})

		// []int8 json
		t.AssertEQ(gconv.Float64s([]uint8(`{"Name":"Earth"}"`)),
			[]float64{123, 34, 78, 97, 109, 101, 34, 58, 34, 69, 97, 114, 116, 104, 34, 125, 34})

		// []interface
		t.AssertEQ(gconv.Float64s([]interface{}{1, 2, 3}), []float64{1, 2, 3})

		// gvar.Var
		t.AssertEQ(gconv.Float64s(
			gvar.New([]float64{1, 2, 3}),
		), []float64{1, 2, 3})
	})
}

// gconv.Floats uses gconv.Float64s.
func TestFloats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range floatTests {
			if test.value == nil {
				t.AssertNil(gconv.Floats(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				float64s  = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []float64{
					test.expect64, test.expect64,
				}
			)
			float64s = reflect.Append(float64s, reflect.ValueOf(test.value))
			float64s = reflect.Append(float64s, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Floats(float64s.Interface()), expects)
			t.AssertEQ(gconv.SliceFloat(float64s.Interface()), expects)
		}
	})
}

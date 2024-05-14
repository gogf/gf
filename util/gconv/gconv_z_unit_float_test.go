// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"math"
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
	{make(chan interface{}), 0, 0},
	{func() {}, 0, 0},
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

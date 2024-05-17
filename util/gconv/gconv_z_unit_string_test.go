// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var stringTests = []struct {
	value  interface{}
	expect string
}{
	{true, "true"},
	{false, "false"},

	{int(0), "0"},
	{int(123), "123"},
	{int8(123), "123"},
	{int16(123), "123"},
	{int32(123), "123"},
	{int64(123), "123"},

	{uint(0), "0"},
	{uint(123), "123"},
	{uint8(123), "123"},
	{uint16(123), "123"},
	{uint32(123), "123"},
	{uint64(123), "123"},

	{uintptr(0), "0"},
	{uintptr(123), "123"},

	{rune(0), "0"},
	{rune(49), "49"},

	{float32(123), "123"},
	{float64(123.456), "123.456"},

	{[]byte(""), ""},

	{"", ""},
	{"true", "true"},
	{"false", "false"},
	{"Neptune", "Neptune"},

	{complex(1, 2), "(1+2i)"},
	{complex(123.456, 789.123), "(123.456+789.123i)"},

	{[3]int{1, 2, 3}, "[1,2,3]"},
	{[]int{1, 2, 3}, "[1,2,3]"},

	{map[int]int{1: 1}, `{"1":1}`},
	{map[string]string{"Earth": "太平洋"}, `{"Earth":"太平洋"}`},

	{struct{}{}, "{}"},
	{nil, ""},

	{gvar.New(123), "123"},
	{gvar.New(123.456), "123.456"},

	{goTime, "1911-10-10 00:00:00 +0000 UTC"},
	{&goTime, "1911-10-10 00:00:00 +0000 UTC"},
	// TODO The String method of gtime not equals to time.Time
	{gfTime, "1911-10-10 00:00:00"},
	{&gfTime, "1911-10-10 00:00:00"},
	//{gfTime, "1911-10-10 00:00:00 +0000 UTC"},
	//{&gfTime, "1911-10-10 00:00:00 +0000 UTC"},
}

var (
	goTime = time.Date(
		1911, 10, 10, 0, 0, 0, 0, time.UTC,
	)
	gfTime = gtime.NewFromTime(goTime)
)

func TestString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range stringTests {
			t.AssertEQ(gconv.String(test.value), test.expect)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Strings(nil), nil)
	})
}

func TestStrings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range stringTests {
			if test.value == nil {
				t.AssertNil(gconv.Strings(test.value))
				continue
			}

			var (
				sliceType = reflect.SliceOf(reflect.TypeOf(test.value))
				strings   = reflect.MakeSlice(sliceType, 0, 0)
				expects   = []string{
					test.expect, test.expect,
				}
			)
			strings = reflect.Append(strings, reflect.ValueOf(test.value))
			strings = reflect.Append(strings, reflect.ValueOf(test.value))

			t.AssertEQ(gconv.Strings(strings.Interface()), expects)
			t.AssertEQ(gconv.SliceStr(strings.Interface()), expects)
		}
	})

	// Test for special types.
	gtest.C(t, func(t *gtest.T) {
		// []int8 json
		t.AssertEQ(gconv.Strings([]uint8(`{"Name":"Earth"}"`)),
			[]string{"123", "34", "78", "97", "109", "101", "34", "58", "34", "69", "97", "114", "116", "104", "34", "125", "34"})
	})
}

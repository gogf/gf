// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var runeTests = []struct {
	value   interface{}
	expect  rune
	expects []rune
}{
	{true, 1, []rune("true")},
	{false, 0, []rune("false")},

	{int(0), 0, []rune("0")},
	{int(123), 123, []rune("123")},
	{int8(123), 123, []rune("123")},
	{int16(123), 123, []rune("123")},
	{int32(123123123), 123123123, []rune("123123123")},
	{int64(123123123123123123), 23327667, []rune("123123123123123123")},

	{uint(0), 0, []rune("0")},
	{uint(123), 123, []rune("123")},
	{uint8(123), 123, []rune("123")},
	{uint16(123), 123, []rune("123")},
	{uint32(123123123), 123123123, []rune("123123123")},
	{uint64(123123123123123123), 23327667, []rune("123123123123123123")},

	{uintptr(0), 0, []rune{48}},
	{uintptr(123), 123, []rune{49, 50, 51}},

	{rune(0), 0, []rune("0")},
	{rune(49), 49, []rune("49")},

	{float32(123), 123, []rune{49, 50, 51}},
	{float64(123.456), 123, []rune{49, 50, 51, 46, 52, 53, 54}},

	{[]rune(""), 0, []rune("")},

	{"Uranus", 0, []rune("Uranus")},

	{complex(1, 2), 0,
		[]rune{40, 49, 43, 50, 105, 41}},

	{[3]int{1, 2, 3}, 0, []rune{91, 49, 44, 50, 44, 51, 93}},
	{[]int{1, 2, 3}, 0, []rune{91, 49, 44, 50, 44, 51, 93}},

	{map[int]int{1: 1}, 0, []rune(`{"1":1}`)},
	{map[string]string{"Earth": "印度洋"}, 0, []rune(`{"Earth":"印度洋"}`)},

	{gvar.New(123), 123, []rune{49, 50, 51}},
	{gvar.New(123.456), 123, []rune{49, 50, 51, 46, 52, 53, 54}},
}

func TestRune(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range runeTests {
			t.AssertEQ(gconv.Rune(test.value), test.expect)
		}
	})
}

func TestRunes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range runeTests {
			t.AssertEQ(gconv.Runes(test.value), test.expects)
		}
	})
}

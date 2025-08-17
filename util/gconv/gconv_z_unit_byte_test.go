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

var byteTests = []struct {
	value   interface{}
	expect  byte
	expects []byte
}{
	{true, 1, []byte{1}},
	{false, 0, []byte{0}},

	{int(0), 0, []byte{0}},
	{int(123), 123, []byte{123}},
	{int8(123), 123, []byte{123}},
	{int16(123), 123, []byte{123, 0}},
	{int32(123123123), 179, []byte{179, 181, 86, 7}},
	{int64(123123123123123123), 179, []byte{179, 243, 99, 1, 212, 107, 181, 1}},

	{uint(0), 0, []byte{0}},
	{uint(123), 123, []byte{123}},
	{uint8(123), 123, []byte{123}},
	{uint16(123), 123, []byte{123, 0}},
	{uint32(123123123), 179, []byte{179, 181, 86, 7}},
	{uint64(123123123123123123), 179, []byte{179, 243, 99, 1, 212, 107, 181, 1}},

	{uintptr(0), 0, []byte{48}},
	{uintptr(123), 123, []byte{49, 50, 51}},

	{rune(0), 0, []byte{0, 0, 0, 0}},
	{rune(49), 49, []byte{49, 0, 0, 0}},

	{float32(123), 123, []byte{0, 0, 246, 66}},
	{float64(123.456), 123, []byte{119, 190, 159, 26, 47, 221, 94, 64}},

	{[]byte(""), 0, []byte("")},

	{"Uranus", 0, []byte("Uranus")},

	{complex(1, 2), 0,
		[]byte{0, 0, 0, 0, 0, 0, 240, 63, 0, 0, 0, 0, 0, 0, 0, 64}},

	{[3]int{1, 2, 3}, 0, []byte{1, 2, 3}},
	{[]int{1, 2, 3}, 0, []byte{1, 2, 3}},

	{map[int]int{1: 1}, 0, []byte(`{"1":1}`)},
	{map[string]string{"Earth": "印度洋"}, 0, []byte(`{"Earth":"印度洋"}`)},

	{gvar.New(123), 123, []byte{123}},
	{gvar.New(123.456), 123, []byte{119, 190, 159, 26, 47, 221, 94, 64}},
}

func TestByte(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range byteTests {
			t.AssertEQ(gconv.Byte(test.value), test.expect)
		}
	})
}

func TestBytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range byteTests {
			t.AssertEQ(gconv.Bytes(test.value), test.expects)
		}
	})

	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Bytes(nil), nil)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

var interfacesTests = []struct {
	value  interface{}
	expect []interface{}
}{
	{[]bool{true, false}, []interface{}{true, false}},

	{[]int{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]int8{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]int16{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]int32{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]int64{0, 1, 2}, []interface{}{0, 1, 2}},

	{[]uint{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]uint8{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]uint16{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]uint32{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]uint64{0, 1, 2}, []interface{}{0, 1, 2}},

	{[]uintptr{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]rune{0, 1, 2}, []interface{}{0, 1, 2}},

	{[]float32{0, 1, 2}, []interface{}{0, 1, 2}},
	{[]float64{0, 1, 2}, []interface{}{0, 1, 2}},

	{[][]byte{[]byte("0"), []byte("1"), []byte("2")},
		[]interface{}{[]byte("0"), []byte("1"), []byte("2")}},
	{[]string{"0", "1", "2"}, []interface{}{"0", "1", "2"}},

	{[]complex64{0, 1, 2}, []interface{}{0 + 0i, 1 + 0i, 2 + 0i}},
	{[]complex128{0, 1, 2}, []interface{}{0 + 0i, 1 + 0i, 2 + 0i}},

	{[]interface{}{0, 1, 2}, []interface{}{0, 1, 2}},
	{nil, nil},
}

func TestInterfaces(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, test := range interfacesTests {
			if test.value == nil {
				t.AssertNil(gconv.Interfaces(test.value))
				continue
			}
			t.AssertEQ(gconv.Interfaces(test.value), test.expect)
		}
	})
}

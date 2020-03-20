// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

type boolStruct struct {
}

func Test_Bool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var i interface{} = nil
		t.AssertEQ(gconv.Bool(i), false)
		t.AssertEQ(gconv.Bool(false), false)
		t.AssertEQ(gconv.Bool(nil), false)
		t.AssertEQ(gconv.Bool(0), false)
		t.AssertEQ(gconv.Bool("0"), false)
		t.AssertEQ(gconv.Bool(""), false)
		t.AssertEQ(gconv.Bool("false"), false)
		t.AssertEQ(gconv.Bool("off"), false)
		t.AssertEQ(gconv.Bool([]byte{}), false)
		t.AssertEQ(gconv.Bool([]string{}), false)
		t.AssertEQ(gconv.Bool([]interface{}{}), false)
		t.AssertEQ(gconv.Bool([]map[int]int{}), false)

		t.AssertEQ(gconv.Bool("1"), true)
		t.AssertEQ(gconv.Bool("on"), true)
		t.AssertEQ(gconv.Bool(1), true)
		t.AssertEQ(gconv.Bool(123.456), true)
		t.AssertEQ(gconv.Bool(boolStruct{}), true)
		t.AssertEQ(gconv.Bool(&boolStruct{}), true)
	})
}

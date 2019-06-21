// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gconv"
	"testing"
)

type boolStruct struct {
}

func Test_Bool(t *testing.T) {
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
		gtest.AssertEQ(gconv.Bool([]interface{}{}), false)
		gtest.AssertEQ(gconv.Bool([]map[int]int{}), false)

		gtest.AssertEQ(gconv.Bool("1"), true)
		gtest.AssertEQ(gconv.Bool("on"), true)
		gtest.AssertEQ(gconv.Bool(1), true)
		gtest.AssertEQ(gconv.Bool(123.456), true)
		gtest.AssertEQ(gconv.Bool(boolStruct{}), true)
		gtest.AssertEQ(gconv.Bool(&boolStruct{}), true)
	})
}

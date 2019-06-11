// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)


func Test_Slice(t *testing.T) {
    gtest.Case(t, func() {
        value := 123.456
        gtest.AssertEQ(gconv.Bytes("123"),      []byte("123"))
        gtest.AssertEQ(gconv.Strings(value),    []string{"123.456"})
        gtest.AssertEQ(gconv.Ints(value),       []int{123})
        gtest.AssertEQ(gconv.Floats(value),     []float64{123.456})
        gtest.AssertEQ(gconv.Interfaces(value), []interface{}{123.456})
    })
}

// 私有属性不会进行转换
func Test_Slice_PrivateAttribute(t *testing.T) {
    type User struct {
        Id   int
        name string
    }
    gtest.Case(t, func() {
        user := &User{1, "john"}
        gtest.Assert(gconv.Interfaces(user), g.Slice{1})
    })
}

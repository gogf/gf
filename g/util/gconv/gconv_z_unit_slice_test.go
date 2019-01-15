// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gconv_test

import (
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gtest"
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

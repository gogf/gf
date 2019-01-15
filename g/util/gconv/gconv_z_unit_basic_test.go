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


func Test_Basic(t *testing.T) {
    gtest.Case(t, func() {
        value := 123.456
        gtest.AssertEQ(gconv.Int(value),     int(123))
        gtest.AssertEQ(gconv.Int8(value),    int8(123))
        gtest.AssertEQ(gconv.Int16(value),   int16(123))
        gtest.AssertEQ(gconv.Int32(value),   int32(123))
        gtest.AssertEQ(gconv.Int64(value),   int64(123))
        gtest.AssertEQ(gconv.Uint(value),    uint(123))
        gtest.AssertEQ(gconv.Uint8(value),   uint8(123))
        gtest.AssertEQ(gconv.Uint16(value),  uint16(123))
        gtest.AssertEQ(gconv.Uint32(value),  uint32(123))
        gtest.AssertEQ(gconv.Uint64(value),  uint64(123))
        gtest.AssertEQ(gconv.Float32(value), float32(123.456))
        gtest.AssertEQ(gconv.Float64(value), float64(123.456))
        gtest.AssertEQ(gconv.Bool(value),    true)
        gtest.AssertEQ(gconv.String(value),  "123.456")
    })
}

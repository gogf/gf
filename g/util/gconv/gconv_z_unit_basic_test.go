// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
)


func Test_Basic(t *testing.T) {
    gtest.Case(t, func() {
        vint   := float32(123.456)
        vint64 := int64(1552578474888)
        gtest.AssertEQ(gconv.Int(vint),       int(123))
        gtest.AssertEQ(gconv.Int8(vint),      int8(123))
        gtest.AssertEQ(gconv.Int16(vint),     int16(123))
        gtest.AssertEQ(gconv.Int32(vint),     int32(123))
        gtest.AssertEQ(gconv.Int64(vint),     int64(123))
        gtest.AssertEQ(gconv.Int64(vint),     int64(123))
        gtest.AssertEQ(gconv.Uint(vint),      uint(123))
        gtest.AssertEQ(gconv.Uint8(vint),     uint8(123))
        gtest.AssertEQ(gconv.Uint16(vint),    uint16(123))
        gtest.AssertEQ(gconv.Uint32(vint),    uint32(123))
        gtest.AssertEQ(gconv.Uint64(vint),    uint64(123))
        gtest.AssertEQ(gconv.Float32(vint),   float32(123.456))
        gtest.AssertEQ(gconv.Float64(vint),   float64(123.456))
        gtest.AssertEQ(gconv.Bool(vint),      true)
        gtest.AssertEQ(gconv.String(vint),    "123.456")
        gtest.AssertEQ(gconv.String(vint64),  "1552578474888")
    })
}

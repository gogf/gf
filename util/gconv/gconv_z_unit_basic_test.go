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

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		f32 := float32(123.456)
		i64 := int64(1552578474888)
		gtest.AssertEQ(gconv.Int(f32), int(123))
		gtest.AssertEQ(gconv.Int8(f32), int8(123))
		gtest.AssertEQ(gconv.Int16(f32), int16(123))
		gtest.AssertEQ(gconv.Int32(f32), int32(123))
		gtest.AssertEQ(gconv.Int64(f32), int64(123))
		gtest.AssertEQ(gconv.Int64(f32), int64(123))
		gtest.AssertEQ(gconv.Uint(f32), uint(123))
		gtest.AssertEQ(gconv.Uint8(f32), uint8(123))
		gtest.AssertEQ(gconv.Uint16(f32), uint16(123))
		gtest.AssertEQ(gconv.Uint32(f32), uint32(123))
		gtest.AssertEQ(gconv.Uint64(f32), uint64(123))
		gtest.AssertEQ(gconv.Float32(f32), float32(123.456))
		gtest.AssertEQ(gconv.Float64(i64), float64(i64))
		gtest.AssertEQ(gconv.Bool(f32), true)
		gtest.AssertEQ(gconv.String(f32), "123.456")
		gtest.AssertEQ(gconv.String(i64), "1552578474888")
	})

	gtest.Case(t, func() {
		s := "-0xFF"
		gtest.Assert(gconv.Int(s), int64(-0xFF))
	})
}

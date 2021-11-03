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

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		f32 := float32(123.456)
		i64 := int64(1552578474888)
		t.AssertEQ(gconv.Int(f32), int(123))
		t.AssertEQ(gconv.Int8(f32), int8(123))
		t.AssertEQ(gconv.Int16(f32), int16(123))
		t.AssertEQ(gconv.Int32(f32), int32(123))
		t.AssertEQ(gconv.Int64(f32), int64(123))
		t.AssertEQ(gconv.Int64(f32), int64(123))
		t.AssertEQ(gconv.Uint(f32), uint(123))
		t.AssertEQ(gconv.Uint8(f32), uint8(123))
		t.AssertEQ(gconv.Uint16(f32), uint16(123))
		t.AssertEQ(gconv.Uint32(f32), uint32(123))
		t.AssertEQ(gconv.Uint64(f32), uint64(123))
		t.AssertEQ(gconv.Float32(f32), float32(123.456))
		t.AssertEQ(gconv.Float64(i64), float64(i64))
		t.AssertEQ(gconv.Bool(f32), true)
		t.AssertEQ(gconv.String(f32), "123.456")
		t.AssertEQ(gconv.String(i64), "1552578474888")
	})

	gtest.C(t, func(t *gtest.T) {
		s := "-0xFF"
		t.Assert(gconv.Int(s), int64(-0xFF))
	})
}

func Test_Duration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		d := gconv.Duration("1s")
		t.Assert(d.String(), "1s")
		t.Assert(d.Nanoseconds(), 1000000000)
	})
}

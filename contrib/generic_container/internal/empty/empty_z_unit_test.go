// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package empty

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type TestInt int

type TestString string

type TestPerson interface {
	Say() string
}

type TestWoman struct {
}

func (woman TestWoman) Say() string {
	return "nice"
}

func TestIsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tmpT1 := "0"
		tmpT2 := func() {}
		tmpT2 = nil
		tmpT3 := make(chan int)
		var (
			tmpT4 TestPerson  = nil
			tmpT5 *TestPerson = nil
			tmpT6 TestPerson  = TestWoman{}
			tmpT7 TestInt     = 0
			tmpT8 TestString  = ""
		)
		tmpF1 := "1"
		tmpF2 := func(a string) string { return "1" }
		tmpF3 := make(chan int, 1)
		tmpF3 <- 1
		var (
			tmpF4 TestPerson = &TestWoman{}
			tmpF5 TestInt    = 1
			tmpF6 TestString = "1"
		)

		// true
		t.Assert(IsEmpty(nil), true)
		t.Assert(IsEmpty(0), true)
		t.Assert(IsEmpty(gconv.Int(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Int8(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Int16(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Int32(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Int64(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Uint64(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Uint(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Uint16(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Uint32(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Uint64(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Float32(tmpT1)), true)
		t.Assert(IsEmpty(gconv.Float64(tmpT1)), true)
		t.Assert(IsEmpty(false), true)
		t.Assert(IsEmpty([]byte("")), true)
		t.Assert(IsEmpty(""), true)
		t.Assert(IsEmpty(g.Map{}), true)
		t.Assert(IsEmpty(g.Slice{}), true)
		t.Assert(IsEmpty(g.Array{}), true)
		t.Assert(IsEmpty(tmpT2), true)
		t.Assert(IsEmpty(tmpT3), true)
		t.Assert(IsEmpty(tmpT3), true)
		t.Assert(IsEmpty(tmpT4), true)
		t.Assert(IsEmpty(tmpT5), true)
		t.Assert(IsEmpty(tmpT6), true)
		t.Assert(IsEmpty(tmpT7), true)
		t.Assert(IsEmpty(tmpT8), true)

		// false
		t.Assert(IsEmpty(gconv.Int(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Int8(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Int16(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Int32(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Int64(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Uint(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Uint8(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Uint16(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Uint32(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Uint64(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Float32(tmpF1)), false)
		t.Assert(IsEmpty(gconv.Float64(tmpF1)), false)
		t.Assert(IsEmpty(true), false)
		t.Assert(IsEmpty(tmpT1), false)
		t.Assert(IsEmpty([]byte("1")), false)
		t.Assert(IsEmpty(g.Map{"a": 1}), false)
		t.Assert(IsEmpty(g.Slice{"1"}), false)
		t.Assert(IsEmpty(g.Array{"1"}), false)
		t.Assert(IsEmpty(tmpF2), false)
		t.Assert(IsEmpty(tmpF3), false)
		t.Assert(IsEmpty(tmpF4), false)
		t.Assert(IsEmpty(tmpF5), false)
		t.Assert(IsEmpty(tmpF6), false)
	})
}

func TestIsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(IsNil(nil), true)
	})
	gtest.C(t, func(t *gtest.T) {
		var i int
		t.Assert(IsNil(i), false)
	})
	gtest.C(t, func(t *gtest.T) {
		var i *int
		t.Assert(IsNil(i), true)
	})
	gtest.C(t, func(t *gtest.T) {
		var i *int
		t.Assert(IsNil(&i), false)
		t.Assert(IsNil(&i, true), true)
	})
}

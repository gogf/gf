// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package empty_test

import (
	"testing"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/internal/empty"
	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gconv"
)

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
		tmpT3 := make(chan int, 0)
		var tmpT4 TestPerson = nil
		var tmpT5 *TestPerson = nil
		tmpF1 := "1"
		tmpF2 := func(a string) string { return "1" }
		tmpF3 := make(chan int, 1)
		tmpF3 <- 1
		var tmpF4 TestPerson = TestWoman{}
		tmpF5 := &tmpF4
		// true
		t.Assert(empty.IsEmpty(nil), true)
		t.Assert(empty.IsEmpty(gconv.Int(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Int8(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Int16(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Int32(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Int64(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Uint64(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Uint(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Uint16(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Uint32(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Uint64(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Float32(tmpT1)), true)
		t.Assert(empty.IsEmpty(gconv.Float64(tmpT1)), true)
		t.Assert(empty.IsEmpty(false), true)
		t.Assert(empty.IsEmpty([]byte("")), true)
		t.Assert(empty.IsEmpty(""), true)
		t.Assert(empty.IsEmpty(g.Map{}), true)
		t.Assert(empty.IsEmpty(g.Slice{}), true)
		t.Assert(empty.IsEmpty(g.Array{}), true)
		t.Assert(empty.IsEmpty(tmpT2), true)
		t.Assert(empty.IsEmpty(tmpT3), true)
		t.Assert(empty.IsEmpty(tmpT3), true)
		t.Assert(empty.IsEmpty(tmpT4), true)
		t.Assert(empty.IsEmpty(tmpT5), true)
		// false
		t.Assert(empty.IsEmpty(gconv.Int(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Int8(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Int16(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Int32(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Int64(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Uint(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Uint8(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Uint16(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Uint32(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Uint64(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Float32(tmpF1)), false)
		t.Assert(empty.IsEmpty(gconv.Float64(tmpF1)), false)
		t.Assert(empty.IsEmpty(true), false)
		t.Assert(empty.IsEmpty(tmpT1), false)
		t.Assert(empty.IsEmpty([]byte("1")), false)
		t.Assert(empty.IsEmpty(g.Map{"a": 1}), false)
		t.Assert(empty.IsEmpty(g.Slice{"1"}), false)
		t.Assert(empty.IsEmpty(g.Array{"1"}), false)
		t.Assert(empty.IsEmpty(tmpF2), false)
		t.Assert(empty.IsEmpty(tmpF3), false)
		t.Assert(empty.IsEmpty(tmpF4), false)
		t.Assert(empty.IsEmpty(tmpF5), false)
	})
}

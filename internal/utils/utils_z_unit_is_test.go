// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestVar_IsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsNil(0), false)
		t.Assert(utils.IsNil(nil), true)
		t.Assert(utils.IsNil(g.Map{}), false)
		t.Assert(utils.IsNil(g.Slice{}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsNil(1), false)
		t.Assert(utils.IsNil(0.1), false)
		t.Assert(utils.IsNil(g.Map{"k": "v"}), false)
		t.Assert(utils.IsNil(g.Slice{0}), false)
	})
}

func TestVar_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsEmpty(0), true)
		t.Assert(utils.IsEmpty(nil), true)
		t.Assert(utils.IsEmpty(g.Map{}), true)
		t.Assert(utils.IsEmpty(g.Slice{}), true)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsEmpty(1), false)
		t.Assert(utils.IsEmpty(0.1), false)
		t.Assert(utils.IsEmpty(g.Map{"k": "v"}), false)
		t.Assert(utils.IsEmpty(g.Slice{0}), false)
	})
}

func TestVar_IsInt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsInt(0), true)
		t.Assert(utils.IsInt(nil), false)
		t.Assert(utils.IsInt(g.Map{}), false)
		t.Assert(utils.IsInt(g.Slice{}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsInt(1), true)
		t.Assert(utils.IsInt(-1), true)
		t.Assert(utils.IsInt(0.1), false)
		t.Assert(utils.IsInt(g.Map{"k": "v"}), false)
		t.Assert(utils.IsInt(g.Slice{0}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsInt(int8(1)), true)
		t.Assert(utils.IsInt(uint8(1)), false)
	})
}

func TestVar_IsUint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsUint(0), false)
		t.Assert(utils.IsUint(nil), false)
		t.Assert(utils.IsUint(g.Map{}), false)
		t.Assert(utils.IsUint(g.Slice{}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsUint(1), false)
		t.Assert(utils.IsUint(-1), false)
		t.Assert(utils.IsUint(0.1), false)
		t.Assert(utils.IsUint(g.Map{"k": "v"}), false)
		t.Assert(utils.IsUint(g.Slice{0}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsUint(int8(1)), false)
		t.Assert(utils.IsUint(uint8(1)), true)
	})
}

func TestVar_IsFloat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsFloat(0), false)
		t.Assert(utils.IsFloat(nil), false)
		t.Assert(utils.IsFloat(g.Map{}), false)
		t.Assert(utils.IsFloat(g.Slice{}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsFloat(1), false)
		t.Assert(utils.IsFloat(-1), false)
		t.Assert(utils.IsFloat(0.1), true)
		t.Assert(utils.IsFloat(float64(1)), true)
		t.Assert(utils.IsFloat(g.Map{"k": "v"}), false)
		t.Assert(utils.IsFloat(g.Slice{0}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsFloat(int8(1)), false)
		t.Assert(utils.IsFloat(uint8(1)), false)
	})
}

func TestVar_IsSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsSlice(0), false)
		t.Assert(utils.IsSlice(nil), false)
		t.Assert(utils.IsSlice(g.Map{}), false)
		t.Assert(utils.IsSlice(g.Slice{}), true)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsSlice(1), false)
		t.Assert(utils.IsSlice(-1), false)
		t.Assert(utils.IsSlice(0.1), false)
		t.Assert(utils.IsSlice(float64(1)), false)
		t.Assert(utils.IsSlice(g.Map{"k": "v"}), false)
		t.Assert(utils.IsSlice(g.Slice{0}), true)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsSlice(int8(1)), false)
		t.Assert(utils.IsSlice(uint8(1)), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsSlice(gvar.New(gtime.Now()).IsSlice()), false)
	})
}

func TestVar_IsMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsMap(0), false)
		t.Assert(utils.IsMap(nil), false)
		t.Assert(utils.IsMap(g.Map{}), true)
		t.Assert(utils.IsMap(g.Slice{}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsMap(1), false)
		t.Assert(utils.IsMap(-1), false)
		t.Assert(utils.IsMap(0.1), false)
		t.Assert(utils.IsMap(float64(1)), false)
		t.Assert(utils.IsMap(g.Map{"k": "v"}), true)
		t.Assert(utils.IsMap(g.Slice{0}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsMap(int8(1)), false)
		t.Assert(utils.IsMap(uint8(1)), false)
	})
}

func TestVar_IsStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsStruct(0), false)
		t.Assert(utils.IsStruct(nil), false)
		t.Assert(utils.IsStruct(g.Map{}), false)
		t.Assert(utils.IsStruct(g.Slice{}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(utils.IsStruct(1), false)
		t.Assert(utils.IsStruct(-1), false)
		t.Assert(utils.IsStruct(0.1), false)
		t.Assert(utils.IsStruct(float64(1)), false)
		t.Assert(utils.IsStruct(g.Map{"k": "v"}), false)
		t.Assert(utils.IsStruct(g.Slice{0}), false)
	})
	gtest.C(t, func(t *gtest.T) {
		a := &struct {
		}{}
		t.Assert(utils.IsStruct(a), true)
		t.Assert(utils.IsStruct(*a), true)
		t.Assert(utils.IsStruct(&a), true)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestVar_IsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsNil(), false)
		t.Assert(g.NewVar(nil).IsNil(), true)
		t.Assert(g.NewVar(g.Map{}).IsNil(), false)
		t.Assert(g.NewVar(g.Slice{}).IsNil(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsNil(), false)
		t.Assert(g.NewVar(0.1).IsNil(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsNil(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsNil(), false)
	})
}

func TestVar_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsEmpty(), true)
		t.Assert(g.NewVar(nil).IsEmpty(), true)
		t.Assert(g.NewVar(g.Map{}).IsEmpty(), true)
		t.Assert(g.NewVar(g.Slice{}).IsEmpty(), true)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsEmpty(), false)
		t.Assert(g.NewVar(0.1).IsEmpty(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsEmpty(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsEmpty(), false)
	})
}

func TestVar_IsInt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsInt(), true)
		t.Assert(g.NewVar(nil).IsInt(), false)
		t.Assert(g.NewVar(g.Map{}).IsInt(), false)
		t.Assert(g.NewVar(g.Slice{}).IsInt(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsInt(), true)
		t.Assert(g.NewVar(-1).IsInt(), true)
		t.Assert(g.NewVar(0.1).IsInt(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsInt(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsInt(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(int8(1)).IsInt(), true)
		t.Assert(g.NewVar(uint8(1)).IsInt(), false)
	})
}

func TestVar_IsUint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsUint(), false)
		t.Assert(g.NewVar(nil).IsUint(), false)
		t.Assert(g.NewVar(g.Map{}).IsUint(), false)
		t.Assert(g.NewVar(g.Slice{}).IsUint(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsUint(), false)
		t.Assert(g.NewVar(-1).IsUint(), false)
		t.Assert(g.NewVar(0.1).IsUint(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsUint(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsUint(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(int8(1)).IsUint(), false)
		t.Assert(g.NewVar(uint8(1)).IsUint(), true)
	})
}

func TestVar_IsFloat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsFloat(), false)
		t.Assert(g.NewVar(nil).IsFloat(), false)
		t.Assert(g.NewVar(g.Map{}).IsFloat(), false)
		t.Assert(g.NewVar(g.Slice{}).IsFloat(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsFloat(), false)
		t.Assert(g.NewVar(-1).IsFloat(), false)
		t.Assert(g.NewVar(0.1).IsFloat(), true)
		t.Assert(g.NewVar(float64(1)).IsFloat(), true)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsFloat(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsFloat(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(int8(1)).IsFloat(), false)
		t.Assert(g.NewVar(uint8(1)).IsFloat(), false)
	})
}

func TestVar_IsSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsSlice(), false)
		t.Assert(g.NewVar(nil).IsSlice(), false)
		t.Assert(g.NewVar(g.Map{}).IsSlice(), false)
		t.Assert(g.NewVar(g.Slice{}).IsSlice(), true)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsSlice(), false)
		t.Assert(g.NewVar(-1).IsSlice(), false)
		t.Assert(g.NewVar(0.1).IsSlice(), false)
		t.Assert(g.NewVar(float64(1)).IsSlice(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsSlice(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsSlice(), true)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(int8(1)).IsSlice(), false)
		t.Assert(g.NewVar(uint8(1)).IsSlice(), false)
	})
}

func TestVar_IsMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsMap(), false)
		t.Assert(g.NewVar(nil).IsMap(), false)
		t.Assert(g.NewVar(g.Map{}).IsMap(), true)
		t.Assert(g.NewVar(g.Slice{}).IsMap(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsMap(), false)
		t.Assert(g.NewVar(-1).IsMap(), false)
		t.Assert(g.NewVar(0.1).IsMap(), false)
		t.Assert(g.NewVar(float64(1)).IsMap(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsMap(), true)
		t.Assert(g.NewVar(g.Slice{0}).IsMap(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(int8(1)).IsMap(), false)
		t.Assert(g.NewVar(uint8(1)).IsMap(), false)
	})
}

func TestVar_IsStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(0).IsStruct(), false)
		t.Assert(g.NewVar(nil).IsStruct(), false)
		t.Assert(g.NewVar(g.Map{}).IsStruct(), false)
		t.Assert(g.NewVar(g.Slice{}).IsStruct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).IsStruct(), false)
		t.Assert(g.NewVar(-1).IsStruct(), false)
		t.Assert(g.NewVar(0.1).IsStruct(), false)
		t.Assert(g.NewVar(float64(1)).IsStruct(), false)
		t.Assert(g.NewVar(g.Map{"k": "v"}).IsStruct(), false)
		t.Assert(g.NewVar(g.Slice{0}).IsStruct(), false)
	})
	gtest.C(t, func(t *gtest.T) {
		a := &struct {
		}{}
		t.Assert(g.NewVar(a).IsStruct(), true)
		t.Assert(g.NewVar(*a).IsStruct(), true)
		t.Assert(g.NewVar(&a).IsStruct(), true)
	})
}

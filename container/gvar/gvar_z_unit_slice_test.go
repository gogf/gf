// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestVarInts(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Ints()[0], arr[0])
	})
}

func TestVarUints(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Uints()[0], arr[0])
	})
}

func TestVarInt64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Int64s()[0], arr[0])
	})
}

func TestVarUint64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Uint64s()[0], arr[0])
	})
}

func TestVarFloats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Floats()[0], arr[0])
	})
}

func TestVarFloat32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []float32{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.AssertEQ(objOne.Float32s(), arr)
	})
}

func TestVarFloat64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.AssertEQ(objOne.Float64s(), arr)
	})
}

func TestVarStrings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []string{"hello", "world"}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Strings()[0], arr[0])
	})
}

func TestVarInterfaces(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Interfaces(), arr)
	})
}

func TestVarSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Slice(), arr)
	})
}

func TestVarArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, false)
		t.Assert(objOne.Array(), arr)
	})
}

func TestVarVars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, false)
		t.Assert(len(objOne.Vars()), 5)
		t.Assert(objOne.Vars()[0].Int(), 1)
		t.Assert(objOne.Vars()[4].Int(), 5)

		objEmpty := gvar.New([]int{})
		t.Assert(objEmpty.Vars(), nil)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package reflection_test

import (
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_OriginValueAndKind(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s = "s"
		out := reflection.OriginValueAndKind(s)
		t.Assert(out.InputKind, reflect.String)
		t.Assert(out.OriginKind, reflect.String)
	})
	gtest.C(t, func(t *gtest.T) {
		var s = "s"
		out := reflection.OriginValueAndKind(&s)
		t.Assert(out.InputKind, reflect.Pointer)
		t.Assert(out.OriginKind, reflect.String)
	})
	gtest.C(t, func(t *gtest.T) {
		var s []int
		out := reflection.OriginValueAndKind(s)
		t.Assert(out.InputKind, reflect.Slice)
		t.Assert(out.OriginKind, reflect.Slice)
	})
	gtest.C(t, func(t *gtest.T) {
		var s []int
		out := reflection.OriginValueAndKind(&s)
		t.Assert(out.InputKind, reflect.Pointer)
		t.Assert(out.OriginKind, reflect.Slice)
	})
}

func Test_OriginTypeAndKind(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var s = "s"
		out := reflection.OriginTypeAndKind(s)
		t.Assert(out.InputKind, reflect.String)
		t.Assert(out.OriginKind, reflect.String)
	})
	gtest.C(t, func(t *gtest.T) {
		var s = "s"
		out := reflection.OriginTypeAndKind(&s)
		t.Assert(out.InputKind, reflect.Pointer)
		t.Assert(out.OriginKind, reflect.String)
	})
	gtest.C(t, func(t *gtest.T) {
		var s []int
		out := reflection.OriginTypeAndKind(s)
		t.Assert(out.InputKind, reflect.Slice)
		t.Assert(out.OriginKind, reflect.Slice)
	})
	gtest.C(t, func(t *gtest.T) {
		var s []int
		out := reflection.OriginTypeAndKind(&s)
		t.Assert(out.InputKind, reflect.Pointer)
		t.Assert(out.OriginKind, reflect.Slice)
	})
}

func Test_IsBasicKind(t *testing.T) {
	// All basic kinds should return true.
	basicKinds := []reflect.Kind{
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool, reflect.String,
	}
	for _, kind := range basicKinds {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(reflection.IsBasicKind(kind), true)
		})
	}
	// Non-basic kinds should return false.
	nonBasicKinds := []reflect.Kind{
		reflect.Array, reflect.Slice, reflect.Map, reflect.Struct,
		reflect.Pointer, reflect.Interface, reflect.Chan, reflect.Func,
		reflect.Complex64, reflect.Complex128, reflect.Uintptr,
	}
	for _, kind := range nonBasicKinds {
		gtest.C(t, func(t *gtest.T) {
			t.Assert(reflection.IsBasicKind(kind), false)
		})
	}
}

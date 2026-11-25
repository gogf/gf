// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// This file documents the known behavior of pointer slice element converting that is
// reported in the review of the "bind struct pointer slice elements by value type" change.
// The cases below intentionally assert the current behavior instead of the expected one,
// so that any future fix has to update these assertions explicitly.

package gconv_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

// Duration is a user-defined named type whose type name collides with the
// built-in "Duration"/"time.Duration" case in doConvert.
type Duration int64

// Time is a user-defined named type whose type name collides with the
// built-in "Time"/"time.Time" case in doConvert.
type Time int64

// TestStructPointerSliceRawMessageStringInput documents how json.RawMessage pointer
// elements are encoded when the input element is a raw JSON string.
//
// Known issue: the element type name is resolved to "json.RawMessage", which routes the
// converting through doConvert's json.Marshal case. The original JSON text is therefore
// re-encoded into a JSON string instead of being passed through to UnmarshalJSON.
func TestStructPointerSliceRawMessageStringInput(t *testing.T) {
	input := map[string]any{
		"Values": []string{`{"a":1}`, `2`},
	}
	gtest.C(t, func(t *gtest.T) {
		// Pointer elements: raw JSON text is wrapped into a JSON string.
		var target struct{ Values []*json.RawMessage }
		err := gconv.NewConverter().Struct(input, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 2)
		t.Assert(string(*target.Values[0]), `"{\"a\":1}"`)
		t.Assert(string(*target.Values[1]), `"2"`)
	})
	gtest.C(t, func(t *gtest.T) {
		// Value elements behave exactly the same, which is the consistency this change brings.
		var target struct{ Values []json.RawMessage }
		err := gconv.NewConverter().Struct(input, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 2)
		t.Assert(string(target.Values[0]), `"{\"a\":1}"`)
		t.Assert(string(target.Values[1]), `"2"`)
	})
}

// TestStructPointerSliceRawMessageObjectInput documents that map input keeps the
// expected json.Marshal semantics for both pointer and value elements.
func TestStructPointerSliceRawMessageObjectInput(t *testing.T) {
	input := map[string]any{
		"Values": []any{map[string]any{"hello": "world"}},
	}
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*json.RawMessage }
		err := gconv.NewConverter().Struct(input, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 1)
		t.Assert(string(*target.Values[0]), `{"hello":"world"}`)
	})
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []json.RawMessage }
		err := gconv.NewConverter().Struct(input, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 1)
		t.Assert(string(target.Values[0]), `{"hello":"world"}`)
	})
}

// TestStructPointerSliceCustomNamedTypeCollision documents that user-defined named types
// which collide with a doConvert case name cannot be converted.
//
// Known issue: elem.Type().Name() yields "Duration"/"Time" without a package path, so the
// converting returns a time.Duration/time.Time value that cannot be assigned to the
// user-defined target type. Non-pointer elements already behaved this way before the change.
func TestStructPointerSliceCustomNamedTypeCollision(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Pointer elements: used to work through the "*gf_test.Duration" default path.
		var target struct{ Values []*Duration }
		err := gconv.NewConverter().Struct(map[string]any{"Values": []any{100, 200}}, &target)
		t.AssertNE(err, nil)
		t.Assert(len(target.Values), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		// Value elements: same failure, which is a pre-existing issue.
		var target struct{ Values []Duration }
		err := gconv.NewConverter().Struct(map[string]any{"Values": []any{100, 200}}, &target)
		t.AssertNE(err, nil)
		t.Assert(len(target.Values), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*Time }
		err := gconv.NewConverter().Struct(map[string]any{"Values": []any{100}}, &target)
		t.AssertNE(err, nil)
		t.Assert(len(target.Values), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []Time }
		err := gconv.NewConverter().Struct(map[string]any{"Values": []any{100}}, &target)
		t.AssertNE(err, nil)
		t.Assert(len(target.Values), 0)
	})
}

// TestStructPointerSliceTimeTypes covers the standard library and framework time types
// for both collection and scalar input, which the added tests of the change do not cover.
func TestStructPointerSliceTimeTypes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*time.Time }
		err := gconv.NewConverter().Struct(map[string]any{
			"Values": []string{"2020-01-01 00:00:00", "2021-02-03 04:05:06"},
		}, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 2)
		t.Assert(target.Values[0].Format("2006-01-02 15:04:05"), "2020-01-01 00:00:00")
		t.Assert(target.Values[1].Format("2006-01-02 15:04:05"), "2021-02-03 04:05:06")
	})
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*gtime.Time }
		err := gconv.NewConverter().Struct(map[string]any{
			"Values": []string{"2020-01-01 00:00:00", "2021-02-03 04:05:06"},
		}, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 2)
		t.Assert(target.Values[0].String(), "2020-01-01 00:00:00")
		t.Assert(target.Values[1].String(), "2021-02-03 04:05:06")
	})
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*time.Duration }
		err := gconv.NewConverter().Struct(map[string]any{
			"Values": []string{"1s", "2m"},
		}, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 2)
		t.Assert(*target.Values[0], time.Second)
		t.Assert(*target.Values[1], 2*time.Minute)
	})
	gtest.C(t, func(t *gtest.T) {
		// Scalar input converts to a single-element slice.
		var target struct{ Values []*time.Time }
		err := gconv.NewConverter().Struct(map[string]any{
			"Values": "2020-01-01 00:00:00",
		}, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 1)
		t.Assert(target.Values[0].Format("2006-01-02 15:04:05"), "2020-01-01 00:00:00")
	})
}

// TestStructPointerSliceNestedPointer documents that double pointer elements cannot be
// converted: the error is reported by the strict converter but silently swallowed by
// gconv.Struct, which is left as-is by the change.
func TestStructPointerSliceNestedPointer(t *testing.T) {
	input := map[string]any{"Values": [][]int{{1, 2}}}
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []**int }
		err := gconv.NewConverter().Struct(input, &target)
		t.AssertNE(err, nil)
		t.Assert(len(target.Values), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		// gconv.Struct continues on error and keeps the field untouched.
		var target struct{ Values []**int }
		err := gconv.Struct(input, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 0)
	})
}

// TestStructPointerSlicePointerElementsOfPointerSlices documents converting into a slice
// whose elements are pointers to slices and maps.
func TestStructPointerSlicePointerElementsOfPointerSlices(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*[]int }
		err := gconv.NewConverter().Struct(map[string]any{
			"Values": []any{[]int{3, 4}},
		}, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 1)
		t.Assert(*target.Values[0], []int{3, 4})
	})
	gtest.C(t, func(t *gtest.T) {
		var target struct{ Values []*map[string]int }
		err := gconv.NewConverter().Struct(map[string]any{
			"Values": []any{map[string]any{"a": 1}},
		}, &target)
		t.AssertNil(err)
		t.Assert(len(target.Values), 1)
		t.Assert((*target.Values[0])["a"], 1)
	})
}

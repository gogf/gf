// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

type testScan struct {
	Src    any
	Dst    any
	Expect any
}

func TestScanBasicTypes(t *testing.T) {
	// Define test data structure
	type User struct {
		Name string
		Age  int
	}
	type UserWithTag struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// Prepare test data
	var testScanData = []testScan{
		// Basic type conversion
		{1, new(int), 1},
		{int8(1), new(int16), int16(1)},
		{int16(1), new(int32), int32(1)},
		{int32(1), new(int64), int64(1)},
		{uint(1), new(int), 1},
		{uint8(1), new(int), 1},
		{uint16(1), new(int), 1},
		{uint32(1), new(int), 1},
		{uint64(1), new(int), 1},
		{float32(1.0), new(int), 1},
		{float64(1.0), new(int), 1},
		{true, new(int), 1},
		{false, new(int), 0},
		{"1", new(int), 1},
		{"true", new(bool), true},
		{"false", new(bool), false},
		{1, new(bool), true},
		{0, new(bool), false},

		// String conversion
		{1, new(string), "1"},
		{1.1, new(string), "1.1"},
		{true, new(string), "true"},
		{false, new(string), "false"},
		{[]byte("hello"), new(string), "hello"},

		// Slice conversion
		{[]int{1, 2, 3}, new([]string), []string{"1", "2", "3"}},
		{[]string{"1", "2", "3"}, new([]int), []int{1, 2, 3}},
		{`["1","2","3"]`, new([]string), []string{"1", "2", "3"}},
		{`[1,2,3]`, new([]int), []int{1, 2, 3}},

		// Map conversion
		{
			map[string]any{"name": "john", "age": 18},
			new(User),
			&User{Name: "john", Age: 18},
		},
		{
			`{"name":"john","age":18}`,
			new(User),
			&User{Name: "john", Age: 18},
		},
		{
			map[string]any{"name": "john", "age": 18},
			new(UserWithTag),
			&UserWithTag{Name: "john", Age: 18},
		},
		{
			map[string]string{"name": "john", "age": "18"},
			new(map[string]any),
			&map[string]any{"name": "john", "age": "18"},
		},

		// Struct conversion
		{
			User{Name: "john", Age: 18},
			new(map[string]any),
			&map[string]any{"Name": "john", "Age": 18},
		},
		{
			&User{Name: "john", Age: 18},
			new(UserWithTag),
			&UserWithTag{Name: "john", Age: 18},
		},

		// Special cases
		{nil, new(any), nil},
		{nil, new(*int), (*int)(nil)},
		{[]byte(nil), new(string), ""},
		{"", new(int), 0},
		{"", new(float64), 0.0},
		{"", new(bool), false},

		// Time type
		{time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local), new(string), "2023-01-02 00:00:00"},

		// Pointer conversion
		{&User{Name: "john"}, new(*User), &User{Name: "john"}},
	}

	// Basic types test.
	gtest.C(t, func(t *gtest.T) {
		for _, v := range testScanData {
			// t.Logf(`%#v`, v)
			err := gconv.Scan(v.Src, v.Dst)
			t.AssertNil(err)
		}
	})

	// int -> **int
	gtest.C(t, func(t *gtest.T) {
		var (
			v = 100
			i *int
		)
		err := gconv.Scan(v, &i)
		t.AssertNil(err)
		t.AssertNE(i, nil)
		t.Assert(*i, v)
	})
	// *int -> **int
	gtest.C(t, func(t *gtest.T) {
		var (
			v = 100
			i *int
		)
		err := gconv.Scan(&v, &i)
		t.AssertNil(err)
		t.AssertNE(i, nil)
		t.Assert(*i, v)
	})
	// string -> **string
	gtest.C(t, func(t *gtest.T) {
		var (
			v = "1000"
			i *string
		)
		err := gconv.Scan(v, &i)
		t.AssertNil(err)
		t.AssertNE(i, nil)
		t.Assert(*i, v)
	})
	// *string -> **string
	gtest.C(t, func(t *gtest.T) {
		var (
			v = "1000"
			i *string
		)
		err := gconv.Scan(&v, &i)
		t.AssertNil(err)
		t.AssertNE(i, nil)
		t.Assert(*i, v)
	})
}

// TestScanReflectValueInput tests that the basic converter functions correctly unwrap
// reflect.Value inputs instead of treating the wrapper as the converted value.
func TestScanReflectValueInput(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// reflect.Value wrapping a bool.
		t.Assert(gconv.Bool(reflect.ValueOf(true)), true)
		t.Assert(gconv.Bool(reflect.ValueOf(false)), false)
		// reflect.Value wrapping an int.
		t.Assert(gconv.Int64(reflect.ValueOf(42)), int64(42))
		t.Assert(gconv.Int(reflect.ValueOf(42)), 42)
		// reflect.Value wrapping a uint.
		t.Assert(gconv.Uint64(reflect.ValueOf(uint(42))), uint64(42))
		// reflect.Value wrapping a float.
		t.Assert(gconv.Float64(reflect.ValueOf(3.14)), 3.14)
		t.Assert(gconv.Float32(reflect.ValueOf(float32(1.5))), float32(1.5))
		// reflect.Value wrapping a string.
		t.Assert(gconv.String(reflect.ValueOf("hello")), "hello")
		// reflect.Value wrapping non-string kinds should unwrap and convert correctly,
		// not return the reflection type name (e.g. "<int Value>").
		t.Assert(gconv.String(reflect.ValueOf(42)), "42")
		t.Assert(gconv.String(reflect.ValueOf(3.14)), "3.14")
		t.Assert(gconv.String(reflect.ValueOf(true)), "true")
	})
}

// TestScanPointerElementSlice tests scanning into slices whose elements are pointers to
// basic types (e.g. []*int, []*string), which is supported by converter_scan.go.
func TestScanPointerElementSlice(t *testing.T) {
	// []string -> []*string
	gtest.C(t, func(t *gtest.T) {
		var (
			src = []string{"a", "b", "c"}
			dst []*string
		)
		err := gconv.Scan(src, &dst)
		t.AssertNil(err)
		t.Assert(len(dst), len(src))
		for i, v := range src {
			t.AssertNE(dst[i], nil)
			t.Assert(*dst[i], v)
		}
	})
	// []int -> []*int
	gtest.C(t, func(t *gtest.T) {
		var (
			src = []int{1, 2, 3}
			dst []*int
		)
		err := gconv.Scan(src, &dst)
		t.AssertNil(err)
		t.Assert(len(dst), len(src))
		for i, v := range src {
			t.AssertNE(dst[i], nil)
			t.Assert(*dst[i], v)
		}
	})
	// []int64 -> []*int64
	gtest.C(t, func(t *gtest.T) {
		var (
			src = []int64{10, 20, 30}
			dst []*int64
		)
		err := gconv.Scan(src, &dst)
		t.AssertNil(err)
		t.Assert(len(dst), len(src))
		for i, v := range src {
			t.AssertNE(dst[i], nil)
			t.Assert(*dst[i], v)
		}
	})
	// []float64 -> []*float64
	gtest.C(t, func(t *gtest.T) {
		var (
			src = []float64{1.1, 2.2, 3.3}
			dst []*float64
		)
		err := gconv.Scan(src, &dst)
		t.AssertNil(err)
		t.Assert(len(dst), len(src))
		for i, v := range src {
			t.AssertNE(dst[i], nil)
			t.Assert(*dst[i], v)
		}
	})
	// []bool -> []*bool
	gtest.C(t, func(t *gtest.T) {
		var (
			src = []bool{true, false, true}
			dst []*bool
		)
		err := gconv.Scan(src, &dst)
		t.AssertNil(err)
		t.Assert(len(dst), len(src))
		for i, v := range src {
			t.AssertNE(dst[i], nil)
			t.Assert(*dst[i], v)
		}
	})
	// []string -> []*int (cross-type pointer element conversion)
	gtest.C(t, func(t *gtest.T) {
		var (
			src = []string{"1", "2", "3"}
			dst []*int
		)
		err := gconv.Scan(src, &dst)
		t.AssertNil(err)
		t.Assert(len(dst), len(src))
		for i, v := range []int{1, 2, 3} {
			t.AssertNE(dst[i], nil)
			t.Assert(*dst[i], v)
		}
	})
}

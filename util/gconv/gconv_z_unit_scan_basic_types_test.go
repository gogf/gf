// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
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
			map[string]interface{}{"name": "john", "age": 18},
			new(User),
			&User{Name: "john", Age: 18},
		},
		{
			`{"name":"john","age":18}`,
			new(User),
			&User{Name: "john", Age: 18},
		},
		{
			map[string]interface{}{"name": "john", "age": 18},
			new(UserWithTag),
			&UserWithTag{Name: "john", Age: 18},
		},
		{
			map[string]string{"name": "john", "age": "18"},
			new(map[string]interface{}),
			&map[string]interface{}{"name": "john", "age": "18"},
		},

		// Struct conversion
		{
			User{Name: "john", Age: 18},
			new(map[string]interface{}),
			&map[string]interface{}{"Name": "john", "Age": 18},
		},
		{
			&User{Name: "john", Age: 18},
			new(UserWithTag),
			&UserWithTag{Name: "john", Age: 18},
		},

		// Special cases
		{nil, new(interface{}), nil},
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

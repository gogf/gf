// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtest_test

import (
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestC(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(float32(123.456), float32(123.456))
	})
}

func TestCase(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(1, 1)
		t.AssertNE(1, 0)
		t.AssertEQ(float32(123.456), float32(123.456))
		t.AssertEQ(float32(123.456), float32(123.456))
	})
}

func TestAssert(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			nilChan chan struct{}
		)
		t.Assert(1, 1)
		t.Assert(nilChan, nil)
		m1 := map[string]string{"k1": "v1", "k2": "v2"}
		m2 := map[string]string{"k2": "v2", "k1": "v1"}
		t.Assert(m1, m2)
	})
}

func TestAssertEQ(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			nilChan chan struct{}
		)
		t.AssertEQ(nilChan, nil)
		t.AssertEQ("0", "0")
		t.AssertEQ(float32(123.456), float32(123.456))
		m1 := map[string]string{"k1": "v1", "k2": "v2"}
		m2 := map[string]string{"k2": "v2", "k1": "v1"}
		t.AssertEQ(m1, m2)
	})
}

func TestAssertNE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c = make(chan struct{}, 1)
		)
		t.AssertNE(nil, c)
		t.AssertNE("0", "1")
		t.AssertNE(float32(123.456), float32(123.4567))
		m1 := map[string]string{"k1": "v1", "k2": "v2"}
		m2 := map[string]string{"k2": "v1", "k1": "v2"}
		t.AssertNE(m1, m2)
	})
}

func TestAssertNQ(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNQ(1, "0")
		t.AssertNQ(float32(123.456), float64(123.4567))
	})
}

func TestAssertGT(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT("b", "a")
		t.AssertGT(1, -1)
		t.AssertGT(uint(1), uint(0))
		t.AssertGT(float32(123.45678), float32(123.4567))
	})
}

func TestAssertGE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGE("b", "a")
		t.AssertGE("a", "a")
		t.AssertGE(1, -1)
		t.AssertGE(1, 1)
		t.AssertGE(uint(1), uint(0))
		t.AssertGE(uint(0), uint(0))
		t.AssertGE(float32(123.45678), float32(123.4567))
		t.AssertGE(float32(123.456), float32(123.456))
	})
}

func TestAssertLT(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertLT("a", "b")
		t.AssertLT(-1, 1)
		t.AssertLT(uint(0), uint(1))
		t.AssertLT(float32(123.456), float32(123.4567))
	})
}

func TestAssertLE(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertLE("a", "b")
		t.AssertLE("a", "a")
		t.AssertLE(-1, 1)
		t.AssertLE(1, 1)
		t.AssertLE(uint(0), uint(1))
		t.AssertLE(uint(0), uint(0))
		t.AssertLE(float32(123.456), float32(123.4567))
		t.AssertLE(float32(123.456), float32(123.456))
	})
}

func TestAssertIN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertIN("a", []string{"a", "b", "c"})
		t.AssertIN(1, []int{1, 2, 3})
	})
}

func TestAssertNI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNI("d", []string{"a", "b", "c"})
		t.AssertNI(4, []int{1, 2, 3})
	})
}

func TestAssertNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			nilChan chan struct{}
		)
		t.AssertNil(nilChan)
		_, err := strconv.ParseInt("123", 10, 64)
		t.AssertNil(err)
	})
}

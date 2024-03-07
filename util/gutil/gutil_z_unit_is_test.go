// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_IsTypeOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsTypeOf(1, 0), true)
		t.Assert(gutil.IsTypeOf(1.1, 0.1), true)
		t.Assert(gutil.IsTypeOf(1.1, 1), false)
		t.Assert(gutil.IsTypeOf(true, false), true)
		t.Assert(gutil.IsTypeOf(true, 1), false)
	})
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name string
		}
		type B struct {
			Name string
		}
		t.Assert(gutil.IsTypeOf(1, A{}), false)
		t.Assert(gutil.IsTypeOf(A{}, B{}), false)
		t.Assert(gutil.IsTypeOf(A{Name: "john"}, &A{Name: "john"}), false)
		t.Assert(gutil.IsTypeOf(A{Name: "john"}, A{Name: "john"}), true)
		t.Assert(gutil.IsTypeOf(A{Name: "john"}, A{}), true)
		t.Assert(gutil.IsTypeOf(&A{Name: "john"}, &A{}), true)
		t.Assert(gutil.IsTypeOf(&A{Name: "john"}, &B{}), false)
		t.Assert(gutil.IsTypeOf(A{Name: "john"}, B{Name: "john"}), false)
	})
}

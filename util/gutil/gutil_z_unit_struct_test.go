// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_StructToSlice(t *testing.T) {
	type A struct {
		K1 int
		K2 string
	}
	gtest.C(t, func(t *gtest.T) {
		a := &A{
			K1: 1,
			K2: "v2",
		}
		s := gutil.StructToSlice(a)
		t.Assert(len(s), 4)
		t.AssertIN(s[0], g.Slice{"K1", "K2", 1, "v2"})
		t.AssertIN(s[1], g.Slice{"K1", "K2", 1, "v2"})
		t.AssertIN(s[2], g.Slice{"K1", "K2", 1, "v2"})
		t.AssertIN(s[3], g.Slice{"K1", "K2", 1, "v2"})
	})
	gtest.C(t, func(t *gtest.T) {
		s := gutil.StructToSlice(1)
		t.Assert(s, nil)
	})
}

func Test_FillStructWithDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			A1 int     `d:"1.01"`
			A2 string  `d:"1.01"`
			A3 float32 `d:"1.01"`
		}
		type B struct {
			B1 bool `d:"true"`
			B2 string
			A  A
		}
		type C struct {
			C1 float64 `d:"1.01"`
			B
			C2 bool
			A  A
		}

		c := C{}
		err := gutil.FillStructWithDefault(&c)
		t.AssertNil(err)

		t.Assert(c.C1, `1.01`)
		t.Assert(c.C2, false)
		t.Assert(c.B1, true)
		t.Assert(c.B2, ``)
		t.Assert(c.A.A1, `1`)
		t.Assert(c.A.A2, `1.01`)
		t.Assert(c.A.A3, `1.01`)
	})
}

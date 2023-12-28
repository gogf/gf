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
		type myInt int
		type Inner1 struct {
			I1V1 int
			I1V2 bool `d:"true"`
		}
		type Inner2 struct {
			I2V1 float64 `d:"1.01"`
		}
		type Inner3 struct {
			Inner1 Inner1
			I3V1   myInt `d:"1"`
		}
		type Inner4 struct {
		}
		type Outer struct {
			O1 int     `d:"1.01"`
			O2 string  `d:"1.01"`
			O3 float32 `d:"1.01"`
			*Inner1
			O4 bool `d:"true"`
			Inner2
			Inner3 Inner3
			Inner4 *Inner4
		}

		outer := Outer{}
		err := gutil.FillStructWithDefault(&outer)
		t.AssertNil(err)

		t.Assert(outer.O1, 1)
		t.Assert(outer.O2, `1.01`)
		t.Assert(outer.O3, `1.01`)
		t.Assert(outer.O4, true)
		t.Assert(outer.Inner1, nil)
		t.Assert(outer.Inner2.I2V1, `1.01`)
		t.Assert(outer.Inner3.I3V1, 1)
		t.Assert(outer.Inner3.Inner1.I1V1, 0)
		t.Assert(outer.Inner3.Inner1.I1V2, true)
		t.Assert(outer.Inner4, nil)
	})
}

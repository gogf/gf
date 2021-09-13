// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"github.com/gogf/gf/frame/g"
	"testing"

	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gutil"
)

func Test_SliceToMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := g.Slice{
			"K1", "v1", "K2", "v2",
		}
		m := gutil.SliceToMap(s)
		t.Assert(len(m), 2)
		t.Assert(m, g.Map{
			"K1": "v1",
			"K2": "v2",
		})
	})
	gtest.C(t, func(t *gtest.T) {
		s := g.Slice{
			"K1", "v1", "K2",
		}
		m := gutil.SliceToMap(s)
		t.Assert(len(m), 0)
		t.Assert(m, nil)
	})
}

func Test_SliceToMapWithColumnAsKey(t *testing.T) {
	m1 := g.Map{"K1": "v1", "K2": 1}
	m2 := g.Map{"K1": "v2", "K2": 2}
	s := g.Slice{m1, m2}
	gtest.C(t, func(t *gtest.T) {
		m := gutil.SliceToMapWithColumnAsKey(s, "K1")
		t.Assert(m, g.MapAnyAny{
			"v1": m1,
			"v2": m2,
		})
	})
	gtest.C(t, func(t *gtest.T) {
		m := gutil.SliceToMapWithColumnAsKey(s, "K2")
		t.Assert(m, g.MapAnyAny{
			1: m1,
			2: m2,
		})
	})
}

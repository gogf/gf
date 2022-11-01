// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package gstr_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_SearchArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		t.AssertEQ(gstr.SearchArray(a, "a"), 0)
		t.AssertEQ(gstr.SearchArray(a, "b"), 1)
		t.AssertEQ(gstr.SearchArray(a, "c"), 2)
		t.AssertEQ(gstr.SearchArray(a, "d"), -1)
	})
}

func Test_InArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		t.AssertEQ(gstr.InArray(a, "a"), true)
		t.AssertEQ(gstr.InArray(a, "b"), true)
		t.AssertEQ(gstr.InArray(a, "c"), true)
		t.AssertEQ(gstr.InArray(a, "d"), false)
	})
}

func Test_PrefixArray(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a := g.SliceStr{"a", "b", "c"}
		gstr.PrefixArray(a, "1-")
		t.AssertEQ(a, g.SliceStr{"1-a", "1-b", "1-c"})
	})
}

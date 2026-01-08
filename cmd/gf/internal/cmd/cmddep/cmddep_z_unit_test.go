// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmddep

import (
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Dep_Tree(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "tree",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_List(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "list",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Mermaid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "mermaid",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Dot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "dot",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "json",
			Depth:    1,
			Internal: true,
			NoStd:    true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Reverse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "tree",
			Depth:    1,
			Internal: true,
			NoStd:    true,
			Reverse:  true,
		})
		t.AssertNil(err)
	})
}

func Test_Dep_Group(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		_, err := Dep.Index(ctx, Input{
			Package:  "./",
			Format:   "mermaid",
			Depth:    1,
			Internal: true,
			NoStd:    true,
			Group:    true,
		})
		t.AssertNil(err)
	})
}

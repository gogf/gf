// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gctx_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestNew(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		t.AssertNE(ctx, nil)
		t.AssertNE(gctx.CtxId(ctx), "")
	})
}

func TestWithCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		ctx = gctx.WithCtx(ctx)
		t.AssertNE(gctx.CtxId(ctx), "")
		t.Assert(ctx.Value("TEST"), 1)
	})
}

func TestSetInitCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		gctx.SetInitCtx(ctx)
		t.AssertNE(gctx.GetInitCtx(), "")
		t.Assert(gctx.GetInitCtx().Value("TEST"), 1)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gctx_test

import (
	"context"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func Test_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		t.AssertNE(ctx, nil)
		t.AssertNE(gctx.CtxId(ctx), "")
	})
}

func Test_WithCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		ctx = gctx.WithCtx(ctx)
		t.AssertNE(gctx.CtxId(ctx), "")
		t.Assert(ctx.Value("TEST"), 1)
	})
}

func Test_WithPrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		ctx = gctx.WithPrefix(ctx, "H-")
		t.Assert(gstr.Contains(gctx.CtxId(ctx), "H-"), true)
		t.Assert(ctx.Value("TEST"), 1)
	})
}

func Test_WithValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.WithValue(context.TODO(), "TEST", 1)
		ctx = gctx.WithCtxId(ctx, "123")
		t.Assert(gctx.CtxId(ctx), "123")
		t.Assert(ctx.Value("TEST"), 1)
	})
}

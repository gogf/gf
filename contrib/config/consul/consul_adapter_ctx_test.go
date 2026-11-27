// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consul_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"

	consul "github.com/gogf/gf/contrib/config/consul/v2"
)

// TestConsulAdapterCtx_WithGetters tests the With* / Get* accessor round-trips
// for ConsulAdapterCtx. These are pure context-value operations that do not
// require a live Consul server.
func TestConsulAdapterCtx_WithGetters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		content := gjson.New(map[string]any{"key": "val"})
		ctx := consul.NewAdapterCtx().
			WithPath("server/config").
			WithContent(content).
			WithOperation(gcfg.OperationUpdate)

		t.Assert(ctx.GetPath(), "server/config")
		t.Assert(ctx.GetOperation(), gcfg.OperationUpdate)
		t.AssertNE(ctx.GetContent(), nil)
		t.Assert(ctx.GetContent().Get("key").String(), "val")
	})
}

func TestConsulAdapterCtx_NewAdapterCtxWithCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parent := context.Background()
		ctx := consul.NewAdapterCtxWithCtx(parent).WithPath("a/b")
		t.Assert(ctx.GetPath(), "a/b")
		t.Assert(ctx.GetOperation(), gcfg.OperationType(""))
	})
}

func TestConsulAdapterCtx_NilContextFallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := consul.NewAdapterCtxWithCtx(nil)
		t.AssertNE(ctx, nil)
		t.Assert(ctx.GetPath(), "")
	})
}

func TestConsulAdapterCtx_GetAdapterCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		original := consul.NewAdapterCtx().WithPath("x/y")
		got := consul.GetAdapterCtx(original.Ctx)
		t.Assert(got.GetPath(), "x/y")
	})
}

func TestConsulAdapterCtx_GetContentDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// GetContent returns an empty gjson.New(nil) when unset, not a Go nil.
		// Path and operation remain at zero values when nothing was set.
		ctx := consul.NewAdapterCtx()
		t.Assert(ctx.GetPath(), "")
		t.Assert(ctx.GetOperation(), gcfg.OperationType(""))
		t.Assert(ctx.GetContent() == nil, false)
		t.Assert(ctx.GetContent().IsNil(), true)
	})
}

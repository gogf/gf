// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/config/nacos/v2"
)

// TestNacosAdapterCtx_WithGetters tests the With* / Get* accessor round-trips
// for NacosAdapterCtx. These are pure context-value operations that do not
// require a live Nacos server.
func TestNacosAdapterCtx_WithGetters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := nacos.NewAdapterCtx().
			WithNamespace("test-ns").
			WithGroup("test-group").
			WithDataId("test-data-id").
			WithContent("some-content").
			WithOperation(gcfg.OperationUpdate)

		t.Assert(ctx.GetNamespace(), "test-ns")
		t.Assert(ctx.GetGroup(), "test-group")
		t.Assert(ctx.GetDataId(), "test-data-id")
		t.Assert(ctx.GetContent(), "some-content")
		t.Assert(ctx.GetOperation(), gcfg.OperationUpdate)
	})
}

func TestNacosAdapterCtx_NewAdapterCtxWithCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parent := context.Background()
		ctx := nacos.NewAdapterCtxWithCtx(parent).
			WithNamespace("ns").
			WithGroup("g")

		t.Assert(ctx.GetNamespace(), "ns")
		t.Assert(ctx.GetGroup(), "g")
		t.Assert(ctx.GetDataId(), "")                        // unset returns empty string
		t.Assert(ctx.GetContent(), "")                       // unset returns empty string
		t.Assert(ctx.GetOperation(), gcfg.OperationType("")) // unset returns zero
	})
}

func TestNacosAdapterCtx_NilContextFallback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// NewAdapterCtxWithCtx must not panic on nil — it uses context.Background()
		ctx := nacos.NewAdapterCtxWithCtx(nil)
		t.AssertNE(ctx, nil)
		t.Assert(ctx.GetNamespace(), "")
	})
}

func TestNacosAdapterCtx_GetAdapterCtx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		original := nacos.NewAdapterCtx().WithNamespace("ns2").WithGroup("g2")
		// GetAdapterCtx wraps a context — verify that an already-enriched ctx round-trips
		got := nacos.GetAdapterCtx(original.Ctx)
		t.Assert(got.GetNamespace(), "ns2")
		t.Assert(got.GetGroup(), "g2")
	})
}

func TestNacosAdapterCtx_MissingValuesReturnZero(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := nacos.NewAdapterCtx()
		t.Assert(ctx.GetNamespace(), "")
		t.Assert(ctx.GetGroup(), "")
		t.Assert(ctx.GetDataId(), "")
		t.Assert(ctx.GetContent(), "")
		t.Assert(ctx.GetOperation(), gcfg.OperationType(""))
	})
}

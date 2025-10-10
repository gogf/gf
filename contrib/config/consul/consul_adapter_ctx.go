// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package consul implements gcfg.Adapter using consul service.
package consul

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// KeyPath is the context key for path
	KeyPath gctx.StrKey = "path"
)

const (
	// OperationUpdate represents the update operation
	OperationUpdate = "update"
)

// ConsulAdapterCtx is the context adapter for Consul configuration
type ConsulAdapterCtx struct {
	Ctx context.Context
}

// NewAdapterCtxWithCtx creates and returns a new ConsulAdapterCtx with the given context.
func NewAdapterCtxWithCtx(ctx context.Context) *ConsulAdapterCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &ConsulAdapterCtx{Ctx: ctx}
}

// NewAdapterCtx creates and returns a new ConsulAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewAdapterCtx(ctx ...context.Context) *ConsulAdapterCtx {
	if len(ctx) > 0 {
		return NewAdapterCtxWithCtx(ctx[0])
	}
	return NewAdapterCtxWithCtx(context.Background())
}

// GetAdapterCtx creates a new ConsulAdapterCtx with the given context
func GetAdapterCtx(ctx context.Context) *ConsulAdapterCtx {
	return NewAdapterCtxWithCtx(ctx)
}

// WithOperation sets the operation in the context
func (a *ConsulAdapterCtx) WithOperation(operation string) *ConsulAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, gcfg.KeyOperation, operation)
	return a
}

// WithPath sets the path in the context
func (a *ConsulAdapterCtx) WithPath(path string) *ConsulAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyPath, path)
	return a
}

// WithContent sets the content in the context
func (a *ConsulAdapterCtx) WithContent(content *gjson.Json) *ConsulAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, gcfg.KeyContent, content)
	return a
}

// GetContent retrieves the content from the context
func (a *ConsulAdapterCtx) GetContent() *gjson.Json {
	if v := a.Ctx.Value(gcfg.KeyContent); v != nil {
		if s, ok := v.(*gjson.Json); ok {
			return s
		}
	}
	return gjson.New(nil)
}

// GetOperation retrieves the operation from the context
func (a *ConsulAdapterCtx) GetOperation() string {
	if v := a.Ctx.Value(gcfg.KeyOperation); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetPath retrieves the path from the context
func (a *ConsulAdapterCtx) GetPath() string {
	if v := a.Ctx.Value(KeyPath); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

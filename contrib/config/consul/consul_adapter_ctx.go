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
	// ContextKeyPath is the context key for path
	ContextKeyPath gctx.StrKey = "path"
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
func (a *ConsulAdapterCtx) WithOperation(operation gcfg.OperationType) *ConsulAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, gcfg.ContextKeyOperation, operation)
	return a
}

// WithPath sets the path in the context
func (a *ConsulAdapterCtx) WithPath(path string) *ConsulAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, ContextKeyPath, path)
	return a
}

// WithContent sets the content in the context
func (a *ConsulAdapterCtx) WithContent(content *gjson.Json) *ConsulAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, gcfg.ContextKeyContent, content)
	return a
}

// GetContent retrieves the content from the context
func (a *ConsulAdapterCtx) GetContent() *gjson.Json {
	if v := a.Ctx.Value(gcfg.ContextKeyContent); v != nil {
		if s, ok := v.(*gjson.Json); ok {
			return s
		}
	}
	return gjson.New(nil)
}

// GetOperation retrieves the operation from the context
func (a *ConsulAdapterCtx) GetOperation() gcfg.OperationType {
	if v := a.Ctx.Value(gcfg.ContextKeyOperation); v != nil {
		if s, ok := v.(gcfg.OperationType); ok {
			return s
		}
	}
	return ""
}

// GetPath retrieves the path from the context
func (a *ConsulAdapterCtx) GetPath() string {
	if v := a.Ctx.Value(ContextKeyPath); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

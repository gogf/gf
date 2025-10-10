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

// NewConsulAdapterCtx creates and returns a new ConsulAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewConsulAdapterCtx(ctx ...context.Context) *ConsulAdapterCtx {
	if len(ctx) > 0 {
		return &ConsulAdapterCtx{
			Ctx: ctx[0],
		}
	}
	return &ConsulAdapterCtx{Ctx: context.Background()}
}

// GetConsulAdapterCtx creates a new ConsulAdapterCtx with the given context
func GetConsulAdapterCtx(ctx context.Context) *ConsulAdapterCtx {
	return &ConsulAdapterCtx{Ctx: ctx}
}

// WithOperation sets the operation in the context
func (a *ConsulAdapterCtx) WithOperation(operation ...string) *ConsulAdapterCtx {
	if len(operation) > 0 {
		a.Ctx = context.WithValue(a.Ctx, gcfg.KeyOperation, operation[0])
	}
	return a
}

// WithPath sets the path in the context
func (a *ConsulAdapterCtx) WithPath(path ...string) *ConsulAdapterCtx {
	if len(path) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyPath, path[0])
	}
	return a
}

// WithSetContent sets the content in the context
func (a *ConsulAdapterCtx) WithSetContent(content ...*gjson.Json) *ConsulAdapterCtx {
	if len(content) > 0 {
		a.Ctx = context.WithValue(a.Ctx, gcfg.KeySetContent, content[0])
	}
	return a
}

// GetSetContent retrieves the content from the context
func (a *ConsulAdapterCtx) GetSetContent() *gjson.Json {
	if v := a.Ctx.Value(gcfg.KeySetContent); v != nil {
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

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

import (
	"context"
)

// AdapterContentCtx is the context for AdapterContent.
type AdapterContentCtx struct {
	// Ctx is the context with configuration values
	Ctx context.Context
}

// NewAdapterContentCtxWithCtx creates and returns a new AdapterContentCtx with the given context.
func NewAdapterContentCtxWithCtx(ctx context.Context) *AdapterContentCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &AdapterContentCtx{Ctx: ctx}
}

// NewAdapterContentCtx creates and returns a new AdapterContentCtx.
// If ctx is provided, it uses that context, otherwise it creates a background context.
func NewAdapterContentCtx(ctx ...context.Context) *AdapterContentCtx {
	if len(ctx) > 0 {
		return NewAdapterContentCtxWithCtx(ctx[0])
	}
	return NewAdapterContentCtxWithCtx(context.Background())
}

// GetAdapterContentCtx creates and returns an AdapterContentCtx with the given context.
func GetAdapterContentCtx(ctx context.Context) *AdapterContentCtx {
	return NewAdapterContentCtxWithCtx(ctx)
}

// WithOperation sets the operation in the context and returns the updated AdapterContentCtx.
func (a *AdapterContentCtx) WithOperation(operation string) *AdapterContentCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyOperation, operation)
	return a
}

// WithContent sets the content in the context and returns the updated AdapterContentCtx.
func (a *AdapterContentCtx) WithContent(content string) *AdapterContentCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyContent, content)
	return a
}

// GetOperation retrieves the operation from the context.
// Returns empty string if not found.
func (a *AdapterContentCtx) GetOperation() string {
	if v := a.Ctx.Value(KeyOperation); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetContent retrieves the content from the context.
// Returns empty string if not found.
func (a *AdapterContentCtx) GetContent() string {
	if v := a.Ctx.Value(KeyContent); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

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

// NewAdapterContentCtx creates and returns a new AdapterContentCtx.
// If ctx is provided, it uses that context, otherwise it creates a background context.
func NewAdapterContentCtx(ctx ...context.Context) *AdapterContentCtx {
	if len(ctx) > 0 {
		return &AdapterContentCtx{
			Ctx: ctx[0],
		}
	}
	return &AdapterContentCtx{Ctx: context.Background()}
}

// GetAdapterContentCtx creates and returns an AdapterContentCtx with the given context.
func GetAdapterContentCtx(ctx context.Context) *AdapterContentCtx {
	return &AdapterContentCtx{Ctx: ctx}
}

// WithOperation sets the operation in the context and returns the updated AdapterContentCtx.
// If operation is not provided, it does nothing.
func (a *AdapterContentCtx) WithOperation(operation ...string) *AdapterContentCtx {
	if len(operation) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyOperation, operation[0])
	}
	return a
}

// WithSetContent sets the content in the context and returns the updated AdapterContentCtx.
// If content is not provided, it does nothing.
func (a *AdapterContentCtx) WithSetContent(content ...string) *AdapterContentCtx {
	if len(content) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeySetContent, content[0])
	}
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

// GetSetContent retrieves the content from the context.
// Returns empty string if not found.
func (a *AdapterContentCtx) GetSetContent() string {
	if v := a.Ctx.Value(KeySetContent); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

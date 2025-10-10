// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides reading, caching and managing for configuration.
package gcfg

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

// AdapterFileCtx is the context for AdapterFile.
type AdapterFileCtx struct {
	// Ctx is the context with configuration values
	Ctx context.Context
}

// NewAdapterFileCtxWithCtx creates and returns a new AdapterFileCtx with the given context.
func NewAdapterFileCtxWithCtx(ctx context.Context) *AdapterFileCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &AdapterFileCtx{Ctx: ctx}
}

// NewAdapterFileCtx creates and returns a new AdapterFileCtx.
// If ctx is provided, it uses that context, otherwise it creates a background context.
func NewAdapterFileCtx(ctx ...context.Context) *AdapterFileCtx {
	if len(ctx) > 0 {
		return NewAdapterFileCtxWithCtx(ctx[0])
	}
	return NewAdapterFileCtxWithCtx(context.Background())
}

// GetAdapterFileCtx creates and returns an AdapterFileCtx with the given context.
func GetAdapterFileCtx(ctx context.Context) *AdapterFileCtx {
	return NewAdapterFileCtxWithCtx(ctx)
}

// WithFileName sets the file name in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithFileName(fileName string) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyFileName, fileName)
	return a
}

// WithFilePath sets the file path in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithFilePath(filePath string) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyFilePath, filePath)
	return a
}

// WithFileType sets the file type in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithFileType(fileType string) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyFileType, fileType)
	return a
}

// WithOperation sets the operation in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithOperation(operation string) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyOperation, operation)
	return a
}

// WithKey sets the key in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithKey(key string) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyKey, key)
	return a
}

// WithValue sets the value in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithValue(value any) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyValue, value)
	return a
}

// WithContent sets the content in the context and returns the updated AdapterFileCtx.
func (a *AdapterFileCtx) WithContent(content any) *AdapterFileCtx {
	a.Ctx = context.WithValue(a.Ctx, KeyContent, content)
	return a
}

// GetFileName retrieves the file name from the context.
// Returns empty string if not found.
func (a *AdapterFileCtx) GetFileName() string {
	if v := a.Ctx.Value(KeyFileName); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetFilePath retrieves the file path from the context.
// Returns empty string if not found.
func (a *AdapterFileCtx) GetFilePath() string {
	if v := a.Ctx.Value(KeyFilePath); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetFileType retrieves the file type from the context.
// Returns empty string if not found.
func (a *AdapterFileCtx) GetFileType() string {
	if v := a.Ctx.Value(KeyFileType); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetOperation retrieves the operation from the context.
// Returns empty string if not found.
func (a *AdapterFileCtx) GetOperation() string {
	if v := a.Ctx.Value(KeyOperation); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetKey retrieves the key from the context.
// Returns empty string if not found.
func (a *AdapterFileCtx) GetKey() string {
	if v := a.Ctx.Value(KeyKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetValue retrieves the value from the context.
// Returns nil if not found.
func (a *AdapterFileCtx) GetValue() *gvar.Var {
	if v := a.Ctx.Value(KeyValue); v != nil {
		return gvar.New(v)
	}
	return nil
}

// GetContent retrieves the set content from the context.
// Returns nil if not found.
func (a *AdapterFileCtx) GetContent() *gvar.Var {
	if v := a.Ctx.Value(KeyContent); v != nil {
		return gvar.New(v)
	}
	return nil
}

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

// NewAdapterFileCtx creates and returns a new AdapterFileCtx.
// If ctx is provided, it uses that context, otherwise it creates a background context.
func NewAdapterFileCtx(ctx ...context.Context) *AdapterFileCtx {
	if len(ctx) > 0 {
		return &AdapterFileCtx{Ctx: ctx[0]}
	}
	return &AdapterFileCtx{Ctx: context.Background()}
}

// GetAdapterFileCtx creates and returns an AdapterFileCtx with the given context.
func GetAdapterFileCtx(ctx context.Context) *AdapterFileCtx {
	return &AdapterFileCtx{Ctx: ctx}
}

// WithFileName sets the file name in the context and returns the updated AdapterFileCtx.
// If fileName is not provided, it does nothing.
func (a *AdapterFileCtx) WithFileName(fileName ...string) *AdapterFileCtx {
	if len(fileName) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyFileName, fileName[0])
	}
	return a
}

// WithFilePath sets the file path in the context and returns the updated AdapterFileCtx.
// If filePath is not provided, it does nothing.
func (a *AdapterFileCtx) WithFilePath(filePath ...string) *AdapterFileCtx {
	if len(filePath) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyFilePath, filePath[0])
	}
	return a
}

// WithFileType sets the file type in the context and returns the updated AdapterFileCtx.
// If fileType is not provided, it does nothing.
func (a *AdapterFileCtx) WithFileType(fileType ...string) *AdapterFileCtx {
	if len(fileType) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyFileType, fileType[0])
	}
	return a
}

// WithOperation sets the operation in the context and returns the updated AdapterFileCtx.
// If operation is not provided, it does nothing.
func (a *AdapterFileCtx) WithOperation(operation ...string) *AdapterFileCtx {
	if len(operation) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyOperation, operation[0])
	}
	return a
}

// WithSetKey sets the set key in the context and returns the updated AdapterFileCtx.
// If setKey is not provided, it does nothing.
func (a *AdapterFileCtx) WithSetKey(setKey ...string) *AdapterFileCtx {
	if len(setKey) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeySetKey, setKey[0])
	}
	return a
}

// WithSetValue sets the value in the context and returns the updated AdapterFileCtx.
// If value is not provided, it does nothing.
func (a *AdapterFileCtx) WithSetValue(value ...any) *AdapterFileCtx {
	if len(value) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeySetValue, value[0])
	}
	return a
}

// WithSetContent sets the content in the context and returns the updated AdapterFileCtx.
// If content is not provided, it does nothing.
func (a *AdapterFileCtx) WithSetContent(content ...any) *AdapterFileCtx {
	if len(content) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeySetContent, content[0])
	}
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

// GetSetKey retrieves the set key from the context.
// Returns empty string if not found.
func (a *AdapterFileCtx) GetSetKey() string {
	if v := a.Ctx.Value(KeySetKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetSetValue retrieves the set value from the context.
// Returns nil if not found.
func (a *AdapterFileCtx) GetSetValue() *gvar.Var {
	if v := a.Ctx.Value(KeySetValue); v != nil {
		return gvar.New(v)
	}
	return nil
}

// GetSetContent retrieves the set content from the context.
// Returns nil if not found.
func (a *AdapterFileCtx) GetSetContent() *gvar.Var {
	if v := a.Ctx.Value(KeySetContent); v != nil {
		return gvar.New(v)
	}
	return nil
}

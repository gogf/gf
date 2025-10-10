// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package polaris implements gcfg.Adapter using polaris service.
package polaris

import (
	"context"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// KeyNamespace is the context key for namespace
	KeyNamespace gctx.StrKey = "namespace"
	// KeyFileGroup is the context key for group
	KeyFileGroup gctx.StrKey = "fileGroup"
	// KeyFileName is the context key for file name
	KeyFileName gctx.StrKey = "fileName"
)

const (
	// OperationUpdate represents the update operation
	OperationUpdate = "update"
)

// PolarisAdapterCtx is the context adapter for polaris configuration
type PolarisAdapterCtx struct {
	Ctx context.Context
}

// NewAdapterCtxWithCtx creates and returns a new PolarisAdapterCtx with the given context.
func NewAdapterCtxWithCtx(ctx context.Context) *PolarisAdapterCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &PolarisAdapterCtx{Ctx: ctx}
}

// NewAdapterCtx creates and returns a new PolarisAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewAdapterCtx(ctx ...context.Context) *PolarisAdapterCtx {
	if len(ctx) > 0 {
		return NewAdapterCtxWithCtx(ctx[0])
	}
	return NewAdapterCtxWithCtx(context.Background())
}

// GetAdapterCtx creates a new PolarisAdapterCtx with the given context
func GetAdapterCtx(ctx context.Context) *PolarisAdapterCtx {
	return NewAdapterCtxWithCtx(ctx)
}

// WithOperation sets the operation in the context
func (n *PolarisAdapterCtx) WithOperation(operation string) *PolarisAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, gcfg.KeyOperation, operation)
	return n
}

// WithNamespace sets the namespace in the context
func (n *PolarisAdapterCtx) WithNamespace(namespace string) *PolarisAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, KeyNamespace, namespace)
	return n
}

// WithFileGroup sets the group in the context
func (n *PolarisAdapterCtx) WithFileGroup(fileGroup string) *PolarisAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, KeyFileGroup, fileGroup)
	return n
}

// WithFileName sets the fileName in the context
func (n *PolarisAdapterCtx) WithFileName(fileName string) *PolarisAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, KeyFileName, fileName)
	return n
}

// WithContent sets the content in the context
func (n *PolarisAdapterCtx) WithContent(content string) *PolarisAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, gcfg.KeyContent, content)
	return n
}

// GetNamespace retrieves the namespace from the context
func (n *PolarisAdapterCtx) GetNamespace() string {
	if v := n.Ctx.Value(KeyNamespace); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetFileGroup retrieves the group from the context
func (n *PolarisAdapterCtx) GetFileGroup() string {
	if v := n.Ctx.Value(KeyFileGroup); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetFileName retrieves the fileName from the context
func (n *PolarisAdapterCtx) GetFileName() string {
	if v := n.Ctx.Value(KeyFileName); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetContent retrieves the content from the context
func (n *PolarisAdapterCtx) GetContent() string {
	if v := n.Ctx.Value(gcfg.KeyContent); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetOperation retrieves the operation from the context
func (n *PolarisAdapterCtx) GetOperation() string {
	if v := n.Ctx.Value(gcfg.KeyOperation); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

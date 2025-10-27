// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package nacos implements gcfg.Adapter using nacos service.
package nacos

import (
	"context"

	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// ContextKeyNamespace is the context key for namespace
	ContextKeyNamespace gctx.StrKey = "namespace"
	// ContextKeyGroup is the context key for group
	ContextKeyGroup gctx.StrKey = "group"
	// ContextKeyDataId is the context key for dataId
	ContextKeyDataId gctx.StrKey = "dataId"
)

// NacosAdapterCtx is the context adapter for Nacos configuration
type NacosAdapterCtx struct {
	Ctx context.Context
}

// NewAdapterCtxWithCtx creates and returns a new NacosAdapterCtx with the given context.
func NewAdapterCtxWithCtx(ctx context.Context) *NacosAdapterCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &NacosAdapterCtx{Ctx: ctx}
}

// NewAdapterCtx creates and returns a new NacosAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewAdapterCtx(ctx ...context.Context) *NacosAdapterCtx {
	if len(ctx) > 0 {
		return NewAdapterCtxWithCtx(ctx[0])
	}
	return NewAdapterCtxWithCtx(context.Background())
}

// GetAdapterCtx creates a new NacosAdapterCtx with the given context
func GetAdapterCtx(ctx context.Context) *NacosAdapterCtx {
	return NewAdapterCtxWithCtx(ctx)
}

// WithOperation sets the operation in the context
func (n *NacosAdapterCtx) WithOperation(operation gcfg.OperationType) *NacosAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, gcfg.ContextKeyOperation, operation)
	return n
}

// WithNamespace sets the namespace in the context
func (n *NacosAdapterCtx) WithNamespace(namespace string) *NacosAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, ContextKeyNamespace, namespace)
	return n
}

// WithGroup sets the group in the context
func (n *NacosAdapterCtx) WithGroup(group string) *NacosAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, ContextKeyGroup, group)
	return n
}

// WithDataId sets the dataId in the context
func (n *NacosAdapterCtx) WithDataId(dataId string) *NacosAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, ContextKeyDataId, dataId)
	return n
}

// WithContent sets the content in the context
func (n *NacosAdapterCtx) WithContent(content string) *NacosAdapterCtx {
	n.Ctx = context.WithValue(n.Ctx, gcfg.ContextKeyContent, content)
	return n
}

// GetNamespace retrieves the namespace from the context
func (n *NacosAdapterCtx) GetNamespace() string {
	if v := n.Ctx.Value(ContextKeyNamespace); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetGroup retrieves the group from the context
func (n *NacosAdapterCtx) GetGroup() string {
	if v := n.Ctx.Value(ContextKeyGroup); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetDataId retrieves the dataId from the context
func (n *NacosAdapterCtx) GetDataId() string {
	if v := n.Ctx.Value(ContextKeyDataId); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetContent retrieves the content from the context
func (n *NacosAdapterCtx) GetContent() string {
	if v := n.Ctx.Value(gcfg.ContextKeyContent); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetOperation retrieves the operation from the context
func (n *NacosAdapterCtx) GetOperation() gcfg.OperationType {
	if v := n.Ctx.Value(gcfg.ContextKeyOperation); v != nil {
		if s, ok := v.(gcfg.OperationType); ok {
			return s
		}
	}
	return ""
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package kubecm implements gcfg.Adapter using kubecm service.
package kubecm

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// ContextKeyNamespace is the context key for namespace
	ContextKeyNamespace gctx.StrKey = "namespace"
	// ContextKeyConfigMap is the context key for configmap
	ContextKeyConfigMap gctx.StrKey = "configMap"
	// ContextKeyDataItem is the context key for dataitem
	ContextKeyDataItem gctx.StrKey = "dataItem"
)

// KubecmAdapterCtx is the context adapter for kubecm configuration
type KubecmAdapterCtx struct {
	Ctx context.Context
}

// NewKubecmAdapterCtx creates and returns a new KubecmAdapterCtx with the given context.
func NewKubecmAdapterCtx(ctx context.Context) *KubecmAdapterCtx {
	if ctx == nil {
		ctx = context.Background()
	}
	return &KubecmAdapterCtx{Ctx: ctx}
}

// NewAdapterCtx creates and returns a new KubecmAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewAdapterCtx(ctx ...context.Context) *KubecmAdapterCtx {
	if len(ctx) > 0 {
		return NewKubecmAdapterCtx(ctx[0])
	}
	return NewKubecmAdapterCtx(context.Background())
}

// GetAdapterCtx creates a new KubecmAdapterCtx with the given context
func GetAdapterCtx(ctx context.Context) *KubecmAdapterCtx {
	return NewKubecmAdapterCtx(ctx)
}

// WithOperation sets the operation in the context
func (a *KubecmAdapterCtx) WithOperation(operation gcfg.OperationType) *KubecmAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, gcfg.ContextKeyOperation, operation)
	return a
}

// WithNamespace sets the namespace in the context
func (a *KubecmAdapterCtx) WithNamespace(namespace string) *KubecmAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, ContextKeyNamespace, namespace)
	return a
}

// WithConfigMap sets the configmap in the context
func (a *KubecmAdapterCtx) WithConfigMap(configMap string) *KubecmAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, ContextKeyConfigMap, configMap)
	return a
}

// WithDataItem sets the dataitem in the context
func (a *KubecmAdapterCtx) WithDataItem(dataItem string) *KubecmAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, ContextKeyDataItem, dataItem)
	return a
}

// WithContent sets the content in the context
func (a *KubecmAdapterCtx) WithContent(content *gjson.Json) *KubecmAdapterCtx {
	a.Ctx = context.WithValue(a.Ctx, gcfg.ContextKeyContent, content)
	return a
}

// GetOperation retrieves the operation from the context
func (a *KubecmAdapterCtx) GetOperation() gcfg.OperationType {
	if v := a.Ctx.Value(gcfg.ContextKeyOperation); v != nil {
		if s, ok := v.(gcfg.OperationType); ok {
			return s
		}
	}
	return ""
}

// GetNamespace retrieves the namespace from the context
func (a *KubecmAdapterCtx) GetNamespace() string {
	if v := a.Ctx.Value(ContextKeyNamespace); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetConfigMap retrieves the configmap from the context
func (a *KubecmAdapterCtx) GetConfigMap() string {
	if v := a.Ctx.Value(ContextKeyConfigMap); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetDataItem retrieves the dataitem from the context
func (a *KubecmAdapterCtx) GetDataItem() string {
	if v := a.Ctx.Value(ContextKeyDataItem); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetContent retrieves the content from the context
func (a *KubecmAdapterCtx) GetContent() *gjson.Json {
	if v := a.Ctx.Value(gcfg.ContextKeyContent); v != nil {
		if s, ok := v.(*gjson.Json); ok {
			return s
		}
	}
	return gjson.New(nil)
}

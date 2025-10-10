// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package apollo implements gcfg.Adapter using apollo service.
package apollo

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	// KeyNamespace is the context key for namespace
	KeyNamespace gctx.StrKey = "namespace"
	// KeyAppId is the context key for appId
	KeyAppId gctx.StrKey = "appId"
	// KeyCluster is the context key for cluster
	KeyCluster gctx.StrKey = "cluster"
)

const (
	// OperationUpdate represents the update operation
	OperationUpdate = "update"
)

// ApolloAdapterCtx is the context adapter for Apollo configuration
type ApolloAdapterCtx struct {
	Ctx context.Context
}

// NewApolloAdapterCtx creates and returns a new ApolloAdapterCtx.
// If context is provided, it will be used; otherwise, a background context is created.
func NewApolloAdapterCtx(ctx ...context.Context) *ApolloAdapterCtx {
	if len(ctx) > 0 {
		return &ApolloAdapterCtx{
			Ctx: ctx[0],
		}
	}
	return &ApolloAdapterCtx{Ctx: context.Background()}
}

// GetApolloAdapterCtx creates a new ApolloAdapterCtx with the given context
func GetApolloAdapterCtx(ctx context.Context) *ApolloAdapterCtx {
	return &ApolloAdapterCtx{Ctx: ctx}
}

// WithOperation sets the operation in the context
func (a *ApolloAdapterCtx) WithOperation(operation ...string) *ApolloAdapterCtx {
	if len(operation) > 0 {
		a.Ctx = context.WithValue(a.Ctx, gcfg.KeyOperation, operation[0])
	}
	return a
}

// WithNamespace sets the namespace in the context
func (a *ApolloAdapterCtx) WithNamespace(namespace ...string) *ApolloAdapterCtx {
	if len(namespace) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyNamespace, namespace[0])
	}
	return a
}

// WithAppId sets the appId in the context
func (a *ApolloAdapterCtx) WithAppId(appId ...string) *ApolloAdapterCtx {
	if len(appId) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyAppId, appId[0])
	}
	return a
}

// WithCluster sets the cluster in the context
func (a *ApolloAdapterCtx) WithCluster(cluster ...string) *ApolloAdapterCtx {
	if len(cluster) > 0 {
		a.Ctx = context.WithValue(a.Ctx, KeyCluster, cluster[0])
	}
	return a
}

// WithSetContent sets the content in the context
func (a *ApolloAdapterCtx) WithSetContent(content ...*gjson.Json) *ApolloAdapterCtx {
	if len(content) > 0 {
		a.Ctx = context.WithValue(a.Ctx, gcfg.KeySetContent, content[0])
	}
	return a
}

// GetNamespace retrieves the namespace from the context
func (a *ApolloAdapterCtx) GetNamespace() string {
	if v := a.Ctx.Value(KeyNamespace); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetAppId retrieves the appId from the context
func (a *ApolloAdapterCtx) GetAppId() string {
	if v := a.Ctx.Value(KeyAppId); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetCluster retrieves the cluster from the context
func (a *ApolloAdapterCtx) GetCluster() string {
	if v := a.Ctx.Value(KeyCluster); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// GetSetContent retrieves the content from the context
func (a *ApolloAdapterCtx) GetSetContent() *gjson.Json {
	if v := a.Ctx.Value(gcfg.KeySetContent); v != nil {
		if s, ok := v.(*gjson.Json); ok {
			return s
		}
	}
	return gjson.New(nil)
}

// GetOperation retrieves the operation from the context
func (a *ApolloAdapterCtx) GetOperation() string {
	if v := a.Ctx.Value(gcfg.KeyOperation); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

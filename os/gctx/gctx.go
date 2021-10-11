// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gctx wraps context.Context and provides extra context features.
package gctx

import (
	"context"
	"github.com/gogf/gf/util/guid"
)

type (
	Ctx    = context.Context // Ctx is short name alias for context.Context.
	StrKey string            // StrKey is a type for warps basic type string as context key.
)

const (
	// CtxKey is custom tracing context key for context id.
	// The context id a unique string for certain context.
	CtxKey StrKey = "GoFrameCtxId"
)

// New creates and returns a context which contains context id.
func New() context.Context {
	return WithCtx(context.Background())
}

// WithCtx creates and returns a context containing context id upon given parent context `ctx`.
func WithCtx(ctx context.Context) context.Context {
	return WithPrefix(ctx, "")
}

// WithPrefix creates and returns a context containing context id upon given parent context `ctx`.
// The generated context id has custom prefix string specified by parameter `prefix`.
func WithPrefix(ctx context.Context, prefix string) context.Context {
	return WithCtxId(ctx, prefix+getUniqueID())
}

// WithCtxId creates and returns a context containing context id upon given parent context `ctx`.
// The generated context id value is specified by parameter `id`.
func WithCtxId(ctx context.Context, id string) context.Context {
	if id == "" {
		return New()
	}
	return context.WithValue(ctx, CtxKey, id)
}

// CtxId retrieves and returns the context id from context.
func CtxId(ctx context.Context) string {
	s, _ := ctx.Value(CtxKey).(string)
	return s
}

// getUniqueID produces a global unique string.
func getUniqueID() string {
	return guid.S()
}

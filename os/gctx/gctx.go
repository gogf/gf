// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gctx wraps context.Context and provides extra context features.
package gctx

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/net/gtrace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type (
	Ctx    = context.Context // Ctx is short name alias for context.Context.
	StrKey string            // StrKey is a type for warps basic type string as context key.
)

var (
	processCtx context.Context // processCtx is the context initialized from process environment.
	initCtx    context.Context // initCtx is the context for init function of packages.
)

func init() {
	// All environment key-value pairs.
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	// OpenTelemetry from environments.
	processCtx = otel.GetTextMapPropagator().Extract(
		context.Background(),
		propagation.MapCarrier(m),
	)
	// Initialize initialization context.
	initCtx = New()
}

// New creates and returns a context which contains context id.
func New() context.Context {
	return WithCtx(processCtx)
}

// WithCtx creates and returns a context containing context id upon given parent context `ctx`.
func WithCtx(ctx context.Context) context.Context {
	if CtxId(ctx) != "" {
		return ctx
	}
	if gtrace.IsUsingDefaultProvider() {
		var span *gtrace.Span
		ctx, span = gtrace.NewSpan(ctx, "gctx.WithCtx")
		defer span.End()
	}
	return ctx
}

// CtxId retrieves and returns the context id from context.
func CtxId(ctx context.Context) string {
	return gtrace.GetTraceID(ctx)
}

// SetInitCtx sets custom initialization context.
// Note that this function cannot be called in multiple goroutines.
func SetInitCtx(ctx context.Context) {
	initCtx = ctx
}

// GetInitCtx returns the initialization context.
func GetInitCtx() context.Context {
	return initCtx
}

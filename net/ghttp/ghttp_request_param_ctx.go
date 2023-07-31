// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gctx"
)

// RequestFromCtx retrieves and returns the Request object from context.
func RequestFromCtx(ctx context.Context) *Request {
	if v := ctx.Value(ctxKeyForRequest); v != nil {
		return v.(*Request)
	}
	return nil
}

// Context is alias for function GetCtx.
// This function overwrites the http.Request.Context function.
// See GetCtx.
func (r *Request) Context() context.Context {
	var ctx = r.Request.Context()
	// Check and inject Request object into context.
	if RequestFromCtx(ctx) == nil {
		// Inject Request object into context.
		ctx = context.WithValue(ctx, ctxKeyForRequest, r)
		// Add default tracing info if using default tracing provider.
		ctx = gctx.WithCtx(ctx)
		// Update the values of the original HTTP request.
		*r.Request = *r.Request.WithContext(ctx)
	}
	return ctx
}

// GetCtx retrieves and returns the request's context.
// Its alias of function Context,to be relevant with function SetCtx.
func (r *Request) GetCtx() context.Context {
	return r.Context()
}

// GetNeverDoneCtx creates and returns a never done context object,
// which forbids the context manually done, to make the context can be propagated to asynchronous goroutines,
// which will not be affected by the HTTP request ends.
//
// This change is considered for common usage habits of developers for context propagation
// in multiple goroutines creation in one HTTP request.
func (r *Request) GetNeverDoneCtx() context.Context {
	return gctx.NeverDone(r.Context())
}

// SetCtx custom context for current request.
func (r *Request) SetCtx(ctx context.Context) {
	*r.Request = *r.WithContext(ctx)
}

// GetCtxVar retrieves and returns a Var with a given key name.
// The optional parameter `def` specifies the default value of the Var if given `key`
// does not exist in the context.
func (r *Request) GetCtxVar(key interface{}, def ...interface{}) *gvar.Var {
	value := r.Context().Value(key)
	if value == nil && len(def) > 0 {
		value = def[0]
	}
	return gvar.New(value)
}

// SetCtxVar sets custom parameter to context with key-value pairs.
func (r *Request) SetCtxVar(key interface{}, value interface{}) {
	var ctx = r.Context()
	ctx = context.WithValue(ctx, key, value)
	*r.Request = *r.Request.WithContext(ctx)
}

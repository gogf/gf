// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gctx"
)

// neverDoneCtx never done.
type neverDoneCtx struct {
	context.Context
}

// Done forbids the context done from parent context.
func (*neverDoneCtx) Done() <-chan struct{} {
	return nil
}

// Deadline forbids the context deadline from parent context.
func (*neverDoneCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

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
	if r.context == nil {
		// It forbids the context manually done,
		// to make the context can be propagated to asynchronous goroutines.
		r.context = &neverDoneCtx{
			r.Request.Context(),
		}
		r.context = gctx.WithCtx(r.context)
	}
	// Inject Request object into context.
	if RequestFromCtx(r.context) == nil {
		r.context = context.WithValue(r.context, ctxKeyForRequest, r)
	}
	return r.context
}

// GetCtx retrieves and returns the request's context.
func (r *Request) GetCtx() context.Context {
	return r.Context()
}

// SetCtx custom context for current request.
func (r *Request) SetCtx(ctx context.Context) {
	r.context = ctx
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
	r.context = context.WithValue(r.Context(), key, value)
}

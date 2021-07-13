// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
)

// MiddlewareHandlerResponse is the default middleware handling handler response object and its error.
func MiddlewareHandlerResponse(r *Request) {
	r.Middleware.Next()
	res, err := r.GetHandlerResponse()
	if err != nil {
		r.Response.Writef(
			`{"code":%d,"message":"%s"}`,
			gerror.Code(err),
			err.Error(),
		)
		return
	}
	if exception := r.Response.WriteJson(res); exception != nil {
		intlog.Error(r.Context(), exception)
	}
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
)

type DefaultHandlerResponse struct {
	Code    int         `json:"code"    dc:"Error code"`
	Message string      `json:"message" dc:"Error message"`
	Data    interface{} `json:"data"    dc:"Result data for certain request according API definition"`
}

// MiddlewareHandlerResponse is the default middleware handling handler response object and its error.
func MiddlewareHandlerResponse(r *Request) {
	r.Middleware.Next()

	// There's custom buffer content, it then exits current handler.
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err         error
		res         interface{}
		ctx         = r.Context()
		internalErr error
	)
	res, err = r.GetHandlerResponse()
	if err != nil {
		code := gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		internalErr = r.Response.WriteJson(DefaultHandlerResponse{
			Code:    code.Code(),
			Message: err.Error(),
			Data:    nil,
		})
		if internalErr != nil {
			intlog.Errorf(ctx, `%+v`, internalErr)
		}
		return
	}
	internalErr = r.Response.WriteJson(DefaultHandlerResponse{
		Code:    gcode.CodeOK.Code(),
		Message: "",
		Data:    res,
	})
	if internalErr != nil {
		intlog.Errorf(ctx, `%+v`, internalErr)
	}
}

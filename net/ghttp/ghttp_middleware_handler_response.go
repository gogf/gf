// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
)

type DefaultHandlerResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// MiddlewareHandlerResponse is the default middleware handling handler response object and its error.
func MiddlewareHandlerResponse(r *Request) {
	r.Middleware.Next()
	var (
		err         error
		res         interface{}
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
			intlog.Error(r.Context(), internalErr)
		}
		return
	}
	internalErr = r.Response.WriteJson(DefaultHandlerResponse{
		Code:    gcode.CodeOK.Code(),
		Message: "",
		Data:    res,
	})
	if internalErr != nil {
		intlog.Error(r.Context(), internalErr)
	}
}

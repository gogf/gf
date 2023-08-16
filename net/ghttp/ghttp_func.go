// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/httputil"
	"github.com/gogf/gf/v2/text/gstr"
)

// SupportedMethods returns all supported HTTP methods.
func SupportedMethods() []string {
	return gstr.SplitAndTrim(supportedHttpMethods, ",")
}

// BuildParams builds the request string for the http client. The `params` can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter `noUrlEncode` specifies whether to ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	return httputil.BuildParams(params, noUrlEncode...)
}

// niceCallFunc calls function `f` with exception capture logic.
func niceCallFunc(f func()) {
	defer func() {
		if exception := recover(); exception != nil {
			switch exception {
			case exceptionExit, exceptionExitAll:
				return

			default:
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					// It's already an error that has stack info.
					panic(v)
				}
				// Create a new error with stack info.
				// Note that there's a skip pointing the start stacktrace
				// of the real error point.
				if v, ok := exception.(error); ok {
					if gerror.Code(v) != gcode.CodeNil {
						panic(v)
					} else {
						panic(gerror.WrapCodeSkip(
							gcode.CodeInternalPanic, 1, v, "exception recovered",
						))
					}
				} else {
					panic(gerror.NewCodeSkipf(
						gcode.CodeInternalPanic, 1, "exception recovered: %+v", exception,
					))
				}
			}
		}
	}()
	f()
}

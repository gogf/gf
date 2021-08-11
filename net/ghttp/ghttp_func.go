// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp/internal/httputil"
)

// BuildParams builds the request string for the http client. The <params> can be type of:
// string/[]byte/map/struct/*struct.
//
// The optional parameter <noUrlEncode> specifies whether ignore the url encoding for the data.
func BuildParams(params interface{}, noUrlEncode ...bool) (encodedParamStr string) {
	return httputil.BuildParams(params, noUrlEncode...)
}

// niceCallFunc calls function <f> with exception capture logic.
func niceCallFunc(f func()) {
	defer func() {
		if exception := recover(); exception != nil {
			switch exception {
			case exceptionExit, exceptionExitAll:
				return

			default:
				if _, ok := exception.(errorStack); ok {
					// It's already an error that has stack info.
					panic(exception)
				} else {
					// Create a new error with stack info.
					// Note that there's a skip pointing the start stacktrace
					// of the real error point.
					if err, ok := exception.(error); ok {
						if gerror.Code(err) != gerror.CodeNil {
							panic(err)
						} else {
							panic(gerror.WrapCodeSkip(gerror.CodeInternalError, 1, err, ""))
						}
					} else {
						panic(gerror.NewCodeSkipf(gerror.CodeInternalError, 1, "%+v", exception))
					}
				}
			}
		}
	}()
	f()
}

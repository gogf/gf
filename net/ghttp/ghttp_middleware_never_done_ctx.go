// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// MiddlewareNeverDoneCtx sets the context never done for current process.
func MiddlewareNeverDoneCtx(r *Request) {
	r.SetCtx(r.GetNeverDoneCtx())
	r.Middleware.Next()
}

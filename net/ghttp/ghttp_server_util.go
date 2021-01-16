// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "net/http"

// WrapF is a helper function for wrapping http.HandlerFunc and returns a ghttp middleware.
func WrapF(f http.HandlerFunc) HandlerFunc {
	return func(r *Request) {
		f(r.Response.Writer, r.Request)
	}
}

// WrapH is a helper function for wrapping http.Handler and returns a ghttp middleware.
func WrapH(h http.Handler) HandlerFunc {
	return func(r *Request) {
		h.ServeHTTP(r.Response.Writer, r.Request)
	}
}

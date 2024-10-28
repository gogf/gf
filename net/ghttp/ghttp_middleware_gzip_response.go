// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"compress/gzip"
	"strings"
)

// MiddlewareGzip compresses the response content using gzip algorithm.
func MiddlewareGzip(r *Request) {

	r.Middleware.Next()

	var buffer strings.Builder
	gzipwriter := gzip.NewWriter(&buffer)

	_, err := gzipwriter.Write(r.Response.Buffer())
	if err != nil {
		r.Response.Write("gzip write error")
		return
	}

	gzipwriter.Flush()
	gzipwriter.Close()

	r.Response.Header().Set("Content-Encoding", "gzip")
	r.Response.ClearBuffer()

	r.Response.Write(buffer.String())
}

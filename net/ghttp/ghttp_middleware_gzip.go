// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"strings"
)

// MiddlewareGzip is a middleware that compresses HTTP response using gzip compression.
// Note that it does not compress responses if:
// 1. The response is already compressed (Content-Encoding header is set)
// 2. The client does not accept gzip compression
// 3. The response body length is too small (less than 1KB)
//
// To disable compression for specific routes, you can use the group middleware:
//
//	group.Group("/api", func(group *ghttp.RouterGroup) {
//	    group.Middleware(ghttp.MiddlewareGzip) // Enable GZIP for /api routes
//	})
func MiddlewareGzip(r *Request) {
	// Skip compression if client doesn't accept gzip
	if !acceptsGzip(r.Request) {
		r.Middleware.Next()
		return
	}

	// Execute the next handlers first
	r.Middleware.Next()

	// Skip if already compressed or empty response
	if r.Response.Header().Get("Content-Encoding") != "" {
		return
	}

	// Get the response buffer and check its length
	buffer := r.Response.Buffer()
	if len(buffer) < 1024 {
		return
	}

	// Try to compress the response
	var (
		compressed bytes.Buffer
		logger     = r.Server.Logger()
		ctx        = r.Context()
	)
	gzipWriter := gzip.NewWriter(&compressed)
	if _, err := gzipWriter.Write(buffer); err != nil {
		logger.Warningf(ctx, "gzip compression failed: %+v", err)
		return
	}
	if err := gzipWriter.Close(); err != nil {
		logger.Warningf(ctx, "gzip writer close failed: %+v", err)
		return
	}

	// Clear the original buffer and set headers
	r.Response.ClearBuffer()
	r.Response.Header().Set("Content-Encoding", "gzip")
	r.Response.Header().Del("Content-Length")

	// Write the compressed data
	r.Response.Write(compressed.Bytes())
}

// acceptsGzip returns true if the client accepts gzip compression.
func acceptsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

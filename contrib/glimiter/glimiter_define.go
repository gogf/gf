// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter implements rate limiting functionality for HTTP requests.
package glimiter

import (
	"net/http"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	DefaultRate      = 100              // DefaultRate is the default rate of token generation per second
	DefaultShards    = 16               // DefaultShards is the default number of shards for concurrent access
	DefaultCapacity  = 1024             // DefaultCapacity is the default capacity of the token bucket
	DefaultExpire    = 10 * time.Second // DefaultExpire is the default expiration time for cached entries
	DefaultKeyPrefix = "Rate-limiter:"
)

// DefaultKeyFunc is the default function to generate key from request, using client IP
func DefaultKeyFunc(r *ghttp.Request) string {
	return DefaultKeyPrefix + r.GetClientIp()
}

// DefaultAllowHandler is the default handler for allowed requests, continues to next middleware
func DefaultAllowHandler(r *ghttp.Request) {
	r.Middleware.Next()
}

// DefaultDenyHandler is the default handler for denied requests, returns 429 status
func DefaultDenyHandler(r *ghttp.Request) {
	r.Response.WriteStatus(http.StatusTooManyRequests, "Too Many Requests")
	r.ExitAll()
}

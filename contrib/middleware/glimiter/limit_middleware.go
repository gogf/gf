// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter provides rate limiter implementations for GoFrame.
package glimiter

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// MiddlewareConfig defines the configuration for rate limit middleware.
type MiddlewareConfig struct {
	// Limiter is the rate limiter instance to use.
	Limiter Limiter

	// KeyFunc generates the rate limit key from request.
	// Default: use client IP address.
	KeyFunc func(r *ghttp.Request) string

	// ErrorHandler handles rate limit exceeded errors.
	// Default: returns 429 Too Many Requests with JSON response.
	ErrorHandler func(r *ghttp.Request)
}

// Middleware creates and returns a rate limit middleware with the given config.
// It automatically sets standard rate limit headers (X-RateLimit-*) in responses.
//
// Example:
//
//	limiter := glimiter.NewMemoryLimiter(100, time.Minute)
//	s.Group("/api", func(group *ghttp.RouterGroup) {
//	    group.Middleware(glimiter.Middleware(glimiter.MiddlewareConfig{
//	        Limiter: limiter,
//	    }))
//	    group.ALL("/user", handler)
//	})
func Middleware(config MiddlewareConfig) ghttp.HandlerFunc {
	// Set defaults
	if config.KeyFunc == nil {
		config.KeyFunc = func(r *ghttp.Request) string {
			return r.GetClientIp()
		}
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = func(r *ghttp.Request) {
			r.Response.WriteStatus(http.StatusTooManyRequests)
			r.Response.WriteJson(ghttp.DefaultHandlerResponse{
				Code:    http.StatusTooManyRequests,
				Message: "Too Many Requests",
			})
		}
	}

	return func(r *ghttp.Request) {
		ctx := r.Context()
		key := config.KeyFunc(r)

		// Check if request is allowed
		allowed, err := config.Limiter.Allow(ctx, key)
		if err != nil {
			r.SetError(gerror.Wrap(err, "rate limiter error"))
			return
		}

		// Get remaining quota for headers
		remaining, _ := config.Limiter.GetRemaining(ctx, key)
		if remaining < 0 {
			remaining = 0
		}

		// Get accurate reset time from limiter
		resetTime, err := config.Limiter.GetResetTime(ctx, key)
		if err != nil {
			// Fallback to approximate time if error occurs
			resetTime = time.Now().Add(config.Limiter.GetWindow())
		}

		// Set rate limit headers
		r.Response.Header().Set("X-RateLimit-Limit", gconv.String(config.Limiter.GetLimit()))
		r.Response.Header().Set("X-RateLimit-Remaining", gconv.String(remaining))
		r.Response.Header().Set("X-RateLimit-Reset", gconv.String(resetTime.Unix()))

		if !allowed {
			// Rate limit exceeded
			config.ErrorHandler(r)
			return
		}

		// Continue to next handler
		r.Middleware.Next()
	}
}

// MiddlewareByIP creates a middleware that limits by IP address.
// This is a convenience function for the most common use case.
func MiddlewareByIP(limiter Limiter) ghttp.HandlerFunc {
	return Middleware(MiddlewareConfig{
		Limiter: limiter,
		KeyFunc: func(r *ghttp.Request) string {
			return r.GetClientIp()
		},
	})
}

// MiddlewareByAPIKey creates a middleware that limits by API key.
// The keyHeader parameter specifies which header contains the API key.
func MiddlewareByAPIKey(limiter Limiter, keyHeader string) ghttp.HandlerFunc {
	return Middleware(MiddlewareConfig{
		Limiter: limiter,
		KeyFunc: func(r *ghttp.Request) string {
			apiKey := r.Header.Get(keyHeader)
			if apiKey == "" {
				// Fallback to IP if no API key
				return r.GetClientIp()
			}
			return fmt.Sprintf("apikey:%s", apiKey)
		},
	})
}

// MiddlewareByUser creates a middleware that limits by user ID.
// The userFunc parameter extracts the user ID from the request context.
func MiddlewareByUser(limiter Limiter, userFunc func(r *ghttp.Request) string) ghttp.HandlerFunc {
	return Middleware(MiddlewareConfig{
		Limiter: limiter,
		KeyFunc: func(r *ghttp.Request) string {
			userID := userFunc(r)
			if userID == "" {
				// Fallback to IP if no user ID
				return r.GetClientIp()
			}
			return fmt.Sprintf("user:%s", userID)
		},
	})
}

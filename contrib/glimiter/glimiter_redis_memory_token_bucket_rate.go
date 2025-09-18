// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter implements rate limiting functionality for HTTP requests.
package glimiter

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

// RedisMemoryTokenBucketRateLimiterOption defines the configuration options for the Redis+Memory rate limiter
// It contains options for both Redis-based and Memory-based rate limiters
type RedisMemoryTokenBucketRateLimiterOption struct {
	RedisLimiterOption  RedisTokenBucketRateLimiterOption  // Configuration options for the Redis-based rate limiter
	MemoryLimiterOption MemoryTokenBucketRateLimiterOption // Configuration options for the Memory-based rate limiter
}

// RedisMemoryTokenBucketRateLimiter implements a rate limiter that combines Redis and Memory-based token bucket algorithms
// It attempts to use Redis-based rate limiting first, and falls back to Memory-based rate limiting if Redis is unavailable
// This ensures high availability even when Redis is down, although it may cause temporary inconsistency between
// Redis and Memory rate limiters during the fallback period.
type RedisMemoryTokenBucketRateLimiter struct {
	RedisTokenBucketRateLimiter  *RedisTokenBucketRateLimiter  // Redis-based token bucket rate limiter
	MemoryTokenBucketRateLimiter *MemoryTokenBucketRateLimiter // Memory-based token bucket rate limiter for fallback
}

// AllowN checks if n requests are allowed to proceed based on the token bucket algorithm
// It first tries to use the Redis-based rate limiter, and falls back to the Memory-based rate limiter if Redis is unavailable
// Returns true if allowed, false otherwise
func (l *RedisMemoryTokenBucketRateLimiter) AllowN(ctx context.Context, key string, n int64) bool {
	res, err := l.RedisTokenBucketRateLimiter.AllowNWithError(ctx, key, n)
	if err != nil {
		l.RedisTokenBucketRateLimiter.option.Logger.Errorf(ctx, "[Redis Memory Token Bucket Rate limiter] redis eval error: %+v", err)
		mRes := l.MemoryTokenBucketRateLimiter.AllowN(ctx, key, n)
		l.RedisTokenBucketRateLimiter.option.Logger.Debugf(ctx, "[Redis Memory Token Bucket Rate limiter] fallback to memory limiter, result: %v", mRes)
		return mRes
	}
	return res
}

// Middleware returns a middleware function that implements rate limiting using the Redis+Memory token bucket algorithm
// It attempts to use Redis-based rate limiting first, and falls back to Memory-based rate limiting if Redis is unavailable
func (l *RedisMemoryTokenBucketRateLimiter) Middleware() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		ctx := r.GetCtx()
		key := l.RedisTokenBucketRateLimiter.option.KeyFunc(r)
		res, err := l.RedisTokenBucketRateLimiter.AllowNWithError(ctx, key, 1)
		if err != nil {
			l.RedisTokenBucketRateLimiter.option.Logger.Errorf(ctx, "[Redis Memory Token Bucket Rate limiter] redis eval error: %+v", err)
			if !l.MemoryTokenBucketRateLimiter.AllowN(ctx, key, 1) {
				l.RedisTokenBucketRateLimiter.option.DenyHandler(r)
			} else {
				l.MemoryTokenBucketRateLimiter.option.AllowHandler(r)
			}
			return
		}
		if res {
			l.RedisTokenBucketRateLimiter.option.AllowHandler(r)
		} else {
			l.RedisTokenBucketRateLimiter.option.DenyHandler(r)
		}
	}
}

// NewRedisMemoryTokenBucketRateLimiter creates a new Redis+Memory token bucket rate limiter with the given options
// It initializes both Redis-based and Memory-based rate limiters
func NewRedisMemoryTokenBucketRateLimiter(option RedisMemoryTokenBucketRateLimiterOption) *RedisMemoryTokenBucketRateLimiter {
	redisTokenBucketRateLimiter := NewRedisTokenBucketRateLimiter(option.RedisLimiterOption)
	memoryTokenBucketRateLimiter := NewMemoryTokenBucketRateLimiter(option.MemoryLimiterOption)
	return &RedisMemoryTokenBucketRateLimiter{
		RedisTokenBucketRateLimiter:  redisTokenBucketRateLimiter,
		MemoryTokenBucketRateLimiter: memoryTokenBucketRateLimiter,
	}
}

// NewRedisMemoryTokenBucketRateLimiterAndMiddleware creates a new Redis+Memory token bucket rate limiter and returns a middleware function
// It returns both the rate limiter instance and the middleware function
func NewRedisMemoryTokenBucketRateLimiterAndMiddleware(option RedisMemoryTokenBucketRateLimiterOption) (*RedisMemoryTokenBucketRateLimiter, ghttp.HandlerFunc) {
	redisMemoryTokenBucketRateLimiter := NewRedisMemoryTokenBucketRateLimiter(option)
	return redisMemoryTokenBucketRateLimiter, redisMemoryTokenBucketRateLimiter.Middleware()
}

// NewRedisMemoryTokenBucketRateLimiterMiddleware creates a new Redis+Memory token bucket rate limiter and returns a middleware function
// It returns only the middleware function
func NewRedisMemoryTokenBucketRateLimiterMiddleware(option RedisMemoryTokenBucketRateLimiterOption) ghttp.HandlerFunc {
	return NewRedisMemoryTokenBucketRateLimiter(option).Middleware()
}

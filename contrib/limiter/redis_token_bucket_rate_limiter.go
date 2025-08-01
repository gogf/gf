// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package limiter implements rate limiting functionality for HTTP requests.
package limiter

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
)

// RedisLimiterLuaScript is the Lua script used for atomic rate limiting operations in Redis
const RedisLimiterLuaScript = `
local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local cost = tonumber(ARGV[3])
local expire = tonumber(ARGV[4])
local now_milliseconds = tonumber(ARGV[5])
local deny_update = tonumber(ARGV[6])

local bucket = redis.call('HMGET', key, 'tokens', 'last_time')
local last_time = tonumber(bucket[2]) or now_milliseconds
local tokens = tonumber(bucket[1]) or capacity

local delta = now_milliseconds - last_time
local incr = delta * rate / 1000
tokens = math.min(capacity, tokens + incr)
tokens = math.floor(tokens)

last_time = now_milliseconds

if tokens >= cost then
    tokens = tokens - cost
    redis.call('HMSET', key, 'tokens', tokens, 'last_time', last_time)
    redis.call('EXPIRE', key, expire)
    return 1
else
    if deny_update == 1 then
        redis.call('HSET', key, 'last_time', last_time)
        redis.call('EXPIRE', key, expire)
    end
    return 0
end
`

// RedisTokenBucketRateLimiterOption defines the configuration options for the Redis-based rate limiter
type RedisTokenBucketRateLimiterOption struct {
	KeyFunc      func(r *ghttp.Request) string // KeyFunc generates the key used for rate limiting based on the request
	AllowHandler func(r *ghttp.Request)        // AllowHandler is called when a request is allowed to proceed
	DenyHandler  func(r *ghttp.Request)        // DenyHandler is called when a request is denied due to rate limiting
	Logger       *glog.Logger                  // Logger for logging
	Redis        *gredis.Redis                 // Redis client instance
	Capacity     int64                         // Capacity is the maximum number of tokens in the bucket
	Rate         int64                         // Rate is the rate at which tokens are added to the bucket per second
	Expire       time.Duration                 // Expire is the expiration time for cached entries
	DenyUpdate   bool                          // DenyUpdate indicates whether to update the cache when a request is denied
}

// RedisTokenBucketRateLimiter implements a Redis-based token bucket rate limiter
type RedisTokenBucketRateLimiter struct {
	redis  *gredis.Redis
	option RedisTokenBucketRateLimiterOption
}

// AllowNWithError checks if n requests are allowed to proceed based on the token bucket algorithm and returns true if allowed, false otherwise
// It returns an error if there was an issue executing the Redis script
func (l *RedisTokenBucketRateLimiter) AllowNWithError(ctx context.Context, key string, n int64) (bool, error) {
	if n < 0 {
		return false, nil
	}
	if n == 0 {
		return true, nil
	}
	denyUpdate := 0
	if l.option.DenyUpdate {
		denyUpdate = 1
	}
	val, err := l.redis.Eval(ctx, RedisLimiterLuaScript, 1, []string{key}, []any{l.option.Capacity, l.option.Rate, n, int(l.option.Expire.Seconds()), time.Now().UnixMilli(), denyUpdate})
	if err != nil {
		return false, err
	}
	return val.Int() == 1, nil
}

// AllowN checks if n requests are allowed to proceed based on the token bucket algorithm
// Returns true if allowed, false otherwise
func (l *RedisTokenBucketRateLimiter) AllowN(ctx context.Context, key string, n int64) bool {
	res, err := l.AllowNWithError(ctx, key, n)
	if err != nil {
		l.option.Logger.Errorf(ctx, "[Redis Token Bucket Rate limiter] redis eval error: %+v", err)
	}
	return res
}

// Allow checks if a single request is allowed to proceed based on the token bucket algorithm
func (l *RedisTokenBucketRateLimiter) Allow(ctx context.Context, key string) bool {
	return l.AllowN(ctx, key, 1)
}

// NewRedisTokenBucketRateLimiter creates a new Redis token bucket rate limiter with the given options
// It sets default values for any unset options
func NewRedisTokenBucketRateLimiter(option RedisTokenBucketRateLimiterOption) *RedisTokenBucketRateLimiter {
	if option.Redis == nil {
		panic("[Redis Token Bucket Rate limiter] redis client is nil")
	}
	if option.Rate <= 0 {
		option.Rate = DefaultRate
	}
	if option.Capacity <= 0 {
		option.Capacity = DefaultCapacity
	}
	if option.Expire <= 0 {
		option.Expire = DefaultExpire
	}
	if option.KeyFunc == nil {
		option.KeyFunc = DefaultKeyFunc
	}
	if option.AllowHandler == nil {
		option.AllowHandler = DefaultAllowHandler
	}
	if option.DenyHandler == nil {
		option.DenyHandler = DefaultDenyHandler
	}
	return &RedisTokenBucketRateLimiter{
		redis:  option.Redis,
		option: option,
	}
}

// Middleware returns a middleware function that implements rate limiting using the Redis token bucket algorithm
func (l *RedisTokenBucketRateLimiter) Middleware() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		key := l.option.KeyFunc(r)
		if !l.AllowN(r.GetCtx(), key, 1) {
			l.option.DenyHandler(r)
			return
		}
		l.option.AllowHandler(r)
	}
}

// NewRedisTokenBucketRateLimiterAndMiddleware creates a new Redis token bucket rate limiter and returns a middleware function
func NewRedisTokenBucketRateLimiterAndMiddleware(option RedisTokenBucketRateLimiterOption) (*RedisTokenBucketRateLimiter, ghttp.HandlerFunc) {
	limiter := NewRedisTokenBucketRateLimiter(option)
	return limiter, limiter.Middleware()
}

// NewRedisTokenBucketRateLimiterMiddleware returns a middleware function that implements rate limiting
// using the Redis token bucket algorithm
func NewRedisTokenBucketRateLimiterMiddleware(option RedisTokenBucketRateLimiterOption) ghttp.HandlerFunc {
	return NewRedisTokenBucketRateLimiter(option).Middleware()
}

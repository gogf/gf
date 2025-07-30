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
local incr = delta * rate
tokens = math.min(capacity, tokens + incr)

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

// RedisRateLimiterOption defines the configuration options for the Redis-based rate limiter
type RedisRateLimiterOption struct {
	KeyFunc      func(r *ghttp.Request) string // KeyFunc generates the key used for rate limiting based on the request
	AllowHandler func(r *ghttp.Request)        // AllowHandler is called when a request is allowed to proceed
	DenyHandler  func(r *ghttp.Request)        // DenyHandler is called when a request is denied due to rate limiting
	Redis        *gredis.Redis                 // Redis client instance
	Capacity     int64                         // Capacity is the maximum number of tokens in the bucket
	Rate         int64                         // Rate is the rate at which tokens are added to the bucket per second
	Expire       time.Duration                 // Expire is the expiration time for cached entries
	DenyUpdate   bool                          // DenyUpdate indicates whether to update the cache when a request is denied
}

// redisTokenBucketRateLimiter implements a Redis-based token bucket rate limiter
type redisTokenBucketRateLimiter struct {
	redis  *gredis.Redis
	option RedisRateLimiterOption
}

// AllowN checks if n requests are allowed to proceed based on the token bucket algorithm
// Returns true if allowed, false otherwise
func (r *redisTokenBucketRateLimiter) AllowN(ctx context.Context, key string, n int64) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}
	denyUpdate := 0
	if r.option.DenyUpdate {
		denyUpdate = 1
	}
	rate := float64(r.option.Rate) / 1000.0
	val, err := r.redis.Eval(ctx, RedisLimiterLuaScript, 1, []string{key}, []any{r.option.Capacity, rate, n, int(r.option.Expire.Seconds()), time.Now().UnixMilli(), denyUpdate})
	if err != nil {
		glog.Errorf(ctx, "[Redis Token Bucket Rate limiter] eval error: %+v", err)
		return false
	}
	return val.Int() == 1
}

// newRedisTokenBucketRateLimiter creates a new Redis token bucket rate limiter with the given options
// It sets default values for any unset options
func newRedisTokenBucketRateLimiter(option RedisRateLimiterOption) *redisTokenBucketRateLimiter {
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
	return &redisTokenBucketRateLimiter{
		redis:  option.Redis,
		option: option,
	}
}

// RedisTokenBucketRateLimiter returns a middleware function that implements rate limiting
// using the Redis token bucket algorithm
func RedisTokenBucketRateLimiter(option RedisRateLimiterOption) ghttp.HandlerFunc {
	limiter := newRedisTokenBucketRateLimiter(option)
	return func(r *ghttp.Request) {
		key := limiter.option.KeyFunc(r)
		if limiter.AllowN(r.Context(), key, 1) {
			limiter.option.AllowHandler(r)
		} else {
			limiter.option.DenyHandler(r)
		}
	}
}

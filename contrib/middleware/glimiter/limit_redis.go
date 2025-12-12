// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter provides rate limiter implementations for GoFrame.
package glimiter

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/util/gconv"
)

// RedisLimiter implements Limiter using Redis.
// It uses sliding window algorithm with Lua script for atomic operations.
type RedisLimiter struct {
	redis  *gredis.Redis
	limit  int
	window time.Duration
}

// Lua script for sliding window rate limiting
// Uses atomic check-then-act with proper locking via Lua script execution
const luaAllowScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local n = tonumber(ARGV[4])

-- Remove expired entries first (cleanup old data)
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- Get current count AFTER cleanup
local current = redis.call('ZCARD', key)

-- Check if adding n would exceed limit
if current + n > limit then
    return 0  -- Deny
end

local counterKey = key .. ':counter'

-- Add entries using timestamp with unique member to ensure each request is counted
for i = 1, n do
    -- Unique member using concatenation to avoid Lua number overflow
    -- Format: "timestamp-index-counter" ensures uniqueness
	local counter = redis.call('INCR', counterKey)
    local member = tostring(now) .. '-' .. tostring(i) .. '-' .. tostring(counter)
    redis.call('ZADD', key, now, member)
end

-- Set expiration (add buffer to avoid premature expiration)
redis.call('PEXPIRE', key, window + 1000)
redis.call('PEXPIRE', counterKey, window + 1000)


return 1
`

// Lua script for getting remaining quota
const luaRemainingScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local now = tonumber(ARGV[3])

-- Remove expired entries
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- Get current count
local current = redis.call('ZCARD', key)

return limit - current
`

// NewRedisLimiter creates and returns a new Redis-based rate limiter.
// The limit parameter specifies the maximum number of requests allowed,
// and window specifies the time window duration.
func NewRedisLimiter(redis *gredis.Redis, limit int, window time.Duration) *RedisLimiter {
	return &RedisLimiter{
		redis:  redis,
		limit:  limit,
		window: window,
	}
}

// Allow implements Limiter.Allow.
func (l *RedisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	return l.AllowN(ctx, key, 1)
}

// AllowN implements Limiter.AllowN.
// It uses Lua script to ensure atomic operations in Redis.
// The Lua script guarantees that check-and-add is atomic, preventing race conditions.
func (l *RedisLimiter) AllowN(ctx context.Context, key string, n int) (bool, error) {
	if n <= 0 {
		return false, fmt.Errorf("n must be positive, got %d", n)
	}

	now := time.Now().UnixMilli()
	windowMs := l.window.Milliseconds()

	result, err := l.redis.Eval(ctx, luaAllowScript, 1, []string{key}, []any{l.limit, windowMs, now, n})
	if err != nil {
		return false, err
	}

	return result.Int() == 1, nil
}

// Wait implements Limiter.Wait.
// It blocks until a request is allowed or context is cancelled.
func (l *RedisLimiter) Wait(ctx context.Context, key string) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Check without consuming quota
			remaining, err := l.GetRemaining(ctx, key)
			if err != nil {
				return err
			}
			if remaining > 0 {
				// Now try to allow
				allowed, err := l.Allow(ctx, key)
				if err != nil {
					return err
				}
				if allowed {
					return nil
				}
			}
		}
	}
}

// GetLimit implements Limiter.GetLimit.
func (l *RedisLimiter) GetLimit() int {
	return l.limit
}

// GetWindow implements Limiter.GetWindow.
func (l *RedisLimiter) GetWindow() time.Duration {
	return l.window
}

// GetRemaining implements Limiter.GetRemaining.
func (l *RedisLimiter) GetRemaining(ctx context.Context, key string) (int, error) {
	now := time.Now().UnixMilli()
	windowMs := l.window.Milliseconds()

	result, err := l.redis.Eval(ctx, luaRemainingScript, 1, []string{key}, []any{l.limit, windowMs, now})
	if err != nil {
		return 0, err
	}

	return gconv.Int(result), nil
}

// Reset implements Limiter.Reset.
func (l *RedisLimiter) Reset(ctx context.Context, key string) error {
	_, err := l.redis.Del(ctx, key)
	return err
}

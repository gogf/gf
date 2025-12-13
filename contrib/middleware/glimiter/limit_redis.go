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

-- If main key was empty (first request or after expiration), reset counter to avoid orphaned state
if current == 0 then
    redis.call('DEL', counterKey)
end

-- Add entries using timestamp with unique member to ensure each request is counted
for i = 1, n do
    -- Unique member using concatenation to avoid Lua number overflow
    -- Format: "timestamp-index-counter" ensures uniqueness
	local counter = redis.call('INCR', counterKey)
    local member = tostring(now) .. '-' .. tostring(i) .. '-' .. tostring(counter)
    redis.call('ZADD', key, now, member)
end

-- Set expiration with buffer to avoid premature expiration
-- Both keys get the same TTL to ensure consistency
-- The buffer prevents edge cases where key expires during active requests
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

// Lua script for getting reset time
// Returns the timestamp of the oldest entry in the window
const luaResetTimeScript = `
local key = KEYS[1]
local window = tonumber(ARGV[1])
local now = tonumber(ARGV[2])

-- Remove expired entries
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)

-- Get the oldest entry's score (timestamp)
local oldest = redis.call('ZRANGE', key, 0, 0, 'WITHSCORES')

if #oldest == 0 then
    -- No entries, return current time
    return now
end

-- oldest[2] contains the score (timestamp) of the oldest entry
-- Reset time is when this oldest entry expires
return tonumber(oldest[2]) + window
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
			// Try to allow directly without checking remaining first
			// This reduces Redis round-trips from two to one per attempt
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
	if result.IsNil() {
		return 0, nil
	}
	return result.Int(), nil
}

// GetResetTime implements Limiter.GetResetTime.
// Returns the time when the oldest request in the sliding window will expire.
func (l *RedisLimiter) GetResetTime(ctx context.Context, key string) (time.Time, error) {
	now := time.Now().UnixMilli()
	windowMs := l.window.Milliseconds()

	result, err := l.redis.Eval(ctx, luaResetTimeScript, 1, []string{key}, []any{windowMs, now})
	if err != nil {
		return time.Time{}, err
	}

	if result.IsNil() {
		return time.Now(), nil
	}

	resetTimeMs := result.Int64()
	return time.UnixMilli(resetTimeMs), nil
}

// Reset implements Limiter.Reset.
// Deletes both the main key and counter key to ensure clean state.
func (l *RedisLimiter) Reset(ctx context.Context, key string) error {
	// Delete both keys atomically
	_, err := l.redis.Del(ctx, key, fmt.Sprintf("%s:counter", key))
	return err
}

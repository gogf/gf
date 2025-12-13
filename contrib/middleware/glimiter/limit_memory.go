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

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/os/gcache"
)

// MemoryLimiter implements Limiter using in-memory cache.
// It uses gcache for storage with automatic expiration.
type MemoryLimiter struct {
	cache  *gcache.Cache
	limit  int
	window time.Duration
}

// memoryCounter holds the counter value for a rate limit key.
// Uses gtype.Int64 for better encapsulation and consistency with GoFrame style.
type memoryCounter struct {
	count *gtype.Int64
}

// NewMemoryLimiter creates and returns a new memory-based rate limiter.
// The limit parameter specifies the maximum number of requests allowed,
// and window specifies the time window duration.
func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	return &MemoryLimiter{
		cache:  gcache.New(),
		limit:  limit,
		window: window,
	}
}

// Allow implements Limiter.Allow.
func (l *MemoryLimiter) Allow(ctx context.Context, key string) (bool, error) {
	return l.AllowN(ctx, key, 1)
}

// AllowN implements Limiter.AllowN.
// It uses CAS (Compare-And-Swap) to ensure atomic check-and-increment operations.
func (l *MemoryLimiter) AllowN(ctx context.Context, key string, n int) (bool, error) {
	if n <= 0 {
		return false, fmt.Errorf("n must be positive, got %d", n)
	}

	// Limit retries to prevent infinite loop in extreme concurrent scenarios
	const maxOuterRetries = 100
	for outerRetry := 0; outerRetry < maxOuterRetries; outerRetry++ {
		// Get or create counter
		value, err := l.cache.GetOrSetFuncLock(ctx, key, func(ctx context.Context) (any, error) {
			return &memoryCounter{count: gtype.NewInt64(0)}, nil
		}, l.window)
		if err != nil {
			return false, err
		}

		counter, ok := value.Val().(*memoryCounter)
		if !ok {
			return false, fmt.Errorf("invalid counter type")
		}

		// Try to increment using CAS (via gtype.Int64)
		// Inner loop handles concurrent updates to the same counter
		const maxCASRetries = 1000
		for casRetry := 0; casRetry < maxCASRetries; casRetry++ {
			current := counter.count.Val()
			if current+int64(n) > int64(l.limit) {
				// Would exceed limit
				return false, nil
			}

			// Try to update using Cas method
			if counter.count.Cas(current, current+int64(n)) {
				return true, nil
			}
			// CAS failed, retry
		}
		// If CAS retries exhausted, it likely means the counter pointer is stale
		// (cache entry was recreated), retry outer loop to get fresh counter
	}

	// Should never reach here under normal circumstances
	return false, fmt.Errorf("exceeded maximum retry attempts")
}

// Wait implements Limiter.Wait.
// It blocks until a request is allowed or context is cancelled.
func (l *MemoryLimiter) Wait(ctx context.Context, key string) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Try to allow directly without checking remaining first
			// This reduces overhead from two operations to one
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
func (l *MemoryLimiter) GetLimit() int {
	return l.limit
}

// GetWindow implements Limiter.GetWindow.
func (l *MemoryLimiter) GetWindow() time.Duration {
	return l.window
}

// GetRemaining implements Limiter.GetRemaining.
func (l *MemoryLimiter) GetRemaining(ctx context.Context, key string) (int, error) {
	value, err := l.cache.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	if value == nil || value.IsNil() {
		return l.limit, nil
	}

	counter, ok := value.Val().(*memoryCounter)
	if !ok {
		return 0, fmt.Errorf("invalid counter type")
	}

	current := counter.count.Val()
	remaining := l.limit - int(current)
	if remaining < 0 {
		remaining = 0
	}
	return remaining, nil
}

// GetResetTime implements Limiter.GetResetTime.
// Returns the expiration time of the cache entry, which is when the rate limit will reset.
func (l *MemoryLimiter) GetResetTime(ctx context.Context, key string) (time.Time, error) {
	expireTime, err := l.cache.GetExpire(ctx, key)
	if err != nil {
		return time.Time{}, err
	}

	// If key doesn't exist or has no expiration, return current time
	if expireTime == 0 {
		return time.Now(), nil
	}

	return time.Now().Add(expireTime), nil
}

// Reset implements Limiter.Reset.
func (l *MemoryLimiter) Reset(ctx context.Context, key string) error {
	_, err := l.cache.Remove(ctx, key)
	return err
}

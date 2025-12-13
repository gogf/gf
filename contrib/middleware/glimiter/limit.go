// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter provides rate limiter implementations for GoFrame.
// It supports both in-memory and Redis-based rate limiting with sliding window algorithm and concurrent-safe operations.
package glimiter

import (
	"context"
	"time"
)

// Limiter defines the interface for rate limiting operations.
type Limiter interface {
	// Allow checks if a single request is allowed under the rate limit.
	// It returns true if allowed, false otherwise.
	Allow(ctx context.Context, key string) (bool, error)

	// AllowN checks if N requests are allowed under the rate limit.
	// It returns true if allowed, false otherwise.
	AllowN(ctx context.Context, key string, n int) (bool, error)

	// Wait blocks until a single request is allowed.
	// It returns error if context is cancelled.
	Wait(ctx context.Context, key string) error

	// GetLimit returns the maximum number of requests allowed in the time window.
	GetLimit() int

	// GetWindow returns the time window duration for rate limiting.
	GetWindow() time.Duration

	// GetRemaining returns the remaining quota for the given key.
	// It returns -1 if the key doesn't exist.
	GetRemaining(ctx context.Context, key string) (int, error)

	// GetResetTime returns the time when the rate limit will reset for the given key.
	// It returns the exact timestamp when the oldest request in the window will expire.
	// For keys that don't exist, it returns the current time.
	GetResetTime(ctx context.Context, key string) (time.Time, error)

	// Reset resets the rate limit for the given key.
	Reset(ctx context.Context, key string) error
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter provides a rate limiter implementation
package glimiter

import "context"

// Adapter is the core adapter for limit features implements.
type Adaper interface {

	// All interfaces add context , you can add special parameters through the content value,
	// and make judgments through RefFilterInfo in Adapter, such as ip, user_id, etc

	// Acquires permits from this RateLimiter, blocking until the request can be granted.
	Acquire(ctx context.Context, reqCount ...int64) bool

	// ResetStatus reset limiter status
	ResetStatus(ctx context.Context)

	// TryAcuqire Acquires permits from this RateLimiter if it can be acquired immediately without delay.
	TryAcuqire(ctx context.Context, reqCount ...int64) bool
}

type localAdapter = Adaper

// Limiter is an adapter implements using local
type Limiter struct {
	localAdapter
}

// NewWithAdapter Create a new limiter with given adapter
func NewWithAdapter(adapter Adaper) *Limiter {
	return &Limiter{adapter}
}

// New Create a new limiter with given rate and burst .
// default adapter is TokenBucketAdapter(use golang.org/x/time/rate)
func New(rate float64, burstCount int64) *Limiter {
	tba := &TokenBucketAdapter{rate, burstCount}
	return &Limiter{tba}
}

// IsAcquire .
func (l *Limiter) IsAcquire(ctx context.Context, reqCount ...int64) bool {
	return l.Acquire(ctx, reqCount...)
}

// Reset .reset limiter status
func (l *Limiter) Reset(ctx context.Context) {
	l.ResetStatus(ctx)
}

// TryReqAcquire Acquires permits from this RateLimiter if it can be acquired immediately without delay.
func (l *Limiter) TryReqAcquire(ctx context.Context, reqCount ...int64) bool {
	return l.TryAcuqire(ctx, reqCount...)
}

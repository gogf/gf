// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter provides a rate limiter implementation
package glimiter

import "context"

type Adaper interface {

	// All interfaces add context , you can add special parameters through the content value,
	// and make judgments through RefFilterInfo in Adapter, such as ip, user_id, etc

	// Acquires permits from this RateLimiter, blocking until the request can be granted.
	Acquire(ctx context.Context, reqCount ...int64) bool

	// Reset limiter status
	ResetStatus(ctx context.Context)

	//Acquires permits from this RateLimiter if it can be acquired immediately without delay.
	TryAcuqire(ctx context.Context, reqCount ...int64) bool
}

type localAdapter = Adaper

type Limiter struct {
	localAdapter
}

// Crate a new limiter with given adapter
func NewWithAdapter(adapter Adaper) *Limiter {
	return &Limiter{adapter}
}

// Crate a new limiter with given rate and burst .
// default adapter is TokenBucketAdapter(use golang.org/x/time/rate)
func New(rate float64, burstCount int64) *Limiter {
	tba := &TokenBucketAdapter{rate, burstCount}
	return &Limiter{tba}
}

func (l *Limiter) IsAcquire(ctx context.Context, reqCount ...int64) bool {
	return l.Acquire(ctx, reqCount...)
}

func (l *Limiter) Reset(ctx context.Context) {
	l.ResetStatus(ctx)
}

func (l *Limiter) TryReqAcquire(ctx context.Context, reqCount ...int64) bool {
	return l.TryAcuqire(ctx, reqCount...)
}

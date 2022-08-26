package glimiter

import "context"

// TokenBucket Adapter provider , Reference: golang.org/x/time/rate
type TokenBucketAdapter struct {
	Limit  float64 // allows events up to rate
	Bursts int64   //permits bursts of at most b tokens.
}

// Acquire reports whether n events may happen at time now.
// Use this method if you intend to drop / skip events that exceed the rate limit.
// Otherwise use Reserve or Wait.
func (t *TokenBucketAdapter) Acquire(ctx context.Context, reqCount ...int64) bool {
	return false
}

// Reset Status of limiter
func (t *TokenBucketAdapter) ResetStatus(ctx context.Context) {

}

// TryAcuqire: Acquires  permits from this RateLimiter if it can be acquired immediately without delay.
func (t *TokenBucketAdapter) TryAcuqire(ctx context.Context, reqCount ...int64) bool {

	return false
}

package glimiter

import "context"

// TokenBucket Adapter provider  , Reference:   golang.org/x/time/rate
type TokenBucketAdapter struct {
	Limit float64 // allows events up to rate
	Bursts  int64 //permits bursts of at most b tokens. 
}


func (t *TokenBucketAdapter) Acquire(ctx context.Context, reqCount ... int64) bool {

	return false
}


func (t *TokenBucketAdapter) ResetStatus(ctx context.Context) {
	
}

func (t *TokenBucketAdapter) TryAcuqire(ctx context.Context, reqCount ...  int64) bool {

	return false
}



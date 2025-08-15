// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glimiter implements rate limiting functionality for HTTP requests.
package glimiter

import (
	"context"
	"hash/fnv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gconv"
)

// TokenBucket represents the token bucket structure
type TokenBucket struct {
	Tokens   int64     // Tokens represents the number of tokens in the bucket
	LastTime time.Time // LastTime represents the last time the bucket was updated
}

// MemoryTokenBucketRateLimiterOption defines the configuration options for the rate limiter
type MemoryTokenBucketRateLimiterOption struct {
	KeyFunc      func(r *ghttp.Request) string // KeyFunc generates the key used for rate limiting based on the request
	AllowHandler func(r *ghttp.Request)        // AllowHandler is called when a request is allowed to proceed
	DenyHandler  func(r *ghttp.Request)        // DenyHandler is called when a request is denied due to rate limiting
	Logger       *glog.Logger                  // Logger for logging
	Shards       int                           // Shards is the number of shards for concurrent access optimization
	LruCapacity  int                           // LruCapacity is the LRU cache capacity
	Capacity     int64                         // Capacity is the maximum number of tokens in the bucket
	Rate         int64                         // Rate is the rate at which tokens are added to the bucket per second
	Expire       time.Duration                 // Expire is the expiration time for cached entries
	DenyUpdate   bool                          // DenyUpdate indicates whether to update the cache when a request is denied
}

// MemoryTokenBucketRateLimiter implements a thread-safe token bucket rate limiter using memory storage
type MemoryTokenBucketRateLimiter struct {
	cache   *gcache.Cache
	option  MemoryTokenBucketRateLimiterOption
	mutexes []sync.Mutex
	shards  int
}

// getShards calculates which shard a key belongs to using FNV-1a hash
func (m *MemoryTokenBucketRateLimiter) getShards(ctx context.Context, key string) int {
	var hash uint64
	h := fnv.New64a()
	_, err := h.Write([]byte(key))
	if err != nil {
		m.option.Logger.Errorf(ctx, "[Token Bucket Rate limiter] hash [%s]error: %+v", key, err)
		hash = 0
	} else {
		hash = h.Sum64()
	}
	return int(hash % uint64(m.shards))
}

// AllowN checks if n tokens can be consumed, and consumes them if possible
func (m *MemoryTokenBucketRateLimiter) AllowN(ctx context.Context, key string, n int64) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}
	shard := m.getShards(ctx, key)
	m.mutexes[shard].Lock()
	defer m.mutexes[shard].Unlock()
	val, err := m.cache.Get(ctx, key)
	if err != nil {
		m.option.Logger.Errorf(ctx, "[Token Bucket Rate limiter] cache get [%s] error: %+v", key, err)
		return false
	}
	var (
		tokens   int64
		lastTime time.Time
	)
	if !val.IsNil() {
		data := val.Val().(*TokenBucket)
		tokens = data.Tokens
		lastTime = data.LastTime
	} else {
		tokens = m.option.Capacity
		lastTime = time.Now()
	}
	delta := time.Since(lastTime)
	if delta < 0 {
		m.option.Logger.Errorf(ctx, "[Memory Token Bucket Rate limiter] delta: [%s] < 0", delta)
		delta = 0
	}
	incr := delta.Nanoseconds() * m.option.Rate / 1e9
	tokens += incr
	if tokens > m.option.Capacity {
		tokens = m.option.Capacity
	}
	if tokens >= n {
		tokens -= n
		bucket := &TokenBucket{
			Tokens:   tokens,
			LastTime: time.Now(),
		}
		err = m.cache.Set(ctx, key, bucket, m.option.Expire)
		if err != nil {
			m.option.Logger.Errorf(ctx, "[Memory Token Bucket Rate limiter] cache set [%s]: [%s]error: %+v", key, gconv.String(bucket), err)
			return false
		}
		return true
	}
	if m.option.DenyUpdate {
		bucket := &TokenBucket{
			Tokens:   tokens,
			LastTime: time.Now(),
		}
		err = m.cache.Set(ctx, key, bucket, m.option.Expire)
		if err != nil {
			m.option.Logger.Errorf(ctx, "[Memory Token Bucket Rate limiter] cache set [%s]: [%s]error: %+v", key, gconv.String(bucket), err)
			return false
		}
	}
	return false
}

// Allow checks if a single token can be consumed
func (m *MemoryTokenBucketRateLimiter) Allow(ctx context.Context, key string) bool {
	return m.AllowN(ctx, key, 1)
}

// Middleware returns one HTTP middleware function that implements rate limiting
func (m *MemoryTokenBucketRateLimiter) Middleware() ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		key := m.option.KeyFunc(r)
		if !m.AllowN(r.GetCtx(), key, 1) {
			m.option.DenyHandler(r)
			return
		}
		m.option.AllowHandler(r)
	}
}

// NewMemoryTokenBucketRateLimiter creates a new memory-based token bucket rate limiter
// It sets default values for any unset options
func NewMemoryTokenBucketRateLimiter(option MemoryTokenBucketRateLimiterOption) *MemoryTokenBucketRateLimiter {
	shards := 16
	if option.Shards <= 0 {
		shards = DefaultShards
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
	if option.Logger == nil {
		option.Logger = glog.New()
	}
	var cache *gcache.Cache
	if option.LruCapacity > 0 {
		cache = gcache.New(option.LruCapacity)
	} else {
		cache = gcache.New()
	}
	return &MemoryTokenBucketRateLimiter{
		cache:   cache,
		option:  option,
		mutexes: make([]sync.Mutex, shards),
		shards:  shards,
	}
}

// NewMemoryTokenBucketRateLimiterAndMiddleware creates a new memory-based token bucket rate limiter and returns a middleware function
func NewMemoryTokenBucketRateLimiterAndMiddleware(option MemoryTokenBucketRateLimiterOption) (*MemoryTokenBucketRateLimiter, ghttp.HandlerFunc) {
	limiter := NewMemoryTokenBucketRateLimiter(option)
	return limiter, limiter.Middleware()
}

// NewMemoryTokenBucketRateLimiterMiddleware creates a new memory-based token bucket rate limiter and returns a middleware function
func NewMemoryTokenBucketRateLimiterMiddleware(option MemoryTokenBucketRateLimiterOption) ghttp.HandlerFunc {
	return NewMemoryTokenBucketRateLimiter(option).Middleware()
}

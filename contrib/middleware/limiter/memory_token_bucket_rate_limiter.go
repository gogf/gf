// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package limiter implements rate limiting functionality for HTTP requests.
package limiter

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

// MemoryRateLimiterOption defines the configuration options for the rate limiter
type MemoryRateLimiterOption struct {
	KeyFunc      func(r *ghttp.Request) string // KeyFunc generates the key used for rate limiting based on the request
	AllowHandler func(r *ghttp.Request)        // AllowHandler is called when a request is allowed to proceed
	DenyHandler  func(r *ghttp.Request)        // DenyHandler is called when a request is denied due to rate limiting
	Shards       int                           // Shards is the number of shards for concurrent access optimization
	LruCapacity  int                           // LruCapacity is the LRU cache capacity
	Capacity     int64                         // Capacity is the maximum number of tokens in the bucket
	Rate         int64                         // Rate is the rate at which tokens are added to the bucket per second
	Expire       time.Duration                 // Expire is the expiration time for cached entries
	DenyUpdate   bool                          // DenyUpdate indicates whether to update the cache when a request is denied
}

// memoryTokenBucketRateLimiter implements a thread-safe token bucket rate limiter using memory storage
type memoryTokenBucketRateLimiter struct {
	cache   *gcache.Cache
	option  MemoryRateLimiterOption
	mutexes []sync.Mutex
	shards  int
}

// getShards calculates which shard a key belongs to using FNV-1a hash
func (t *memoryTokenBucketRateLimiter) getShards(ctx context.Context, key string) int {
	var hash uint64
	h := fnv.New64a()
	_, err := h.Write([]byte(key))
	if err != nil {
		glog.Errorf(ctx, "[Token Bucket Rate limiter] hash [%s]error: %+v", key, err)
		hash = 0
	} else {
		hash = h.Sum64()
	}
	return int(hash % uint64(t.shards))
}

// AllowN checks if n tokens can be consumed, and consumes them if possible
func (t *memoryTokenBucketRateLimiter) AllowN(ctx context.Context, key string, n int64) bool {
	if n < 0 {
		return false
	}
	if n == 0 {
		return true
	}
	shard := t.getShards(ctx, key)
	t.mutexes[shard].Lock()
	defer t.mutexes[shard].Unlock()
	val, err := t.cache.Get(ctx, key)
	if err != nil {
		glog.Errorf(ctx, "[Token Bucket Rate limiter] cache get [%s] error: %+v", key, err)
		return false
	}
	var (
		tokens   int64
		lastTime time.Time
	)
	if val != nil {
		data := val.Map()
		tokens = data[Tokens].(int64)
		lastTime = data[LastTime].(time.Time)
	} else {
		tokens = t.option.Capacity
		lastTime = time.Now()
	}
	delta := time.Since(lastTime)
	if delta < 0 {
		delta = 0
	}
	incr := delta.Nanoseconds() * t.option.Rate / 1e9
	tokens += incr
	if tokens > t.option.Capacity {
		tokens = t.option.Capacity
	}
	if tokens >= n {
		tokens -= n
		bucket := map[string]any{
			Tokens:   tokens,
			LastTime: time.Now(),
		}
		err = t.cache.Set(ctx, key, bucket, t.option.Expire)
		if err != nil {
			glog.Errorf(ctx, "[Token Bucket Rate limiter] cache set [%s]: [%s]error: %+v", key, gconv.String(bucket), err)
			return false
		}
		return true
	}
	if t.option.DenyUpdate {
		bucket := map[string]any{
			Tokens:   tokens,
			LastTime: time.Now(),
		}
		err = t.cache.Set(ctx, key, bucket, t.option.Expire)
		if err != nil {
			glog.Errorf(ctx, "[Token Bucket Rate limiter] cache set [%s]: [%s]error: %+v", key, gconv.String(bucket), err)
			return false
		}
	}
	return false
}

// newMemoryTokenBucketRateLimiter creates a new memory-based token bucket rate limiter
func newMemoryTokenBucketRateLimiter(option MemoryRateLimiterOption) *memoryTokenBucketRateLimiter {
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
	var cache *gcache.Cache
	if option.LruCapacity > 0 {
		cache = gcache.New(option.LruCapacity)
	} else {
		cache = gcache.New()
	}
	return &memoryTokenBucketRateLimiter{
		cache:   cache,
		option:  option,
		mutexes: make([]sync.Mutex, shards),
		shards:  shards,
	}
}

// MemoryTokenBucketRateLimiter returns one HTTP middleware function that implements rate limiting
// using the token bucket algorithm with in-memory storage
func MemoryTokenBucketRateLimiter(option MemoryRateLimiterOption) ghttp.HandlerFunc {
	limiter := newMemoryTokenBucketRateLimiter(option)
	return func(r *ghttp.Request) {
		key := limiter.option.KeyFunc(r)
		if !limiter.AllowN(r.GetCtx(), key, 1) {
			limiter.option.DenyHandler(r)
			return
		}
		limiter.option.AllowHandler(r)
	}
}

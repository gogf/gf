// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/contrib/middleware/glimiter/v2"
)

var (
	redisConfig = &gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      10, // Use a separate database for testing
	}
)

// createRedisLimiter creates a new Redis limiter for testing.
// It creates a unique key prefix to avoid test interference.
func createRedisLimiter(limit int, window time.Duration) (*glimiter.RedisLimiter, *gredis.Redis, error) {
	redis, err := gredis.New(redisConfig)
	if err != nil {
		return nil, nil, err
	}
	limiter := glimiter.NewRedisLimiter(redis, limit, window)
	return limiter, redis, nil
}

func TestRedisLimiter_Allow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(5, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-allow-" + guid.S()

		// First 5 requests should be allowed
		for i := 0; i < 5; i++ {
			allowed, err := limiter.Allow(ctx, key)
			t.AssertNil(err)
			t.Assert(allowed, true)
		}

		// 6th request should be denied
		allowed, err := limiter.Allow(ctx, key)
		t.AssertNil(err)
		t.Assert(allowed, false)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_AllowN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(10, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-allown-" + guid.S()

		// Allow 3 requests at once
		allowed, err := limiter.AllowN(ctx, key, 3)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Allow another 5 requests
		allowed, err = limiter.AllowN(ctx, key, 5)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Trying to allow 3 more should fail (would exceed 10)
		allowed, err = limiter.AllowN(ctx, key, 3)
		t.AssertNil(err)
		t.Assert(allowed, false)

		// But 2 should still work
		allowed, err = limiter.AllowN(ctx, key, 2)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_GetRemaining(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(10, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-remaining-" + guid.S()

		// Initially, all quota is available
		remaining, err := limiter.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 10)

		// Use 3 quota
		limiter.AllowN(ctx, key, 3)
		remaining, err = limiter.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 7)

		// Use 5 more
		limiter.AllowN(ctx, key, 5)
		remaining, err = limiter.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 2)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_Reset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(5, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-reset-" + guid.S()

		// Use all quota
		for i := 0; i < 5; i++ {
			limiter.Allow(ctx, key)
		}

		// Should be denied
		allowed, _ := limiter.Allow(ctx, key)
		t.Assert(allowed, false)

		// Reset
		err = limiter.Reset(ctx, key)
		t.AssertNil(err)

		// Should be allowed again
		allowed, err = limiter.Allow(ctx, key)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_SlidingWindow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(5, 500*time.Millisecond)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-sliding-" + guid.S()

		// Use all quota
		for i := 0; i < 5; i++ {
			allowed, err := limiter.Allow(ctx, key)
			t.AssertNil(err)
			t.Assert(allowed, true)
		}

		// Should be denied
		allowed, _ := limiter.Allow(ctx, key)
		t.Assert(allowed, false)

		// Wait for sliding window to partially expire
		time.Sleep(550 * time.Millisecond)

		// Should be allowed again (old requests expired)
		allowed, err = limiter.Allow(ctx, key)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(100, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-concurrent-" + guid.S()

		var (
			wg           sync.WaitGroup
			successCount int64
			goroutineNum = 20
			requestPerGo = 10
		)

		// Launch concurrent requests
		for i := 0; i < goroutineNum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < requestPerGo; j++ {
					allowed, err := limiter.Allow(ctx, key)
					if err == nil && allowed {
						atomic.AddInt64(&successCount, 1)
					}
				}
			}()
		}

		wg.Wait()

		// Exactly 100 requests should succeed (Redis Lua script ensures atomicity)
		t.Assert(successCount, 100)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_Wait(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(2, 300*time.Millisecond)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-wait-" + guid.S()

		// Use all quota
		limiter.Allow(ctx, key)
		limiter.Allow(ctx, key)

		// Wait should block then succeed after expiration
		start := time.Now()
		err = limiter.Wait(ctx, key)
		elapsed := time.Since(start)

		t.AssertNil(err)
		// Should wait at least 200ms for window to expire
		t.Assert(elapsed > 200*time.Millisecond, true)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_WaitCancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(1, time.Hour)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		key := "test-wait-cancel-" + guid.S()

		// Use quota
		limiter.Allow(ctx, key)

		// Wait should fail due to context timeout
		err = limiter.Wait(ctx, key)
		t.AssertNE(err, nil)

		// Cleanup
		limiter.Reset(context.Background(), key)
	})
}

func TestRedisLimiter_MultipleKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(3, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key1 := "test-multi-key1-" + guid.S()
		key2 := "test-multi-key2-" + guid.S()

		// Use all quota for key1
		for i := 0; i < 3; i++ {
			limiter.Allow(ctx, key1)
		}

		// key1 should be denied
		allowed, _ := limiter.Allow(ctx, key1)
		t.Assert(allowed, false)

		// key2 should still be allowed (different key)
		allowed, err = limiter.Allow(ctx, key2)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Cleanup
		limiter.Reset(ctx, key1)
		limiter.Reset(ctx, key2)
	})
}

func TestRedisLimiter_InvalidN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(10, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-invalid-" + guid.S()

		// Test with n = 0
		allowed, err := limiter.AllowN(ctx, key, 0)
		t.AssertNE(err, nil)
		t.Assert(allowed, false)

		// Test with n < 0
		allowed, err = limiter.AllowN(ctx, key, -1)
		t.AssertNE(err, nil)
		t.Assert(allowed, false)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_GetLimitAndWindow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(100, 5*time.Minute)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		t.Assert(limiter.GetLimit(), 100)
		t.Assert(limiter.GetWindow(), 5*time.Minute)
	})
}

func TestRedisLimiter_Distributed(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create two limiter instances sharing the same Redis
		limiter1, redis1, err := createRedisLimiter(10, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis1.Close(context.Background())

		limiter2, redis2, err := createRedisLimiter(10, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis2.Close(context.Background())

		ctx := context.Background()
		key := "test-distributed-" + guid.S()

		// Use 5 quota from limiter1
		for i := 0; i < 5; i++ {
			allowed, err := limiter1.Allow(ctx, key)
			t.AssertNil(err)
			t.Assert(allowed, true)
		}

		// Check remaining from limiter2
		remaining, err := limiter2.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 5)

		// Use 3 more from limiter2
		for i := 0; i < 3; i++ {
			allowed, err := limiter2.Allow(ctx, key)
			t.AssertNil(err)
			t.Assert(allowed, true)
		}

		// Check remaining from limiter1
		remaining, err = limiter1.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 2)

		// Cleanup
		limiter1.Reset(ctx, key)
	})
}

func TestRedisLimiter_Precision(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(5, 1*time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-precision-" + guid.S()

		// Use quota at different times within window
		for i := 0; i < 3; i++ {
			allowed, err := limiter.Allow(ctx, key)
			t.AssertNil(err)
			t.Assert(allowed, true)
			time.Sleep(100 * time.Millisecond)
		}

		// Should still have quota
		remaining, err := limiter.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 2)

		// Wait for first requests to expire (sliding window)
		time.Sleep(800 * time.Millisecond)

		// First requests should be expired now
		remaining, err = limiter.GetRemaining(ctx, key)
		t.AssertNil(err)
		// Should have more quota as old requests expired
		t.Assert(remaining > 2, true)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

func TestRedisLimiter_LargeN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter, redis, err := createRedisLimiter(1000, time.Second)
		if err != nil {
			t.Skip("Redis not available:", err)
			return
		}
		defer redis.Close(context.Background())

		ctx := context.Background()
		key := "test-large-" + guid.S()

		// Allow large number at once
		allowed, err := limiter.AllowN(ctx, key, 500)
		t.AssertNil(err)
		t.Assert(allowed, true)

		remaining, err := limiter.GetRemaining(ctx, key)
		t.AssertNil(err)
		t.Assert(remaining, 500)

		// Allow another 500
		allowed, err = limiter.AllowN(ctx, key, 500)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Should be at limit
		allowed, err = limiter.AllowN(ctx, key, 1)
		t.AssertNil(err)
		t.Assert(allowed, false)

		// Cleanup
		limiter.Reset(ctx, key)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/middleware/glimiter/v2"
)

func TestMemoryLimiter_Allow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(5, time.Second)
		ctx := context.Background()

		// First 5 requests should be allowed
		for i := 0; i < 5; i++ {
			allowed, err := limiter.Allow(ctx, "test-key")
			t.AssertNil(err)
			t.Assert(allowed, true)
		}

		// 6th request should be denied
		allowed, err := limiter.Allow(ctx, "test-key")
		t.AssertNil(err)
		t.Assert(allowed, false)
	})
}

func TestMemoryLimiter_AllowN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(10, time.Second)
		ctx := context.Background()

		// Allow 3 requests at once
		allowed, err := limiter.AllowN(ctx, "test-key", 3)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Allow another 5 requests
		allowed, err = limiter.AllowN(ctx, "test-key", 5)
		t.AssertNil(err)
		t.Assert(allowed, true)

		// Trying to allow 3 more should fail (would exceed 10)
		allowed, err = limiter.AllowN(ctx, "test-key", 3)
		t.AssertNil(err)
		t.Assert(allowed, false)

		// But 2 should still work
		allowed, err = limiter.AllowN(ctx, "test-key", 2)
		t.AssertNil(err)
		t.Assert(allowed, true)
	})
}

func TestMemoryLimiter_GetRemaining(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(10, time.Second)
		ctx := context.Background()

		// Initially, all quota is available
		remaining, err := limiter.GetRemaining(ctx, "test-key")
		t.AssertNil(err)
		t.Assert(remaining, 10)

		// Use 3 quota
		limiter.AllowN(ctx, "test-key", 3)
		remaining, err = limiter.GetRemaining(ctx, "test-key")
		t.AssertNil(err)
		t.Assert(remaining, 7)

		// Use 5 more
		limiter.AllowN(ctx, "test-key", 5)
		remaining, err = limiter.GetRemaining(ctx, "test-key")
		t.AssertNil(err)
		t.Assert(remaining, 2)
	})
}

func TestMemoryLimiter_Reset(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(5, time.Second)
		ctx := context.Background()

		// Use all quota
		for i := 0; i < 5; i++ {
			limiter.Allow(ctx, "test-key")
		}

		// Should be denied
		allowed, _ := limiter.Allow(ctx, "test-key")
		t.Assert(allowed, false)

		// Reset
		err := limiter.Reset(ctx, "test-key")
		t.AssertNil(err)

		// Should be allowed again
		allowed, err = limiter.Allow(ctx, "test-key")
		t.AssertNil(err)
		t.Assert(allowed, true)
	})
}

func TestMemoryLimiter_Expiration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(5, 200*time.Millisecond)
		ctx := context.Background()

		// Use all quota
		for i := 0; i < 5; i++ {
			limiter.Allow(ctx, "test-key")
		}

		// Should be denied
		allowed, _ := limiter.Allow(ctx, "test-key")
		t.Assert(allowed, false)

		// Wait for expiration
		time.Sleep(250 * time.Millisecond)

		// Should be allowed again after expiration
		allowed, err := limiter.Allow(ctx, "test-key")
		t.AssertNil(err)
		t.Assert(allowed, true)
	})
}

func TestMemoryLimiter_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(100, time.Second)
		ctx := context.Background()

		var (
			wg           sync.WaitGroup
			successCount = gtype.NewInt64(0)
			goroutineNum = 20
			requestPerGo = 10
		)

		// Launch concurrent requests
		for i := 0; i < goroutineNum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < requestPerGo; j++ {
					allowed, err := limiter.Allow(ctx, "test-key")
					if err == nil && allowed {
						successCount.Add(1)
					}
				}
			}()
		}

		wg.Wait()

		// Exactly 100 requests should succeed
		t.Assert(successCount.Val(), 100)
	})
}

func TestMemoryLimiter_Wait(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(2, 200*time.Millisecond)
		ctx := context.Background()

		// Use all quota
		limiter.Allow(ctx, "test-key")
		limiter.Allow(ctx, "test-key")

		// Wait should block then succeed after expiration
		start := time.Now()
		err := limiter.Wait(ctx, "test-key")
		elapsed := time.Since(start)

		t.AssertNil(err)
		t.Assert(elapsed > 150*time.Millisecond, true)
	})
}

func TestMemoryLimiter_WaitCancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(1, time.Hour)
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Use quota
		limiter.Allow(ctx, "test-key")

		// Wait should fail due to context timeout
		err := limiter.Wait(ctx, "test-key")
		t.AssertNE(err, nil)
	})
}

func TestMemoryLimiter_MultipleKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(3, time.Second)
		ctx := context.Background()

		// Use all quota for key1
		for i := 0; i < 3; i++ {
			limiter.Allow(ctx, "key1")
		}

		// key1 should be denied
		allowed, _ := limiter.Allow(ctx, "key1")
		t.Assert(allowed, false)

		// key2 should still be allowed (different key)
		allowed, err := limiter.Allow(ctx, "key2")
		t.AssertNil(err)
		t.Assert(allowed, true)
	})
}

func TestMemoryLimiter_InvalidN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(10, time.Second)
		ctx := context.Background()

		// Test with n = 0
		allowed, err := limiter.AllowN(ctx, "test-key", 0)
		t.AssertNE(err, nil)
		t.Assert(allowed, false)

		// Test with n < 0
		allowed, err = limiter.AllowN(ctx, "test-key", -1)
		t.AssertNE(err, nil)
		t.Assert(allowed, false)
	})
}

func TestMemoryLimiter_GetLimitAndWindow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewMemoryLimiter(100, 5*time.Minute)

		t.Assert(limiter.GetLimit(), 100)
		t.Assert(limiter.GetWindow(), 5*time.Minute)
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	"context"
	"github.com/gogf/gf/contrib/glimiter/v2"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestMemoryTokenBucketRateLimiter_Allow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		memoryLimiter := glimiter.NewMemoryTokenBucketRateLimiter(glimiter.MemoryTokenBucketRateLimiterOption{
			Rate:     10,
			Capacity: 20,
			Expire:   time.Second * 5,
		})

		ctx := context.Background()
		key := "test_key"

		for i := 0; i < 10; i++ {
			allowed := memoryLimiter.Allow(ctx, key)
			t.Assert(allowed, true)
		}

		allowed := memoryLimiter.AllowN(ctx, key, 15)
		t.Assert(allowed, false)

		time.Sleep(time.Second * 2)

		allowed = memoryLimiter.AllowN(ctx, key, 10)
		t.Assert(allowed, true)
	})
}

func TestMemoryTokenBucketRateLimiter_AllowN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		memoryLimiter := glimiter.NewMemoryTokenBucketRateLimiter(glimiter.MemoryTokenBucketRateLimiterOption{
			Rate:     5,
			Capacity: 10,
		})

		ctx := context.Background()
		key := "test_key_n"

		allowed := memoryLimiter.AllowN(ctx, key, 5)
		t.Assert(allowed, true)

		allowed = memoryLimiter.AllowN(ctx, key, 8)
		t.Assert(allowed, false)

		allowed = memoryLimiter.AllowN(ctx, key, 0)
		t.Assert(allowed, true)

		allowed = memoryLimiter.AllowN(ctx, key, -1)
		t.Assert(allowed, false)
	})
}

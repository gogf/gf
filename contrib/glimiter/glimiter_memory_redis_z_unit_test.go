// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	"context"
	"github.com/gogf/gf/contrib/glimiter/v2"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"time"
)

func TestRedisMemoryTokenBucketRateLimiter_Allow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiterInstance := glimiter.NewRedisMemoryTokenBucketRateLimiter(glimiter.RedisMemoryTokenBucketRateLimiterOption{
			RedisLimiterOption: glimiter.RedisTokenBucketRateLimiterOption{
				Redis:    re,
				Rate:     10,
				Capacity: 20,
				Expire:   time.Second * 5,
			},
			MemoryLimiterOption: glimiter.MemoryTokenBucketRateLimiterOption{
				Rate:     10,
				Capacity: 20,
				Expire:   time.Second * 5,
			},
		})

		ctx := context.Background()
		key := "test_hybrid_key"

		for i := 0; i < 10; i++ {
			allowed := limiterInstance.AllowN(ctx, key, 1)
			t.Assert(allowed, true)
		}

		allowed := limiterInstance.AllowN(ctx, key, 15)
		t.Assert(allowed, false)

		time.Sleep(time.Second * 2)

		// 再次测试允许通过的请求
		allowed = limiterInstance.AllowN(ctx, key, 10)
		t.Assert(allowed, true)
	})
}

func TestRedisMemoryTokenBucketRateLimiter_Middleware(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiterInstance := glimiter.NewRedisMemoryTokenBucketRateLimiter(glimiter.RedisMemoryTokenBucketRateLimiterOption{
			RedisLimiterOption: glimiter.RedisTokenBucketRateLimiterOption{
				Redis:    re,
				Rate:     1,
				Capacity: 2,
			},
			MemoryLimiterOption: glimiter.MemoryTokenBucketRateLimiterOption{
				Rate:     1,
				Capacity: 2,
			},
		})

		middleware := limiterInstance.Middleware()
		t.AssertNE(middleware, nil)
	})
}

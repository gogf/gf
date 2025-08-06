// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glimiter_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/glimiter/v2"
)

func TestRedisTokenBucketRateLimiter_Allow(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redisLimiter := glimiter.NewRedisTokenBucketRateLimiter(glimiter.RedisTokenBucketRateLimiterOption{
			Redis:    re,
			Rate:     10,
			Capacity: 20,
			Expire:   time.Second * 5,
		})

		ctx := context.Background()
		key := "test_redis_key"

		for i := 0; i < 10; i++ {
			allowed := redisLimiter.Allow(ctx, key)
			t.Assert(allowed, true)
		}

		allowed := redisLimiter.AllowN(ctx, key, 15)
		t.Assert(allowed, false)
		time.Sleep(time.Second * 2)
		allowed = redisLimiter.AllowN(ctx, key, 10)
		t.Assert(allowed, true)
	})
}

func TestRedisTokenBucketRateLimiter_AllowN(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiter := glimiter.NewRedisTokenBucketRateLimiter(glimiter.RedisTokenBucketRateLimiterOption{
			Redis:    re,
			Rate:     5,
			Capacity: 10,
		})

		ctx := context.Background()
		key := "test_redis_key_n"

		allowed := limiter.AllowN(ctx, key, 5)
		t.Assert(allowed, true)

		allowed = limiter.AllowN(ctx, key, 8)
		t.Assert(allowed, false)

		allowed = limiter.AllowN(ctx, key, 0)
		t.Assert(allowed, true)

		allowed = limiter.AllowN(ctx, key, -1)
		t.Assert(allowed, false)
	})
}

func TestRedisTokenBucketRateLimiter_Middleware(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		limiterInstance := glimiter.NewRedisTokenBucketRateLimiter(glimiter.RedisTokenBucketRateLimiterOption{
			Redis:    re,
			Rate:     1,
			Capacity: 2,
		})
		middleware := limiterInstance.Middleware()
		t.AssertNE(middleware, nil)
	})
}

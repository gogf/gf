// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	cacheRedis  = gcache.New()
	redisConfig = &gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      2,
	}
)

func init() {
	redis, err := gredis.New(redisConfig)
	if err != nil {
		panic(err)
	}
	cacheRedis.SetAdapter(gcache.NewAdapterRedis(redis))
}

func Test_AdapterRedis_Basic1(t *testing.T) {
	// Set
	size := 10
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < size; i++ {
			t.AssertNil(cacheRedis.Set(ctx, i, i*10, 0))
		}
		for i := 0; i < size; i++ {
			v, _ := cacheRedis.Get(ctx, i)
			t.Assert(v, i*10)
		}
		n, _ := cacheRedis.Size(ctx)
		t.Assert(n, size)
	})
	// Data
	gtest.C(t, func(t *gtest.T) {
		data, _ := cacheRedis.Data(ctx)
		t.Assert(len(data), size)
		t.Assert(data["0"], "0")
		t.Assert(data["1"], "10")
		t.Assert(data["9"], "90")
	})
	// Clear
	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(cacheRedis.Clear(ctx))
		n, _ := cacheRedis.Size(ctx)
		t.Assert(n, 0)
	})
}

func Test_AdapterRedis_Basic2(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	size := 10
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < size; i++ {
			t.AssertNil(cacheRedis.Set(ctx, i, i*10, -1))
		}
		for i := 0; i < size; i++ {
			v, _ := cacheRedis.Get(ctx, i)
			t.Assert(v, nil)
		}
		n, _ := cacheRedis.Size(ctx)
		t.Assert(n, 0)
	})
}

func Test_AdapterRedis_Basic3(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	size := 10
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < size; i++ {
			t.AssertNil(cacheRedis.Set(ctx, i, i*10, time.Second))
		}
		for i := 0; i < size; i++ {
			v, _ := cacheRedis.Get(ctx, i)
			t.Assert(v, i*10)
		}
		n, _ := cacheRedis.Size(ctx)
		t.Assert(n, size)
	})
	time.Sleep(time.Second * 2)
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < size; i++ {
			v, _ := cacheRedis.Get(ctx, i)
			t.Assert(v, nil)
		}
		n, _ := cacheRedis.Size(ctx)
		t.Assert(n, 0)
	})
}

func Test_AdapterRedis_Update(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "key"
			value1 = "value1"
			value2 = "value2"
		)
		t.AssertNil(cacheRedis.Set(ctx, key, value1, time.Second))
		v, _ := cacheRedis.Get(ctx, key)
		t.Assert(v, value1)

		d, _ := cacheRedis.GetExpire(ctx, key)
		t.Assert(d > time.Millisecond*500, true)
		t.Assert(d <= time.Second, true)

		_, _, err := cacheRedis.Update(ctx, key, value2)
		t.AssertNil(err)

		v, _ = cacheRedis.Get(ctx, key)
		t.Assert(v, value2)
		d, _ = cacheRedis.GetExpire(ctx, key)
		t.Assert(d > time.Millisecond*500, true)
		t.Assert(d <= time.Second, true)
	})
}

func Test_AdapterRedis_UpdateExpire(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key   = "key"
			value = "value"
		)
		t.AssertNil(cacheRedis.Set(ctx, key, value, time.Second))
		v, _ := cacheRedis.Get(ctx, key)
		t.Assert(v, value)

		d, _ := cacheRedis.GetExpire(ctx, key)
		t.Assert(d > time.Millisecond*500, true)
		t.Assert(d <= time.Second, true)

		_, err := cacheRedis.UpdateExpire(ctx, key, time.Second*2)
		t.AssertNil(err)

		d, _ = cacheRedis.GetExpire(ctx, key)
		t.Assert(d > time.Second, true)
		t.Assert(d <= 2*time.Second, true)
	})
}

func Test_AdapterRedis_SetIfNotExist(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "key"
			value1 = "value1"
			value2 = "value2"
		)
		t.AssertNil(cacheRedis.Set(ctx, key, value1, time.Second))
		v, _ := cacheRedis.Get(ctx, key)
		t.Assert(v, value1)

		r, _ := cacheRedis.SetIfNotExist(ctx, key, value2, time.Second*2)

		t.Assert(r, false)

		v, _ = cacheRedis.Get(ctx, key)
		t.Assert(v, value1)

		d, _ := cacheRedis.GetExpire(ctx, key)
		t.Assert(d > time.Millisecond*500, true)
		t.Assert(d <= time.Second, true)
	})
}

func Test_AdapterRedis_GetOrSetFunc(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "key"
			value1 = "valueFunc"
		)
		v, err := cacheRedis.GetOrSetFunc(ctx, key, func(ctx context.Context) (value interface{}, err error) {
			value = value1
			return
		}, 0)
		t.AssertNil(err)
		t.Assert(v, value1)
	})
}

func Test_AdapterRedis_GetOrSetFuncLock(t *testing.T) {
	defer cacheRedis.Clear(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "key"
			value1 = "valueFuncLock"
		)
		v, err := cacheRedis.GetOrSetFuncLock(ctx, key, func(ctx context.Context) (value interface{}, err error) {
			value = value1
			return
		}, time.Second*60)
		t.AssertNil(err)
		t.Assert(v, value1)
	})
}

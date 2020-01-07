// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/test/gtest"
)

//clear 用于清除全局缓存，因gcache api 暂未暴露 Clear 方法
//暂定所有测试用例key的集合为1,2,3，避免不同测试用例间因全局cache共享带来的问题，每个测试用例在测试gcache.XXX之前，先调用clear()
func clear() {
	gcache.Removes(g.Slice{1, 2, 3})
}

func TestCache_Set(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.Set(1, 11, 0)
		gtest.Assert(cache.Get(1), 11)
		gtest.Assert(cache.Contains(1), true)

		clear()
		gcache.Set(1, 11, 0)
		gtest.Assert(gcache.Get(1), 11)
		gtest.Assert(gcache.Contains(1), true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.Set(2, 22, 100*time.Millisecond)
		gtest.Assert(cache.Get(2), 22)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(cache.Get(2), nil)
		time.Sleep(3 * time.Second)
		gtest.Assert(cache.Size(), 0)
		cache.Close()
	})

	gtest.Case(t, func() {
		cache := gcache.New()
		cache.Set(1, 11, 100*time.Millisecond)
		gtest.Assert(cache.Get(1), 11)
		time.Sleep(200 * time.Millisecond)
		gtest.Assert(cache.Get(1), nil)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		for i := 0; i < 10; i++ {
			cache.Set(i, i*10, 0)
		}
		gtest.Assert(len(cache.Keys()), 10)
		gtest.Assert(len(cache.Values()), 10)
		gtest.AssertIN(0, cache.Keys())
		gtest.AssertIN(90, cache.Values())
	})
}

func TestCache_LRU(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New(2)
		for i := 0; i < 10; i++ {
			cache.Set(i, i, 0)
		}
		gtest.Assert(cache.Size(), 10)
		gtest.Assert(cache.Get(6), 6)
		time.Sleep(4 * time.Second)
		gtest.Assert(cache.Size(), 2)
		gtest.Assert(cache.Get(6), 6)
		gtest.Assert(cache.Get(1), nil)
		cache.Close()
	})
}

func TestCache_LRU_expire(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New(2)
		cache.Set(1, nil, 1000)
		gtest.Assert(cache.Size(), 1)
		gtest.Assert(cache.Get(1), nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.SetIfNotExist(1, 11, 0)
		gtest.Assert(cache.Get(1), 11)
		cache.SetIfNotExist(1, 22, 0)
		gtest.Assert(cache.Get(1), 11)
		cache.SetIfNotExist(2, 22, 0)
		gtest.Assert(cache.Get(2), 22)

		clear()
		gcache.SetIfNotExist(1, 11, 0)
		gtest.Assert(gcache.Get(1), 11)
		gcache.SetIfNotExist(1, 22, 0)
		gtest.Assert(gcache.Get(1), 11)
	})
}

func TestCache_Sets(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		gtest.Assert(cache.Get(1), 11)

		clear()
		gcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		gtest.Assert(gcache.Get(1), 11)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.GetOrSet(1, 11, 0)
		gtest.Assert(cache.Get(1), 11)
		cache.GetOrSet(1, 111, 0)
		gtest.Assert(cache.Get(1), 11)

		clear()
		gcache.GetOrSet(1, 11, 0)
		gtest.Assert(gcache.Get(1), 11)
		gcache.GetOrSet(1, 111, 0)
		gtest.Assert(gcache.Get(1), 11)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.GetOrSetFunc(1, func() interface{} {
			return 11
		}, 0)
		gtest.Assert(cache.Get(1), 11)
		cache.GetOrSetFunc(1, func() interface{} {
			return 111
		}, 0)
		gtest.Assert(cache.Get(1), 11)

		clear()
		gcache.GetOrSetFunc(1, func() interface{} {
			return 11
		}, 0)
		gtest.Assert(gcache.Get(1), 11)
		gcache.GetOrSetFunc(1, func() interface{} {
			return 111
		}, 0)
		gtest.Assert(gcache.Get(1), 11)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.GetOrSetFuncLock(1, func() interface{} {
			return 11
		}, 0)
		gtest.Assert(cache.Get(1), 11)
		cache.GetOrSetFuncLock(1, func() interface{} {
			return 111
		}, 0)
		gtest.Assert(cache.Get(1), 11)

		clear()
		gcache.GetOrSetFuncLock(1, func() interface{} {
			return 11
		}, 0)
		gtest.Assert(gcache.Get(1), 11)
		gcache.GetOrSetFuncLock(1, func() interface{} {
			return 111
		}, 0)
		gtest.Assert(gcache.Get(1), 11)
	})
}

func TestCache_Clear(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		cache.Clear()
		gtest.Assert(cache.Size(), 0)
	})
}

func TestCache_SetConcurrency(t *testing.T) {
	gtest.Case(t, func() {
		cache := gcache.New()
		pool := grpool.New(4)
		go func() {
			for {
				pool.Add(func() {
					cache.SetIfNotExist(1, 11, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			//t.Log("first part end")
		}

		go func() {
			for {
				pool.Add(func() {
					cache.SetIfNotExist(1, nil, 10)
				})
			}
		}()
		select {
		case <-time.After(2 * time.Second):
			//t.Log("second part end")
		}
	})
}

func TestCache_Basic(t *testing.T) {
	gtest.Case(t, func() {
		{
			cache := gcache.New()
			cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			gtest.Assert(cache.Contains(1), true)
			gtest.Assert(cache.Get(1), 11)
			data := cache.Data()
			gtest.Assert(data[1], 11)
			gtest.Assert(data[2], 22)
			gtest.Assert(data[3], nil)
			gtest.Assert(cache.Size(), 2)
			keys := cache.Keys()
			gtest.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs := cache.KeyStrings()
			gtest.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values := cache.Values()
			gtest.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1 := cache.Remove(1)
			gtest.Assert(removeData1, 11)
			gtest.Assert(cache.Size(), 1)
			cache.Removes(g.Slice{2})
			gtest.Assert(cache.Size(), 0)
		}

		clear()
		{
			gcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			gtest.Assert(gcache.Contains(1), true)
			gtest.Assert(gcache.Get(1), 11)
			data := gcache.Data()
			gtest.Assert(data[1], 11)
			gtest.Assert(data[2], 22)
			gtest.Assert(data[3], nil)
			gtest.Assert(gcache.Size(), 2)
			keys := gcache.Keys()
			gtest.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs := gcache.KeyStrings()
			gtest.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values := gcache.Values()
			gtest.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1 := gcache.Remove(1)
			gtest.Assert(removeData1, 11)
			gtest.Assert(gcache.Size(), 1)
			gcache.Removes(g.Slice{2})
			gtest.Assert(gcache.Size(), 0)
		}
	})
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/grand"
	"testing"
	"time"

	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/grpool"
	"github.com/gogf/gf/test/gtest"
)

func TestCache_GCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gcache.Set(1, 11, 0)
		defer gcache.Removes(g.Slice{1, 2, 3})
		t.Assert(gcache.Get(1), 11)
		t.Assert(gcache.Contains(1), true)
	})
}

func TestCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := gcache.New()
		defer c.Close()
		c.Set(1, 11, 0)
		t.Assert(c.Get(1), 11)
		t.Assert(c.Contains(1), true)
	})
}

func TestCache_GetVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := gcache.New()
		defer c.Close()
		c.Set(1, 11, 0)
		t.Assert(c.Get(1), 11)
		t.Assert(c.Contains(1), true)
		t.Assert(c.GetVar(1).Int(), 11)
		t.Assert(c.GetVar(2).Int(), 0)
		t.Assert(c.GetVar(2).IsNil(), true)
		t.Assert(c.GetVar(2).IsEmpty(), true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Set(2, 22, 100*time.Millisecond)
		t.Assert(cache.Get(2), 22)
		time.Sleep(200 * time.Millisecond)
		t.Assert(cache.Get(2), nil)
		time.Sleep(3 * time.Second)
		t.Assert(cache.Size(), 0)
		cache.Close()
	})

	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Set(1, 11, 100*time.Millisecond)
		t.Assert(cache.Get(1), 11)
		time.Sleep(200 * time.Millisecond)
		t.Assert(cache.Get(1), nil)
	})
}

func TestCache_SetVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Set(1, 11, 3*time.Second)
		expireBefore, _ := cache.GetExpire(1)
		cache.SetVar(1, 12)
		expireAfter, _ := cache.GetExpire(1)
		t.Assert(cache.GetVar(1), 12)
		t.Assert(expireBefore, expireAfter)
	})
}

func BenchmarkMemCache_GetSetExpire(b *testing.B) {
	cache := gcache.New()
	cache.Set(1, 11, 3*time.Second)
	if expire, ok := cache.GetExpire(1); ok {
		glog.Println(expire)
	}
	for i := 0; i < b.N; i++ {
		r := time.Duration(grand.N(5, 10))
		cache.SetExpire(1, r*time.Second)
		//cache.SetExpire(1, 7*time.Second)
		if _, ok := cache.GetExpire(1); !ok {
			panic(fmt.Sprintf("[ERROR] %s", "GetExpire error"))
		}
	}
}

func TestCache_Keys_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		for i := 0; i < 10; i++ {
			cache.Set(i, i*10, 0)
		}
		t.Assert(len(cache.Keys()), 10)
		t.Assert(len(cache.Values()), 10)
		t.AssertIN(0, cache.Keys())
		t.AssertIN(90, cache.Values())
	})
}

func TestCache_LRU(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New(2)
		for i := 0; i < 10; i++ {
			cache.Set(i, i, 0)
		}
		t.Assert(cache.Size(), 10)
		t.Assert(cache.Get(6), 6)
		time.Sleep(4 * time.Second)
		t.Assert(cache.Size(), 2)
		t.Assert(cache.Get(6), 6)
		t.Assert(cache.Get(1), nil)
		cache.Close()
	})
}

func TestCache_LRU_expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New(2)
		cache.Set(1, nil, 1000)
		t.Assert(cache.Size(), 1)
		t.Assert(cache.Get(1), nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.SetIfNotExist(1, 11, 0)
		t.Assert(cache.Get(1), 11)
		cache.SetIfNotExist(1, 22, 0)
		t.Assert(cache.Get(1), 11)
		cache.SetIfNotExist(2, 22, 0)
		t.Assert(cache.Get(2), 22)

		gcache.Removes(g.Slice{1, 2, 3})
		gcache.SetIfNotExist(1, 11, 0)
		t.Assert(gcache.Get(1), 11)
		gcache.SetIfNotExist(1, 22, 0)
		t.Assert(gcache.Get(1), 11)
	})
}

func TestCache_Sets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		t.Assert(cache.Get(1), 11)

		gcache.Removes(g.Slice{1, 2, 3})
		gcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		t.Assert(gcache.Get(1), 11)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.GetOrSet(1, 11, 0)
		t.Assert(cache.Get(1), 11)
		cache.GetOrSet(1, 111, 0)
		t.Assert(cache.Get(1), 11)

		gcache.Removes(g.Slice{1, 2, 3})
		gcache.GetOrSet(1, 11, 0)
		t.Assert(gcache.Get(1), 11)
		gcache.GetOrSet(1, 111, 0)
		t.Assert(gcache.Get(1), 11)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.GetOrSetFunc(1, func() interface{} {
			return 11
		}, 0)
		t.Assert(cache.Get(1), 11)
		cache.GetOrSetFunc(1, func() interface{} {
			return 111
		}, 0)
		t.Assert(cache.Get(1), 11)

		gcache.Removes(g.Slice{1, 2, 3})
		gcache.GetOrSetFunc(1, func() interface{} {
			return 11
		}, 0)
		t.Assert(gcache.Get(1), 11)
		gcache.GetOrSetFunc(1, func() interface{} {
			return 111
		}, 0)
		t.Assert(gcache.Get(1), 11)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.GetOrSetFuncLock(1, func() interface{} {
			return 11
		}, 0)
		t.Assert(cache.Get(1), 11)
		cache.GetOrSetFuncLock(1, func() interface{} {
			return 111
		}, 0)
		t.Assert(cache.Get(1), 11)

		gcache.Removes(g.Slice{1, 2, 3})
		gcache.GetOrSetFuncLock(1, func() interface{} {
			return 11
		}, 0)
		t.Assert(gcache.Get(1), 11)
		gcache.GetOrSetFuncLock(1, func() interface{} {
			return 111
		}, 0)
		t.Assert(gcache.Get(1), 11)
	})
}

func TestCache_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		cache.Clear()
		t.Assert(cache.Size(), 0)
	})
}

func TestCache_SetConcurrency(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
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
	gtest.C(t, func(t *gtest.T) {
		{
			cache := gcache.New()
			cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			t.Assert(cache.Contains(1), true)
			t.Assert(cache.Get(1), 11)
			data := cache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			t.Assert(cache.Size(), 2)
			keys := cache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs := cache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values := cache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1 := cache.Remove(1)
			t.Assert(removeData1, 11)
			t.Assert(cache.Size(), 1)
			cache.Removes(g.Slice{2})
			t.Assert(cache.Size(), 0)
		}

		gcache.Removes(g.Slice{1, 2, 3})
		{
			gcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			t.Assert(gcache.Contains(1), true)
			t.Assert(gcache.Get(1), 11)
			data := gcache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			t.Assert(gcache.Size(), 2)
			keys := gcache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs := gcache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values := gcache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1 := gcache.Remove(1)
			t.Assert(removeData1, 11)
			t.Assert(gcache.Size(), 1)
			gcache.Removes(g.Slice{2})
			t.Assert(gcache.Size(), 0)
		}
	})
}

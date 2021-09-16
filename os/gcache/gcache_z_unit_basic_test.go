// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*" -benchmem

package gcache_test

import (
	"context"
	"github.com/gogf/gf/util/guid"
	"math"
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
		t.AssertNil(gcache.Set(1, 11, 0))
		defer gcache.Remove(g.Slice{1, 2, 3}...)
		v, _ := gcache.Get(1)
		t.Assert(v, 11)
		b, _ := gcache.Contains(1)
		t.Assert(b, true)
	})
}

func TestCache_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := gcache.New()
		defer c.Close()
		t.Assert(c.Set(1, 11, 0), nil)
		v, _ := c.Get(1)
		t.Assert(v, 11)
		b, _ := c.Contains(1)
		t.Assert(b, true)
	})
}

func TestCache_GetVar(t *testing.T) {
	c := gcache.New()
	defer c.Close()
	gtest.C(t, func(t *gtest.T) {
		t.Assert(c.Set(1, 11, 0), nil)
		v, _ := c.Get(1)
		t.Assert(v, 11)
		b, _ := c.Contains(1)
		t.Assert(b, true)
	})
	gtest.C(t, func(t *gtest.T) {
		v, _ := c.GetVar(1)
		t.Assert(v.Int(), 11)
		v, _ = c.GetVar(2)
		t.Assert(v.Int(), 0)
		t.Assert(v.IsNil(), true)
		t.Assert(v.IsEmpty(), true)
	})
}

func TestCache_Set_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		t.Assert(cache.Set(2, 22, 100*time.Millisecond), nil)
		v, _ := cache.Get(2)
		t.Assert(v, 22)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(2)
		t.Assert(v, nil)
		time.Sleep(3 * time.Second)
		n, _ := cache.Size()
		t.Assert(n, 0)
		t.Assert(cache.Close(), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		t.Assert(cache.Set(1, 11, 100*time.Millisecond), nil)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
		time.Sleep(200 * time.Millisecond)
		v, _ = cache.Get(1)
		t.Assert(v, nil)
	})
}

func TestCache_Update(t *testing.T) {
	// gcache
	gtest.C(t, func(t *gtest.T) {
		key := guid.S()
		t.AssertNil(gcache.Set(key, 11, 3*time.Second))
		expire1, _ := gcache.GetExpire(key)
		oldValue, exist, err := gcache.Update(key, 12)
		t.AssertNil(err)
		t.Assert(oldValue, 11)
		t.Assert(exist, true)

		expire2, _ := gcache.GetExpire(key)
		v, _ := gcache.GetVar(key)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
	// gcache.Cache
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		t.AssertNil(cache.Set(1, 11, 3*time.Second))

		oldValue, exist, err := cache.Update(1, 12)
		t.AssertNil(err)
		t.Assert(oldValue, 11)
		t.Assert(exist, true)

		expire1, _ := cache.GetExpire(1)
		expire2, _ := cache.GetExpire(1)
		v, _ := cache.GetVar(1)
		t.Assert(v, 12)
		t.Assert(math.Ceil(expire1.Seconds()), math.Ceil(expire2.Seconds()))
	})
}

func TestCache_UpdateExpire(t *testing.T) {
	// gcache
	gtest.C(t, func(t *gtest.T) {
		key := guid.S()
		t.AssertNil(gcache.Set(key, 11, 3*time.Second))
		defer gcache.Remove(key)
		oldExpire, _ := gcache.GetExpire(key)
		newExpire := 10 * time.Second
		oldExpire2, err := gcache.UpdateExpire(key, newExpire)
		t.AssertNil(err)
		t.Assert(oldExpire2, oldExpire)

		e, _ := gcache.GetExpire(key)
		t.AssertNE(e, oldExpire)
		e, _ = gcache.GetExpire(key)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
	// gcache.Cache
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		t.AssertNil(cache.Set(1, 11, 3*time.Second))
		oldExpire, _ := cache.GetExpire(1)
		newExpire := 10 * time.Second
		oldExpire2, err := cache.UpdateExpire(1, newExpire)
		t.AssertNil(err)
		t.Assert(oldExpire2, oldExpire)

		e, _ := cache.GetExpire(1)
		t.AssertNE(e, oldExpire)

		e, _ = cache.GetExpire(1)
		t.Assert(math.Ceil(e.Seconds()), 10)
	})
}

func TestCache_Keys_Values(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := gcache.New()
		for i := 0; i < 10; i++ {
			t.Assert(c.Set(i, i*10, 0), nil)
		}
		var (
			keys, _   = c.Keys()
			values, _ = c.Values()
		)
		t.Assert(len(keys), 10)
		t.Assert(len(values), 10)
		t.AssertIN(0, keys)
		t.AssertIN(90, values)
	})
}

func TestCache_LRU(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New(2)
		for i := 0; i < 10; i++ {
			t.AssertNil(cache.Set(i, i, 0))
		}
		n, _ := cache.Size()
		t.Assert(n, 10)
		v, _ := cache.Get(6)
		t.Assert(v, 6)
		time.Sleep(4 * time.Second)
		n, _ = cache.Size()
		t.Assert(n, 2)
		v, _ = cache.Get(6)
		t.Assert(v, 6)
		v, _ = cache.Get(1)
		t.Assert(v, nil)
		t.Assert(cache.Close(), nil)
	})
}

func TestCache_LRU_expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New(2)
		t.Assert(cache.Set(1, nil, 1000), nil)
		n, _ := cache.Size()
		t.Assert(n, 1)
		v, _ := cache.Get(1)

		t.Assert(v, nil)
	})
}

func TestCache_SetIfNotExist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		ok, err := cache.SetIfNotExist(1, 11, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ := cache.Get(1)
		t.Assert(v, 11)

		ok, err = cache.SetIfNotExist(1, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, _ = cache.Get(1)
		t.Assert(v, 11)

		ok, err = cache.SetIfNotExist(2, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ = cache.Get(2)
		t.Assert(v, 22)

		gcache.Remove(g.Slice{1, 2, 3}...)
		ok, err = gcache.SetIfNotExist(1, 11, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ = gcache.Get(1)
		t.Assert(v, 11)

		ok, err = gcache.SetIfNotExist(1, 22, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, _ = gcache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_SetIfNotExistFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		exist, err := cache.SetIfNotExistFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, _ := cache.Get(1)
		t.Assert(v, 11)

		exist, err = cache.SetIfNotExistFunc(1, func() (interface{}, error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
	gtest.C(t, func(t *gtest.T) {
		gcache.Remove(g.Slice{1, 2, 3}...)

		ok, err := gcache.SetIfNotExistFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(ok, true)

		v, _ := gcache.Get(1)
		t.Assert(v, 11)

		ok, err = gcache.SetIfNotExistFunc(1, func() (interface{}, error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(ok, false)

		v, _ = gcache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_SetIfNotExistFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		exist, err := cache.SetIfNotExistFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, _ := cache.Get(1)
		t.Assert(v, 11)

		exist, err = cache.SetIfNotExistFuncLock(1, func() (interface{}, error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
	gtest.C(t, func(t *gtest.T) {
		gcache.Remove(g.Slice{1, 2, 3}...)

		exist, err := gcache.SetIfNotExistFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, true)

		v, _ := gcache.Get(1)
		t.Assert(v, 11)

		exist, err = gcache.SetIfNotExistFuncLock(1, func() (interface{}, error) {
			return 22, nil
		}, 0)
		t.AssertNil(err)
		t.Assert(exist, false)

		v, _ = gcache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_Sets(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		t.AssertNil(cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0))
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		gcache.Remove(g.Slice{1, 2, 3}...)
		t.AssertNil(gcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0))
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		value, err := cache.GetOrSet(1, 11, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, _ := cache.Get(1)
		t.Assert(v, 11)
		value, err = cache.GetOrSet(1, 111, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})

	gtest.C(t, func(t *gtest.T) {
		gcache.Remove(g.Slice{1, 2, 3}...)
		value, err := gcache.GetOrSet(1, 11, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, err := gcache.Get(1)
		t.AssertNil(err)
		t.Assert(v, 11)

		value, err = gcache.GetOrSet(1, 111, 0)
		t.AssertNil(err)
		t.Assert(value, 11)

		v, err = gcache.Get(1)
		t.AssertNil(err)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.GetOrSetFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		cache.GetOrSetFunc(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		gcache.Remove(g.Slice{1, 2, 3}...)

		gcache.GetOrSetFunc(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		gcache.GetOrSetFunc(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_GetOrSetFuncLock(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ := cache.Get(1)
		t.Assert(v, 11)

		cache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		gcache.Remove(g.Slice{1, 2, 3}...)
		gcache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 11, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)

		gcache.GetOrSetFuncLock(1, func() (interface{}, error) {
			return 111, nil
		}, 0)
		v, _ = cache.Get(1)
		t.Assert(v, 11)
	})
}

func TestCache_Clear(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		cache.Clear()
		n, _ := cache.Size()
		t.Assert(n, 0)
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
			b, _ := cache.Contains(1)
			t.Assert(b, true)
			v, _ := cache.Get(1)
			t.Assert(v, 11)
			data, _ := cache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := cache.Size()
			t.Assert(n, 2)
			keys, _ := cache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs, _ := cache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values, _ := cache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1, _ := cache.Remove(1)
			t.Assert(removeData1, 11)
			n, _ = cache.Size()
			t.Assert(n, 1)

			cache.Remove(2)
			n, _ = cache.Size()
			t.Assert(n, 0)
		}

		gcache.Remove(g.Slice{1, 2, 3}...)
		{
			gcache.Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
			b, _ := gcache.Contains(1)
			t.Assert(b, true)
			v, _ := gcache.Get(1)
			t.Assert(v, 11)
			data, _ := gcache.Data()
			t.Assert(data[1], 11)
			t.Assert(data[2], 22)
			t.Assert(data[3], nil)
			n, _ := gcache.Size()
			t.Assert(n, 2)
			keys, _ := gcache.Keys()
			t.Assert(gset.NewFrom(g.Slice{1, 2}).Equal(gset.NewFrom(keys)), true)
			keyStrs, _ := gcache.KeyStrings()
			t.Assert(gset.NewFrom(g.Slice{"1", "2"}).Equal(gset.NewFrom(keyStrs)), true)
			values, _ := gcache.Values()
			t.Assert(gset.NewFrom(g.Slice{11, 22}).Equal(gset.NewFrom(values)), true)
			removeData1, _ := gcache.Remove(1)
			t.Assert(removeData1, 11)
			n, _ = gcache.Size()
			t.Assert(n, 1)
			gcache.Remove(2)
			n, _ = gcache.Size()
			t.Assert(n, 0)
		}
	})
}

func TestCache_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		t.AssertNil(cache.Set(1, 11, 0))
		t.AssertNil(cache.Set(2, 22, 0))
		t.AssertNil(cache.Set(3, 33, 0))
		t.AssertNil(cache.Removes(g.Slice{2, 3}))

		ok, err := cache.Contains(1)
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = cache.Contains(2)
		t.AssertNil(err)
		t.Assert(ok, false)
	})

	gtest.C(t, func(t *gtest.T) {
		t.AssertNil(gcache.Set(1, 11, 0))
		t.AssertNil(gcache.Set(2, 22, 0))
		t.AssertNil(gcache.Set(3, 33, 0))
		t.AssertNil(gcache.Removes(g.Slice{2, 3}))

		ok, err := gcache.Contains(1)
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = gcache.Contains(2)
		t.AssertNil(err)
		t.Assert(ok, false)
	})
}

func TestCache_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cache := gcache.New()
		cache.Ctx(context.Background()).Sets(g.MapAnyAny{1: 11, 2: 22}, 0)
		b, _ := cache.Contains(1)
		t.Assert(b, true)
		v, _ := cache.Get(1)
		t.Assert(v, 11)
	})
}

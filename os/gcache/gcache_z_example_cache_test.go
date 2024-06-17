// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcache_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

func ExampleNew() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly.
	c := gcache.New()

	// Set cache without expiration
	c.Set(ctx, "k1", "v1", 0)

	// Get cache
	v, _ := c.Get(ctx, "k1")
	fmt.Println(v)

	// Get cache size
	n, _ := c.Size(ctx)
	fmt.Println(n)

	// Does the specified key name exist in the cache
	b, _ := c.Contains(ctx, "k1")
	fmt.Println(b)

	// Delete and return the deleted key value
	fmt.Println(c.Remove(ctx, "k1"))

	// Close the cache object and let the GC reclaim resources

	c.Close(ctx)

	// Output:
	// v1
	// 1
	// true
	// v1 <nil>
}

func ExampleCache_Set() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set cache without expiration
	c.Set(ctx, "k1", g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9}, 0)

	// Get cache
	fmt.Println(c.Get(ctx, "k1"))

	// Output:
	// [1,2,3,4,5,6,7,8,9] <nil>
}

func ExampleCache_SetIfNotExist() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Write when the key name does not exist, and set the expiration time to 1000 milliseconds
	k1, err := c.SetIfNotExist(ctx, "k1", "v1", 1000*time.Millisecond)
	fmt.Println(k1, err)

	// Returns false when the key name already exists
	k2, err := c.SetIfNotExist(ctx, "k1", "v2", 1000*time.Millisecond)
	fmt.Println(k2, err)

	// Print the current list of key values
	keys1, _ := c.Keys(ctx)
	fmt.Println(keys1)

	// It does not expire if `duration` == 0. It deletes the `key` if `duration` < 0 or given `value` is nil.
	c.SetIfNotExist(ctx, "k1", 0, -10000)

	// Wait 1.5 second for K1: V1 to expire automatically
	time.Sleep(1500 * time.Millisecond)

	// Print the current key value pair again and find that K1: V1 has expired
	keys2, _ := c.Keys(ctx)
	fmt.Println(keys2)

	// Output:
	// true <nil>
	// false <nil>
	// [k1]
	// [<nil>]
}

func ExampleCache_SetMap() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// map[interface{}]interface{}
	data := g.MapAnyAny{
		"k1": "v1",
		"k2": "v2",
		"k3": "v3",
	}

	// Sets batch sets cache with key-value pairs by `data`, which is expired after `duration`.
	// It does not expire if `duration` == 0. It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
	c.SetMap(ctx, data, 1000*time.Millisecond)

	// Gets the specified key value
	v1, _ := c.Get(ctx, "k1")
	v2, _ := c.Get(ctx, "k2")
	v3, _ := c.Get(ctx, "k3")

	fmt.Println(v1, v2, v3)

	// Output:
	// v1 v2 v3
}

func ExampleCache_Size() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Add 10 elements without expiration
	for i := 0; i < 10; i++ {
		c.Set(ctx, i, i, 0)
	}

	// Size returns the number of items in the cache.
	n, _ := c.Size(ctx)
	fmt.Println(n)

	// Output:
	// 10
}

func ExampleCache_Update() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Sets batch sets cache with key-value pairs by `data`, which is expired after `duration`.
	// It does not expire if `duration` == 0. It deletes the keys of `data` if `duration` < 0 or given `value` is nil.
	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2", "k3": "v3"}, 0)

	// Print the current key value pair
	k1, _ := c.Get(ctx, "k1")
	fmt.Println(k1)
	k2, _ := c.Get(ctx, "k2")
	fmt.Println(k2)
	k3, _ := c.Get(ctx, "k3")
	fmt.Println(k3)

	// Update updates the value of `key` without changing its expiration and returns the old value.
	re, exist, _ := c.Update(ctx, "k1", "v11")
	fmt.Println(re, exist)

	// The returned value `exist` is false if the `key` does not exist in the cache.
	// It does nothing if `key` does not exist in the cache.
	re1, exist1, _ := c.Update(ctx, "k4", "v44")
	fmt.Println(re1, exist1)

	kup1, _ := c.Get(ctx, "k1")
	fmt.Println(kup1)
	kup2, _ := c.Get(ctx, "k2")
	fmt.Println(kup2)
	kup3, _ := c.Get(ctx, "k3")
	fmt.Println(kup3)

	// Output:
	// v1
	// v2
	// v3
	// v1 true
	//  false
	// v11
	// v2
	// v3
}

func ExampleCache_UpdateExpire() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.Set(ctx, "k1", "v1", 1000*time.Millisecond)
	expire, _ := c.GetExpire(ctx, "k1")
	fmt.Println(expire)

	// UpdateExpire updates the expiration of `key` and returns the old expiration duration value.
	// It returns -1 and does nothing if the `key` does not exist in the cache.
	c.UpdateExpire(ctx, "k1", 500*time.Millisecond)

	expire1, _ := c.GetExpire(ctx, "k1")
	fmt.Println(expire1)

	// May Output:
	// 1s
	// 500ms
}

func ExampleCache_Values() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Write value
	c.Set(ctx, "k1", g.Map{"k1": "v1", "k2": "v2"}, 0)
	// c.Set(ctx, "k2", "Here is Value2", 0)
	// c.Set(ctx, "k3", 111, 0)

	// Values returns all values in the cache as slice.
	data, _ := c.Values(ctx)
	fmt.Println(data)

	// May Output:
	// [map[k1:v1 k2:v2]]
}

func ExampleCache_Close() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set Cache
	c.Set(ctx, "k1", "v", 0)
	data, _ := c.Get(ctx, "k1")
	fmt.Println(data)

	// Close closes the cache if necessary.
	c.Close(ctx)

	data1, _ := c.Get(ctx, "k1")

	fmt.Println(data1)

	// Output:
	// v
	// v
}

func ExampleCache_Contains() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set Cache
	c.Set(ctx, "k", "v", 0)

	// Contains returns true if `key` exists in the cache, or else returns false.
	// return true
	data, _ := c.Contains(ctx, "k")
	fmt.Println(data)

	// return false
	data1, _ := c.Contains(ctx, "k1")
	fmt.Println(data1)

	// Output:
	// true
	// false
}

func ExampleCache_Data() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1"}, 0)

	data, _ := c.Data(ctx)
	fmt.Println(data)

	// Set Cache
	c.Set(ctx, "k5", "v5", 0)
	data1, _ := c.Get(ctx, "k1")
	fmt.Println(data1)

	// Output:
	// map[k1:v1]
	// v1
}

func ExampleCache_Get() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set Cache Object
	c.Set(ctx, "k1", "v1", 0)

	// Get retrieves and returns the associated value of given `key`.
	// It returns nil if it does not exist, its value is nil or it's expired.
	data, _ := c.Get(ctx, "k1")
	fmt.Println(data)

	// Output:
	// v1
}

func ExampleCache_GetExpire() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set cache without expiration
	c.Set(ctx, "k", "v", 10000*time.Millisecond)

	// GetExpire retrieves and returns the expiration of `key` in the cache.
	// It returns 0 if the `key` does not expire. It returns -1 if the `key` does not exist in the cache.
	expire, _ := c.GetExpire(ctx, "k")
	fmt.Println(expire)

	// May Output:
	// 10s
}

func ExampleCache_GetOrSet() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// GetOrSet retrieves and returns the value of `key`, or sets `key`-`value` pair and returns `value`
	// if `key` does not exist in the cache.
	data, _ := c.GetOrSet(ctx, "k", "v", 10000*time.Millisecond)
	fmt.Println(data)

	data1, _ := c.Get(ctx, "k")
	fmt.Println(data1)

	// Output:
	// v
	// v

}

func ExampleCache_GetOrSetFunc() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// GetOrSetFunc retrieves and returns the value of `key`, or sets `key` with result of function `f`
	// and returns its result if `key` does not exist in the cache.
	c.GetOrSetFunc(ctx, "k1", func(ctx context.Context) (value interface{}, err error) {
		return "v1", nil
	}, 10000*time.Millisecond)
	v, _ := c.Get(ctx, "k1")
	fmt.Println(v)

	// If func returns nil, no action is taken
	c.GetOrSetFunc(ctx, "k2", func(ctx context.Context) (value interface{}, err error) {
		return nil, nil
	}, 10000*time.Millisecond)
	v1, _ := c.Get(ctx, "k2")
	fmt.Println(v1)

	// Output:
	// v1
}

func ExampleCache_GetOrSetFuncLock() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Modify locking Note that the function `f` should be executed within writing mutex lock for concurrent safety purpose.
	c.GetOrSetFuncLock(ctx, "k1", func(ctx context.Context) (value interface{}, err error) {
		return "v1", nil
	}, 0)
	v, _ := c.Get(ctx, "k1")
	fmt.Println(v)

	// Modification failed
	c.GetOrSetFuncLock(ctx, "k1", func(ctx context.Context) (value interface{}, err error) {
		return "update v1", nil
	}, 0)
	v, _ = c.Get(ctx, "k1")
	fmt.Println(v)

	c.Remove(ctx, g.Slice{"k1"}...)

	// Output:
	// v1
	// v1
}

func ExampleCache_Keys() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1"}, 0)

	// Print the current list of key values
	keys1, _ := c.Keys(ctx)
	fmt.Println(keys1)

	// Output:
	// [k1]
}

func ExampleCache_KeyStrings() {
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2"}, 0)

	// KeyStrings returns all keys in the cache as string slice.
	keys, _ := c.KeyStrings(ctx)
	fmt.Println(keys)

	// May Output:
	// [k1 k2]
}

func ExampleCache_Remove() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2"}, 0)

	// Remove deletes one or more keys from cache, and returns its value.
	// If multiple keys are given, it returns the value of the last deleted item.
	remove, _ := c.Remove(ctx, "k1")
	fmt.Println(remove)

	data, _ := c.Data(ctx)
	fmt.Println(data)

	// Output:
	// v1
	// map[k2:v2]
}

func ExampleCache_Removes() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2", "k3": "v3", "k4": "v4"}, 0)

	// Remove deletes one or more keys from cache, and returns its value.
	// If multiple keys are given, it returns the value of the last deleted item.
	c.Removes(ctx, g.Slice{"k1", "k2", "k3"})

	data, _ := c.Data(ctx)
	fmt.Println(data)

	// Output:
	// map[k4:v4]
}

func ExampleCache_Clear() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2", "k3": "v3", "k4": "v4"}, 0)

	// clears all data of the cache.
	c.Clear(ctx)

	data, _ := c.Data(ctx)
	fmt.Println(data)

	// Output:
	// map[]
}

func ExampleCache_MustGet() {
	// Intercepting panic exception information
	// err is empty, so panic is not performed
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover...:", r)
		}
	}()

	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set Cache Object
	c.Set(ctx, "k1", "v1", 0)

	// MustGet acts like Get, but it panics if any error occurs.
	k2 := c.MustGet(ctx, "k2")
	fmt.Println(k2)

	k1 := c.MustGet(ctx, "k1")
	fmt.Println(k1)

	// Output:
	// v1
}

func ExampleCache_MustGetOrSet() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// MustGetOrSet acts like GetOrSet, but it panics if any error occurs.
	k1 := c.MustGetOrSet(ctx, "k1", "v1", 0)
	fmt.Println(k1)

	k2 := c.MustGetOrSet(ctx, "k1", "v2", 0)
	fmt.Println(k2)

	// Output:
	// v1
	// v1
}

func ExampleCache_MustGetOrSetFunc() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// MustGetOrSetFunc acts like GetOrSetFunc, but it panics if any error occurs.
	c.MustGetOrSetFunc(ctx, "k1", func(ctx context.Context) (value interface{}, err error) {
		return "v1", nil
	}, 10000*time.Millisecond)
	v := c.MustGet(ctx, "k1")
	fmt.Println(v)

	c.MustGetOrSetFunc(ctx, "k2", func(ctx context.Context) (value interface{}, err error) {
		return nil, nil
	}, 10000*time.Millisecond)
	v1 := c.MustGet(ctx, "k2")
	fmt.Println(v1)

	// Output:
	// v1
	//
}

func ExampleCache_MustGetOrSetFuncLock() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// MustGetOrSetFuncLock acts like GetOrSetFuncLock, but it panics if any error occurs.
	c.MustGetOrSetFuncLock(ctx, "k1", func(ctx context.Context) (value interface{}, err error) {
		return "v1", nil
	}, 0)
	v := c.MustGet(ctx, "k1")
	fmt.Println(v)

	// Modification failed
	c.MustGetOrSetFuncLock(ctx, "k1", func(ctx context.Context) (value interface{}, err error) {
		return "update v1", nil
	}, 0)
	v = c.MustGet(ctx, "k1")
	fmt.Println(v)

	// Output:
	// v1
	// v1
}

func ExampleCache_MustContains() {

	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set Cache
	c.Set(ctx, "k", "v", 0)

	// MustContains returns true if `key` exists in the cache, or else returns false.
	// return true
	data := c.MustContains(ctx, "k")
	fmt.Println(data)

	// return false
	data1 := c.MustContains(ctx, "k1")
	fmt.Println(data1)

	// Output:
	// true
	// false
}

func ExampleCache_MustGetExpire() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Set cache without expiration
	c.Set(ctx, "k", "v", 10000*time.Millisecond)

	// MustGetExpire acts like GetExpire, but it panics if any error occurs.
	expire := c.MustGetExpire(ctx, "k")
	fmt.Println(expire)

	// May Output:
	// 10s
}

func ExampleCache_MustSize() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Add 10 elements without expiration
	for i := 0; i < 10; i++ {
		c.Set(ctx, i, i, 0)
	}

	// Size returns the number of items in the cache.
	n := c.MustSize(ctx)
	fmt.Println(n)

	// Output:
	// 10
}

func ExampleCache_MustData() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2"}, 0)

	data := c.MustData(ctx)
	fmt.Println(data)

	// May Output:
	// map[k1:v1 k2:v2]
}

func ExampleCache_MustKeys() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2"}, 0)

	// MustKeys acts like Keys, but it panics if any error occurs.
	keys1 := c.MustKeys(ctx)
	fmt.Println(keys1)

	// May Output:
	// [k1 k2]

}

func ExampleCache_MustKeyStrings() {
	c := gcache.New()

	c.SetMap(ctx, g.MapAnyAny{"k1": "v1", "k2": "v2"}, 0)

	// MustKeyStrings returns all keys in the cache as string slice.
	// MustKeyStrings acts like KeyStrings, but it panics if any error occurs.
	keys := c.MustKeyStrings(ctx)
	fmt.Println(keys)

	// May Output:
	// [k1 k2]
}

func ExampleCache_MustValues() {
	// Create a cache object,
	// Of course, you can also easily use the gcache package method directly
	c := gcache.New()

	// Write value
	c.Set(ctx, "k1", "v1", 0)

	// MustValues returns all values in the cache as slice.
	data := c.MustValues(ctx)
	fmt.Println(data)

	// Output:
	// [v1]
}

func ExampleCache_SetAdapter() {
	var (
		err         error
		ctx         = gctx.New()
		cache       = gcache.New()
		redisConfig = &gredis.Config{
			Address: "127.0.0.1:6379",
			Db:      9,
		}
		cacheKey   = `key`
		cacheValue = `value`
	)
	// Create redis client object.
	redis, err := gredis.New(redisConfig)
	if err != nil {
		panic(err)
	}
	// Create redis cache adapter and set it to cache object.
	cache.SetAdapter(gcache.NewAdapterRedis(redis))

	// Set and Get using cache object.
	err = cache.Set(ctx, cacheKey, cacheValue, time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(cache.MustGet(ctx, cacheKey).String())

	// Get using redis client.
	fmt.Println(redis.MustDo(ctx, "GET", cacheKey).String())

	// May Output:
	// value
	// value
}

func ExampleCache_GetAdapter() {
	var (
		err         error
		ctx         = gctx.New()
		cache       = gcache.New()
		redisConfig = &gredis.Config{
			Address: "127.0.0.1:6379",
			Db:      10,
		}
		cacheKey   = `key`
		cacheValue = `value`
	)
	redis, err := gredis.New(redisConfig)
	if err != nil {
		panic(err)
	}
	cache.SetAdapter(gcache.NewAdapterRedis(redis))

	// Set and Get using cache object.
	err = cache.Set(ctx, cacheKey, cacheValue, time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(cache.MustGet(ctx, cacheKey).String())

	// Get using redis client.
	v, err := cache.GetAdapter().(*gcache.AdapterRedis).Get(ctx, cacheKey)
	fmt.Println(err)
	fmt.Println(v.String())

	// May Output:
	// value
	// <nil>
	// value
}

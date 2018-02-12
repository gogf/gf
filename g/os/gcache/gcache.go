// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 单进程高速缓存.
package gcache

import (
    "sync"
    "time"
    "math"
    "sync/atomic"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gset"
)

// 缓存对象
type Cache struct {
    mu         sync.RWMutex
    data       *gmap.StringInterfaceMap // 存放真实的键值对数据
    eksets     *gmap.IntInterfaceMap    // 存放过期的键名数据
    closed     int32                    // 缓存对象是否关闭
}

// 缓存数据项
type CacheItem struct {
    v interface{} // 缓存键值
    e int64       // 过期时间
}

// 全局缓存管理对象
var cache *Cache = New()

// Cache对象按照缓存键名首字母做了分组
func New() *Cache {
    c := &Cache {
        data    : gmap.NewStringInterfaceMap(),
        eksets  : gmap.NewIntInterfaceMap(),
    }
    c.autoClearLoop()
    return c
}

// (使用全局KV缓存对象)设置kv缓存键值对，过期时间单位为毫秒
func Set(key string, value interface{}, expired int64)  {
    cache.Set(key, value, expired)
}

// (使用全局KV缓存对象)批量设置kv缓存键值对，过期时间单位为毫秒
func BatchSet(data map[string]interface{}, expire int64)  {
    cache.BatchSet(data, expire)
}

// (使用全局KV缓存对象)获取指定键名的值
func Get(key string) interface{} {
    return cache.Get(key)
}

// (使用全局KV缓存对象)删除指定键值对
func Remove(key string) {
    cache.Remove(key)
}

// (使用全局KV缓存对象)批量删除指定键值对
func BatchRemove(keys []string) {
    cache.BatchRemove(keys)
}

// 计算过期缓存的键名
func (c *Cache) getExpireKey(expire int64) int {
    return int(math.Ceil(float64(expire/1000)) + 1)*1000
}

// 获取或者创建一个过期键名存放Set
func (c *Cache) getOrNewExpireSet(ek int) *gset.StringSet {
    c.mu.RLock()
    defer c.mu.RUnlock()
    ekset := c.eksets.Get(ek)
    if ekset == nil {
        s := gset.NewStringSet()
        c.eksets.Set(ek, s)
        return s
    } else {
        return ekset.(*gset.StringSet)
    }
}


// 设置kv缓存键值对，过期时间单位为毫秒，expire<=0表示不过期
func (c *Cache) Set(key string, value interface{}, expire int64) {
    // 查找老的键值过期时间
    olde := int64(0)
    item := c.data.Get(key)
    if item != nil {
        olde = item.(CacheItem).e
    }
    // 保存新的键值对
    newe := int64(math.MaxInt64)
    if expire > 0 {
        newe = gtime.Millisecond() + int64(expire)
    }
    // 删除旧的过期键名
    if olde > 0 && newe != olde {
        c.eksets.Get(c.getExpireKey(olde)).(*gset.StringSet).Remove(key)
    }
    // 保存新的过期键名
    c.getOrNewExpireSet(c.getExpireKey(newe)).Add(key)
    // 最后才真实保存数据
    c.data.Set(key, CacheItem{v: value, e: newe})
}

// 批量设置
func (c *Cache) BatchSet(data map[string]interface{}, expire int64)  {
    l    := make([]string, len(data))
    m    := make(map[string]interface{})
    newe := int64(math.MaxInt64)
    if expire > 0 {
        newe = gtime.Millisecond() + int64(expire)
    }
    for k, v := range m {
        m[k]          = CacheItem{v: v, e: newe}
        l[len(l) - 1] = k
    }
    c.getOrNewExpireSet(c.getExpireKey(newe)).BatchAdd(l)
    c.data.BatchSet(m)
}

// 获取指定键名的值
func (c *Cache) Get(key string) interface{} {
    r := c.data.Get(key)
    if r != nil {
        item := r.(CacheItem)
        if item.e > gtime.Millisecond() {
            return item.v
        }
    }
    return nil
}

// 删除指定键值对
func (c *Cache) Remove(key string) {
    c.data.Remove(key)
}

// 批量删除键值对
func (c *Cache) BatchRemove(keys []string) {
    c.data.BatchRemove(keys)
}

// 获得所有的键名，组成字符串数组返回
func (c *Cache) Keys() []string {
    return c.data.Keys()
}

// 获得所有的值，组成数组返回
func (c *Cache) Values() []interface{} {
    return c.data.Values()
}

// 获得缓存对象的键值对数量
func (c *Cache) Size() int {
    return c.data.Size()
}

// 删除缓存对象
func (c *Cache) Close()  {
    atomic.AddInt32(&c.closed, 1)
}

// 自动清理过期键值对
// 每隔1秒清除过去3秒的键值对数据
func (c *Cache) autoClearLoop() {
    for atomic.LoadInt32(&c.closed) == 0 {
        ek  := c.getExpireKey(gtime.Millisecond())
        eks := []int{ek - 2000, ek - 3000, ek - 4000}
        for _, v := range eks {
            if r := c.eksets.Get(v); r != nil {
                c.data.BatchRemove(r.(*gset.StringSet).Slice())
            }
        }
        c.eksets.BatchRemove(eks)
        time.Sleep(time.Second)
    }
}

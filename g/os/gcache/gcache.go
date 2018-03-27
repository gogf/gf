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
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/container/gqueue"
)

// 缓存对象
type Cache struct {
    dmu        sync.RWMutex              // data锁
    emu        sync.RWMutex              // ekmap锁
    smu        sync.RWMutex              // eksets锁
    data       map[string]_CacheItem     // 缓存数据
    ekmap      map[string]int64          // 键名对应的分组过期时间
    eksets     map[int64]*gset.StringSet // 分组过期时间对应的键名列表
    eventQueue *gqueue.Queue             // 异步处理队列
    stopEvents chan struct{}             // 关闭时间通知
}

// 缓存数据项
type _CacheItem struct {
    v interface{} // 缓存键值
    e int64       // 过期时间
}

// 异步队列数据项
type _EventItem struct {
    k string      // 键名
    e int64       // 过期时间
}

// 全局缓存管理对象
var cache *Cache = New()

// Cache对象按照缓存键名首字母做了分组
func New() *Cache {
    c := &Cache {
        data       : make(map[string]_CacheItem),
        ekmap      : make(map[string]int64),
        eksets     : make(map[int64]*gset.StringSet),
        eventQueue : gqueue.New(),
        stopEvents : make(chan struct{}, 2),
    }
    go c.autoSyncLoop()
    go c.autoClearLoop()
    return c
}

// (使用全局KV缓存对象)设置kv缓存键值对，过期时间单位为毫秒
func Set(key string, value interface{}, expire int64)  {
    cache.Set(key, value, expire)
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

// 获得所有的键名，组成字符串数组返回
func Keys() []string {
    return cache.Keys()
}

// 获得所有的值，组成数组返回
func Values() []interface{} {
    return cache.Values()
}

// 获得缓存对象的键值对数量
func Size() int {
    return cache.Size()
}

// 计算过期缓存的键名
func (c *Cache) makeExpireKey(expire int64) int64 {
    return int64(math.Ceil(float64(expire/1000) + 1)*1000)
}

// 获取一个过期键名存放Set,如果没有则返回nil
func (c *Cache) getExpireSet(expire int64) *gset.StringSet {
    c.smu.RLock()
    if ekset, ok := c.eksets[expire]; ok {
        c.smu.RUnlock()
        return ekset
    }
    c.smu.RUnlock()
    return nil
}

// 获取或者创建一个过期键名存放Set(由于是异步单线程执行，因此不会出现创建set时的覆盖问题)
func (c *Cache) getOrNewExpireSet(expire int64) *gset.StringSet {
    if ekset := c.getExpireSet(expire); ekset == nil {
        set := gset.NewStringSet()
        c.smu.Lock()
        c.eksets[expire] = set
        c.smu.Unlock()
        return set
    } else {
        return ekset
    }
}

// 设置kv缓存键值对，过期时间单位为毫秒，expire<=0表示不过期
func (c *Cache) Set(key string, value interface{}, expire int64) {
    e := int64(math.MaxInt64)
    if expire > 0 {
        e = gtime.Millisecond() + int64(expire)
    }
    c.dmu.Lock()
    c.data[key] = _CacheItem{v: value, e: e}
    c.dmu.Unlock()
    c.eventQueue.PushBack(_EventItem{k: key, e:e})
}

// 批量设置
func (c *Cache) BatchSet(data map[string]interface{}, expire int64)  {
    e := int64(math.MaxInt64)
    if expire > 0 {
        e = gtime.Millisecond() + int64(expire)
    }
    for k, v := range data {
        c.dmu.Lock()
        c.data[k] = _CacheItem{v: v, e: e}
        c.dmu.Unlock()
        c.eventQueue.PushBack(_EventItem{k: k, e:e})
    }
}

// 获取指定键名的值
func (c *Cache) Get(key string) interface{} {
    c.dmu.RLock()
    item, ok := c.data[key]
    c.dmu.RUnlock()
    if ok {
        if item.e > gtime.Millisecond() {
            return item.v
        }
    }
    return nil
}

// 删除指定键值对
func (c *Cache) Remove(key string) {
    c.Set(key, nil, -1)
}

// 批量删除键值对
func (c *Cache) BatchRemove(keys []string) {
    for _, key := range keys {
        c.dmu.Lock()
        c.data[key] = _CacheItem{v: nil, e: -1}
        c.dmu.Unlock()
        c.eventQueue.PushBack(_EventItem{k: key, e: -1})
    }
}

// 获得所有的键名，组成字符串数组返回
func (c *Cache) Keys() []string {
    keys := make([]string, 0)
    c.dmu.RLock()
    for k, _ := range c.data {
        keys = append(keys, k)
    }
    c.dmu.RUnlock()
    return keys
}

// 获得所有的值，组成数组返回
func (c *Cache) Values() []interface{} {
    values := make([]interface{}, 0)
    c.dmu.RLock()
    for _, v := range c.data {
        values = append(values, v.v)
    }
    c.dmu.RUnlock()
    return values
}

// 获得缓存对象的键值对数量
func (c *Cache) Size() int {
    c.dmu.RLock()
    length := len(c.data)
    c.dmu.RUnlock()
    return length
}

// 删除缓存对象
func (c *Cache) Close()  {
    c.stopEvents <- struct{}{}
    c.eventQueue.Close()
}

// 数据自动同步循环
func (c *Cache) autoSyncLoop() {
    for {
        if r := c.eventQueue.PopFront(); r != nil {
            item := r.(_EventItem)
            newe := c.makeExpireKey(item.e)
            // 查询键名是否已经存在过期时间
            c.emu.RLock()
            olde, ok := c.ekmap[item.k];
            c.emu.RUnlock()
            // 是否需要删除旧的过期时间map中对应的键名
            if ok && newe != olde {
                if ekset := c.getExpireSet(olde); ekset != nil {
                    ekset.Remove(item.k)
                }
            }
            c.getOrNewExpireSet(newe).Add(item.k)
            // 重新设置对应键名的过期时间
            c.emu.Lock()
            c.ekmap[item.k] = newe
            c.emu.Unlock()
        } else {
            break
        }
    }
}

// 自动清理过期键值对
// 每隔1秒清除过去3秒的键值对数据
func (c *Cache) autoClearLoop() {
    for {
        select {
            case <- c.stopEvents:
                return
            default:
                ek  := c.makeExpireKey(gtime.Millisecond())
                eks := []int64{ek - 2000, ek - 3000, ek - 4000}
                for _, v := range eks {
                    if ekset := c.getExpireSet(v); ekset != nil {
                        ekset.Iterator(func(key string) {
                            // 删除缓存数据
                            c.dmu.Lock()
                            delete(c.data,  key)
                            c.dmu.Unlock()

                            // 删除异步处理数据项
                            c.emu.Lock()
                            delete(c.ekmap, key)
                            c.emu.Unlock()
                        })
                    }
                    // 删除异步处理键名set
                    c.smu.Lock()
                    delete(c.eksets, v)
                    c.smu.Unlock()
                }
                time.Sleep(time.Second)
        }
    }
}

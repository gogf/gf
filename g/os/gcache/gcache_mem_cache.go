// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcache

import (
    "gitee.com/johng/gf/g/container/glist"
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gtimer"
    "gitee.com/johng/gf/g/util/gconv"
    "math"
    "sync"
)


// 缓存对象
type memCache struct {
    dataMu       sync.RWMutex
    expireTimeMu sync.RWMutex
    expireSetMu  sync.RWMutex

    cap          int                            // 控制缓存池大小，超过大小则按照LRU算法进行缓存过期处理(默认为0表示不进行限制)
    data         map[interface{}]memCacheItem   // 缓存数据(所有的缓存数据存放哈希表)
    expireTimes  map[interface{}]int64          // 键名对应的分组过期时间(用于相同键名过期时间快速更新)，键值为1秒级时间戳
    expireSets   map[int64]*gset.Set            // 分组过期时间对应的键名列表(用于自动过期快速删除)，键值为1秒级时间戳

    lru          *memCacheLru                   // LRU缓存限制(只有限定cap池大小时才启用)
    lruGetList   *glist.List                    // Get操作的LRU记录
    eventList    *glist.List                    // 异步处理队列
    closed       *gtype.Bool                    // 关闭事件通知
}

// 缓存数据项
type memCacheItem struct {
    v interface{} // 键值
    e int64       // 过期时间
}

// 异步队列数据项
type memCacheEvent struct {
    k interface{} // 键名
    e int64       // 过期时间
}

const (
    // 当数据不过期时，默认设置的过期属性值，相当于：math.MaxInt64/1000000
    gDEFAULT_MAX_EXPIRE = 9223372036854
)

// 创建底层的缓存对象
func newMemCache(lruCap...int) *memCache {
    c := &memCache {
        lruGetList  : glist.New(),
        data        : make(map[interface{}]memCacheItem),
        expireTimes : make(map[interface{}]int64),
        expireSets  : make(map[int64]*gset.Set),
        eventList   : glist.New(),
        closed      : gtype.NewBool(),
    }
    if len(lruCap) > 0 {
        c.cap = lruCap[0]
        c.lru = newMemCacheLru(c)
    }
    return c
}

// 计算过期缓存的键名(将毫秒换算成秒的整数毫秒，按照1秒进行分组)
func (c *memCache) makeExpireKey(expire int64) int64 {
    return int64(math.Ceil(float64(expire/1000) + 1)*1000)
}

// 获取一个过期键名存放Set, 如果没有则返回nil
func (c *memCache) getExpireSet(expire int64) (expireSet *gset.Set) {
    c.expireSetMu.RLock()
    expireSet, _ = c.expireSets[expire]
    c.expireSetMu.RUnlock()
    return
}

// 获取或者创建一个过期键名存放Set(由于是异步单线程执行，因此不会出现创建set时的覆盖问题)
func (c *memCache) getOrNewExpireSet(expire int64) (expireSet *gset.Set) {
    if expireSet = c.getExpireSet(expire); expireSet == nil {
        expireSet = gset.New()
        c.expireSetMu.Lock()
        // 写锁二次检索确认
        if es, ok := c.expireSets[expire]; ok {
            expireSet = es
        } else {
            c.expireSets[expire] = expireSet
        }
        c.expireSetMu.Unlock()
    }
    return
}

// 设置kv缓存键值对，过期时间单位为毫秒，expire<=0表示不过期
func (c *memCache) Set(key interface{}, value interface{}, expire int) {
    expireTime := c.getInternalExpire(expire)
    c.dataMu.Lock()
    c.data[key] = memCacheItem{v : value, e : expireTime}
    c.dataMu.Unlock()
    c.eventList.PushBack(&memCacheEvent{k : key, e : expireTime})
}

// 设置kv缓存键值对，内部会对键名的存在性使用写锁进行二次检索确认，如果存在则不再写入；返回键名对应的键值。
// 在高并发下有用，防止数据写入的并发逻辑错误。
func (c *memCache) doSetWithLockCheck(key interface{}, value interface{}, expire int) interface{} {
    expireTimestamp := c.getInternalExpire(expire)
    c.dataMu.Lock()
    if v, ok := c.data[key]; ok && !v.IsExpired() {
        c.dataMu.Unlock()
        return v.v
    }
    if f, ok := value.(func() interface {}); ok {
        value = f()
    }
    c.data[key] = memCacheItem{v : value, e : expireTimestamp}
    c.dataMu.Unlock()
    c.eventList.PushBack(&memCacheEvent{k : key, e : expireTimestamp})
    return value
}

// 根据给定expire参数计算内部使用的expire过期时间
func (c *memCache) getInternalExpire(expire int) int64 {
    if expire != 0 {
        return gtime.Millisecond() + int64(expire)
    } else {
        return gDEFAULT_MAX_EXPIRE
    }
}

// 当键名不存在时写入，并返回true；否则返回false。
func (c *memCache) SetIfNotExist(key interface{}, value interface{}, expire int) bool {
    if !c.Contains(key) {
        c.doSetWithLockCheck(key, value, expire)
        return true
    }
    return false
}

// 批量设置
func (c *memCache) BatchSet(data map[interface{}]interface{}, expire int) {
    expireTime := c.getInternalExpire(expire)
    for k, v := range data {
        c.dataMu.Lock()
        c.data[k] = memCacheItem{v: v, e: expireTime}
        c.dataMu.Unlock()
        c.eventList.PushBack(&memCacheEvent{k: k, e: expireTime})
    }
}

// 获取指定键名的值
func (c *memCache) Get(key interface{}) interface{} {
    c.dataMu.RLock()
    item, ok := c.data[key]
    c.dataMu.RUnlock()
    if ok && !item.IsExpired() {
        // 增加LRU(Least Recently Used)操作记录
        if c.cap > 0 {
            c.lruGetList.PushBack(key)
        }
        return item.v
    }
    return nil
}

// 当键名存在时返回其键值，否则写入指定的键值
func (c *memCache) GetOrSet(key interface{}, value interface{}, expire int) interface{} {
    if v := c.Get(key); v == nil {
        return c.doSetWithLockCheck(key, value, expire)
    } else {
        return v
    }
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func (c *memCache) GetOrSetFunc(key interface{}, f func() interface{}, expire int) interface{} {
    if v := c.Get(key); v == nil {
        // 可能存在多个goroutine被阻塞在这里，f可能是并发运行
        return c.doSetWithLockCheck(key, f(), expire)
    } else {
        return v
    }
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func (c *memCache) GetOrSetFuncLock(key interface{}, f func() interface{}, expire int) interface{} {
    if v := c.Get(key); v == nil {
        return c.doSetWithLockCheck(key, f, expire)
    } else {
        return v
    }
}

// 是否存在指定的键名，true表示存在，false表示不存在。
func (c *memCache) Contains(key interface{}) bool {
    return c.Get(key) != nil
}

// 删除指定键值对，并返回被删除的键值
func (c *memCache) Remove(key interface{}) (value interface{}) {
    c.dataMu.RLock()
    item, ok := c.data[key]
    c.dataMu.RUnlock()
    if ok {
        value = item.v
        c.dataMu.Lock()
        delete(c.data, key)
        c.dataMu.Unlock()
        c.eventList.PushBack(&memCacheEvent{k: key, e: gtime.Millisecond() - 1000})
    }
    return
}

// 批量删除键值对，并返回被删除的键值对数据
func (c *memCache) BatchRemove(keys []interface{}) {
    for _, key := range keys {
        c.Remove(key)
    }
}

// 返回缓存的所有数据键值对(不包含已过期数据)
func (c *memCache) Data() map[interface{}]interface{} {
    m := make(map[interface{}]interface{})
    c.dataMu.RLock()
    for k, v := range c.data {
        if !v.IsExpired() {
            m[k] = v.v
        }
    }
    c.dataMu.RUnlock()
    return m
}

// 获得所有的键名，组成数组返回
func (c *memCache) Keys() []interface{} {
    keys := make([]interface{}, 0)
    c.dataMu.RLock()
    for k, v := range c.data {
        if !v.IsExpired() {
            keys = append(keys, k)
        }
    }
    c.dataMu.RUnlock()
    return keys
}

// 获得所有的键名，组成字符串数组返回
func (c *memCache) KeyStrings() []string {
    return gconv.Strings(c.Keys())
}

// 获得所有的值，组成数组返回
func (c *memCache) Values() []interface{} {
    values := make([]interface{}, 0)
    c.dataMu.RLock()
    for _, v := range c.data {
        if !v.IsExpired() {
            values = append(values, v.v)
        }
    }
    c.dataMu.RUnlock()
    return values
}

// 获得缓存对象的键值对数量
func (c *memCache) Size() (size int) {
    c.dataMu.RLock()
    size = len(c.data)
    c.dataMu.RUnlock()
    return
}

// 删除缓存对象
func (c *memCache) Close()  {
    if c.cap > 0 {
        c.lru.Close()
    }
    c.closed.Set(true)
}

// 数据异步任务循环:
// 1、将事件列表中的数据异步处理，并同步结果到expireTimes和expireSets属性中；
// 2、清理过期键值对数据；
func (c *memCache) syncEventAndClearExpired() {
    event         := (*memCacheEvent)(nil)
    oldExpireTime := int64(0)
    newExpireTime := int64(0)
    if c.closed.Val() {
        gtimer.Exit()
        return
    }
    // ========================
    // 数据同步处理
    // ========================
    for {
        v := c.eventList.PopFront()
        if v == nil {
            break
        }
        event = v.(*memCacheEvent)
        // 获得旧的过期时间分组
        c.expireTimeMu.RLock()
        oldExpireTime = c.expireTimes[event.k]
        c.expireTimeMu.RUnlock()
        // 计算新的过期时间分组
        newExpireTime = c.makeExpireKey(event.e)
        if newExpireTime != oldExpireTime {
            c.getOrNewExpireSet(newExpireTime).Add(event.k)
            if oldExpireTime != 0 {
                c.getOrNewExpireSet(oldExpireTime).Remove(event.k)
            }
            // 重新设置对应键名的过期时间
            c.expireTimeMu.Lock()
            c.expireTimes[event.k] = newExpireTime
            c.expireTimeMu.Unlock()
        }
        // 写入操作也会增加到LRU(Least Recently Used)操作记录
        if c.cap > 0 {
            c.lru.Push(event.k)
        }
    }
    // 异步处理读取操作的LRU列表
    if c.cap > 0 && c.lruGetList.Len() > 0 {
        for {
            if v := c.lruGetList.PopFront(); v != nil {
                c.lru.Push(v)
            } else {
                break
            }
        }
    }
    // ========================
    // 缓存过期处理
    // ========================
    ek := c.makeExpireKey(gtime.Millisecond())
    eks := []int64{ek - 1000, ek - 2000, ek - 3000, ek - 4000, ek - 5000}
    for _, expireTime := range eks {
        if expireSet := c.getExpireSet(expireTime); expireSet != nil {
            // 遍历Set，执行数据过期删除
            expireSet.Iterator(func(key interface{}) bool {
                c.clearByKey(key)
                return true
            })
            // Set数据处理完之后删除该Set
            c.expireSetMu.Lock()
            delete(c.expireSets, expireTime)
            c.expireSetMu.Unlock()
        }
    }
}

// 删除对应键名的缓存数据
func (c *memCache) clearByKey(key interface{}, force...bool) {
    // 删除缓存数据
    c.dataMu.Lock()
    // 删除核对，真正的过期才删除
    if item, ok := c.data[key]; (ok && item.IsExpired()) || (len(force) > 0 && force[0]) {
        delete(c.data, key)
    }
    c.dataMu.Unlock()

    // 删除异步处理数据项
    c.expireTimeMu.Lock()
    delete(c.expireTimes, key)
    c.expireTimeMu.Unlock()

    // 删除LRU管理对象中指定键名
    if c.cap > 0 {
        c.lru.Remove(key)
    }
}

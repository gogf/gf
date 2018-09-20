// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcache

import (
    "math"
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/os/gtime"
    "sync"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/container/glist"
    "time"
)

// 缓存对象
type memCache struct {
    dmu        sync.RWMutex                   // data锁(自定义锁的目的是除去键值的断言转换造成的性能损耗)
    emu        sync.RWMutex                   // ekmap锁(expire key map)
    smu        sync.RWMutex                   // eksets锁(expire key sets)
    lru        *memCacheLru                   // LRU缓存限制(只有限定池大小时才启用)
    cap        int                            // 控制缓存池大小，超过大小则按照LRU算法进行缓存过期处理(默认为0表示不进行限制)
    data       map[interface{}]memCacheItem   // 缓存数据(所有的缓存数据存放哈希表)
    ekmap      map[interface{}]int64          // 键名对应的分组过期时间(用于相同键名过期时间快速更新)，键值为10秒级时间戳
    eksets     map[int64]*gset.Set            // 分组过期时间对应的键名列表(用于自动过期快速删除)，键值为10秒级时间戳
    eventList  *glist.List                    // 异步处理队列
    lruGetList *glist.List                    // 获取方法的LRU列表
    stopChan   chan struct{}                  // 关闭时间通知
}

// 缓存数据项
type memCacheItem struct {
    v interface{} // 缓存键值
    e int64       // 过期时间
}

// 异步队列数据项
type memCacheEvent struct {
    k interface{} // 键名
    e int64       // 过期时间
}

// 创建底层的缓存对象
func newMemCache(lruCap...int) *memCache {
    c := &memCache {
        lru        : newMemCacheLru(),
        data       : make(map[interface{}]memCacheItem),
        ekmap      : make(map[interface{}]int64),
        eksets     : make(map[int64]*gset.Set),
        stopChan   : make(chan struct{}),
        eventList  : glist.New(),
        lruGetList : glist.New(),
    }
    if len(lruCap) > 0 {
       c.cap = lruCap[0]
    }
    return c
}

// 计算过期缓存的键名(将毫秒换算成秒的整数毫秒)
func (c *memCache) makeExpireKey(expire int64) int64 {
    return int64(math.Ceil(float64(expire/10000) + 1)*10000)
}

// 获取一个过期键名存放Set,如果没有则返回nil
func (c *memCache) getExpireSet(expire int64) *gset.Set {
    c.smu.RLock()
    if ekset, ok := c.eksets[expire]; ok {
        c.smu.RUnlock()
        return ekset
    }
    c.smu.RUnlock()
    return nil
}

// 获取或者创建一个过期键名存放Set(由于是异步单线程执行，因此不会出现创建set时的覆盖问题)
func (c *memCache) getOrNewExpireSet(expire int64) *gset.Set {
    if ekset := c.getExpireSet(expire); ekset == nil {
        set := gset.New()
        c.smu.Lock()
        // 二次检索确认
        if ekset, ok := c.eksets[expire]; !ok {
            c.eksets[expire] = set
        } else {
            set = ekset
        }
        c.smu.Unlock()
        return set
    } else {
        return ekset
    }
}

// 设置kv缓存键值对，过期时间单位为毫秒，expire<=0表示不过期
func (c *memCache) Set(key interface{}, value interface{}, expire int) {
    expireTimestamp := c.getInternalExpire(expire)
    c.dmu.Lock()
    c.data[key] = memCacheItem{v : value, e : expireTimestamp}
    c.dmu.Unlock()
    c.eventList.PushBack(memCacheEvent{k : key, e : expireTimestamp})
}

// 设置kv缓存键值对，内部会对键名的存在性使用写锁进行二次检索确认，如果存在则不再写入；返回键名对应的键值。
// 在高并发下有用，防止数据写入的并发逻辑错误。
func (c *memCache) doSetWithLockCheck(key interface{}, value interface{}, expire int) interface{} {
    expireTimestamp := c.getInternalExpire(expire)
    c.dmu.Lock()
    if v, ok := c.data[key]; ok && !v.IsExpired() {
        c.dmu.Unlock()
        return v
    }
    if f, ok := value.(func() interface {}); ok {
        value = f()
    }
    c.data[key] = memCacheItem{v : value, e : expireTimestamp}
    c.dmu.Unlock()
    c.eventList.PushBack(memCacheEvent{k : key, e : expireTimestamp})
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
func (c *memCache) BatchSet(data map[interface{}]interface{}, expire int)  {
    expireTimestamp := c.getInternalExpire(expire)
    for k, v := range data {
        c.dmu.Lock()
        c.data[k] = memCacheItem{v: v, e: expireTimestamp}
        c.dmu.Unlock()
        c.eventList.PushBack(memCacheEvent{k: k, e: expireTimestamp})
    }
}

// 获取指定键名的值
func (c *memCache) Get(key interface{}) interface{} {
    c.dmu.RLock()
    item, ok := c.data[key]
    c.dmu.RUnlock()
    if ok && !item.IsExpired() {
        // LRU(Least Recently Used)操作记录
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
        v = f()
        c.doSetWithLockCheck(key, v, expire)
        return v
    } else {
        return v
    }
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func (c *memCache) GetOrSetFuncLock(key interface{}, f func() interface{}, expire int) interface{} {
    if v := c.Get(key); v == nil {
        c.doSetWithLockCheck(key, f, expire)
        return v
    } else {
        return v
    }
}

// 是否存在指定的键名，true表示存在，false表示不存在。
func (c *memCache) Contains(key interface{}) bool {
    return c.Get(key) != nil
}

// 删除指定键值对，并返回被删除的键值
func (c *memCache) Remove(key interface{}) interface{} {
    c.dmu.Lock()
    item, ok := c.data[key]
    if ok {
        delete(c.data, key)
    }
    c.dmu.Unlock()
    return item.v
}

// 批量删除键值对，并返回被删除的键值对数据
func (c *memCache) BatchRemove(keys []interface{}) {
    c.dmu.Lock()
    for _, key := range keys {
        delete(c.data, key)
    }
    c.dmu.Unlock()
}

// 获得所有的键名，组成数组返回
func (c *memCache) Keys() []interface{} {
    keys := make([]interface{}, 0)
    c.dmu.RLock()
    for k, v := range c.data {
        if !v.IsExpired() {
            keys = append(keys, k)
        }
    }
    c.dmu.RUnlock()
    return keys
}

// 获得所有的键名，组成字符串数组返回
func (c *memCache) KeyStrings() []string {
    return gconv.Strings(c.Keys())
}

// 获得所有的值，组成数组返回
func (c *memCache) Values() []interface{} {
    values := make([]interface{}, 0)
    c.dmu.RLock()
    for _, v := range c.data {
        if !v.IsExpired() {
            values = append(values, v.v)
        }
    }
    c.dmu.RUnlock()
    return values
}

// 获得缓存对象的键值对数量
func (c *memCache) Size() int {
    c.dmu.RLock()
    length := len(c.data)
    c.dmu.RUnlock()
    return length
}

// 删除缓存对象
func (c *memCache) Close()  {
    close(c.stopChan)
    c.lru.Close()
}

// 数据自动同步循环
func (c *memCache) autoSyncLoop() {
    newe := int64(0)
    for {
        select {
            case <-c.stopChan:
                return
            default:
                for {
                    v := c.eventList.PopFront()
                    if v == nil {
                        break
                    }
                    item := v.(memCacheEvent)
                    nowm := gtime.Millisecond()
                    // 如果用户设置的时间比当前时间还小，那么表示要自动清除了，
                    // 这里赋值一个当前时间-10秒的时间，在自动清理的goroutine中会自动检测删除该key
                    if item.e < nowm {
                        newe = c.makeExpireKey(nowm) - 10000
                    } else {
                        newe = c.makeExpireKey(item.e)
                    }
                    // 添加该key到对应的过期集合中
                    // 注意：这里不需要检查存在性，
                    // 因为在key过期的时候，会和原始的键值对中的过期时间做核对。
                    c.getOrNewExpireSet(newe).Add(item.k)
                    // 重新设置对应键名的过期时间
                    c.emu.Lock()
                    c.ekmap[item.k] = newe
                    c.emu.Unlock()
                    // LRU(Least Recently Used)操作记录
                    if c.cap > 0 {
                        c.lru.Push(item.k)
                    }
                }
                if c.cap > 0 {
                    // 优先级高的lru key放后面，读取列表
                    for {
                        if v := c.lruGetList.PopFront(); v != nil {
                            c.lru.Push(v)
                        } else {
                            break
                        }
                    }
                }
                time.Sleep(10 * time.Second)
        }
    }
}

// LRU缓存淘汰处理+自动清理过期键值对
// 每隔10秒清除过去60秒的键值对数据
func (c *memCache) autoClearLoop() {
   for {
       select {
           case <- c.stopChan:
               return
           default:
               // 缓存过期处理
               ek  := c.makeExpireKey(gtime.Millisecond())
               eks := []int64{ek - 10000, ek - 20000, ek - 30000, ek - 40000, ek - 50000, ek - 60000}
               for _, v := range eks {
                   if ekset := c.getExpireSet(v); ekset != nil {
                       ekset.Iterator(func(v interface{}) bool {
                           return c.clearByKey(v)
                       })
                   }
                   // 数据处理完之后从集合中删除该时间段
                   c.smu.Lock()
                   delete(c.eksets, v)
                   c.smu.Unlock()
               }
               // LRU缓存淘汰清理
               if c.cap > 0 {
                   for i := c.Size() - c.cap; i > 0; i-- {
                       if s := c.lru.Pop(); s != nil {
                           c.clearByKey(s, true)
                       }
                   }
               }
               time.Sleep(10*time.Second)
       }
   }
}

// 删除对应键名的缓存数据
func (c *memCache) clearByKey(key interface{}, force...bool) bool {
    // 删除缓存数据
    c.dmu.Lock()
    // 删除核对，真正的过期才删除
    if item, ok := c.data[key]; (ok && item.IsExpired()) || (len(force) > 0 && force[0]) {
        delete(c.data, key)
    }
    c.dmu.Unlock()

    // 删除异步处理数据项
    c.emu.Lock()
    delete(c.ekmap, key)
    c.emu.Unlock()

    // 删除LRU管理对象中指定键名
    c.lru.Remove(key)

    return true
}

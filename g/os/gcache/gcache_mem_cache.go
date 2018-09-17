// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gcache

import (
    "time"
    "math"
    "gitee.com/johng/gf/g/container/gset"
    "gitee.com/johng/gf/g/os/gtime"
    "sync"
    "gitee.com/johng/gf/g/container/gtype"
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
)

// 缓存对象
type memCache struct {
    dmu        sync.RWMutex                   // data锁(自定义锁的目的是除去键值的断言转换造成的性能损耗)
    emu        sync.RWMutex                   // ekmap锁(expire key map)
    smu        sync.RWMutex                   // eksets锁(expire key sets)
    lru        *memCacheLru                   // LRU缓存限制(只有限定池大小时才启用)
    cap        *gtype.Int                     // 控制缓存池大小，超过大小则按照LRU算法进行缓存过期处理(默认为0表示不进行限制)
    data       map[interface{}]memCacheItem   // 缓存数据(所有的缓存数据存放哈希表)
    ekmap      map[interface{}]int64          // 键名对应的分组过期时间(用于相同键名过期时间快速更新)，键值为10秒级时间戳
    eksets     map[int64]*gset.Set            // 分组过期时间对应的键名列表(用于自动过期快速删除)，键值为10秒级时间戳
    eventChan  chan memCacheEvent             // 异步处理队列
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

const (
    // 这个数值不能太大，否则初始化会占用太多无意义的内存
    // 60W，这个数值是创始人的机器上支持基准测试的参考结果
    gEVENT_QUEUE_SIZE = 600000
)

// 创建底层的缓存对象
func newMemCache() *memCache {
    c := &memCache {
        lru        : newMemCacheLru(),
        cap        : gtype.NewInt(),
        data       : make(map[interface{}]memCacheItem),
        ekmap      : make(map[interface{}]int64),
        eksets     : make(map[int64]*gset.Set),
        stopChan   : make(chan struct{}),
        eventChan  : make(chan memCacheEvent, gEVENT_QUEUE_SIZE),
    }
    go c.autoSyncLoop()
    go c.autoClearLoop()
    return c
}

// 设置缓存池大小，内部依靠LRU算法进行缓存淘汰处理
func (c *memCache) SetCap(cap int) {
    c.cap.Set(cap)
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
        c.eksets[expire] = set
        c.smu.Unlock()
        return set
    } else {
        return ekset
    }
}

// 设置kv缓存键值对，过期时间单位为毫秒，expire<=0表示不过期
func (c *memCache) Set(key interface{}, value interface{}, expire int) {
    var e int64
    if expire != 0 {
        e = gtime.Millisecond() + int64(expire)
    } else {
        e = gDEFAULT_MAX_EXPIRE
    }
    c.dmu.Lock()
    c.data[key] = memCacheItem{v : value, e : e}
    c.dmu.Unlock()
    c.eventChan <- memCacheEvent{k : key, e : e}
}

// 当键名不存在时写入，并返回true；否则返回false。
// 常用来做对并发性要求不高的内存锁。
func (c *memCache) SetIfNotExist(key interface{}, value interface{}, expire int) bool {
    if !c.Contains(key) {
        c.Set(key, value, expire)
        return true
    }
    return false
}

// 批量设置
func (c *memCache) BatchSet(data map[interface{}]interface{}, expire int)  {
    var e int64
    if expire != 0 {
        e = gtime.Millisecond() + int64(expire)
    } else {
        e = gDEFAULT_MAX_EXPIRE
    }
    for k, v := range data {
        c.dmu.Lock()
        c.data[k] = memCacheItem{v: v, e: e}
        c.dmu.Unlock()
        c.eventChan <- memCacheEvent{k: k, e:e}
    }
}

// 获取指定键名的值
func (c *memCache) Get(key interface{}) interface{} {
    c.dmu.RLock()
    item, ok := c.data[key]
    c.dmu.RUnlock()
    if ok && !item.IsExpired() {
        return item.v
    }
    return nil
}

// 当键名存在时返回其键值，否则写入指定的键值
func (c *memCache) GetOrSet(key interface{}, value interface{}, expire int) interface{} {
    if v := c.Get(key); v == nil {
        c.Set(key, value, expire)
        return value
    } else {
        return v
    }
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func (c *memCache) GetOrSetFunc(key interface{}, f func() interface{}, expire int) interface{} {
    if v := c.Get(key); v == nil {
        v = f()
        c.Set(key, v, expire)
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
func (c *memCache) BatchRemove(keys []interface{}) map[interface{}]interface{} {
    m := make(map[interface{}]interface{})
    for _, key := range keys {
        if v := c.Remove(key); v != nil {
            m[key] = v
        }
    }
    return m
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
    close(c.eventChan)
    c.lru.Close()
}

// 数据自动同步循环
func (c *memCache) autoSyncLoop() {
    for {
        if len(c.eventChan) > gEVENT_QUEUE_SIZE - 1000 {
            fmt.Println("full")
        }
        item := <- c.eventChan
        if item.k == nil {
            break
        }
        // 添加该key到对应的过期集合中
        // 注意：这里不需要检查存在性，
        // 因为在key过期的时候，会和原始的键值对中的过期时间做核对
        newe := c.makeExpireKey(item.e)
        c.getOrNewExpireSet(newe).Add(item.k)
        // 重新设置对应键名的过期时间
        c.emu.Lock()
        c.ekmap[item.k] = newe
        c.emu.Unlock()
        // LRU操作记录(只有新增和修改操作才会记录到LRU管理对象中，删除不会)
        if c.cap.Val() > 0 {
            c.lru.Push(item.k)
        }
    }
}

// LRU缓存淘汰处理+自动清理过期键值对
// 每隔10秒清除过去30秒的键值对数据
func (c *memCache) autoClearLoop() {
    for {
        select {
            case <- c.stopChan:
                return
            default:
                // 缓存过期处理
                ek  := c.makeExpireKey(gtime.Millisecond())
                eks := []int64{ek - 10000, ek - 20000, ek - 30000}
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
                // LRU缓存淘汰处理
                if c.cap.Val() > 0 {
                    for i := c.Size() - c.cap.Val(); i > 0; i-- {
                        if s := c.lru.Pop(); s != "" {
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

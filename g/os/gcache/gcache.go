package gcache

import (
    "sync"
    "gf/g/util/gtime"
    "time"
    "gf/g/encoding/ghash"
)

const (
    gCACHE_GROUP_SIZE = 4 // 缓存分区大小，不能超过uint8的最大值
)

type Cache struct {
    sync.RWMutex
    m map[uint8]*CacheMap     // 以分区大小数字作为索引
}

type CacheMap  struct {
    sync.RWMutex
    deleted bool             // 对象是否已删除，以便判断停止goroutine
    m map[string]CacheItem   // 键值对
}

type CacheItem struct {
    v interface{} // 缓存键值
    e int64       // 过期时间
}

// 全局缓存管理对象
var cache *Cache = New()

// Cache对象按照缓存键名首字母做了分组
func New() *Cache {
    c := &Cache {
        m : make(map[uint8]*CacheMap),
    }
    var i uint8 = 0
    for ; i < gCACHE_GROUP_SIZE; i++ {
        m := &CacheMap {
            m : make(map[string]CacheItem),
        }
        c.m[i] = m
        go m.autoClearLoop()
    }
    return c
}

// (使用全局KV缓存对象)设置kv缓存键值对，过期时间单位为毫秒
func Set(k string, v interface{}, expired int64)  {
    cache.Set(k, v, expired)
}

// (使用全局KV缓存对象)批量设置kv缓存键值对，过期时间单位为毫秒
func BatchSet(m map[string]interface{}, expired int64)  {
    cache.BatchSet(m, expired)
}

// (使用全局KV缓存对象)获取指定键名的值
func Get(k string) interface{} {
    return cache.Get(k)
}

// (使用全局KV缓存对象)删除指定键值对
func Remove(k string) {
    cache.Remove(k)
}

// (使用全局KV缓存对象)批量删除指定键值对
func BatchRemove(l []string) {
    cache.BatchRemove(l)
}

// 设置kv缓存键值对，过期时间单位为毫秒
func (c *Cache) Set(k string, v interface{}, expired int64)  {
    c.RLock()
    c.m[c.getIndex(k)].Set(k, v, expired)
    c.RUnlock()
}

// 批量设置
func (c *Cache) BatchSet(m map[string]interface{}, expired int64)  {
    c.RLock()
    for k, v := range m {
        c.m[c.getIndex(k)].Set(k, v, expired)
    }
    c.RUnlock()
}

// 获取指定键名的值
func (c *Cache) Get(k string) interface{} {
    c.RLock()
    r := c.m[c.getIndex(k)].Get(k)
    c.RUnlock()
    return r
}

// 删除指定键值对
func (c *Cache) Remove(k string) {
    c.RLock()
    c.m[c.getIndex(k)].Remove(k)
    c.RUnlock()
}

// 批量删除键值对
func (c *Cache) BatchRemove(l []string) {
    c.RLock()
    for _, k := range l {
        c.m[c.getIndex(k)].Remove(k)
    }
    c.RUnlock()
}

// 获得所有的键名，组成字符串数组返回
func (c *Cache) Keys() []string {
    l := make([]string, 0)
    c.RLock()
    for _, cm := range c.m {
        cm.RLock()
        for k2, _ := range cm.m {
            l = append(l, k2)
        }
        cm.RUnlock()
    }
    c.RUnlock()
    return l
}

// 获得所有的值，组成数组返回
func (c *Cache) Values() []interface{} {
    l := make([]interface{}, 0)
    c.RLock()
    for _, cm := range c.m {
        cm.RLock()
        for _, v2 := range cm.m {
            l = append(l, v2.v)
        }
        cm.RUnlock()
    }
    c.RUnlock()
    return l
}

// 获得缓存对象的键值对数量
func (c *Cache) Size() int {
    var size int
    c.RLock()
    for _, cm := range c.m {
        cm.RLock()
        size += len(cm.m)
        cm.RUnlock()
    }
    c.RUnlock()
    return size
}

// 删除缓存对象
func (c *Cache) Close()  {
    c.RLock()
    for _, cm := range c.m {
        cm.Lock()
        cm.deleted = true
        cm.Unlock()
    }
    c.RUnlock()

    c.Lock()
    c.m = nil
    c.Unlock()
}

// 计算缓存的索引
func (c *Cache) getIndex(k string) uint8 {
    return uint8(ghash.BKDRHash([]byte(k)) % gCACHE_GROUP_SIZE)
}

// 设置kv缓存键值对，过期时间单位为毫秒
func (cm *CacheMap) Set(k string, v interface{}, expired int64)  {
    var e int64
    if expired > 0 {
        e = gtime.Millisecond() + int64(expired)
    }
    cm.Lock()
    cm.m[k] = CacheItem{v: v, e: e}
    cm.Unlock()
}

// 获取指定键名的值
func (cm *CacheMap) Get(k string) interface{} {
    var v interface{}
    cm.RLock()
    if r, ok := cm.m[k]; ok {
        if r.e > 0 && r.e < gtime.Millisecond() {
            v = nil
        } else {
            v = r.v
        }
    }
    cm.RUnlock()
    return v
}

// 删除指定键值对
func (cm *CacheMap) Remove(k string) {
    cm.Lock()
    delete(cm.m, k)
    cm.Unlock()
}

// 自动清理过期键值对(每间隔60秒执行)
func (cm *CacheMap) autoClearLoop() {
    for !cm.deleted {
        expired := make([]string, 0)
        cm.RLock()
        for k, v := range cm.m {
            if v.e > 0 && v.e < gtime.Millisecond() {
                expired = append(expired, k)
            }
        }
        cm.RUnlock()

        if len(expired) > 0 {
            cm.Lock()
            for _, k := range expired {
                delete(cm.m, k)
            }
            cm.Unlock()
        }
        time.Sleep(60 * time.Second)
    }
}



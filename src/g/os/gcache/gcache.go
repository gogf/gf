package gcache

import (
    "sync"
    "g/util/gtime"
    "time"
    "g/encoding/gmd5"
    "strings"
)

type Cache struct {
    sync.RWMutex
    m map[string]*CacheMap
}

type CacheMap  struct {
    sync.RWMutex
    m map[string]*CacheItem
}

type CacheItem struct {
    Value   interface{}
    Expired int64
}

var cache *Cache = New()

// Cache对象按照缓存键名首字母做了分组
func New() *Cache {
    c := &Cache{
        m : make(map[string]*CacheMap),
    }
    // 0 - 9
    for i := 48; i <= 57; i++ {
        m := &CacheMap {
            m : make(map[string]*CacheItem),
        }
        c.m[string(i)] = m
        go m.autoClearLoop()
    }
    // a - z
    for i := 97; i <= 122; i++ {
        m := &CacheMap {
            m : make(map[string]*CacheItem),
        }
        c.m[string(i)] = m
        go m.autoClearLoop()
    }
    return c
}

// (使用全局KV缓存对象)设置kv缓存键值对，过期时间单位为秒
func Set(k string, v interface{}, expired int64)  {
    cache.Set(k, v, expired)
}

// (使用全局KV缓存对象)获取指定键名的值
func Get(k string) interface{} {
    return cache.Get(k)
}

// (使用全局KV缓存对象)删除指定键值对
func Remove(k string) {
    cache.Remove(k)
}

// 设置kv缓存键值对，过期时间单位为秒
func (c *Cache) Set(k string, v interface{}, expired int64)  {
    c.RLock()
    c.m[c.getIndex(k)].Set(k, v, expired)
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
    c.Lock()
    delete(c.m, k)
    c.Unlock()
}

// 获得键名的索引字母
func (c *Cache) getIndex(k string) string {
    s := gmd5.EncodeString(k)
    return strings.ToLower(string(s[0]))
}

// 设置kv缓存键值对，过期时间单位为秒
func (cm *CacheMap) Set(k string, v interface{}, expired int64)  {
    t := gtime.Second() + expired
    if expired == 0 {
        t = t + 86400*365*10
    }
    cm.Lock()
    cm.m[k] = &CacheItem{Value: v, Expired: t}
    cm.Unlock()
}

// 获取指定键名的值
func (cm *CacheMap) Get(k string) interface{} {
    cm.RLock()
    r, _ := cm.m[k]
    cm.RUnlock()
    if r != nil {
        if r.Expired < gtime.Second() {
            cm.Remove(k)
            return nil
        } else {
            return r.Value
        }
    }
    return nil
}

// 删除指定键值对
func (cm *CacheMap) Remove(k string) {
    cm.Lock()
    delete(cm.m, k)
    cm.Unlock()
}

// 自动清理过期键值对(每间隔60秒执行)
func (cm *CacheMap) autoClearLoop() {
    for {
        cm.Lock()
        for k, v := range cm.m {
            if v.Expired < gtime.Second() {
                delete(cm.m, k)
            }
        }
        cm.Unlock()
        time.Sleep(60 * time.Second)
    }
}



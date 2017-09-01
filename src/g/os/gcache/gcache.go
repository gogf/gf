package gcache

import (
    "sync"
    "g/util/gtime"
    "time"
)

type Cache struct {
    sync.RWMutex
    m map[string]*CacheItem
}

type CacheItem struct {
    Value   interface{}
    Expired int64
}

var cache *Cache = New()

func New() *Cache {
    r := &Cache{
        m : make(map[string]*CacheItem),
    }
    go r.autoClearLoop()
    return r
}

func Set(k string, v interface{}, expired int64)  {
    cache.Set(k, v, expired)
}

func Get(k string) interface{} {
    return cache.Get(k)
}

func Remove(k string) {
    cache.Remove(k)
}

// 设置kv缓存键值对，过期时间单位为秒
func (c *Cache) Set(k string, v interface{}, expired int64)  {
    c.Lock()
    c.m[k] = &CacheItem{Value: v, Expired: gtime.Second() + expired}
    c.Unlock()
}

// 获取指定键名的值
func (c *Cache) Get(k string) interface{} {
    c.RLock()
    r, _ := c.m[k]
    c.RUnlock()
    if r != nil {
        if r.Expired < gtime.Second() {
            c.Remove(k)
            return nil
        }
    }
    return r.Value
}

// 删除指定键值对
func (c *Cache) Remove(k string) {
    c.Lock()
    delete(c.m, k)
    c.RUnlock()
}

// 自动清理过期键值对
func (c *Cache) autoClearLoop() {
    for {
        c.Lock()
        for k, v := range c.m {
            if v.Expired < gtime.Second() {
                delete(c.m, k)
            }
        }
        c.Unlock()
        time.Sleep(60 * time.Second)
    }
}



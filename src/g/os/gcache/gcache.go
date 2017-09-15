package gcache

import (
    "sync"
    "g/util/gtime"
    "time"
    "g/encoding/gmd5"
    "strings"
    "g/encoding/gjson"
)

type Cache struct {
    sync.RWMutex
    m map[string]*CacheMap // 以键名首字母为索引
}

type CacheMap  struct {
    sync.RWMutex
    deleted bool              // 对象是否已删除，以便判断停止goroutine
    m1 map[string]interface{} // 不过期的键值对
    m2 map[string]CacheItem   // 有过期时间的键值对
}

type CacheItem struct {
    v interface{}
    e int64
}

var cache *Cache = New()

// Cache对象按照缓存键名首字母做了分组
func New() *Cache {
    c := &Cache {
        m : make(map[string]*CacheMap),
    }
    // 0 - 9
    for i := 48; i <= 57; i++ {
        m := &CacheMap {
            m1 : make(map[string]interface{}),
            m2 : make(map[string]CacheItem),
        }
        c.m[string(i)] = m
        go m.autoClearLoop()
    }
    // a - z
    for i := 97; i <= 122; i++ {
        m := &CacheMap {
            m1 : make(map[string]interface{}),
            m2 : make(map[string]CacheItem),
        }
        c.m[string(i)] = m
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
        for k1, _ := range cm.m1 {
            l = append(l, k1)
        }
        for k2, _ := range cm.m2 {
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
        for _, v1 := range cm.m1 {
            l = append(l, v1)
        }
        for _, v2 := range cm.m2 {
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
        size += len(cm.m1)
        size += len(cm.m2)
        cm.RUnlock()
    }
    c.RUnlock()
    return size
}

// 删除缓存对象
func (c *Cache) Destroy()  {
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

// 清空缓存对象（相当于新建一个新的缓存对象，旧的丢给GC处理）
func (c *Cache) Clear() {
    c.Destroy()
    c.Lock()
    c.m = New().m
    c.Unlock()
}

// 将数据导出为JSON字符串
func (c *Cache) Export() string {
    data := make(map[string]interface{})
    c.RLock()
    for _, cm := range c.m {
        cm.RLock()
        for k1, v1 := range cm.m1 {
            data[k1] = make(map[string]interface{})
            data[k1] = map[string]interface{}{
                "v": v1,
                "e": 0,
            }
        }
        for k2, v2 := range cm.m2 {
            data[k2] = make(map[string]interface{})
            data[k2] = map[string]interface{}{
                "v": v2.v,
                "e": v2.e,
            }
        }
        cm.RUnlock()
    }
    c.RUnlock()
    return gjson.Encode(data)
}

// 将导出的JSON字符串导入到缓存对象中
func (c *Cache) Import(s string) {
    data := make(map[string]map[string]interface{})
    gjson.DecodeTo(s, &data)
    for k, m := range data {
        // Set的第三个参数是过期时间数，因此这里需要计算导入的时候还剩多少时间过期
        expire := gtime.Millisecond() - int64(m["e"].(float64))
        if expire > 0 {
            c.Set(k, m["v"], expire)
        }
    }
}

func (c *Cache) getIndex(k string) string {
    s := gmd5.EncodeString(k)
    return strings.ToLower(string(s[0]))
}

// 设置kv缓存键值对，过期时间单位为毫秒
func (cm *CacheMap) Set(k string, v interface{}, expired int64)  {
    cm.Lock()
    if expired == 0 {
        cm.m1[k] = v
        if _, ok := cm.m2[k]; ok {
            delete(cm.m2, k)
        }
    } else {
        cm.m2[k] = CacheItem{v: v, e: gtime.Millisecond() + int64(expired)}
        if _, ok := cm.m1[k]; ok {
            delete(cm.m1, k)
        }
    }
    cm.Unlock()
}

// 获取指定键名的值
func (cm *CacheMap) Get(k string) interface{} {
    var v interface{}
    cm.RLock()
    if r1, ok := cm.m1[k]; ok {
        v = r1
    } else if r2, ok := cm.m2[k]; ok {
        if r2.e < gtime.Millisecond() {
            v = nil
        } else {
            v = r2.v
        }
    }
    cm.RUnlock()
    return v
}

// 删除指定键值对
func (cm *CacheMap) Remove(k string) {
    cm.Lock()
    delete(cm.m1, k)
    delete(cm.m2, k)
    cm.Unlock()
}

// 自动清理过期键值对(每间隔60秒执行)
func (cm *CacheMap) autoClearLoop() {
    for !cm.deleted {
        expired := make([]string, 0)
        cm.RLock()
        for k, v := range cm.m2 {
            if v.e < gtime.Millisecond() {
                expired = append(expired, k)
            }
        }
        cm.RUnlock()
        if len(expired) > 0 {
            cm.Lock()
            for _, k := range expired {
                delete(cm.m2, k)
            }
            cm.Unlock()
        }
        time.Sleep(60 * time.Second)
    }
}



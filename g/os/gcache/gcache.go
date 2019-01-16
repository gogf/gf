// Copyright 2017-2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gcache provides high performance and concurrent-safe in-memory cache for process.
// 
// 缓存模块,
// 并发安全的单进程高速缓存.
package gcache

// 全局缓存管理对象
var cache = New()

// (使用全局KV缓存对象)设置kv缓存键值对，过期时间单位为**毫秒**
func Set(key interface{}, value interface{}, expire int)  {
    cache.Set(key, value, expire)
}

// 当键名不存在时写入，并返回true；否则返回false。
// 常用来做对并发性要求不高的内存锁。
func SetIfNotExist(key interface{}, value interface{}, expire int) bool {
    return cache.SetIfNotExist(key, value, expire)
}

// (使用全局KV缓存对象)批量设置kv缓存键值对，过期时间单位为**毫秒**
func BatchSet(data map[interface{}]interface{}, expire int)  {
    cache.BatchSet(data, expire)
}

// (使用全局KV缓存对象)获取指定键名的值
func Get(key interface{}) interface{} {
    return cache.Get(key)
}

// 当键名存在时返回其键值，否则写入指定的键值
func GetOrSet(key interface{}, value interface{}, expire int) interface{} {
    return cache.GetOrSet(key, value, expire)
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func GetOrSetFunc(key interface{}, f func() interface{}, expire int) interface{} {
    return cache.GetOrSetFunc(key, f, expire)
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func GetOrSetFuncLock(key interface{}, f func() interface{}, expire int) interface{} {
    return cache.GetOrSetFuncLock(key, f, expire)
}

// 是否存在指定的键名，true表示存在，false表示不存在。
func Contains(key interface{}) bool {
    return cache.Contains(key)
}

// (使用全局KV缓存对象)删除指定键值对
func Remove(key interface{}) interface{} {
    return cache.Remove(key)
}

// (使用全局KV缓存对象)批量删除指定键值对
func BatchRemove(keys []interface{}) {
    cache.BatchRemove(keys)
}

// 返回缓存的所有数据键值对(不包含已过期数据)
func Data() map[interface{}]interface{} {
    return cache.Data()
}

// 获得所有的键名，组成数组返回
func Keys() []interface{} {
    return cache.Keys()
}

// 获得所有的键名，组成字符串数组返回
func KeyStrings() []string {
    return cache.KeyStrings()
}

// 获得所有的值，组成数组返回
func Values() []interface{} {
    return cache.Values()
}

// 获得缓存对象的键值对数量
func Size() int {
    return cache.Size()
}

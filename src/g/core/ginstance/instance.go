package ginstance

import "g/core/types/gnmap"

// 单例对象存储器
var instances = gnmap.NewSafeMap()

// 获取单例对象
func Get(k string) interface{} {
    if v, ok := instances.Get(k); ok {
        return v
    } else {
        return nil
    }
}

// 设置单例对象
func Set(k string, v interface{}) {
    instances.Put(k, v)
}
package gconfig

import (
    "g/core/types/gmap"
)

// 配置对象
var config = gmap.NewStringInterfaceMap()

// 获取配置
func Get(k string) interface{} {
    if v, ok := config.Get(k); ok {
        return v
    } else {
        return nil
    }
}

// 设置配置
func Set(k string, v interface{}) {
    config.Set(k, v)
}

// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfcache provides reading and caching for file contents.
// 
// 文件缓存.
package gfcache

import (
    "github.com/gogf/gf/g/internal/cmdenv"
    "github.com/gogf/gf/g/os/gcache"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/gfsnotify"
)

const (
    // 默认的缓存超时时间(60秒)
    gDEFAULT_CACHE_EXPIRE = 60
)

var (
    // 默认的缓存时间(秒)
    cacheExpire = cmdenv.Get("gf.gfcache.expire", gDEFAULT_CACHE_EXPIRE).Int()*1000
)

// 获得文件内容 string，expire参数为缓存过期时间，单位为秒。
func GetContents(path string, expire...int) string {
    return string(GetBinContents(path, expire...))
}

// 获得文件内容 []byte，expire参数为缓存过期时间，单位为秒。
func GetBinContents(path string, expire...int) []byte {
    k := cacheKey(path)
    e := cacheExpire
    if len(expire) > 0 {
        e = expire[0]
    }
    r := gcache.GetOrSetFuncLock(k, func() interface{} {
        b := gfile.GetBinContents(path)
        if b != nil {
            // 添加文件监控，如果文件有任何变化，立即清空缓存
            gfsnotify.Add(path, func(event *gfsnotify.Event) {
                gcache.Remove(k)
                gfsnotify.Exit()
            })
        }
        return b
    }, e*1000)
    if r != nil {
        return r.([]byte)
    }
    return nil
}

// 生成缓存键名
func cacheKey(path string) string {
    return "gf.gfcache:" + path
}
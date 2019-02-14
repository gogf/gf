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
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/internal/cmdenv"
)

type Cache struct {
    cap    *gtype.Int               // 缓存容量(byte)，设置为0表示不限制
    size   *gtype.Int               // 缓存大小(Byte)
    cache  *gmap.StringInterfaceMap // 缓存对象
}

const (
    // 默认的缓存容量(10MB)
    gDEFAULT_CACHE_CAP = 10*1024*1024
)

var (
    // 默认的缓存容量
    cacheCap = cmdenv.Get("gf.gfcache.cap", gDEFAULT_CACHE_CAP).Int()
    // 默认的文件缓存对象
    cache    = New()
)

func New(cap ... int) *Cache {
    c := cacheCap
    if len(cap) > 0 {
        c = cap[0]
    }
    return &Cache {
        cap    : gtype.NewInt(c),
        size   : gtype.NewInt(),
        cache  : gmap.NewStringInterfaceMap(),
    }
}

// 获得已缓存的文件大小(byte)
func GetSize() int {
    return cache.GetSize()
}

// 获得文件内容 string
func GetContents(path string) string {
    return cache.GetContents(path)
}

// 获得文件内容 []byte
func GetBinContents(path string) []byte {
    return cache.GetBinContents(path)
}
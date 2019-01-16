// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gfcache provides reading and caching for file contents.
// 
// 文件缓存.
package gfcache

import (
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gcache"
)

type Cache struct {
    cap    *gtype.Int         // 缓存容量(byte)，设置为0表示不限制
    size   *gtype.Int         // 缓存大小(Byte)
    cache  *gcache.Cache      // 缓存对象
}

const (
    // 默认的缓存容量(不限制)
    gDEFAULT_CACHE_CAP = 0
)

var (
    // 默认的文件缓存对象
    cache = New()
)

func New(cap ... int) *Cache {
    c := gDEFAULT_CACHE_CAP
    if len(cap) > 0 {
        c = cap[0]
    }
    return &Cache {
        cap    : gtype.NewInt(c),
        size   : gtype.NewInt(),
        cache  : gcache.New(),
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
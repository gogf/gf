// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gfcache provides reading and caching for file contents.
package gfcache

import (
	"github.com/gogf/gf/g/internal/cmdenv"
	"github.com/gogf/gf/g/os/gcache"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfsnotify"
)

const (
	// Default expire time for file content caching in seconds.
	gDEFAULT_CACHE_EXPIRE = 60
)

var (
	// Default expire time for file content caching in seconds.
	cacheExpire = cmdenv.Get("gf.gfcache.expire", gDEFAULT_CACHE_EXPIRE).Int() * 1000
)

// GetContents returns string content of given file by <path> from cache.
// If there's no content in the cache, it will read it from disk file specified by <path>.
// The parameter <expire> specifies the caching time for this file content in seconds.
func GetContents(path string, expire ...int) string {
	return string(GetBinContents(path, expire...))
}

// GetBinContents returns []byte content of given file by <path> from cache.
// If there's no content in the cache, it will read it from disk file specified by <path>.
// The parameter <expire> specifies the caching time for this file content in seconds.
func GetBinContents(path string, expire ...int) []byte {
	k := cacheKey(path)
	e := cacheExpire
	if len(expire) > 0 {
		e = expire[0]
	}
	r := gcache.GetOrSetFuncLock(k, func() interface{} {
		b := gfile.GetBinContents(path)
		if b != nil {
			// Adding this <path> to gfsnotify,
			// it will clear its cache if there's any changes of the file.
			_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
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

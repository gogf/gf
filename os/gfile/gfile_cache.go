// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfsnotify"
)

const (
	defaultCacheExpire    = "1m"             // defaultCacheExpire is the expire time for file content caching in seconds.
	commandEnvKeyForCache = "gf.gfile.cache" // commandEnvKeyForCache is the configuration key for command argument or environment configuring cache expire duration.
)

var (
	// Default expire time for file content caching.
	cacheExpire = getCacheExpire()

	// internalCache is the memory cache for internal usage.
	internalCache = gcache.New()
)

func getCacheExpire() time.Duration {
	d, err := time.ParseDuration(command.GetOptWithEnv(commandEnvKeyForCache, defaultCacheExpire))
	if err != nil {
		panic(err)
	}
	return d
}

// GetContentsWithCache returns string content of given file by `path` from cache.
// If there's no content in the cache, it will read it from disk file specified by `path`.
// The parameter `expire` specifies the caching time for this file content in seconds.
func GetContentsWithCache(path string, duration ...time.Duration) string {
	return string(GetBytesWithCache(path, duration...))
}

// GetBytesWithCache returns []byte content of given file by `path` from cache.
// If there's no content in the cache, it will read it from disk file specified by `path`.
// The parameter `expire` specifies the caching time for this file content in seconds.
func GetBytesWithCache(path string, duration ...time.Duration) []byte {
	var (
		ctx      = context.Background()
		expire   = cacheExpire
		cacheKey = commandEnvKeyForCache + path
	)

	if len(duration) > 0 {
		expire = duration[0]
	}
	r, _ := internalCache.GetOrSetFuncLock(ctx, cacheKey, func(ctx context.Context) (interface{}, error) {
		b := GetBytes(path)
		if b != nil {
			// Adding this `path` to gfsnotify,
			// it will clear its cache if there's any changes of the file.
			_, _ = gfsnotify.Add(path, func(event *gfsnotify.Event) {
				_, err := internalCache.Remove(ctx, cacheKey)
				if err != nil {
					intlog.Errorf(ctx, `%+v`, err)
				}
				gfsnotify.Exit()
			})
		}
		return b, nil
	}, expire)
	if r != nil {
		return r.Bytes()
	}
	return nil
}

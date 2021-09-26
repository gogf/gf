// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import "github.com/go-redis/redis/v8"

// localAdapterGoRedisPoolStats is statistics of redis connection pool.
type localAdapterGoRedisPoolStats struct {
	*redis.PoolStats
}

// ActiveCount is the number of connections in the pool. The count includes
// idle connections and connections in use.
func (s *localAdapterGoRedisPoolStats) ActiveCount() int64 {
	if s == nil {
		return -1
	}
	return int64(s.PoolStats.TotalConns)
}

// IdleCount is the number of idle connections in the pool.
func (s *localAdapterGoRedisPoolStats) IdleCount() int64 {
	if s == nil {
		return -1
	}
	return int64(s.PoolStats.IdleConns)
}

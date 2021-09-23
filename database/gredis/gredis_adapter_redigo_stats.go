// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import "github.com/gomodule/redigo/redis"

// localAdapterRedigoPoolStats is statistics of redis connection pool.
type localAdapterRedigoPoolStats struct {
	redis.PoolStats
}

// ActiveCount is the number of connections in the pool. The count includes
// idle connections and connections in use.
func (s *localAdapterRedigoPoolStats) ActiveCount() int64 {
	if s == nil {
		return -1
	}
	return int64(s.PoolStats.ActiveCount)
}

// IdleCount is the number of idle connections in the pool.
func (s *localAdapterRedigoPoolStats) IdleCount() int64 {
	if s == nil {
		return -1
	}
	return int64(s.PoolStats.IdleCount)
}

// WaitCount is the total number of connections waited for.
// This value is currently not guaranteed to be 100% accurate.
func (s *localAdapterRedigoPoolStats) WaitCount() int64 {
	if s == nil {
		return -1
	}
	return s.PoolStats.WaitCount
}

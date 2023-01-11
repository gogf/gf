// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	sentinelCtx    = context.TODO()
	sentinelConfig = &gredis.Config{
		Address:    `127.0.0.1:26379,127.0.0.1:26380,127.0.0.1:26381`,
		MasterName: `mymaster`,
		Pass:       "111111",
	}
)

func TestConn_sentinel_master(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sentinelConfig.SlaveOnly = false
		redis, err := gredis.New(sentinelConfig)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(sentinelCtx)

		conn, err := redis.Conn(sentinelCtx)
		t.AssertNil(err)
		defer conn.Close(sentinelCtx)

		_, err = conn.Do(sentinelCtx, "set", "test", "123")
		t.AssertNil(err)
		defer conn.Do(sentinelCtx, "del", "test")

		r, err := conn.Do(sentinelCtx, "get", "test")
		t.AssertNil(err)
		t.Assert(r.String(), "123")
	})
}

func TestConn_sentinel_slave(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sentinelConfig.SlaveOnly = true
		redis, err := gredis.New(sentinelConfig)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(sentinelCtx)

		conn, err := redis.Conn(sentinelCtx)
		t.AssertNil(err)
		defer conn.Close(sentinelCtx)

		_, err = conn.Do(sentinelCtx, "set", "test", "123")
		t.AssertNQ(err, nil)
	})
}

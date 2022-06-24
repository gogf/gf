/**
 * package gredis
 *
 * @Author 曾洪亮<zenghongl@126.com>
 * @Email  zenghongl@126.com
 * User: whoSafe
 * Date: 2022/6/24
 * Time: 11:19
 */

package gredis_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	sentinelCtx    = context.TODO()
	sentinelConfig = &gredis.Config{
		Address:    `192.168.41.162:26379,192.168.41.174:26379,192.168.41.192:26379`,
		MasterName: `master`,
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

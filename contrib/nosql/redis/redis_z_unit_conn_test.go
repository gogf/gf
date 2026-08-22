// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

func TestConn_DoWithTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Do(ctx, "set", "test", "123")
		t.AssertNil(err)
		defer conn.Do(ctx, "del", "test")

		r, err := conn.Do(ctx, "get", "test")
		t.AssertNil(err)
		t.Assert(r.String(), "123")
	})
}

func TestConn_Transaction(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = redis.Do(ctx, "set", "test:tx:probe", "real")
		t.AssertNil(err)
		defer redis.Do(ctx, "del", "test:tx:probe")

		_, err = redis.Do(ctx, "del", "test:tx:set")
		t.AssertNil(err)

		_, err = conn.Do(ctx, "multi")
		t.AssertNil(err)

		// While the transaction is open on the pinned connection, commands sent
		// through the shared client must not be affected by the MULTI state.
		v, err := redis.Do(ctx, "get", "test:tx:probe")
		t.AssertNil(err)
		t.Assert(v.String(), "real")

		_, err = conn.Do(ctx, "sadd", "test:tx:set", "m1")
		t.AssertNil(err)
		_, err = conn.Do(ctx, "sadd", "test:tx:set", "m2")
		t.AssertNil(err)

		v, err = conn.Do(ctx, "exec")
		t.AssertNil(err)
		t.Assert(v.Strings(), []string{"1", "1"})
		defer redis.Do(ctx, "del", "test:tx:set")

		// The connection is clean after EXEC.
		v, err = conn.Do(ctx, "scard", "test:tx:set")
		t.AssertNil(err)
		t.Assert(v.Int(), 2)
	})
}

func TestConn_ReceiveVarWithTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		sub, err := conn.Subscribe(ctx, "gf")
		t.AssertNil(err)
		t.Assert(sub[0].Channel, "gf")

		_, err = redis.Publish(ctx, "gf", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gf")
		t.Assert(msg.Payload, "test")
	})
}

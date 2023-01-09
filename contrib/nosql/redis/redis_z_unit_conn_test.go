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

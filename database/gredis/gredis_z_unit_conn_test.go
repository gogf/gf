// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"context"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

var (
	ctx = context.TODO()
)

func TestConn_DoWithTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Do(ctx, "set", "test", "123", &gredis.Option{ReadTimeout: time.Second})
		t.Assert(err, nil)
		defer conn.Do(ctx, "del", "test")

		r, err := conn.Do(ctx, "get", "test", &gredis.Option{ReadTimeout: time.Second})
		t.Assert(err, nil)
		t.Assert(r.String(), "123")
	})
}

func TestConn_ReceiveVarWithTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Do(ctx, "Subscribe", "gf", &gredis.Option{ReadTimeout: time.Second})
		t.AssertNil(err)

		v, err := redis.Do(ctx, "PUBLISH", "gf", "test", &gredis.Option{ReadTimeout: time.Second})
		t.Assert(err, nil)
		t.Assert(v.String(), "1")

		v, _ = conn.Receive(ctx)
		t.Assert(len(v.Strings()), 3)
		t.Assert(v.Strings()[2], "test")
	})
}

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
		t.AssertNil(redis)

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err := conn.Do(ctx, "set", "test", "123", &gredis.Option{ReadTimeout: time.Second})
		t.Assert(err, nil)
		defer conn.DoWithTimeout(time.Second, "del", "test")

		r, err := conn.DoWithTimeout(time.Second, "get", "test")
		t.Assert(err, nil)
		t.Assert(r, "123")
	})
}

func TestConn_ReceiveVarWithTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		t.AssertNE(redis, nil)
		conn := redis.Conn()
		defer conn.Close()

		_, err := conn.DoVarWithTimeout(time.Second, "Subscribe", "gf")
		t.Assert(err, nil)

		v, err := redis.DoVarWithTimeout(time.Second, "PUBLISH", "gf", "test")
		t.Assert(err, nil)
		t.Assert(v.String(), "1")

		v, _ = conn.ReceiveVar()
		t.Assert(len(v.Strings()), 3)
		t.Assert(v.Strings()[2], "test")
	})
}

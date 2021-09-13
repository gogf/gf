// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func TestConn_DoWithTimeout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		t.AssertNE(redis, nil)
		conn := redis.Conn()
		defer conn.Close()

		_, err := conn.DoWithTimeout(time.Second, "set", "test", "123")
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

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

func Test_GroupPubSub_Publish(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		conn, subs, err := redis.Subscribe(ctx, "gf")
		t.AssertNil(err)
		t.Assert(subs[0].Channel, "gf")

		defer conn.Close(ctx)

		_, err = redis.Publish(ctx, "gf", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gf")
		t.Assert(msg.Payload, "test")
	})
}

func Test_GroupPubSub_Subscribe(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		conn, subs, err := redis.Subscribe(ctx, "aa", "bb", "gf")
		t.AssertNil(err)
		t.Assert(len(subs), 3)
		t.Assert(subs[0].Channel, "aa")
		t.Assert(subs[1].Channel, "bb")
		t.Assert(subs[2].Channel, "gf")

		defer conn.Close(ctx)

		_, err = redis.Publish(ctx, "gf", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gf")
		t.Assert(msg.Payload, "test")
	})
}

func Test_GroupPubSub_PSubscribe(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		conn, subs, err := redis.PSubscribe(ctx, "aa", "bb", "g?")
		t.AssertNil(err)
		t.Assert(len(subs), 3)
		t.Assert(subs[0].Channel, "aa")
		t.Assert(subs[1].Channel, "bb")
		t.Assert(subs[2].Channel, "g?")

		defer conn.Close(ctx)

		_, err = redis.Publish(ctx, "gf", "test")
		t.AssertNil(err)

		msg, err := conn.ReceiveMessage(ctx)
		t.AssertNil(err)
		t.Assert(msg.Channel, "gf")
		t.Assert(msg.Payload, "test")
	})
}

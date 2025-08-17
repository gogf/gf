// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_GroupString_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = guid.S()
			v1 = guid.S()
			k2 = guid.S()
			v2 = guid.S()
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, k2, v2)
		t.AssertNil(err)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)
		r2, err := redis.GroupString().Get(ctx, k2)
		t.AssertNil(err)
		t.Assert(r2.String(), v2)
	})
	// With Option.
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		_, err = redis.GroupString().Set(ctx, k1, v2, gredis.SetOption{
			NX: true,
			TTLOption: gredis.TTLOption{
				EX: gconv.PtrInt64(60),
			},
		})
		t.AssertNil(err)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		_, err = redis.GroupString().Set(ctx, k1, v2, gredis.SetOption{
			XX: true,
			TTLOption: gredis.TTLOption{
				EX: gconv.PtrInt64(60),
			},
		})
		t.AssertNil(err)

		r2, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r2.String(), v2)
	})
}

func Test_GroupString_SetNX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		_, err = redis.GroupString().SetNX(ctx, k1, v2)
		t.AssertNil(err)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)
	})
}

func Test_GroupString_SetEX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
		)
		err := redis.GroupString().SetEX(ctx, k1, v1, 1)
		t.AssertNil(err)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		time.Sleep(time.Second * 2)

		r2, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r2.String(), "")
	})
}

func Test_GroupString_GetDel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().GetDel(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		r2, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r2.String(), "")
	})
}

func Test_GroupString_GetEX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
		)
		err := redis.GroupString().SetEX(ctx, k1, v1, 1)
		t.AssertNil(err)

		r1, err := redis.GroupString().GetEX(ctx, k1, gredis.GetEXOption{
			Persist: true,
		})
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		time.Sleep(2 * time.Second)

		r2, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r2.String(), v1)
	})
}

func Test_GroupString_GetSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			k2 = "k2"
			v2 = "v2"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		r2, err := redis.GroupString().GetSet(ctx, k1, v2)
		t.AssertNil(err)
		t.Assert(r2.String(), v1)

		r3, err := redis.GroupString().GetSet(ctx, k2, v2)
		t.AssertNil(err)
		t.Assert(r3.String(), "")

		r4, err := redis.GroupString().GetSet(ctx, k2, v2)
		t.AssertNil(err)
		t.Assert(r4.String(), v2)
	})
}

func Test_GroupString_StrLen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().StrLen(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1, 2)
	})
}

func Test_GroupString_Append(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().Append(ctx, k1, v2)
		t.AssertNil(err)
		t.Assert(r1, len(v1+v2))

		r2, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r2.String(), v1+v2)
	})
}

func Test_GroupString_SetRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "v1"
			v2 = "v2"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().SetRange(ctx, k1, 2, v2)
		t.AssertNil(err)
		t.Assert(r1, len(v1+v2))

		r2, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r2.String(), v1+v2)
	})
}

func Test_GroupString_GetRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = "hello gf"
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().GetRange(ctx, k1, 6, 8)
		t.AssertNil(err)
		t.Assert(r1, "gf")
	})
}

func Test_GroupString_Incr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = 1
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().Incr(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1, 2)
	})
}

func Test_GroupString_IncrBy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = 1
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().IncrBy(ctx, k1, 10)
		t.AssertNil(err)
		t.Assert(r1, 11)
	})
}

func Test_GroupString_IncrByFloat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = 1
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().IncrByFloat(ctx, k1, 1.01)
		t.AssertNil(err)
		t.Assert(r1, 2.01)
	})
}

func Test_GroupString_Decr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = 10
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().Decr(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1, 9)
	})
}

func Test_GroupString_DecrBy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = "k1"
			v1 = 10
		)
		_, err := redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)

		r1, err := redis.GroupString().DecrBy(ctx, k1, 3)
		t.AssertNil(err)
		t.Assert(r1, 7)
	})
}

func Test_GroupString_MSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = guid.S()
			v1 = guid.S()
			k2 = guid.S()
			v2 = guid.S()
		)
		err := redis.GroupString().MSet(ctx, map[string]interface{}{
			k1: v1,
			k2: v2,
		})
		t.AssertNil(err)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		r2, err := redis.GroupString().Get(ctx, k2)
		t.AssertNil(err)
		t.Assert(r2.String(), v2)
	})
}

func Test_GroupString_MSetNX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = guid.S()
			v1 = guid.S()
			k2 = guid.S()
			v2 = guid.S()
		)
		ok, err := redis.GroupString().MSetNX(ctx, map[string]interface{}{
			k1: v1,
		})
		t.AssertNil(err)
		t.Assert(ok, true)

		ok, err = redis.GroupString().MSetNX(ctx, map[string]interface{}{
			k1: v1,
			k2: v2,
		})
		t.AssertNil(err)
		t.Assert(ok, false)

		r1, err := redis.GroupString().Get(ctx, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		r2, err := redis.GroupString().Get(ctx, k2)
		t.AssertNil(err)
		t.Assert(r2.String(), "")
	})
}

func Test_GroupString_MGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1 = guid.S()
			v1 = guid.S()
			k2 = guid.S()
			v2 = guid.S()
		)
		err := redis.GroupString().MSet(ctx, map[string]interface{}{
			k1: v1,
			k2: v2,
		})
		t.AssertNil(err)

		r1, err := redis.GroupString().MGet(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(len(r1), 2)
		t.Assert(r1[k1].String(), v1)
		t.Assert(r1[k2].String(), v2)
	})
}

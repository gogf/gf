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
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

var (
	TestKey   = "mykey"
	TestValue = "hello"
)

func Test_GroupGeneric_Copy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1     = guid.S()
			v1     = guid.S()
			k2     = guid.S()
			result int64
			err    error
		)
		_, err = redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)
		result, err = redis.GroupGeneric().Copy(ctx, k1, k2)
		t.AssertEQ(result, int64(1))
		t.AssertNil(err)
		v2, err := redis.GroupString().Get(ctx, k2)
		t.AssertNil(err)
		t.Assert(v2.String(), v1)
	})
	// With Option.
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1     = guid.S()
			v1     = guid.S()
			k2     = guid.S()
			result int64
			err    error
		)
		_, err = redis.GroupString().Set(ctx, k1, v1)
		t.AssertNil(err)
		result, err = redis.GroupGeneric().Copy(ctx, k1, k2, gredis.CopyOption{
			DB:      1,
			REPLACE: true,
		})
		t.AssertEQ(result, int64(1))
		t.AssertNil(err)
		v2, err := redis.GroupString().Get(ctx, k2)
		t.AssertNil(err)
		t.Assert(v2.String(), v1)
	})
}

func Test_GroupGeneric_Exists(t *testing.T) {
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
		result, err := redis.GroupGeneric().Exists(ctx, k1)
		t.AssertEQ(result, int64(1))
		t.AssertNil(err)
		result, err = redis.GroupGeneric().Exists(ctx, "nosuchkey")
		t.AssertEQ(result, int64(0))
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, k2, v2)
		t.AssertNil(err)
		result, err = redis.GroupGeneric().Exists(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(result, int64(2))
	})
}

func Test_GroupGeneric_Type(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupList().LPush(ctx, "k2", "v2")
		t.AssertNil(err)
		_, err = redis.GroupSet().SAdd(ctx, "k3", "v3")
		t.AssertNil(err)

		t1, err := redis.GroupGeneric().Type(ctx, "k1")
		t.AssertNil(err)
		t.AssertEQ(t1, "string")
		t2, err := redis.GroupGeneric().Type(ctx, "k2")
		t.AssertNil(err)
		t.AssertEQ(t2, "list")
		t3, err := redis.GroupGeneric().Type(ctx, "k3")
		t.AssertNil(err)
		t.AssertEQ(t3, "set")
	})
}

func Test_GroupGeneric_Unlink(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)

		result, err := redis.GroupGeneric().Unlink(ctx, "k1", "k2", "k3")
		t.AssertNil(err)
		t.AssertEQ(result, int64(2))
		v1, err := redis.GroupString().Get(ctx, "k1")
		t.AssertNil(err)
		t.AssertEQ(v1.String(), "")
		v2, err := redis.GroupString().Get(ctx, "k2")
		t.AssertNil(err)
		t.AssertEQ(v2.String(), "")
	})
}

func Test_GroupGeneric_Rename(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		err = redis.GroupGeneric().Rename(ctx, "k1", "k2")
		t.AssertNil(err)
		v1, err := redis.GroupString().Get(ctx, "k1")
		t.AssertNil(err)
		t.AssertEQ(v1.String(), "")
		v2, err := redis.GroupString().Get(ctx, "k2")
		t.AssertNil(err)
		t.AssertEQ(v2.String(), "v1")
	})
}

func Test_GroupGeneric_RenameNX(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)
		result, err := redis.GroupGeneric().RenameNX(ctx, "k1", "k2")
		t.AssertNil(err)
		t.AssertEQ(result, int64(0))
		result, err = redis.GroupGeneric().RenameNX(ctx, "k1", "k3")
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		v2, err := redis.GroupString().Get(ctx, "k2")
		t.AssertNil(err)
		t.AssertEQ(v2.String(), "v2")
		v3, err := redis.GroupString().Get(ctx, "k3")
		t.AssertNil(err)
		t.AssertEQ(v3.String(), "v1")
	})
}

func Test_GroupGeneric_Move(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushAll(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Move(ctx, "k1", 0)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
	})
}

func Test_GroupGeneric_Del(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Del(ctx, "k1", "k2", "k3")
		t.AssertNil(err)
		t.AssertEQ(result, int64(2))
		v1, err := redis.GroupString().Get(ctx, "k1")
		t.AssertNil(err)
		t.AssertEQ(v1.String(), "")
		v2, err := redis.GroupString().Get(ctx, "k2")
		t.AssertNil(err)
		t.AssertEQ(v2.String(), "")
	})
}

func Test_GroupGeneric_RandomKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		key, err := redis.GroupGeneric().RandomKey(ctx)
		t.AssertNil(err)
		t.AssertEQ(key, "")

		_, err = redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)

		key, err = redis.GroupGeneric().RandomKey(ctx)
		t.AssertNil(err)
		t.AssertIN(key, []string{"k1", "k2"})
	})
}

func Test_GroupGeneric_DBSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		dbSize, err := redis.GroupGeneric().DBSize(ctx)
		t.AssertNil(err)
		t.AssertEQ(dbSize, int64(0))

		_, err = redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)

		dbSize, err = redis.GroupGeneric().DBSize(ctx)
		t.AssertNil(err)
		t.AssertEQ(dbSize, int64(2))
	})
}

func Test_GroupGeneric_Keys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		err := redis.GroupString().MSet(ctx, map[string]interface{}{
			"firstname": "Jack",
			"lastname":  "Stuntman",
			"age":       35,
		})
		t.AssertNil(err)
		keys, err := redis.GroupGeneric().Keys(ctx, "*name*")
		t.AssertNil(err)
		t.AssertIN(keys, []string{"lastname", "firstname"})
		keys, err = redis.GroupGeneric().Keys(ctx, "a??")
		t.AssertNil(err)
		t.AssertEQ(keys, []string{"age"})
		keys, err = redis.GroupGeneric().Keys(ctx, "*")
		t.AssertNil(err)
		t.AssertIN(keys, []string{"lastname", "firstname", "age"})
	})
}

func Test_GroupGeneric_FlushDB(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)

		dbSize, err := redis.GroupGeneric().DBSize(ctx)
		t.AssertNil(err)
		t.AssertEQ(dbSize, int64(2))

		err = redis.GroupGeneric().FlushDB(ctx)
		t.AssertNil(err)

		dbSize, err = redis.GroupGeneric().DBSize(ctx)
		t.AssertNil(err)
		t.AssertEQ(dbSize, int64(0))
	})
}

func Test_GroupGeneric_FlushAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		_, err = redis.GroupString().Set(ctx, "k2", "v2")
		t.AssertNil(err)

		dbSize, err := redis.GroupGeneric().DBSize(ctx)
		t.AssertNil(err)
		t.AssertEQ(dbSize, int64(2))

		err = redis.GroupGeneric().FlushAll(ctx)
		t.AssertNil(err)

		dbSize, err = redis.GroupGeneric().DBSize(ctx)
		t.AssertNil(err)
		t.AssertEQ(dbSize, int64(0))
	})
}

func Test_GroupGeneric_Expire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Expire(ctx, TestKey, 1)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		ttl, err := redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(ttl, int64(1))
	})
	// With Option.
	// Starting with Redis version 7.0.0: Added options: NX, XX, GT and LT.
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		ttl, err := redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(ttl, int64(-1))
		result, err := redis.GroupGeneric().Expire(ctx, TestKey, 1, gredis.ExpireOption{XX: true})
		t.AssertNil(err)
		t.AssertEQ(result, int64(0))
		ttl, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(ttl, int64(-1))
		result, err = redis.GroupGeneric().Expire(ctx, TestKey, 1, gredis.ExpireOption{NX: true})
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		ttl, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(ttl, int64(1))
	})
}

func Test_GroupGeneric_ExpireAt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Exists(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().ExpireAt(ctx, TestKey, time.Now().Add(time.Millisecond*100))
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		time.Sleep(time.Millisecond * 100)
		result, err = redis.GroupGeneric().Exists(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(0))
	})
	// With Option.
	// Starting with Redis version 7.0.0: Added options: NX, XX, GT and LT.
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		ttl, err := redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(ttl, int64(-1))
		result, err := redis.GroupGeneric().ExpireAt(ctx, TestKey, time.Now().Add(time.Millisecond*100), gredis.ExpireOption{XX: true})
		t.AssertNil(err)
		t.AssertEQ(result, int64(0))
		result, err = redis.GroupGeneric().ExpireAt(ctx, TestKey, time.Now().Add(time.Minute), gredis.ExpireOption{NX: true})
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		ttl, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertGT(ttl, int64(0))
	})
}

func Test_GroupGeneric_ExpireTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		expireTime := time.Now().Add(time.Minute)
		result, err := redis.GroupGeneric().ExpireAt(ctx, TestKey, expireTime)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		resultTime, err := redis.GroupGeneric().ExpireTime(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(resultTime.Int64(), expireTime.Unix())

		_, err = redis.GroupString().Set(ctx, "noExpireKey", TestValue)
		t.AssertNil(err)
		resultTime, err = redis.GroupGeneric().ExpireTime(ctx, "noExpireKey")
		t.AssertNil(err)
		t.AssertEQ(resultTime.Int64(), int64(-1))

		resultTime, err = redis.GroupGeneric().ExpireTime(ctx, "noExistKey")
		t.AssertNil(err)
		t.AssertEQ(resultTime.Int64(), int64(-2))
	})
}

func Test_GroupGeneric_TTL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Expire(ctx, TestKey, 10)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(10))
	})
}

func Test_GroupGeneric_Persist(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Expire(ctx, TestKey, 10)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(10))
		result, err = redis.GroupGeneric().Persist(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(-1))
	})
}

func Test_GroupGeneric_PExpire(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().PExpire(ctx, TestKey, 2500)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().PTTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertLE(result, int64(2500))
	})

	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().PExpire(ctx, TestKey, 2500, gredis.ExpireOption{
			NX: true,
		})
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))

		result, err = redis.GroupGeneric().PExpire(ctx, TestKey, 2500, gredis.ExpireOption{
			NX: true,
		})
		t.AssertNil(err)
		t.AssertEQ(result, int64(0))

		result, err = redis.GroupGeneric().PTTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertLE(result, int64(2500))
	})
}

func Test_GroupGeneric_PExpireAt(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().PExpireAt(ctx, TestKey, time.Now().Add(-time.Hour))
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().TTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(-2))
		result, err = redis.GroupGeneric().PTTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(result, int64(-2))
	})
}

func Test_GroupGeneric_PExpireTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		expireTime := time.Now().Add(time.Hour)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().PExpireAt(ctx, TestKey, expireTime)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		resultTime, err := redis.GroupGeneric().PExpireTime(ctx, TestKey)
		t.AssertNil(err)
		t.AssertEQ(resultTime.Int64(), gtime.NewFromTime(expireTime).TimestampMilli())
	})
}

func Test_GroupGeneric_PTTL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		_, err := redis.GroupString().Set(ctx, TestKey, TestValue)
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Expire(ctx, TestKey, 1)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
		result, err = redis.GroupGeneric().PTTL(ctx, TestKey)
		t.AssertNil(err)
		t.AssertLE(result, int64(1000))
	})
}

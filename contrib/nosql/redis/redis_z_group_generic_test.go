// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"testing"
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
		defer redis.FlushDB(ctx)
		_, err := redis.GroupString().Set(ctx, "k1", "v1")
		t.AssertNil(err)
		result, err := redis.GroupGeneric().Move(ctx, "k1", 0)
		t.AssertNil(err)
		t.AssertEQ(result, int64(1))
	})
}

func Test_GroupGeneric_Del(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//defer redis.FlushDB(ctx)
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

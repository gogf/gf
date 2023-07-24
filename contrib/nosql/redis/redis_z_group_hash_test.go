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

func Test_GroupHash_HSet(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key         = "myhash"
			field1      = "field1"
			field1Value = "Hello"
			fields      = map[string]interface{}{
				field1: field1Value,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HGet(ctx, key, field1)
		t.AssertNil(err)
		t.Assert(r1.String(), field1Value)
	})
}

func Test_GroupHash_HSetNX(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			field1      = "field1"
			field1Value = "Hello"
			key         = "myhash"
		)
		r1, err := redis.HSetNX(ctx, key, field1, field1Value)
		t.AssertNil(err)
		t.Assert(r1, 1)

		r2, err := redis.HSetNX(ctx, key, field1, "World")
		t.AssertNil(err)
		t.Assert(r2, 0)
	})
}

func Test_GroupHash_HStrLen(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key         = "myhash"
			field1      = "field1"
			field1Value = "Hello"
			field2      = "field2"
			field2Value = "Hello World"
			fields      = map[string]interface{}{
				field1: field1Value,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		fieldValueLen, err := redis.HStrLen(ctx, key, field1)
		t.AssertNil(err)
		t.Assert(5, fieldValueLen)

		fields[field2] = field2Value
		_, err = redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		fieldValueLen, err = redis.HStrLen(ctx, key, field2)
		t.AssertNil(err)
		t.Assert(11, fieldValueLen)
	})
}

func Test_GroupHash_HExists(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key         = "myhash"
			field1      = "field1"
			field1Value = "Hello"
			fields      = map[string]interface{}{
				field1: field1Value,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HExists(ctx, key, field1)
		t.AssertNil(err)
		t.Assert(1, r1)

		r2, err := redis.HExists(ctx, key, "name")
		t.AssertNil(err)
		t.Assert(0, r2)
	})
}

func Test_GroupHash_HDel(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "myhash"
			k1     = "k1"
			v1     = "v1"
			k2     = "k2"
			v2     = "v2"
			k3     = "k3"
			v3     = "v3"
			fields = map[string]interface{}{
				k1: v1,
				k2: v2,
				k3: v3,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HDel(ctx, key, k1)
		t.AssertNil(err)
		t.Assert(1, r1)

		r2, err := redis.HDel(ctx, key, k1)
		t.AssertNil(err)
		t.Assert(0, r2)

		r3, err := redis.HDel(ctx, key, k2, k3)
		t.AssertNil(err)
		t.Assert(2, r3)
	})
}

func Test_GroupHash_HLen(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key         = "myhash"
			field1      = "field1"
			field1Value = "Hello"
			fields      = map[string]interface{}{
				field1: field1Value,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		fieldLen, err := redis.HLen(ctx, key)
		t.AssertNil(err)
		t.Assert(1, fieldLen)

		fields = map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
		}
		fieldLen, err = redis.HSet(ctx, key, fields)
		t.AssertNil(err)
		t.Assert(2, fieldLen)
	})
}

func Test_GroupHash_HIncrBy(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key         = "myhash"
			field1      = "field1"
			field1Value = 1
			fields      = map[string]interface{}{
				field1: field1Value,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HIncrBy(ctx, key, field1, 2)
		t.AssertNil(err)
		t.Assert(3, r1)

		r2, err := redis.HGet(ctx, key, field1)
		t.AssertNil(err)
		t.Assert(3, r2.Int64())

		r3, err := redis.HIncrBy(ctx, key, field1, -1)
		t.AssertNil(err)
		t.Assert(2, r3)
	})
}

func Test_GroupHash_HIncrByFloat(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key         = "myhash"
			field1      = "field1"
			field1Value = 10.50
			fields      = map[string]interface{}{
				field1: field1Value,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HIncrByFloat(ctx, key, field1, 0.1)
		t.AssertNil(err)
		t.Assert(10.60, r1)

		r2, err := redis.HGet(ctx, key, field1)
		t.AssertNil(err)
		t.Assert(10.60, r2.Float64())

		r3, err := redis.HIncrByFloat(ctx, key, field1, -5)
		t.AssertNil(err)
		t.Assert(5.60, r3)
	})
}

func Test_GroupHash_HMSet(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "myhash"
			k1     = "k1"
			v1     = "v1"
			k2     = "k2"
			v2     = "v2"
			fields = map[string]interface{}{
				k1: v1,
				k2: v2,
			}
		)
		err := redis.HMSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HGet(ctx, key, k1)
		t.AssertNil(err)
		t.Assert(r1.String(), v1)

		r2, err := redis.HGet(ctx, key, k2)
		t.AssertNil(err)
		t.Assert(r2.String(), v2)
	})
}

func Test_GroupHash_HMGet(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "myhash"
			k1     = "k1"
			v1     = "v1"
			k2     = "k2"
			v2     = "v2"
			fields = map[string]interface{}{
				k1: v1,
				k2: v2,
			}
		)
		err := redis.HMSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HMGet(ctx, key, k1, k2)
		t.AssertNil(err)
		t.Assert(r1, []string{v1, v2})
	})
}

func Test_GroupHash_HKeys(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "myhash"
			k1     = "k1"
			v1     = "v1"
			fields = map[string]interface{}{
				k1: v1,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HKeys(ctx, key)
		t.AssertNil(err)
		t.Assert(r1, []string{k1})
	})
}

func Test_GroupHash_HVals(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "myhash"
			k1     = "k1"
			v1     = "v1"
			fields = map[string]interface{}{
				k1: v1,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HVals(ctx, key)
		t.AssertNil(err)
		t.Assert(r1, []string{v1})
	})
}

func Test_GroupHash_HGetAll(t *testing.T) {
	defer redis.FlushAll(ctx)
	gtest.C(t, func(t *gtest.T) {
		var (
			key    = "myhash"
			k1     = "k1"
			v1     = "v1"
			k2     = "k2"
			v2     = "v2"
			fields = map[string]interface{}{
				k1: v1,
				k2: v2,
			}
		)
		_, err := redis.HSet(ctx, key, fields)
		t.AssertNil(err)

		r1, err := redis.HGetAll(ctx, key)
		t.Assert(r1.Map(), fields)
	})
}

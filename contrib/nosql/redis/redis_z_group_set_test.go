// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"testing"
)

func Test_GroupSet_SAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)
		num, err := redis.GroupSet().SAdd(ctx, k1, members...)
		t.AssertNil(err)
		t.Assert(num, 3)
	})
}

func Test_GroupSet_SIsMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num, err := redis.GroupSet().SAdd(ctx, k1, members...)
		t.AssertNil(err)
		t.Assert(num, 3)

		num, err = redis.GroupSet().SIsMember(ctx, k1, "v1")
		t.AssertNil(err)
		t.Assert(1, num)
	})
}

func Test_GroupSet_SPop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num, err := redis.GroupSet().SAdd(ctx, k1, members...)
		t.AssertNil(err)
		t.Assert(num, 3)

		m1, err := redis.GroupSet().SPop(ctx, k1, 2)
		t.AssertNil(err)
		t.AssertIN(m1, []string{"v1", "v2", "v3"})
	})
}

func Test_GroupSet_SRandMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num, err := redis.GroupSet().SAdd(ctx, k1, members...)
		t.AssertNil(err)
		t.Assert(num, 3)

		r, err := redis.GroupSet().SRandMember(ctx, k1, 1)
		t.AssertNil(err)
		t.AssertIN(r, []string{"v1", "v2", "v3"})
	})
}

func Test_GroupSet_SRem(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num, err := redis.GroupSet().SAdd(ctx, k1, members...)
		t.AssertNil(err)
		t.Assert(num, 3)

		n, err := redis.GroupSet().SRem(ctx, k1, "v1")
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

func Test_GroupSet_SMove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}
			k2       = guid.S()
			members2 = []interface{}{
				"v4",
				"v5",
				"v6",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		n, err := redis.GroupSet().SMove(ctx, k1, k2, "v2")
		t.AssertNil(err)
		t.Assert(n, 1)

		m1s, err := redis.GroupSet().SMembers(ctx, k1)
		t.Assert(2, len(m1s))

		m2s, err := redis.GroupSet().SMembers(ctx, k2)
		t.Assert(4, len(m2s))

	})
}

func Test_GroupSet_SCard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		n, err := redis.GroupSet().SCard(ctx, k1)
		t.AssertNil(err)
		t.Assert(n, 3)
	})
}

func Test_GroupSet_SMembers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		r1, err := redis.GroupSet().SMembers(ctx, k1)
		t.AssertNil(err)
		t.AssertIN(r1, []string{"v1", "v2", "v3"})
	})
}

func Test_GroupSet_SMIsMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		//TODO
		_, err = redis.GroupSet().SMIsMember(ctx, k1, "v1")
		t.AssertNil(err)

	})
}

func Test_GroupSet_SInter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v4",
				"v3",
				"v6",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		n, err := redis.GroupSet().SInter(ctx, k1, k2)
		t.AssertNil(err)
		t.AssertIN("v3", n)
		t.AssertNI("v4", n)

	})
}

func Test_GroupSet_SInterStore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v3",
				"v4",
				"v6",
			}

			k3 = guid.S()
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		_, err = redis.GroupSet().SInterStore(ctx, k3, k1, k2)
		t.AssertNil(err)

		member3, err := redis.GroupSet().SMembers(ctx, k3)
		t.AssertNil(err)
		t.AssertIN("v3", member3)
	})
}

func Test_GroupSet_SUnion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v3",
				"v5",
				"v6",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		union, err := redis.GroupSet().SUnion(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(len(union), 5)
	})
}

func Test_GroupSet_SUnionStore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v3",
				"v5",
				"v6",
			}

			k3 = guid.S()
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		union, err := redis.GroupSet().SUnionStore(ctx, k3, k1, k2)
		t.AssertNil(err)

		member3, err := redis.GroupSet().SMembers(ctx, k3)
		t.AssertNil(err)
		t.Assert(len(member3), union)
	})
}

func Test_GroupSet_SDiff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v3",
				"v5",
				"v6",
			}
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		diff, err := redis.GroupSet().SDiff(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(len(diff), 2)
	})
}

func Test_GroupSet_SDiffStore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v1",
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v3",
				"v5",
				"v6",
			}

			k3 = guid.S()
		)

		num1, err := redis.GroupSet().SAdd(ctx, k1, members1...)
		t.AssertNil(err)
		t.Assert(num1, 3)

		num2, err := redis.GroupSet().SAdd(ctx, k2, members2...)
		t.AssertNil(err)
		t.Assert(num2, 3)

		diffStore, err := redis.GroupSet().SDiffStore(ctx, k3, k1, k2)
		t.AssertNil(err)

		members3, err := redis.GroupSet().SMembers(ctx, k3)
		t.AssertNil(err)
		t.Assert(len(members3), diffStore)

	})
}

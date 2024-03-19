// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func TestGroupSetSAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v2",
				"v3",
			}
		)
		num, err := redis.GroupSet().SAdd(ctx, k1, "v1", members...)
		t.Assert(num, 3)
		t.AssertNil(err)
	})
}

func TestGroupSetSIsMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members...)
		t.AssertNil(err)

		num, err := redis.GroupSet().SIsMember(ctx, k1, "v1")
		t.AssertNil(err)
		t.Assert(1, num)
	})
}

func TestGroupSetSPop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members...)
		t.AssertNil(err)

		m1, err := redis.GroupSet().SPop(ctx, k1, 2)
		t.AssertNil(err)
		t.AssertIN(m1, []string{"v1", "v2", "v3"})
	})
}

func TestGroupSetSRandMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members...)
		t.AssertNil(err)

		r, err := redis.GroupSet().SRandMember(ctx, k1, 1)
		t.AssertNil(err)
		t.AssertIN(r, []string{"v1", "v2", "v3"})
	})
}

func TestGroupSetSRem(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1      = guid.S()
			members = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members...)
		t.AssertNil(err)

		n, err := redis.GroupSet().SRem(ctx, k1, "v1")
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

func TestGroupSetSMove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}
			k2       = guid.S()
			members2 = []interface{}{
				"v5",
				"v6",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v4", members2...)
		t.AssertNil(err)

		n, err := redis.GroupSet().SMove(ctx, k1, k2, "v2")
		t.AssertNil(err)
		t.Assert(n, 1)

		m1s, err := redis.GroupSet().SMembers(ctx, k1)
		t.Assert(2, len(m1s))

		m2s, err := redis.GroupSet().SMembers(ctx, k2)
		t.Assert(4, len(m2s))

	})
}

func TestGroupSetSCard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		n, err := redis.GroupSet().SCard(ctx, k1)
		t.AssertNil(err)
		t.Assert(n, 3)
	})
}

func TestGroupSetSMembers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		r1, err := redis.GroupSet().SMembers(ctx, k1)
		t.AssertNil(err)
		t.AssertIN(r1, []string{"v1", "v2", "v3"})
	})
}

func TestGroupSetSMIsMember(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SMIsMember(ctx, k1, "v1")
		t.AssertNil(err)

	})
}

func TestGroupSetSInter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v3",
				"v6",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v4", members2...)
		t.AssertNil(err)

		n, err := redis.GroupSet().SInter(ctx, k1, k2)
		t.AssertNil(err)
		t.AssertIN("v3", n)
		t.AssertNI("v4", n)

	})
}

func TestGroupSetSInterStore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v4",
				"v6",
			}

			k3 = guid.S()
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v3", members2...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SInterStore(ctx, k3, k1, k2)
		t.AssertNil(err)

		member3, err := redis.GroupSet().SMembers(ctx, k3)
		t.AssertNil(err)
		t.AssertIN("v3", member3)
	})
}

func TestGroupSetSUnion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v5",
				"v6",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v3", members2...)
		t.AssertNil(err)

		union, err := redis.GroupSet().SUnion(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(len(union), 5)
	})
}

func TestGroupSetSUnionStore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v5",
				"v6",
			}

			k3 = guid.S()
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v3", members2...)
		t.AssertNil(err)

		union, err := redis.GroupSet().SUnionStore(ctx, k3, k1, k2)
		t.AssertNil(err)

		member3, err := redis.GroupSet().SMembers(ctx, k3)
		t.AssertNil(err)
		t.Assert(len(member3), union)
	})
}

func TestGroupSetSDiff(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v5",
				"v6",
			}
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v3", members2...)
		t.AssertNil(err)

		diff, err := redis.GroupSet().SDiff(ctx, k1, k2)
		t.AssertNil(err)
		t.Assert(len(diff), 2)
	})
}

func TestGroupSetSDiffStore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k1       = guid.S()
			members1 = []interface{}{
				"v2",
				"v3",
			}

			k2       = guid.S()
			members2 = []interface{}{
				"v5",
				"v6",
			}

			k3 = guid.S()
		)

		_, err := redis.GroupSet().SAdd(ctx, k1, "v1", members1...)
		t.AssertNil(err)

		_, err = redis.GroupSet().SAdd(ctx, k2, "v3", members2...)
		t.AssertNil(err)

		diffStore, err := redis.GroupSet().SDiffStore(ctx, k3, k1, k2)
		t.AssertNil(err)

		members3, err := redis.GroupSet().SMembers(ctx, k3)
		t.AssertNil(err)
		t.Assert(len(members3), diffStore)

	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_GroupSortedSet_ZADD(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			k1 = guid.S()
			m1 = guid.S()

			option gredis.ZAddOption
			member gredis.ZAddMember
		)

		member = gredis.ZAddMember{
			Score:  float64(grand.Intn(1000000)),
			Member: m1,
		}

		t.Logf("k1: %s, member: %#v", k1, member)
		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, &option, member)
		t.AssertNil(err)

		// _, err = redis.GroupSortedSet().ZRem(ctx, k1, member.Member)
		// t.AssertNil(err)
	})
}

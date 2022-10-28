// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package group_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Set_Get(t *testing.T) {
	defer redis.FlushDB(ctx)
	gtest.C(t, func(t *gtest.T) {
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
		t.Assert(r2.String(), "v2")
	})
	// With Option.
	gtest.C(t, func(t *gtest.T) {

	})
}

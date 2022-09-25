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
	"github.com/gogf/gf/v2/util/guid"
)

func Test_GroupString_Set_Get(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		var (
			key   = guid.S()
			value = guid.S()
		)
		_, err = redis.String().Set(ctx, key, value)
		t.AssertNil(err)

		v, err := redis.String().Get(ctx, key)
		t.AssertNil(err)
		t.Assert(v.String(), value)
	})
}

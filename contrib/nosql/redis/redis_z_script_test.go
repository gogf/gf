// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	redis2 "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func Test_Script_Eval(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		script := redis2.NewScript(`return ARGV[1]`)
		v, err := script.Run(ctx, redis, nil, "hello")
		t.AssertNil(err)
		t.Assert(v.String(), "hello")
	})
}

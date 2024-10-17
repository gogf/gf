// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"testing"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_mustMergeOptionToArgs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var args []interface{}
		newArgs := mustMergeOptionToArgs(args, gredis.SetOption{
			NX:  true,
			Get: true,
		})
		t.Assert(newArgs, []interface{}{"NX", "GET"})
	})
	gtest.C(t, func(t *gtest.T) {
		var args []interface{}
		newArgs := mustMergeOptionToArgs(args, gredis.SetOption{
			NX:  true,
			Get: true,
			EX:  60,
		})
		t.Assert(newArgs, []interface{}{"NX", "GET", "EX", 60})
	})
}

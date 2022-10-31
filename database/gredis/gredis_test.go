// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_mustMergeOptionToArgs(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var args []interface{}
		newArgs := mustMergeOptionToArgs(args, SetOption{
			NX:  true,
			Get: true,
		})
		t.Assert(newArgs, []interface{}{"NX", "Get"})
	})
	gtest.C(t, func(t *gtest.T) {
		var args []interface{}
		newArgs := mustMergeOptionToArgs(args, SetOption{
			NX:  true,
			Get: true,
			TTLOption: TTLOption{
				EX: gconv.PtrInt64(60),
			},
		})
		t.Assert(newArgs, []interface{}{"EX", 60, "NX", "Get"})
	})
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/gins"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_SetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gins.Set("test-user", 1)
		t.Assert(gins.Get("test-user"), 1)
		t.Assert(gins.Get("none-exists"), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gins.GetOrSet("test-1", 1), 1)
		t.Assert(gins.Get("test-1"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gins.GetOrSetFunc("test-2", func() interface{} {
			return 2
		}), 2)
		t.Assert(gins.Get("test-2"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gins.GetOrSetFuncLock("test-3", func() interface{} {
			return 3
		}), 3)
		t.Assert(gins.Get("test-3"), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gins.SetIfNotExist("test-4", 4), true)
		t.Assert(gins.Get("test-4"), 4)
		t.Assert(gins.SetIfNotExist("test-4", 5), false)
		t.Assert(gins.Get("test-4"), 4)
	})
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins_test

import (
	"github.com/gogf/gf/frame/gins"
	"testing"

	"github.com/gogf/gf/test/gtest"
)

func Test_SetGet(t *testing.T) {
	gtest.Case(t, func() {
		gins.Set("test-user", 1)
		gtest.Assert(gins.Get("test-user"), 1)
		gtest.Assert(gins.Get("none-exists"), nil)
	})
	gtest.Case(t, func() {
		gtest.Assert(gins.GetOrSet("test-1", 1), 1)
		gtest.Assert(gins.Get("test-1"), 1)
	})
	gtest.Case(t, func() {
		gtest.Assert(gins.GetOrSetFunc("test-2", func() interface{} {
			return 2
		}), 2)
		gtest.Assert(gins.Get("test-2"), 2)
	})
	gtest.Case(t, func() {
		gtest.Assert(gins.GetOrSetFuncLock("test-3", func() interface{} {
			return 3
		}), 3)
		gtest.Assert(gins.Get("test-3"), 3)
	})
	gtest.Case(t, func() {
		gtest.Assert(gins.SetIfNotExist("test-4", 4), true)
		gtest.Assert(gins.Get("test-4"), 4)
		gtest.Assert(gins.SetIfNotExist("test-4", 5), false)
		gtest.Assert(gins.Get("test-4"), 4)
	})
}

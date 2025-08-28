// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package instance_test

import (
	"testing"

	"github.com/gogf/gf/v2/internal/instance"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_SetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		instance.Set("test-user", 1)
		t.Assert(instance.Get("test-user"), 1)
		t.Assert(instance.Get("none-exists"), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(instance.GetOrSet("test-1", 1), 1)
		t.Assert(instance.Get("test-1"), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(instance.GetOrSetFunc("test-2", func() any {
			return 2
		}), 2)
		t.Assert(instance.Get("test-2"), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(instance.GetOrSetFuncLock("test-3", func() any {
			return 3
		}), 3)
		t.Assert(instance.Get("test-3"), 3)
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(instance.SetIfNotExist("test-4", 4), true)
		t.Assert(instance.Get("test-4"), 4)
		t.Assert(instance.SetIfNotExist("test-4", 5), false)
		t.Assert(instance.Get("test-4"), 4)
	})
}

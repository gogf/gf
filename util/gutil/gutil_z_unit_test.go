// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"testing"

	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gutil"
)

func Test_Dump(t *testing.T) {
	gtest.Case(t, func() {
		gutil.Dump(map[int]int{
			100: 100,
		})
	})

	gtest.Case(t, func() {
		gutil.Dump(map[string]interface{}{"": func() {}})
	})

	gtest.Case(t, func() {
		gutil.Dump([]byte("gutil Dump test"))
	})
}

func Test_TryCatch(t *testing.T) {

	gtest.Case(t, func() {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")
		})
	})

	gtest.Case(t, func() {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")

		}, func(err interface{}) {
			gtest.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {

	gtest.Case(t, func() {
		gtest.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {

	gtest.Case(t, func() {
		defer func() {
			gtest.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

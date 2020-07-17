// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gutil_test

import (
	"testing"

	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gutil"
)

func Test_Dump(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.Dump(map[int]int{
			100: 100,
		})
	})
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")
		})
	})

	gtest.C(t, func(t *gtest.T) {
		gutil.TryCatch(func() {
			panic("gutil TryCatch test")

		}, func(err interface{}) {
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

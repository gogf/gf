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
	gtest.C(t, func(t *gtest.T) {
		gutil.Dump(map[int]int{
			100: 100,
		})
	})
<<<<<<< HEAD
=======
}

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := `gutil Try test`
		t.Assert(gutil.Try(func() {
			panic(s)
		}), s)
	})
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
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

<<<<<<< HEAD
		}, func(err interface{}) {
=======
		}, func(err error) {
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
			t.Assert(err, "gutil TryCatch test")
		})
	})
}

func Test_IsEmpty(t *testing.T) {
<<<<<<< HEAD

=======
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gutil.IsEmpty(1), false)
	})
}

func Test_Throw(t *testing.T) {
<<<<<<< HEAD

=======
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			t.Assert(recover(), "gutil Throw test")
		}()

		gutil.Throw("gutil Throw test")
	})
}

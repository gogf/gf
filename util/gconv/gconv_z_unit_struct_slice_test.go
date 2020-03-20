// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func Test_Struct_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []int
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []int32
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []int64
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []uint
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []uint32
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []uint64
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []float32
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Scores []float64
		}
		user := new(User)
		array := g.Slice{1, 2, 3}
		err := gconv.Struct(g.Map{"scores": array}, user)
		t.Assert(err, nil)
		t.Assert(user.Scores, array)
	})
}

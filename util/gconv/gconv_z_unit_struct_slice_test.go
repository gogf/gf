// Copyright 2019 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gconv_test

import (
	"testing"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gconv"
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

func Test_Struct_SliceWithTag(t *testing.T) {
	type User struct {
		Uid      int    `json:"id"`
		NickName string `json:"name"`
	}
	gtest.C(t, func(t *gtest.T) {
		var users []User
		params := g.Slice{
			g.Map{
				"id":   1,
				"name": "name1",
			},
			g.Map{
				"id":   2,
				"name": "name2",
			},
		}
		err := gconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		params := g.Slice{
			g.Map{
				"id":   1,
				"name": "name1",
			},
			g.Map{
				"id":   2,
				"name": "name2",
			},
		}
		err := gconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Uid, 1)
		t.Assert(users[0].NickName, "name1")
		t.Assert(users[1].Uid, 2)
		t.Assert(users[1].NickName, "name2")
	})
}

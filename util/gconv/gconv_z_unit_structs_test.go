// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Structs_WithTag(t *testing.T) {
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

func Test_Structs_WithoutTag(t *testing.T) {
	type User struct {
		Uid      int
		NickName string
	}
	gtest.C(t, func(t *gtest.T) {
		var users []User
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
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
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
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

func Test_Structs_SliceParameter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			NickName string
		}
		var users []User
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := gconv.Structs(params, users)
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			NickName string
		}
		type A struct {
			Users []User
		}
		var a A
		params := g.Slice{
			g.Map{
				"uid":       1,
				"nick-name": "name1",
			},
			g.Map{
				"uid":       2,
				"nick-name": "name2",
			},
		}
		err := gconv.Structs(params, a.Users)
		t.AssertNE(err, nil)
	})
}

func Test_Structs_DirectReflectSet(t *testing.T) {
	type A struct {
		Id   int
		Name string
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			a = []*A{
				{Id: 1, Name: "john"},
				{Id: 2, Name: "smith"},
			}
			b []*A
		)
		err := gconv.Structs(a, &b)
		t.Assert(err, nil)
		t.AssertEQ(a, b)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			a = []A{
				{Id: 1, Name: "john"},
				{Id: 2, Name: "smith"},
			}
			b []A
		)
		err := gconv.Structs(a, &b)
		t.Assert(err, nil)
		t.AssertEQ(a, b)
	})
}

func Test_Structs_IntSliceAttribute(t *testing.T) {
	type A struct {
		Id []int
	}
	type B struct {
		*A
		Name string
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			array []*B
		)
		err := gconv.Structs(g.Slice{
			g.Map{"id": nil, "name": "john"},
			g.Map{"id": nil, "name": "smith"},
		}, &array)
		t.Assert(err, nil)
		t.Assert(len(array), 2)
		t.Assert(array[0].Name, "john")
		t.Assert(array[1].Name, "smith")
	})
}

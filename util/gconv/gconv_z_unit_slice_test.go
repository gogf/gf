// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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

func Test_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := 123.456
		t.AssertEQ(gconv.Bytes("123"), []byte("123"))
		t.AssertEQ(gconv.Strings(value), []string{"123.456"})
		t.AssertEQ(gconv.Ints(value), []int{123})
		t.AssertEQ(gconv.Floats(value), []float64{123.456})
		t.AssertEQ(gconv.Interfaces(value), []interface{}{123.456})
	})
}

func Test_Strings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := []*g.Var{
			g.NewVar(1),
			g.NewVar(2),
			g.NewVar(3),
		}
		t.AssertEQ(gconv.Strings(array), []string{"1", "2", "3"})
	})
}

func Test_Slice_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int    `json:"id"`
		name string `json:"name"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{1, "john"}
		t.Assert(gconv.Interfaces(user), g.Slice{"id", 1})
	})
}

func Test_Slice_Structs(t *testing.T) {
	type Base struct {
		Age int
	}
	type User struct {
		Id   int
		Name string
		Base
	}

	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gconv.Structs(params, &users)
		t.Assert(err, nil)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})
}

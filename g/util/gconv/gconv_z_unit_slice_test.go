// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gconv"
)

func Test_Slice(t *testing.T) {
	gtest.Case(t, func() {
		value := 123.456
		gtest.AssertEQ(gconv.Bytes("123"), []byte("123"))
		gtest.AssertEQ(gconv.Strings(value), []string{"123.456"})
		gtest.AssertEQ(gconv.Ints(value), []int{123})
		gtest.AssertEQ(gconv.Floats(value), []float64{123.456})
		gtest.AssertEQ(gconv.Interfaces(value), []interface{}{123.456})
	})
}

// 私有属性不会进行转换
func Test_Slice_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gtest.Case(t, func() {
		user := &User{1, "john"}
		gtest.Assert(gconv.Interfaces(user), g.Slice{1})
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

	gtest.Case(t, func() {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gconv.Structs(params, &users)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, params[0]["id"])
		gtest.Assert(users[0].Name, params[0]["name"])
		gtest.Assert(users[0].Age, 0)

		gtest.Assert(users[1].Id, params[1]["id"])
		gtest.Assert(users[1].Name, params[1]["name"])
		gtest.Assert(users[1].Age, 0)
	})

	gtest.Case(t, func() {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gconv.StructsDeep(params, &users)
		gtest.Assert(err, nil)
		gtest.Assert(len(users), 2)
		gtest.Assert(users[0].Id, params[0]["id"])
		gtest.Assert(users[0].Name, params[0]["name"])
		gtest.Assert(users[0].Age, params[0]["age"])

		gtest.Assert(users[1].Id, params[1]["id"])
		gtest.Assert(users[1].Name, params[1]["name"])
		gtest.Assert(users[1].Age, params[1]["age"])
	})
}

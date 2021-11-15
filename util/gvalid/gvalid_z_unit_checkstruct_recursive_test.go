// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gvalid"
)

func Test_CheckStruct_Recursive_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id   int
			Name string `v:"required"`
			Pass Pass
		}
		user := &User{
			Name: "",
			Pass: Pass{
				Pass1: "1",
				Pass2: "2",
			},
		}
		err := gvalid.CheckStruct(context.TODO(), user, nil)
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], g.Map{"required": "The Name field is required"})
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `1` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `2` must be the same as field Pass1"})
	})
}

func Test_CheckStruct_Recursive_Struct_WithData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id   int
			Name string `v:"required"`
			Pass Pass
		}
		user := &User{}
		data := g.Map{
			"Name": "john",
			"Pass": g.Map{
				"Pass1": 100,
				"Pass2": 200,
			},
		}
		err := g.Validator().Data(data).CheckStruct(context.TODO(), user)
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], nil)
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `100` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `200` must be the same as field Pass1"})
	})
}

func Test_CheckStruct_Recursive_SliceStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id     int
			Name   string `v:"required"`
			Passes []Pass
		}
		user := &User{
			Name: "",
			Passes: []Pass{
				{
					Pass1: "1",
					Pass2: "2",
				},
				{
					Pass1: "3",
					Pass2: "4",
				},
			},
		}
		err := gvalid.CheckStruct(context.TODO(), user, nil)
		g.Dump(err.Items())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], g.Map{"required": "The Name field is required"})
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `1` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `2` must be the same as field Pass1"})
	})
}

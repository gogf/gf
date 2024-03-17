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

func Test_Struct_Custom_Mapping1_Attribute(t *testing.T) {
	type TestData struct {
		Name string
		Age  string
	}
	type dataMapKey = string
	type structFieldName = string
	gtest.C(t, func(t *gtest.T) {
		data := map[string]any{
			"a":    "123",
			"name": "456",
		}
		mapping := map[dataMapKey]structFieldName{
			"a":    "Age",
			"name": "Age",

			"age": "Name",
			"n":   "Name",
		}

		var input = &TestData{}
		err := gconv.Struct(data, input, mapping)
		t.AssertNil(err)
		t.AssertEQ(input.Age, "123")
		t.AssertNE(input.Name, "456")
	})
}

func Test_Struct_Custom_Mapping2_Attribute(t *testing.T) {
	type User struct {
		Uid  int
		Name string
	}
	gtest.C(t, func(t *gtest.T) {

		var (
			user   = new(User)
			params = g.Map{
				"uid": 1,
				//"myname": "john",
				"name": "smith",
			}
		)
		err := gconv.Scan(params, user, g.MapStrStr{
			"myname": "Name",
		})

		t.AssertNil(err)
		t.Assert(user, &User{
			Uid:  1,
			Name: "smith",
		})
	})

}

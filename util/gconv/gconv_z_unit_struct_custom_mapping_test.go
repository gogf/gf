// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Struct_Custom_Mapping_Attribute(t *testing.T) {
	type TestData struct {
		Name string
		Age  string
		Sex  string
	}
	data := map[string]any{
		"age":  "123",
		"name": "123",
	}
	type dataMapKey = string
	type structFieldName = string
	mapping := map[dataMapKey]structFieldName{
		"age":  "Age",
		"name": "Sex",
	}

	gtest.C(t, func(t *gtest.T) {
		var input = &TestData{
			Name: "",
			Age:  "",
		}
		err := gconv.Struct(data, input, mapping)
		t.AssertNil(err)
		t.AssertEQ(input, &TestData{
			Age: "123",
			Sex: "123",
		})
	})

}

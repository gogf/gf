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

func Test_StructTag(t *testing.T) {
	type User struct {
		Uid   int
		Name  string
		Pass1 string `orm:"password1"`
		Pass2 string `orm:"password2"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		params1 := g.Map{
			"uid":       1,
			"Name":      "john",
			"password1": "123",
			"password2": "456",
		}
		if err := gconv.Struct(params1, user); err != nil {
			t.Error(err)
		}
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "",
			Pass2: "",
		})
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		params1 := g.Map{
			"uid":       1,
			"Name":      "john",
			"password1": "123",
			"password2": "456",
		}
		if err := gconv.StructTag(params1, user, "orm"); err != nil {
			t.Error(err)
		}
		t.Assert(user, &User{
			Uid:   1,
			Name:  "john",
			Pass1: "123",
			Pass2: "456",
		})
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		params2 := g.Map{
			"uid":       2,
			"name":      "smith",
			"password1": "111",
			"password2": "222",
		}
		if err := gconv.StructTag(params2, user, "orm"); err != nil {
			t.Error(err)
		}
		t.Assert(user, &User{
			Uid:   2,
			Name:  "smith",
			Pass1: "111",
			Pass2: "222",
		})
	})
}

func Test_Struct_HTTP_Params_Parse_Issue1488(t *testing.T) {
	type Request struct {
		Type         []int  `p:"type"`
		Keyword      string `p:"keyword"`
		Limit        int    `p:"per_page" d:"10"`
		Page         int    `p:"page" d:"1"`
		Order        string
		CreatedAtLte string
		CreatedAtGte string
		CreatorID    []int
	}
	data := g.Map{
		"type[]":   []int{1, 2, 3},
		"Limit":    10,
		"per_page": 6,
	}
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i < 10; i++ {
			var req = &Request{}
			err := gconv.Struct(data, req)
			t.AssertNil(err)
			// maybe 10
			t.AssertEQ(req.Limit, 6)
		}
	})
}

func Test_StructTag_MultiAttribute(t *testing.T) {
	type ConfigMap struct {
		TypeMeta   string `json:",inline"`
		ObjectMeta string `json:" metadata,omitempty" `
		Immutable  bool   `json:" immutable,omitempty" `
	}

	gtest.C(t, func(t *gtest.T) {
		data := g.Map{
			"TypeMeta":  "1",
			"metadata":  "2",
			"immutable": true,
		}
		var cfg ConfigMap
		err := gconv.Struct(data, &cfg)
		t.AssertNil(err)

		t.Assert(cfg, &ConfigMap{
			TypeMeta:   "1",
			ObjectMeta: "2",
			Immutable:  true,
		})
	})

}

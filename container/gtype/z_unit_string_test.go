// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
)

func Test_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := gtype.NewString("abc")
		iClone := i.Clone()
		t.AssertEQ(iClone.Set("123"), "abc")
		t.AssertEQ(iClone.Val(), "123")

		//空参测试
		i1 := gtype.NewString()
		t.AssertEQ(i1.Val(), "")
	})
}

func Test_String_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "i love gf"
		i1 := gtype.NewString(s)
		b1, err1 := json.Marshal(i1)
		b2, err2 := json.Marshal(i1.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.NewString()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), s)
	})
}

func Test_String_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gtype.String
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}

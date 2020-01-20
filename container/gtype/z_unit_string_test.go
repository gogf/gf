// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"encoding/json"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
)

func Test_String(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewString("abc")
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set("123"), "abc")
		gtest.AssertEQ(iClone.Val(), "123")

		//空参测试
		i1 := gtype.NewString()
		gtest.AssertEQ(i1.Val(), "")
	})
}

func Test_String_JSON(t *testing.T) {
	gtest.Case(t, func() {
		s := "i love gf"
		i1 := gtype.NewString(s)
		b1, err1 := json.Marshal(i1)
		b2, err2 := json.Marshal(i1.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewString()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), s)
	})
}

func Test_String_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.String
	}
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Var.Val(), "123")
	})
}

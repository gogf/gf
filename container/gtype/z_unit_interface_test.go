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

func Test_Interface(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t := Temp{Name: "gf", Age: 18}
		t1 := Temp{Name: "gf", Age: 19}
		i := gtype.New(t)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(t1), t)
		t.AssertEQ(iClone.Val().(Temp), t1)

		//空参测试
		i1 := gtype.New()
		t.AssertEQ(i1.Val(), nil)
	})
}

func Test_Interface_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "i love gf"
		i := gtype.New(s)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.New()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), s)
	})
}

func Test_Interface_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.Interface
	}
	gtest.C(t, func(t *gtest.T) {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &t)
		t.Assert(err, nil)
		t.Assert(t.Name, "john")
		t.Assert(t.Var.Val(), "123")
	})
}

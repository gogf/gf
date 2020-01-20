// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"testing"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/test/gtest"
)

func Test_Bool(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(false), true)
		gtest.AssertEQ(iClone.Val(), false)

		i1 := gtype.NewBool(false)
		iClone1 := i1.Clone()
		gtest.AssertEQ(iClone1.Set(true), false)
		gtest.AssertEQ(iClone1.Val(), true)

		//空参测试
		i2 := gtype.NewBool()
		gtest.AssertEQ(i2.Val(), false)
	})
}

func Test_Bool_JSON(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)
	})
	gtest.Case(t, func() {
		i := gtype.NewBool(false)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		var err error
		i := gtype.NewBool()
		err = json.Unmarshal([]byte("true"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), true)
		err = json.Unmarshal([]byte("false"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), false)
		err = json.Unmarshal([]byte("1"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), true)
		err = json.Unmarshal([]byte("0"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), false)
	})

	gtest.Case(t, func() {
		i := gtype.NewBool(true)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewBool()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i.Val())
	})
	gtest.Case(t, func() {
		i := gtype.NewBool(false)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewBool()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i.Val())
	})
}

func Test_Bool_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.Bool
	}
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "true",
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Var.Val(), true)
	})
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "false",
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Var.Val(), false)
	})
}

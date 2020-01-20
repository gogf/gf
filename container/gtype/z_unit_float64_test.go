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
	"math"
	"testing"
)

func Test_Float64(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewFloat64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(0.1), float64(0))
		gtest.AssertEQ(iClone.Val(), float64(0.1))
		//空参测试
		i1 := gtype.NewFloat64()
		gtest.AssertEQ(i1.Val(), float64(0))
	})
}

func Test_Float64_JSON(t *testing.T) {
	gtest.Case(t, func() {
		v := math.MaxFloat64
		i := gtype.NewFloat64(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewFloat64()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), v)
	})
}

func Test_Float64_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.Float64
	}
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123.456",
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Var.Val(), "123.456")
	})
}

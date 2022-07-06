// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"math"
	"testing"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Float64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := gtype.NewFloat64(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(0.1), float64(0))
		t.AssertEQ(iClone.Val(), float64(0.1))
		// 空参测试
		i1 := gtype.NewFloat64()
		t.AssertEQ(i1.Val(), float64(0))

		//
		i2 := gtype.NewFloat64(1.1)
		t.AssertEQ(i2.Add(3.3), 4.4)
		t.AssertEQ(i2.Cas(4.5, 5.5), false)
		t.AssertEQ(i2.Cas(4.4, 5.5), true)
		t.AssertEQ(i2.String(), "5.5")

		copyVal := i2.DeepCopy()
		i2.Set(6.6)
		t.AssertNE(copyVal, iClone.Val())
		i2 = nil
		copyVal = i2.DeepCopy()
		t.AssertNil(copyVal)
	})
}

func Test_Float64_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := math.MaxFloat64
		i := gtype.NewFloat64(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.NewFloat64()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.AssertNil(err)
		t.Assert(i2.Val(), v)
	})
}

func Test_Float64_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gtype.Float64
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123.456",
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123.456")
	})
}

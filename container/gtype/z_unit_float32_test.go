// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gtype_test

import (
	"github.com/jin502437344/gf/container/gtype"
	"github.com/jin502437344/gf/internal/json"
	"github.com/jin502437344/gf/test/gtest"
	"github.com/jin502437344/gf/util/gconv"
	"math"
	"testing"
)

func Test_Float32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := gtype.NewFloat32(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(0.1), float32(0))
		t.AssertEQ(iClone.Val(), float32(0.1))

		//空参测试
		i1 := gtype.NewFloat32()
		t.AssertEQ(i1.Val(), float32(0))
	})
}

func Test_Float32_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := float32(math.MaxFloat32)
		i := gtype.NewFloat32(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())

		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.NewFloat32()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}

func Test_Float32_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gtype.Float32
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123.456",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123.456")
	})
}

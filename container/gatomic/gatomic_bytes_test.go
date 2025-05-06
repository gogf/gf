// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gatomic_test

import (
	"testing"

	"github.com/gogf/gf/v3/container/gatomic"
	"github.com/gogf/gf/v3/internal/json"
	"github.com/gogf/gf/v3/test/gtest"
	"github.com/gogf/gf/v3/util/gconv"
)

func Test_Bytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := gatomic.NewBytes([]byte("abc"))
		iClone := i.Clone()
		t.AssertEQ(iClone.Set([]byte("123")), []byte("abc"))
		t.AssertEQ(iClone.Val(), []byte("123"))

		// empty param test
		i1 := gatomic.NewBytes()
		t.AssertEQ(i1.Val(), nil)

		i2 := gatomic.NewBytes([]byte("abc"))
		t.Assert(i2.String(), "abc")

		copyVal := i2.DeepCopy()
		i2.Set([]byte("def"))
		t.AssertNE(copyVal, iClone.Val())
		i2 = nil
		copyVal = i2.DeepCopy()
		t.AssertNil(copyVal)
	})
}

func Test_Bytes_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		b := []byte("i love gf")
		i := gatomic.NewBytes(b)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gatomic.NewBytes()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.AssertNil(err)
		t.Assert(i2.Val(), b)
	})
}

func Test_Bytes_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gatomic.Bytes
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}

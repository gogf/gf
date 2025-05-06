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

func Test_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := gatomic.NewString("abc")
		iClone := i.Clone()
		t.AssertEQ(iClone.Set("123"), "abc")
		t.AssertEQ(iClone.Val(), "123")
		t.AssertEQ(iClone.String(), "123")
		//
		copyVal := iClone.DeepCopy()
		iClone.Set("124")
		t.AssertNE(copyVal, iClone.Val())
		iClone = nil
		copyVal = iClone.DeepCopy()
		t.AssertNil(copyVal)
		// empty param test
		i1 := gatomic.NewString()
		t.AssertEQ(i1.Val(), "")
	})
}

func Test_String_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := "i love gf"
		i1 := gatomic.NewString(s)
		b1, err1 := json.Marshal(i1)
		b2, err2 := json.Marshal(i1.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gatomic.NewString()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.AssertNil(err)
		t.Assert(i2.Val(), s)
	})
}

func Test_String_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gatomic.String
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

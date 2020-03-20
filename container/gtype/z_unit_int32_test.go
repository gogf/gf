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
	"sync"
	"testing"
)

func Test_Int32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt32(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(1), int32(0))
		t.AssertEQ(iClone.Val(), int32(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		t.AssertEQ(int32(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewInt32()
		t.AssertEQ(i1.Val(), int32(0))
	})
}

func Test_Int32_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := int32(math.MaxInt32)
		i := gtype.NewInt32(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.NewInt32()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}

func Test_Int32_UnmarshalValue(t *testing.T) {
	type Var struct {
		Name string
		Var  *gtype.Int32
	}
	gtest.C(t, func(t *gtest.T) {
		var v *Var
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}

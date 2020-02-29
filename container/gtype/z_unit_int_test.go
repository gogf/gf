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
	"sync"
	"testing"
)

func Test_Int(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), 0)
		gtest.AssertEQ(iClone.Val(), 1)
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(addTimes, i.Val())

		//空参测试
		i1 := gtype.NewInt()
		gtest.AssertEQ(i1.Val(), 0)
	})
}

func Test_Int_JSON(t *testing.T) {
	gtest.Case(t, func() {
		v := 666
		i := gtype.NewInt(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewInt()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), v)
	})
}

func Test_Int_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.Int
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

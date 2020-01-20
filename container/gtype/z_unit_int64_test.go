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

func Test_Int64(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt64(0)
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(1), int64(0))
		gtest.AssertEQ(iClone.Val(), int64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(int64(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewInt64()
		gtest.AssertEQ(i1.Val(), int64(0))
	})
}

func Test_Int64_JSON(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewInt64(math.MaxInt64)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)

		i2 := gtype.NewInt64()
		err := json.Unmarshal(b2, &i2)
		gtest.Assert(err, nil)
		gtest.Assert(i2.Val(), i)
	})
}

func Test_Int64_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.Int64
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

// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtype_test

import (
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"sync"
	"testing"

	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/test/gtest"
)

func Test_Byte(t *testing.T) {
	gtest.Case(t, func() {
		var wg sync.WaitGroup
		addTimes := 127
		i := gtype.NewByte(byte(0))
		iClone := i.Clone()
		gtest.AssertEQ(iClone.Set(byte(1)), byte(0))
		gtest.AssertEQ(iClone.Val(), byte(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		gtest.AssertEQ(byte(addTimes), i.Val())

		//空参测试
		i1 := gtype.NewByte()
		gtest.AssertEQ(i1.Val(), byte(0))
	})
}

func Test_Byte_JSON(t *testing.T) {
	gtest.Case(t, func() {
		i := gtype.NewByte(49)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
		gtest.Assert(b1, b2)
	})
	// Unmarshal
	gtest.Case(t, func() {
		var err error
		i := gtype.NewByte()
		err = json.Unmarshal([]byte("49"), &i)
		gtest.Assert(err, nil)
		gtest.Assert(i.Val(), "49")
	})
}

func Test_Byte_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gtype.Byte
	}
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "2",
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Var.Val(), "2")
	})
}

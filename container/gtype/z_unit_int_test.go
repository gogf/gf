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
	"sync"
	"testing"
)

func Test_Int(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gtype.NewInt(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(1), 0)
		t.AssertEQ(iClone.Val(), 1)
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		t.AssertEQ(addTimes, i.Val())

		//空参测试
		i1 := gtype.NewInt()
		t.AssertEQ(i1.Val(), 0)
	})
}

func Test_Int_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := 666
		i := gtype.NewInt(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.NewInt()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}

func Test_Int_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gtype.Int
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "123",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.Val(), "123")
	})
}

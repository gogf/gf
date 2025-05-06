// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gatomic_test

import (
	"math"
	"sync"
	"testing"

	"github.com/gogf/gf/v3/container/gatomic"
	"github.com/gogf/gf/v3/internal/json"
	"github.com/gogf/gf/v3/test/gtest"
	"github.com/gogf/gf/v3/util/gconv"
)

type Temp struct {
	Name string
	Age  int
}

func Test_Uint64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var wg sync.WaitGroup
		addTimes := 1000
		i := gatomic.NewUint64(0)
		iClone := i.Clone()
		t.AssertEQ(iClone.Set(1), uint64(0))
		t.AssertEQ(iClone.Val(), uint64(1))
		for index := 0; index < addTimes; index++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				i.Add(1)
			}()
		}
		wg.Wait()
		t.AssertEQ(uint64(addTimes), i.Val())

		// empty param test
		i1 := gatomic.NewUint64()
		t.AssertEQ(i1.Val(), uint64(0))

		i2 := gatomic.NewUint64(11)
		t.AssertEQ(i2.Add(1), uint64(12))
		t.AssertEQ(i2.Cas(11, 13), false)
		t.AssertEQ(i2.Cas(12, 13), true)
		t.AssertEQ(i2.String(), "13")

		copyVal := i2.DeepCopy()
		i2.Set(14)
		t.AssertNE(copyVal, iClone.Val())
		i2 = nil
		copyVal = i2.DeepCopy()
		t.AssertNil(copyVal)
	})
}

func Test_Uint64_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		i := gatomic.NewUint64(math.MaxUint64)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gatomic.NewUint64()
		err := json.UnmarshalUseNumber(b2, &i2)
		t.AssertNil(err)
		t.Assert(i2.Val(), i)
	})
}

func Test_Uint64_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gatomic.Uint64
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

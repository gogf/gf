// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gring_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gring"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

type Student struct {
	position int
	name     string
	upgrade  bool
}

func TestRing_Val(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		//定义cap 为3的ring类型数据
		r := gring.New(3, true)
		//分别给3个元素初始化赋值
		r.Put(&Student{1, "jimmy", true})
		r.Put(&Student{2, "tom", true})
		r.Put(&Student{3, "alon", false})

		//元素取值并判断和预设值是否相等
		t.Assert(r.Val().(*Student).name, "jimmy")
		//从当前位置往后移两个元素
		r.Move(2)
		t.Assert(r.Val().(*Student).name, "alon")
		//更新元素值
		//测试 value == nil
		r.Set(nil)
		t.Assert(r.Val(), nil)
		//测试value != nil
		r.Set(&Student{3, "jack", true})
	})
}
func TestRing_CapLen(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r := gring.New(10)
		t.Assert(r.Cap(), 10)
		t.Assert(r.Len(), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		r := gring.New(10)
		r.Put("goframe")
		//cap长度 10
		t.Assert(r.Cap(), 10)
		//已有数据项 1
		t.Assert(r.Len(), 1)
	})
}

func TestRing_Position(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r := gring.New(2)
		r.Put(1)
		r.Put(2)
		//往后移动1个元素
		r.Next()
		t.Assert(r.Val(), 2)
		//往前移动1个元素
		r.Prev()
		t.Assert(r.Val(), 1)

	})
}

func TestRing_Link(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r := gring.New(3)
		r.Put(1)
		r.Put(2)
		r.Put(3)
		s := gring.New(2)
		s.Put("a")
		s.Put("b")

		rs := r.Link(s)
		t.Assert(rs.Move(2).Val(), "b")
	})
}

func TestRing_Unlink(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r := gring.New(5)
		for i := 1; i <= 5; i++ {
			r.Put(i)
		}
		t.Assert(r.Val(), 1)
		// 1 2 3 4
		// 删除当前位置往后的2个数据，返回被删除的数据
		// 重新计算s len
		s := r.Unlink(2) // 2 3
		t.Assert(s.Val(), 2)
		t.Assert(s.Len(), 2)
	})
}

func TestRing_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ringLen := 5
		r := gring.New(ringLen)
		for i := 0; i < ringLen; i++ {
			r.Put(i + 1)
		}
		r.Move(2)              // 3
		array := r.SliceNext() // [3 4 5 1 2]
		t.Assert(array[0], 3)
		t.Assert(len(array), 5)

		//判断array是否等于[3 4 5 1 2]
		ra := []int{3, 4, 5, 1, 2}
		t.Assert(ra, array)

		//第3个元素设为nil
		r.Set(nil)
		array2 := r.SliceNext() //[4 5 1 2]
		//返回当前位置往后不为空的元素数组，长度为4
		t.Assert(array2, g.Slice{nil, 4, 5, 1, 2})

		array3 := r.SlicePrev() //[2 1 5 4]
		t.Assert(array3, g.Slice{nil, 2, 1, 5, 4})

		s := gring.New(ringLen)
		for i := 0; i < ringLen; i++ {
			s.Put(i + 1)
		}
		array4 := s.SlicePrev() // []
		t.Assert(array4, g.Slice{1, 5, 4, 3, 2})
	})
}

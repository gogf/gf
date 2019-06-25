// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray_test

import (
	"github.com/gogf/gf/g/container/garray"
	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gconv"
	"strings"
	"testing"
	"time"
)

func Test_IntArray_Basic(t *testing.T) {
	gtest.Case(t, func() {
		expect := []int{0, 1, 2, 3}
		expect2 := []int{}
		array := garray.NewIntArrayFrom(expect)
		array2 := garray.NewIntArrayFrom(expect2)
		gtest.Assert(array.Slice(), expect)
		array.Set(0, 100)
		gtest.Assert(array.Get(0), 100)
		gtest.Assert(array.Get(1), 1)
		gtest.Assert(array.Search(100), 0)
		gtest.Assert(array2.Search(7), -1)
		gtest.Assert(array.Contains(100), true)
		gtest.Assert(array.Remove(0), 100)
		gtest.Assert(array.Contains(100), false)
		array.Append(4)
		gtest.Assert(array.Len(), 4)
		array.InsertBefore(0, 100)
		array.InsertAfter(0, 200)
		gtest.Assert(array.Slice(), []int{100, 200, 1, 2, 3, 4})
		array.InsertBefore(5, 300)
		array.InsertAfter(6, 400)
		gtest.Assert(array.Slice(), []int{100, 200, 1, 2, 3, 300, 4, 400})
		gtest.Assert(array.Clear().Len(), 0)
	})
}

func TestIntArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		expect1 := []int{0, 1, 2, 3}
		expect2 := []int{3, 2, 1, 0}
		array := garray.NewIntArray()
		array2 := garray.NewIntArray(true)
		for i := 3; i >= 0; i-- {
			array.Append(i)
			array2.Append(i)
		}

		array.Sort()
		gtest.Assert(array.Slice(), expect1)
		array.Sort(true)
		gtest.Assert(array.Slice(), expect2)
		array2.Sort(true)
		gtest.Assert(array2.Slice(), expect2)

	})
}

func TestIntArray_Unique(t *testing.T) {
	gtest.Case(t, func() {
		expect := []int{1, 1, 2, 3}
		array := garray.NewIntArrayFrom(expect)
		gtest.Assert(array.Unique().Slice(), []int{1, 2, 3})
	})
}

func TestIntArray_PushAndPop(t *testing.T) {
	gtest.Case(t, func() {
		expect := []int{0, 1, 2, 3}
		array := garray.NewIntArrayFrom(expect)
		gtest.Assert(array.Slice(), expect)
		gtest.Assert(array.PopLeft(), 0)
		gtest.Assert(array.PopRight(), 3)
		gtest.AssertIN(array.PopRand(), []int{1, 2})
		gtest.AssertIN(array.PopRand(), []int{1, 2})
		gtest.Assert(array.Len(), 0)
		array.PushLeft(1).PushRight(2)
		gtest.Assert(array.Slice(), []int{1, 2})
	})
}

func TestIntArray_PopLeftsAndPopRights(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []int{0, 1, 2, 3, 4, 5, 6}
		value2 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(value1)
		array2 := garray.NewIntArrayFrom(value2)
		gtest.Assert(array1.PopLefts(2), []int{0, 1})
		gtest.Assert(array1.Slice(), []int{2, 3, 4, 5, 6})
		gtest.Assert(array1.PopRights(2), []int{5, 6})
		gtest.Assert(array1.Slice(), []int{2, 3, 4})
		gtest.Assert(array1.PopRights(20), []int{2, 3, 4})
		gtest.Assert(array1.Slice(), []int{})
		gtest.Assert(array2.PopLefts(20), []int{0, 1, 2, 3, 4, 5, 6})
		gtest.Assert(array2.Slice(), []int{})
	})
}

func TestIntArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(value1)
		array2 := garray.NewIntArrayFrom(value1, true)
		gtest.Assert(array1.Range(0, 1), []int{0})
		gtest.Assert(array1.Range(1, 2), []int{1})
		gtest.Assert(array1.Range(0, 2), []int{0, 1})
		gtest.Assert(array1.Range(10, 2), nil)
		gtest.Assert(array1.Range(-1, 10), value1)
		gtest.Assert(array1.Range(8, 2), nil)

		gtest.Assert(array2.Range(2, 4), []int{2, 3})
	})
}

func TestIntArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		n1 := []int{1, 2, 4, 3}
		n2 := []int{7, 8, 9}
		n3 := []int{3, 6}

		s1 := []string{"a", "b", "c"}
		in1 := []interface{}{1, "a", 2, "b"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}

		a1 := garray.NewIntArrayFrom(n1)
		b1 := garray.NewStringArrayFrom(s1)
		b2 := garray.NewIntArrayFrom(n3)
		b3 := garray.NewArrayFrom(in1)
		b4 := garray.NewSortedStringArrayFrom(s1)
		b5 := garray.NewSortedIntArrayFrom(n3)
		b6 := garray.NewSortedArrayFrom(in1, func1)

		gtest.Assert(a1.Merge(n2).Len(), 7)
		gtest.Assert(a1.Merge(n3).Len(), 9)
		gtest.Assert(a1.Merge(b1).Len(), 12)
		gtest.Assert(a1.Merge(b2).Len(), 14)
		gtest.Assert(a1.Merge(b3).Len(), 18)
		gtest.Assert(a1.Merge(b4).Len(), 21)
		gtest.Assert(a1.Merge(b5).Len(), 23)
		gtest.Assert(a1.Merge(b6).Len(), 27)
	})
}

func TestIntArray_Fill(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0}
		a2 := []int{0}
		array1 := garray.NewIntArrayFrom(a1)
		array2 := garray.NewIntArrayFrom(a2)
		gtest.Assert(array1.Fill(1, 2, 100).Slice(), []int{0, 100, 100})
		gtest.Assert(array2.Fill(0, 2, 100).Slice(), []int{100, 100})
		gtest.Assert(array2.Fill(-1, 2, 100).Slice(), []int{100, 100})
	})
}

func TestIntArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewIntArrayFrom(a1)
		chunks := array1.Chunk(2)
		chunks2 := array1.Chunk(0)
		gtest.Assert(chunks2, nil)
		gtest.Assert(chunks[0], []int{1, 2})
		gtest.Assert(chunks[1], []int{3, 4})
		gtest.Assert(chunks[2], []int{5})
	})
}

func TestIntArray_Pad(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Pad(3, 1).Slice(), []int{0, 1, 1})
		gtest.Assert(array1.Pad(-4, 1).Slice(), []int{1, 0, 1, 1})
		gtest.Assert(array1.Pad(3, 1).Slice(), []int{1, 0, 1, 1})
	})
}

func TestIntArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		array2 := garray.NewIntArrayFrom(a1,true)
		gtest.Assert(array1.SubSlice(6), []int{6})
		gtest.Assert(array1.SubSlice(5), []int{5, 6})
		gtest.Assert(array1.SubSlice(8), nil)
		gtest.Assert(array1.SubSlice(0, 2), []int{0, 1})
		gtest.Assert(array1.SubSlice(2, 2), []int{2, 3})
		gtest.Assert(array1.SubSlice(5, 8), []int{5, 6})
		gtest.Assert(array1.SubSlice(-1, 1), []int{6})
		gtest.Assert(array1.SubSlice(-1, 9), []int{6})
		gtest.Assert(array1.SubSlice(-2, 3), []int{5, 6})
		gtest.Assert(array1.SubSlice(-7, 3), []int{0, 1, 2})
		gtest.Assert(array1.SubSlice(-8, 3), nil)
		gtest.Assert(array1.SubSlice(-1, -3), []int{3, 4, 5})
		gtest.Assert(array1.SubSlice(-9, 3), nil)
		gtest.Assert(array1.SubSlice(1, -1), []int{0})
		gtest.Assert(array1.SubSlice(1, -3), nil)
		gtest.Assert(array2.SubSlice(1, 2), []int{1, 2})
	})
}

func TestIntArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(len(array1.Rands(2)), 2)
		gtest.Assert(len(array1.Rands(10)), 7)
		gtest.AssertIN(array1.Rands(1)[0], a1)
		gtest.AssertIN(array1.Rand(), a1)
	})
}

func TestIntArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{100, 200, 300, 400, 500, 600}
		array := garray.NewIntArrayFrom(a1)
		ns1 := array.PopRands(2)
		gtest.AssertIN(ns1, []int{100, 200, 300, 400, 500, 600})
		gtest.AssertIN(len(ns1), 2)

		ns2 := array.PopRands(7)
		gtest.AssertIN(len(ns2), 6)
		gtest.AssertIN(ns2, []int{100, 200, 300, 400, 500, 600})

	})
}

func TestIntArray_Shuffle(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestIntArray_Reverse(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Reverse().Slice(), []int{6, 5, 4, 3, 2, 1, 0})
	})
}

func TestIntArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestNewSortedIntArrayFrom(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 3, 2, 1, 4, 5, 6}
		array1 := garray.NewSortedIntArrayFrom(a1, true)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestNewSortedIntArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 5, 2, 1, 4, 3, 6}
		array1 := garray.NewSortedIntArrayFromCopy(a1, false)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestSortedIntArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 1, 2, 3}
		a2 := []int{4, 5, 6}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := array1.SetArray(a2)

		gtest.Assert(array2.Len(), 3)
		gtest.Assert(array2.Search(3), -1)
		gtest.Assert(array2.Search(5), 1)
		gtest.Assert(array2.Search(6), 2)
	})
}

func TestSortedIntArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{0, 3, 2, 1}

		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := array1.Sort()

		gtest.Assert(array2.Len(), 4)
		gtest.Assert(array2, []int{0, 1, 2, 3})

	})
}

func TestSortedIntArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 0}
		array1 := garray.NewSortedIntArrayFrom(a1)
		gtest.Assert(array1.Get(0), 0)
		gtest.Assert(array1.Get(1), 1)
		gtest.Assert(array1.Get(3), 5)

	})
}

func TestSortedIntArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 0}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.Remove(2)
		gtest.Assert(i1, 3)
		gtest.Assert(array1.Search(5), 2)

		// 再次删除剩下的数组中的第一个
		i2 := array1.Remove(0)
		gtest.Assert(i2, 0)
		gtest.Assert(array1.Search(5), 1)

		a2 := []int{1, 3, 4}
		array2 := garray.NewSortedIntArrayFrom(a2)
		i3 := array2.Remove(1)
		gtest.Assert(array2.Search(1), 0)
		gtest.Assert(i3, 3)
		i3 = array2.Remove(1)
		gtest.Assert(array2.Search(4), -1)
		gtest.Assert(i3, 4)

	})
}

func TestSortedIntArray_PopLeft(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.PopLeft()
		gtest.Assert(i1, 1)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Search(1), -1)

	})
}

func TestSortedIntArray_PopRight(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.PopRight()
		gtest.Assert(i1, 5)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Search(5), -1)
	})
}

func TestSortedIntArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		i1 := array1.PopRand()
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Search(i1), -1)
		gtest.AssertIN(i1, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopRands(2)
		gtest.Assert(array1.Len(), 2)
		gtest.AssertIN(ns1, []int{1, 3, 5, 2})

		a2 := []int{1, 3, 5, 2}
		array2 := garray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopRands(5)
		gtest.Assert(array2.Len(), 0)
		gtest.Assert(len(ns2), 4)
		gtest.AssertIN(ns2, []int{1, 3, 5, 2})

	})
}

func TestSortedIntArray_PopLefts(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopLefts(2)
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(ns1, []int{1, 2})

		a2 := []int{1, 3, 5, 2}
		array2 := garray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopLefts(5)
		gtest.Assert(array2.Len(), 0)
		gtest.AssertIN(ns2, []int{1, 3, 5, 2})

	})
}

func TestSortedIntArray_PopRights(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.PopRights(2)
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(ns1, []int{3, 5})

		a2 := []int{1, 3, 5, 2}
		array2 := garray.NewSortedIntArrayFrom(a2)
		ns2 := array2.PopRights(5)
		gtest.Assert(array2.Len(), 0)
		gtest.AssertIN(ns2, []int{1, 3, 5, 2})
	})
}

func TestSortedIntArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5, 2, 6, 7}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := garray.NewSortedIntArrayFrom(a1, true)
		ns1 := array1.Range(1, 4)
		gtest.Assert(len(ns1), 3)
		gtest.Assert(ns1, []int{2, 3, 5})

		ns2 := array1.Range(5, 4)
		gtest.Assert(len(ns2), 0)

		ns3 := array1.Range(-1, 4)
		gtest.Assert(len(ns3), 4)

		nsl := array1.Range(5, 8)
		gtest.Assert(len(nsl), 1)

		ns4 := array2.Range(2, 5)
		gtest.Assert(len(ns4), 3)
	})
}

func TestSortedIntArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		n1 := array1.Sum()
		gtest.Assert(n1, 9)
	})
}

func TestSortedIntArray_Contains(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		//gtest.Assert(array1.Contains(3),true) //todo 这一行应该返回true
		gtest.Assert(array1.Contains(4), false)
	})
}

func TestSortedIntArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array2.Len(), 3)
		gtest.Assert(array2, array1)
	})
}

func TestSortedIntArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 3, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedIntArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Chunk(2) //按每几个元素切成一个数组
		ns2 := array1.Chunk(-1)
		gtest.Assert(len(ns1), 3)
		gtest.Assert(ns1[0], []int{1, 2})
		gtest.Assert(ns1[2], []int{5})
		gtest.Assert(len(ns2), 0)

	})
}

func TestSortedIntArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.SubSlice(1, 2)
		gtest.Assert(len(ns1), 2)
		gtest.Assert(ns1, []int{2, 3})

		ns2 := array1.SubSlice(7, 2)
		gtest.Assert(len(ns2), 0)

		ns3 := array1.SubSlice(3, 5)
		gtest.Assert(len(ns3), 2)
		gtest.Assert(ns3, []int{4, 5})

		ns4 := array1.SubSlice(3, 1)
		gtest.Assert(len(ns4), 1)
		gtest.Assert(ns4, []int{4})

		array3 := garray.NewSortedIntArrayFrom(a1,true)
		gtest.Assert(array3.SubSlice(2, 2), []int{3, 4})
		gtest.Assert(array3.SubSlice(-1, 2), []int{5})
		gtest.Assert(array3.SubSlice(-9, 2), nil)
		gtest.Assert(array3.SubSlice(4, -2), []int{3,4})
		gtest.Assert(array3.SubSlice(1, -3), nil)

	})
}

func TestSortedIntArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Rand() //按每几个元素切成一个数组
		gtest.AssertIN(ns1, a1)
	})
}

func TestSortedIntArray_Rands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.Rands(2) //按每几个元素切成一个数组
		gtest.AssertIN(ns1, a1)
		gtest.Assert(len(ns1), 2)

		ns2 := array1.Rands(6) //按每几个元素切成一个数组
		gtest.AssertIN(ns2, a1)
		gtest.Assert(len(ns2), 5)
	})
}

func TestSortedIntArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5, 3}
		array1 := garray.NewSortedIntArrayFrom(a1)
		ns1 := array1.CountValues() //按每几个元素切成一个数组
		gtest.Assert(len(ns1), 5)
		gtest.Assert(ns1[2], 1)
		gtest.Assert(ns1[3], 2)
	})
}

func TestSortedIntArray_SetUnique(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 4, 5, 3}
		array1 := garray.NewSortedIntArrayFrom(a1)
		array1.SetUnique(true)
		gtest.Assert(array1.Len(), 5)
		gtest.Assert(array1, []int{1, 2, 3, 4, 5})

	})
}

func TestIntArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		a2 := []int{6, 7}
		array1 := garray.NewIntArrayFrom(a1)
		array1.SetArray(a2)
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(array1, []int{6, 7})

	})
}

func TestIntArray_Replace(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		a2 := []int{6, 7}
		a3 := []int{9, 10, 11, 12, 13}
		array1 := garray.NewIntArrayFrom(a1)
		array1.Replace(a2)
		gtest.Assert(array1, []int{6, 7, 3, 5})

		array1.Replace(a3)
		gtest.Assert(array1, []int{9, 10, 11, 12})
	})
}

func TestIntArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)
	})
}

func TestIntArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array1, array2)
	})
}

func TestArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Get(2), 3)
		gtest.Assert(array1.Len(), 4)
	})
}

func TestIntArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5}
		array1 := garray.NewIntArrayFrom(a1)
		gtest.Assert(array1.Sum(), 11)
	})
}

func TestIntArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5, 3}
		array1 := garray.NewIntArrayFrom(a1)
		m1 := array1.CountValues()
		gtest.Assert(len(m1), 4)
		gtest.Assert(m1[1], 1)
		gtest.Assert(m1[3], 2)
	})
}

func TestNewIntArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5, 3}
		array1 := garray.NewIntArrayFromCopy(a1)
		gtest.Assert(array1.Len(), 5)
		gtest.Assert(array1, a1)
	})
}

func TestIntArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []int{1, 2, 3, 5, 4}
		array1 := garray.NewIntArrayFrom(a1)
		n1 := array1.Remove(1)
		gtest.Assert(n1, 2)
		gtest.Assert(array1.Len(), 4)

		n1 = array1.Remove(0)
		gtest.Assert(n1, 1)
		gtest.Assert(array1.Len(), 3)

		n1 = array1.Remove(2)
		gtest.Assert(n1, 4)
		gtest.Assert(array1.Len(), 2)
	})
}

func TestSortedIntArray_LockFunc(t *testing.T) {
	n1 := []int{1, 2, 4, 3}
	a1 := garray.NewSortedIntArrayFrom(n1)

	ch1 := make(chan int64, 2)
	go a1.LockFunc(func(n1 []int) { //互斥锁
		for i := 1; i <= 4; i++ {
			gtest.Assert(i, n1[i-1])
		}
		n1[3] = 7
		time.Sleep(1 * time.Second) //暂停一秒
	})

	go func() {
		time.Sleep(10 * time.Millisecond) //故意暂停0.01秒,等另一个goroutine执行锁后，再开始执行.
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		a1.Len()
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
	}()

	t1 := <-ch1
	t2 := <-ch1
	// 相差大于0.6秒，说明在读取a1.len时，发生了等待。  防止ci抖动,以豪秒为单位
	gtest.AssertGT(t2-t1, 600)
	gtest.Assert(a1.Contains(7), true)
}

func TestSortedIntArray_RLockFunc(t *testing.T) {
	n1 := []int{1, 2, 4, 3}
	a1 := garray.NewSortedIntArrayFrom(n1)

	ch1 := make(chan int64, 2)
	go a1.RLockFunc(func(n1 []int) { //读锁
		for i := 1; i <= 4; i++ {
			gtest.Assert(i, n1[i-1])
		}
		n1[3] = 7
		time.Sleep(1 * time.Second) //暂停一秒
	})

	go func() {
		time.Sleep(10 * time.Millisecond) //故意暂停0.01秒,等另一个goroutine执行锁后，再开始执行.
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		a1.Len()
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
	}()

	t1 := <-ch1
	t2 := <-ch1
	// 由于另一个goroutine加的读锁，其它可读,所以ch1的操作间隔是很小的.a.len 操作并没有等待,
	// 防止ci抖动,以豪秒为单位
	gtest.AssertLT(t2-t1, 2)
	gtest.Assert(a1.Contains(7), true)
}

func TestSortedIntArray_Merge(t *testing.T) {
	n1 := []int{1, 2, 4, 3}
	n2 := []int{7, 8, 9}
	n3 := []int{3, 6}

	s1 := []string{"a", "b", "c"}
	in1 := []interface{}{1, "a", 2, "b"}

	a1 := garray.NewSortedIntArrayFrom(n1)
	b1 := garray.NewStringArrayFrom(s1)
	b2 := garray.NewIntArrayFrom(n3)
	b3 := garray.NewArrayFrom(in1)
	b4 := garray.NewSortedStringArrayFrom(s1)
	b5 := garray.NewSortedIntArrayFrom(n3)

	gtest.Assert(a1.Merge(n2).Len(), 7)
	gtest.Assert(a1.Merge(n3).Len(), 9)
	gtest.Assert(a1.Merge(b1).Len(), 12)
	gtest.Assert(a1.Merge(b2).Len(), 14)
	gtest.Assert(a1.Merge(b3).Len(), 18)
	gtest.Assert(a1.Merge(b4).Len(), 21)
	gtest.Assert(a1.Merge(b5).Len(), 23)
}

func TestSortedArray_LockFunc(t *testing.T) {
	n1 := []interface{}{1, 2, 4, 3}

	func1 := func(v1, v2 interface{}) int {
		return strings.Compare(gconv.String(v1), gconv.String(v2))
	}
	a1 := garray.NewSortedArrayFrom(n1, func1)

	ch1 := make(chan int64, 2)
	go a1.LockFunc(func(n1 []interface{}) { //互斥锁
		n1[3] = 7
		time.Sleep(1 * time.Second) //暂停一秒
	})

	go func() {
		time.Sleep(10 * time.Millisecond) //故意暂停0.01秒,等另一个goroutine执行锁后，再开始执行.
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		a1.Len()
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
	}()

	t1 := <-ch1
	t2 := <-ch1
	// 相差大于0.6秒，说明在读取a1.len时，发生了等待。  防止ci抖动,以豪秒为单位
	gtest.AssertGT(t2-t1, 600)
	gtest.Assert(a1.Contains(7), true)
}

func TestSortedArray_RLockFunc(t *testing.T) {
	n1 := []interface{}{1, 2, 4, 3}

	func1 := func(v1, v2 interface{}) int {
		return strings.Compare(gconv.String(v1), gconv.String(v2))
	}
	a1 := garray.NewSortedArrayFrom(n1, func1)

	ch1 := make(chan int64, 2)
	go a1.RLockFunc(func(n1 []interface{}) { //互斥锁
		n1[3] = 7
		time.Sleep(1 * time.Second) //暂停一秒
	})

	go func() {
		time.Sleep(10 * time.Millisecond) //故意暂停0.01秒,等另一个goroutine执行锁后，再开始执行.
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		a1.Len()
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
	}()

	t1 := <-ch1
	t2 := <-ch1
	// 由于另一个goroutine加的读锁，其它可读,所以ch1的操作间隔是很小的.a.len 操作并没有等待,
	// 防止ci抖动,以豪秒为单位
	gtest.AssertLT(t2-t1, 20)
	gtest.Assert(a1.Contains(7), true)
}

func TestSortedArray_Merge(t *testing.T) {
	n1 := []interface{}{1, 2, 4, 3}
	n2 := []int{7, 8, 9}
	n3 := []int{3, 6}

	s1 := []string{"a", "b", "c"}
	in1 := []interface{}{1, "a", 2, "b"}

	func1 := func(v1, v2 interface{}) int {
		return strings.Compare(gconv.String(v1), gconv.String(v2))
	}

	a1 := garray.NewSortedArrayFrom(n1, func1)
	b1 := garray.NewStringArrayFrom(s1)
	b2 := garray.NewIntArrayFrom(n3)
	b3 := garray.NewArrayFrom(in1)
	b4 := garray.NewSortedStringArrayFrom(s1)
	b5 := garray.NewSortedIntArrayFrom(n3)

	gtest.Assert(a1.Merge(n2).Len(), 7)
	gtest.Assert(a1.Merge(n3).Len(), 9)
	gtest.Assert(a1.Merge(b1).Len(), 12)
	gtest.Assert(a1.Merge(b2).Len(), 14)
	gtest.Assert(a1.Merge(b3).Len(), 18)
	gtest.Assert(a1.Merge(b4).Len(), 21)
	gtest.Assert(a1.Merge(b5).Len(), 23)
}

func TestIntArray_SortFunc(t *testing.T) {
	n1 := []int{1, 2, 3, 5, 4}
	a1 := garray.NewIntArrayFrom(n1)

	func1 := func(v1, v2 int) bool {
		if v1 > v2 {
			return false
		}
		return true
	}
	func2 := func(v1, v2 int) bool {
		if v1 > v2 {
			return true
		}
		return true
	}
	a2 := a1.SortFunc(func1)
	gtest.Assert(a2, []int{1, 2, 3, 4, 5})
	a3 := a1.SortFunc(func2)
	gtest.Assert(a3, []int{5, 4, 3, 2, 1})

}

func TestIntArray_LockFunc(t *testing.T) {
	n1 := []int{1, 2, 4, 3}
	a1 := garray.NewIntArrayFrom(n1)
	ch1 := make(chan int64, 2)
	go a1.LockFunc(func(n1 []int) { //互斥锁
		n1[3] = 7
		time.Sleep(1 * time.Second) //暂停一秒
	})

	go func() {
		time.Sleep(10 * time.Millisecond) //故意暂停0.01秒,等另一个goroutine执行锁后，再开始执行.
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		a1.Len()
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
	}()

	t1 := <-ch1
	t2 := <-ch1
	// 相差大于0.6秒，说明在读取a1.len时，发生了等待。  防止ci抖动,以豪秒为单位
	gtest.AssertGT(t2-t1, 600)
	gtest.Assert(a1.Contains(7), true)
}

func TestIntArray_RLockFunc(t *testing.T) {
	n1 := []int{1, 2, 4, 3}
	a1 := garray.NewIntArrayFrom(n1)

	ch1 := make(chan int64, 2)
	go a1.RLockFunc(func(n1 []int) { //互斥锁
		n1[3] = 7
		time.Sleep(1 * time.Second) //暂停一秒
	})

	go func() {
		time.Sleep(10 * time.Millisecond) //故意暂停0.01秒,等另一个goroutine执行锁后，再开始执行.
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		a1.Len()
		ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
	}()

	t1 := <-ch1
	t2 := <-ch1
	// 由于另一个goroutine加的读锁，其它可读,所以ch1的操作间隔是很小的.a.len 操作并没有等待,
	// 防止ci抖动,以豪秒为单位
	gtest.AssertLT(t2-t1, 20)
	gtest.Assert(a1.Contains(7), true)
}

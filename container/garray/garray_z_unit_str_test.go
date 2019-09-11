// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray_test

import (
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
)

func Test_StrArray_Basic(t *testing.T) {
	gtest.Case(t, func() {
		expect := []string{"0", "1", "2", "3"}
		array := garray.NewStrArrayFrom(expect)
		array2 := garray.NewStrArrayFrom(expect, true)
		array3 := garray.NewStrArrayFrom([]string{})
		gtest.Assert(array.Slice(), expect)
		array.Set(0, "100")
		gtest.Assert(array.Get(0), 100)
		gtest.Assert(array.Get(1), 1)
		gtest.Assert(array.Search("100"), 0)
		gtest.Assert(array.Contains("100"), true)
		gtest.Assert(array.Remove(0), 100)
		gtest.Assert(array.Contains("100"), false)
		array.Append("4")
		gtest.Assert(array.Len(), 4)
		array.InsertBefore(0, "100")
		array.InsertAfter(0, "200")
		gtest.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "4"})
		array.InsertBefore(5, "300")
		array.InsertAfter(6, "400")
		gtest.Assert(array.Slice(), []string{"100", "200", "1", "2", "3", "300", "4", "400"})
		gtest.Assert(array.Clear().Len(), 0)
		gtest.Assert(array2.Slice(), expect)
		gtest.Assert(array3.Search("100"), -1)
	})
}

func TestStrArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		expect1 := []string{"0", "1", "2", "3"}
		expect2 := []string{"3", "2", "1", "0"}
		array := garray.NewStrArray()
		for i := 3; i >= 0; i-- {
			array.Append(gconv.String(i))
		}
		array.Sort()
		gtest.Assert(array.Slice(), expect1)
		array.Sort(true)
		gtest.Assert(array.Slice(), expect2)
	})
}

func TestStrArray_Unique(t *testing.T) {
	gtest.Case(t, func() {
		expect := []string{"1", "1", "2", "3"}
		array := garray.NewStrArrayFrom(expect)
		gtest.Assert(array.Unique().Slice(), []string{"1", "2", "3"})
	})
}

func TestStrArray_PushAndPop(t *testing.T) {
	gtest.Case(t, func() {
		expect := []string{"0", "1", "2", "3"}
		array := garray.NewStrArrayFrom(expect)
		gtest.Assert(array.Slice(), expect)
		gtest.Assert(array.PopLeft(), "0")
		gtest.Assert(array.PopRight(), "3")
		gtest.AssertIN(array.PopRand(), []string{"1", "2"})
		gtest.AssertIN(array.PopRand(), []string{"1", "2"})
		gtest.Assert(array.Len(), 0)
		array.PushLeft("1").PushRight("2")
		gtest.Assert(array.Slice(), []string{"1", "2"})
	})
}

func TestStrArray_PopLeftsAndPopRights(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		value2 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(value1)
		array2 := garray.NewStrArrayFrom(value2)
		gtest.Assert(array1.PopLefts(2), []interface{}{"0", "1"})
		gtest.Assert(array1.Slice(), []interface{}{"2", "3", "4", "5", "6"})
		gtest.Assert(array1.PopRights(2), []interface{}{"5", "6"})
		gtest.Assert(array1.Slice(), []interface{}{"2", "3", "4"})
		gtest.Assert(array1.PopRights(20), []interface{}{"2", "3", "4"})
		gtest.Assert(array1.Slice(), []interface{}{})
		gtest.Assert(array2.PopLefts(20), []interface{}{"0", "1", "2", "3", "4", "5", "6"})
		gtest.Assert(array2.Slice(), []interface{}{})
	})
}

func TestString_Range(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(value1)
		array2 := garray.NewStrArrayFrom(value1, true)
		gtest.Assert(array1.Range(0, 1), []interface{}{"0"})
		gtest.Assert(array1.Range(1, 2), []interface{}{"1"})
		gtest.Assert(array1.Range(0, 2), []interface{}{"0", "1"})
		gtest.Assert(array1.Range(-1, 10), value1)
		gtest.Assert(array1.Range(10, 1), nil)
		gtest.Assert(array2.Range(0, 1), []interface{}{"0"})
	})
}

func TestStrArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		a11 := []string{"0", "1", "2", "3"}
		a21 := []string{"4", "5", "6", "7"}
		array1 := garray.NewStrArrayFrom(a11)
		array2 := garray.NewStrArrayFrom(a21)
		gtest.Assert(array1.Merge(array2).Slice(), []string{"0", "1", "2", "3", "4", "5", "6", "7"})

		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStrArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStrArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := garray.NewStrArrayFrom(s1)

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i1).Len(), 9)
		gtest.Assert(a1.Merge(i2).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestStrArray_Fill(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0"}
		a2 := []string{"0"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := garray.NewStrArrayFrom(a2)
		gtest.Assert(array1.Fill(1, 2, "100").Slice(), []string{"0", "100", "100"})
		gtest.Assert(array2.Fill(0, 2, "100").Slice(), []string{"100", "100"})
		s1 := array2.Fill(-1, 2, "100")
		gtest.Assert(s1.Len(), 2)
	})
}

func TestStrArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"1", "2", "3", "4", "5"}
		array1 := garray.NewStrArrayFrom(a1)
		chunks := array1.Chunk(2)
		gtest.Assert(len(chunks), 3)
		gtest.Assert(chunks[0], []string{"1", "2"})
		gtest.Assert(chunks[1], []string{"3", "4"})
		gtest.Assert(chunks[2], []string{"5"})
		gtest.Assert(len(array1.Chunk(0)), 0)
	})
}

func TestStrArray_Pad(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(array1.Pad(3, "1").Slice(), []string{"0", "1", "1"})
		gtest.Assert(array1.Pad(-4, "1").Slice(), []string{"1", "0", "1", "1"})
		gtest.Assert(array1.Pad(3, "1").Slice(), []string{"1", "0", "1", "1"})
	})
}

func TestStrArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := garray.NewStrArrayFrom(a1, true)
		gtest.Assert(array1.SubSlice(0, 2), []string{"0", "1"})
		gtest.Assert(array1.SubSlice(2, 2), []string{"2", "3"})
		gtest.Assert(array1.SubSlice(5, 8), []string{"5", "6"})
		gtest.Assert(array1.SubSlice(8, 2), nil)
		gtest.Assert(array1.SubSlice(1, -2), nil)
		gtest.Assert(array1.SubSlice(-5, 2), []string{"2", "3"})
		gtest.Assert(array1.SubSlice(-10, 1), nil)
		gtest.Assert(array2.SubSlice(0, 2), []string{"0", "1"})
	})
}

func TestStrArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(len(array1.Rands(2)), "2")
		gtest.Assert(len(array1.Rands(10)), "7")
		gtest.AssertIN(array1.Rands(1)[0], a1)
		gtest.Assert(len(array1.Rand()), 1)
		gtest.AssertIN(array1.Rand(), a1)
	})
}

func TestStrArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "b", "c", "d", "e", "f", "g"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.AssertIN(array1.PopRands(1), []string{"a", "b", "c", "d", "e", "f", "g"})
		gtest.AssertIN(array1.PopRands(1), []string{"a", "b", "c", "d", "e", "f", "g"})
		gtest.AssertNI(array1.PopRands(1), array1.Slice())
		gtest.AssertNI(array1.PopRands(1), array1.Slice())
		gtest.Assert(len(array1.PopRands(10)), 3)
	})
}

func TestStrArray_Shuffle(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestStrArray_Reverse(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(array1.Reverse().Slice(), []string{"6", "5", "4", "3", "2", "1", "0"})
	})
}

func TestStrArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestNewStrArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := garray.NewStrArrayFromCopy(a1)
		a3 := garray.NewStrArrayFromCopy(a1, true)
		gtest.Assert(a2.Contains("1"), true)
		gtest.Assert(a2.Len(), 7)
		gtest.Assert(a2, a3)
	})
}

func TestStrArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"a", "b", "c", "d"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(array1.Contains("2"), true)
		gtest.Assert(array1.Len(), 7)

		array1 = array1.SetArray(a2)
		gtest.Assert(array1.Contains("2"), false)
		gtest.Assert(array1.Contains("c"), true)
		gtest.Assert(array1.Len(), 4)
	})
}

func TestStrArray_Replace(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"a", "b", "c", "d"}
		a3 := []string{"o", "p", "q", "x", "y", "z", "w", "r", "v"}
		array1 := garray.NewStrArrayFrom(a1)
		gtest.Assert(array1.Contains("2"), true)
		gtest.Assert(array1.Len(), 7)

		array1 = array1.Replace(a2)
		gtest.Assert(array1.Contains("2"), false)
		gtest.Assert(array1.Contains("c"), true)
		gtest.Assert(array1.Contains("5"), true)
		gtest.Assert(array1.Len(), 7)

		array1 = array1.Replace(a3)
		gtest.Assert(array1.Contains("2"), false)
		gtest.Assert(array1.Contains("c"), false)
		gtest.Assert(array1.Contains("5"), false)
		gtest.Assert(array1.Contains("p"), true)
		gtest.Assert(array1.Contains("r"), false)
		gtest.Assert(array1.Len(), 7)

	})
}

func TestStrArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		a2 := []string{"0", "a", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := garray.NewStrArrayFrom(a2)
		gtest.Assert(array1.Sum(), 21)
		gtest.Assert(array2.Sum(), 18)
	})
}

func TestStrArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		str1 := array1.PopRand()
		gtest.Assert(strings.Contains("0,1,2,3,4,5,6", str1), true)
		gtest.Assert(array1.Len(), 6)
	})
}

func TestStrArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "5", "6"}
		array1 := garray.NewStrArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array2, array1)
		gtest.Assert(array2.Len(), 7)
	})
}

func TestStrArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"0", "1", "2", "3", "4", "4", "6"}
		array1 := garray.NewStrArrayFrom(a1)

		m1 := array1.CountValues()
		gtest.Assert(len(m1), 6)
		gtest.Assert(m1["2"], 1)
		gtest.Assert(m1["4"], 2)
	})
}

func TestNewSortedStrArrayFrom(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayFrom(a1, true)
		gtest.Assert(s1, []string{"a", "b", "c", "d"})
		s2 := garray.NewSortedStrArrayFrom(a1, false)
		gtest.Assert(s2, []string{"a", "b", "c", "d"})
	})
}

func TestNewSortedStrArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		s1 := garray.NewSortedStrArrayFromCopy(a1, true)
		gtest.Assert(s1.Len(), 4)
		gtest.Assert(s1, []string{"a", "b", "c", "d"})
	})
}

func TestSortedStrArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		a2 := []string{"f", "g", "h"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.SetArray(a2)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("d"), false)
		gtest.Assert(array1.Contains("b"), false)
		gtest.Assert(array1.Contains("g"), true)
	})
}

func TestSortedStrArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)

		gtest.Assert(array1, []string{"a", "b", "c", "d"})
		array1.Sort()
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains("c"), true)
		gtest.Assert(array1, []string{"a", "b", "c", "d"})
	})
}

func TestSortedStrArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Get(2), "c")
		gtest.Assert(array1.Get(0), "a")
	})
}

func TestSortedStrArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Remove(2), "c")
		gtest.Assert(array1.Get(2), "d")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("c"), false)

		gtest.Assert(array1.Remove(0), "a")
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(array1.Contains("a"), false)

		// 此时array1里的元素只剩下2个
		gtest.Assert(array1.Remove(1), "d")
		gtest.Assert(array1.Len(), 1)
	})
}

func TestSortedStrArray_PopLeft(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopLeft()
		gtest.Assert(s1, "a")
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains("a"), false)
	})
}

func TestSortedStrArray_PopRight(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRight()
		gtest.Assert(s1, "e")
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains("e"), false)
	})
}

func TestSortedStrArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRand()
		gtest.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1.Contains(s1), false)
	})
}

func TestSortedStrArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRands(2)
		gtest.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(len(s1), 2)

		s1 = array1.PopRands(4)
		gtest.Assert(len(s1), 3)
		gtest.AssertIN(s1, []string{"e", "a", "d", "c", "b"})
	})
}

func TestSortedStrArray_PopLefts(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopLefts(2)
		gtest.Assert(s1, []string{"a", "b"})
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(len(s1), 2)

		s1 = array1.PopLefts(4)
		gtest.Assert(len(s1), 3)
		gtest.Assert(s1, []string{"c", "d", "e"})
	})
}

func TestSortedStrArray_PopRights(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.PopRights(2)
		gtest.Assert(s1, []string{"f", "g"})
		gtest.Assert(array1.Len(), 5)
		gtest.Assert(len(s1), 2)
		s1 = array1.PopRights(6)
		gtest.Assert(len(s1), 5)
		gtest.Assert(s1, []string{"a", "b", "c", "d", "e"})
		gtest.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a1, true)
		s1 := array1.Range(2, 4)
		gtest.Assert(len(s1), 2)
		gtest.Assert(s1, []string{"c", "d"})

		s1 = array1.Range(-1, 2)
		gtest.Assert(len(s1), 2)
		gtest.Assert(s1, []string{"a", "b"})

		s1 = array1.Range(4, 8)
		gtest.Assert(len(s1), 3)
		gtest.Assert(s1, []string{"e", "f", "g"})
		gtest.Assert(array1.Range(10, 2), nil)

		s2 := array2.Range(2, 4)
		gtest.Assert(s2, []string{"c", "d"})

	})
}

func TestSortedStrArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		a2 := []string{"1", "2", "3", "4", "a"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a2)
		gtest.Assert(array1.Sum(), 0)
		gtest.Assert(array2.Sum(), 10)
	})
}

func TestSortedStrArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.Clone()
		gtest.Assert(array1, array2)
		array1.Remove(1)
		gtest.Assert(array2.Len(), 7)
	})
}

func TestSortedStrArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)
	})
}

func TestSortedStrArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := garray.NewSortedStrArrayFrom(a1, true)
		s1 := array1.SubSlice(1, 3)
		gtest.Assert(len(s1), 3)
		gtest.Assert(s1, []string{"b", "c", "d"})
		gtest.Assert(array1.Len(), 7)

		s2 := array1.SubSlice(1, 10)
		gtest.Assert(len(s2), 6)

		s3 := array1.SubSlice(10, 2)
		gtest.Assert(len(s3), 0)

		s3 = array1.SubSlice(-5, 2)
		gtest.Assert(s3, []string{"c", "d"})

		s3 = array1.SubSlice(-10, 2)
		gtest.Assert(s3, nil)

		s3 = array1.SubSlice(1, -2)
		gtest.Assert(s3, nil)

		gtest.Assert(array2.SubSlice(1, 3), []string{"b", "c", "d"})
	})
}

func TestSortedStrArray_Len(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "c", "b", "f", "g"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Len(), 7)

	})
}

func TestSortedStrArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.AssertIN(array1.Rand(), []string{"e", "a", "d"})
	})
}

func TestSortedStrArray_Rands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		s1 := array1.Rands(2)

		gtest.AssertIN(s1, []string{"e", "a", "d"})
		gtest.Assert(len(s1), 2)

		s1 = array1.Rands(4)
		gtest.AssertIN(s1, []string{"e", "a", "d"})
		gtest.Assert(len(s1), 3)
	})
}

func TestSortedStrArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		gtest.Assert(array1.Join(","), "a,d,e")
		gtest.Assert(array1.Join("."), "a.d.e")
	})
}

func TestSortedStrArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		m1 := array1.CountValues()
		gtest.Assert(m1["a"], 2)
		gtest.Assert(m1["d"], 1)

	})
}

func TestSortedStrArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.Chunk(2)
		gtest.Assert(len(array2), 3)
		gtest.Assert(len(array2[0]), 2)
		gtest.Assert(array2[1], []string{"c", "d"})
		gtest.Assert(array1.Chunk(0), nil)
	})
}

func TestSortedStrArray_SetUnique(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewSortedStrArrayFrom(a1)
		array2 := array1.SetUnique(true)
		gtest.Assert(array2.Len(), 4)
		gtest.Assert(array2, []string{"a", "c", "d", "e"})
	})
}

func TestStrArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []string{"e", "a", "d", "a", "c"}
		array1 := garray.NewStrArrayFrom(a1)
		s1 := array1.Remove(1)
		gtest.Assert(s1, "a")
		gtest.Assert(array1.Len(), 4)
		s1 = array1.Remove(3)
		gtest.Assert(s1, "c")
		gtest.Assert(array1.Len(), 3)
	})
}

func TestStrArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []string) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		gtest.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候，并没有阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []string) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		gtest.AssertGT(t2-t1, 20) //go1加的读写互斥锁，所go2读的时候被阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewSortedStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []string) { //读锁
			time.Sleep(2 * time.Second) //暂停1秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		gtest.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候，并没有阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedStrArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStrArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStrArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := garray.NewSortedStrArrayFrom(s1)

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i1).Len(), 9)
		gtest.Assert(a1.Merge(i2).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestStrArray_SortFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "d", "c", "b"}
		a1 := garray.NewStrArrayFrom(s1)
		func1 := func(v1, v2 string) bool {
			return v1 < v2
		}
		a11 := a1.SortFunc(func1)
		gtest.Assert(a11, []string{"a", "b", "c", "d"})
	})
}

func TestStrArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []string{"a", "b", "c", "d"}
		a1 := garray.NewStrArrayFrom(s1, true)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []string) { //读写锁
			time.Sleep(2 * time.Second) //暂停2秒
			n1[2] = "g"
			ch2 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		})

		//go2
		go func() {
			time.Sleep(100 * time.Millisecond) //故意暂停0.01秒,等go1执行锁后，再开始执行.
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
			a1.Len()
			ch1 <- gconv.Int64(time.Now().UnixNano() / 1000 / 1000)
		}()

		t1 := <-ch1
		t2 := <-ch1
		<-ch2 //等待go1完成

		// 防止ci抖动,以豪秒为单位
		gtest.AssertGT(t2-t1, 20) //go1加的读写互斥锁，所go2读的时候被阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

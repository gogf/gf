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

func Test_Array_Basic(t *testing.T) {
	gtest.Case(t, func() {
		expect := []interface{}{0, 1, 2, 3}
		array := garray.NewArrayFrom(expect)
		array2 := garray.NewArrayFrom(expect)
		array3 := garray.NewArrayFrom([]interface{}{})
		gtest.Assert(array.Slice(), expect)
		array.Set(0, 100)
		gtest.Assert(array.Get(0), 100)
		gtest.Assert(array.Get(1), 1)
		gtest.Assert(array.Search(100), 0)
		gtest.Assert(array3.Search(100), -1)
		gtest.Assert(array.Contains(100), true)
		gtest.Assert(array.Remove(0), 100)

		gtest.Assert(array2.Remove(3), 3)
		gtest.Assert(array2.Remove(1), 1)

		gtest.Assert(array.Contains(100), false)
		array.Append(4)
		gtest.Assert(array.Len(), 4)
		array.InsertBefore(0, 100)
		array.InsertAfter(0, 200)
		gtest.Assert(array.Slice(), []interface{}{100, 200, 2, 2, 3, 4})
		array.InsertBefore(5, 300)
		array.InsertAfter(6, 400)
		gtest.Assert(array.Slice(), []interface{}{100, 200, 2, 2, 3, 300, 4, 400})
		gtest.Assert(array.Clear().Len(), 0)
	})
}

func TestArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		expect1 := []interface{}{0, 1, 2, 3}
		expect2 := []interface{}{3, 2, 1, 0}
		array := garray.NewArray()
		for i := 3; i >= 0; i-- {
			array.Append(i)
		}
		array.SortFunc(func(v1, v2 interface{}) bool {
			return v1.(int) < v2.(int)
		})
		gtest.Assert(array.Slice(), expect1)
		array.SortFunc(func(v1, v2 interface{}) bool {
			return v1.(int) > v2.(int)
		})
		gtest.Assert(array.Slice(), expect2)
	})
}

func TestArray_Unique(t *testing.T) {
	gtest.Case(t, func() {
		expect := []interface{}{1, 1, 2, 3}
		array := garray.NewArrayFrom(expect)
		gtest.Assert(array.Unique().Slice(), []interface{}{1, 2, 3})
	})
}

func TestArray_PushAndPop(t *testing.T) {
	gtest.Case(t, func() {
		expect := []interface{}{0, 1, 2, 3}
		array := garray.NewArrayFrom(expect)
		gtest.Assert(array.Slice(), expect)
		gtest.Assert(array.PopLeft(), 0)
		gtest.Assert(array.PopRight(), 3)
		gtest.AssertIN(array.PopRand(), []interface{}{1, 2})
		gtest.AssertIN(array.PopRand(), []interface{}{1, 2})
		gtest.Assert(array.Len(), 0)
		array.PushLeft(1).PushRight(2)
		gtest.Assert(array.Slice(), []interface{}{1, 2})
	})
}

func TestArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{100, 200, 300, 400, 500, 600}
		array := garray.NewFromCopy(a1)
		gtest.AssertIN(array.PopRands(2), a1)
	})
}

func TestArray_PopLeftsAndPopRights(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		value2 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(value1)
		array2 := garray.NewArrayFrom(value2)
		gtest.Assert(array1.PopLefts(2), []interface{}{0, 1})
		gtest.Assert(array1.Slice(), []interface{}{2, 3, 4, 5, 6})
		gtest.Assert(array1.PopRights(2), []interface{}{5, 6})
		gtest.Assert(array1.Slice(), []interface{}{2, 3, 4})
		gtest.Assert(array1.PopRights(20), []interface{}{2, 3, 4})
		gtest.Assert(array1.Slice(), []interface{}{})
		gtest.Assert(array2.PopLefts(20), []interface{}{0, 1, 2, 3, 4, 5, 6})
		gtest.Assert(array2.Slice(), []interface{}{})
	})
}

func TestArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		value1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(value1)
		array2 := garray.NewArrayFrom(value1, true)
		gtest.Assert(array1.Range(0, 1), []interface{}{0})
		gtest.Assert(array1.Range(1, 2), []interface{}{1})
		gtest.Assert(array1.Range(0, 2), []interface{}{0, 1})
		gtest.Assert(array1.Range(-1, 10), value1)
		gtest.Assert(array1.Range(10, 2), nil)
		gtest.Assert(array2.Range(1, 3), []interface{}{1, 2})
	})
}

func TestArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		i1 := []interface{}{0, 1, 2, 3}
		i2 := []interface{}{4, 5, 6, 7}
		array1 := garray.NewArrayFrom(i1)
		array2 := garray.NewArrayFrom(i2)
		gtest.Assert(array1.Merge(array2).Slice(), []interface{}{0, 1, 2, 3, 4, 5, 6, 7})

		//s1 := []string{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i3 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i4 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStringArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStringArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})
		a1 := garray.NewArrayFrom(i1)

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i3).Len(), 9)
		gtest.Assert(a1.Merge(i4).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)
	})
}

func TestArray_Fill(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0}
		a2 := []interface{}{0}
		array1 := garray.NewArrayFrom(a1)
		array2 := garray.NewArrayFrom(a2, true)
		gtest.Assert(array1.Fill(1, 2, 100).Slice(), []interface{}{0, 100, 100})
		gtest.Assert(array2.Fill(0, 2, 100).Slice(), []interface{}{100, 100})
		gtest.Assert(array2.Fill(-1, 2, 100).Slice(), []interface{}{100, 100})
	})
}

func TestArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{1, 2, 3, 4, 5}
		array1 := garray.NewArrayFrom(a1)
		chunks := array1.Chunk(2)
		gtest.Assert(len(chunks), 3)
		gtest.Assert(chunks[0], []interface{}{1, 2})
		gtest.Assert(chunks[1], []interface{}{3, 4})
		gtest.Assert(chunks[2], []interface{}{5})
		gtest.Assert(array1.Chunk(0), nil)
	})
}

func TestArray_Pad(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0}
		array1 := garray.NewArrayFrom(a1)
		gtest.Assert(array1.Pad(3, 1).Slice(), []interface{}{0, 1, 1})
		gtest.Assert(array1.Pad(-4, 1).Slice(), []interface{}{1, 0, 1, 1})
		gtest.Assert(array1.Pad(3, 1).Slice(), []interface{}{1, 0, 1, 1})
	})
}

func TestArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(a1)
		array2 := garray.NewArrayFrom(a1, true)
		gtest.Assert(array1.SubSlice(0, 2), []interface{}{0, 1})
		gtest.Assert(array1.SubSlice(2, 2), []interface{}{2, 3})
		gtest.Assert(array1.SubSlice(5, 8), []interface{}{5, 6})
		gtest.Assert(array1.SubSlice(9, 1), nil)
		gtest.Assert(array1.SubSlice(-2, 2), []interface{}{5, 6})
		gtest.Assert(array1.SubSlice(-9, 2), nil)
		gtest.Assert(array1.SubSlice(1, -2), nil)
		gtest.Assert(array2.SubSlice(0, 2), []interface{}{0, 1})
	})
}

func TestArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(a1)
		gtest.Assert(len(array1.Rands(2)), 2)
		gtest.Assert(len(array1.Rands(10)), 7)
		gtest.AssertIN(array1.Rands(1)[0], a1)
	})

	gtest.Case(t, func() {
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewArrayFrom(s1)
		i1 := a1.Rand()
		gtest.Assert(a1.Contains(i1), true)
		gtest.Assert(a1.Len(), 4)
	})
}

func TestArray_Shuffle(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(a1)
		gtest.Assert(array1.Shuffle().Len(), 7)
	})
}

func TestArray_Reverse(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(a1)
		gtest.Assert(array1.Reverse().Slice(), []interface{}{6, 5, 4, 3, 2, 1, 0})
	})
}

func TestArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		array1 := garray.NewArrayFrom(a1)
		gtest.Assert(array1.Join("."), "0.1.2.3.4.5.6")
	})
}

func TestArray_Replace(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		a2 := []interface{}{"a", "b", "c"}
		a3 := []interface{}{"m", "n", "p", "z", "x", "y", "d", "u"}
		array1 := garray.NewArrayFrom(a1)
		array2 := array1.Replace(a2)
		gtest.Assert(array2.Len(), 7)
		gtest.Assert(array2.Contains("b"), true)
		gtest.Assert(array2.Contains(4), true)
		gtest.Assert(array2.Contains("v"), false)
		array3 := array1.Replace(a3)
		gtest.Assert(array3.Len(), 7)
		gtest.Assert(array3.Contains(4), false)
		gtest.Assert(array3.Contains("p"), true)
		gtest.Assert(array3.Contains("u"), false)
	})
}

func TestArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3, 4, 5, 6}
		a2 := []interface{}{"a", "b", "c"}

		array1 := garray.NewArrayFrom(a1)
		array1 = array1.SetArray(a2)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("b"), true)
		gtest.Assert(array1.Contains("5"), false)
	})
}

func TestArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3}
		a2 := []interface{}{"a", "b", "c"}
		a3 := []interface{}{"a", "1", "2"}

		array1 := garray.NewArrayFrom(a1)
		array2 := garray.NewArrayFrom(a2)
		array3 := garray.NewArrayFrom(a3)

		gtest.Assert(array1.Sum(), 6)
		gtest.Assert(array2.Sum(), 0)
		gtest.Assert(array3.Sum(), 3)

	})
}

func TestArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{0, 1, 2, 3}
		array1 := garray.NewArrayFrom(a1)
		array2 := array1.Clone()

		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array2.Sum(), 6)
		gtest.AssertEQ(array1, array2)

	})
}

func TestArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "b", "c", "d", "e", "d"}
		array1 := garray.NewArrayFrom(a1)
		array2 := array1.CountValues()
		gtest.Assert(len(array2), 5)
		gtest.Assert(array2["b"], 1)
		gtest.Assert(array2["d"], 2)
	})
}

func TestSortedArray_NewSortedArrayFrom(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		a2 := []interface{}{"h", "j", "i", "k"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 interface{}) int {
			return -1
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a2, func2)

		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "f"})

		gtest.Assert(array2.Len(), 4)
		gtest.Assert(array2, []interface{}{"k", "i", "j", "h"})
	})
}

func TestNewSortedArrayFromCopy(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		func2 := func(v1, v2 interface{}) int {
			return -1
		}
		array1 := garray.NewSortedArrayFromCopy(a1, func1)
		array2 := garray.NewSortedArrayFromCopy(a1, func2)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "f"})
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array2, []interface{}{"c", "f", "a"})
	})
}

func TestSortedArray_SetArray(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		a2 := []interface{}{"e", "h", "g", "k"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}

		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.SetArray(a2)
		gtest.Assert(array1.Len(), 4)
		gtest.Assert(array1, []interface{}{"e", "g", "h", "k"})
	})

}

func TestSortedArray_Sort(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.Sort()
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "f"})
	})

}

func TestSortedArray_Get(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "f", "c"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		gtest.Assert(array1.Get(2), "f")
		gtest.Assert(array1.Get(1), "c")
	})

}

func TestSortedArray_Remove(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Remove(1)
		gtest.Assert(gconv.String(i1), "b")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1.Contains("b"), false)

		i2 := array1.Remove(0)
		gtest.Assert(gconv.String(i2), "a")
		gtest.Assert(array1.Len(), 2)
		gtest.Assert(array1.Contains("a"), false)

		i3 := array1.Remove(1)
		gtest.Assert(gconv.String(i3), "d")
		gtest.Assert(array1.Len(), 1)
		gtest.Assert(array1.Contains("d"), false)
	})

}

func TestSortedArray_PopLeft(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopLeft()
		gtest.Assert(gconv.String(i1), "a")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"b", "c", "d"})
	})

}

func TestSortedArray_PopRight(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRight()
		gtest.Assert(gconv.String(i1), "d")
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "b", "c"})
	})

}

func TestSortedArray_PopRand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRand()
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 3)

	})
}

func TestSortedArray_PopRands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRands(2)
		gtest.Assert(len(i1), 2)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 2)

		i2 := array1.PopRands(3)
		gtest.Assert(len(i1), 2)
		gtest.AssertIN(i2, []interface{}{"a", "d", "c", "b"})
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_PopLefts(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopLefts(2)
		gtest.Assert(len(i1), 2)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b", "e", "f"})
		gtest.Assert(array1.Len(), 4)

		i2 := array1.PopLefts(5)
		gtest.Assert(len(i2), 4)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c", "b", "e", "f"})
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_PopRights(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.PopRights(2)
		gtest.Assert(len(i1), 2)
		gtest.Assert(i1, []interface{}{"e", "f"})
		gtest.Assert(array1.Len(), 4)

		i2 := array1.PopRights(10)
		gtest.Assert(len(i2), 4)

	})
}

func TestSortedArray_Range(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a1, func1, true)
		i1 := array1.Range(2, 5)
		gtest.Assert(i1, []interface{}{"c", "d", "e"})
		gtest.Assert(array1.Len(), 6)

		i2 := array1.Range(7, 5)
		gtest.Assert(len(i2), 0)
		i2 = array1.Range(-1, 2)
		gtest.Assert(i2, []interface{}{"a", "b"})

		i2 = array1.Range(4, 10)
		gtest.Assert(len(i2), 2)
		gtest.Assert(i2, []interface{}{"e", "f"})

		gtest.Assert(array2.Range(1, 3), []interface{}{"b", "c"})

	})
}

func TestSortedArray_Sum(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}
		a2 := []interface{}{"1", "2", "3", "b", "e", "f"}
		a3 := []interface{}{"4", "5", "6"}
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a2, func1)
		array3 := garray.NewSortedArrayFrom(a3, func1)
		gtest.Assert(array1.Sum(), 0)
		gtest.Assert(array2.Sum(), 6)
		gtest.Assert(array3.Sum(), 15)

	})
}

func TestSortedArray_Clone(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := array1.Clone()
		gtest.Assert(array1, array2)
		array1.Remove(1)
		gtest.AssertNE(array1, array2)

	})
}

func TestSortedArray_Clear(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e", "f"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		gtest.Assert(array1.Len(), 6)
		array1.Clear()
		gtest.Assert(array1.Len(), 0)

	})
}

func TestSortedArray_Chunk(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Chunk(2)
		gtest.Assert(len(i1), 3)
		gtest.Assert(i1[0], []interface{}{"a", "b"})
		gtest.Assert(i1[2], []interface{}{"e"})

		i1 = array1.Chunk(0)
		gtest.Assert(len(i1), 0)
	})
}

func TestSortedArray_SubSlice(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "b", "e"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array2 := garray.NewSortedArrayFrom(a1, func1, true)
		i1 := array1.SubSlice(2, 3)
		gtest.Assert(len(i1), 3)
		gtest.Assert(i1, []interface{}{"c", "d", "e"})

		i1 = array1.SubSlice(2, 6)
		gtest.Assert(len(i1), 3)
		gtest.Assert(i1, []interface{}{"c", "d", "e"})

		i1 = array1.SubSlice(7, 2)
		gtest.Assert(len(i1), 0)

		s1 := array1.SubSlice(1, -2)
		gtest.Assert(s1, nil)

		s1 = array1.SubSlice(-9, 2)
		gtest.Assert(s1, nil)
		gtest.Assert(array2.SubSlice(1, 3), []interface{}{"b", "c", "d"})

	})
}

func TestSortedArray_Rand(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Rand()
		gtest.AssertIN(i1, []interface{}{"a", "d", "c"})
		gtest.Assert(array1.Len(), 3)
	})
}

func TestSortedArray_Rands(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		i1 := array1.Rands(2)
		gtest.AssertIN(i1, []interface{}{"a", "d", "c"})
		gtest.Assert(len(i1), 2)
		gtest.Assert(array1.Len(), 3)

		i1 = array1.Rands(4)
		gtest.Assert(len(i1), 3)
	})
}

func TestSortedArray_Join(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		gtest.Assert(array1.Join(","), "a,c,d")
		gtest.Assert(array1.Join("."), "a.c.d")

	})
}

func TestSortedArray_CountValues(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		m1 := array1.CountValues()
		gtest.Assert(len(m1), 3)
		gtest.Assert(m1["c"], 2)
		gtest.Assert(m1["a"], 1)

	})
}

func TestSortedArray_SetUnique(t *testing.T) {
	gtest.Case(t, func() {
		a1 := []interface{}{"a", "d", "c", "c"}

		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		array1 := garray.NewSortedArrayFrom(a1, func1)
		array1.SetUnique(true)
		gtest.Assert(array1.Len(), 3)
		gtest.Assert(array1, []interface{}{"a", "c", "d"})
	})
}

func TestSortedArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, func1)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []interface{}) { //读写锁
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

func TestSortedArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			return strings.Compare(gconv.String(v1), gconv.String(v2))
		}
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewSortedArrayFrom(s1, func1)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.RLockFunc(func(n1 []interface{}) { //读写锁
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
		gtest.AssertLT(t2-t1, 20) //go1加的读锁，所go2读的时候不会被阻塞。
		gtest.Assert(a1.Contains("g"), true)
	})
}

func TestSortedArray_Merge(t *testing.T) {
	gtest.Case(t, func() {
		func1 := func(v1, v2 interface{}) int {
			if gconv.Int(v1) < gconv.Int(v2) {
				return 0
			}
			return 1
		}

		s1 := []interface{}{"a", "b", "c", "d"}
		s2 := []string{"e", "f"}
		i1 := garray.NewIntArrayFrom([]int{1, 2, 3})
		i2 := garray.NewArrayFrom([]interface{}{3})
		s3 := garray.NewStringArrayFrom([]string{"g", "h"})
		s4 := garray.NewSortedArrayFrom([]interface{}{4, 5}, func1)
		s5 := garray.NewSortedStringArrayFrom(s2)
		s6 := garray.NewSortedIntArrayFrom([]int{1, 2, 3})

		a1 := garray.NewSortedArrayFrom(s1, func1)

		gtest.Assert(a1.Merge(s2).Len(), 6)
		gtest.Assert(a1.Merge(i1).Len(), 9)
		gtest.Assert(a1.Merge(i2).Len(), 10)
		gtest.Assert(a1.Merge(s3).Len(), 12)
		gtest.Assert(a1.Merge(s4).Len(), 14)
		gtest.Assert(a1.Merge(s5).Len(), 16)
		gtest.Assert(a1.Merge(s6).Len(), 19)

	})
}

func TestArray_LockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewArrayFrom(s1)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 3)
		//go1
		go a1.LockFunc(func(n1 []interface{}) { //读写锁
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

func TestArray_RLockFunc(t *testing.T) {
	gtest.Case(t, func() {
		s1 := []interface{}{"a", "b", "c", "d"}
		a1 := garray.NewArrayFrom(s1)

		ch1 := make(chan int64, 3)
		ch2 := make(chan int64, 1)
		//go1
		go a1.RLockFunc(func(n1 []interface{}) { //读锁
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

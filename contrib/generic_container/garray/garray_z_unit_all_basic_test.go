// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go

package garray

import (
	"testing"

	"github.com/gogf/gf/contrib/generic_container/v2/comparator"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Array_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var array StdArray[int]
		expect := []int{2, 3, 1}
		array.Append(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array StdArray[int]
		expect := []int{2, 3, 1}
		array.Append(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array StdArray[string]
		expect := []string{"b", "a"}
		array.Append("b", "a")
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array SortedArray[int]
		array.SetComparator(comparator.ComparatorInt)
		expect := []int{1, 2, 3}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array SortedArray[int]
		array.SetComparator(comparator.ComparatorInt)
		expect := []int{1, 2, 3}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
	gtest.C(t, func(t *gtest.T) {
		var array SortedArray[string]
		array.SetComparator(comparator.ComparatorString)
		expect := []string{"a", "b", "c"}
		array.Add("c", "a", "b")
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedIntArray_Var(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var array SortedArray[int]
		array.SetComparator(comparator.ComparatorInt)
		expect := []int{1, 2, 3}
		array.Add(2, 3, 1)
		t.Assert(array.Slice(), expect)
	})
}

func Test_IntArray_Unique(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []int{1, 2, 3, 4, 5, 6}
		array := NewArray[int]()
		array.Append(1, 1, 2, 3, 3, 4, 4, 5, 5, 6, 6)
		array.Unique()
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedIntArray1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		array := NewSortedArray[int](comparator.ComparatorInt)
		for i := 10; i > -1; i-- {
			array.Add(i)
		}
		t.Assert(array.Slice(), expect)
		t.Assert(array.Add().Slice(), expect)
	})
}

func Test_SortedIntArray2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		array := NewSortedArray[int](comparator.ComparatorInt)
		for i := 0; i <= 10; i++ {
			array.Add(i)
		}
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedStrArray1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array1 := NewSortedArray[string](comparator.ComparatorString)
		array2 := NewSortedArray[string](comparator.ComparatorString, true)
		for i := 10; i > -1; i-- {
			array1.Add(gconv.String(i))
			array2.Add(gconv.String(i))
		}
		t.Assert(array1.Slice(), expect)
		t.Assert(array2.Slice(), expect)
	})

}

func Test_SortedStrArray2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array := NewSortedArray[string](comparator.ComparatorString)
		for i := 0; i <= 10; i++ {
			array.Add(gconv.String(i))
		}
		t.Assert(array.Slice(), expect)
		array.Add()
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedArray1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array := NewSortedArray(comparator.ComparatorString)
		for i := 10; i > -1; i-- {
			array.Add(gconv.String(i))
		}
		t.Assert(array.Slice(), expect)
	})
}

func Test_SortedArray2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		expect := []string{"0", "1", "10", "2", "3", "4", "5", "6", "7", "8", "9"}
		array := NewSortedArray[string](comparator.ComparatorString)
		array2 := NewSortedArray[string](comparator.ComparatorString, true)
		for i := 0; i <= 10; i++ {
			array.Add(gconv.String(i))
			array2.Add(gconv.String(i))
		}
		t.Assert(array.Slice(), expect)
		t.Assert(array.Add().Slice(), expect)
		t.Assert(array2.Slice(), expect)
	})
}

func TestNewFromCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a1 := []interface{}{"100", "200", "300", "400", "500", "600"}
		array1 := NewFromCopy(a1)
		t.AssertIN(array1.PopRands(2), a1)
		t.Assert(len(array1.PopRands(1)), 1)
		t.Assert(len(array1.PopRands(9)), 3)
	})
}

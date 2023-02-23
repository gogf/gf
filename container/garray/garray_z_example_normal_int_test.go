// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"fmt"

	"github.com/gogf/gf/v2/internal/empty"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

func ExampleIntArray_Walk() {
	var array garray.IntArray
	tables := g.SliceInt{10, 20}
	prefix := 99
	array.Append(tables...)
	// Add prefix for given table names.
	array.Walk(func(value int) int {
		return prefix + value
	})
	fmt.Println(array.Slice())

	// Output:
	// [109 119]
}

func ExampleNewIntArray() {
	s := garray.NewIntArray()
	s.Append(10)
	s.Append(20)
	s.Append(15)
	s.Append(30)
	fmt.Println(s.Slice())

	// Output:
	// [10 20 15 30]
}

func ExampleNewIntArraySize() {
	s := garray.NewIntArraySize(3, 5)
	s.Set(0, 10)
	s.Set(1, 20)
	s.Set(2, 15)
	s.Set(3, 30)
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [10 20 15] 3 5
}

func ExampleNewIntArrayRange() {
	s := garray.NewIntArrayRange(1, 5, 1)
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [1 2 3 4 5] 5 8
}

func ExampleNewIntArrayFrom() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [10 20 15 30] 4 4
}

func ExampleNewIntArrayFromCopy() {
	s := garray.NewIntArrayFromCopy(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [10 20 15 30] 4 4
}

func ExampleIntArray_At() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30})
	sAt := s.At(2)
	fmt.Println(sAt)

	// Output:
	// 15
}

func ExampleIntArray_Get() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30})
	sGet, sBool := s.Get(3)
	fmt.Println(sGet, sBool)
	sGet, sBool = s.Get(99)
	fmt.Println(sGet, sBool)

	// Output:
	// 30 true
	// 0 false
}

func ExampleIntArray_Set() {
	s := garray.NewIntArraySize(3, 5)
	s.Set(0, 10)
	s.Set(1, 20)
	s.Set(2, 15)
	s.Set(3, 30)
	fmt.Println(s.Slice())

	// Output:
	// [10 20 15]
}

func ExampleIntArray_SetArray() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s.Slice())

	// Output:
	// [10 20 15 30]
}

func ExampleIntArray_Replace() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s.Slice())
	s.Replace(g.SliceInt{12, 13})
	fmt.Println(s.Slice())

	// Output:
	// [10 20 15 30]
	// [12 13 15 30]
}

func ExampleIntArray_Sum() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	a := s.Sum()
	fmt.Println(a)

	// Output:
	// 75
}

func ExampleIntArray_Sort() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	a := s.Sort()
	fmt.Println(a)

	// Output:
	// [10,15,20,30]
}

func ExampleIntArray_SortFunc() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.SortFunc(func(v1, v2 int) bool {
		// fmt.Println(v1,v2)
		return v1 > v2
	})
	fmt.Println(s)
	s.SortFunc(func(v1, v2 int) bool {
		return v1 < v2
	})
	fmt.Println(s)

	// Output:
	// [10,20,15,30]
	// [30,20,15,10]
	// [10,15,20,30]
}

func ExampleIntArray_InsertBefore() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	s.InsertBefore(1, 99)
	fmt.Println(s.Slice())

	// Output:
	// [10 99 20 15 30]
}

func ExampleIntArray_InsertAfter() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	s.InsertAfter(1, 99)
	fmt.Println(s.Slice())

	// Output:
	// [10 20 99 15 30]
}

func ExampleIntArray_Remove() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.Remove(1)
	fmt.Println(s.Slice())

	// Output:
	// [10,20,15,30]
	// [10 15 30]
}

func ExampleIntArray_RemoveValue() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.RemoveValue(20)
	fmt.Println(s.Slice())

	// Output:
	// [10,20,15,30]
	// [10 15 30]
}

func ExampleIntArray_PushLeft() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.PushLeft(96, 97, 98, 99)
	fmt.Println(s.Slice())

	// Output:
	// [10,20,15,30]
	// [96 97 98 99 10 20 15 30]
}

func ExampleIntArray_PushRight() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.PushRight(96, 97, 98, 99)
	fmt.Println(s.Slice())

	// Output:
	// [10,20,15,30]
	// [10 20 15 30 96 97 98 99]
}

func ExampleIntArray_PopLeft() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.PopLeft()
	fmt.Println(s.Slice())

	// Output:
	// [10,20,15,30]
	// [20 15 30]
}

func ExampleIntArray_PopRight() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30})
	fmt.Println(s)
	s.PopRight()
	fmt.Println(s.Slice())

	// Output:
	// [10,20,15,30]
	// [10 20 15]
}

func ExampleIntArray_PopRand() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60, 70})
	fmt.Println(s)
	r, _ := s.PopRand()
	fmt.Println(s)
	fmt.Println(r)

	// May Output:
	// [10,20,15,30,40,50,60,70]
	// [10,20,15,30,40,60,70]
	// 50
}

func ExampleIntArray_PopRands() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.PopRands(2)
	fmt.Println(s)
	fmt.Println(r)

	// May Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40]
	// [50 60]
}

func ExampleIntArray_PopLefts() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.PopLefts(2)
	fmt.Println(s)
	fmt.Println(r)

	// Output:
	// [10,20,15,30,40,50,60]
	// [15,30,40,50,60]
	// [10 20]
}

func ExampleIntArray_PopRights() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.PopRights(2)
	fmt.Println(s)
	fmt.Println(r)

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40]
	// [50 60]
}

func ExampleIntArray_Range() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.Range(2, 5)
	fmt.Println(r)

	// Output:
	// [10,20,15,30,40,50,60]
	// [15 30 40]
}

func ExampleIntArray_SubSlice() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.SubSlice(3, 4)
	fmt.Println(r)

	// Output:
	// [10,20,15,30,40,50,60]
	// [30 40 50 60]
}

func ExampleIntArray_Append() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	s.Append(96, 97, 98)
	fmt.Println(s)

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40,50,60,96,97,98]
}

func ExampleIntArray_Len() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Len())

	// Output:
	// [10,20,15,30,40,50,60]
	// 7
}

func ExampleIntArray_Slice() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s.Slice())

	// Output:
	// [10 20 15 30 40 50 60]
}

func ExampleIntArray_Interfaces() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	r := s.Interfaces()
	fmt.Println(r)

	// Output:
	// [10 20 15 30 40 50 60]
}

func ExampleIntArray_Clone() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.Clone()
	fmt.Println(r)

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40,50,60]
}

func ExampleIntArray_Clear() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Clear())
	fmt.Println(s)

	// Output:
	// [10,20,15,30,40,50,60]
	// []
	// []
}

func ExampleIntArray_Contains() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s.Contains(20))
	fmt.Println(s.Contains(21))

	// Output:
	// true
	// false
}

func ExampleIntArray_Search() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s.Search(20))
	fmt.Println(s.Search(21))

	// Output:
	// 1
	// -1
}

func ExampleIntArray_Unique() {
	s := garray.NewIntArray()
	s.SetArray(g.SliceInt{10, 20, 15, 15, 20, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Unique())

	// Output:
	// [10,20,15,15,20,50,60]
	// [10,20,15,50,60]
}

func ExampleIntArray_LockFunc() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	s.LockFunc(func(array []int) {
		for i := 0; i < len(array)-1; i++ {
			fmt.Println(array[i])
		}
	})

	// Output:
	// 10
	// 20
	// 15
	// 30
	// 40
	// 50
}

func ExampleIntArray_RLockFunc() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	s.RLockFunc(func(array []int) {
		for i := 0; i < len(array); i++ {
			fmt.Println(array[i])
		}
	})

	// Output:
	// 10
	// 20
	// 15
	// 30
	// 40
	// 50
	// 60
}

func ExampleIntArray_Merge() {
	s1 := garray.NewIntArray()
	s2 := garray.NewIntArray()
	s1.SetArray(g.SliceInt{10, 20, 15})
	s2.SetArray(g.SliceInt{40, 50, 60})
	fmt.Println(s1)
	fmt.Println(s2)
	s1.Merge(s2)
	fmt.Println(s1)

	// Output:
	// [10,20,15]
	// [40,50,60]
	// [10,20,15,40,50,60]
}

func ExampleIntArray_Fill() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	s.Fill(2, 3, 99)
	fmt.Println(s)

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,99,99,99,50,60]
}

func ExampleIntArray_Chunk() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	r := s.Chunk(3)
	fmt.Println(r)

	// Output:
	// [10,20,15,30,40,50,60]
	// [[10 20 15] [30 40 50] [60]]
}

func ExampleIntArray_Pad() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	s.Pad(8, 99)
	fmt.Println(s)
	s.Pad(-10, 89)
	fmt.Println(s)

	// Output:
	// [10,20,15,30,40,50,60,99]
	// [89,89,10,20,15,30,40,50,60,99]
}

func ExampleIntArray_Rand() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Rand())

	// May Output:
	// [10,20,15,30,40,50,60]
	// 10 true
}

func ExampleIntArray_Rands() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Rands(3))

	// May Output:
	// [10,20,15,30,40,50,60]
	// [20 50 20]
}

func ExampleIntArray_Shuffle() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Shuffle())

	// May Output:
	// [10,20,15,30,40,50,60]
	// [10,40,15,50,20,60,30]
}

func ExampleIntArray_Reverse() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Reverse())

	// Output:
	// [10,20,15,30,40,50,60]
	// [60,50,40,30,15,20,10]
}

func ExampleIntArray_Join() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.Join(","))

	// Output:
	// [10,20,15,30,40,50,60]
	// 10,20,15,30,40,50,60
}

func ExampleIntArray_CountValues() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 15, 40, 40, 40})
	fmt.Println(s.CountValues())

	// Output:
	// map[10:1 15:2 20:1 40:3]
}

func ExampleIntArray_Iterator() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	s.Iterator(func(k int, v int) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 10
	// 1 20
	// 2 15
	// 3 30
	// 4 40
	// 5 50
	// 6 60
}

func ExampleIntArray_IteratorAsc() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	s.IteratorAsc(func(k int, v int) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 10
	// 1 20
	// 2 15
	// 3 30
	// 4 40
	// 5 50
	// 6 60
}

func ExampleIntArray_IteratorDesc() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	s.IteratorDesc(func(k int, v int) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 6 60
	// 5 50
	// 4 40
	// 3 30
	// 2 15
	// 1 20
	// 0 10
}

func ExampleIntArray_String() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s)
	fmt.Println(s.String())

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40,50,60]
}

func ExampleIntArray_MarshalJSON() {
	type Student struct {
		Id     int
		Name   string
		Scores garray.IntArray
	}
	var array garray.IntArray
	array.SetArray(g.SliceInt{98, 97, 96})
	s := Student{
		Id:     1,
		Name:   "john",
		Scores: array,
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))

	// Output:
	// {"Id":1,"Name":"john","Scores":[98,97,96]}
}

func ExampleIntArray_UnmarshalJSON() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[98,96,97]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *garray.IntArray
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// Output:
	// {1 john [98,96,97]}
}

func ExampleIntArray_UnmarshalValue() {
	type Student struct {
		Name   string
		Scores *garray.IntArray
	}

	var s *Student
	gconv.Struct(g.Map{
		"name":   "john",
		"scores": g.SliceInt{96, 98, 97},
	}, &s)
	fmt.Println(s)

	// Output:
	// &{john [96,98,97]}
}

func ExampleIntArray_Filter() {
	array1 := garray.NewIntArrayFrom(g.SliceInt{10, 40, 50, 0, 0, 0, 60})
	array2 := garray.NewIntArrayFrom(g.SliceInt{10, 4, 51, 5, 45, 50, 56})
	fmt.Println(array1.Filter(func(value int, index int) bool {
		return empty.IsEmpty(value)
	}))
	fmt.Println(array2.Filter(func(value int, index int) bool {
		return value%2 == 0
	}))
	fmt.Println(array2.Filter(func(value int, index int) bool {
		return value%2 == 1
	}))

	// Output:
	// [10,40,50,60]
	// [51,5,45]
	// [10,4,50,56]
}

func ExampleIntArray_FilterEmpty() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 40, 50, 0, 0, 0, 60})
	fmt.Println(s)
	fmt.Println(s.FilterEmpty())

	// Output:
	// [10,40,50,0,0,0,60]
	// [10,40,50,60]
}

func ExampleIntArray_IsEmpty() {
	s := garray.NewIntArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
	fmt.Println(s.IsEmpty())
	s1 := garray.NewIntArray()
	fmt.Println(s1.IsEmpty())

	// Output:
	// false
	// true
}

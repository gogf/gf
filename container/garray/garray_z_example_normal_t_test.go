// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func ExampleTArray_Walk() {
	{
		var intArray garray.TArray[int]
		intTables := g.SliceInt{10, 20}
		intPrefix := 99
		intArray.Append(intTables...)
		// Add prefix for given table names.
		intArray.Walk(func(value int) int {
			return intPrefix + value
		})
		fmt.Println(intArray.Slice())
	}

	{
		var strArray garray.TArray[string]
		strTables := g.SliceStr{"user", "user_detail"}
		strPrefix := "gf_"
		strArray.Append(strTables...)
		// Add prefix for given table names.
		strArray.Walk(func(value string) string {
			return strPrefix + value
		})
		fmt.Println(strArray.Slice())
	}

	// Output:
	// [109 119]
	// [gf_user gf_user_detail]
}

func ExampleNewTArray() {
	{
		intArr := garray.NewTArray[int]()
		intArr.Append(10)
		intArr.Append(20)
		intArr.Append(15)
		intArr.Append(30)
		fmt.Println(intArr.Slice())
	}

	{
		strArr := garray.NewTArray[string]()
		strArr.Append("We")
		strArr.Append("are")
		strArr.Append("GF")
		strArr.Append("fans")
		fmt.Println(strArr.Slice())
	}

	// Output:
	// [10 20 15 30]
	// [We are GF fans]
}

func ExampleNewTArraySize() {
	{
		intArr := garray.NewTArraySize[int](3, 5)
		intArr.Set(0, 10)
		intArr.Set(1, 20)
		intArr.Set(2, 15)
		intArr.Set(3, 30)
		fmt.Println(intArr.Slice(), intArr.Len(), cap(intArr.Slice()))
	}

	{
		strArr := garray.NewTArraySize[string](3, 5)
		strArr.Set(0, "We")
		strArr.Set(1, "are")
		strArr.Set(2, "GF")
		strArr.Set(3, "fans")
		fmt.Println(strArr.Slice(), strArr.Len(), cap(strArr.Slice()))
	}

	// Output:
	// [10 20 15] 3 5
	// [We are GF] 3 5
}

func ExampleNewTArrayFrom() {
	{
		intArr := garray.NewTArrayFrom[int](g.SliceInt{10, 20, 15, 30})
		fmt.Println(intArr.Slice(), intArr.Len(), cap(intArr.Slice()))
	}

	{
		strArr := garray.NewTArrayFrom[string](g.SliceStr{"We", "are", "GF", "fans", "!"})
		fmt.Println(strArr.Slice(), strArr.Len(), cap(strArr.Slice()))
	}

	// Output:
	// [10 20 15 30] 4 4
	// [We are GF fans !] 5 5
}

func ExampleNewTArrayFromCopy() {
	{
		intArr := garray.NewTArrayFromCopy(g.SliceInt{10, 20, 15, 30})
		fmt.Println(intArr.Slice(), intArr.Len(), cap(intArr.Slice()))
	}

	{
		strArr := garray.NewTArrayFromCopy(g.SliceStr{"a", "b", "c", "d", "e"})
		fmt.Println(strArr.Slice(), strArr.Len(), cap(strArr.Slice()))
	}

	// Output:
	// [10 20 15 30] 4 4
	// [a b c d e] 5 5
}

func ExampleTArray_At() {
	{
		intArr := garray.NewTArrayFrom[int](g.SliceInt{10, 20, 15, 30})
		isAt := intArr.At(2)
		fmt.Println(isAt)
	}

	{
		strArr := garray.NewTArrayFrom[string](g.SliceStr{"We", "are", "GF", "fans", "!"})
		ssAt := strArr.At(2)
		fmt.Println(ssAt)
	}

	// Output:
	// 15
	// GF
}

func ExampleTArray_Get() {
	{
		s := garray.NewTArrayFrom[int](g.SliceInt{10, 20, 15, 30})
		sGet, sBool := s.Get(3)
		fmt.Println(sGet, sBool)
		sGet, sBool = s.Get(99)
		fmt.Println(sGet, sBool)
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"We", "are", "GF", "fans", "!"})
		sGet, sBool := s.Get(3)
		fmt.Println(sGet, sBool)
	}

	// Output:
	// 30 true
	// 0 false
	// fans true
}

func ExampleTArray_Set() {
	{
		s := garray.NewTArraySize[int](3, 5)
		s.Set(0, 10)
		s.Set(1, 20)
		s.Set(2, 15)
		s.Set(3, 30)
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArraySize[string](3, 5)
		s.Set(0, "We")
		s.Set(1, "are")
		s.Set(2, "GF")
		s.Set(3, "fans")
		fmt.Println(s.Slice())
	}

	// Output:
	// [10 20 15]
	// [We are GF]
}

func ExampleTArray_SetArray() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"We", "are", "GF", "fans", "!"})
		fmt.Println(s.Slice())
	}

	// Output:
	// [10 20 15 30]
	// [We are GF fans !]
}

func ExampleTArray_Replace() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s.Slice())
		s.Replace(g.SliceInt{12, 13})
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"We", "are", "GF", "fans", "!"})
		fmt.Println(s.Slice())
		s.Replace(g.SliceStr{"Happy", "coding"})
		fmt.Println(s.Slice())
	}

	// Output:
	// [10 20 15 30]
	// [12 13 15 30]
	// [We are GF fans !]
	// [Happy coding GF fans !]
}

func ExampleTArray_Sum() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		a := s.Sum()
		fmt.Println(a)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"3", "5", "10"})
		a := s.Sum()
		fmt.Println(a)
	}
	// Output:
	// 75
	// 18
}

func ExampleTArray_SortFunc() {
	{
		s := garray.NewTArrayFrom[int](g.SliceInt{10, 20, 15, 30})
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
	}

	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"b", "c", "a"})
		fmt.Println(s)
		s.SortFunc(func(v1, v2 string) bool {
			return gstr.Compare(v1, v2) > 0
		})
		fmt.Println(s)
		s.SortFunc(func(v1, v2 string) bool {
			return gstr.Compare(v1, v2) < 0
		})
		fmt.Println(s)
	}

	// Output:
	// [10,20,15,30]
	// [30,20,15,10]
	// [10,15,20,30]
	// ["b","c","a"]
	// ["c","b","a"]
	// ["a","b","c"]
}

func ExampleTArray_InsertBefore() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		s.InsertBefore(1, 99)
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.InsertBefore(1, "here")
		fmt.Println(s.Slice())
	}

	// Output:
	// [10 99 20 15 30]
	// [a here b c d]
}

func ExampleTArray_InsertAfter() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		s.InsertAfter(1, 99)
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.InsertAfter(1, "here")
		fmt.Println(s.Slice())
	}
	// Output:
	// [10 20 99 15 30]
	// [a b here c d]
}

func ExampleTArray_Remove() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s)
		s.Remove(1)
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.Remove(1)
		fmt.Println(s.Slice())
	}
	// Output:
	// [10,20,15,30]
	// [10 15 30]
	// [a c d]
}

func ExampleTArray_RemoveValue() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s)
		s.RemoveValue(20)
		fmt.Println(s.Slice())
	}

	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.RemoveValue("b")
		fmt.Println(s.Slice())
	}

	// Output:
	// [10,20,15,30]
	// [10 15 30]
	// [a c d]
}

func ExampleTArray_PushLeft() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s)
		s.PushLeft(96, 97, 98, 99)
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.PushLeft("We", "are", "GF", "fans")
		fmt.Println(s.Slice())
	}

	// Output:
	// [10,20,15,30]
	// [96 97 98 99 10 20 15 30]
	// [We are GF fans a b c d]
}

func ExampleTArray_PushRight() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s)
		s.PushRight(96, 97, 98, 99)
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.PushRight("We", "are", "GF", "fans")
		fmt.Println(s.Slice())
	}
	// Output:
	// [10,20,15,30]
	// [10 20 15 30 96 97 98 99]
	// [a b c d We are GF fans]
}

func ExampleTArray_PopLeft() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s)
		s.PopLeft()
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.PopLeft()
		fmt.Println(s.Slice())
	}
	// Output:
	// [10,20,15,30]
	// [20 15 30]
	// [b c d]
}

func ExampleTArray_PopRight() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30})
		fmt.Println(s)
		s.PopRight()
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d"})
		s.PopRight()
		fmt.Println(s.Slice())
	}
	// Output:
	// [10,20,15,30]
	// [10 20 15]
	// [a b c]
}

func ExampleTArray_PopRand() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60, 70})
		fmt.Println(s)
		r, _ := s.PopRand()
		fmt.Println(s)
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r, _ := s.PopRand()
		fmt.Println(r)
	}

	// May Output:
	// [10,20,15,30,40,50,60,70]
	// [10,20,15,30,40,60,70]
	// 50
	// e
}

func ExampleTArray_PopRands() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.PopRands(2)
		fmt.Println(s)
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.PopRands(2)
		fmt.Println(r)
	}
	// May Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40]
	// [50 60]
	// [e c]
}

func ExampleTArray_PopLefts() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.PopLefts(2)
		fmt.Println(s)
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.PopLefts(2)
		fmt.Println(r)
		fmt.Println(s)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// [15,30,40,50,60]
	// [10 20]
	// [a b]
	// ["c","d","e","f","g","h"]
}

func ExampleTArray_PopRights() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.PopRights(2)
		fmt.Println(s)
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.PopRights(2)
		fmt.Println(r)
		fmt.Println(s)
	}

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40]
	// [50 60]
	// [g h]
	// ["a","b","c","d","e","f"]
}

func ExampleTArray_Range() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.Range(2, 5)
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.Range(2, 5)
		fmt.Println(r)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// [15 30 40]
	// [c d e]
}

func ExampleTArray_SubSlice() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.SubSlice(3, 4)
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.SubSlice(3, 4)
		fmt.Println(r)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// [30 40 50 60]
	// [d e f g]
}

func ExampleTArray_Append() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		s.Append(96, 97, 98)
		fmt.Println(s)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"We", "are", "GF", "fans"})
		s.Append("a", "b", "c")
		fmt.Println(s)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40,50,60,96,97,98]
	// ["We","are","GF","fans","a","b","c"]
}

func ExampleTArray_Len() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Len())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Len())
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// 7
	// 8
}

func ExampleTArray_Slice() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s.Slice())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Slice())
	}
	// Output:
	// [10 20 15 30 40 50 60]
	// [a b c d e f g h]
}

func ExampleTArray_Interfaces() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		r := s.Interfaces()
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.Interfaces()
		fmt.Println(r)
	}
	// Output:
	// [10 20 15 30 40 50 60]
	// [a b c d e f g h]
}

func ExampleTArray_Clone() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.Clone()
		fmt.Println(r)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.Clone()
		fmt.Println(r)
		fmt.Println(s)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40,50,60]
	// ["a","b","c","d","e","f","g","h"]
	// ["a","b","c","d","e","f","g","h"]
}

func ExampleTArray_Clear() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Clear())
		fmt.Println(s)
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s)
		fmt.Println(s.Clear())
		fmt.Println(s)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// []
	// []
	// ["a","b","c","d","e","f","g","h"]
	// []
	// []
}

func ExampleTArray_Contains() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s.Contains(20))
		fmt.Println(s.Contains(21))
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Contains("e"))
		fmt.Println(s.Contains("z"))
	}
	// Output:
	// true
	// false
	// true
	// false
}

func ExampleTArray_Search() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s.Search(20))
		fmt.Println(s.Search(21))
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Search("e"))
		fmt.Println(s.Search("z"))
	}
	// Output:
	// 1
	// -1
	// 4
	// -1
}

func ExampleTArray_Unique() {
	{
		s := garray.NewTArray[int]()
		s.SetArray(g.SliceInt{10, 20, 15, 15, 20, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Unique())
	}
	{
		s := garray.NewTArray[string]()
		s.SetArray(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
		fmt.Println(s.Unique())
	}
	// Output:
	// [10,20,15,15,20,50,60]
	// [10,20,15,50,60]
	// ["a","b","c","d"]
}

func ExampleTArray_LockFunc() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		s.LockFunc(func(array []int) {
			for i := 0; i < len(array)-1; i++ {
				fmt.Println(array[i])
			}
		})
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		s.LockFunc(func(array []string) {
			array[len(array)-1] = "GF fans"
		})
		fmt.Println(s)
	}
	// Output:
	// 10
	// 20
	// 15
	// 30
	// 40
	// 50
	// ["a","b","GF fans"]
}

func ExampleTArray_RLockFunc() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		s.RLockFunc(func(array []int) {
			for i := 0; i < len(array); i++ {
				fmt.Println(array[i])
			}
		})
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e"})
		s.RLockFunc(func(array []string) {
			for i := 0; i < len(array); i++ {
				fmt.Println(array[i])
			}
		})
	}
	// Output:
	// 10
	// 20
	// 15
	// 30
	// 40
	// 50
	// 60
	// a
	// b
	// c
	// d
	// e
}

func ExampleTArray_Merge() {
	{
		s1 := garray.NewTArray[int]()
		s2 := garray.NewTArray[int]()
		s1.SetArray(g.SliceInt{10, 20, 15})
		s2.SetArray(g.SliceInt{40, 50, 60})
		fmt.Println(s1)
		fmt.Println(s2)
		s1.Merge(s2)
		fmt.Println(s1)
	}
	{
		s1 := garray.NewTArray[string]()
		s2 := garray.NewTArray[string]()
		s1.SetArray(g.SliceStr{"a", "b", "c"})
		s2.SetArray(g.SliceStr{"d", "e", "f"})
		s1.Merge(s2)
		fmt.Println(s1)
	}
	// Output:
	// [10,20,15]
	// [40,50,60]
	// [10,20,15,40,50,60]
	// ["a","b","c","d","e","f"]
}

func ExampleTArray_Fill() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		s.Fill(2, 3, 99)
		fmt.Println(s)
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		s.Fill(2, 3, "here")
		fmt.Println(s)
	}
	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,99,99,99,50,60]
	// ["a","b","here","here","here","f","g","h"]
}

func ExampleTArray_Chunk() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		r := s.Chunk(3)
		fmt.Println(r)
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		r := s.Chunk(3)
		fmt.Println(r)
	}

	// Output:
	// [10,20,15,30,40,50,60]
	// [[10 20 15] [30 40 50] [60]]
	// [[a b c] [d e f] [g h]]
}

func ExampleTArray_Pad() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		s.Pad(8, 99)
		fmt.Println(s)
		s.Pad(-10, 89)
		fmt.Println(s)
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		s.Pad(7, "here")
		fmt.Println(s)
		s.Pad(-10, "there")
		fmt.Println(s)
	}

	// Output:
	// [10,20,15,30,40,50,60,99]
	// [89,89,10,20,15,30,40,50,60,99]
	// ["a","b","c","here","here","here","here"]
	// ["there","there","there","a","b","c","here","here","here","here"]
}

func ExampleTArray_Rand() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Rand())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Rand())
	}

	// May Output:
	// [10,20,15,30,40,50,60]
	// 10 true
	// c true
}

func ExampleTArray_Rands() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Rands(3))
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Rands(3))
	}

	// May Output:
	// [10,20,15,30,40,50,60]
	// [20 50 20]
	// [e h e]
}

func ExampleTArray_Shuffle() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Shuffle())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Shuffle())
	}

	// May Output:
	// [10,20,15,30,40,50,60]
	// [10,40,15,50,20,60,30]
	// ["a","c","e","d","b","g","f","h"]
}

func ExampleTArray_Reverse() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Reverse())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
		fmt.Println(s.Reverse())
	}

	// Output:
	// [10,20,15,30,40,50,60]
	// [60,50,40,30,15,20,10]
	// ["h","g","f","e","d","c","b","a"]
}

func ExampleTArray_Join() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.Join(","))
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		fmt.Println(s.Join(","))
	}

	// Output:
	// [10,20,15,30,40,50,60]
	// 10,20,15,30,40,50,60
	// a,b,c
}

func ExampleTArray_CountValues() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 15, 40, 40, 40})
		fmt.Println(s.CountValues())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
		fmt.Println(s.CountValues())
	}

	// Output:
	// map[10:1 15:2 20:1 40:3]
	// map[a:1 b:1 c:3 d:2]
}

func ExampleTArray_Iterator() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		s.Iterator(func(k int, v int) bool {
			fmt.Println(k, v)
			return true
		})
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		s.Iterator(func(k int, v string) bool {
			fmt.Println(k, v)
			return true
		})
	}

	// Output:
	// 0 10
	// 1 20
	// 2 15
	// 3 30
	// 4 40
	// 5 50
	// 6 60
	// 0 a
	// 1 b
	// 2 c
}

func ExampleTArray_IteratorAsc() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		s.IteratorAsc(func(k int, v int) bool {
			fmt.Println(k, v)
			return true
		})
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		s.IteratorAsc(func(k int, v string) bool {
			fmt.Println(k, v)
			return true
		})
	}

	// Output:
	// 0 10
	// 1 20
	// 2 15
	// 3 30
	// 4 40
	// 5 50
	// 6 60
	// 0 a
	// 1 b
	// 2 c
}

func ExampleTArray_IteratorDesc() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		s.IteratorDesc(func(k int, v int) bool {
			fmt.Println(k, v)
			return true
		})
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		s.IteratorDesc(func(k int, v string) bool {
			fmt.Println(k, v)
			return true
		})
	}

	// Output:
	// 6 60
	// 5 50
	// 4 40
	// 3 30
	// 2 15
	// 1 20
	// 0 10
	// 2 c
	// 1 b
	// 0 a
}

func ExampleTArray_String() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s)
		fmt.Println(s.String())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
		fmt.Println(s.String())
	}

	// Output:
	// [10,20,15,30,40,50,60]
	// [10,20,15,30,40,50,60]
	// ["a","b","c"]
}

func ExampleTArray_MarshalJSON() {
	{
		type Student struct {
			Id     int
			Name   string
			Scores garray.TArray[int]
		}
		var array garray.TArray[int]
		array.SetArray(g.SliceInt{98, 97, 96})
		s := Student{
			Id:     1,
			Name:   "john",
			Scores: array,
		}
		b, _ := json.Marshal(s)
		fmt.Println(string(b))
	}
	{
		type Student struct {
			Id      int
			Name    string
			Lessons []string
		}
		s := Student{
			Id:      1,
			Name:    "john",
			Lessons: []string{"Math", "English", "Music"},
		}
		b, _ := json.Marshal(s)
		fmt.Println(string(b))
	}

	// Output:
	// {"Id":1,"Name":"john","Scores":[98,97,96]}
	// {"Id":1,"Name":"john","Lessons":["Math","English","Music"]}
}

func ExampleTArray_UnmarshalJSON() {
	{
		b := []byte(`{"Id":1,"Name":"john","Scores":[98,96,97]}`)
		type Student struct {
			Id     int
			Name   string
			Scores *garray.TArray[int]
		}
		s := Student{}
		json.Unmarshal(b, &s)
		fmt.Println(s)
	}
	{
		b := []byte(`{"Id":1,"Name":"john","Lessons":["Math","English","Sport"]}`)
		type Student struct {
			Id      int
			Name    string
			Lessons *garray.TArray[string]
		}
		s := Student{}
		json.Unmarshal(b, &s)
		fmt.Println(s)
	}

	// Output:
	// {1 john [98,96,97]}
	// {1 john ["Math","English","Sport"]}
}

func ExampleTArray_UnmarshalValue() {
	{
		type Student struct {
			Name   string
			Scores *garray.TArray[int]
		}

		var s *Student
		gconv.Struct(g.Map{
			"name":   "john",
			"scores": g.SliceInt{96, 98, 97},
		}, &s)
		fmt.Println(s)
	}
	{
		type Student struct {
			Name    string
			Lessons *garray.TArray[string]
		}
		var s *Student
		gconv.Struct(g.Map{
			"name":    "john",
			"lessons": []byte(`["Math","English","Sport"]`),
		}, &s)
		fmt.Println(s)

		var s1 *Student
		gconv.Struct(g.Map{
			"name":    "john",
			"lessons": g.SliceStr{"Math", "English", "Sport"},
		}, &s1)
		fmt.Println(s1)
	}

	// Output:
	// &{john [96,98,97]}
	// &{john ["Math","English","Sport"]}
	// &{john ["Math","English","Sport"]}
}

func ExampleTArray_Filter() {
	{
		array1 := garray.NewTArrayFrom(g.SliceInt{10, 40, 50, 0, 0, 0, 60})
		array2 := garray.NewTArrayFrom(g.SliceInt{10, 4, 51, 5, 45, 50, 56})
		fmt.Println(array1.Filter(func(index int, value int) bool {
			return empty.IsEmpty(value)
		}))
		fmt.Println(array2.Filter(func(index int, value int) bool {
			return value%2 == 0
		}))
		fmt.Println(array2.Filter(func(index int, value int) bool {
			return value%2 == 1
		}))
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"Math", "English", "Sport"})
		s1 := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "", "c", "", "", "d"})
		fmt.Println(s1.Filter(func(index int, value string) bool {
			return empty.IsEmpty(value)
		}))

		fmt.Println(s.Filter(func(index int, value string) bool {
			return strings.Contains(value, "h")
		}))
	}

	// Output:
	// [10,40,50,60]
	// [51,5,45]
	// []
	// ["a","b","c","d"]
	// ["Sport"]
}

func ExampleTArray_FilterEmpty() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 40, 50, 0, 0, 0, 60})
		fmt.Println(s)
		fmt.Println(s.FilterEmpty())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "", "c", "", "", "d"})
		fmt.Println(s.FilterEmpty())
	}

	// Output:
	// [10,40,50,0,0,0,60]
	// [10,40,50,60]
	// ["a","b","c","d"]
}

func ExampleTArray_IsEmpty() {
	{
		s := garray.NewTArrayFrom(g.SliceInt{10, 20, 15, 30, 40, 50, 60})
		fmt.Println(s.IsEmpty())
		s1 := garray.NewTArray[int]()
		fmt.Println(s1.IsEmpty())
	}
	{
		s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "", "c", "", "", "d"})
		fmt.Println(s.IsEmpty())
		s1 := garray.NewTArray[string]()
		fmt.Println(s1.IsEmpty())
	}

	// Output:
	// false
	// true
	// false
	// true
}

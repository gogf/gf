// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

func ExampleSortedTArray_Walk() {
	var array garray.SortedTArray[string]
	array.SetComparator(gutil.ComparatorT)
	tables := g.SliceStr{"user", "user_detail"}
	prefix := "gf_"
	array.Append(tables...)
	// Add prefix for given table names.
	array.Walk(func(value string) string {
		return prefix + value
	})
	fmt.Println(array.Slice())

	// Output:
	// [gf_user gf_user_detail]
}

func ExampleNewSortedTArray() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.Append("b")
	s.Append("d")
	s.Append("c")
	s.Append("a")
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
}

func ExampleNewSortedTArraySize() {
	s := garray.NewSortedTArraySize[string](3, gutil.ComparatorT)
	s.SetArray([]string{"b", "d", "a", "c"})
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [a b c d] 4 4
}

func ExampleNewSortedTArrayFromCopy() {
	s := garray.NewSortedTArrayFromCopy(g.SliceStr{"b", "d", "c", "a"}, gutil.ComparatorT)
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
}

func ExampleSortedTArray_At() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "d", "c", "a"}, gutil.ComparatorT)
	sAt := s.At(2)
	fmt.Println(s)
	fmt.Println(sAt)

	// Output:
	// ["a","b","c","d"]
	// c

}

func ExampleSortedTArray_Get() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "d", "c", "a", "e"}, gutil.ComparatorT)
	sGet, sBool := s.Get(3)
	fmt.Println(s)
	fmt.Println(sGet, sBool)

	// Output:
	// ["a","b","c","d","e"]
	// d true
}

func ExampleSortedTArray_SetArray() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray([]string{"b", "d", "a", "c"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
}

func ExampleSortedTArray_SetUnique() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray([]string{"b", "d", "a", "c", "c", "a"})
	fmt.Println(s.SetUnique(true))

	// Output:
	// ["a","b","c","d"]
}

func ExampleSortedTArray_Sum() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray([]string{"5", "3", "2"})
	fmt.Println(s)
	a := s.Sum()
	fmt.Println(a)

	// Output:
	// [2,3,5]
	// 10
}

func ExampleSortedTArray_Sort() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"b", "d", "a", "c"})
	fmt.Println(s)
	a := s.Sort()
	fmt.Println(a)

	// Output:
	// ["a","b","c","d"]
	// ["a","b","c","d"]
}

func ExampleSortedTArray_Remove() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s.Slice())
	s.Remove(1)
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
	// [a c d]
}

func ExampleSortedTArray_RemoveValue() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s.Slice())
	s.RemoveValue("b")
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
	// [a c d]
}

func ExampleSortedTArray_PopLeft() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	r, _ := s.PopLeft()
	fmt.Println(r)
	fmt.Println(s.Slice())

	// Output:
	// a
	// [b c d]
}

func ExampleSortedTArray_PopRight() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s.Slice())
	r, _ := s.PopRight()
	fmt.Println(r)
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
	// d
	// [a b c]
}

func ExampleSortedTArray_PopRights() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.PopRights(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [g h]
	// ["a","b","c","d","e","f"]
}

func ExampleSortedTArray_Rand() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r, _ := s.PopRand()
	fmt.Println(r)
	fmt.Println(s)

	// May Output:
	// b
	// ["a","c","d","e","f","g","h"]
}

func ExampleSortedTArray_PopRands() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.PopRands(2)
	fmt.Println(r)
	fmt.Println(s)

	// May Output:
	// [d a]
	// ["b","c","e","f","g","h"]
}

func ExampleSortedTArray_PopLefts() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.PopLefts(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [a b]
	// ["c","d","e","f","g","h"]
}

func ExampleSortedTArray_Range() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Range(2, 5)
	fmt.Println(r)

	// Output:
	// [c d e]
}

func ExampleSortedTArray_SubSlice() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.SubSlice(3, 4)
	fmt.Println(s.Slice())
	fmt.Println(r)

	// Output:
	// [a b c d e f g h]
	// [d e f g]
}

func ExampleSortedTArray_Add() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.Add("b", "d", "c", "a")
	fmt.Println(s)

	// Output:
	// ["a","b","c","d"]
}

func ExampleSortedTArray_Append() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s)
	s.Append("f", "e", "g")
	fmt.Println(s)

	// Output:
	// ["a","b","c","d"]
	// ["a","b","c","d","e","f","g"]
}

func ExampleSortedTArray_Len() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.Len())

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// 8
}

func ExampleSortedTArray_Slice() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d e f g h]
}

func ExampleSortedTArray_Interfaces() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Interfaces()
	fmt.Println(r)

	// Output:
	// [a b c d e f g h]
}

func ExampleSortedTArray_Clone() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Clone()
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// ["a","b","c","d","e","f","g","h"]
}

func ExampleSortedTArray_Clear() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.Clear())
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// []
	// []
}

func ExampleSortedTArray_Contains() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s.Contains("e"))
	fmt.Println(s.Contains("E"))
	fmt.Println(s.Contains("z"))

	// Output:
	// true
	// false
	// false
}

func ExampleSortedTArray_Search() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.Search("e"))
	fmt.Println(s.Search("E"))
	fmt.Println(s.Search("z"))

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// 4
	// -1
	// -1
}

func ExampleSortedTArray_Unique() {
	s := garray.NewSortedTArray[string](gutil.ComparatorT)
	s.SetArray(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s)
	fmt.Println(s.Unique())

	// Output:
	// ["a","b","c","c","c","d","d"]
	// ["a","b","c","d"]
}

func ExampleSortedTArray_LockFunc() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	s.LockFunc(func(array []string) {
		array[len(array)-1] = "GF fans"
	})
	fmt.Println(s)

	// Output:
	// ["GF fans","a","b"]
}

func ExampleSortedTArray_RLockFunc() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	s.RLockFunc(func(array []string) {
		array[len(array)-1] = "GF fans"
		fmt.Println(array[len(array)-1])
	})
	fmt.Println(s)

	// Output:
	// GF fans
	// ["a","b","GF fans"]
}

func ExampleSortedTArray_Merge() {
	s1 := garray.NewSortedTArray[string](gutil.ComparatorT)
	s2 := garray.NewSortedTArray[string](gutil.ComparatorT)
	s1.SetArray(g.SliceStr{"b", "c", "a"})
	s2.SetArray(g.SliceStr{"e", "d", "f"})
	fmt.Println(s1)
	fmt.Println(s2)
	s1.Merge(s2)
	fmt.Println(s1)

	// Output:
	// ["a","b","c"]
	// ["d","e","f"]
	// ["a","b","c","d","e","f"]
}

func ExampleSortedTArray_Chunk() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"}, gutil.ComparatorT)
	r := s.Chunk(3)
	fmt.Println(r)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleSortedTArray_Rands() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"}, gutil.ComparatorT)
	fmt.Println(s)
	fmt.Println(s.Rands(3))

	// May Output:
	// ["a","b","c","d","e","f","g","h"]
	// [h g c]
}

func ExampleSortedTArray_Join() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"}, gutil.ComparatorT)
	fmt.Println(s.Join(","))

	// Output:
	// a,b,c,d,e,f,g,h
}

func ExampleSortedTArray_CountValues() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"}, gutil.ComparatorT)
	fmt.Println(s.CountValues())

	// Output:
	// map[a:1 b:1 c:3 d:2]
}

func ExampleSortedTArray_Iterator() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	s.Iterator(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleSortedTArray_IteratorAsc() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	s.IteratorAsc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleSortedTArray_IteratorDesc() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	s.IteratorDesc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 2 c
	// 1 b
	// 0 a
}

func ExampleSortedTArray_String() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	fmt.Println(s.String())

	// Output:
	// ["a","b","c"]
}

func ExampleSortedTArray_MarshalJSON() {
	type Student struct {
		ID     int
		Name   string
		Levels garray.SortedTArray[string]
	}
	r := garray.NewSortedTArrayFrom(g.SliceStr{"b", "c", "a"}, gutil.ComparatorT)
	s := Student{
		ID:     1,
		Name:   "john",
		Levels: *r,
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))

	// Output:
	// {"ID":1,"Name":"john","Levels":["a","b","c"]}
}

func ExampleSortedTArray_UnmarshalJSON() {
	b := []byte(`{"Id":1,"Name":"john","Lessons":["Math","English","Sport"]}`)
	type Student struct {
		Id      int
		Name    string
		Lessons *garray.StrArray
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// Output:
	// {1 john ["Math","English","Sport"]}
}

func ExampleSortedTArray_UnmarshalValue() {
	type Student struct {
		Name    string
		Lessons *garray.StrArray
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

	// Output:
	// &{john ["Math","English","Sport"]}
	// &{john ["Math","English","Sport"]}
}

func ExampleSortedTArray_Filter() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "a", "", "c", "", "", "d"}, gutil.ComparatorT)
	fmt.Println(s)
	fmt.Println(s.Filter(func(index int, value string) bool {
		return empty.IsEmpty(value)
	}))

	// Output:
	// ["","","","a","b","c","d"]
	// ["a","b","c","d"]
}

func ExampleSortedTArray_FilterEmpty() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "a", "", "c", "", "", "d"}, gutil.ComparatorT)
	fmt.Println(s)
	fmt.Println(s.FilterEmpty())

	// Output:
	// ["","","","a","b","c","d"]
	// ["a","b","c","d"]
}

func ExampleSortedTArray_IsEmpty() {
	s := garray.NewSortedTArrayFrom(g.SliceStr{"b", "a", "", "c", "", "", "d"}, gutil.ComparatorT)
	fmt.Println(s.IsEmpty())
	s1 := garray.NewSortedTArray[string](gutil.ComparatorT)
	fmt.Println(s1.IsEmpty())

	// Output:
	// false
	// true
}

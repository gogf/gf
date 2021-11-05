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
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/util/gconv"
)

func ExampleSortedStrArray_Walk() {
	var array garray.SortedStrArray
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

func ExampleNewSortedStrArray() {
	s := garray.NewSortedStrArray()
	s.Append("b")
	s.Append("d")
	s.Append("c")
	s.Append("a")
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
}

func ExampleNewSortedStrArraySize() {
	s := garray.NewSortedStrArraySize(3)
	s.SetArray([]string{"b", "d", "a", "c"})
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [a b c d] 4 4
}

func ExampleNewStrArrayFromCopy() {
	s := garray.NewSortedStrArrayFromCopy(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
}

func ExampleSortedStrArray_At() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "d", "c", "a"})
	sAt := s.At(2)
	fmt.Println(s)
	fmt.Println(sAt)

	// Output:
	// ["a","b","c","d"]
	// c

}

func ExampleSortedStrArray_Get() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "d", "c", "a", "e"})
	sGet, sBool := s.Get(3)
	fmt.Println(s)
	fmt.Println(sGet, sBool)

	// Output:
	// ["a","b","c","d","e"]
	// d true
}

func ExampleSortedStrArray_SetArray() {
	s := garray.NewSortedStrArray()
	s.SetArray([]string{"b", "d", "a", "c"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
}

func ExampleSortedStrArray_SetUnique() {
	s := garray.NewSortedStrArray()
	s.SetArray([]string{"b", "d", "a", "c", "c", "a"})
	fmt.Println(s.SetUnique(true))

	// Output:
	// ["a","b","c","d"]
}

func ExampleSortedStrArray_Sum() {
	s := garray.NewSortedStrArray()
	s.SetArray([]string{"5", "3", "2"})
	fmt.Println(s)
	a := s.Sum()
	fmt.Println(a)

	// Output:
	// ["2","3","5"]
	// 10
}

func ExampleSortedStrArray_Sort() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"b", "d", "a", "c"})
	fmt.Println(s)
	a := s.Sort()
	fmt.Println(a)

	// Output:
	// ["a","b","c","d"]
	// ["a","b","c","d"]
}

func ExampleSortedStrArray_Remove() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s.Slice())
	s.Remove(1)
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
	// [a c d]
}

func ExampleSortedStrArray_RemoveValue() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s.Slice())
	s.RemoveValue("b")
	fmt.Println(s.Slice())

	// Output:
	// [a b c d]
	// [a c d]
}

func ExampleSortedStrArray_PopLeft() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	r, _ := s.PopLeft()
	fmt.Println(r)
	fmt.Println(s.Slice())

	// Output:
	// a
	// [b c d]
}

func ExampleSortedStrArray_PopRight() {
	s := garray.NewSortedStrArray()
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

func ExampleSortedStrArray_PopRights() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.PopRights(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [g h]
	// ["a","b","c","d","e","f"]
}

func ExampleSortedStrArray_Rand() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r, _ := s.PopRand()
	fmt.Println(r)
	fmt.Println(s)

	// May Output:
	// b
	// ["a","c","d","e","f","g","h"]
}

func ExampleSortedStrArray_PopRands() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.PopRands(2)
	fmt.Println(r)
	fmt.Println(s)

	// May Output:
	// [d a]
	// ["b","c","e","f","g","h"]
}

func ExampleSortedStrArray_PopLefts() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.PopLefts(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [a b]
	// ["c","d","e","f","g","h"]
}

func ExampleSortedStrArray_Range() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Range(2, 5)
	fmt.Println(r)

	// Output:
	// [c d e]
}

func ExampleSortedStrArray_SubSlice() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.SubSlice(3, 4)
	fmt.Println(s.Slice())
	fmt.Println(r)

	// Output:
	// [a b c d e f g h]
	// [d e f g]
}

func ExampleSortedStrArray_Add() {
	s := garray.NewSortedStrArray()
	s.Add("b", "d", "c", "a")
	fmt.Println(s)

	// Output:
	// ["a","b","c","d"]
}

func ExampleSortedStrArray_Append() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"b", "d", "c", "a"})
	fmt.Println(s)
	s.Append("f", "e", "g")
	fmt.Println(s)

	// Output:
	// ["a","b","c","d"]
	// ["a","b","c","d","e","f","g"]
}

func ExampleSortedStrArray_Len() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.Len())

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// 8
}

func ExampleSortedStrArray_Slice() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d e f g h]
}

func ExampleSortedStrArray_Interfaces() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Interfaces()
	fmt.Println(r)

	// Output:
	// [a b c d e f g h]
}

func ExampleSortedStrArray_Clone() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Clone()
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// ["a","b","c","d","e","f","g","h"]
}

func ExampleSortedStrArray_Clear() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.Clear())
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// []
	// []
}

func ExampleSortedStrArray_Contains() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s.Contains("e"))
	fmt.Println(s.Contains("E"))
	fmt.Println(s.Contains("z"))

	// Output:
	// true
	// false
	// false
}

func ExampleSortedStrArray_ContainsI() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.ContainsI("E"))
	fmt.Println(s.ContainsI("z"))

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// true
	// false
}

func ExampleSortedStrArray_Search() {
	s := garray.NewSortedStrArray()
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

func ExampleSortedStrArray_Unique() {
	s := garray.NewSortedStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s)
	fmt.Println(s.Unique())

	// Output:
	// ["a","b","c","c","c","d","d"]
	// ["a","b","c","d"]
}

func ExampleSortedStrArray_LockFunc() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
	s.LockFunc(func(array []string) {
		array[len(array)-1] = "GF fans"
	})
	fmt.Println(s)

	// Output:
	// ["a","b","GF fans"]
}

func ExampleSortedStrArray_RLockFunc() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
	s.RLockFunc(func(array []string) {
		array[len(array)-1] = "GF fans"
		fmt.Println(array[len(array)-1])
	})
	fmt.Println(s)

	// Output:
	// GF fans
	// ["a","b","GF fans"]
}

func ExampleSortedStrArray_Merge() {
	s1 := garray.NewSortedStrArray()
	s2 := garray.NewSortedStrArray()
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

func ExampleSortedStrArray_Chunk() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	r := s.Chunk(3)
	fmt.Println(r)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleSortedStrArray_Rands() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s)
	fmt.Println(s.Rands(3))

	// May Output:
	// ["a","b","c","d","e","f","g","h"]
	// [h g c]
}

func ExampleSortedStrArray_Join() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"c", "b", "a", "d", "f", "e", "h", "g"})
	fmt.Println(s.Join(","))

	// Output:
	// a,b,c,d,e,f,g,h
}

func ExampleSortedStrArray_CountValues() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s.CountValues())

	// Output:
	// map[a:1 b:1 c:3 d:2]
}

func ExampleSortedStrArray_Iterator() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
	s.Iterator(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleSortedStrArray_IteratorAsc() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
	s.IteratorAsc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleSortedStrArray_IteratorDesc() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
	s.IteratorDesc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 2 c
	// 1 b
	// 0 a
}

func ExampleSortedStrArray_String() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
	fmt.Println(s.String())

	// Output:
	// ["a","b","c"]
}

func ExampleSortedStrArray_MarshalJSON() {
	type Student struct {
		ID     int
		Name   string
		Levels garray.SortedStrArray
	}
	r := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "c", "a"})
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

func ExampleSortedStrArray_UnmarshalJSON() {
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

func ExampleSortedStrArray_UnmarshalValue() {
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

func ExampleSortedStrArray_FilterEmpty() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "a", "", "c", "", "", "d"})
	fmt.Println(s)
	fmt.Println(s.FilterEmpty())

	// Output:
	// ["","","","a","b","c","d"]
	// ["a","b","c","d"]
}

func ExampleSortedStrArray_IsEmpty() {
	s := garray.NewSortedStrArrayFrom(g.SliceStr{"b", "a", "", "c", "", "", "d"})
	fmt.Println(s.IsEmpty())
	s1 := garray.NewSortedStrArray()
	fmt.Println(s1.IsEmpty())

	// Output:
	// false
	// true
}

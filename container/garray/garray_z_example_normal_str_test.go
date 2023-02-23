// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/internal/empty"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

func ExampleStrArray_Walk() {
	var array garray.StrArray
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

func ExampleNewStrArray() {
	s := garray.NewStrArray()
	s.Append("We")
	s.Append("are")
	s.Append("GF")
	s.Append("fans")
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans]
}

func ExampleNewStrArraySize() {
	s := garray.NewStrArraySize(3, 5)
	s.Set(0, "We")
	s.Set(1, "are")
	s.Set(2, "GF")
	s.Set(3, "fans")
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [We are GF] 3 5
}

func ExampleNewStrArrayFrom() {
	s := garray.NewStrArrayFrom(g.SliceStr{"We", "are", "GF", "fans", "!"})
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [We are GF fans !] 5 5
}

func ExampleStrArray_At() {
	s := garray.NewStrArrayFrom(g.SliceStr{"We", "are", "GF", "fans", "!"})
	sAt := s.At(2)
	fmt.Println(sAt)

	// Output:
	// GF
}

func ExampleStrArray_Get() {
	s := garray.NewStrArrayFrom(g.SliceStr{"We", "are", "GF", "fans", "!"})
	sGet, sBool := s.Get(3)
	fmt.Println(sGet, sBool)

	// Output:
	// fans true
}

func ExampleStrArray_Set() {
	s := garray.NewStrArraySize(3, 5)
	s.Set(0, "We")
	s.Set(1, "are")
	s.Set(2, "GF")
	s.Set(3, "fans")
	fmt.Println(s.Slice())

	// Output:
	// [We are GF]
}

func ExampleStrArray_SetArray() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"We", "are", "GF", "fans", "!"})
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans !]
}

func ExampleStrArray_Replace() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"We", "are", "GF", "fans", "!"})
	fmt.Println(s.Slice())
	s.Replace(g.SliceStr{"Happy", "coding"})
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans !]
	// [Happy coding GF fans !]
}

func ExampleStrArray_Sum() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"3", "5", "10"})
	a := s.Sum()
	fmt.Println(a)

	// Output:
	// 18
}

func ExampleStrArray_Sort() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"b", "d", "a", "c"})
	a := s.Sort()
	fmt.Println(a)

	// Output:
	// ["a","b","c","d"]
}

func ExampleStrArray_SortFunc() {
	s := garray.NewStrArrayFrom(g.SliceStr{"b", "c", "a"})
	fmt.Println(s)
	s.SortFunc(func(v1, v2 string) bool {
		return gstr.Compare(v1, v2) > 0
	})
	fmt.Println(s)
	s.SortFunc(func(v1, v2 string) bool {
		return gstr.Compare(v1, v2) < 0
	})
	fmt.Println(s)

	// Output:
	// ["b","c","a"]
	// ["c","b","a"]
	// ["a","b","c"]
}

func ExampleStrArray_InsertBefore() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.InsertBefore(1, "here")
	fmt.Println(s.Slice())

	// Output:
	// [a here b c d]
}

func ExampleStrArray_InsertAfter() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.InsertAfter(1, "here")
	fmt.Println(s.Slice())

	// Output:
	// [a b here c d]
}

func ExampleStrArray_Remove() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.Remove(1)
	fmt.Println(s.Slice())

	// Output:
	// [a c d]
}

func ExampleStrArray_RemoveValue() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.RemoveValue("b")
	fmt.Println(s.Slice())

	// Output:
	// [a c d]
}

func ExampleStrArray_PushLeft() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PushLeft("We", "are", "GF", "fans")
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans a b c d]
}

func ExampleStrArray_PushRight() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PushRight("We", "are", "GF", "fans")
	fmt.Println(s.Slice())

	// Output:
	// [a b c d We are GF fans]
}

func ExampleStrArray_PopLeft() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PopLeft()
	fmt.Println(s.Slice())

	// Output:
	// [b c d]
}

func ExampleStrArray_PopRight() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PopRight()
	fmt.Println(s.Slice())

	// Output:
	// [a b c]
}

func ExampleStrArray_PopRand() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r, _ := s.PopRand()
	fmt.Println(r)

	// May Output:
	// e
}

func ExampleStrArray_PopRands() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.PopRands(2)
	fmt.Println(r)

	// May Output:
	// [e c]
}

func ExampleStrArray_PopLefts() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.PopLefts(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [a b]
	// ["c","d","e","f","g","h"]
}

func ExampleStrArray_PopRights() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.PopRights(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [g h]
	// ["a","b","c","d","e","f"]
}

func ExampleStrArray_Range() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Range(2, 5)
	fmt.Println(r)

	// Output:
	// [c d e]
}

func ExampleStrArray_SubSlice() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.SubSlice(3, 4)
	fmt.Println(r)

	// Output:
	// [d e f g]
}

func ExampleStrArray_Append() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"We", "are", "GF", "fans"})
	s.Append("a", "b", "c")
	fmt.Println(s)

	// Output:
	// ["We","are","GF","fans","a","b","c"]
}

func ExampleStrArray_Len() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Len())

	// Output:
	// 8
}

func ExampleStrArray_Slice() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d e f g h]
}

func ExampleStrArray_Interfaces() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Interfaces()
	fmt.Println(r)

	// Output:
	// [a b c d e f g h]
}

func ExampleStrArray_Clone() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Clone()
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// ["a","b","c","d","e","f","g","h"]
}

func ExampleStrArray_Clear() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s)
	fmt.Println(s.Clear())
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// []
	// []
}

func ExampleStrArray_Contains() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Contains("e"))
	fmt.Println(s.Contains("z"))

	// Output:
	// true
	// false
}

func ExampleStrArray_ContainsI() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.ContainsI("E"))
	fmt.Println(s.ContainsI("z"))

	// Output:
	// true
	// false
}

func ExampleStrArray_Search() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Search("e"))
	fmt.Println(s.Search("z"))

	// Output:
	// 4
	// -1
}

func ExampleStrArray_Unique() {
	s := garray.NewStrArray()
	s.SetArray(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s.Unique())

	// Output:
	// ["a","b","c","d"]
}

func ExampleStrArray_LockFunc() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	s.LockFunc(func(array []string) {
		array[len(array)-1] = "GF fans"
	})
	fmt.Println(s)

	// Output:
	// ["a","b","GF fans"]
}

func ExampleStrArray_RLockFunc() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e"})
	s.RLockFunc(func(array []string) {
		for i := 0; i < len(array); i++ {
			fmt.Println(array[i])
		}
	})

	// Output:
	// a
	// b
	// c
	// d
	// e
}

func ExampleStrArray_Merge() {
	s1 := garray.NewStrArray()
	s2 := garray.NewStrArray()
	s1.SetArray(g.SliceStr{"a", "b", "c"})
	s2.SetArray(g.SliceStr{"d", "e", "f"})
	s1.Merge(s2)
	fmt.Println(s1)

	// Output:
	// ["a","b","c","d","e","f"]
}

func ExampleStrArray_Fill() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	s.Fill(2, 3, "here")
	fmt.Println(s)

	// Output:
	// ["a","b","here","here","here","f","g","h"]
}

func ExampleStrArray_Chunk() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Chunk(3)
	fmt.Println(r)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleStrArray_Pad() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	s.Pad(7, "here")
	fmt.Println(s)
	s.Pad(-10, "there")
	fmt.Println(s)

	// Output:
	// ["a","b","c","here","here","here","here"]
	// ["there","there","there","a","b","c","here","here","here","here"]
}

func ExampleStrArray_Rand() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Rand())

	// May Output:
	// c true
}

func ExampleStrArray_Rands() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Rands(3))

	// May Output:
	// [e h e]
}

func ExampleStrArray_Shuffle() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Shuffle())

	// May Output:
	// ["a","c","e","d","b","g","f","h"]
}

func ExampleStrArray_Reverse() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Reverse())

	// Output:
	// ["h","g","f","e","d","c","b","a"]
}

func ExampleStrArray_Join() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	fmt.Println(s.Join(","))

	// Output:
	// a,b,c
}

func ExampleStrArray_CountValues() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s.CountValues())

	// Output:
	// map[a:1 b:1 c:3 d:2]
}

func ExampleStrArray_Iterator() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	s.Iterator(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleStrArray_IteratorAsc() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	s.IteratorAsc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleStrArray_IteratorDesc() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	s.IteratorDesc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 2 c
	// 1 b
	// 0 a
}

func ExampleStrArray_String() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "c"})
	fmt.Println(s.String())

	// Output:
	// ["a","b","c"]
}

func ExampleStrArray_MarshalJSON() {
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

	// Output:
	// {"Id":1,"Name":"john","Lessons":["Math","English","Music"]}
}

func ExampleStrArray_UnmarshalJSON() {
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

func ExampleStrArray_UnmarshalValue() {
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

func ExampleStrArray_Filter() {
	s := garray.NewStrArrayFrom(g.SliceStr{"Math", "English", "Sport"})
	s1 := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "", "c", "", "", "d"})
	fmt.Println(s1.Filter(func(value string, index int) bool {
		return empty.IsEmpty(value)
	}))

	fmt.Println(s.Filter(func(value string, index int) bool {
		return strings.Contains(value, "h")
	}))

	// Output:
	// ["a","b","c","d"]
	// ["Sport"]
}

func ExampleStrArray_FilterEmpty() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "", "c", "", "", "d"})
	fmt.Println(s.FilterEmpty())

	// Output:
	// ["a","b","c","d"]
}

func ExampleStrArray_IsEmpty() {
	s := garray.NewStrArrayFrom(g.SliceStr{"a", "b", "", "c", "", "", "d"})
	fmt.Println(s.IsEmpty())
	s1 := garray.NewStrArray()
	fmt.Println(s1.IsEmpty())

	// Output:
	// false
	// true
}

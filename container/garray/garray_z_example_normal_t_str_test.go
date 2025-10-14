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

func ExampleTArrayStr_Walk() {
	var array garray.TArray[string]
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

func ExampleNewTArrayStr() {
	s := garray.NewTArray[string]()
	s.Append("We")
	s.Append("are")
	s.Append("GF")
	s.Append("fans")
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans]
}

func ExampleNewTArrayStrSize() {
	s := garray.NewTArraySize[string](3, 5)
	s.Set(0, "We")
	s.Set(1, "are")
	s.Set(2, "GF")
	s.Set(3, "fans")
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [We are GF] 3 5
}

func ExampleNewTArrayStrFrom() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"We", "are", "GF", "fans", "!"})
	fmt.Println(s.Slice(), s.Len(), cap(s.Slice()))

	// Output:
	// [We are GF fans !] 5 5
}

func ExampleTArrayStr_At() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"We", "are", "GF", "fans", "!"})
	sAt := s.At(2)
	fmt.Println(sAt)

	// Output:
	// GF
}

func ExampleTArrayStr_Get() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"We", "are", "GF", "fans", "!"})
	sGet, sBool := s.Get(3)
	fmt.Println(sGet, sBool)

	// Output:
	// fans true
}

func ExampleTArrayStr_Set() {
	s := garray.NewTArraySize[string](3, 5)
	s.Set(0, "We")
	s.Set(1, "are")
	s.Set(2, "GF")
	s.Set(3, "fans")
	fmt.Println(s.Slice())

	// Output:
	// [We are GF]
}

func ExampleTArrayStr_SetArray() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"We", "are", "GF", "fans", "!"})
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans !]
}

func ExampleTArrayStr_Replace() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"We", "are", "GF", "fans", "!"})
	fmt.Println(s.Slice())
	s.Replace(g.SliceStr{"Happy", "coding"})
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans !]
	// [Happy coding GF fans !]
}

func ExampleTArrayStr_Sum() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"3", "5", "10"})
	a := s.Sum()
	fmt.Println(a)

	// Output:
	// 18
}

func ExampleTArrayStr_SortFunc() {
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

	// Output:
	// ["b","c","a"]
	// ["c","b","a"]
	// ["a","b","c"]
}

func ExampleTArrayStr_InsertBefore() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.InsertBefore(1, "here")
	fmt.Println(s.Slice())

	// Output:
	// [a here b c d]
}

func ExampleTArrayStr_InsertAfter() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.InsertAfter(1, "here")
	fmt.Println(s.Slice())

	// Output:
	// [a b here c d]
}

func ExampleTArrayStr_Remove() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.Remove(1)
	fmt.Println(s.Slice())

	// Output:
	// [a c d]
}

func ExampleTArrayStr_RemoveValue() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.RemoveValue("b")
	fmt.Println(s.Slice())

	// Output:
	// [a c d]
}

func ExampleTArrayStr_PushLeft() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PushLeft("We", "are", "GF", "fans")
	fmt.Println(s.Slice())

	// Output:
	// [We are GF fans a b c d]
}

func ExampleTArrayStr_PushRight() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PushRight("We", "are", "GF", "fans")
	fmt.Println(s.Slice())

	// Output:
	// [a b c d We are GF fans]
}

func ExampleTArrayStr_PopLeft() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PopLeft()
	fmt.Println(s.Slice())

	// Output:
	// [b c d]
}

func ExampleTArrayStr_PopRight() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d"})
	s.PopRight()
	fmt.Println(s.Slice())

	// Output:
	// [a b c]
}

func ExampleTArrayStr_PopRand() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r, _ := s.PopRand()
	fmt.Println(r)

	// May Output:
	// e
}

func ExampleTArrayStr_PopRands() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.PopRands(2)
	fmt.Println(r)

	// May Output:
	// [e c]
}

func ExampleTArrayStr_PopLefts() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.PopLefts(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [a b]
	// ["c","d","e","f","g","h"]
}

func ExampleTArrayStr_PopRights() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.PopRights(2)
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// [g h]
	// ["a","b","c","d","e","f"]
}

func ExampleTArrayStr_Range() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Range(2, 5)
	fmt.Println(r)

	// Output:
	// [c d e]
}

func ExampleTArrayStr_SubSlice() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.SubSlice(3, 4)
	fmt.Println(r)

	// Output:
	// [d e f g]
}

func ExampleTArrayStr_Append() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"We", "are", "GF", "fans"})
	s.Append("a", "b", "c")
	fmt.Println(s)

	// Output:
	// ["We","are","GF","fans","a","b","c"]
}

func ExampleTArrayStr_Len() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Len())

	// Output:
	// 8
}

func ExampleTArrayStr_Slice() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Slice())

	// Output:
	// [a b c d e f g h]
}

func ExampleTArrayStr_Interfaces() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Interfaces()
	fmt.Println(r)

	// Output:
	// [a b c d e f g h]
}

func ExampleTArrayStr_Clone() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Clone()
	fmt.Println(r)
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// ["a","b","c","d","e","f","g","h"]
}

func ExampleTArrayStr_Clear() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s)
	fmt.Println(s.Clear())
	fmt.Println(s)

	// Output:
	// ["a","b","c","d","e","f","g","h"]
	// []
	// []
}

func ExampleTArrayStr_Contains() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Contains("e"))
	fmt.Println(s.Contains("z"))

	// Output:
	// true
	// false
}

func ExampleTArrayStr_Search() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Search("e"))
	fmt.Println(s.Search("z"))

	// Output:
	// 4
	// -1
}

func ExampleTArrayStr_Unique() {
	s := garray.NewTArray[string]()
	s.SetArray(g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s.Unique())

	// Output:
	// ["a","b","c","d"]
}

func ExampleTArrayStr_LockFunc() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	s.LockFunc(func(array []string) {
		array[len(array)-1] = "GF fans"
	})
	fmt.Println(s)

	// Output:
	// ["a","b","GF fans"]
}

func ExampleTArrayStr_RLockFunc() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e"})
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

func ExampleTArrayStr_Merge() {
	s1 := garray.NewTArray[string]()
	s2 := garray.NewTArray[string]()
	s1.SetArray(g.SliceStr{"a", "b", "c"})
	s2.SetArray(g.SliceStr{"d", "e", "f"})
	s1.Merge(s2)
	fmt.Println(s1)

	// Output:
	// ["a","b","c","d","e","f"]
}

func ExampleTArrayStr_Fill() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	s.Fill(2, 3, "here")
	fmt.Println(s)

	// Output:
	// ["a","b","here","here","here","f","g","h"]
}

func ExampleTArrayStr_Chunk() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	r := s.Chunk(3)
	fmt.Println(r)

	// Output:
	// [[a b c] [d e f] [g h]]
}

func ExampleTArrayStr_Pad() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	s.Pad(7, "here")
	fmt.Println(s)
	s.Pad(-10, "there")
	fmt.Println(s)

	// Output:
	// ["a","b","c","here","here","here","here"]
	// ["there","there","there","a","b","c","here","here","here","here"]
}

func ExampleTArrayStr_Rand() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Rand())

	// May Output:
	// c true
}

func ExampleTArrayStr_Rands() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Rands(3))

	// May Output:
	// [e h e]
}

func ExampleTArrayStr_Shuffle() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Shuffle())

	// May Output:
	// ["a","c","e","d","b","g","f","h"]
}

func ExampleTArrayStr_Reverse() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "d", "e", "f", "g", "h"})
	fmt.Println(s.Reverse())

	// Output:
	// ["h","g","f","e","d","c","b","a"]
}

func ExampleTArrayStr_Join() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	fmt.Println(s.Join(","))

	// Output:
	// a,b,c
}

func ExampleTArrayStr_CountValues() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c", "c", "c", "d", "d"})
	fmt.Println(s.CountValues())

	// Output:
	// map[a:1 b:1 c:3 d:2]
}

func ExampleTArrayStr_Iterator() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	s.Iterator(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleTArrayStr_IteratorAsc() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	s.IteratorAsc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
}

func ExampleTArrayStr_IteratorDesc() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	s.IteratorDesc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 2 c
	// 1 b
	// 0 a
}

func ExampleTArrayStr_String() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "c"})
	fmt.Println(s.String())

	// Output:
	// ["a","b","c"]
}

func ExampleTArrayStr_MarshalJSON() {
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

func ExampleTArrayStr_UnmarshalJSON() {
	b := []byte(`{"Id":1,"Name":"john","Lessons":["Math","English","Sport"]}`)
	type Student struct {
		Id      int
		Name    string
		Lessons *garray.TArray[string]
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// Output:
	// {1 john ["Math","English","Sport"]}
}

func ExampleTArrayStr_UnmarshalValue() {
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

	// Output:
	// &{john ["Math","English","Sport"]}
	// &{john ["Math","English","Sport"]}
}

func ExampleTArrayStr_Filter() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"Math", "English", "Sport"})
	s1 := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "", "c", "", "", "d"})
	fmt.Println(s1.Filter(func(index int, value string) bool {
		return empty.IsEmpty(value)
	}))

	fmt.Println(s.Filter(func(index int, value string) bool {
		return strings.Contains(value, "h")
	}))

	// Output:
	// ["a","b","c","d"]
	// ["Sport"]
}

func ExampleTArrayStr_FilterEmpty() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "", "c", "", "", "d"})
	fmt.Println(s.FilterEmpty())

	// Output:
	// ["a","b","c","d"]
}

func ExampleTArrayStr_IsEmpty() {
	s := garray.NewTArrayFrom[string](g.SliceStr{"a", "b", "", "c", "", "", "d"})
	fmt.Println(s.IsEmpty())
	s1 := garray.NewTArray[string]()
	fmt.Println(s1.IsEmpty())

	// Output:
	// false
	// true
}

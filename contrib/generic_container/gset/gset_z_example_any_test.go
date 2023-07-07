// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gset

import (
	"fmt"

	"github.com/gogf/gf/contrib/generic_container/v2/internal/json"
	"github.com/gogf/gf/v2/frame/g"
)

func ExampleSet_Intersect() {
	s1 := NewFrom[int]([]int{1, 2, 3})
	s2 := NewFrom[int]([]int{4, 5, 6})
	s3 := NewFrom[int]([]int{1, 2, 3, 4, 5, 6, 7})

	fmt.Println(s3.Intersect(s1).Slice())
	fmt.Println(s3.Diff(s1).Slice())
	fmt.Println(s1.Union(s2).Slice())
	fmt.Println(s1.Complement(s3).Slice())

	// May Output:
	// [2 3 1]
	// [5 6 7 4]
	// [6 1 2 3 4 5]
	// [4 5 6 7]
}

func ExampleSet_Diff() {
	s1 := NewFrom([]int{1, 2, 3})
	s2 := NewFrom([]int{4, 5, 6})
	s3 := NewFrom([]int{1, 2, 3, 4, 5, 6, 7})

	fmt.Println(s3.Intersect(s1).Slice())
	fmt.Println(s3.Diff(s1).Slice())
	fmt.Println(s1.Union(s2).Slice())
	fmt.Println(s1.Complement(s3).Slice())

	// May Output:
	// [2 3 1]
	// [5 6 7 4]
	// [6 1 2 3 4 5]
	// [4 5 6 7]
}

func ExampleSet_Union() {
	s1 := NewFrom([]int{1, 2, 3})
	s2 := NewFrom([]int{4, 5, 6})
	s3 := NewFrom([]int{1, 2, 3, 4, 5, 6, 7})

	fmt.Println(s3.Intersect(s1).Slice())
	fmt.Println(s3.Diff(s1).Slice())
	fmt.Println(s1.Union(s2).Slice())
	fmt.Println(s1.Complement(s3).Slice())

	// May Output:
	// [2 3 1]
	// [5 6 7 4]
	// [6 1 2 3 4 5]
	// [4 5 6 7]
}

func ExampleSet_Complement() {
	s1 := NewFrom([]int{1, 2, 3})
	s2 := NewFrom([]int{4, 5, 6})
	s3 := NewFrom([]int{1, 2, 3, 4, 5, 6, 7})

	fmt.Println(s3.Intersect(s1).Slice())
	fmt.Println(s3.Diff(s1).Slice())
	fmt.Println(s1.Union(s2).Slice())
	fmt.Println(s1.Complement(s3).Slice())

	// May Output:
	// [2 3 1]
	// [5 6 7 4]
	// [6 1 2 3 4 5]
	// [4 5 6 7]
}

func ExampleSet_IsSubsetOf() {
	var s1, s2 HashSet[int]
	s1.Add([]int{1, 2, 3}...)
	s2.Add([]int{2, 3}...)
	fmt.Println(s1.IsSubsetOf(&s2))
	fmt.Println(s2.IsSubsetOf(&s1))

	// Output:
	// false
	// true
}

func ExampleSet_AddIfNotExist() {
	var set HashSet[int]
	fmt.Println(set.AddIfNotExist(1))
	fmt.Println(set.AddIfNotExist(1))
	fmt.Println(set.Slice())

	// Output:
	// true
	// false
	// [1]
}

func ExampleSet_Pop() {
	var set HashSet[int]
	set.Add(1, 2, 3, 4)
	fmt.Println(set.Pop())
	fmt.Println(set.Pops(2))
	fmt.Println(set.Size())

	// May Output:
	// 1
	// [2 3]
	// 1
}

func ExampleSet_Pops() {
	var set HashSet[int]
	set.Add(1, 2, 3, 4)
	fmt.Println(set.Pop())
	fmt.Println(set.Pops(2))
	fmt.Println(set.Size())

	// May Output:
	// 1
	// [2 3]
	// 1
}

func ExampleSet_Join() {
	var set HashSet[string]
	set.Add("a", "b", "c", "d")
	fmt.Println(set.Join(","))

	// May Output:
	// a,b,c,d
}

func ExampleSet_Contains() {
	var set HashSet[string]
	set.Add("a")
	fmt.Println(set.Contains("a"))
	fmt.Println(set.Contains("A"))
	fmt.Println(set.ContainsI("A"))

	// Output:
	// true
	// false
	// true
}

func ExampleSet_ContainsI() {
	var set HashSet[string]
	set.Add("a")
	fmt.Println(set.Contains("a"))
	fmt.Println(set.Contains("A"))
	fmt.Println(set.ContainsI("A"))

	// Output:
	// true
	// false
	// true
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func ExampleIntSet_UnmarshalJSON() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[100,99,98]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *HashSet[int]
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// May Output:
	// {1 john [100,99,98]}
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func ExampleIntSet_UnmarshalValue() {
	b := []byte(`{"Id":1,"Name":"john","Scores":100,99,98}`)
	type Student struct {
		Id     int
		Name   string
		Scores *HashSet[int]
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// May Output:
	// {1 john [100,99,98]}
}

// Walk applies a user supplied function `f` to every item of set.
func ExampleIntSet_Walk() {
	var (
		set   HashSet[int]
		names = g.SliceInt{1, 0}
		delta = 10
	)
	set.Add(names...)
	// Add prefix for given table names.
	set.Walk(func(item int) int {
		return delta + item
	})
	fmt.Println(set.Slice())

	// May Output:
	// [12 60]
}

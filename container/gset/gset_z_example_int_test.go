// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gset_test

import (
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
)

// New create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func ExampleNewIntSet() {
	intSet := gset.NewIntSet()
	intSet.Add([]int{1, 2, 3}...)
	fmt.Println(intSet.Slice())

	// May Output:
	// [2 1 3]
}

// NewIntSetFrom  returns a new set from `items`.
func ExampleNewFrom() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3})
	fmt.Println(intSet.Slice())

	// May Output:
	// [2 1 3]
}

// Add adds one or multiple items to the set.
func ExampleIntSet_Add() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3})
	intSet.Add(1)
	fmt.Println(intSet.Slice())
	fmt.Println(intSet.AddIfNotExist(1))

	// May Output:
	// [1 2 3]
	// false
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set,
// or else it does nothing and returns false.
func ExampleIntSet_AddIfNotExist() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3})
	intSet.Add(1)
	fmt.Println(intSet.Slice())
	fmt.Println(intSet.AddIfNotExist(1))

	// May Output:
	// [1 2 3]
	// false
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and function `f` returns true,
// or else it does nothing and returns false.
// Note that, the function `f` is executed without writing lock.
func ExampleIntSet_AddIfNotExistFunc() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3})
	intSet.Add(1)
	fmt.Println(intSet.Slice())
	fmt.Println(intSet.AddIfNotExistFunc(5, func() bool {
		return true
	}))

	// May Output:
	// [1 2 3]
	// true
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and function `f` returns true,
// or else it does nothing and returns false.
// Note that, the function `f` is executed without writing lock.
func ExampleIntSet_AddIfNotExistFuncLock() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3})
	intSet.Add(1)
	fmt.Println(intSet.Slice())
	fmt.Println(intSet.AddIfNotExistFuncLock(4, func() bool {
		return true
	}))

	// May Output:
	// [1 2 3]
	// true
}

// Clear deletes all items of the set.
func ExampleIntSet_Clear() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3})
	fmt.Println(intSet.Size())
	intSet.Clear()
	fmt.Println(intSet.Size())

	// Output:
	// 3
	// 0
}

// Complement returns a new set which is the complement from `set` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `set`.
// It returns the difference between `full` and `set` if the given set `full` is not the full set of `set`.
func ExampleIntSet_Complement() {
	intSet := gset.NewIntSetFrom([]int{1, 2, 3, 4, 5})
	s := gset.NewIntSetFrom([]int{1, 2, 3})
	fmt.Println(s.Complement(intSet).Slice())

	// May Output:
	// [4 5]
}

// Contains checks whether the set contains `item`.
func ExampleIntSet_Contains() {
	var set1 gset.IntSet
	set1.Add(1, 4, 5, 6, 7)
	fmt.Println(set1.Contains(1))

	var set2 gset.IntSet
	set2.Add(1, 4, 5, 6, 7)
	fmt.Println(set2.Contains(8))

	// Output:
	// true
	// false
}

// Diff returns a new set which is the difference set from `set` to `other`.
// Which means, all the items in `newSet` are in `set` but not in `other`.
func ExampleIntSet_Diff() {
	s1 := gset.NewIntSetFrom([]int{1, 2, 3})
	s2 := gset.NewIntSetFrom([]int{1, 2, 3, 4})
	fmt.Println(s2.Diff(s1).Slice())

	// Output:
	// [4]
}

// Equal checks whether the two sets equal.
func ExampleIntSet_Equal() {
	s1 := gset.NewIntSetFrom([]int{1, 2, 3})
	s2 := gset.NewIntSetFrom([]int{1, 2, 3, 4})
	fmt.Println(s2.Equal(s1))

	s3 := gset.NewIntSetFrom([]int{1, 2, 3})
	s4 := gset.NewIntSetFrom([]int{1, 2, 3})
	fmt.Println(s3.Equal(s4))

	// Output:
	// false
	// true
}

// Intersect returns a new set which is the intersection from `set` to `other`.
// Which means, all the items in `newSet` are in `set` and also in `other`.
func ExampleIntSet_Intersect() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3}...)
	var s2 gset.IntSet
	s2.Add([]int{1, 2, 3, 4}...)
	fmt.Println(s2.Intersect(s1).Slice())

	// May Output:
	// [1 2 3]
}

// IsSubsetOf checks whether the current set is a sub-set of `other`
func ExampleIntSet_IsSubsetOf() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	var s2 gset.IntSet
	s2.Add([]int{1, 2, 4}...)
	fmt.Println(s2.IsSubsetOf(s1))

	// Output:
	// true
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func ExampleIntSet_Iterator() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	s1.Iterator(func(v int) bool {
		fmt.Println("Iterator", v)
		return true
	})
	// May Output:
	// Iterator 2
	// Iterator 3
	// Iterator 1
	// Iterator 4
}

// Join joins items with a string `glue`.
func ExampleIntSet_Join() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	fmt.Println(s1.Join(","))

	// May Output:
	// 3,4,1,2
}

// LockFunc locks writing with callback function `f`.
func ExampleIntSet_LockFunc() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2}...)
	s1.LockFunc(func(m map[int]struct{}) {
		m[3] = struct{}{}
	})
	fmt.Println(s1.Slice())

	// May Output
	// [2 3 1]
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func ExampleIntSet_MarshalJSON() {
	type Student struct {
		Id     int
		Name   string
		Scores *gset.IntSet
	}
	s := Student{
		Id:     1,
		Name:   "john",
		Scores: gset.NewIntSetFrom([]int{100, 99, 98}),
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))

	// May Output:
	// {"Id":1,"Name":"john","Scores":[100,99,98]}
}

// Merge adds items from `others` sets into `set`.
func ExampleIntSet_Merge() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)

	s2 := gset.NewIntSet()
	fmt.Println(s1.Merge(s2).Slice())

	// May Output:
	// [1 2 3 4]
}

// Pops randomly pops an item from set.
func ExampleIntSet_Pop() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)

	fmt.Println(s1.Pop())

	// May Output:
	// 1
}

// Pops randomly pops `size` items from set.
// It returns all items if size == -1.
func ExampleIntSet_Pops() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	for _, v := range s1.Pops(2) {
		fmt.Println(v)
	}

	// May Output:
	// 1
	// 2
}

// RLockFunc locks reading with callback function `f`.
func ExampleIntSet_RLockFunc() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	s1.RLockFunc(func(m map[int]struct{}) {
		fmt.Println(m)
	})

	// Output:
	// map[1:{} 2:{} 3:{} 4:{}]
}

// Remove deletes `item` from set.
func ExampleIntSet_Remove() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	s1.Remove(1)
	fmt.Println(s1.Slice())

	// May Output:
	// [3 4 2]
}

// Size returns the size of the set.
func ExampleIntSet_Size() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	fmt.Println(s1.Size())

	// Output:
	// 4
}

// Slice returns the an of items of the set as slice.
func ExampleIntSet_Slice() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	fmt.Println(s1.Slice())

	// May Output:
	// [1, 2, 3, 4]
}

// String returns items as a string, which implements like json.Marshal does.
func ExampleIntSet_String() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	fmt.Println(s1.String())

	// May Output:
	// [1,2,3,4]
}

// Sum sums items. Note: The items should be converted to int type,
// or you'd get a result that you unexpected.
func ExampleIntSet_Sum() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	fmt.Println(s1.Sum())

	// Output:
	// 10
}

// Union returns a new set which is the union of `set` and `other`.
// Which means, all the items in `newSet` are in `set` or in `other`.
func ExampleIntSet_Union() {
	s1 := gset.NewIntSet()
	s1.Add([]int{1, 2, 3, 4}...)
	s2 := gset.NewIntSet()
	s2.Add([]int{1, 2, 4}...)
	fmt.Println(s1.Union(s2).Slice())

	// May Output:
	// [3 4 1 2]
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func ExampleIntSet_UnmarshalJSON() {
	b := []byte(`{"Id":1,"Name":"john","Scores":[100,99,98]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *gset.IntSet
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
		Scores *gset.IntSet
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
		set   gset.IntSet
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

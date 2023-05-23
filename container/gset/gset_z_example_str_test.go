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

// NewStrSet create and returns a new set, which contains un-repeated items.
// The parameter `safe` is used to specify whether using set in concurrent-safety,
// which is false in default.
func ExampleNewStrSet() {
	strSet := gset.NewStrSet(true)
	strSet.Add([]string{"str1", "str2", "str3"}...)
	fmt.Println(strSet.Slice())

	// May Output:
	// [str3 str1 str2]
}

// NewStrSetFrom returns a new set from `items`.
func ExampleNewStrSetFrom() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	fmt.Println(strSet.Slice())

	// May Output:
	// [str1 str2 str3]
}

// Add adds one or multiple items to the set.
func ExampleStrSet_Add() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	strSet.Add("str")
	fmt.Println(strSet.Slice())
	fmt.Println(strSet.AddIfNotExist("str"))

	// Mya Output:
	// [str str1 str2 str3]
	// false
}

// AddIfNotExist checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set,
// or else it does nothing and returns false.
func ExampleStrSet_AddIfNotExist() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	strSet.Add("str")
	fmt.Println(strSet.Slice())
	fmt.Println(strSet.AddIfNotExist("str"))

	// Mya Output:
	// [str str1 str2 str3]
	// false
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and function `f` returns true,
// or else it does nothing and returns false.
// Note that, the function `f` is executed without writing lock.
func ExampleStrSet_AddIfNotExistFunc() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	strSet.Add("str")
	fmt.Println(strSet.Slice())
	fmt.Println(strSet.AddIfNotExistFunc("str5", func() bool {
		return true
	}))

	// May Output:
	// [str1 str2 str3 str]
	// true
}

// AddIfNotExistFunc checks whether item exists in the set,
// it adds the item to set and returns true if it does not exists in the set and function `f` returns true,
// or else it does nothing and returns false.
// Note that, the function `f` is executed without writing lock.
func ExampleStrSet_AddIfNotExistFuncLock() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	strSet.Add("str")
	fmt.Println(strSet.Slice())
	fmt.Println(strSet.AddIfNotExistFuncLock("str4", func() bool {
		return true
	}))

	// May Output:
	// [str1 str2 str3 str]
	// true
}

// Clear deletes all items of the set.
func ExampleStrSet_Clear() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	fmt.Println(strSet.Size())
	strSet.Clear()
	fmt.Println(strSet.Size())

	// Output:
	// 3
	// 0
}

// Complement returns a new set which is the complement from `set` to `full`.
// Which means, all the items in `newSet` are in `full` and not in `set`.
// It returns the difference between `full` and `set` if the given set `full` is not the full set of `set`.
func ExampleStrSet_Complement() {
	strSet := gset.NewStrSetFrom([]string{"str1", "str2", "str3", "str4", "str5"}, true)
	s := gset.NewStrSetFrom([]string{"str1", "str2", "str3"}, true)
	fmt.Println(s.Complement(strSet).Slice())

	// May Output:
	// [str4 str5]
}

// Contains checks whether the set contains `item`.
func ExampleStrSet_Contains() {
	var set gset.StrSet
	set.Add("a")
	fmt.Println(set.Contains("a"))
	fmt.Println(set.Contains("A"))

	// Output:
	// true
	// false
}

// ContainsI checks whether a value exists in the set with case-insensitively.
// Note that it internally iterates the whole set to do the comparison with case-insensitively.
func ExampleStrSet_ContainsI() {
	var set gset.StrSet
	set.Add("a")
	fmt.Println(set.ContainsI("a"))
	fmt.Println(set.ContainsI("A"))

	// Output:
	// true
	// true
}

// Diff returns a new set which is the difference set from `set` to `other`.
// Which means, all the items in `newSet` are in `set` but not in `other`.
func ExampleStrSet_Diff() {
	s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
	s2 := gset.NewStrSetFrom([]string{"a", "b", "c", "d"}, true)
	fmt.Println(s2.Diff(s1).Slice())

	// Output:
	// [d]
}

// Equal checks whether the two sets equal.
func ExampleStrSet_Equal() {
	s1 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
	s2 := gset.NewStrSetFrom([]string{"a", "b", "c", "d"}, true)
	fmt.Println(s2.Equal(s1))

	s3 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
	s4 := gset.NewStrSetFrom([]string{"a", "b", "c"}, true)
	fmt.Println(s3.Equal(s4))

	// Output:
	// false
	// true
}

// Intersect returns a new set which is the intersection from `set` to `other`.
// Which means, all the items in `newSet` are in `set` and also in `other`.
func ExampleStrSet_Intersect() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c"}...)
	var s2 gset.StrSet
	s2.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s2.Intersect(s1).Slice())

	// May Output:
	// [c a b]
}

// IsSubsetOf checks whether the current set is a sub-set of `other`
func ExampleStrSet_IsSubsetOf() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	var s2 gset.StrSet
	s2.Add([]string{"a", "b", "d"}...)
	fmt.Println(s2.IsSubsetOf(s1))

	// Output:
	// true
}

// Iterator iterates the set readonly with given callback function `f`,
// if `f` returns true then continue iterating; or false to stop.
func ExampleStrSet_Iterator() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	s1.Iterator(func(v string) bool {
		fmt.Println("Iterator", v)
		return true
	})

	// May Output:
	// Iterator a
	// Iterator b
	// Iterator c
	// Iterator d
}

// Join joins items with a string `glue`.
func ExampleStrSet_Join() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s1.Join(","))

	// May Output:
	// b,c,d,a
}

// LockFunc locks writing with callback function `f`.
func ExampleStrSet_LockFunc() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"1", "2"}...)
	s1.LockFunc(func(m map[string]struct{}) {
		m["3"] = struct{}{}
	})
	fmt.Println(s1.Slice())

	// May Output
	// [2 3 1]

}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func ExampleStrSet_MarshalJSON() {
	type Student struct {
		Id     int
		Name   string
		Scores *gset.StrSet
	}
	s := Student{
		Id:     1,
		Name:   "john",
		Scores: gset.NewStrSetFrom([]string{"100", "99", "98"}, true),
	}
	b, _ := json.Marshal(s)
	fmt.Println(string(b))

	// May Output:
	// {"Id":1,"Name":"john","Scores":["100","99","98"]}
}

// Merge adds items from `others` sets into `set`.
func ExampleStrSet_Merge() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)

	s2 := gset.NewStrSet(true)
	fmt.Println(s1.Merge(s2).Slice())

	// May Output:
	// [d a b c]
}

// Pops randomly pops an item from set.
func ExampleStrSet_Pop() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)

	fmt.Println(s1.Pop())

	// May Output:
	// a
}

// Pops randomly pops `size` items from set.
// It returns all items if size == -1.
func ExampleStrSet_Pops() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	for _, v := range s1.Pops(2) {
		fmt.Println(v)
	}

	// May Output:
	// a
	// b
}

// RLockFunc locks reading with callback function `f`.
func ExampleStrSet_RLockFunc() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	s1.RLockFunc(func(m map[string]struct{}) {
		fmt.Println(m)
	})

	// Output:
	// map[a:{} b:{} c:{} d:{}]
}

// Remove deletes `item` from set.
func ExampleStrSet_Remove() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	s1.Remove("a")
	fmt.Println(s1.Slice())

	// May Output:
	// [b c d]
}

// Size returns the size of the set.
func ExampleStrSet_Size() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s1.Size())

	// Output:
	// 4
}

// Slice returns the an of items of the set as slice.
func ExampleStrSet_Slice() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s1.Slice())

	// May Output:
	// [a,b,c,d]
}

// String returns items as a string, which implements like json.Marshal does.
func ExampleStrSet_String() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s1.String())

	// May Output:
	// "a","b","c","d"
}

// Sum sums items. Note: The items should be converted to int type,
// or you'd get a result that you unexpected.
func ExampleStrSet_Sum() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"1", "2", "3", "4"}...)
	fmt.Println(s1.Sum())

	// Output:
	// 10
}

// Union returns a new set which is the union of `set` and `other`.
// Which means, all the items in `newSet` are in `set` or in `other`.
func ExampleStrSet_Union() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	s2 := gset.NewStrSet(true)
	s2.Add([]string{"a", "b", "d"}...)
	fmt.Println(s1.Union(s2).Slice())

	// May Output:
	// [a b c d]
}

// UnmarshalJSON implements the interface UnmarshalJSON for json.Unmarshal.
func ExampleStrSet_UnmarshalJSON() {
	b := []byte(`{"Id":1,"Name":"john","Scores":["100","99","98"]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *gset.StrSet
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// May Output:
	// {1 john "99","98","100"}
}

// UnmarshalValue is an interface implement which sets any type of value for set.
func ExampleStrSet_UnmarshalValue() {
	b := []byte(`{"Id":1,"Name":"john","Scores":["100","99","98"]}`)
	type Student struct {
		Id     int
		Name   string
		Scores *gset.StrSet
	}
	s := Student{}
	json.Unmarshal(b, &s)
	fmt.Println(s)

	// May Output:
	// {1 john "99","98","100"}
}

// Walk applies a user supplied function `f` to every item of set.
func ExampleStrSet_Walk() {
	var (
		set    gset.StrSet
		names  = g.SliceStr{"user", "user_detail"}
		prefix = "gf_"
	)
	set.Add(names...)
	// Add prefix for given table names.
	set.Walk(func(item string) string {
		return prefix + item
	})
	fmt.Println(set.Slice())

	// May Output:
	// [gf_user gf_user_detail]
}

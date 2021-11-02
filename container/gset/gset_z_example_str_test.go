// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gset_test

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/frame/g"
)

func ExampleStrSet_NewStrSet() {
	strSet := gset.NewStrSet(true)
	strSet.Add([]string{"str1", "str2", "str3"}...)

	// Mya Output:
	//Iterator  str1
	//Iterator  str2
	//Iterator  str3
}

func ExampleStrSet_Add() {
	var strSet gset.StrSet
	strSet.Add([]string{"str1", "str2", "str3"}...)

	// Mya Output:
	//Iterator  str1
	//Iterator  str2
	//Iterator  str3
}

func ExampleStrSet_AddIfNotExist() {
	var strSet gset.StrSet
	fmt.Println(strSet.AddIfNotExist("str"))

	// Output:
	// true
}

func ExampleStrSet_AddIfNotExistFunc() {
	var strSet gset.StrSet
	fmt.Println(strSet.AddIfNotExistFunc("str", func() bool {
		return true
	}))

	// Output:
	// true
}

func ExampleStrSet_AddIfNotExistFuncLock() {
	var strSet gset.StrSet
	fmt.Println(strSet.AddIfNotExistFuncLock("str", func() bool {
		return true
	}))

	// Output:
	// true
}

func ExampleStrSet_Clear() {
	var strSet gset.StrSet
	strSet.Add([]string{"str1", "str2", "str3"}...)

	strSet.Clear()

	fmt.Println(strSet.Size())

	// Output:
	// 0
}

func ExampleStrSet_Complement() {
	strSet := gset.NewStrSet(true)
	strSet.Add([]string{"str1", "str2", "str3", "str4", "str5"}...)

	var s gset.StrSet
	s.Add([]string{"str1", "str2", "str3"}...)

	fmt.Println(s.Complement(strSet).Slice())

	// May Output:
	// [str4 str5]
}

func ExampleStrSet_Contains() {
	var set gset.StrSet
	set.Add("a")
	fmt.Println(set.Contains("a"))
	fmt.Println(set.Contains("A"))

	// Output:
	// true
	// false
}

func ExampleStrSet_ContainsI() {
	var set gset.StrSet
	set.Add("a")
	fmt.Println(set.ContainsI("a"))
	fmt.Println(set.ContainsI("A"))

	// Output:
	// true
	// true
}

func ExampleStrSet_Diff() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c"}...)
	var s2 gset.StrSet
	s2.Add([]string{"a", "b", "c", "d"}...)
	// 差集
	fmt.Println(s2.Diff(s1).Slice())

	// Output:
	// [d]
}

func ExampleStrSet_Equal() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c"}...)
	var s2 gset.StrSet
	s2.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s2.Equal(s1))

	// Output:
	// false
}

func ExampleStrSet_Intersect() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c"}...)
	var s2 gset.StrSet
	s2.Add([]string{"a", "b", "c", "d"}...)
	// 交集
	fmt.Println(s2.Intersect(s1).Slice())

	// May Output:
	// [c a b]
}

func ExampleStrSet_IsSubsetOf() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	var s2 gset.StrSet
	s2.Add([]string{"a", "b", "d"}...)
	fmt.Println(s2.IsSubsetOf(s1))

	// Output:
	// true
}

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

func ExampleStrSet_Join() {
	s1 := gset.NewStrSet(true)
	s1.Add([]string{"a", "b", "c", "d"}...)
	fmt.Println(s1.Join(","))

	// May Output:
	// b,c,d,a
}

func ExampleStrSet_LockFunc() {

}

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

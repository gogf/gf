// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glist_test

import (
	"container/list"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/glist"
)

func Example_basic() {
	n := 10
	l := glist.New()
	for i := 0; i < n; i++ {
		l.PushBack(i)
	}
	fmt.Println(l.Len())
	fmt.Println(l.FrontAll())
	fmt.Println(l.BackAll())
	for i := 0; i < n; i++ {
		fmt.Print(l.PopFront())
	}
	l.Clear()
	fmt.Println()
	fmt.Println(l.Len())

	// Output:
	//10
	//[0 1 2 3 4 5 6 7 8 9]
	//[9 8 7 6 5 4 3 2 1 0]
	//0123456789
	//0
}

func Example_iterate() {
	// concurrent-safe list.
	l := glist.NewFrom(garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from head.
	l.RLockFunc(func(list *list.List) {
		length := list.Len()
		if length > 0 {
			for i, e := 0, list.Front(); i < length; i, e = i+1, e.Next() {
				fmt.Print(e.Value)
			}
		}
	})
	fmt.Println()
	// iterate reading from tail.
	l.RLockFunc(func(list *list.List) {
		length := list.Len()
		if length > 0 {
			for i, e := 0, list.Back(); i < length; i, e = i+1, e.Prev() {
				fmt.Print(e.Value)
			}
		}
	})

	fmt.Println()

	// iterate reading from head using IteratorAsc.
	l.IteratorAsc(func(e *glist.Element) bool {
		fmt.Print(e.Value)
		return true
	})
	fmt.Println()
	// iterate reading from tail using IteratorDesc.
	l.IteratorDesc(func(e *glist.Element) bool {
		fmt.Print(e.Value)
		return true
	})

	fmt.Println()

	// iterate writing from head.
	l.LockFunc(func(list *list.List) {
		length := list.Len()
		if length > 0 {
			for i, e := 0, list.Front(); i < length; i, e = i+1, e.Next() {
				if e.Value == 6 {
					e.Value = "M"
					break
				}
			}
		}
	})
	fmt.Println(l)

	//output:
	//12345678910
	//10987654321
	//12345678910
	//10987654321
	//[1,2,3,4,5,M,7,8,9,10]
}

func Example_popItem() {
	l := glist.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	fmt.Println(l.PopBack())
	fmt.Println(l.PopBacks(2))
	fmt.Println(l.PopFront())
	fmt.Println(l.PopFronts(2))

	// Output:
	// 9
	// [8 7]
	// 1
	// [2 3]
}

func Example_join() {
	var l glist.List
	l.PushBacks(g.Slice{"a", "b", "c", "d"})

	fmt.Println(l.Join(","))

	// Output:
	// a,b,c,d
}

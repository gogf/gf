// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glist_test

import (
	"container/list"
	"fmt"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/gogf/gf/v2/container/glist"
)

func ExampleNew() {
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
	// 10
	// [0 1 2 3 4 5 6 7 8 9]
	// [9 8 7 6 5 4 3 2 1 0]
	// 0123456789
	// 0
}

func ExampleList_RLockFunc() {
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
	// Output:
	// 12345678910
	// 10987654321
}

func ExampleList_IteratorAsc() {
	// concurrent-safe list.
	l := glist.NewFrom(garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from head using IteratorAsc.
	l.IteratorAsc(func(e *glist.Element) bool {
		fmt.Print(e.Value)
		return true
	})

	// Output:
	// 12345678910
}

func ExampleList_IteratorDesc() {
	// concurrent-safe list.
	l := glist.NewFrom(garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from tail using IteratorDesc.
	l.IteratorDesc(func(e *glist.Element) bool {
		fmt.Print(e.Value)
		return true
	})
	// Output:
	// 10987654321
}

func ExampleList_LockFunc() {
	// concurrent-safe list.
	l := glist.NewFrom(garray.NewArrayRange(1, 10, 1).Slice(), true)
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

	// Output:
	// [1,2,3,4,5,M,7,8,9,10]
}

func ExampleList_PopBack() {
	l := glist.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	fmt.Println(l.PopBack())

	// Output:
	// 9
}
func ExampleList_PopBacks() {
	l := glist.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	fmt.Println(l.PopBacks(2))

	// Output:
	// [9 8]
}

func ExampleList_PopFront() {
	l := glist.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	fmt.Println(l.PopFront())

	// Output:
	// 1
}

func ExampleList_PopFronts() {
	l := glist.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	fmt.Println(l.PopFronts(2))

	// Output:
	// [1 2]
}

func ExampleList_Join() {
	var l glist.List
	l.PushBacks(g.Slice{"a", "b", "c", "d"})

	fmt.Println(l.Join(","))

	// Output:
	// a,b,c,d
}

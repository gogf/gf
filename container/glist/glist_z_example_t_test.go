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
	"github.com/gogf/gf/v2/container/glist"
	"github.com/gogf/gf/v2/frame/g"
)

func ExampleNewT() {
	n := 10
	l := glist.NewT[any]()
	for i := 0; i < n; i++ {
		l.PushBack(i)
	}

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.FrontAll())
	fmt.Println(l.BackAll())

	for i := 0; i < n; i++ {
		fmt.Print(l.PopFront())
	}

	fmt.Println()
	fmt.Println(l.Len())

	// Output:
	// 10
	// [0,1,2,3,4,5,6,7,8,9]
	// [0 1 2 3 4 5 6 7 8 9]
	// [9 8 7 6 5 4 3 2 1 0]
	// 0123456789
	// 0
}

func ExampleNewTFrom() {
	n := 10
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 10, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.FrontAll())
	fmt.Println(l.BackAll())

	for i := 0; i < n; i++ {
		fmt.Print(l.PopFront())
	}

	fmt.Println()
	fmt.Println(l.Len())

	// Output:
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
	// [1 2 3 4 5 6 7 8 9 10]
	// [10 9 8 7 6 5 4 3 2 1]
	// 12345678910
	// 0
}

func ExampleTList_PushFront() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushFront(0)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [0,1,2,3,4,5]
}

func ExampleTList_PushBack() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushBack(6)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [1,2,3,4,5,6]
}

func ExampleTList_PushFronts() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushFronts(g.Slice{0, -1, -2, -3, -4})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 10
	// [-4,-3,-2,-1,0,1,2,3,4,5]
}

func ExampleTList_PushBacks() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushBacks(g.Slice{6, 7, 8, 9, 10})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
}

func ExampleTList_PopBack() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopBack())
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// 4
	// [1,2,3,4]
}

func ExampleTList_PopFront() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopFront())
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 1
	// 4
	// [2,3,4,5]
}

func ExampleTList_PopBacks() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopBacks(2))
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [5 4]
	// 3
	// [1,2,3]
}

func ExampleTList_PopFronts() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopFronts(2))
	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1 2]
	// 3
	// [3,4,5]
}

func ExampleTList_PopBackAll() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopBackAll())
	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// [5 4 3 2 1]
	// 0
}

func ExampleTList_PopFrontAll() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)
	fmt.Println(l.PopFrontAll())
	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1 2 3 4 5]
	// 0
}

func ExampleTList_FrontAll() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.FrontAll())

	// Output:
	// [1,2,3,4,5]
	// [1 2 3 4 5]
}

func ExampleTList_BackAll() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.BackAll())

	// Output:
	// [1,2,3,4,5]
	// [5 4 3 2 1]
}

func ExampleTList_FrontValue() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.FrontValue())

	// Output:
	// [1,2,3,4,5]
	// 1
}

func ExampleTList_BackValue() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.BackValue())

	// Output:
	// [1,2,3,4,5]
	// 5
}

func ExampleTList_Front() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Front().Value)
	fmt.Println(l)

	e := l.Front()
	l.InsertBefore(e, 0)
	l.InsertAfter(e, "a")

	fmt.Println(l)

	// Output:
	// 1
	// [1,2,3,4,5]
	// [0,1,a,2,3,4,5]
}

func ExampleTList_Back() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Back().Value)
	fmt.Println(l)

	e := l.Back()
	l.InsertBefore(e, "a")
	l.InsertAfter(e, 6)

	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1,2,3,4,a,5,6]
}

func ExampleTList_Len() {
	l := glist.NewTFrom[any](g.Slice{1, 2, 3, 4, 5})

	fmt.Println(l.Len())

	// Output:
	// 5
}

func ExampleTList_Size() {
	l := glist.NewTFrom[any](g.Slice{1, 2, 3, 4, 5})

	fmt.Println(l.Size())

	// Output:
	// 5
}

func ExampleTList_MoveBefore() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	e := l.PushBack(6)
	fmt.Println(l.Size())
	fmt.Println(l)

	l.MoveBefore(e, l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e = &glist.TElement[any]{Value: 7}
	l.MoveBefore(e, l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [1,2,3,4,5,6]
	// 6
	// [6,1,2,3,4,5]
	// 6
	// [6,1,2,3,4,5]
}

func ExampleTList_MoveAfter() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	e := l.PushFront(0)
	fmt.Println(l.Size())
	fmt.Println(l)

	l.MoveAfter(e, l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e = &glist.TElement[any]{Value: -1}
	l.MoveAfter(e, l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 6
	// [0,1,2,3,4,5]
	// 6
	// [1,2,3,4,5,0]
	// 6
	// [1,2,3,4,5,0]
}

func ExampleTList_MoveToFront() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	l.MoveToFront(l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e := &glist.TElement[any]{Value: 6}
	l.MoveToFront(e)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [5,1,2,3,4]
	// 5
	// [5,1,2,3,4]
}

func ExampleTList_MoveToBack() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	l.MoveToBack(l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e := &glist.TElement[any]{Value: 0}
	l.MoveToBack(e)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [2,3,4,5,1]
	// 5
	// [2,3,4,5,1]
}

func ExampleTList_PushBackList() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	other := glist.NewTFrom[any](g.Slice{6, 7, 8, 9, 10})

	fmt.Println(other.Size())
	fmt.Println(other)

	l.PushBackList(other)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [6,7,8,9,10]
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
}

func ExampleTList_PushFrontList() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	other := glist.NewTFrom[any](g.Slice{-4, -3, -2, -1, 0})

	fmt.Println(other.Size())
	fmt.Println(other)

	l.PushFrontList(other)

	fmt.Println(l.Size())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 5
	// [-4,-3,-2,-1,0]
	// 10
	// [-4,-3,-2,-1,0,1,2,3,4,5]
}

func ExampleTList_InsertAfter() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.InsertAfter(l.Front(), "a")
	l.InsertAfter(l.Back(), "b")

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 7
	// [1,a,2,3,4,5,b]
}

func ExampleTList_InsertBefore() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.InsertBefore(l.Front(), "a")
	l.InsertBefore(l.Back(), "b")

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 7
	// [a,1,2,3,4,b,5]
}

func ExampleTList_Remove() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	fmt.Println(l.Remove(l.Front()))
	fmt.Println(l.Remove(l.Back()))

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 1
	// 5
	// 3
	// [2,3,4]
}

func ExampleTList_Removes() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.Removes([]*glist.TElement[any]{l.Front(), l.Back()})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 3
	// [2,3,4]
}

func ExampleTList_RemoveAll() {
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.RemoveAll()

	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// 0
}

func ExampleTList_RLockFunc() {
	// concurrent-safe list.
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 10, 1).Slice(), true)
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

func ExampleTList_IteratorAsc() {
	// concurrent-safe list.
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from head using IteratorAsc.
	l.IteratorAsc(func(e *glist.TElement[any]) bool {
		fmt.Print(e.Value)
		return true
	})

	// Output:
	// 12345678910
}

func ExampleTList_IteratorDesc() {
	// concurrent-safe list.
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from tail using IteratorDesc.
	l.IteratorDesc(func(e *glist.TElement[any]) bool {
		fmt.Print(e.Value)
		return true
	})
	// Output:
	// 10987654321
}

func ExampleTList_LockFunc() {
	// concurrent-safe list.
	l := glist.NewTFrom[any](garray.NewArrayRange(1, 10, 1).Slice(), true)
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

func ExampleTList_Join() {
	var l glist.TList[any]
	l.PushBacks(g.Slice{"a", "b", "c", "d"})

	fmt.Println(l.Join(","))

	// Output:
	// a,b,c,d
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glist

import (
	"fmt"

	"github.com/gogf/gf/contrib/generic_container/v2/garray"
)

func ExampleNew() {
	n := 10
	l := New[int]()
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
	fmt.Println(l.PopFront())
	fmt.Println(l.Len())

	// Output:
	// 10
	// [0,1,2,3,4,5,6,7,8,9]
	// [0 1 2 3 4 5 6 7 8 9]
	// [9 8 7 6 5 4 3 2 1 0]
	// 0123456789
	// 0
	// 0
	// 0
}

func ExampleNewFrom() {
	n := 10
	l := NewFrom[int](garray.NewArrayRange(1, 10, 1).Slice())

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

func ExampleList_PushFront() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PushBack() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PushFronts() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushFronts([]int{0, -1, -2, -3, -4})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 10
	// [-4,-3,-2,-1,0,1,2,3,4,5]
}

func ExampleList_PushBacks() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.PushBacks([]int{6, 7, 8, 9, 10})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 10
	// [1,2,3,4,5,6,7,8,9,10]
}

func ExampleList_PopBack() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PopFront() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PopBacks() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PopFronts() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PopBackAll() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_PopFrontAll() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_FrontAll() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.FrontAll())

	// Output:
	// [1,2,3,4,5]
	// [1 2 3 4 5]
}

func ExampleList_BackAll() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.BackAll())

	// Output:
	// [1,2,3,4,5]
	// [5 4 3 2 1]
}

func ExampleList_FrontValue() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.FrontValue())

	// Output:
	// [1,2,3,4,5]
	// 1
}

func ExampleList_BackValue() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l)
	fmt.Println(l.BackValue())

	// Output:
	// [1,2,3,4,5]
	// 5
}

func ExampleList_Front() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Front().Value)
	fmt.Println(l)

	e := l.Front()
	l.InsertBefore(e, 0)
	l.InsertAfter(e, 9)

	fmt.Println(l)

	// Output:
	// 1
	// [1,2,3,4,5]
	// [0,1,9,2,3,4,5]
}

func ExampleList_Back() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Back().Value)
	fmt.Println(l)

	e := l.Back()
	l.InsertBefore(e, 9)
	l.InsertAfter(e, 6)

	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// [1,2,3,4,9,5,6]
}

func ExampleList_Len() {
	l := NewFrom[int]([]int{1, 2, 3, 4, 5})

	fmt.Println(l.Len())

	// Output:
	// 5
}

func ExampleList_Size() {
	l := NewFrom[int]([]int{1, 2, 3, 4, 5})

	fmt.Println(l.Size())

	// Output:
	// 5
}

func ExampleList_MoveBefore() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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
	e = &Element[int]{Value: 7}
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

func ExampleList_MoveAfter() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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
	e = &Element[int]{Value: -1}
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

func ExampleList_MoveToFront() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	l.MoveToFront(l.Back())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e := &Element[int]{Value: 6}
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

func ExampleList_MoveToBack() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	// element of `l`
	l.MoveToBack(l.Front())

	fmt.Println(l.Size())
	fmt.Println(l)

	// not element of `l`
	e := &Element[int]{Value: 0}
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

func ExampleList_PushBackList() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	other := NewFrom[int]([]int{6, 7, 8, 9, 10})

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

func ExampleList_PushFrontList() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Size())
	fmt.Println(l)

	other := NewFrom[int]([]int{-4, -3, -2, -1, 0})

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

func ExampleList_InsertAfter() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.InsertAfter(l.Front(), 8)
	l.InsertAfter(l.Back(), 9)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 7
	// [1,8,2,3,4,5,9]
}

func ExampleList_InsertBefore() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.InsertBefore(l.Front(), 8)
	l.InsertBefore(l.Back(), 9)

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 7
	// [8,1,2,3,4,9,5]
}

func ExampleList_Remove() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

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

func ExampleList_Removes() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.Removes([]*Element[int]{l.Front(), l.Back()})

	fmt.Println(l.Len())
	fmt.Println(l)

	// Output:
	// 5
	// [1,2,3,4,5]
	// 3
	// [2,3,4]
}

func ExampleList_RemoveAll() {
	l := NewFrom[int](garray.NewArrayRange(1, 5, 1).Slice())

	fmt.Println(l.Len())
	fmt.Println(l)

	l.RemoveAll()

	fmt.Println(l.Len())

	// Output:
	// 5
	// [1,2,3,4,5]
	// 0
}

func ExampleList_IteratorAsc() {
	// concurrent-safe list.
	l := NewFrom[int](garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from head using IteratorAsc.
	l.IteratorAsc(func(e *Element[int]) bool {
		fmt.Print(e.Value)
		return true
	})

	// Output:
	// 12345678910
}

func ExampleList_IteratorDesc() {
	// concurrent-safe list.
	l := NewFrom[int](garray.NewArrayRange(1, 10, 1).Slice(), true)
	// iterate reading from tail using IteratorDesc.
	l.IteratorDesc(func(e *Element[int]) bool {
		fmt.Print(e.Value)
		return true
	})
	// Output:
	// 10987654321
}

func ExampleList_Join() {
	var l List[string]
	l.PushBacks([]string{"a", "b", "c", "d"})

	fmt.Println(l.Join(","))

	// Output:
	// a,b,c,d
}

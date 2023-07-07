// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/contrib/generic_container/v2/internal/empty"
	"github.com/gogf/gf/v2/util/gconv"
)

type exampleElement struct {
	code    int
	message string
}

func (e *exampleElement) Compare(other *exampleElement) int {
	if e == nil && other == nil {
		return 0
	}
	if e == nil && other != nil {
		return -1
	}
	if e != nil && other == nil {
		return 1
	}
	result := e.code - other.code
	if result != 0 {
		return result
	}
	result = strings.Compare(e.message, other.message)
	return result
}

func (e *exampleElement) String() string {
	if e == nil {
		return ""
	}
	if e.code == 0 && e.message == "" {
		return ""
	}
	result := ""
	if e.code != 0 {
		result += gconv.String(e.code)
	}
	if e.message != "" {
		result += "\"" + e.message + "\""
	}
	return result
}

func ExampleNew() {
	// A normal array.
	a := New[int]()

	// Adding items.
	for i := 0; i < 10; i++ {
		a.Append(i)
	}

	// Print the array length.
	fmt.Println(a.Len())

	// Print the array items.
	fmt.Println(a.Slice())

	// Retrieve item by index.
	fmt.Println(a.Get(6))

	// Check item existence.
	fmt.Println(a.Contains(6))
	fmt.Println(a.Contains(100))

	// Insert item before specified index.
	a.InsertAfter(9, 11)
	// Insert item after specified index.
	a.InsertBefore(10, 10)

	fmt.Println(a.Slice())

	// Modify item by index.
	a.Set(0, 100)
	fmt.Println(a.Slice())

	fmt.Println(a.At(0))

	// Search item and return its index.
	fmt.Println(a.Search(5))

	// Remove item by index.
	a.Remove(0)
	fmt.Println(a.Slice())

	// Empty the array, removes all items of it.
	fmt.Println(a.Slice())
	a.Clear()
	fmt.Println(a.Slice())

	// Output:
	// 10
	// [0 1 2 3 4 5 6 7 8 9]
	// 6 true
	// true
	// false
	// [0 1 2 3 4 5 6 7 8 9 10 11]
	// [100 1 2 3 4 5 6 7 8 9 10 11]
	// 100
	// 5
	// [1 2 3 4 5 6 7 8 9 10 11]
	// [1 2 3 4 5 6 7 8 9 10 11]
	// []
}

func ExampleArray_Iterator() {
	array := NewArrayFrom[string]([]string{"a", "b", "c"})
	// Iterator is alias of IteratorAsc, which iterates the array readonly in ascending order
	//  with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	array.Iterator(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})
	// IteratorDesc iterates the array readonly in descending order with given callback function `f`.
	// If `f` returns true, then it continues iterating; or false to stop.
	array.IteratorDesc(func(k int, v string) bool {
		fmt.Println(k, v)
		return true
	})

	// Output:
	// 0 a
	// 1 b
	// 2 c
	// 2 c
	// 1 b
	// 0 a
}

func ExampleArray_Reverse() {
	array := NewFrom[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Reverse makes array with elements in reverse order.
	fmt.Println(array.Reverse().Slice())

	// Output:
	// [9 8 7 6 5 4 3 2 1]
}

func ExampleArray_Shuffle() {
	array := NewFrom[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Shuffle randomly shuffles the array.
	fmt.Println(array.Shuffle().Slice())
}

func ExampleArray_Rands() {
	array := NewFrom[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Randomly retrieve and return 2 items from the array.
	// It does not delete the items from array.
	fmt.Println(array.Rands(2))

	// Randomly pick and return one item from the array.
	// It deletes the picked up item from array.
	fmt.Println(array.PopRand())
}

func ExampleArray_PopRand() {
	array := NewFrom[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Randomly retrieve and return 2 items from the array.
	// It does not delete the items from array.
	fmt.Println(array.Rands(2))

	// Randomly pick and return one item from the array.
	// It deletes the picked up item from array.
	fmt.Println(array.PopRand())
}

func ExampleArray_Join() {
	array := NewFrom[string]([]string{"a", "b", "c", "d"})
	fmt.Println(array.Join(","))

	// Output:
	// a,b,c,d
}

func ExampleArray_Chunk() {
	array := NewFrom[int]([]int{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Chunk splits an array into multiple arrays,
	// the size of each array is determined by `size`.
	// The last chunk may contain less than size elements.
	fmt.Println(array.Chunk(2))

	// Output:
	// [[1 2] [3 4] [5 6] [7 8] [9]]
}

func ExampleArray_PopLeft() {
	array := NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Any Pop* functions pick, delete and return the item from array.

	fmt.Println(array.PopLeft())
	fmt.Println(array.PopLefts(2))
	fmt.Println(array.PopRight())
	fmt.Println(array.PopRights(2))

	// Output:
	// 1 true
	// [2 3]
	// 9 true
	// [7 8]
}

func ExampleArray_PopLefts() {
	array := NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Any Pop* functions pick, delete and return the item from array.

	fmt.Println(array.PopLeft())
	fmt.Println(array.PopLefts(2))
	fmt.Println(array.PopRight())
	fmt.Println(array.PopRights(2))

	// Output:
	// 1 true
	// [2 3]
	// 9 true
	// [7 8]
}

func ExampleArray_PopRight() {
	array := NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Any Pop* functions pick, delete and return the item from array.

	fmt.Println(array.PopLeft())
	fmt.Println(array.PopLefts(2))
	fmt.Println(array.PopRight())
	fmt.Println(array.PopRights(2))

	// Output:
	// 1 true
	// [2 3]
	// 9 true
	// [7 8]
}

func ExampleArray_PopRights() {
	array := NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Any Pop* functions pick, delete and return the item from array.

	fmt.Println(array.PopLeft())
	fmt.Println(array.PopLefts(2))
	fmt.Println(array.PopRight())
	fmt.Println(array.PopRights(2))

	// Output:
	// 1 true
	// [2 3]
	// 9 true
	// [7 8]
}

func ExampleArray_Contains() {
	var array StdArray[string]
	array.Append("a")
	fmt.Println(array.Contains("a"))
	fmt.Println(array.Contains("A"))
	fmt.Println(array.ContainsI("A"))

	// Output:
	// true
	// false
	// true
}

func ExampleArray_Merge() {
	array1 := NewFrom[int]([]int{1, 2})
	array2 := NewFrom[int]([]int{3, 4})
	slice1 := []int{5, 6}
	slice2 := []int{7, 8}
	slice3 := []int{9, 0}
	fmt.Println(array1.Slice())
	array1.Merge(array1)
	array1.Merge(array2)
	array1.MergeSlice(slice1)
	array1.MergeSlice(slice2)
	array1.MergeSlice(slice3)
	fmt.Println(array1.Slice())

	// Output:
	// [1 2]
	// [1 2 1 2 3 4 5 6 7 8 9 0]
}

func ExampleArray_Filter() {
	array1 := NewFrom[*exampleElement]([]*exampleElement{
		{code: 0},
		{code: 1},
		{code: 2},
		nil,
		{message: "john"},
	})
	array2 := NewFrom[*exampleElement]([]*exampleElement{
		{code: 0},
		{code: 1},
		{code: 2},
		nil,
		{message: "john"},
	})
	fmt.Println(array1.Filter(func(index int, value *exampleElement) bool {
		return empty.IsNil(value)
	}).Slice())
	fmt.Println(array2.Filter(func(index int, value *exampleElement) bool {
		return empty.IsEmpty(value)
	}).Slice())

	// Output:
	// [ 1 2 "john"]
	// [1 2 "john"]
}

func ExampleArray_FilterEmpty() {
	array1 := NewFrom[*exampleElement]([]*exampleElement{
		{code: 0},
		{code: 1},
		{code: 2},
		nil,
		{message: "john"},
	})
	array2 := NewFrom[*exampleElement]([]*exampleElement{
		{code: 0},
		{code: 1},
		{code: 2},
		nil,
		{message: "john"},
	})
	fmt.Printf("%v\n", array1.FilterNil().Slice())
	fmt.Printf("%v\n", array2.FilterEmpty().Slice())

	// Output:
	// [ 1 2 "john"]
	// [1 2 "john"]
}

func ExampleArray_FilterNil() {
	array1 := NewFrom[*exampleElement]([]*exampleElement{
		{code: 0},
		{code: 1},
		{code: 2},
		nil,
		{message: "john"},
	})
	array2 := NewFrom[*exampleElement]([]*exampleElement{
		{code: 0},
		{code: 1},
		{code: 2},
		nil,
		{message: "john"},
	})
	fmt.Println(array1.FilterNil().Slice())
	fmt.Println(array2.FilterEmpty().Slice())

	// Output:
	// [ 1 2 "john"]
	// [1 2 "john"]
}

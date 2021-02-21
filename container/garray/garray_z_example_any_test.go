// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package garray_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/garray"
)

func ExampleNew() {
	// A normal array.
	a := garray.New()

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
	// 5
	// [1 2 3 4 5 6 7 8 9 10 11]
	// [1 2 3 4 5 6 7 8 9 10 11]
	// []
}

func ExampleArray_Iterator() {
	array := garray.NewArrayFrom(g.Slice{"a", "b", "c"})
	// Iterator is alias of IteratorAsc, which iterates the array readonly in ascending order
	//  with given callback function <f>.
	// If <f> returns true, then it continues iterating; or false to stop.
	array.Iterator(func(k int, v interface{}) bool {
		fmt.Println(k, v)
		return true
	})
	// IteratorDesc iterates the array readonly in descending order with given callback function <f>.
	// If <f> returns true, then it continues iterating; or false to stop.
	array.IteratorDesc(func(k int, v interface{}) bool {
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
	array := garray.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Reverse makes array with elements in reverse order.
	fmt.Println(array.Reverse().Slice())

	// Output:
	// [9 8 7 6 5 4 3 2 1]
}

func ExampleArray_Shuffle() {
	array := garray.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Shuffle randomly shuffles the array.
	fmt.Println(array.Shuffle().Slice())
}

func ExampleArray_Rands() {
	array := garray.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Randomly retrieve and return 2 items from the array.
	// It does not delete the items from array.
	fmt.Println(array.Rands(2))

	// Randomly pick and return one item from the array.
	// It deletes the picked up item from array.
	fmt.Println(array.PopRand())
}

func ExampleArray_PopRand() {
	array := garray.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Randomly retrieve and return 2 items from the array.
	// It does not delete the items from array.
	fmt.Println(array.Rands(2))

	// Randomly pick and return one item from the array.
	// It deletes the picked up item from array.
	fmt.Println(array.PopRand())
}

func ExampleArray_Join() {
	array := garray.NewFrom(g.Slice{"a", "b", "c", "d"})
	fmt.Println(array.Join(","))

	// Output:
	// a,b,c,d
}

func ExampleArray_Chunk() {
	array := garray.NewFrom(g.Slice{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Chunk splits an array into multiple arrays,
	// the size of each array is determined by <size>.
	// The last chunk may contain less than size elements.
	fmt.Println(array.Chunk(2))

	// Output:
	// [[1 2] [3 4] [5 6] [7 8] [9]]
}

func ExampleArray_PopLeft() {
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

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
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

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
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

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
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

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
	var array garray.StrArray
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
	array1 := garray.NewFrom(g.Slice{1, 2})
	array2 := garray.NewFrom(g.Slice{3, 4})
	slice1 := g.Slice{5, 6}
	slice2 := []int{7, 8}
	slice3 := []string{"9", "0"}
	fmt.Println(array1.Slice())
	array1.Merge(array1)
	array1.Merge(array2)
	array1.Merge(slice1)
	array1.Merge(slice2)
	array1.Merge(slice3)
	fmt.Println(array1.Slice())

	// Output:
	// [1 2]
	// [1 2 1 2 3 4 5 6 7 8 9 0]
}

func ExampleArray_FilterEmpty() {
	array1 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	array2 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	fmt.Printf("%#v\n", array1.FilterNil().Slice())
	fmt.Printf("%#v\n", array2.FilterEmpty().Slice())

	// Output:
	// []interface {}{0, 1, 2, "", []interface {}{}, "john"}
	// []interface {}{1, 2, "john"}
}

func ExampleArray_FilterNil() {
	array1 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	array2 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	fmt.Printf("%#v\n", array1.FilterNil().Slice())
	fmt.Printf("%#v\n", array2.FilterEmpty().Slice())

	// Output:
	// []interface {}{0, 1, 2, "", []interface {}{}, "john"}
	// []interface {}{1, 2, "john"}
}

func ExampleArray_Immutable() {
	array1 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	array2 := array1.FilterEmpty().Immutable()
	fmt.Printf("%#v\n", array2.Slice())

	// Output:
	// []interface {}{1, 2, "john"}
}

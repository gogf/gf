// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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

func Example_Basic() {
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
	// 6
	// true
	// false
	// [0 1 2 3 4 5 6 7 8 9 10 11]
	// [100 1 2 3 4 5 6 7 8 9 10 11]
	// 5
	// [1 2 3 4 5 6 7 8 9 10 11]
	// [1 2 3 4 5 6 7 8 9 10 11]
	// []
}

func Example_Rand() {
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Randomly retrieve and return 2 items from the array.
	// It does not delete the items from array.
	fmt.Println(array.Rands(2))

	// Randomly pick and return one item from the array.
	// It deletes the picked up item from array.
	fmt.Println(array.PopRand())
}

func Example_PopItem() {
	array := garray.NewFrom([]interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9})

	// Any Pop* functions pick, delete and return the item from array.

	fmt.Println(array.PopLeft())
	fmt.Println(array.PopLefts(2))
	fmt.Println(array.PopRight())
	fmt.Println(array.PopRights(2))

	// Output:
	// 1
	// [2 3]
	// 9
	// [7 8]
}

func Example_MergeArray() {
	array1 := garray.NewFrom([]interface{}{1, 2})
	array2 := garray.NewFrom([]interface{}{3, 4})
	slice1 := []interface{}{5, 6}
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

func Example_Filter() {
	array1 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	array2 := garray.NewFrom(g.Slice{0, 1, 2, nil, "", g.Slice{}, "john"})
	fmt.Printf("%#v\n", array1.FilterNil().Slice())
	fmt.Printf("%#v\n", array2.FilterEmpty().Slice())

	// Output:
	// []interface {}{0, 1, 2, "", []interface {}{}, "john"}
	// []interface {}{1, 2, "john"}
}

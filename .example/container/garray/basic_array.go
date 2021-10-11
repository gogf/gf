package main

import (
	"fmt"

	"github.com/gogf/gf/v2/container/garray"
)

func main() {
	// Create a int array, which is concurrent-unsafe in default.
	a := garray.NewIntArray()

	// Appending items.
	for i := 0; i < 10; i++ {
		a.Append(i)
	}

	// Get the length of the array.
	fmt.Println(a.Len())

	// Get the slice of the array.
	fmt.Println(a.Slice())

	// Get the item of specified index.
	fmt.Println(a.Get(6))

	// Insert after/before specified index.
	a.InsertAfter(9, 11)
	a.InsertBefore(10, 10)
	fmt.Println(a.Slice())

	a.Set(0, 100)
	fmt.Println(a.Slice())

	// Searching the item and returning the index.
	fmt.Println(a.Search(5))

	// Remove item of specified index.
	a.Remove(0)
	fmt.Println(a.Slice())

	// Clearing the array.
	fmt.Println(a.Slice())
	a.Clear()
	fmt.Println(a.Slice())
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gring_test

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gring"
)

func ExampleNew() {
	// Non concurrent safety
	gring.New(10)

	// Concurrent safety
	gring.New(10, true)

	// Output:
}

func ExampleRing_Cap() {
	r1 := gring.New(10)
	for i := 0; i < 5; i++ {
		r1.Set(i).Next()
	}
	fmt.Println("Cap:", r1.Cap())

	r2 := gring.New(10, true)
	for i := 0; i < 10; i++ {
		r2.Set(i).Next()
	}
	fmt.Println("Cap:", r2.Cap())

	// Output:
	// Cap: 10
	// Cap: 10
}

func ExampleRing_Len() {
	r1 := gring.New(10)
	for i := 0; i < 5; i++ {
		r1.Set(i).Next()
	}
	fmt.Println("Cap:", r1.Len())

	r2 := gring.New(10, true)
	for i := 0; i < 10; i++ {
		r2.Set(i).Next()
	}
	fmt.Println("Cap:", r2.Len())

	// Output:
	// Cap: 5
	// Cap: 10
}

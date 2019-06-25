// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glist_test

import (
	"fmt"
	"github.com/gogf/gf/g/container/glist"
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

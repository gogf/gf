// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gchan_test

import (
	"fmt"
	"github.com/gogf/gf/g/container/gchan"
)

func Example_basic() {
	n := 10
	c := gchan.New(n)
	for i := 0; i < n; i++ {
		c.Push(i)
	}
	fmt.Println(c.Len(), c.Cap())
	for i := 0; i < n; i++ {
		fmt.Print(c.Pop())
	}
	c.Close()

	// Output:
	//10 10
	//0123456789
}

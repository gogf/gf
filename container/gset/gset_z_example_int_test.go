// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gset_test

import (
	"fmt"
	"github.com/gogf/gf/container/gset"
)

func ExampleIntSet_Contains() {
	var set gset.IntSet
	set.Add(1)
	fmt.Println(set.Contains(1))
	fmt.Println(set.Contains(2))

	// Output:
	// true
	// false
}

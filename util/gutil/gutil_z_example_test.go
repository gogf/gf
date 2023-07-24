// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil_test

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"
)

func ExampleSliceInsertBefore() {
	s1 := g.Slice{
		0, 1, 2, 3, 4,
	}
	s2 := gutil.SliceInsertBefore(s1, 1, 8, 9)
	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// [0 1 2 3 4]
	// [0 8 9 1 2 3 4]
}

func ExampleSliceInsertAfter() {
	s1 := g.Slice{
		0, 1, 2, 3, 4,
	}
	s2 := gutil.SliceInsertAfter(s1, 1, 8, 9)
	fmt.Println(s1)
	fmt.Println(s2)

	// Output:
	// [0 1 2 3 4]
	// [0 1 8 9 2 3 4]
}

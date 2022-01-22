// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gipv4_test

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gipv4"
)

func ExampleGetFreePort() {
	fmt.Println(gipv4.GetFreePort())

	// May Output:
	// 57429 <nil>
}

func ExampleGetFreePorts() {
	fmt.Println(gipv4.GetFreePorts(2))

	// Output:
	// [57743 57744] <nil>
}

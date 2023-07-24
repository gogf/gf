// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp_test

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gudp"
)

func ExampleGetFreePort() {
	fmt.Println(gudp.GetFreePort())

	// May Output:
	// 57429 <nil>
}

func ExampleGetFreePorts() {
	fmt.Println(gudp.GetFreePorts(2))

	// May Output:
	// [57743 57744] <nil>
}

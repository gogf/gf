// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package gfile_test

import (
	"fmt"

	"github.com/wangyougui/gf/v2/os/gfile"
)

func ExampleHome() {
	// user's home directory
	homePath, _ := gfile.Home()
	fmt.Println(homePath)

	// May Output:
	// C:\Users\hailaz
}

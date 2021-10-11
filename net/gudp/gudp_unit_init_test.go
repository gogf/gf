// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp_test

import (
	"github.com/gogf/gf/v2/container/garray"
)

var (
	ports = garray.NewIntArray(true)
)

func init() {
	for i := 9000; i <= 10000; i++ {
		ports.Append(i)
	}
}

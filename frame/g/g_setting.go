// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"github.com/gogf/gf/v2/internal/utils"
)

// SetDebug enables/disables the GoFrame internal logging manually.
// Note that this function is not concurrent safe, be aware of the DATA RACE,
// which means you should call this function in your boot but not the runtime.
func SetDebug(enabled bool) {
	utils.SetDebugEnabled(enabled)
}

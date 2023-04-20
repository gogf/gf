// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package allyes

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
)

const (
	EnvName = "GF_CLI_ALL_YES"
)

// Init initializes the package manually.
func Init() {
	if gcmd.GetOpt("y") != nil {
		genv.MustSet(EnvName, "1")
	}
}

// Check checks whether option allow all yes for command.
func Check() bool {
	return genv.Get(EnvName).String() == "1"
}

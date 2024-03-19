// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package allyes

import (
	"github.com/wangyougui/gf/v2/os/gcmd"
	"github.com/wangyougui/gf/v2/os/genv"
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

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

import (
	"github.com/gogf/gf/v2/internal/command"
)

const (
	// Debug key for checking if in debug mode.
	commandEnvKeyForDebugKey = "gf.debug"
)

// isDebugEnabled marks whether GoFrame debug mode is enabled.
var isDebugEnabled = false

func init() {
	// Debugging configured.
	value := command.GetOptWithEnv(commandEnvKeyForDebugKey)
	if value == "" || value == "0" || value == "false" {
		isDebugEnabled = false
	} else {
		isDebugEnabled = true
	}
}

// IsDebugEnabled checks and returns whether debug mode is enabled.
// The debug mode is enabled when command argument "gf.debug" or environment "GF_DEBUG" is passed.
func IsDebugEnabled() bool {
	return isDebugEnabled
}

// SetDebugEnabled enables/disables the internal debug info.
func SetDebugEnabled(enabled bool) {
	isDebugEnabled = enabled
}

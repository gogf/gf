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

	// StackFilterKeyForGoFrame is the stack filtering key for all GoFrame module paths.
	// Eg: .../pkg/mod/github.com/gogf/gf/v2@v2.0.0-20211011134327-54dd11f51122/debug/gdebug/gdebug_caller.go
	StackFilterKeyForGoFrame = "github.com/gogf/gf/v"
)

var (
	// isDebugEnabled marks whether GoFrame debug mode is enabled.
	isDebugEnabled = false
)

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

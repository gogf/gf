// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package consts defines constants that are shared all among packages of framework.
package consts

const (
	ConfigNodeNameDatabase        = "database"
	ConfigNodeNameLogger          = "logger"
	ConfigNodeNameRedis           = "redis"
	ConfigNodeNameViewer          = "viewer"
	ConfigNodeNameServer          = "server"     // General version configuration item name.
	ConfigNodeNameServerSecondary = "httpserver" // New version configuration item name support from v2.

	// StackFilterKeyForGoFrame is the stack filtering key for all GoFrame module paths.
	// Eg: .../pkg/mod/github.com/gogf/gf/v2@v2.0.0-20211011134327-54dd11f51122/debug/gdebug/gdebug_caller.go
	StackFilterKeyForGoFrame = "github.com/gogf/gf/"
)

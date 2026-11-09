// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// URL schemes used by the ghttp package.
const (
	schemeHttp  = "http"
	schemeHttps = "https"
)

// Keys of the router serve tree.
const (
	routerKeyFuzz = "*fuzz"
	routerKeyList = "*list"
)

// gprocMessageExit is the process communication message for graceful shutdown.
const gprocMessageExit = "exit"

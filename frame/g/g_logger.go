// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"github.com/gogf/gf/os/glog"
)

// SetLogLevel sets the logging level globally.
// Deprecated, use functions of package glog or g.Log() instead.
func SetLogLevel(level int) {
	glog.SetLevel(level)
}

// GetLogLevel returns the global logging level.
// Deprecated, use functions of package glog or g.Log() instead.
func GetLogLevel() int {
	return glog.GetLevel()
}

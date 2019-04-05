// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
    "github.com/gogf/gf/g/os/glog"
)

// SetDebug disables/enables debug level for logging globally.
func SetDebug(debug bool) {
    glog.SetDebug(debug)
}

// SetLogLevel sets the logging level globally.
func SetLogLevel(level int) {
    glog.SetLevel(level)
}

// GetLogLevel returns the global logging level.
func GetLogLevel() int {
    return glog.GetLevel()
}
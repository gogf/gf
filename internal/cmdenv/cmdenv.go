// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package cmdenv provides access to certain variable for both command options and environment.
package cmdenv

import (
	"github.com/gogf/gf/os/gcmd"
	"os"
	"strings"

	"github.com/gogf/gf/container/gvar"
)

// Get returns the command line argument of the specified <key>.
// If the argument does not exist, then it returns the environment variable with specified <key>.
// It returns the default value <def> if none of them exists.
//
// Fetching Rules:
// 1. Command line arguments are in lowercase format, eg: gf.<package name>.<variable name>;
// 2. Environment arguments are in uppercase format, eg: GF_<package name>_<variable name>；
func Get(key string, def ...interface{}) *gvar.Var {
	value := interface{}(nil)
	if len(def) > 0 {
		value = def[0]
	}
	cmdKey := strings.ToLower(strings.Replace(key, "_", ".", -1))
	if v := gcmd.GetOpt(cmdKey); v != "" {
		value = v
	} else {
		envKey := strings.ToUpper(strings.Replace(key, ".", "_", -1))
		if v := os.Getenv(envKey); v != "" {
			value = v
		}
	}
	return gvar.New(value)
}

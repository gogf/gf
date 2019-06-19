// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package genv provides operations for environment variables of system.
package genv

import "os"

// All returns a copy of strings representing the environment,
// in the form "key=value".
func All() []string {
	return os.Environ()
}

// Get returns the value of the environment variable named by the <key>.
// It returns given <def> if the variable does not exist in the environment.
func Get(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

// Set sets the value of the environment variable named by the <key>.
// It returns an error, if any.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Remove deletes a single environment variable.
func Remove(key string) error {
	return os.Unsetenv(key)
}

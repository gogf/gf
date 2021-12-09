// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package genv provides operations for environment variables of system.
package genv

import (
	"os"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/internal/utils"
)

// All returns a copy of strings representing the environment,
// in the form "key=value".
func All() []string {
	return os.Environ()
}

// Map returns a copy of strings representing the environment as a map.
func Map() map[string]string {
	m := make(map[string]string)
	i := 0
	for _, s := range os.Environ() {
		i = strings.IndexByte(s, '=')
		m[s[0:i]] = s[i+1:]
	}
	return m
}

// Get creates and returns a Var with the value of the environment variable
// named by the `key`. It uses the given `def` if the variable does not exist
// in the environment.
func Get(key string, def ...interface{}) *gvar.Var {
	v, ok := os.LookupEnv(key)
	if !ok {
		if len(def) > 0 {
			return gvar.New(def[0])
		}
		return nil
	}
	return gvar.New(v)
}

// Set sets the value of the environment variable named by the `key`.
// It returns an error, if any.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// SetMap sets the environment variables using map.
func SetMap(m map[string]string) error {
	for k, v := range m {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Contains checks whether the environment variable named `key` exists.
func Contains(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// Remove deletes one or more environment variables.
func Remove(key ...string) error {
	var err error
	for _, v := range key {
		err = os.Unsetenv(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetWithCmd returns the environment value specified `key`.
// If the environment value does not exist, then it retrieves and returns the value from command line options.
// It returns the default value `def` if none of them exists.
//
// Fetching Rules:
// 1. Environment arguments are in uppercase format, eg: GF_<package name>_<variable name>；
// 2. Command line arguments are in lowercase format, eg: gf.<package name>.<variable name>;
func GetWithCmd(key string, def ...interface{}) *gvar.Var {
	envKey := utils.FormatEnvKey(key)
	if v := os.Getenv(envKey); v != "" {
		return gvar.New(v)
	}
	cmdKey := utils.FormatCmdKey(key)
	if v := command.GetOpt(cmdKey); v != "" {
		return gvar.New(v)
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// Build builds a map to an environment variable slice.
func Build(m map[string]string) []string {
	array := make([]string, len(m))
	index := 0
	for k, v := range m {
		array[index] = k + "=" + v
		index++
	}
	return array
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package command provides console operations, like options/arguments reading.
package command

import (
	"os"
	"regexp"
	"strings"
)

var (
	defaultParsedArgs    = make([]string, 0)
	defaultParsedOptions = make(map[string]string)
	argumentRegex        = regexp.MustCompile(`^\-{1,2}([\w\?\.\-]+)(=){0,1}(.*)$`)
)

// Init does custom initialization.
func Init(args ...string) {
	if len(args) == 0 {
		if len(defaultParsedArgs) == 0 && len(defaultParsedOptions) == 0 {
			args = os.Args
		} else {
			return
		}
	} else {
		defaultParsedArgs = make([]string, 0)
		defaultParsedOptions = make(map[string]string)
	}
	// Parsing os.Args with default algorithm.
	defaultParsedArgs, defaultParsedOptions = ParseUsingDefaultAlgorithm(args...)
}

// ParseUsingDefaultAlgorithm parses arguments using default algorithm.
func ParseUsingDefaultAlgorithm(args ...string) (parsedArgs []string, parsedOptions map[string]string) {
	parsedArgs = make([]string, 0)
	parsedOptions = make(map[string]string)
	for i := 0; i < len(args); {
		array := argumentRegex.FindStringSubmatch(args[i])
		if len(array) > 2 {
			if array[2] == "=" {
				parsedOptions[array[1]] = array[3]
			} else if i < len(args)-1 {
				if len(args[i+1]) > 0 && args[i+1][0] == '-' {
					// Eg: gf gen -d -n 1
					parsedOptions[array[1]] = array[3]
				} else {
					// Eg: gf gen -n 2
					parsedOptions[array[1]] = args[i+1]
					i += 2
					continue
				}
			} else {
				// Eg: gf gen -h
				parsedOptions[array[1]] = array[3]
			}
		} else {
			parsedArgs = append(parsedArgs, args[i])
		}
		i++
	}
	return
}

// GetOpt returns the option value named `name`.
func GetOpt(name string, def ...string) string {
	Init()
	if v, ok := defaultParsedOptions[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetOptAll returns all parsed options.
func GetOptAll() map[string]string {
	Init()
	return defaultParsedOptions
}

// ContainsOpt checks whether option named `name` exist in the arguments.
func ContainsOpt(name string) bool {
	Init()
	_, ok := defaultParsedOptions[name]
	return ok
}

// GetArg returns the argument at `index`.
func GetArg(index int, def ...string) string {
	Init()
	if index < len(defaultParsedArgs) {
		return defaultParsedArgs[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetArgAll returns all parsed arguments.
func GetArgAll() []string {
	Init()
	return defaultParsedArgs
}

// GetOptWithEnv returns the command line argument of the specified `key`.
// If the argument does not exist, then it returns the environment variable with specified `key`.
// It returns the default value `def` if none of them exists.
//
// Fetching Rules:
// 1. Command line arguments are in lowercase format, eg: gf.package.variable;
// 2. Environment arguments are in uppercase format, eg: GF_PACKAGE_VARIABLE；
func GetOptWithEnv(key string, def ...string) string {
	cmdKey := strings.ToLower(strings.Replace(key, "_", ".", -1))
	if ContainsOpt(cmdKey) {
		return GetOpt(cmdKey)
	} else {
		envKey := strings.ToUpper(strings.Replace(key, ".", "_", -1))
		if r, ok := os.LookupEnv(envKey); ok {
			return r
		} else {
			if len(def) > 0 {
				return def[0]
			}
		}
	}
	return ""
}

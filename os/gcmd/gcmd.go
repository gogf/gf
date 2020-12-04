// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package gcmd provides console operations, like options/arguments reading and command running.
package gcmd

import (
	"os"

	"github.com/gogf/gf/container/gvar"

	"github.com/gogf/gf/text/gregex"
)

var (
	defaultParsedArgs     = make([]string, 0)
	defaultParsedOptions  = make(map[string]string)
	defaultCommandFuncMap = make(map[string]func())
)

// Custom initialization.
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
	for i := 0; i < len(args); {
		array, _ := gregex.MatchString(`^\-{1,2}([\w\?\.\-]+)(=){0,1}(.*)$`, args[i])
		if len(array) > 2 {
			if array[2] == "=" {
				defaultParsedOptions[array[1]] = array[3]
			} else if i < len(args)-1 {
				if args[i+1][0] == '-' {
					defaultParsedOptions[array[1]] = array[3]
				} else {
					defaultParsedOptions[array[1]] = args[i+1]
					i += 2
					continue
				}
			} else {
				defaultParsedArgs = append(defaultParsedArgs, args[i])
			}
		} else {
			defaultParsedArgs = append(defaultParsedArgs, args[i])
		}
		i++
	}
}

// GetOpt returns the option value named <name>.
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

// GetOptVar returns the option value named <name> as gvar.Var.
func GetOptVar(name string, def ...string) *gvar.Var {
	Init()
	return gvar.New(GetOpt(name, def...))
}

// GetOptAll returns all parsed options.
func GetOptAll() map[string]string {
	Init()
	return defaultParsedOptions
}

// ContainsOpt checks whether option named <name> exist in the arguments.
func ContainsOpt(name string, def ...string) bool {
	Init()
	_, ok := defaultParsedOptions[name]
	return ok
}

// GetArg returns the argument at <index>.
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

// GetArgVar returns the argument at <index> as gvar.Var.
func GetArgVar(index int, def ...string) *gvar.Var {
	Init()
	return gvar.New(GetArg(index, def...))
}

// GetArgAll returns all parsed arguments.
func GetArgAll() []string {
	Init()
	return defaultParsedArgs
}

// BuildOptions builds the options as string.
func BuildOptions(m map[string]string, prefix ...string) string {
	options := ""
	leadStr := "-"
	if len(prefix) > 0 {
		leadStr = prefix[0]
	}
	for k, v := range m {
		if len(options) > 0 {
			options += " "
		}
		options += leadStr + k
		if v != "" {
			options += "=" + v
		}
	}
	return options
}

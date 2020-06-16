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
func doInit() {
	if len(defaultParsedArgs) > 0 {
		return
	}
	// Parsing os.Args with default algorithm.
	// The option should use '=' to separate its name and value in default.
	for _, arg := range os.Args {
		array, _ := gregex.MatchString(`^\-{1,2}([\w\?\.\-]+)={0,1}(.*)$`, arg)
		if len(array) == 3 {
			defaultParsedOptions[array[1]] = array[2]
		} else {
			defaultParsedArgs = append(defaultParsedArgs, arg)
		}
	}
}

// GetOpt returns the option value named <name>.
func GetOpt(name string, def ...string) string {
	doInit()
	if v, ok := defaultParsedOptions[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetOptVar returns the option value named <name> as gvar.Var.
func GetOptVar(name string, def ...string) gvar.Var {
	doInit()
	return gvar.New(GetOpt(name, def...))
}

// GetOptAll returns all parsed options.
func GetOptAll() map[string]string {
	doInit()
	return defaultParsedOptions
}

// ContainsOpt checks whether option named <name> exist in the arguments.
func ContainsOpt(name string, def ...string) bool {
	doInit()
	_, ok := defaultParsedOptions[name]
	return ok
}

// GetArg returns the argument at <index>.
func GetArg(index int, def ...string) string {
	doInit()
	if index < len(defaultParsedArgs) {
		return defaultParsedArgs[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetArgVar returns the argument at <index> as gvar.Var.
func GetArgVar(index int, def ...string) gvar.Var {
	doInit()
	return gvar.New(GetArg(index, def...))
}

// GetArgAll returns all parsed arguments.
func GetArgAll() []string {
	doInit()
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

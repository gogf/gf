// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

// Package gcmd provides console operations, like options/arguments reading and command running.
package gcmd

import (
	"os"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	CtxKeyParser         gctx.StrKey = `CtxKeyParser`
	CtxKeyCommand        gctx.StrKey = `CtxKeyCommand`
	CtxKeyArgumentsIndex gctx.StrKey = `CtxKeyArgumentsIndex`
)

const (
	helpOptionName        = "help"
	helpOptionNameShort   = "h"
	maxLineChars          = 120
	tracingInstrumentName = "github.com/gogf/gf/v2/os/gcmd.Command"
	tagNameName           = "name"
	tagNameShort          = "short"
)

// Init does custom initialization.
func Init(args ...string) {
	command.Init(args...)
}

// GetOpt returns the option value named `name` as gvar.Var.
func GetOpt(name string, def ...string) *gvar.Var {
	if v := command.GetOpt(name, def...); v != "" {
		return gvar.New(v)
	}
	if command.ContainsOpt(name) {
		return gvar.New("")
	}
	return nil
}

// GetOptAll returns all parsed options.
func GetOptAll() map[string]string {
	return command.GetOptAll()
}

// GetArg returns the argument at `index` as gvar.Var.
func GetArg(index int, def ...string) *gvar.Var {
	if v := command.GetArg(index, def...); v != "" {
		return gvar.New(v)
	}
	return nil
}

// GetArgAll returns all parsed arguments.
func GetArgAll() []string {
	return command.GetArgAll()
}

// GetOptWithEnv returns the command line argument of the specified `key`.
// If the argument does not exist, then it returns the environment variable with specified `key`.
// It returns the default value `def` if none of them exists.
//
// Fetching Rules:
// 1. Command line arguments are in lowercase format, eg: gf.`package name`.<variable name>;
// 2. Environment arguments are in uppercase format, eg: GF_`package name`_<variable name>；
func GetOptWithEnv(key string, def ...interface{}) *gvar.Var {
	cmdKey := utils.FormatCmdKey(key)
	if command.ContainsOpt(cmdKey) {
		return gvar.New(GetOpt(cmdKey))
	} else {
		envKey := utils.FormatEnvKey(key)
		if r, ok := os.LookupEnv(envKey); ok {
			return gvar.New(r)
		} else {
			if len(def) > 0 {
				return gvar.New(def[0])
			}
		}
	}
	return nil
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

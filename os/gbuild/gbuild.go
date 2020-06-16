// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gbuild manages the build-in variables from "gf build".
package gbuild

import (
	"encoding/json"
	"github.com/gogf/gf"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/util/gconv"
	"runtime"
)

var (
	builtInVarStr = ""                       // Raw variable base64 string.
	builtInVarMap = map[string]interface{}{} // Binary custom variable map decoded.
)

func init() {
	if builtInVarStr != "" {
		err := json.Unmarshal(gbase64.MustDecodeString(builtInVarStr), &builtInVarMap)
		if err != nil {
			intlog.Error(err)
		}
		builtInVarMap["gfVersion"] = gf.VERSION
		builtInVarMap["goVersion"] = runtime.Version()
		intlog.Printf("build variables: %+v", builtInVarMap)
	} else {
		intlog.Print("no build variables")
	}
}

// Info returns the basic built information of the binary as map.
// Note that it should be used with gf-cli tool "gf build",
// which injects necessary information into the binary.
func Info() map[string]string {
	return map[string]string{
		"gf":   GetString("gfVersion"),
		"go":   GetString("goVersion"),
		"git":  GetString("builtGit"),
		"time": GetString("builtTime"),
	}
}

// Get retrieves and returns the build-in binary variable with given name.
func Get(name string, def ...interface{}) interface{} {
	if v, ok := builtInVarMap[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// Get retrieves and returns the build-in binary variable of given name as gvar.Var.
func GetVar(name string, def ...interface{}) gvar.Var {
	return gvar.New(Get(name, def...))
}

// GetString retrieves and returns the build-in binary variable of given name as string.
func GetString(name string, def ...interface{}) string {
	return gconv.String(Get(name, def...))
}

// Map returns the custom build-in variable map.
func Map() map[string]interface{} {
	return builtInVarMap
}

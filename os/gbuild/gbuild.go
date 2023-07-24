// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gbuild manages the build-in variables from "gf build".
package gbuild

import (
	"context"
	"runtime"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
)

// BuildInfo maintains the built info of current binary.
type BuildInfo struct {
	GoFrame string                 // Built used GoFrame version.
	Golang  string                 // Built used Golang version.
	Git     string                 // Built used git repo. commit id and datetime.
	Time    string                 // Built datetime.
	Version string                 // Built version.
	Data    map[string]interface{} // All custom built data key-value pairs.
}

const (
	gfVersion    = `gfVersion`
	goVersion    = `goVersion`
	BuiltGit     = `builtGit`
	BuiltTime    = `builtTime`
	BuiltVersion = `builtVersion`
)

var (
	builtInVarStr = ""                       // Raw variable base64 string, which is injected by go build flags.
	builtInVarMap = map[string]interface{}{} // Binary custom variable map decoded.
)

func init() {
	// The `builtInVarStr` is injected by go build flags.
	if builtInVarStr != "" {
		err := json.UnmarshalUseNumber(gbase64.MustDecodeString(builtInVarStr), &builtInVarMap)
		if err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
		builtInVarMap[gfVersion] = gf.VERSION
		builtInVarMap[goVersion] = runtime.Version()
		intlog.Printf(context.TODO(), "build variables: %+v", builtInVarMap)
	} else {
		intlog.Print(context.TODO(), "no build variables")
	}
}

// Info returns the basic built information of the binary as map.
// Note that it should be used with gf-cli tool "gf build",
// which automatically injects necessary information into the binary.
func Info() BuildInfo {
	return BuildInfo{
		GoFrame: Get(gfVersion).String(),
		Golang:  Get(goVersion).String(),
		Git:     Get(BuiltGit).String(),
		Time:    Get(BuiltTime).String(),
		Version: Get(BuiltVersion).String(),
		Data:    Data(),
	}
}

// Get retrieves and returns the build-in binary variable with given name.
func Get(name string, def ...interface{}) *gvar.Var {
	if v, ok := builtInVarMap[name]; ok {
		return gvar.New(v)
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// Data returns the custom build-in variables as map.
func Data() map[string]interface{} {
	return builtInVarMap
}

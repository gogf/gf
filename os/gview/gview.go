// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gview implements a template engine based on text/template.
//
// Reserved template variable names:
// I18nLanguage: Assign this variable to define i18n language for each page.
package gview

import (
	"context"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
)

// View object for template engine.
type View struct {
	searchPaths  *garray.StrArray       // Searching array for path, NOT concurrent-safe for performance purpose.
	data         map[string]interface{} // Global template variables.
	funcMap      map[string]interface{} // Global template function map.
	fileCacheMap *gmap.StrAnyMap        // File cache map.
	config       Config                 // Extra configuration for the view.
}

type (
	Params  = map[string]interface{} // Params is type for template params.
	FuncMap = map[string]interface{} // FuncMap is type for custom template functions.
)

const (
	commandEnvKeyForPath = "gf.gview.path"
)

var (
	// Default view object.
	defaultViewObj *View
)

// checkAndInitDefaultView checks and initializes the default view object.
// The default view object will be initialized just once.
func checkAndInitDefaultView() {
	if defaultViewObj == nil {
		defaultViewObj = New()
	}
}

// ParseContent parses the template content directly using the default view object
// and returns the parsed content.
func ParseContent(ctx context.Context, content string, params ...Params) (string, error) {
	checkAndInitDefaultView()
	return defaultViewObj.ParseContent(ctx, content, params...)
}

// New returns a new view object.
// The parameter `path` specifies the template directory path to load template files.
func New(path ...string) *View {
	var (
		ctx = context.TODO()
	)
	view := &View{
		searchPaths:  garray.NewStrArray(),
		data:         make(map[string]interface{}),
		funcMap:      make(map[string]interface{}),
		fileCacheMap: gmap.NewStrAnyMap(true),
		config:       DefaultConfig(),
	}
	if len(path) > 0 && len(path[0]) > 0 {
		if err := view.SetPath(path[0]); err != nil {
			intlog.Errorf(context.TODO(), `%+v`, err)
		}
	} else {
		// Customized dir path from env/cmd.
		if envPath := gcmd.GetOptWithEnv(commandEnvKeyForPath).String(); envPath != "" {
			if gfile.Exists(envPath) {
				if err := view.SetPath(envPath); err != nil {
					intlog.Errorf(context.TODO(), `%+v`, err)
				}
			} else {
				if errorPrint() {
					glog.Errorf(ctx, "Template directory path does not exist: %s", envPath)
				}
			}
		} else {
			// Dir path of working dir.
			if pwdPath := gfile.Pwd(); pwdPath != "" {
				if err := view.SetPath(pwdPath); err != nil {
					intlog.Errorf(context.TODO(), `%+v`, err)
				}
			}
			// Dir path of binary.
			if selfPath := gfile.SelfDir(); selfPath != "" && gfile.Exists(selfPath) {
				if err := view.AddPath(selfPath); err != nil {
					intlog.Errorf(context.TODO(), `%+v`, err)
				}
			}
			// Dir path of main package.
			if mainPath := gfile.MainPkgPath(); mainPath != "" && gfile.Exists(mainPath) {
				if err := view.AddPath(mainPath); err != nil {
					intlog.Errorf(context.TODO(), `%+v`, err)
				}
			}
		}
	}
	view.SetDelimiters("{{", "}}")
	// default build-in variables.
	view.data["GF"] = map[string]interface{}{
		"version": gf.VERSION,
	}
	// default build-in functions.
	view.BindFuncMap(FuncMap{
		"eq":         view.buildInFuncEq,
		"ne":         view.buildInFuncNe,
		"lt":         view.buildInFuncLt,
		"le":         view.buildInFuncLe,
		"gt":         view.buildInFuncGt,
		"ge":         view.buildInFuncGe,
		"text":       view.buildInFuncText,
		"html":       view.buildInFuncHtmlEncode,
		"htmlencode": view.buildInFuncHtmlEncode,
		"htmldecode": view.buildInFuncHtmlDecode,
		"encode":     view.buildInFuncHtmlEncode,
		"decode":     view.buildInFuncHtmlDecode,
		"url":        view.buildInFuncUrlEncode,
		"urlencode":  view.buildInFuncUrlEncode,
		"urldecode":  view.buildInFuncUrlDecode,
		"date":       view.buildInFuncDate,
		"substr":     view.buildInFuncSubStr,
		"strlimit":   view.buildInFuncStrLimit,
		"concat":     view.buildInFuncConcat,
		"replace":    view.buildInFuncReplace,
		"compare":    view.buildInFuncCompare,
		"hidestr":    view.buildInFuncHideStr,
		"highlight":  view.buildInFuncHighlight,
		"toupper":    view.buildInFuncToUpper,
		"tolower":    view.buildInFuncToLower,
		"nl2br":      view.buildInFuncNl2Br,
		"include":    view.buildInFuncInclude,
		"dump":       view.buildInFuncDump,
		"map":        view.buildInFuncMap,
		"maps":       view.buildInFuncMaps,
		"json":       view.buildInFuncJson,
		"xml":        view.buildInFuncXml,
		"ini":        view.buildInFuncIni,
		"yaml":       view.buildInFuncYaml,
		"yamli":      view.buildInFuncYamlIndent,
		"toml":       view.buildInFuncToml,
		"plus":       view.buildInFuncPlus,
		"minus":      view.buildInFuncMinus,
		"times":      view.buildInFuncTimes,
		"divide":     view.buildInFuncDivide,
	})
	return view
}

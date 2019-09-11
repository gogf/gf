// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gview implements a template engine based on text/template.
//
// Reserved template variable names:
//     I18nLanguage: Assign this variable to define i18n language for each page.
package gview

import (
	"errors"
	"fmt"
	"sync"

	"github.com/gogf/gf/i18n/gi18n"

	"github.com/gogf/gf/os/gres"

	"github.com/gogf/gf"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/internal/cmdenv"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gspath"
)

// View object for template engine.
type View struct {
	mu          sync.RWMutex
	paths       *garray.StrArray       // Searching path array.
	data        map[string]interface{} // Global template variables.
	funcMap     map[string]interface{} // Global template function map.
	i18nManager *gi18n.Manager         // I18n manager for this view.
	delimiters  []string               // Customized template delimiters.
}

// Params is type for template params.
type Params = map[string]interface{}

// FuncMap is type for custom template functions.
type FuncMap = map[string]interface{}

// Default view object.
var defaultViewObj *View

// checkAndInitDefaultView checks and initializes the default view object.
// The default view object will be initialized just once.
func checkAndInitDefaultView() {
	if defaultViewObj == nil {
		defaultViewObj = New()
	}
}

// ParseContent parses the template content directly using the default view object
// and returns the parsed content.
func ParseContent(content string, params Params) (string, error) {
	checkAndInitDefaultView()
	return defaultViewObj.ParseContent(content, params)
}

// New returns a new view object.
// The parameter <path> specifies the template directory path to load template files.
func New(path ...string) *View {
	view := &View{
		paths:       garray.NewStrArray(true),
		data:        make(map[string]interface{}),
		funcMap:     make(map[string]interface{}),
		i18nManager: gi18n.Instance(),
		delimiters:  make([]string, 2),
	}
	if len(path) > 0 && len(path[0]) > 0 {
		view.SetPath(path[0])
	} else {
		// Customized dir path from env/cmd.
		if envPath := cmdenv.Get("gf.gview.path").String(); envPath != "" {
			if gfile.Exists(envPath) {
				view.SetPath(envPath)
			} else {
				if errorPrint() {
					glog.Errorf("Template directory path does not exist: %s", envPath)
				}
			}
		} else {
			// Dir path of working dir.
			view.SetPath(gfile.Pwd())
			// Dir path of binary.
			if selfPath := gfile.SelfDir(); selfPath != "" && gfile.Exists(selfPath) {
				view.AddPath(selfPath)
			}
			// Dir path of main package.
			if mainPath := gfile.MainPkgPath(); mainPath != "" && gfile.Exists(mainPath) {
				view.AddPath(mainPath)
			}
		}
	}
	view.SetDelimiters("{{", "}}")
	// default build-in variables.
	view.data["GF"] = map[string]interface{}{
		"version": gf.VERSION,
	}
	// default build-in functions.
	view.BindFunc("eq", view.funcEq)
	view.BindFunc("ne", view.funcNe)
	view.BindFunc("lt", view.funcLt)
	view.BindFunc("le", view.funcLe)
	view.BindFunc("gt", view.funcGt)
	view.BindFunc("ge", view.funcGe)
	view.BindFunc("text", view.funcText)
	view.BindFunc("html", view.funcHtmlEncode)
	view.BindFunc("htmlencode", view.funcHtmlEncode)
	view.BindFunc("htmldecode", view.funcHtmlDecode)
	view.BindFunc("url", view.funcUrlEncode)
	view.BindFunc("urlencode", view.funcUrlEncode)
	view.BindFunc("urldecode", view.funcUrlDecode)
	view.BindFunc("date", view.funcDate)
	view.BindFunc("substr", view.funcSubStr)
	view.BindFunc("strlimit", view.funcStrLimit)
	view.BindFunc("compare", view.funcCompare)
	view.BindFunc("hidestr", view.funcHideStr)
	view.BindFunc("highlight", view.funcHighlight)
	view.BindFunc("toupper", view.funcToUpper)
	view.BindFunc("tolower", view.funcToLower)
	view.BindFunc("nl2br", view.funcNl2Br)
	view.BindFunc("include", view.funcInclude)
	return view
}

// SetPath sets the template directory path for template file search.
// The parameter <path> can be absolute or relative path, but absolute path is suggested.
func (view *View) SetPath(path string) error {
	isDir := false
	realPath := ""
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			view.paths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := gspath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		err := errors.New(fmt.Sprintf(`[gview] SetPath failed: path "%s" does not exist`, path))
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// Should be a directory.
	if !isDir {
		err := errors.New(fmt.Sprintf(`[gview] SetPath failed: path "%s" should be directory type`, path))
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// Repeated path check.
	if view.paths.Search(realPath) != -1 {
		return nil
	}
	view.paths.Clear()
	view.paths.Append(realPath)
	//glog.Debug("[gview] SetPath:", realPath)
	return nil
}

// AddPath adds a absolute or relative path to the search paths.
func (view *View) AddPath(path string) error {
	isDir := false
	realPath := ""
	if file := gres.Get(path); file != nil {
		realPath = path
		isDir = file.FileInfo().IsDir()
	} else {
		// Absolute path.
		realPath = gfile.RealPath(path)
		if realPath == "" {
			// Relative path.
			view.paths.RLockFunc(func(array []string) {
				for _, v := range array {
					if path, _ := gspath.Search(v, path); path != "" {
						realPath = path
						break
					}
				}
			})
		}
		if realPath != "" {
			isDir = gfile.IsDir(realPath)
		}
	}
	// Path not exist.
	if realPath == "" {
		err := errors.New(fmt.Sprintf(`[gview] AddPath failed: path "%s" does not exist`, path))
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// realPath should be type of folder.
	if !isDir {
		err := errors.New(fmt.Sprintf(`[gview] AddPath failed: path "%s" should be directory type`, path))
		if errorPrint() {
			glog.Error(err)
		}
		return err
	}
	// Repeated path check.
	if view.paths.Search(realPath) != -1 {
		return nil
	}
	view.paths.Append(realPath)
	return nil
}

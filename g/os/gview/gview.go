// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gview implements a template engine based on text/template.
package gview

import (
    "bytes"
    "errors"
    "fmt"
    "github.com/gogf/gf"
    "github.com/gogf/gf/g/container/garray"
    "github.com/gogf/gf/g/encoding/ghash"
    "github.com/gogf/gf/g/encoding/ghtml"
    "github.com/gogf/gf/g/encoding/gurl"
    "github.com/gogf/gf/g/internal/cmdenv"
    "github.com/gogf/gf/g/os/gfcache"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gspath"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/os/gview/internal/text/template"
    "github.com/gogf/gf/g/text/gstr"
    "github.com/gogf/gf/g/util/gconv"
    "strings"
    "sync"
)

type View struct {
    mu         sync.RWMutex
    paths      *garray.StringArray     // 模板查找目录(绝对路径)
    data       map[string]interface{}  // 模板变量
    funcmap    map[string]interface{}  // FuncMap
    delimiters []string                // 模板变量分隔符号
}

// Template params type.
type Params  = map[string]interface{}

// Customized template function map type.
type FuncMap = map[string]interface{}

// Default view object.
var defaultViewObj *View

// checkAndInitDefaultView checks and initializes the default view object.
// The default view object will be initialized just once.
func checkAndInitDefaultView() {
    if defaultViewObj == nil {
        defaultViewObj = New(gfile.Pwd())
    }
}

// ParseContent parses the template content directly using the default view object
// and returns the parsed content.
func ParseContent(content string, params Params) ([]byte, error) {
    checkAndInitDefaultView()
    return defaultViewObj.ParseContent(content, params)
}

// New returns a new view object.
// The parameter <path> specifies the template directory path to load template files.
func New(path...string) *View {
    view := &View {
        paths      : garray.NewStringArray(),
        data       : make(map[string]interface{}),
        funcmap    : make(map[string]interface{}),
        delimiters : make([]string, 2),
    }
    if len(path) > 0 && len(path[0]) > 0 {
        view.SetPath(path[0])
    } else {
        // Customized dir path from env/cmd.
        if envPath := cmdenv.Get("gf.gview.path").String(); envPath != "" {
            if gfile.Exists(envPath) {
	            view.SetPath(envPath)
            } else {
                glog.Errorfln("Template directory path does not exist: %s", envPath)
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
    view.data["GF"] = map[string]interface{} {
        "version" : gf.VERSION,
    }
    // default build-in functions.
    view.BindFunc("text",        view.funcText)
    view.BindFunc("html",        view.funcHtmlEncode)
    view.BindFunc("htmlencode",  view.funcHtmlEncode)
    view.BindFunc("htmldecode",  view.funcHtmlDecode)
    view.BindFunc("url",         view.funcUrlEncode)
    view.BindFunc("urlencode",   view.funcUrlEncode)
    view.BindFunc("urldecode",   view.funcUrlDecode)
    view.BindFunc("date",        view.funcDate)
    view.BindFunc("substr",      view.funcSubStr)
    view.BindFunc("strlimit",    view.funcStrLimit)
    view.BindFunc("compare",     view.funcCompare)
    view.BindFunc("hidestr",     view.funcHideStr)
    view.BindFunc("highlight",   view.funcHighlight)
    view.BindFunc("toupper",     view.funcToUpper)
    view.BindFunc("tolower",     view.funcToLower)
    view.BindFunc("nl2br",       view.funcNl2Br)
    view.BindFunc("include",     view.funcInclude)
    return view
}

// SetPath sets the template directory path for template file search.
// The param <path> can be absolute or relative path, but absolute path is suggested.
func (view *View) SetPath(path string) error {
	// Absolute path.
    realPath := gfile.RealPath(path)
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
	// Path not exist.
    if realPath == "" {
        buffer := bytes.NewBuffer(nil)
        if view.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gview] SetPath failed: cannot find directory \"%s\" in following paths:", path))
            view.paths.RLockFunc(func(array []string) {
                for k, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s",k + 1,  v))
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf(`[gview] SetPath failed: path "%s" does not exist`, path))
        }
        err := errors.New(buffer.String())
        glog.Error(err)
        return err
    }
	// Should be a directory.
    if !gfile.IsDir(realPath) {
        err := errors.New(fmt.Sprintf(`[gview] SetPath failed: path "%s" should be directory type`, path))
        glog.Error(err)
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
	// Absolute path.
    realPath := gfile.RealPath(path)
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
    // Path not exist.
    if realPath == "" {
        buffer := bytes.NewBuffer(nil)
        if view.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gview] AddPath failed: cannot find directory \"%s\" in following paths:", path))
            view.paths.RLockFunc(func(array []string) {
                for k, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s",k + 1,  v))
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf(`[gview] AddPath failed: path "%s" does not exist`, path))
        }
        err := errors.New(buffer.String())
        glog.Error(err)
        return err
    }
    // realPath should be type of folder.
    if !gfile.IsDir(realPath) {
        err := errors.New(fmt.Sprintf(`[gview] AddPath failed: path "%s" should be directory type`, path))
        glog.Error(err)
        return err
    }
	// Repeated path check.
    if view.paths.Search(realPath) != -1 {
        return nil
    }
    view.paths.Append(realPath)
    //glog.Debug("[gview] AddPath:", realPath)
    return nil
}

// Assign binds multiple template variables to current view object.
// Each goroutine will take effect after the call, so it is concurrent-safe.
func (view *View) Assigns(data Params) {
    view.mu.Lock()
    for k, v := range data {
        view.data[k] = v
    }
    view.mu.Unlock()
}

// Assign binds a template variable to current view object.
// Each goroutine will take effect after the call, so it is concurrent-safe.
func (view *View) Assign(key string, value interface{}) {
    view.mu.Lock()
    view.data[key] = value
    view.mu.Unlock()
}

// ParseContent parses given template file <file>
// with given template parameters <params> and function map <funcMap>
// and returns the parsed content in []byte.
func (view *View) Parse(file string, params Params, funcMap...map[string]interface{}) ([]byte, error) {
    path := ""
    view.paths.RLockFunc(func(array []string) {
        for _, v := range array {
            if path, _ = gspath.Search(v, file); path != "" {
                break
            }
        }
    })
    if path == "" {
        buffer := bytes.NewBuffer(nil)
        if view.paths.Len() > 0 {
            buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" in following paths:", file))
            view.paths.RLockFunc(func(array []string) {
                for k, v := range array {
                    buffer.WriteString(fmt.Sprintf("\n%d. %s",k + 1,  v))
                }
            })
        } else {
            buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" with no path set/add", file))
        }
        glog.Error(buffer.String())
        return nil, errors.New(fmt.Sprintf(`tpl "%s" not found`, file))
    }
    content := gfcache.GetContents(path)
    view.mu.RLock()
    defer view.mu.RUnlock()
    buffer := bytes.NewBuffer(nil)
    tplObj := template.New(path).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcmap)
    if len(funcMap) > 0 {
	    tplObj = tplObj.Funcs(funcMap[0])
    }
    if tpl, err := tplObj.Parse(content); err != nil {
        return nil, err
    } else {
	    // Note that the template variable assignment cannot change the value
	    // of the existing <params> or view.data because both variables are pointers.
	    // It's need to merge the values of the two maps into a new map.
        vars := (map[string]interface{})(nil)
        if len(view.data) > 0 {
            if len(params) > 0 {
                vars = make(map[string]interface{}, len(view.data) + len(params))
                for k, v := range params {
                    vars[k] = v
                }
                for k, v := range view.data {
                    vars[k] = v
                }
            } else {
                vars = view.data
            }
        } else {
            vars = params
        }
        if err := tpl.Execute(buffer, vars); err != nil {
            return nil, err
        }
    }
    return buffer.Bytes(), nil
}

// ParseContent parses given template content <content>
// with given template parameters <params> and function map <funcMap>
// and returns the parsed content in []byte.
func (view *View) ParseContent(content string, params Params, funcMap...map[string]interface{}) ([]byte, error) {
    view.mu.RLock()
    defer view.mu.RUnlock()
    name   := gconv.String(ghash.BKDRHash64([]byte(content)))
    buffer := bytes.NewBuffer(nil)
    tplObj := template.New(name).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcmap)
    if len(funcMap) > 0 {
	    tplObj = tplObj.Funcs(funcMap[0])
    }
    if tpl, err := tplObj.Parse(content); err != nil {
        return nil, err
    } else {
	    // Note that the template variable assignment cannot change the value
	    // of the existing <params> or view.data because both variables are pointers.
	    // It's need to merge the values of the two maps into a new map.
        vars := (map[string]interface{})(nil)
        if len(view.data) > 0 {
            if len(params) > 0 {
                vars = make(map[string]interface{}, len(view.data) + len(params))
                for k, v := range params {
                    vars[k] = v
                }
                for k, v := range view.data {
                    vars[k] = v
                }
            } else {
                vars = view.data
            }
        } else {
            vars = params
        }
        if err := tpl.Execute(buffer, vars); err != nil {
            return nil, err
        }
    }
    return buffer.Bytes(), nil
}

// SetDelimiters sets customized delimiters for template parsing.
func (view *View) SetDelimiters(left, right string) {
    view.delimiters[0] = left
    view.delimiters[1] = right
}

// BindFunc registers customized template function named <name>
// with given function <function> to current view object.
// The <name> is the function name which can be called in template content.
func (view *View) BindFunc(name string, function interface{}) {
    view.mu.Lock()
    view.funcmap[name] = function
    view.mu.Unlock()
}

// BindFuncMap registers customized template functions by map to current view object.
// The key of map is the template function name
// and the value of map is the address of customized function.
func (view *View) BindFuncMap(funcMap FuncMap) {
	view.mu.Lock()
	for k, v := range funcMap {
		view.funcmap[k] = v
	}
	view.mu.Unlock()
}

// Build-in template function: include
func (view *View) funcInclude(file string, data...map[string]interface{}) string {
    var m map[string]interface{} = nil
    if len(data) > 0 {
        m = data[0]
    }
    content, err := view.Parse(file, m)
    if err != nil {
        return err.Error()
    }
    return string(content)
}

// Build-in template function: text
func (view *View) funcText(html interface{}) string {
    return ghtml.StripTags(gconv.String(html))
}

// Build-in template function: html
func (view *View) funcHtmlEncode(html interface{}) string {
    return ghtml.Entities(gconv.String(html))
}

// Build-in template function: htmldecode
func (view *View) funcHtmlDecode(html interface{}) string {
    return ghtml.EntitiesDecode(gconv.String(html))
}

// Build-in template function: url
func (view *View) funcUrlEncode(url interface{}) string {
    return gurl.Encode(gconv.String(url))
}

// Build-in template function: urldecode
func (view *View) funcUrlDecode(url interface{}) string {
    if content, err := gurl.Decode(gconv.String(url)); err == nil {
        return content
    } else {
        return err.Error()
    }
}

// Build-in template function: date
func (view *View) funcDate(format string, timestamp...interface{}) string {
    t := int64(0)
    if len(timestamp) > 0 {
        t = gconv.Int64(timestamp[0])
    }
    if t == 0 {
        t = gtime.Millisecond()
    }
    return gtime.NewFromTimeStamp(t).Format(format)
}

// Build-in template function: compare
func (view *View) funcCompare(value1, value2 interface{}) int {
    return strings.Compare(gconv.String(value1), gconv.String(value2))
}

// Build-in template function: substr
func (view *View) funcSubStr(start, end int, str interface{}) string {
    return gstr.SubStr(gconv.String(str), start, end)
}

// Build-in template function: strlimit
func (view *View) funcStrLimit(length int, suffix string, str interface{}) string {
    return gstr.StrLimit(gconv.String(str), length, suffix)
}

// Build-in template function: highlight
func (view *View) funcHighlight(key string, color string, str interface{}) string {
    return gstr.Replace(gconv.String(str), key, fmt.Sprintf(`<span style="color:%s;">%s</span>`, color, key))
}

// Build-in template function: hidestr
func (view *View) funcHideStr(percent int, hide string, str interface{}) string {
    return gstr.HideStr(gconv.String(str), percent, hide)
}

// Build-in template function: toupper
func (view *View) funcToUpper(str interface{}) string {
    return gstr.ToUpper(gconv.String(str))
}

// Build-in template function: toupper
func (view *View) funcToLower(str interface{}) string {
    return gstr.ToLower(gconv.String(str))
}

// Build-in template function: nl2br
func (view *View) funcNl2Br(str interface{}) string {
    return gstr.Nl2Br(gconv.String(str))
}



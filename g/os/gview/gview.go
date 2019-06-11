<<<<<<< HEAD
// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 视图管理
package gview

import (
    "sync"
    "bytes"
    "errors"
    "html/template"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/os/gspath"
)

// 视图对象
type View struct {
    mu       sync.RWMutex
    paths    *gspath.SPath           // 模板查找目录(绝对路径)
    funcmap  map[string]interface{}  // FuncMap
    contents *gmap.StringStringMap   // 已解析的模板文件内容
}

// 视图表
var viewMap = gmap.NewStringInterfaceMap()

// 默认的视图对象
var viewObj = Get(".")

// 输出到模板页面时保留HTML标签原意，不做自动escape处理
func HTML(content string) template.HTML {
    return template.HTML(content)
}

// 直接解析模板内容，返回解析后的内容
func ParseContent(content string, params map[string]interface{}) ([]byte, error) {
    return viewObj.ParseContent(content, params)
}

// 获取或者创建一个视图对象
func Get(path string) *View {
    if r := viewMap.Get(path); r != nil {
        return r.(*View)
    }
    v := New(path)
    viewMap.Set(path, v)
    return v
}

// 生成一个视图对象
func New(path string) *View {
    s := gspath.New()
    s.Set(path)
    view := &View {
        paths    : s,
        funcmap  : make(map[string]interface{}),
        contents : gmap.NewStringStringMap(),
    }
    view.BindFunc("include", view.funcInclude)
    return view
}

// 设置模板目录绝对路径
func (view *View) SetPath(path string) error {
    return view.paths.Set(path)
}

// 添加模板目录搜索路径
func (view *View) AddPath(path string) error {
    return view.paths.Add(path)
}

// 解析模板，返回解析后的内容
func (view *View) Parse(file string, params map[string]interface{}) ([]byte, error) {
    path    := view.paths.Search(file)
    content := view.contents.Get(path)
    if content == "" {
        content = gfile.GetContents(path)
        if content != "" {
            view.addMonitor(path)
            view.contents.Set(path, content)
        }
    }
    if content == "" {
        return nil, errors.New("tpl \"" + file + "\" not found")
    }
    // 执行模板解析，互斥锁主要是用于funcmap
    view.mu.RLock()
    defer view.mu.RUnlock()
    buffer := bytes.NewBuffer(nil)
    if tpl, err := template.New(path).Funcs(view.funcmap).Parse(content); err != nil {
        return nil, err
    } else {
        if err := tpl.Execute(buffer, params); err != nil {
            return nil, err
        }
    }
    return buffer.Bytes(), nil
}

// 直接解析模板内容，返回解析后的内容
func (view *View) ParseContent(content string, params map[string]interface{}) ([]byte, error) {
    view.mu.RLock()
    defer view.mu.RUnlock()
    name   := gconv.String(ghash.BKDRHash64([]byte(content)))
    buffer := bytes.NewBuffer(nil)
    if tpl, err := template.New(name).Funcs(view.funcmap).Parse(content); err != nil {
        return nil, err
    } else {
        if err := tpl.Execute(buffer, params); err != nil {
            return nil, err
        }
    }
    return buffer.Bytes(), nil
}

// 绑定自定义函数，该函数是全局有效，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) BindFunc(name string, function interface{}) {
    view.mu.Lock()
    view.funcmap[name] = function
    view.mu.Unlock()
}

// 模板内置方法：include
func (view *View) funcInclude(file string, data...map[string]interface{}) template.HTML {
    var m map[string]interface{} = nil
    if len(data) > 0 {
        m = data[0]
    }
    content, err := view.Parse(file, m)
    if err != nil {
        return template.HTML(err.Error())
    }
    return template.HTML(content)
}

// 添加模板文件监控
func (view *View) addMonitor(path string) {
    if view.contents.Get(path) == "" {
        gfsnotify.Add(path, func(event *gfsnotify.Event) {
            if event.IsRemove() {
                gfsnotify.Remove(event.Path)
                return
            }
            view.contents.Remove(event.Path)
        })
    }
}
=======
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
	"github.com/gogf/gf/g/internal/cmdenv"
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gspath"
	"sync"
)

type View struct {
    mu         sync.RWMutex
    paths      *garray.StringArray     // Searching path array.
    data       map[string]interface{}  // Global template variables.
    funcMap    map[string]interface{}  // Global template function map.
    delimiters []string                // Customized template delimiters.
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
func New(path...string) *View {
    view := &View {
        paths      : garray.NewStringArray(),
        data       : make(map[string]interface{}),
        funcMap    : make(map[string]interface{}),
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
                glog.Errorf("Template directory path does not exist: %s", envPath)
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
    view.BindFunc("eq",          view.funcEq)
    view.BindFunc("ne",          view.funcNe)
    view.BindFunc("lt",          view.funcLt)
    view.BindFunc("le",          view.funcLe)
    view.BindFunc("gt",          view.funcGt)
    view.BindFunc("ge",          view.funcGe)
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



>>>>>>> upstream/master

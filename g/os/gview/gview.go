// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gview implements a template engine based on text/template.
// 
// 模板引擎.
package gview

import (
    "bytes"
    "errors"
    "fmt"
    "gitee.com/johng/gf"
    "gitee.com/johng/gf/g/container/garray"
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/encoding/ghtml"
    "gitee.com/johng/gf/g/encoding/gurl"
    "gitee.com/johng/gf/g/os/gfcache"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gspath"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/os/gview/internal/text/template"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/util/gstr"
    "strings"
    "sync"
)

// 视图对象
type View struct {
    mu         sync.RWMutex
    paths      *garray.StringArray     // 模板查找目录(绝对路径)
    data       map[string]interface{}  // 模板变量
    funcmap    map[string]interface{}  // FuncMap
    delimiters []string                // 模板变量分隔符号
}

// 模板变量
type Params  = map[string]interface{}

// 函数映射表
type FuncMap = map[string]interface{}

// 默认的视图对象
var viewObj *View

// 初始化默认的视图对象
func checkAndInitDefaultView() {
    if viewObj == nil {
        // gfile.MainPkgPath() 用以判断是否开发环境
        mainPkgPath := gfile.MainPkgPath()
        if gfile.MainPkgPath() == "" {
            viewObj = New(gfile.SelfDir())
        } else {
            viewObj = New(mainPkgPath)
        }
    }
}

// 直接解析模板内容，返回解析后的内容
func ParseContent(content string, params Params) ([]byte, error) {
    checkAndInitDefaultView()
    return viewObj.ParseContent(content, params)
}

// 生成一个视图对象
func New(path...string) *View {
    view := &View {
        paths      : garray.NewStringArray(0, 1),
        data       : make(map[string]interface{}),
        funcmap    : make(map[string]interface{}),
        delimiters : make([]string, 2),
    }
    if len(path) > 0 && len(path[0]) > 0 {
        view.SetPath(path[0])
    }
    view.SetDelimiters("{{", "}}")
    // 内置变量
    view.data["GF"] = map[string]interface{} {
        "version" : gf.VERSION,
    }
    // 内置方法
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

// 设置模板目录绝对路径
func (view *View) SetPath(path string) error {
    realPath := gfile.RealPath(path)
    if realPath == "" {
        err := errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
        glog.Error(fmt.Sprintf(`[gview] SetPath failed: %s`, err.Error()))
        return err
    }
    view.paths.Clear()
    view.paths.Append(realPath)
    glog.Debug("[gview] SetPath:", realPath)
    return nil
}

// 添加模板目录搜索路径
func (view *View) AddPath(path string) error {
    realPath := gfile.RealPath(path)
    if realPath == "" {
        err := errors.New(fmt.Sprintf(`path "%s" does not exist`, path))
        glog.Error(fmt.Sprintf(`[gview] AddPath failed: %s`, err.Error()))
        return err
    }
    view.paths.Append(realPath)
    glog.Debug("[gview] AddPath:", realPath)
    return nil
}

// 批量绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assigns(data Params) {
    view.mu.Lock()
    for k, v := range data {
        view.data[k] = v
    }
    view.mu.Unlock()
}

// 绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assign(key string, value interface{}) {
    view.mu.Lock()
    view.data[key] = value
    view.mu.Unlock()
}

// 解析模板，返回解析后的内容
func (view *View) Parse(file string, params Params, funcmap...map[string]interface{}) ([]byte, error) {
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
        buffer.WriteString(fmt.Sprintf("[gview] cannot find template file \"%s\" in following paths:", file))
        view.paths.RLockFunc(func(array []string) {
            for k, v := range array {
                buffer.WriteString(fmt.Sprintf("\n%d. %s",k + 1,  v))
            }
        })
        glog.Error(buffer.String())
        return nil, errors.New(fmt.Sprintf(`tpl "%s" not found`, file))
    }
    content := gfcache.GetContents(path)
    // 执行模板解析，互斥锁主要是用于funcmap
    view.mu.RLock()
    defer view.mu.RUnlock()
    buffer := bytes.NewBuffer(nil)
    tplobj := template.New(path).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcmap)
    if len(funcmap) > 0 {
        tplobj = tplobj.Funcs(funcmap[0])
    }
    if tpl, err := tplobj.Parse(content); err != nil {
        return nil, err
    } else {
        // 注意模板变量赋值不能改变已有的params或者view.data的值，因为这两个变量都是指针
        // 因此在必要条件下，需要合并两个map的值到一个新的map
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

// 直接解析模板内容，返回解析后的内容
func (view *View) ParseContent(content string, params Params, funcmap...map[string]interface{}) ([]byte, error) {
    view.mu.RLock()
    defer view.mu.RUnlock()
    name   := gconv.String(ghash.BKDRHash64([]byte(content)))
    buffer := bytes.NewBuffer(nil)
    tplobj := template.New(name).Delims(view.delimiters[0], view.delimiters[1]).Funcs(view.funcmap)
    if len(funcmap) > 0 {
        tplobj = tplobj.Funcs(funcmap[0])
    }
    if tpl, err := tplobj.Parse(content); err != nil {
        return nil, err
    } else {
        // 注意模板变量赋值不能改变已有的params或者view.data的值，因为这两个变量都是指针
        // 因此在必要条件下，需要合并两个map的值到一个新的map
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

// 设置模板变量解析分隔符号
func (view *View) SetDelimiters(left, right string) {
    view.delimiters[0] = left
    view.delimiters[1] = right
}

// 绑定自定义函数，该函数是全局有效，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) BindFunc(name string, function interface{}) {
    view.mu.Lock()
    view.funcmap[name] = function
    view.mu.Unlock()
}

// 模板内置方法：include
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

// 模板内置方法：text
func (view *View) funcText(html interface{}) string {
    return ghtml.StripTags(gconv.String(html))
}

// 模板内置方法：html
func (view *View) funcHtmlEncode(html interface{}) string {
    return ghtml.Entities(gconv.String(html))
}

// 模板内置方法：htmldecode
func (view *View) funcHtmlDecode(html interface{}) string {
    return ghtml.EntitiesDecode(gconv.String(html))
}

// 模板内置方法：url
func (view *View) funcUrlEncode(url interface{}) string {
    return gurl.Encode(gconv.String(url))
}

// 模板内置方法：urldecode
func (view *View) funcUrlDecode(url interface{}) string {
    if content, err := gurl.Decode(gconv.String(url)); err == nil {
        return content
    } else {
        return err.Error()
    }
}

// 模板内置方法：date
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

// 模板内置方法：compare
func (view *View) funcCompare(value1, value2 interface{}) int {
    return strings.Compare(gconv.String(value1), gconv.String(value2))
}

// 模板内置方法：substr
func (view *View) funcSubStr(start, end int, str interface{}) string {
    return gstr.SubStr(gconv.String(str), start, end)
}

// 模板内置方法：strlimit
func (view *View) funcStrLimit(length int, suffix string, str interface{}) string {
    return gstr.StrLimit(gconv.String(str), length, suffix)
}

// 模板内置方法：highlight
func (view *View) funcHighlight(key string, color string, str interface{}) string {
    return gstr.Replace(gconv.String(str), key, fmt.Sprintf(`<span style="color:%s;">%s</span>`, color, key))
}

// 模板内置方法：hidestr
func (view *View) funcHideStr(percent int, hide string, str interface{}) string {
    return gstr.HideStr(gconv.String(str), percent, hide)
}

// 模板内置方法：toupper
func (view *View) funcToUpper(str interface{}) string {
    return gstr.ToUpper(gconv.String(str))
}

// 模板内置方法：toupper
func (view *View) funcToLower(str interface{}) string {
    return gstr.ToLower(gconv.String(str))
}

// 模板内置方法：nl2br
func (view *View) funcNl2Br(str interface{}) string {
    return gstr.Nl2Br(gconv.String(str))
}



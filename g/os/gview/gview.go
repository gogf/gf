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
    "strings"
    "html/template"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/encoding/ghash"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/os/gfsnotify"
)

// 视图对象
type View struct {
    mu       sync.RWMutex
    path     *gtype.String           // 模板目录(绝对路径)
    funcmap  map[string]interface{}  // FuncMap
    contents *gmap.StringStringMap   // 已解析的模板文件内容
}

// 视图表
var viewMap = gmap.NewStringInterfaceMap()

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
    view := &View {
        path     : gtype.NewString(path),
        funcmap  : make(map[string]interface{}),
        contents : gmap.NewStringStringMap(),
    }
    view.BindFunc("include", view.funcInclude)
    return view
}

// 设置模板目录绝对路径
func (view *View) SetPath(path string) {
    view.path.Set(path)
}

// 获取模板目录绝对路径
func (view *View) GetPath() string {
    return view.path.Val()
}

// 解析模板，返回解析后的内容
func (view *View) Parse(file string, params map[string]interface{}) ([]byte, error) {
    path    := strings.TrimRight(view.GetPath(), gfile.Separator) + gfile.Separator + file
    content := view.contents.Get(path)
    if content == "" {
        content = gfile.GetContents(path)
        if content != "" {
            view.mu.Lock()
            view.addMonitor(path)
            view.contents.Set(path, content)
            view.mu.Unlock()
        }
    }
    if content == "" {
        return nil, errors.New("invalid tpl \"" + file + "\"")
    }
    // 执行模板解析
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

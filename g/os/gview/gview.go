// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 视图管理
package gview

import (
    "gitee.com/johng/gf/g/container/gmap"
    "html/template"
    "gitee.com/johng/gf/g/os/gfile"
    "sync"
    "strings"
    "bytes"
    "errors"
)

// 视图对象
type View struct {
    mu     sync.RWMutex
    path   string                   // 模板目录(绝对路径)
    tpls   *gmap.StringInterfaceMap // 已解析的模板对象指针，防止重复解析
    suffix string                   // 模板文件名后缀
}

// 模板对象
type Template struct {
    mu        sync.RWMutex            // 并发互斥锁
    path      string                  // 模板文件(绝对路径)
    data      map[string]interface{}  // 全局的模板变量
    content   string                  // 模板文件内容(解析之后保存到内存中)
    funcmap   map[string]interface{}  // FuncMap
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
    return &View {
        path   : path,
        tpls   : gmap.NewStringInterfaceMap(),
        suffix : "tpl",
    }
}

// 设置模板目录绝对路径
func (view *View) SetPath(path string) {
    view.mu.Lock()
    defer view.mu.Unlock()
    view.path = path
}

// 获取模板目录绝对路径
func (view *View) GetPath() string {
    view.mu.RLock()
    defer view.mu.RUnlock()
    return view.path
}

// 直接解析模板，返回解析后的内容
func (view *View) Parse(file string, params map[string]interface{}) ([]byte, error) {
    if t, err := view.Template(file); err == nil {
        return t.Parse(params)
    } else {
        return nil, err
    }
}

// 根据文件名称生成一个模板对象，或者获取一个现有的模板对象
func (view *View) Template(file string) (*Template, error) {
    path := strings.TrimRight(view.GetPath(), gfile.Separator) + gfile.Separator + file + "." + view.suffix
    if t := view.tpls.Get(path); t != nil {
        return t.(*Template), nil
    }
    if !gfile.Exists(path) {
        return nil, errors.New("template '" + path + "' does not exist")
    }
    if !gfile.IsReadable(path) {
        return nil, errors.New("template '" + path + "' is not readable")
    }
    t := &Template{
        path     : path,
        data     : make(map[string]interface{}),
        content  : gfile.GetContents(path),
        funcmap  : make(map[string]interface{}),
    }
    // 绑定内置inluce方法
    t.BindFunc("include", t.funcInclude)
    // 将模板对象注册到内部map中
    view.tpls.Set(path, t)
    return t, nil
}

// 绑定自定义函数，该函数是全局有效，即调用之后每个线程都会生效，因此有并发安全控制
func (t *Template) BindFunc(name string, function interface{}) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.funcmap[name] = function
}

// 批量绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (t *Template) Assigns(data map[string]interface{}) {
    t.mu.Lock()
    defer t.mu.Unlock()
    for k, v := range data {
        t.data[k] = v
    }
}

// 绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (t *Template) Assign(k string, v interface{}) {
    t.mu.Lock()
    defer t.mu.Unlock()
    t.data[k] = v
}

// 返回解析后的模板内容，可以额外指定模板变量，如果没有可以传入nil
// 函数内部的底层template必须每次调用都新生成一个，防止错误：html/template: cannot Parse after Execute
func (t *Template) Parse(data map[string]interface{}) ([]byte, error) {
    t.mu.RLock()
    defer t.mu.RUnlock()
    buffer := bytes.NewBuffer(nil)
    if tpl, err := template.New(t.path).Funcs(t.funcmap).Parse(t.content); err != nil {
        return nil, err
    } else {
        m := t.data
        for k, v := range data {
            m[k] = v
        }
        if err := tpl.Execute(buffer, m); err != nil {
            return nil, err
        }
    }
    return buffer.Bytes(), nil
}

// 模板内置方法：include
func (t *Template) funcInclude(file string) template.HTML {
    content, err := t.Parse(nil)
    if err != nil {
        return template.HTML(err.Error())
    }
    return template.HTML(content)
}


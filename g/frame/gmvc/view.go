// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmvc

import (
    "sync"
    "github.com/gogf/gf/g/os/gview"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/frame/gins"
)

// 基于控制器注册的MVC视图基类(一个请求一个视图对象，用完即销毁)
type View struct {
    mu       sync.RWMutex              // 并发互斥锁
    view     *gview.View               // 底层视图对象
    data     gview.Params              // 视图数据/模板变量
    response *ghttp.Response           // 数据返回对象
}

// 创建一个MVC请求中使用的视图对象
func NewView(w *ghttp.Response) *View {
    return &View {
        view     : gins.View(),
        data     : make(gview.Params),
        response : w,
    }
}

// 批量绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assigns(data gview.Params) {
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

// 解析模板，并返回解析后的内容
func (view *View) Parse(file string) (string, error) {
    view.mu.RLock()
    defer view.mu.RUnlock()
    buffer, err := view.response.ParseTpl(file, view.data)
    return buffer, err
}

// 直接解析模板内容，并返回解析后的内容
func (view *View) ParseContent(content string) (string, error) {
    view.mu.RLock()
    defer view.mu.RUnlock()
    buffer, err := view.response.ParseTplContent(content, view.data)
    return buffer, err
}

// 使用自定义方法对模板变量执行加锁修改操作
func (view *View) LockFunc(f func(data gview.Params)) {
    view.mu.Lock()
    defer view.mu.Unlock()
    f(view.data)
}

// 使用自定义方法对模板变量执行加锁读取操作
func (view *View) RLockFunc(f func(data gview.Params)) {
    view.mu.RLock()
    defer view.mu.RUnlock()
    f(view.data)
}

// 解析并显示指定模板
func (view *View) Display(file...string) error {
    name := "index.tpl"
    if len(file) > 0 {
        name = file[0]
    }
    if content, err := view.Parse(name); err != nil {
        view.response.Write("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        view.response.Write(content)
    }
    return nil
}

// 解析并显示模板内容
func (view *View) DisplayContent(content string) error {
    if content, err := view.ParseContent(content); err != nil {
        view.response.Write("Tpl Parsing Error: " + err.Error())
        return err
    } else {
        view.response.Write(content)
    }
    return nil
}
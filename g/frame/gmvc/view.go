// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gmvc

import (
    "sync"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gins"
)

// 基于控制器注册的MVC视图基类(一个请求一个视图对象，用完即销毁)
type View struct {
    mu       sync.RWMutex              // 并发互斥锁
    view     *gview.View               // 底层视图对象
    data     map[string]interface{}    // 视图数据/模板变量
    response *ghttp.Response           // 数据返回对象
}

// 创建一个MVC请求中使用的视图对象
func NewView(w *ghttp.Response) *View {
    return &View {
        view     : gins.View(),
        data     : make(map[string]interface{}),
        response : w,
    }
}

// 批量绑定模板变量，即调用之后每个线程都会生效，因此有并发安全控制
func (view *View) Assigns(data map[string]interface{}) {
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
func (view *View) Parse(file string) ([]byte, error) {
    view.mu.RLock()
    buffer, err := view.view.Parse(file, view.data)
    view.mu.RUnlock()
    return buffer, err
}

// 直接解析模板内容，并返回解析后的内容
func (view *View) ParseContent(content string) ([]byte, error) {
    view.mu.RLock()
    buffer, err := view.view.ParseContent(content, view.data)
    view.mu.RUnlock()
    return buffer, err
}

// 使用自定义方法对模板变量执行加锁修改操作
func (view *View) LockFunc(f func(vars map[string]interface{})) {
    view.mu.Lock()
    f(view.data)
    view.mu.Unlock()
}

// 使用自定义方法对模板变量执行加锁读取操作
func (view *View) RLockFunc(f func(vars map[string]interface{})) {
    view.mu.RLock()
    f(view.data)
    view.mu.RUnlock()
}

// 解析指定模板
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
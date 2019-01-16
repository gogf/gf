// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gmvc provides basic object classes for MVC.
package gmvc

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

// 控制器基类
type Controller struct {
    Request  *ghttp.Request  // 请求数据对象
    Response *ghttp.Response // 返回数据对象(r.Response)
    Server   *ghttp.Server   // Web Server对象(r.Server)
    Cookie   *ghttp.Cookie   // COOKIE操作对象(r.Cookie)
    Session  *ghttp.Session  // SESSION操作对象
    View     *View           // 视图对象
}

// 控制器初始化接口方法
func (c *Controller) Init(r *ghttp.Request) {
    c.Request  = r
    c.Response = r.Response
    c.Server   = r.Server
    c.View     = NewView(r.Response)
    c.Cookie   = r.Cookie
    c.Session  = r.Session
}

// 控制器结束请求接口方法
func (c *Controller) Shut() {

}

// 退出请求执行
func (c *Controller) Exit() {
    c.Request.Exit()
}



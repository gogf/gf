// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// MVC控制器基类
package gmvc

import (
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/net/gsession"
)

// 控制器基类
type Controller struct {
    Server   *ghttp.Server         // Web Server对象
    Request  *ghttp.ClientRequest  // 请求数据对象
    Response *ghttp.ServerResponse // 返回数据对象
    Cookie   *ghttp.Cookie         // COOKIE操作对象
    Session  *gsession.Session     // SESSION操作对象
    View     *View                 // 视图对象
}

// 控制器初始化接口方法
func (c *Controller) Init(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    c.Server   = s
    c.Request  = r
    c.Response = w
    c.Cookie   = ghttp.NewCookie(c.Request, c.Response)
    c.View     = NewView(c)
    c.Session  = gsession.Get(c.Cookie.SessionId())
}

// 控制器结束请求接口方法
func (c *Controller) Shut() {

}



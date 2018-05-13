// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// pprof封装.

package ghttp

import (
    "strings"
    "gitee.com/johng/gf/g/os/gview"
)

// 用于服务管理的对象
type utilAdmin struct {}

// 服务管理首页
func (p *utilAdmin) Index(r *Request) {
    data := map[string]interface{}{
        "uri" : strings.TrimRight(r.URL.Path, "/"),
    }
    buffer, _ := gview.ParseContent(`
            <html>
            <head>
                <title>gf ghttp admin</title>
            </head>
            <body>
                <p><a href="{{$.uri}}/restart">restart</a></p>
                <p><a href="{{$.uri}}/shutdown">shutdown</a></p>
            </body>
            </html>
            `, data)
    r.Response.Write(buffer)
}

// 服务重启
func (p *utilAdmin) Restart(r *Request) {
    r.Response.Write("restart server")
    r.Server.Restart()
}

// 服务关闭
func (p *utilAdmin) Shutdown(r *Request) {
    r.Response.Write("shutdown server")
    r.Server.Shutdown()
}


// 开启服务管理支持
func (s *Server) EnableAdmin(pattern...string) {
    p := "/debug/admin"
    if len(pattern) > 0 {
        p = pattern[0]
    }
    s.BindObject(p, &utilAdmin{})
}
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// pprof封装.

package ghttp

import (
	"github.com/gogf/gf/g/os/gproc"
	"github.com/gogf/gf/g/os/gtimer"
	"github.com/gogf/gf/g/os/gview"
	"os"
	"strings"
	"time"
)

// 服务管理首页
func (p *utilAdmin) Index(r *Request) {
	data := map[string]interface{}{
		"pid": gproc.Pid(),
		"uri": strings.TrimRight(r.URL.Path, "/"),
	}
	buffer, _ := gview.ParseContent(`
            <html>
            <head>
                <title>GoFrame Web Server Admin</title>
            </head>
            <body>
                <p>PID: {{.pid}}</p>
                <p><a href="{{$.uri}}/restart">Restart</a></p>
                <p><a href="{{$.uri}}/shutdown">Shutdown</a></p>
            </body>
            </html>
    `, data)
	r.Response.Write(buffer)
}

// 服务重启
func (p *utilAdmin) Restart(r *Request) {
	var err error = nil
	// 必须检查可执行文件的权限
	path := r.GetQueryString("newExeFilePath")
	if path == "" {
		path = os.Args[0]
	}
	// 执行重启操作
	if len(path) > 0 {
		err = RestartAllServer(path)
	} else {
		err = RestartAllServer()
	}
	if err == nil {
		r.Response.Write("server restarted")
	} else {
		r.Response.Write(err.Error())
	}
}

// 服务关闭
func (p *utilAdmin) Shutdown(r *Request) {
	r.Server.Shutdown()
	if err := ShutdownAllServer(); err == nil {
		r.Response.Write("server shutdown")
	} else {
		r.Response.Write(err.Error())
	}
}

// 开启服务管理支持
func (s *Server) EnableAdmin(pattern ...string) {
	p := "/debug/admin"
	if len(pattern) > 0 {
		p = pattern[0]
	}
	s.BindObject(p, &utilAdmin{})
}

// 关闭当前Web Server
func (s *Server) Shutdown() error {
	// 非终端信号下，异步1秒后再执行关闭，
	// 目的是让接口能够正确返回结果，否则接口会报错(因为web server关闭了)
	gtimer.SetTimeout(time.Second, func() {
		// 只关闭当前的Web Server
		for _, v := range s.servers {
			v.close()
		}
	})
	return nil
}

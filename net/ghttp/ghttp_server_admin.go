// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/os/gview"
)

// utilAdmin is the controller for administration.
type utilAdmin struct{}

// Index shows the administration page.
func (p *utilAdmin) Index(r *Request) {
	data := map[string]any{
		"pid":  gproc.Pid(),
		"path": gfile.SelfPath(),
		"uri":  strings.TrimRight(r.URL.Path, "/"),
	}
	buffer, _ := gview.ParseContent(r.Context(), `
            <html>
            <head>
                <title>GoFrame Web Server Admin</title>
            </head>
            <body>
                <p>Pid: {{.pid}}</p>
                <p>File Path: {{.path}}</p>
                <p>
<a href="{{$.uri}}/restart">Restart</a>
please make sure it is running using standalone binary not from IDE or "go run"
</p>
                <p>
<a href="{{$.uri}}/shutdown">Shutdown</a>
graceful shutdown the server
</p>
            </body>
            </html>
    `, data)
	r.Response.Write(buffer)
}

// Restart restarts all the servers in the process.
func (p *utilAdmin) Restart(r *Request) {
	var (
		ctx = r.Context()
		err error
	)
	// Custom start binary path when this process exits.
	path := r.GetQuery("newExeFilePath").String()
	if path == "" {
		path = os.Args[0]
	}
	if err = RestartAllServer(ctx, path); err == nil {
		r.Response.WriteExit("server restarted")
	} else {
		r.Response.WriteExit(err.Error())
	}
}

// Shutdown shuts down all the servers.
func (p *utilAdmin) Shutdown(r *Request) {
	gtimer.SetTimeout(r.Context(), time.Second, func(ctx context.Context) {
		// It shuts down the server after 1 second, which is not triggered by system signal,
		// to ensure the response successfully to the client.
		_ = r.Server.Shutdown()
	})
	r.Response.WriteExit("server shutdown")
}

// EnableAdmin enables the administration feature for the process.
// The optional parameter `pattern` specifies the URI for the administration page.
func (s *Server) EnableAdmin(pattern ...string) {
	p := "/debug/admin"
	if len(pattern) > 0 {
		p = pattern[0]
	}
	s.BindObject(p, &utilAdmin{})
}

// Shutdown shuts the current server.
func (s *Server) Shutdown() error {
	var ctx = context.TODO()
	// Remove plugins.
	if len(s.plugins) > 0 {
		for _, p := range s.plugins {
			s.Logger().Printf(ctx, `remove plugin: %s`, p.Name())
			if err := p.Remove(); err != nil {
				s.Logger().Errorf(ctx, "%+v", err)
			}
		}
	}

	s.doServiceDeregister()
	// Only shut down current servers.
	// It may have multiple underlying http servers.
	for _, v := range s.servers {
		v.Shutdown(ctx)
	}
	s.Logger().Infof(ctx, "pid[%d]: all servers shutdown", gproc.Pid())
	return nil
}

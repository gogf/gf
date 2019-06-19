// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// pprof封装.

package ghttp

import (
	"github.com/gogf/gf/g/os/gview"
	netpprof "net/http/pprof"
	runpprof "runtime/pprof"
	"strings"
)

// 用于pprof的对象
type utilPprof struct{}

func (p *utilPprof) Index(r *Request) {
	profiles := runpprof.Profiles()
	action := r.Get("action")
	data := map[string]interface{}{
		"uri":      strings.TrimRight(r.URL.Path, "/") + "/",
		"profiles": profiles,
	}
	if len(action) == 0 {
		buffer, _ := gview.ParseContent(`
            <html>
            <head>
                <title>gf ghttp pprof</title>
            </head>
            {{$uri := .uri}}
            <body>
                profiles:<br>
                <table>
                    {{range .profiles}}<tr><td align=right>{{.Count}}<td><a href="{{$uri}}{{.Name}}?debug=1">{{.Name}}</a>{{end}}
                </table>
                <br><a href="{{$uri}}goroutine?debug=2">full goroutine stack dump</a><br>
            </body>
            </html>
            `, data)
		r.Response.Write(buffer)
		return
	}
	for _, p := range profiles {
		if p.Name() == action {
			p.WriteTo(r.Response.Writer, r.GetRequestInt("debug"))
			break
		}
	}
}

func (p *utilPprof) Cmdline(r *Request) {
	netpprof.Cmdline(r.Response.Writer, r.Request)
}

func (p *utilPprof) Profile(r *Request) {
	netpprof.Profile(r.Response.Writer, r.Request)
}

func (p *utilPprof) Symbol(r *Request) {
	netpprof.Symbol(r.Response.Writer, r.Request)
}

func (p *utilPprof) Trace(r *Request) {
	netpprof.Trace(r.Response.Writer, r.Request)
}

// 开启pprof支持
func (s *Server) EnablePprof(pattern ...string) {
	p := "/debug/pprof"
	if len(pattern) > 0 {
		p = pattern[0]
	}
	up := &utilPprof{}
	_, _, uri, _ := s.parsePattern(p)
	uri = strings.TrimRight(uri, "/")
	s.BindHandler(uri+"/*action", up.Index)
	s.BindHandler(uri+"/cmdline", up.Cmdline)
	s.BindHandler(uri+"/profile", up.Profile)
	s.BindHandler(uri+"/symbol", up.Symbol)
	s.BindHandler(uri+"/trace", up.Trace)
}

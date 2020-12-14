// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	netpprof "net/http/pprof"
	runpprof "runtime/pprof"
	"strings"

	"github.com/gogf/gf/os/gview"
)

// utilPProf is the PProf interface implementer.
type utilPProf struct{}

const (
	gDEFAULT_PPROF_PATTERN = "/debug/pprof"
)

// EnablePProf enables PProf feature for server.
func (s *Server) EnablePProf(pattern ...string) {
	s.Domain(defaultDomainName).EnablePProf(pattern...)
}

// EnablePProf enables PProf feature for server of specified domain.
func (d *Domain) EnablePProf(pattern ...string) {
	p := gDEFAULT_PPROF_PATTERN
	if len(pattern) > 0 && pattern[0] != "" {
		p = pattern[0]
	}
	up := &utilPProf{}
	_, _, uri, _ := d.server.parsePattern(p)
	uri = strings.TrimRight(uri, "/")
	d.Group(uri, func(group *RouterGroup) {
		group.ALL("/*action", up.Index)
		group.ALL("/cmdline", up.Cmdline)
		group.ALL("/profile", up.Profile)
		group.ALL("/symbol", up.Symbol)
		group.ALL("/trace", up.Trace)
	})
}

// Index shows the PProf index page.
func (p *utilPProf) Index(r *Request) {
	profiles := runpprof.Profiles()
	action := r.GetString("action")
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

// Cmdline responds with the running program's
// command line, with arguments separated by NUL bytes.
// The package initialization registers it as /debug/pprof/cmdline.
func (p *utilPProf) Cmdline(r *Request) {
	netpprof.Cmdline(r.Response.Writer, r.Request)
}

// Profile responds with the pprof-formatted cpu profile.
// Profiling lasts for duration specified in seconds GET parameter, or for 30 seconds if not specified.
// The package initialization registers it as /debug/pprof/profile.
func (p *utilPProf) Profile(r *Request) {
	netpprof.Profile(r.Response.Writer, r.Request)
}

// Symbol looks up the program counters listed in the request,
// responding with a table mapping program counters to function names.
// The package initialization registers it as /debug/pprof/symbol.
func (p *utilPProf) Symbol(r *Request) {
	netpprof.Symbol(r.Response.Writer, r.Request)
}

// Trace responds with the execution trace in binary form.
// Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified.
// The package initialization registers it as /debug/pprof/trace.
func (p *utilPProf) Trace(r *Request) {
	netpprof.Trace(r.Response.Writer, r.Request)
}

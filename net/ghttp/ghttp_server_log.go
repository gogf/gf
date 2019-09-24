// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"fmt"

	"github.com/gogf/gf/errors/gerror"

	"github.com/gogf/gf/os/gtime"
)

const (
	gPATH_FILTER_KEY = "/net/ghttp/ghttp"
)

// 处理服务错误信息，主要是panic，http请求的status由access log进行管理
func (s *Server) handleAccessLog(r *Request) {
	if !s.IsAccessLogEnabled() {
		return
	}
	// 自定义错误处理
	if v := s.GetLogHandler(); v != nil {
		v(r)
		return
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	content := fmt.Sprintf(`%d "%s %s %s %s %s"`,
		r.Response.Status,
		r.Method, scheme, r.Host, r.URL.String(), r.Proto,
	)
	content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime-r.EnterTime)/1000)
	content += fmt.Sprintf(`, %s, "%s", "%s"`, r.GetClientIp(), r.Referer(), r.UserAgent())
	s.logger.Cat("access").StackWithFilter(gPATH_FILTER_KEY).Stdout(s.config.LogStdout).Println(content)
}

// 处理服务错误信息，主要是panic，http请求的status由access log进行管理
func (s *Server) handleErrorLog(err error, r *Request) {
	// 错误输出默认是开启的
	if !s.IsErrorLogEnabled() {
		return
	}

	// 自定义错误处理
	if v := s.GetLogHandler(); v != nil {
		v(r, err)
		return
	}

	// 错误日志信息
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	content := fmt.Sprintf(`%v, "%s %s %s %s %s"`, err, r.Method, scheme, r.Host, r.URL.String(), r.Proto)
	if r.LeaveTime > r.EnterTime {
		content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime-r.EnterTime)/1000)
	} else {
		content += fmt.Sprintf(` %.3f`, float64(gtime.Microsecond()-r.EnterTime)/1000)
	}
	content += fmt.Sprintf(`, %s, "%s", "%s"`, r.GetClientIp(), r.Referer(), r.UserAgent())
	if s.config.ErrorStack {
		if stack := gerror.Stack(err); stack != "" {
			content += "\n" + stack
		}
	}
	s.logger.Cat("error").Stack(false).Stdout(s.config.LogStdout).Error(content)
}

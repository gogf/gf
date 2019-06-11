// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// 默认错误日志封装.

package ghttp

import (
    "fmt"
    "github.com/gogf/gf/g/os/gtime"
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
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`, r.GetClientIp(), r.Referer(), r.UserAgent())
    s.logger.Cat("access").Backtrace(false, 2).Stdout(s.config.LogStdout).Println(content)
}

// 处理服务错误信息，主要是panic，http请求的status由access log进行管理
func (s *Server) handleErrorLog(error interface{}, r *Request) {
    // 错误输出默认是开启的
    if !s.IsErrorLogEnabled() {
        return
    }

    // 自定义错误处理
    if v := s.GetLogHandler(); v != nil {
        v(r, error)
        return
    }

    // 错误日志信息
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
    content := fmt.Sprintf(`%v, "%s %s %s %s %s"`, error, r.Method, scheme, r.Host, r.URL.String(), r.Proto)
    if r.LeaveTime > r.EnterTime {
        content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    } else {
        content += fmt.Sprintf(` %.3f`, float64(gtime.Microsecond() - r.EnterTime)/1000)
    }
    content += fmt.Sprintf(`, %s, "%s", "%s"`,  r.GetClientIp(), r.Referer(), r.UserAgent())
    s.logger.Cat("error").Backtrace(true, 2).Stdout(s.config.LogStdout).Error(content)
}

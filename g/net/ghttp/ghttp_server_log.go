// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 默认错误日志封装.

package ghttp

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
    "net/http"
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
    content := fmt.Sprintf(`"%s %s %s %s" %s %s`,
        r.Method, r.Host, r.URL.String(), r.Proto,
        gconv.String(r.Response.Status),
        gconv.String(r.Response.Length),
    )
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`, r.GetClientIp(), r.Referer(), r.UserAgent())
    s.accessLogger.Println(content)
}

// 处理服务错误信息，主要是panic，http请求的status由access log进行管理
func (s *Server) handleErrorLog(error interface{}, r *Request) {
    r.Response.WriteStatus(http.StatusInternalServerError)

    if !s.IsErrorLogEnabled() {
        return
    }
    // 自定义错误处理
    if v := s.GetLogHandler(); v != nil {
        v(r, error)
        return
    }

    content := fmt.Sprintf(`%v, "%s %s %s %s"`, error, r.Method, r.Host, r.URL.String(), r.Proto)
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`,  r.GetClientIp(), r.Referer(), r.UserAgent())
    s.errorLogger.Error(content)
}

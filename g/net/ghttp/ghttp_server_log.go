// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 默认错误日志封装.

package ghttp

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfile"
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
    content := fmt.Sprintf(`"%s %s %s %s" %d`,
        r.Method, r.Host, r.URL.String(), r.Proto,
        r.Response.Status,
    )
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`, r.GetClientIp(), r.Referer(), r.UserAgent())
    s.logger.Cat("access").Backtrace(false, 2).Println(content)
}

// 处理服务错误信息，主要是panic，http请求的status由access log进行管理
func (s *Server) handleErrorLog(error interface{}, r *Request) {
    r.Response.WriteStatus(http.StatusInternalServerError)

    // 错误输出默认是开启的
    if !s.IsErrorLogEnabled() && gfile.MainPkgPath() == "" {
        return
    }

    // 自定义错误处理
    if v := s.GetLogHandler(); v != nil {
        v(r, error)
        return
    }

    // 错误日志信息
    content := fmt.Sprintf(`%v, "%s %s %s %s"`, error, r.Method, r.Host, r.URL.String(), r.Proto)
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`,  r.GetClientIp(), r.Referer(), r.UserAgent())

    if s.logger.GetPath() == "" {
        // 错误信息特殊处理，在未开启日志文件保存时强制强制输出到终端
        s.logger.Cat("error").Backtrace(true, 2).StdPrint(true).Error(content)
    } else {
        s.logger.Cat("error").Backtrace(true, 2).Error(content)
        // 开发环境下(MainPkgPath)自动输出错误信息到标准输出
        if gfile.MainPkgPath() != "" {
            s.logger.Cat("error").Backtrace(true, 2).StdPrint(true).Error(content)
        }
    }
}

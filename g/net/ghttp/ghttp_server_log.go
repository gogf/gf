<<<<<<< HEAD
// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
=======
// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
>>>>>>> upstream/master
// 默认错误日志封装.

package ghttp

import (
    "fmt"
<<<<<<< HEAD
    "gitee.com/johng/gf/g/util/gconv"
    "net/http"
=======
    "github.com/gogf/gf/g/os/gtime"
>>>>>>> upstream/master
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
<<<<<<< HEAD
    content := fmt.Sprintf(`"%s %s %s %s" %s %s`,
        r.Method, r.Host, r.URL.String(), r.Proto,
        gconv.String(r.Response.Status),
        gconv.String(r.Response.Length),
    )
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`, r.GetClientIp(), r.Referer(), r.UserAgent())
    s.accessLogger.Println(content)
=======
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
>>>>>>> upstream/master
}

// 处理服务错误信息，主要是panic，http请求的status由access log进行管理
func (s *Server) handleErrorLog(error interface{}, r *Request) {
<<<<<<< HEAD
    r.Response.WriteStatus(http.StatusInternalServerError)

    if !s.IsErrorLogEnabled() {
        return
    }
=======
    // 错误输出默认是开启的
    if !s.IsErrorLogEnabled() {
        return
    }

>>>>>>> upstream/master
    // 自定义错误处理
    if v := s.GetLogHandler(); v != nil {
        v(r, error)
        return
    }

<<<<<<< HEAD
    content := fmt.Sprintf(`%v, "%s %s %s %s"`, error, r.Method, r.Host, r.URL.String(), r.Proto)
    content += fmt.Sprintf(` %.3f`, float64(r.LeaveTime - r.EnterTime)/1000)
    content += fmt.Sprintf(`, %s, "%s", "%s"`,  r.GetClientIp(), r.Referer(), r.UserAgent())
    s.errorLogger.Error(content)
=======
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
>>>>>>> upstream/master
}

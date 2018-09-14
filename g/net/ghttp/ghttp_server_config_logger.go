// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 配置管理数据结构定义.

package ghttp

import (
    "gitee.com/johng/gf/g/os/glog"
)

// 设置日志目录
func (s *Server)SetLogPath(path string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error("cannot be changed while running")
    }
    if len(path) == 0 {
        return
    }
    s.config.LogPath = path
}

// 设置是否开启access log日志功能
func (s *Server)SetAccessLogEnabled(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error("cannot be changed while running")
    }
    s.config.AccessLogEnabled = enabled
}

// 设置是否开启error log日志功能
func (s *Server)SetErrorLogEnabled(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error("cannot be changed while running")
    }
    s.config.ErrorLogEnabled = enabled
}

// 设置日志写入的回调函数
func (s *Server) SetLogHandler(handler LogHandler) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error("cannot be changed while running")
    }
    s.config.LogHandler = handler
}

// 获取日志写入的回调函数
func (s *Server) GetLogHandler() LogHandler {
    return s.config.LogHandler
}

// 获取日志目录
func (s *Server)GetLogPath() string {
    return s.config.LogPath
}

// access log日志功能是否开启
func (s *Server)IsAccessLogEnabled() bool {
    return s.config.AccessLogEnabled
}

// error log日志功能是否开启
func (s *Server)IsErrorLogEnabled() bool {
    return s.config.ErrorLogEnabled
}

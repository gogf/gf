// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
    "github.com/gogf/gf/g/os/glog"
)

// 设置日志目录，只有在设置了日志目录的情况下才会输出日志到日志文件中。
// 日志文件路径格式为：
// 1. 请求日志: access/YYYY-MM-DD.log
// 2. 错误日志: error/YYYY-MM-DD.log
func (s *Server)SetLogPath(path string) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    if len(path) == 0 {
        return
    }
    s.config.LogPath = path
    s.logger.SetPath(path)
}

// 设置日志内容是否输出到终端，默认情况下只有错误日志才会自动输出到终端。
// 如果需要输出请求日志到终端，默认情况下使用SetAccessLogEnabled方法开启请求日志特性即可。
func (s *Server)SetLogStdout(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.LogStdout = enabled
}

// 设置是否开启access log日志功能
func (s *Server)SetAccessLogEnabled(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.AccessLogEnabled = enabled
}

// 设置是否开启error log日志功能
func (s *Server)SetErrorLogEnabled(enabled bool) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.ErrorLogEnabled = enabled
}

// 设置日志写入的回调函数
func (s *Server) SetLogHandler(handler LogHandler) {
    if s.Status() == SERVER_STATUS_RUNNING {
        glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
        return
    }
    s.config.LogHandler = handler
}

// 获取日志写入的回调函数
func (s *Server)GetLogHandler() LogHandler {
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

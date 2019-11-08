// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// 设置日志目录，只有在设置了日志目录的情况下才会输出日志到日志文件中。
func (s *Server) SetLogPath(path string) {
	if len(path) == 0 {
		return
	}
	s.config.LogPath = path
	s.config.ErrorLogEnabled = true
	s.config.AccessLogEnabled = true
}

// 设置日志内容是否输出到终端，默认情况下只有错误日志才会自动输出到终端。
// 如果需要输出请求日志到终端，默认情况下使用SetAccessLogEnabled方法开启请求日志特性即可。
func (s *Server) SetLogStdout(enabled bool) {
	s.config.LogStdout = enabled
}

// 设置是否开启access log日志功能
func (s *Server) SetAccessLogEnabled(enabled bool) {
	s.config.AccessLogEnabled = enabled
}

// 设置是否开启error log日志功能
func (s *Server) SetErrorLogEnabled(enabled bool) {
	s.config.ErrorLogEnabled = enabled
}

// 设置是否开启error stack打印功能
func (s *Server) SetErrorStack(enabled bool) {
	s.config.ErrorStack = enabled
}

// 获取日志目录
func (s *Server) GetLogPath() string {
	return s.config.LogPath
}

// access log日志功能是否开启
func (s *Server) IsAccessLogEnabled() bool {
	return s.config.AccessLogEnabled
}

// error log日志功能是否开启
func (s *Server) IsErrorLogEnabled() bool {
	return s.config.ErrorLogEnabled
}

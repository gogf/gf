// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/internal/intlog"

// SetLogPath sets the log path for server.
// It logs content to file only if the log path is set.
func (s *Server) SetLogPath(path string) {
	if len(path) == 0 {
		return
	}
	intlog.Print("SetLogPath:", path)
	s.config.LogPath = path
	s.config.ErrorLogEnabled = true
	s.config.AccessLogEnabled = true
}

// SetLogStdout sets whether output the logging content to stdout.
func (s *Server) SetLogStdout(enabled bool) {
	s.config.LogStdout = enabled
}

// SetAccessLogEnabled enables/disables the access log.
func (s *Server) SetAccessLogEnabled(enabled bool) {
	s.config.AccessLogEnabled = enabled
}

// SetErrorLogEnabled enables/disables the error log.
func (s *Server) SetErrorLogEnabled(enabled bool) {
	s.config.ErrorLogEnabled = enabled
}

// SetErrorStack enables/disables the error stack feature.
func (s *Server) SetErrorStack(enabled bool) {
	s.config.ErrorStack = enabled
}

// GetLogPath returns the log path.
func (s *Server) GetLogPath() string {
	return s.config.LogPath
}

// IsAccessLogEnabled checks whether the access log enabled.
func (s *Server) IsAccessLogEnabled() bool {
	return s.config.AccessLogEnabled
}

// IsErrorLogEnabled checks whether the error log enabled.
func (s *Server) IsErrorLogEnabled() bool {
	return s.config.ErrorLogEnabled
}

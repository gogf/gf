// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/os/glog"

// 设置http server参数 - SessionMaxAge
func (s *Server) SetSessionMaxAge(age int64) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.SessionMaxAge = age
}

// 设置http server参数 - SessionIdName
func (s *Server) SetSessionIdName(name string) {
	if s.Status() == SERVER_STATUS_RUNNING {
		glog.Error(gCHANGE_CONFIG_WHILE_RUNNING_ERROR)
		return
	}
	s.config.SessionIdName = name
}

// 获取http server参数 - SessionMaxAge
func (s *Server) GetSessionMaxAge() int64 {
	return s.config.SessionMaxAge
}

// 获取http server参数 - SessionIdName
func (s *Server) GetSessionIdName() string {
	return s.config.SessionIdName
}

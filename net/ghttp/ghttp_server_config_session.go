// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"time"

	"github.com/gogf/gf/os/gsession"
)

// 设置http server参数 - SessionMaxAge
func (s *Server) SetSessionMaxAge(ttl time.Duration) {
	s.config.SessionMaxAge = ttl
	s.sessionManager.SetTTL(ttl)
}

// 设置http server参数 - SessionIdName
func (s *Server) SetSessionIdName(name string) {
	s.config.SessionIdName = name
}

// 设置http server参数 - SessionStorage
func (s *Server) SetSessionStorage(storage gsession.Storage) {
	s.config.SessionStorage = storage
	s.sessionManager.SetStorage(storage)
}

// 获取http server参数 - SessionMaxAge
func (s *Server) GetSessionMaxAge() time.Duration {
	return s.config.SessionMaxAge
}

// 获取http server参数 - SessionIdName
func (s *Server) GetSessionIdName() string {
	return s.config.SessionIdName
}

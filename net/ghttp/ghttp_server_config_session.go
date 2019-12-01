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

func (s *Server) SetSessionMaxAge(ttl time.Duration) {
	s.config.SessionMaxAge = ttl
}

func (s *Server) SetSessionIdName(name string) {
	s.config.SessionIdName = name
}

func (s *Server) SetSessionStorage(storage gsession.Storage) {
	s.config.SessionStorage = storage
}

func (s *Server) GetSessionMaxAge() time.Duration {
	return s.config.SessionMaxAge
}

func (s *Server) GetSessionIdName() string {
	return s.config.SessionIdName
}

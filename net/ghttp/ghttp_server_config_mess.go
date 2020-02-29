// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func (s *Server) SetNameToUriType(t int) {
	s.config.NameToUriType = t
}

func (s *Server) SetDumpRouterMap(enabled bool) {
	s.config.DumpRouterMap = enabled
}

func (s *Server) SetFormParsingMemory(maxMemory int64) {
	s.config.FormParsingMemory = maxMemory
}

// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

func (s *Server) SetDenyIps(ips []string) {
	s.config.DenyIps = ips
}

func (s *Server) SetAllowIps(ips []string) {
	s.config.AllowIps = ips
}

func (s *Server) SetDenyRoutes(routes []string) {
	s.config.DenyRoutes = routes
}

// 设置URI重写规则
func (s *Server) SetRewrite(uri string, rewrite string) {
	s.config.Rewrites[uri] = rewrite
}

// 设置URI重写规则（批量）
func (s *Server) SetRewriteMap(rewrites map[string]string) {
	for k, v := range rewrites {
		s.config.Rewrites[k] = v
	}
}

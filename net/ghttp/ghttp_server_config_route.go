// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// SetRewrite sets rewrites for static URI for server.
func (s *Server) SetRewrite(uri string, rewrite string) {
	s.config.Rewrites[uri] = rewrite
}

// SetRewriteMap sets the rewrite map for server.
func (s *Server) SetRewriteMap(rewrites map[string]string) {
	for k, v := range rewrites {
		s.config.Rewrites[k] = v
	}
}

// SetRouteOverWrite sets the RouteOverWrite for server.
func (s *Server) SetRouteOverWrite(enabled bool) {
	s.config.RouteOverWrite = enabled
}

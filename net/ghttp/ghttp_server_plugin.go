// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// Plugin is the interface for server plugin.
type Plugin interface {
	Install(s *Server) error
	Remove() error
}

// Plugin adds plugin for server.
func (s *Server) Plugin(plugin ...Plugin) {
	for _, p := range plugin {
		if err := p.Install(s); err != nil {
			s.Logger().Fatal(err)
		}
	}
}

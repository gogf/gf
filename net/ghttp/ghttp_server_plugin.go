// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// Plugin is the interface for server plugin.
type Plugin interface {
	Name() string            // Name returns the name of the plugin.
	Author() string          // Author returns the author of the plugin.
	Version() string         // Version returns the version of the plugin, like "v1.0.0".
	Description() string     // Description returns the description of the plugin.
	Install(s *Server) error // Install installs the plugin before server starts.
	Remove() error           // Remove removes the plugin.
}

// Plugin adds plugin to server.
func (s *Server) Plugin(plugin ...Plugin) {
	s.plugins = append(s.plugins, plugin...)
}

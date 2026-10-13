// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Plugin registration and installed-resource cleanup for HTTP servers.

package ghttp

import "context"

// Plugin is the interface for server plugin.
type Plugin interface {
	Name() string            // Name returns the name of the plugin.
	Author() string          // Author returns the author of the plugin.
	Version() string         // Version returns the version of the plugin, like "v1.0.0".
	Description() string     // Description returns the description of the plugin.
	Install(s *Server) error // Install installs the plugin BEFORE the server starts.
	Remove() error           // Remove removes the plugin when server shuts down.
}

// Plugin adds plugin to the server.
func (s *Server) Plugin(plugin ...Plugin) {
	s.plugins = append(s.plugins, plugin...)
}

// removePlugins detaches installed plugins so shutdown never removes uninstalled ones.
func (s *Server) removePlugins(ctx context.Context) {
	s.pluginMu.Lock()
	plugins := s.installedPlugins
	s.installedPlugins = nil
	s.pluginMu.Unlock()
	for _, p := range plugins {
		s.Logger().Printf(ctx, `remove plugin: %s`, p.Name())
		if err := p.Remove(); err != nil {
			s.Logger().Errorf(ctx, "%+v", err)
		}
	}
}

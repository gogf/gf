// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

// SetNameToUriType sets the NameToUriType for server.
func (s *Server) SetNameToUriType(t int) {
	s.config.NameToUriType = t
}

// SetDumpRouterMap sets the DumpRouterMap for server.
// If DumpRouterMap is enabled, it automatically dumps the route map when server starts.
func (s *Server) SetDumpRouterMap(enabled bool) {
	s.config.DumpRouterMap = enabled
}

// SetClientMaxBodySize sets the ClientMaxBodySize for server.
func (s *Server) SetClientMaxBodySize(maxSize int64) {
	s.config.ClientMaxBodySize = maxSize
}

// SetFormParsingMemory sets the FormParsingMemory for server.
func (s *Server) SetFormParsingMemory(maxMemory int64) {
	s.config.FormParsingMemory = maxMemory
}

// SetGraceful sets the Graceful for server.
func (s *Server) SetGraceful(graceful bool) {
	s.config.Graceful = graceful
	// note: global setting.
	gracefulEnabled = graceful
}

// GetGraceful returns the Graceful for server.
func (s *Server) GetGraceful() bool {
	return s.config.Graceful
}

// SetGracefulTimeout sets the GracefulTimeout for server.
func (s *Server) SetGracefulTimeout(gracefulTimeout int) {
	s.config.GracefulTimeout = gracefulTimeout
}

// GetGracefulTimeout returns the GracefulTimeout for server.
func (s *Server) GetGracefulTimeout() int {
	return s.config.GracefulTimeout
}

// SetGracefulShutdownTimeout sets the GracefulShutdownTimeout for server.
func (s *Server) SetGracefulShutdownTimeout(gracefulShutdownTimeout int) {
	s.config.GracefulShutdownTimeout = gracefulShutdownTimeout
}

// GetGracefulShutdownTimeout returns the GracefulShutdownTimeout for server.
func (s *Server) GetGracefulShutdownTimeout() int {
	return s.config.GracefulShutdownTimeout
}

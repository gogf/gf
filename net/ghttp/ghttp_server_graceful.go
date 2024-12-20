// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/v2/net/ghttp/internal/graceful"

// newGracefulServer creates and returns a graceful http server with a given address.
// The optional parameter `fd` specifies the file descriptor which is passed from parent server.
func (s *Server) newGracefulServer(address string, fd int) *graceful.Server {
	var (
		loggerWriter = &errorLogger{logger: s.config.Logger}
		serverConfig = graceful.ServerConfig{
			Listeners:               s.config.Listeners,
			Handler:                 s.config.Handler,
			ReadTimeout:             s.config.ReadTimeout,
			WriteTimeout:            s.config.WriteTimeout,
			IdleTimeout:             s.config.IdleTimeout,
			GracefulShutdownTimeout: s.config.GracefulTimeout,
			MaxHeaderBytes:          s.config.MaxHeaderBytes,
			KeepAlive:               s.config.KeepAlive,
			Logger:                  s.config.Logger,
		}
	)
	return graceful.New(address, fd, loggerWriter, serverConfig)
}

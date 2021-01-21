// Copyright 2017 gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// +build !windows

package ghttp

import (
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gproc"
	"os"
	"syscall"
)

// registerSignalHandler handles all signal from system.
func registerSignalHandler() {
	gproc.AddSigHandler(func(sig os.Signal) {
		// Shutdown the servers with force.
		shutdownWebServers(sig.String())
	}, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT)

	gproc.AddSigHandler(func(sig os.Signal) {
		// Shutdown the servers gracefully.
		// Especially from K8S when running server in POD.
		shutdownWebServersGracefully(sig.String())
	}, syscall.SIGTERM)

	gproc.AddSigHandler(func(sig os.Signal) {
		// Restart the servers.
		if err := restartWebServers(sig.String()); err != nil {
			intlog.Error(err)
		}
	}, syscall.SIGUSR1)
}

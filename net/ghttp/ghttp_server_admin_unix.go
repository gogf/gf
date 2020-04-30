// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// +build !windows

package ghttp

import (
	"github.com/gogf/gf/internal/intlog"
	"os"
	"os/signal"
	"syscall"
)

// procSignalChan is the channel for listening the signal.
var procSignalChan = make(chan os.Signal)

// handleProcessSignal handles all signal from system.
func handleProcessSignal() {
	var sig os.Signal
	signal.Notify(
		procSignalChan,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGABRT,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	for {
		sig = <-procSignalChan
		intlog.Printf(`signal received: %s`, sig.String())
		switch sig {
		// Stop the servers.
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGABRT:
			shutdownWebServers(sig.String())
			return

		// Restart the servers.
		case syscall.SIGUSR1:
			restartWebServers(sig.String())
			return

		default:
		}
	}
}

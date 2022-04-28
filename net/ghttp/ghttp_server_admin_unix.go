// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

//go:build !windows
// +build !windows

package ghttp

import (
	"context"
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
		intlog.Printf(context.TODO(), `signal received: %s`, sig.String())
		switch sig {
		// Shutdown the servers.
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
			shutdownWebServers(sig.String())
			return

		// Shutdown the servers gracefully.
		// Especially from K8S when running server in POD.
		case syscall.SIGTERM:
			shutdownWebServersGracefully(sig.String())
			return

		// Restart the servers.
		case syscall.SIGUSR1:
			if err := restartWebServers(sig.String()); err != nil {
				intlog.Error(context.TODO(), err)
			}
			return

		default:
		}
	}
}

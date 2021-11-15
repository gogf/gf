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
	"os"
	"os/signal"
	"syscall"

	"github.com/gogf/gf/v2/internal/intlog"
)

// procSignalChan is the channel for listening the signal.
var procSignalChan = make(chan os.Signal)

// handleProcessSignal handles all signal from system.
func handleProcessSignal() {
	var (
		ctx = context.TODO()
		sig os.Signal
	)
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
		intlog.Printf(ctx, `signal received: %s`, sig.String())
		switch sig {
		// Shutdown the servers.
		case syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT:
			shutdownWebServers(ctx, sig.String())
			return

		// Shutdown the servers gracefully.
		// Especially from K8S when running server in POD.
		case syscall.SIGTERM:
			shutdownWebServersGracefully(ctx, sig.String())
			return

		// Restart the servers.
		case syscall.SIGUSR1:
			if err := restartWebServers(ctx, sig.String()); err != nil {
				intlog.Error(ctx, err)
			}
			return

		default:
		}
	}
}

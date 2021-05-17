// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"github.com/gogf/gf/internal/intlog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// SigHandler defines a function type for signal handling.
type SigHandler func(sig os.Signal)

var (
	signalHandlerMap  = make(map[os.Signal][]SigHandler)
	shutdownSignalMap = map[os.Signal]struct{}{
		syscall.SIGINT:  {},
		syscall.SIGQUIT: {},
		syscall.SIGKILL: {},
		syscall.SIGTERM: {},
		syscall.SIGABRT: {},
	}
)

func init() {
	for sig, _ := range shutdownSignalMap {
		signalHandlerMap[sig] = make([]SigHandler, 0)
	}
}

// AddSigHandler adds custom signal handler for custom one or more signals.
func AddSigHandler(handler SigHandler, signals ...os.Signal) {
	for _, sig := range signals {
		signalHandlerMap[sig] = append(signalHandlerMap[sig], handler)
	}
}

// AddSigHandlerShutdown adds custom signal handler for shutdown signals:
// syscall.SIGINT,
// syscall.SIGQUIT,
// syscall.SIGKILL,
// syscall.SIGTERM,
// syscall.SIGABRT.
func AddSigHandlerShutdown(handler SigHandler) {
	for sig, _ := range shutdownSignalMap {
		signalHandlerMap[sig] = append(signalHandlerMap[sig], handler)
	}
}

// Listen blocks and does signal listening and handling.
func Listen() {
	signals := make([]os.Signal, 0)
	for sig, _ := range signalHandlerMap {
		signals = append(signals, sig)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	var sig os.Signal
	for {
		wg := sync.WaitGroup{}
		sig = <-sigChan
		intlog.Printf(`signal received: %s`, sig.String())
		if handlers, ok := signalHandlerMap[sig]; ok {
			for _, handler := range handlers {
				wg.Add(1)
				go func(handler SigHandler, sig os.Signal) {
					defer wg.Done()
					handler(sig)
				}(handler, sig)
			}
		}
		// If it is shutdown signal, it exits this signal listening.
		if _, ok := shutdownSignalMap[sig]; ok {
			// Wait until signal handlers done.
			wg.Wait()
			return
		}
	}
}

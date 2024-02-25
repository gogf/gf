// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"context"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// SigHandler defines a function type for signal handling.
type SigHandler func(sig os.Signal)

var (
	// Use internal variable to guarantee concurrent safety
	// when multiple Listen happen.
	listenOnce        = sync.Once{}
	waitChan          = make(chan struct{})
	signalChan        = make(chan os.Signal, 1)
	signalHandlerMu   sync.Mutex
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
	for sig := range shutdownSignalMap {
		signalHandlerMap[sig] = make([]SigHandler, 0)
	}
}

// AddSigHandler adds custom signal handler for custom one or more signals.
func AddSigHandler(handler SigHandler, signals ...os.Signal) {
	signalHandlerMu.Lock()
	defer signalHandlerMu.Unlock()
	for _, sig := range signals {
		signalHandlerMap[sig] = append(signalHandlerMap[sig], handler)
	}
	notifySignals()
}

// AddSigHandlerShutdown adds custom signal handler for shutdown signals:
// syscall.SIGINT,
// syscall.SIGQUIT,
// syscall.SIGKILL,
// syscall.SIGTERM,
// syscall.SIGABRT.
func AddSigHandlerShutdown(handler ...SigHandler) {
	signalHandlerMu.Lock()
	defer signalHandlerMu.Unlock()
	for _, h := range handler {
		for sig := range shutdownSignalMap {
			signalHandlerMap[sig] = append(signalHandlerMap[sig], h)
		}
	}
	notifySignals()
}

// Listen blocks and does signal listening and handling.
func Listen() {
	listenOnce.Do(func() {
		go listen()
	})

	<-waitChan
}

func listen() {
	defer close(waitChan)

	var (
		ctx = context.Background()
		wg  = sync.WaitGroup{}
		sig os.Signal
	)
	for {
		sig = <-signalChan
		intlog.Printf(ctx, `signal received: %s`, sig.String())
		if handlers := getHandlersBySignal(sig); len(handlers) > 0 {
			for _, handler := range handlers {
				wg.Add(1)
				var (
					currentHandler = handler
					currentSig     = sig
				)
				gutil.TryCatch(ctx, func(ctx context.Context) {
					defer wg.Done()
					currentHandler(currentSig)
				}, func(ctx context.Context, exception error) {
					intlog.Errorf(ctx, `execute signal handler failed: %+v`, exception)
				})
			}
		}
		// If it is shutdown signal, it exits this signal listening.
		if _, ok := shutdownSignalMap[sig]; ok {
			intlog.Printf(
				ctx,
				`receive shutdown signal "%s", waiting all signal handler done`,
				sig.String(),
			)
			// Wait until signal handlers done.
			wg.Wait()
			intlog.Print(ctx, `all signal handler done, exit process`)
			return
		}
	}
}

func notifySignals() {
	var signals = make([]os.Signal, 0)
	for s := range signalHandlerMap {
		signals = append(signals, s)
	}
	signal.Notify(signalChan, signals...)
}

func getHandlersBySignal(sig os.Signal) []SigHandler {
	signalHandlerMu.Lock()
	defer signalHandlerMu.Unlock()
	return signalHandlerMap[sig]
}

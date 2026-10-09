// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gproc

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gutil"
)

// SigHandler defines a function type for signal handling.
type SigHandler func(sig os.Signal)

// signalListenEnded marks that the signal listening loop has exited, after which nothing
// drains signalChan any more. It is guarded by signalHandlerMu.
var signalListenEnded bool

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
		_, isShutdownSignal := shutdownSignalMap[sig]
		if isShutdownSignal {
			// This listening loop returns right after the shutdown handlers are done, so
			// from this point on nothing drains signalChan any more. Restore the default
			// behavior before running them: otherwise every signal received afterwards is
			// silently discarded and the process can only be stopped by SIGKILL, with not
			// even SIGQUIT able to dump the goroutine stacks.
			// It also restores the conventional escape hatch: a second shutdown signal
			// terminates the process even when a shutdown handler blocks.
			endSignalListening()
		}
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
		if isShutdownSignal {
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

// endSignalListening restores the default behavior for the listened signals and marks
// the listening as ended, so that handlers added afterwards cannot re-arm signal.Notify
// on a channel that no longer has a reader.
//
// It holds signalHandlerMu because notifySignals runs under that same lock: without it a
// concurrent AddSigHandler could observe signalListenEnded as false and call signal.Notify
// right after signal.Stop, silently swallowing signals again.
func endSignalListening() {
	signalHandlerMu.Lock()
	defer signalHandlerMu.Unlock()
	signalListenEnded = true
	signal.Stop(signalChan)
}

func notifySignals() {
	// The listening loop has exited and nothing drains signalChan any more. Re-arming
	// signal.Notify here would silently discard every signal received from now on.
	if signalListenEnded {
		intlog.Print(
			context.Background(),
			`signal listening has ended, newly added signal handlers will not be notified`,
		)
		return
	}
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

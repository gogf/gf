// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gapp provides application-level lifecycle management for multiple servers.
//
// It defines a unified Server interface and an App struct that coordinates
// the startup and shutdown of all registered servers, including signal handling
// for graceful shutdown.
package gapp

import (
	"context"
	"os"
	"sync"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
)

// Server is the interface for server lifecycle management.
// It defines the minimal contract for servers that can be managed by App.
type Server interface {
	// Start starts the server in non-blocking way.
	Start() error

	// Stop stops the server.
	// The parameter graceful indicates whether to stop gracefully.
	// When graceful is true, the server waits for in-flight requests to complete
	// before shutting down. When graceful is false, the server is forcibly closed.
	Stop(graceful bool) error
}

// App manages the lifecycle of multiple Server instances.
// It provides unified startup, shutdown, and signal handling
// for all registered servers.
type App struct {
	mu           sync.RWMutex
	servers      []Server
	options      []Option                    // registered options to apply during Boot
	hooks        []func(ctx context.Context) // cleanup functions collected during Boot
	booted       bool                        // whether Boot has been called successfully
	booting      bool                        // whether Boot is in progress
	lastBootErr  error                       // result of the last failed Boot for concurrent waiters
	bootCond     sync.Cond                   // coordinates concurrent Boot callers
	lifecycleCtx context.Context             // context from first successful Boot; used when Start/Stop get nil ctx
	logger       *glog.Logger
	stopOnce     sync.Once // ensures Stop is executed only once
}

// New creates and returns a new App instance with optional initial servers.
func New(servers ...Server) *App {
	app := &App{
		servers: make([]Server, 0),
		options: make([]Option, 0),
		hooks:   make([]func(ctx context.Context), 0),
		logger:  glog.New(),
	}
	app.bootCond.L = &app.mu
	if len(servers) > 0 {
		app.servers = append(app.servers, servers...)
	}
	return app
}

// shutdownModeGracefully identifies a graceful shutdown in log messages.
const shutdownModeGracefully = "gracefully"

// shutdownModeForcefully identifies a forceful shutdown in log messages.
const shutdownModeForcefully = "forcefully"

// normalizeCtx maps nil to gctx.GetInitCtx so callers use the framework default propagation context.
func normalizeCtx(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return gctx.GetInitCtx()
}

// resolveCtx returns explicit ctx when non-nil, otherwise lifecycleCtx after Boot or GetInitCtx.
func (app *App) resolveCtx(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	app.mu.RLock()
	lc := app.lifecycleCtx
	app.mu.RUnlock()
	if lc != nil {
		return lc
	}
	return gctx.GetInitCtx()
}

// Add adds one or more Server instances to the App.
func (app *App) Add(servers ...Server) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.servers = append(app.servers, servers...)
}

// Option registers one or more Options to be applied during Boot.
// Options are applied in registration order.
func (app *App) Option(opts ...Option) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.options = append(app.options, opts...)
}

// Boot applies all registered Options in registration order.
// Each Option's Apply method is called, and any returned cleanup
// functions are collected for later execution during Stop().
//
// Boot is idempotent: calling it multiple times is safe and subsequent
// calls after the first are no-ops.
//
// If an Option's Apply returns an error, Boot rolls back by running
// any already-collected cleanup hooks in reverse order and returns
// the error.
//
// If ctx is nil, gctx.GetInitCtx() is used for Apply and rollback.
func (app *App) Boot(ctx context.Context) error {
	app.mu.Lock()
	for {
		if app.booted {
			app.mu.Unlock()
			return nil
		}
		if app.booting {
			app.bootCond.Wait()
			if app.booted {
				app.mu.Unlock()
				return nil
			}
			if !app.booting {
				err := app.lastBootErr
				app.mu.Unlock()
				return err
			}
			continue
		}
		app.booting = true
		app.lastBootErr = nil
		break
	}
	baseCtx := normalizeCtx(ctx)

	// Release the lock before applying options. applyOptions reads from
	// app.options by index so that Options registered dynamically during
	// Apply() (via app.Option()) are also processed.
	app.mu.Unlock()

	// Ensure booting is always reset and waiters are unblocked,
	// even if applyOptions panics during rollback.
	var err error
	defer func() {
		app.mu.Lock()
		app.booting = false
		if err == nil {
			app.booted = true
			if app.lifecycleCtx == nil {
				app.lifecycleCtx = baseCtx
			}
		} else {
			app.lastBootErr = gerror.WrapCode(gcode.CodeInternalError, err, "app boot failed")
			err = app.lastBootErr
		}
		app.bootCond.Broadcast()
		app.mu.Unlock()
	}()

	err = app.applyOptions(baseCtx, 0)

	if err == nil {
		app.logger.Infof(baseCtx, "app booted successfully")
	}
	return err
}

// applyOptions runs registered Options and collects cleanup hooks.
// It reads from app.options by index so that Options registered
// dynamically during Apply() (via app.Option()) are also processed.
func (app *App) applyOptions(baseCtx context.Context, startIdx int) error {
	for {
		app.mu.RLock()
		if startIdx >= len(app.options) {
			app.mu.RUnlock()
			break
		}
		opt := app.options[startIdx]
		app.mu.RUnlock()

		hook, err := opt.Apply(baseCtx, app)
		if err != nil {
			// Rollback: run already-collected hooks in reverse.
			app.runHooksReverse(baseCtx)
			return err
		}
		if hook != nil {
			app.mu.Lock()
			app.hooks = append(app.hooks, hook)
			app.mu.Unlock()
		}

		startIdx++
	}
	return nil
}

// Booted returns whether Boot has been called successfully.
func (app *App) Booted() bool {
	app.mu.RLock()
	defer app.mu.RUnlock()
	return app.booted
}

// runHooksReverse atomically swaps out all collected cleanup hooks, then runs
// them in reverse order. The hooks slice is cleared before execution so that
// hooks appended concurrently (e.g. by applyOptions during Boot) are not
// accidentally discarded, and a panic in one hook does not cause double-cleanup
// on a subsequent call.
func (app *App) runHooksReverse(ctx context.Context) {
	app.mu.Lock()
	hooks := app.hooks
	app.hooks = app.hooks[:0]
	app.mu.Unlock()

	for i := len(hooks) - 1; i >= 0; i-- {
		func() {
			defer func() {
				if r := recover(); r != nil {
					app.logger.Errorf(ctx, "cleanup hook panicked: %v", r)
				}
			}()
			hooks[i](ctx)
		}()
	}
}

// rollbackStartedServers force-stops servers whose Start succeeded, in reverse registration order.
func (app *App) rollbackStartedServers(ctx context.Context, servers []Server, startOK []bool) {
	for i := len(servers) - 1; i >= 0; i-- {
		if !startOK[i] {
			continue
		}
		if err := servers[i].Stop(false); err != nil {
			app.logger.Errorf(ctx, "server rollback stop failed during start: %v", err)
		}
	}
}

// Start starts all registered servers in non-blocking way.
// It starts all servers concurrently and returns the first error encountered,
// or nil if all servers started successfully.
//
// When one or more servers fail to start after others have succeeded, each
// server that did start successfully is force-stopped in reverse registration
// order before the error is returned.
//
// When ctx becomes done before all Server.Start calls complete, servers that
// already started successfully are force-stopped in reverse order and this
// method returns ctx.Err().
//
// If Boot has not been called yet, it is called automatically first using the
// same resolved context as this call. If ctx is nil, lifecycleCtx from Boot
// or gctx.GetInitCtx() applies.
func (app *App) Start(ctx context.Context) error {
	ctx = app.resolveCtx(ctx)
	if err := app.Boot(ctx); err != nil {
		return err
	}

	app.mu.RLock()
	servers := make([]Server, len(app.servers))
	copy(servers, app.servers)
	app.mu.RUnlock()

	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		firstErr error
		startOK  = make([]bool, len(servers))
	)

	for i := range servers {
		wg.Add(1)
		go func(idx int, s Server) {
			defer wg.Done()
			if err := s.Start(); err != nil {
				mu.Lock()
				if firstErr == nil {
					firstErr = err
				}
				mu.Unlock()
				app.logger.Errorf(ctx, "server start failed: %v", err)
				return
			}
			mu.Lock()
			startOK[idx] = true
			mu.Unlock()
		}(i, servers[i])
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	var cancelErr error
	select {
	case <-done:
	case <-ctx.Done():
		wg.Wait()
		cancelErr = ctx.Err()
	}

	if cancelErr != nil {
		app.rollbackStartedServers(ctx, servers, startOK)
		return cancelErr
	}

	if firstErr != nil {
		app.rollbackStartedServers(ctx, servers, startOK)
		return gerror.WrapCode(gcode.CodeInternalError, firstErr, "app start failed")
	}
	app.logger.Infof(ctx, "all servers started successfully")
	return nil
}

// Run starts all registered servers and blocks until a shutdown signal is received
// or ctx is canceled. It handles OS signals (SIGINT, SIGTERM, etc.) for graceful
// shutdown when a signal is received. When ctx completes, graceful shutdown runs
// the same Stop path via gctx.NeverDone so propagation metadata is kept without
// canceling cleanup I/O prematurely.
//
// The lifecycle order is: Boot (apply Options) -> Start (start servers) ->
// block on signal or ctx -> Stop (cleanup hooks + stop servers).
//
// Run returns an error if Boot or Start fails, allowing the caller to decide
// whether to exit. On successful startup, Run blocks until shutdown completes
// and returns any error from graceful Stop; nil means shutdown completed cleanly.
//
// App and its registered servers follow a single lifecycle: after the first
// successful Stop, subsequent Stop calls are no-ops. Create a new App instance
// if you need a fresh lifecycle (for example in tests).
//
// If ctx is nil, gctx.GetInitCtx() resolves the root context before Boot.
func (app *App) Run(ctx context.Context) error {
	root := app.resolveCtx(ctx)
	var (
		stopOnce    sync.Once
		exitCh      = make(chan struct{})
		shutdownErr error
	)
	doShutdown := func() {
		stopOnce.Do(func() {
			shutdownErr = app.Stop(gctx.NeverDone(root), true)
			if shutdownErr != nil {
				app.logger.Errorf(root, "graceful shutdown failed: %v", shutdownErr)
			}
			close(exitCh)
		})
	}

	if err := app.Boot(root); err != nil {
		return gerror.WrapCode(gcode.CodeInternalError, err, "app boot failed")
	}

	if err := app.Start(root); err != nil {
		if stopErr := app.Stop(root, false); stopErr != nil {
			app.logger.Errorf(root, "forceful stop after start failure also failed: %v", stopErr)
		}
		return gerror.WrapCode(gcode.CodeInternalError, err, "app start failed")
	}

	gproc.AddSigHandlerShutdown(func(sig os.Signal) {
		app.logger.Infof(root, "received shutdown signal: %s, shutting down...", sig.String())
		doShutdown()
	})

	app.waitForRunExit(root, doShutdown, exitCh)
	return shutdownErr
}

// waitForRunExit blocks until shutdown completes from either context cancellation
// or an OS shutdown signal. Signal handlers are registered before this call;
// StartListen starts the background listener without blocking on waitChan so
// context-driven shutdown can return without leaving a goroutine stuck in Listen.
func (app *App) waitForRunExit(ctx context.Context, doShutdown func(), shutdownDone <-chan struct{}) {
	go func() {
		<-ctx.Done()
		doShutdown()
	}()

	gproc.StartListen()

	<-shutdownDone
}

// Stop stops all registered servers.
// The parameter graceful indicates whether to stop gracefully.
// Cleanup hooks from Options are run in reverse order first,
// then servers are stopped in reverse registration order.
//
// Stop is idempotent: calling it multiple times is safe and only the first
// call executes the shutdown logic. After Stop completes, the App cannot
// shut down servers again; register servers and call Start on a new App for
// another lifecycle round.
//
// If ctx is nil, lifecycleCtx or gctx.GetInitCtx() is used for hooks and logs.
func (app *App) Stop(ctx context.Context, graceful bool) error {
	var firstErr error
	app.stopOnce.Do(func() {
		opCtx := app.resolveCtx(ctx)

		// Run cleanup hooks in reverse order first.
		app.runHooksReverse(opCtx)

		// Copy server list under the lock, then release before stopping
		// to avoid holding the read lock during potentially long server shutdown.
		app.mu.RLock()
		servers := make([]Server, len(app.servers))
		copy(servers, app.servers)
		app.mu.RUnlock()

		if len(servers) == 0 {
			return
		}

		// Stop in reverse order.
		for i := len(servers) - 1; i >= 0; i-- {
			server := servers[i]
			if err := server.Stop(graceful); err != nil {
				if firstErr == nil {
					firstErr = err
				}
				app.logger.Errorf(opCtx, "server stop failed: %v", err)
			}
		}

		mode := shutdownModeGracefully
		if !graceful {
			mode = shutdownModeForcefully
		}
		app.logger.Infof(opCtx, "all servers stopped %s", mode)
	})
	if firstErr != nil {
		return gerror.WrapCode(gcode.CodeInternalError, firstErr, "app stop failed")
	}
	return nil
}

// Servers returns a copy of the registered servers list.
func (app *App) Servers() []Server {
	app.mu.RLock()
	defer app.mu.RUnlock()

	servers := make([]Server, len(app.servers))
	copy(servers, app.servers)
	return servers
}

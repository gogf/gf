// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// WorkerServer implementation for long-running background tasks.

package gjob

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gapp"
)

// Compile-time check that WorkerServer implements gapp.Server.
var _ gapp.Server = (*WorkerServer)(nil)

// WorkerHandler is the handler function for a background worker task.
// It receives a context that is cancelled when the server stops,
// and returns an optional cleanup function that is called on shutdown.
type WorkerHandler func(ctx context.Context) func()

// WorkerTask defines a background worker task that runs in its own goroutine.
// When the server stops, the task's context is cancelled and the
// handler's cleanup function (if any) is called.
type WorkerTask struct {
	// Name is the unique name of the worker task.
	Name string

	// Handler is the function that implements the worker logic.
	Handler WorkerHandler
}

// WorkerServer manages a set of background worker tasks that implement the gapp.Server interface.
// Each task runs in its own goroutine and is terminated when the server stops.
//
// WorkerServer follows a single lifecycle: after Stop, Start returns CodeInvalidOperation.
// Create a new WorkerServer for another lifecycle round.
// When Stop is called, the internal context is cancelled to signal all tasks
// to stop. The server always waits for task goroutines to finish before
// returning so lifecycle managers (including gapp rollback) do not proceed early.
type WorkerServer struct {
	mu        sync.Mutex
	parentCtx context.Context
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	tasks     []WorkerTask
	started   bool
	stopped   bool
}

// NewWorkerServer creates and returns a new WorkerServer with the given tasks.
// The ctx parameter defines the parent lifecycle context for worker tasks; a nil
// value is normalized to gctx.GetInitCtx().
// The context.WithCancel is deferred until Start is called to avoid leaking
// resources if the server is created but never started.
func NewWorkerServer(ctx context.Context, tasks ...WorkerTask) *WorkerServer {
	return &WorkerServer{
		parentCtx: normalizeCtx(ctx),
		tasks:     tasks,
	}
}

// Add adds one or more worker tasks to the server.
// Tasks must be added before Start is called.
// Returns CodeInvalidOperation when the server has already started or stopped.
func (s *WorkerServer) Add(tasks ...WorkerTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return gerror.NewCode(gcode.CodeInvalidOperation, "cannot add worker tasks after server started")
	}
	if s.stopped {
		return gerror.NewCode(gcode.CodeInvalidOperation, "cannot add worker tasks after server stopped")
	}
	s.tasks = append(s.tasks, tasks...)
	return nil
}

// Start starts all registered worker tasks in non-blocking way.
// Each task is launched in its own goroutine with panic recovery.
func (s *WorkerServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return gerror.NewCode(gcode.CodeInvalidOperation, "worker server already started")
	}
	if s.stopped {
		return gerror.NewCode(gcode.CodeInvalidOperation, "worker server already stopped")
	}

	s.ctx, s.cancel = context.WithCancel(s.parentCtx)

	for i := range s.tasks {
		task := s.tasks[i]
		s.wg.Add(1)
		go s.runTask(task)
	}

	s.started = true
	return nil
}

// Stop stops the worker server by cancelling the internal context and waiting
// for all task goroutines to finish. The graceful parameter is accepted for
// gapp.Server compatibility; worker shutdown always waits after cancellation.
func (s *WorkerServer) Stop(_ bool) error {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return nil
	}
	s.started = false
	s.stopped = true
	s.mu.Unlock()

	s.cancel()
	s.wg.Wait()

	return nil
}

// runTask runs a single worker task in a goroutine with panic recovery.
func (s *WorkerServer) runTask(task WorkerTask) {
	defer s.wg.Done()

	if err := g.Try(s.ctx, func(ctx context.Context) {
		cleanup := task.Handler(ctx)
		if cleanup != nil {
			defer cleanup()
		}
		<-ctx.Done()
	}); err != nil {
		g.Log().Errorf(s.ctx, "worker task %s exited with error: %v", task.Name, err)
	}
}

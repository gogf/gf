// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// CronServer implementation for scheduled cron tasks.

package gjob

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/os/gcron"
)

// CronHandler is the handler function for a cron task.
type CronHandler func(ctx context.Context) error

// CronTask defines a cron task that runs on a specified schedule.
type CronTask struct {
	// Name is the unique name of the cron task.
	Name string

	// Spec is the cron expression that defines the schedule.
	// It supports second-level precision, e.g. "*/2 * * * * *".
	Spec string

	// Handler is the function that is called on each schedule tick.
	Handler CronHandler
}

// CronServer manages a set of cron tasks that implement the gapp.Server interface.
// Tasks are registered with the internal gcron.Cron scheduler as singleton jobs
// and run on their specified schedules.
//
// CronServer follows a single lifecycle: after Stop, Start returns CodeInvalidOperation.
// Create a new CronServer for another lifecycle round.
// When Stop is called, the scheduler is stopped and no more tasks are executed.
type CronServer struct {
	mu        sync.Mutex
	parentCtx context.Context
	ctx       context.Context
	cancel    context.CancelFunc
	cron      *gcron.Cron
	tasks     []CronTask
	started   bool
	stopped   bool
}

// NewCronServer creates and returns a new CronServer with the given tasks.
// The ctx parameter defines the parent lifecycle context for cron jobs; a nil
// value is normalized to gctx.GetInitCtx().
// The context.WithCancel is deferred until Start is called to avoid leaking
// resources if the server is created but never started.
func NewCronServer(ctx context.Context, tasks ...CronTask) *CronServer {
	return &CronServer{
		parentCtx: normalizeCtx(ctx),
		tasks:     tasks,
	}
}

// Add adds one or more cron tasks to the server.
// Tasks must be added before Start is called.
// Returns CodeInvalidOperation when the server has already started or stopped.
func (s *CronServer) Add(tasks ...CronTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return gerror.NewCode(gcode.CodeInvalidOperation, "cannot add cron tasks after server started")
	}
	if s.stopped {
		return gerror.NewCode(gcode.CodeInvalidOperation, "cannot add cron tasks after server stopped")
	}
	s.tasks = append(s.tasks, tasks...)
	return nil
}

// Start starts the cron scheduler and registers all tasks as singleton jobs.
// Registration uses a fresh cron instance so a partial failure does not leave
// the server in an unrecoverable state; callers may retry Start after fixing tasks.
func (s *CronServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return gerror.NewCode(gcode.CodeInvalidOperation, "cron server already started")
	}
	if s.stopped {
		return gerror.NewCode(gcode.CodeInvalidOperation, "cron server already stopped")
	}

	s.ctx, s.cancel = context.WithCancel(s.parentCtx)

	cron := gcron.New()
	for i := range s.tasks {
		task := s.tasks[i]
		handler := func(ctx context.Context) {
			if err := task.Handler(ctx); err != nil {
				g.Log().Errorf(ctx, "cron task %s handle error: %v", task.Name, err)
			}
		}

		_, err := cron.AddSingleton(s.ctx, task.Spec, handler, task.Name)
		if err != nil {
			cron.Close()
			s.cancel()
			return err
		}
	}

	cron.Start()
	s.cron = cron
	s.started = true
	return nil
}

// Stop stops the cron scheduler.
// When graceful is true, it waits for in-flight jobs to finish before returning.
// When graceful is false, running jobs are stopped immediately.
func (s *CronServer) Stop(graceful bool) error {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return nil
	}
	s.started = false
	s.stopped = true
	s.mu.Unlock()

	if graceful {
		s.cron.StopGracefully()
	} else {
		s.cron.Stop()
	}
	s.cron.Close()
	s.cancel()
	return nil
}

// Compile-time check that CronServer implements gapp.Server.
var _ gapp.Server = (*CronServer)(nil)

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjob_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	gjob "github.com/gogf/gf/contrib/scheduler/gjob/v2"
)

func TestWorkerServerImplementsGappServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var _ gapp.Server = gjob.NewWorkerServer(nil)
	})
}

func TestNewWorkerServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewWorkerServer(nil)
		t.AssertNE(s, nil)
	})
}

func TestNewWorkerServerWithTasks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{Name: "task1", Handler: func(ctx context.Context) func() { return nil }},
		)
		t.AssertNE(s, nil)
	})
}

func TestWorkerServerAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewWorkerServer(nil)
		err := s.Add(
			gjob.WorkerTask{Name: "task1", Handler: func(ctx context.Context) func() { return nil }},
			gjob.WorkerTask{Name: "task2", Handler: func(ctx context.Context) func() { return nil }},
		)
		t.AssertNil(err)
	})
}

func TestWorkerServerAddAfterStart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{Name: "task1", Handler: func(ctx context.Context) func() { return nil }},
		)
		err := s.Start()
		t.AssertNil(err)

		err = s.Add(gjob.WorkerTask{Name: "task2", Handler: func(ctx context.Context) func() { return nil }})
		t.AssertNE(err, nil)

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestWorkerServerStartStop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			started int32
			cleaned int32
		)

		s := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{
				Name: "test-worker",
				Handler: func(ctx context.Context) func() {
					atomic.StoreInt32(&started, 1)
					return func() {
						atomic.StoreInt32(&cleaned, 1)
					}
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		// Wait for the handler to be called.
		time.Sleep(100 * time.Millisecond)
		t.Assert(atomic.LoadInt32(&started), 1)

		err = s.Stop(true)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&cleaned), 1)
	})
}

func TestWorkerServerStopForceful(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			started int32
			exited  int32
		)

		s := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{
				Name: "test-worker",
				Handler: func(ctx context.Context) func() {
					atomic.StoreInt32(&started, 1)
					<-ctx.Done()
					atomic.StoreInt32(&exited, 1)
					return nil
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		t.Assert(atomic.LoadInt32(&started), 1)

		err = s.Stop(false)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&exited), 1)
	})
}

func TestWorkerServerStartTwice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewWorkerServer(nil)
		err := s.Start()
		t.AssertNil(err)
		err = s.Start()
		t.AssertNE(err, nil)
	})
}

func TestWorkerServerMultipleTasks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu      sync.Mutex
			started []string
			cleaned []string
		)

		makeHandler := func(name string) gjob.WorkerHandler {
			return func(ctx context.Context) func() {
				mu.Lock()
				started = append(started, name)
				mu.Unlock()
				return func() {
					mu.Lock()
					cleaned = append(cleaned, name)
					mu.Unlock()
				}
			}
		}

		s := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{Name: "task1", Handler: makeHandler("task1")},
			gjob.WorkerTask{Name: "task2", Handler: makeHandler("task2")},
		)

		err := s.Start()
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)

		err = s.Stop(true)
		t.AssertNil(err)

		mu.Lock()
		t.Assert(len(started), 2)
		t.Assert(len(cleaned), 2)
		mu.Unlock()
	})
}

func TestWorkerServerWithContextCancellation(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			handlerCalled int32
			cleaned       int32
		)

		s := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{
				Name: "test-ctx",
				Handler: func(ctx context.Context) func() {
					atomic.StoreInt32(&handlerCalled, 1)
					return func() {
						atomic.StoreInt32(&cleaned, 1)
					}
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		time.Sleep(100 * time.Millisecond)
		t.Assert(atomic.LoadInt32(&handlerCalled), 1)

		err = s.Stop(true)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&cleaned), 1)
	})
}

func TestWorkerServerPropagatesLifecycleContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		traceCtx := gctx.WithSpan(context.Background(), "worker-lifecycle")
		traceID := gctx.CtxId(traceCtx)

		var (
			mu    sync.Mutex
			gotID string
		)
		done := make(chan struct{})
		s := gjob.NewWorkerServer(traceCtx, gjob.WorkerTask{
			Name: "ctx-worker",
			Handler: func(ctx context.Context) func() {
				mu.Lock()
				gotID = gctx.CtxId(ctx)
				mu.Unlock()
				close(done)
				return nil
			},
		})

		err := s.Start()
		t.AssertNil(err)

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("worker handler was not invoked")
		}

		mu.Lock()
		t.Assert(gotID, traceID)
		mu.Unlock()

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerPropagatesLifecycleContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		traceCtx := gctx.WithSpan(context.Background(), "cron-lifecycle")
		traceID := gctx.CtxId(traceCtx)

		var (
			mu    sync.Mutex
			gotID string
		)
		done := make(chan struct{})
		s := gjob.NewCronServer(traceCtx, gjob.CronTask{
			Name: "ctx-cron",
			Spec: "*/1 * * * * *",
			Handler: func(ctx context.Context) error {
				mu.Lock()
				if gotID == "" {
					gotID = gctx.CtxId(ctx)
					close(done)
				}
				mu.Unlock()
				return nil
			},
		})

		err := s.Start()
		t.AssertNil(err)

		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("cron handler was not invoked")
		}

		mu.Lock()
		t.Assert(gotID, traceID)
		mu.Unlock()

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerImplementsGappServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var _ gapp.Server = gjob.NewCronServer(nil)
	})
}

func TestWorkerServerStopBeforeStart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewWorkerServer(nil)
		err := s.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerStopGracefulWaitsInFlightJob(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			running int32
			done    int32
		)

		s := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name: "slow-cron",
				Spec: "*/1 * * * * *",
				Handler: func(ctx context.Context) error {
					atomic.StoreInt32(&running, 1)
					time.Sleep(300 * time.Millisecond)
					atomic.StoreInt32(&done, 1)
					return nil
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		for i := 0; i < 50 && atomic.LoadInt32(&running) == 0; i++ {
			time.Sleep(20 * time.Millisecond)
		}
		t.Assert(atomic.LoadInt32(&running), int32(1))
		t.Assert(atomic.LoadInt32(&done), int32(0))

		err = s.Stop(true)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&done), int32(1))
	})
}

func TestNewCronServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewCronServer(nil)
		t.AssertNE(s, nil)
	})
}

func TestCronServerAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewCronServer(nil)
		err := s.Add(
			gjob.CronTask{Name: "task1", Spec: "*/1 * * * * *", Handler: func(ctx context.Context) error { return nil }},
		)
		t.AssertNil(err)
	})
}

func TestCronServerAddAfterStart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewCronServer(nil,
			gjob.CronTask{Name: "task1", Spec: "*/1 * * * * *", Handler: func(ctx context.Context) error { return nil }},
		)
		err := s.Start()
		t.AssertNil(err)

		err = s.Add(gjob.CronTask{Name: "task2", Spec: "*/1 * * * * *", Handler: func(ctx context.Context) error { return nil }})
		t.AssertNE(err, nil)

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerStartStop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var called int32

		s := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name: "test-cron",
				Spec: "*/1 * * * * *",
				Handler: func(ctx context.Context) error {
					atomic.AddInt32(&called, 1)
					return nil
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		// Wait for at least one tick.
		time.Sleep(1500 * time.Millisecond)
		t.Assert(atomic.LoadInt32(&called) >= 1, true)

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerStartTwice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewCronServer(nil,
			gjob.CronTask{Name: "t", Spec: "*/1 * * * * *", Handler: func(ctx context.Context) error { return nil }},
		)

		err := s.Start()
		t.AssertNil(err)
		err = s.Start()
		t.AssertNE(err, nil)

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerStopForceful(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var called int32

		s := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name: "test-cron-force",
				Spec: "*/1 * * * * *",
				Handler: func(ctx context.Context) error {
					atomic.AddInt32(&called, 1)
					return nil
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		time.Sleep(1500 * time.Millisecond)

		err = s.Stop(false)
		t.AssertNil(err)
	})
}

func TestCronServerInvalidSpec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name:    "invalid",
				Spec:    "invalid-spec",
				Handler: func(ctx context.Context) error { return nil },
			},
		)

		err := s.Start()
		t.AssertNE(err, nil)
	})
}

func TestCronServerStartPartialFailureRetry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name:    "good",
				Spec:    "*/1 * * * * *",
				Handler: func(ctx context.Context) error { return nil },
			},
			gjob.CronTask{
				Name:    "bad",
				Spec:    "invalid-spec",
				Handler: func(ctx context.Context) error { return nil },
			},
		)

		err := s.Start()
		t.AssertNE(err, nil)

		// A second Start must fail with the same validation error, not duplicate-name.
		err = s.Start()
		t.AssertNE(err, nil)
		t.AssertIN("invalid pattern", err.Error())

		// After fixing tasks, Start on a fresh server succeeds.
		fixed := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name:    "good",
				Spec:    "*/1 * * * * *",
				Handler: func(ctx context.Context) error { return nil },
			},
		)
		err = fixed.Start()
		t.AssertNil(err)
		err = fixed.Stop(true)
		t.AssertNil(err)
	})
}

func TestCronServerHandlerError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var called int32

		s := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name: "error-cron",
				Spec: "*/1 * * * * *",
				Handler: func(ctx context.Context) error {
					atomic.AddInt32(&called, 1)
					return context.DeadlineExceeded
				},
			},
		)

		err := s.Start()
		t.AssertNil(err)

		// Wait for at least one tick, handler should not crash the server.
		time.Sleep(1500 * time.Millisecond)
		t.Assert(atomic.LoadInt32(&called) >= 1, true)

		err = s.Stop(true)
		t.AssertNil(err)
	})
}

func TestWorkerAndCronWithGapp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			workerStarted int32
			cronCalled    int32
		)

		workerSrv := gjob.NewWorkerServer(nil,
			gjob.WorkerTask{
				Name: "app-worker",
				Handler: func(ctx context.Context) func() {
					atomic.StoreInt32(&workerStarted, 1)
					return nil
				},
			},
		)

		cronSrv := gjob.NewCronServer(nil,
			gjob.CronTask{
				Name: "app-cron",
				Spec: "*/1 * * * * *",
				Handler: func(ctx context.Context) error {
					atomic.AddInt32(&cronCalled, 1)
					return nil
				},
			},
		)

		app := gapp.New(workerSrv, cronSrv)

		err := app.Start(context.Background())
		t.AssertNil(err)

		time.Sleep(1500 * time.Millisecond)

		t.Assert(atomic.LoadInt32(&workerStarted), 1)
		t.Assert(atomic.LoadInt32(&cronCalled) >= 1, true)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
	})
}

const testJobConfig = `
scheduler:
  job:
    - name: worker1
      type: worker
      enable: true
    - name: cron1
      type: cron
      enable: true
      spec: "*/1 * * * * *"
    - name: disabled-worker
      type: worker
      enable: false
    - name: no-type
      enable: true
    - name: invalid-type
      type: unknown
      enable: true
`

func TestNewServersFromConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Set up test configuration.
		adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
		if !ok {
			t.Fatal("expected gcfg.AdapterFile")
		}
		adapter.SetContent(testJobConfig)
		defer adapter.SetContent("")

		var (
			workerStarted int32
			cronCalled    int32
		)

		handlers := gjob.HandlerMap{
			"worker1": gjob.WorkerHandler(func(ctx context.Context) func() {
				atomic.StoreInt32(&workerStarted, 1)
				return nil
			}),
			"cron1": gjob.CronHandler(func(ctx context.Context) error {
				atomic.AddInt32(&cronCalled, 1)
				return nil
			}),
		}

		servers := gjob.NewServersFromConfig(context.Background(), handlers)
		// Should create 2 servers: one WorkerServer and one CronServer.
		t.Assert(len(servers), 2)

		app := gapp.New(servers...)
		err := app.Start(context.Background())
		t.AssertNil(err)

		time.Sleep(1500 * time.Millisecond)

		t.Assert(atomic.LoadInt32(&workerStarted), 1)
		t.Assert(atomic.LoadInt32(&cronCalled) >= 1, true)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
	})
}

func TestNewServersFromConfigNoHandlers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
		if !ok {
			t.Fatal("expected gcfg.AdapterFile")
		}
		adapter.SetContent(testJobConfig)
		defer adapter.SetContent("")

		// No handlers provided.
		servers := gjob.NewServersFromConfig(context.Background(), gjob.HandlerMap{})
		t.Assert(len(servers), 0)
	})
}

func TestNewServersFromConfigWrongHandlerType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
		if !ok {
			t.Fatal("expected gcfg.AdapterFile")
		}
		adapter.SetContent(testJobConfig)
		defer adapter.SetContent("")

		// Provide wrong handler types.
		handlers := gjob.HandlerMap{
			"worker1": func() {}, // Not a WorkerHandler
			"cron1":   func() {}, // Not a CronHandler
		}

		servers := gjob.NewServersFromConfig(context.Background(), handlers)
		// All handlers have wrong types, so no servers should be created.
		t.Assert(len(servers), 0)
	})
}

func TestNewServersFromConfigEmptyConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
		if !ok {
			t.Fatal("expected gcfg.AdapterFile")
		}
		adapter.SetContent("")
		defer adapter.SetContent("")

		servers := gjob.NewServersFromConfig(context.Background(), gjob.HandlerMap{
			"worker1": gjob.WorkerHandler(func(ctx context.Context) func() { return nil }),
		})
		t.Assert(len(servers), 0)
	})
}

func TestNewServersFromConfigCronMissingSpec(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		const cronMissingSpecConfig = `
scheduler:
  job:
    - name: cron1
      type: cron
      enable: true
`
		adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
		if !ok {
			t.Fatal("expected gcfg.AdapterFile")
		}
		adapter.SetContent(cronMissingSpecConfig)
		defer adapter.SetContent("")

		handlers := gjob.HandlerMap{
			"cron1": gjob.CronHandler(func(ctx context.Context) error { return nil }),
		}

		servers := gjob.NewServersFromConfig(context.Background(), handlers)
		t.Assert(len(servers), 0)
	})
}

func TestNewServersFromConfigOnlyCron(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		const cronOnlyConfig = `
scheduler:
  job:
    - name: cron1
      type: cron
      enable: true
      spec: "*/1 * * * * *"
`
		adapter, ok := g.Cfg().GetAdapter().(*gcfg.AdapterFile)
		if !ok {
			t.Fatal("expected gcfg.AdapterFile")
		}
		adapter.SetContent(cronOnlyConfig)
		defer adapter.SetContent("")

		handlers := gjob.HandlerMap{
			"cron1": gjob.CronHandler(func(ctx context.Context) error { return nil }),
		}

		servers := gjob.NewServersFromConfig(context.Background(), handlers)
		// Only one CronServer should be created.
		t.Assert(len(servers), 1)
	})
}

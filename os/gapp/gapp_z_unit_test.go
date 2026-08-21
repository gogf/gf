// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gapp_test

import (
	"context"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/os/gapp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
)

// mockServer is a mock Server implementation for testing.
type mockServer struct {
	mu          sync.Mutex
	started     bool
	stopped     bool
	gracefulVal bool
	startErr    error
	stopErr     error
	name        string
	startDelay  time.Duration
	stopFunc    func(graceful bool) error
}

func (m *mockServer) Start() error {
	if m.startDelay > 0 {
		time.Sleep(m.startDelay)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.started = true
	return m.startErr
}

func (m *mockServer) Stop(graceful bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopped = true
	m.gracefulVal = graceful
	if m.stopFunc != nil {
		return m.stopFunc(graceful)
	}
	return m.stopErr
}

func TestNew(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := gapp.New()
		t.AssertNE(app, nil)
		t.Assert(len(app.Servers()), 0)
	})
}

func TestNewWithServers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1"}
		s2 := &mockServer{name: "s2"}
		app := gapp.New(s1, s2)
		t.Assert(len(app.Servers()), 2)
	})
}

func TestAdd(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1"}
		app := gapp.New(s1)
		t.Assert(len(app.Servers()), 1)

		s2 := &mockServer{name: "s2"}
		s3 := &mockServer{name: "s3"}
		app.Add(s2, s3)
		t.Assert(len(app.Servers()), 3)
	})
}

func TestAppStartStop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1"}
		s2 := &mockServer{name: "s2"}
		app := gapp.New(s1, s2)

		err := app.Start(context.Background())
		t.AssertNil(err)
		t.Assert(s1.started, true)
		t.Assert(s2.started, true)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
		t.Assert(s1.stopped, true)
		t.Assert(s2.stopped, true)
		t.Assert(s1.gracefulVal, true)
		t.Assert(s2.gracefulVal, true)
	})
}

func TestAppStopForceful(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1"}
		app := gapp.New(s1)

		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), false)
		t.AssertNil(err)
		t.Assert(s1.stopped, true)
		t.Assert(s1.gracefulVal, false)
	})
}

func TestAppStopReverseOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu      sync.Mutex
			stopSeq []string
		)

		s1 := &mockServer{name: "s1", stopFunc: func(_ bool) error {
			mu.Lock()
			stopSeq = append(stopSeq, "s1")
			mu.Unlock()
			return nil
		}}
		s2 := &mockServer{name: "s2", stopFunc: func(_ bool) error {
			mu.Lock()
			stopSeq = append(stopSeq, "s2")
			mu.Unlock()
			return nil
		}}
		s3 := &mockServer{name: "s3", stopFunc: func(_ bool) error {
			mu.Lock()
			stopSeq = append(stopSeq, "s3")
			mu.Unlock()
			return nil
		}}

		app := gapp.New(s1, s2, s3)
		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)

		mu.Lock()
		t.Assert(len(stopSeq), 3)
		t.Assert(stopSeq[0], "s3")
		t.Assert(stopSeq[1], "s2")
		t.Assert(stopSeq[2], "s1")
		mu.Unlock()
	})
}

func TestAppStartError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1", startErr: errTestStart}
		app := gapp.New(s1)

		err := app.Start(context.Background())
		t.AssertNE(err, nil)
	})
}

func TestAppStartPartialFailureRollback(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1"}
		s2 := &mockServer{name: "s2", startErr: errTestStart}
		app := gapp.New(s1, s2)

		err := app.Start(context.Background())
		t.AssertNE(err, nil)
		t.Assert(s1.started, true)
		t.Assert(s2.started, true)
		t.Assert(s1.stopped, true)
		t.Assert(s1.gracefulVal, false)
	})
}

func TestAppStopError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1", stopErr: errTestStop}
		app := gapp.New(s1)

		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNE(err, nil)
	})
}

func TestAppStopEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := gapp.New()
		err := app.Stop(context.Background(), true)
		t.AssertNil(err)
	})
}

func TestHTTPServerAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := ghttp.GetServer("gapp-test-http")
		s.SetPort(0)
		s.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write("ok")
		})

		adapter := gapp.NewHTTPServerAdapter(s)
		err := adapter.Start()
		t.AssertNil(err)

		time.Sleep(time.Millisecond * 100)

		err = adapter.Stop(true)
		t.AssertNil(err)
	})
}

func TestHTTPServerAdapterForceful(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := ghttp.GetServer("gapp-test-http-force")
		s.SetPort(0)
		s.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write("ok")
		})

		adapter := gapp.NewHTTPServerAdapter(s)
		err := adapter.Start()
		t.AssertNil(err)

		time.Sleep(time.Millisecond * 100)

		err = adapter.Stop(false)
		t.AssertNil(err)
	})
}

func TestTCPServerAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gtcp.NewServer(":0", func(conn *gtcp.Conn) {
			defer conn.Close()
			for {
				data, err := conn.Recv(-1)
				if err != nil {
					break
				}
				conn.Send(data)
			}
		})

		adapter := gapp.NewTCPServerAdapter(s)
		err := adapter.Start()
		t.AssertNil(err)

		err = adapter.Stop(true)
		t.AssertNil(err)
	})
}

func TestUDPServerAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := gudp.NewServer(":0", func(conn *gudp.ServerConn) {
			defer conn.Close()
			for {
				data, remote, err := conn.Recv(-1)
				if err != nil {
					if err != io.EOF {
						break
					}
					break
				}
				if err = conn.Send(data, remote); err != nil {
					break
				}
			}
		})

		adapter := gapp.NewUDPServerAdapter(s)
		err := adapter.Start()
		t.AssertNil(err)

		err = adapter.Stop(true)
		t.AssertNil(err)
	})
}

func TestAppWithHTTPServer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := ghttp.GetServer("gapp-test-app-http")
		s.SetPort(0)
		s.BindHandler("/", func(r *ghttp.Request) {
			r.Response.Write("ok")
		})

		adapter := gapp.NewHTTPServerAdapter(s)
		app := gapp.New(adapter)

		err := app.Start(context.Background())
		t.AssertNil(err)

		time.Sleep(time.Millisecond * 100)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
	})
}

func TestNewOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var called bool
		opt := gapp.NewOption(func(ctx context.Context, app *gapp.App) {
			called = true
		})
		app := gapp.New()
		hook, err := opt.Apply(context.TODO(), app)
		t.AssertNil(err)
		t.Assert(called, true)
		t.Assert(hook, nil)
	})
}

func TestNewOptionWithHook(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			initCalled    bool
			cleanupCalled bool
		)
		opt := gapp.NewOptionWithHook(func(ctx context.Context, app *gapp.App) (func(ctx context.Context), error) {
			initCalled = true
			return func(ctx context.Context) {
				cleanupCalled = true
			}, nil
		})
		app := gapp.New()
		hook, err := opt.Apply(context.TODO(), app)
		t.AssertNil(err)
		t.Assert(initCalled, true)
		t.AssertNE(hook, nil)

		// Call the cleanup hook.
		hook(context.TODO())
		t.Assert(cleanupCalled, true)
	})
}

func TestNewOptionWithHookNilCleanup(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var initCalled bool
		opt := gapp.NewOptionWithHook(func(ctx context.Context, app *gapp.App) (func(ctx context.Context), error) {
			initCalled = true
			return nil, nil
		})
		app := gapp.New()
		hook, err := opt.Apply(context.TODO(), app)
		t.AssertNil(err)
		t.Assert(initCalled, true)
		t.Assert(hook, nil)
	})
}

func TestAppOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var called bool
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
			called = true
		}))
		err := app.Boot(context.TODO())
		t.AssertNil(err)
		t.Assert(called, true)
	})
}

func TestAppBoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu    sync.Mutex
			order []string
		)
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
			mu.Lock()
			order = append(order, "first")
			mu.Unlock()
		}))
		app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
			mu.Lock()
			order = append(order, "second")
			mu.Unlock()
		}))

		err := app.Boot(context.TODO())
		t.AssertNil(err)

		mu.Lock()
		t.Assert(len(order), 2)
		t.Assert(order[0], "first")
		t.Assert(order[1], "second")
		mu.Unlock()
	})
}

func TestAppBootIdempotent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var callCount int
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
			callCount++
		}))

		err := app.Boot(context.TODO())
		t.AssertNil(err)
		t.Assert(callCount, 1)

		// Second call should be a no-op.
		err = app.Boot(context.TODO())
		t.AssertNil(err)
		t.Assert(callCount, 1)
	})
}

func TestAppBootError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := gapp.New()
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			return nil, errTestBoot
		}))

		err := app.Boot(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(app.Booted(), false)
	})
}

func TestAppBootRollbackOnError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu           sync.Mutex
			cleanupOrder []string
		)
		app := gapp.New()

		// First option succeeds and registers a cleanup.
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			return func(ctx context.Context) {
				mu.Lock()
				cleanupOrder = append(cleanupOrder, "first")
				mu.Unlock()
			}, nil
		}))

		// Second option fails, triggering rollback.
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			return nil, errTestBoot
		}))

		err := app.Boot(context.TODO())
		t.AssertNE(err, nil)

		mu.Lock()
		t.Assert(len(cleanupOrder), 1)
		t.Assert(cleanupOrder[0], "first")
		mu.Unlock()
	})
}

func TestAppStartAutoBoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var bootCalled bool
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
			bootCalled = true
		}))

		// Start without explicitly calling Boot.
		t.Assert(app.Booted(), false)
		err := app.Start(context.Background())
		t.AssertNil(err)
		t.Assert(bootCalled, true)
		t.Assert(app.Booted(), true)
	})
}

func TestAppStopRunsHooks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu           sync.Mutex
			cleanupOrder []string
		)
		s1 := &mockServer{name: "s1"}
		app := gapp.New(s1)

		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			return func(ctx context.Context) {
				mu.Lock()
				cleanupOrder = append(cleanupOrder, "hook1")
				mu.Unlock()
			}, nil
		}))
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			return func(ctx context.Context) {
				mu.Lock()
				cleanupOrder = append(cleanupOrder, "hook2")
				mu.Unlock()
			}, nil
		}))

		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)

		mu.Lock()
		// Hooks run in reverse order.
		t.Assert(len(cleanupOrder), 2)
		t.Assert(cleanupOrder[0], "hook2")
		t.Assert(cleanupOrder[1], "hook1")
		mu.Unlock()
	})
}

func TestMultipleOptionsOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu           sync.Mutex
			applyOrder   []string
			cleanupOrder []string
		)
		app := gapp.New()

		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			mu.Lock()
			applyOrder = append(applyOrder, "first")
			mu.Unlock()
			return func(ctx context.Context) {
				mu.Lock()
				cleanupOrder = append(cleanupOrder, "first")
				mu.Unlock()
			}, nil
		}))
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			mu.Lock()
			applyOrder = append(applyOrder, "second")
			mu.Unlock()
			return func(ctx context.Context) {
				mu.Lock()
				cleanupOrder = append(cleanupOrder, "second")
				mu.Unlock()
			}, nil
		}))
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			mu.Lock()
			applyOrder = append(applyOrder, "third")
			mu.Unlock()
			return func(ctx context.Context) {
				mu.Lock()
				cleanupOrder = append(cleanupOrder, "third")
				mu.Unlock()
			}, nil
		}))

		s1 := &mockServer{name: "s1"}
		app.Add(s1)

		err := app.Start(context.Background())
		t.AssertNil(err)

		mu.Lock()
		t.Assert(len(applyOrder), 3)
		t.Assert(applyOrder[0], "first")
		t.Assert(applyOrder[1], "second")
		t.Assert(applyOrder[2], "third")
		mu.Unlock()

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)

		mu.Lock()
		// Cleanups run in reverse order.
		t.Assert(len(cleanupOrder), 3)
		t.Assert(cleanupOrder[0], "third")
		t.Assert(cleanupOrder[1], "second")
		t.Assert(cleanupOrder[2], "first")
		mu.Unlock()
	})
}

func TestOptionCanAddServers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, a *gapp.App) {
			s := &mockServer{name: "dynamic"}
			a.Add(s)
		}))

		t.Assert(len(app.Servers()), 0)
		err := app.Boot(context.TODO())
		t.AssertNil(err)
		t.Assert(len(app.Servers()), 1)
	})
}

func TestBootedAccessor(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := gapp.New()
		t.Assert(app.Booted(), false)

		err := app.Boot(context.TODO())
		t.AssertNil(err)
		t.Assert(app.Booted(), true)
	})
}

func TestAppBootHookRunsBeforeServerStop(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu    sync.Mutex
			order []string
		)
		s1 := &mockServer{name: "s1", stopFunc: func(_ bool) error {
			mu.Lock()
			order = append(order, "server")
			mu.Unlock()
			return nil
		}}
		app := gapp.New(s1)
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, a *gapp.App) (func(ctx context.Context), error) {
			return func(ctx context.Context) {
				mu.Lock()
				order = append(order, "hook")
				mu.Unlock()
			}, nil
		}))

		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)

		mu.Lock()
		// Hook runs before server stop.
		t.Assert(len(order), 2)
		t.Assert(order[0], "hook")
		t.Assert(order[1], "server")
		mu.Unlock()
	})
}

func TestAppBootConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var applyCount int32
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, _ *gapp.App) {
			atomic.AddInt32(&applyCount, 1)
			time.Sleep(50 * time.Millisecond)
		}))

		const callers = 20
		var wg sync.WaitGroup
		wg.Add(callers)
		for i := 0; i < callers; i++ {
			go func() {
				defer wg.Done()
				t.AssertNil(app.Boot(context.Background()))
			}()
		}
		wg.Wait()

		t.Assert(atomic.LoadInt32(&applyCount), int32(1))
		t.Assert(app.Booted(), true)
	})
}

func TestAppBootConcurrentFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var applyCount int32
		app := gapp.New()
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, _ *gapp.App) (func(context.Context), error) {
			atomic.AddInt32(&applyCount, 1)
			time.Sleep(50 * time.Millisecond)
			return nil, errTestBoot
		}))

		const callers = 10
		errs := make([]error, callers)
		var wg sync.WaitGroup
		wg.Add(callers)
		for i := 0; i < callers; i++ {
			idx := i
			go func() {
				defer wg.Done()
				errs[idx] = app.Boot(context.Background())
			}()
		}
		wg.Wait()

		t.Assert(atomic.LoadInt32(&applyCount), int32(1))
		t.Assert(app.Booted(), false)
		for i := range errs {
			t.AssertNE(errs[i], nil)
		}
	})
}

func TestAppContextPropagationBootHook(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			mu            sync.Mutex
			bootApplyID   string
			stopCleanupID string
		)

		traceCtx := gctx.WithSpan(context.Background(), "test-app-span")
		traceID := gctx.CtxId(traceCtx)

		app := gapp.New()
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, _ *gapp.App) (func(context.Context), error) {
			mu.Lock()
			bootApplyID = gctx.CtxId(ctx)
			mu.Unlock()

			return func(ctx context.Context) {
				mu.Lock()
				stopCleanupID = gctx.CtxId(ctx)
				mu.Unlock()
			}, nil
		}))

		s1 := &mockServer{name: "s1"}

		app.Add(s1)

		err := app.Boot(traceCtx)
		t.AssertNil(err)
		t.Assert(bootApplyID, traceID)

		err = app.Start(traceCtx)
		t.AssertNil(err)

		err = app.Stop(traceCtx, true)
		t.AssertNil(err)
		mu.Lock()
		t.Assert(stopCleanupID, traceID)
		mu.Unlock()
	})
}

func TestAppStartRespectsContextCancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1", startDelay: 200 * time.Millisecond}
		app := gapp.New(s1)

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Millisecond)
		defer cancel()

		err := app.Start(ctx)
		t.Assert(errors.Is(err, context.DeadlineExceeded), true)
		t.Assert(s1.stopped, true)
	})
}

func TestAppNilBootUsesNormalizedContext(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var mu sync.Mutex
		gotCtx := false
		app := gapp.New()
		app.Option(gapp.NewOption(func(ctx context.Context, _ *gapp.App) {
			mu.Lock()
			gotCtx = ctx != nil
			mu.Unlock()
		}))

		err := app.Boot(nil)
		t.AssertNil(err)
		mu.Lock()
		t.Assert(gotCtx, true)
		mu.Unlock()
	})
}

func TestAppLifecycleContextDefaultsAfterBoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var mu sync.Mutex
		stopCleanupID := ""

		traceBoot := gctx.WithSpan(context.Background(), "default-lifecycle")
		traceID := gctx.CtxId(traceBoot)

		app := gapp.New(&mockServer{name: "s1"})
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, _ *gapp.App) (func(context.Context), error) {
			return func(ctx context.Context) {
				mu.Lock()
				stopCleanupID = gctx.CtxId(ctx)
				mu.Unlock()
			}, nil
		}))

		err := app.Boot(traceBoot)
		t.AssertNil(err)

		err = app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(nil, true)
		t.AssertNil(err)

		mu.Lock()
		t.Assert(stopCleanupID, traceID)
		mu.Unlock()
	})
}

func TestAppRunBootFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		app := gapp.New(&mockServer{name: "s1"})
		app.Option(gapp.NewOptionWithHook(func(ctx context.Context, _ *gapp.App) (func(context.Context), error) {
			return nil, errTestBoot
		}))

		err := app.Run(context.Background())
		t.AssertNE(err, nil)
	})
}

func TestAppRunStartFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1", startErr: errTestStart}
		app := gapp.New(s1)

		err := app.Run(context.Background())
		t.AssertNE(err, nil)
	})
}

func TestAppRunContextCancelGracefulShutdown(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1"}
		app := gapp.New(s1)

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() {
			done <- app.Run(ctx)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case err := <-done:
			t.AssertNil(err)
		case <-time.After(2 * time.Second):
			t.Fatal("Run did not return after context cancellation")
		}
		t.Assert(s1.started, true)
		t.Assert(s1.stopped, true)
		t.Assert(s1.gracefulVal, true)
	})
}

func TestAppRunReturnsShutdownError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s1 := &mockServer{name: "s1", stopErr: errTestStop}
		app := gapp.New(s1)

		ctx, cancel := context.WithCancel(context.Background())
		done := make(chan error, 1)
		go func() {
			done <- app.Run(ctx)
		}()

		time.Sleep(50 * time.Millisecond)
		cancel()

		select {
		case err := <-done:
			t.AssertNE(err, nil)
		case <-time.After(2 * time.Second):
			t.Fatal("Run did not return after context cancellation")
		}
		t.Assert(s1.stopped, true)
	})
}

func TestAppStopIdempotent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var stopCount int32
		s1 := &mockServer{name: "s1", stopFunc: func(_ bool) error {
			atomic.AddInt32(&stopCount, 1)
			return nil
		}}
		app := gapp.New(s1)

		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&stopCount), int32(1))
	})
}

func TestAppStartAfterStopSecondStopIsNoOp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var stopCount int32
		s1 := &mockServer{name: "s1", stopFunc: func(_ bool) error {
			atomic.AddInt32(&stopCount, 1)
			return nil
		}}
		app := gapp.New(s1)

		err := app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&stopCount), int32(1))

		err = app.Start(context.Background())
		t.AssertNil(err)

		err = app.Stop(context.Background(), true)
		t.AssertNil(err)
		t.Assert(atomic.LoadInt32(&stopCount), int32(1))
	})
}

var (
	errTestStart = newTestError("start failed")
	errTestStop  = newTestError("stop failed")
	errTestBoot  = newTestError("boot failed")
)

type testError string

func newTestError(msg string) error { return testError(msg) }
func (e testError) Error() string   { return string(e) }

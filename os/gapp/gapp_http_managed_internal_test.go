// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover managed HTTP startup failure without terminating the application.

package gapp

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// managedHTTPFailOnceServer fails only its first startup to exercise HTTP rollback.
type managedHTTPFailOnceServer struct {
	failed bool
}

// Start returns one failure, then allows application startup retry.
func (s *managedHTTPFailOnceServer) Start() error {
	if !s.failed {
		s.failed = true
		return gerror.New("other server startup failed")
	}
	return nil
}

// Stop has no resources to release.
func (*managedHTTPFailOnceServer) Stop(bool) error { return nil }

// TestHTTPManagedApplicationRollbackRetry serves requests at the reported retry port.
func TestHTTPManagedApplicationRollbackRetry(t *testing.T) {
	const helperEnv = "GF_GAPP_MANAGED_HTTP_ROLLBACK_HELPER"
	if os.Getenv(helperEnv) == "" {
		gtest.C(t, func(t *gtest.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHTTPManagedApplicationRollbackRetry$")
			command.Env = append(os.Environ(), helperEnv+"=1")
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("managed HTTP rollback subprocess failed: %v\n%s", err, output)
			}
		})
		return
	}
	gtest.C(t, func(t *gtest.T) {
		var (
			name   = guid.S()
			server = ghttp.GetServer(name)
			app    = New(NewHTTPServerAdapter(server), &managedHTTPFailOnceServer{})
			client = &http.Client{Timeout: time.Second}
		)
		server.SetAddr("127.0.0.1:0")
		server.SetRegistrar(nil)
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.SetDumpRouterMap(false)
		server.BindHandler("/", func(r *ghttp.Request) { r.Response.Write("ok") })
		app.logger = glog.NewWithWriter(io.Discard)
		t.AssertNE(app.Start(context.Background()), nil)
		t.AssertNil(app.Start(context.Background()))
		defer func() { t.AssertNil(app.Stop(context.Background(), false)) }()
		t.Assert(len(server.GetListenedPorts()), 1)
		response, err := client.Get("http://" + server.GetListenedAddress())
		t.AssertNil(err)
		t.Assert(response.StatusCode, http.StatusOK)
		t.AssertNil(response.Body.Close())
		t.Assert(ghttp.GetServer(name) == server, true)
		// Global administration deliberately rate-limits actions for two seconds.
		time.Sleep(2100 * time.Millisecond)
		t.AssertNil(ghttp.ShutdownAllServer(context.Background()))
		connection, err := net.DialTimeout("tcp", server.GetListenedAddress(), time.Second)
		if err == nil {
			t.AssertNil(connection.Close())
			t.Fatal("global shutdown left the retried HTTP listener open")
		}
	})
}

// managedHTTPHelperEnv selects the subprocess that exercises formerly fatal startup.
const managedHTTPHelperEnv = "GF_GAPP_MANAGED_HTTP_HELPER"

// TestHTTPManagedStartupFailureCleansResources isolates any accidental os.Exit in a subprocess.
func TestHTTPManagedStartupFailureCleansResources(t *testing.T) {
	if os.Getenv(managedHTTPHelperEnv) == "" {
		gtest.C(t, func(t *gtest.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHTTPManagedStartupFailureCleansResources$", "-test.v")
			command.Env = append(os.Environ(), managedHTTPHelperEnv+"=1")
			if output, err := command.CombinedOutput(); err != nil {
				t.Fatalf("managed HTTP startup subprocess failed: %v\n%s", err, output)
			}
		})
		return
	}

	gtest.C(t, func(t *gtest.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		defer func() { t.AssertNil(listener.Close()) }()
		var (
			server        = ghttp.GetServer("gapp-managed-http-startup-failure")
			app           = New(NewHTTPServerAdapter(server))
			cleanupCalled bool
			ctx, cancel   = context.WithTimeout(context.Background(), 10*time.Second)
		)
		defer cancel()
		server.SetAddr(listener.Addr().String())
		server.BindHandler("/", func(r *ghttp.Request) { r.Response.Write("ok") })
		app.Option(NewOptionWithHook(func(context.Context, *App) (func(context.Context), error) {
			return func(context.Context) { cleanupCalled = true }, nil
		}))
		err = app.Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(strings.Contains(err.Error(), listener.Addr().String()), true)
		t.Assert(cleanupCalled, true)
	})
}

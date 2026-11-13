// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests HTTP internal route reuse and singleton ownership across managed retries.

package ghttp

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
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/test/gtest"
)

// TestHTTPManagedInternalRouteRetry isolates formerly fatal internal route duplication.
func TestHTTPManagedInternalRouteRetry(t *testing.T) {
	const helperEnv = "GF_TEST_HTTP_INTERNAL_ROUTE_RETRY"
	feature := os.Getenv(helperEnv)
	if feature == "" {
		for _, feature := range []string{"openapi", "pprof", "pprof-root", "swagger"} {
			t.Run(feature, func(t *testing.T) {
				gtest.C(t, func(t *gtest.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
					defer cancel()
					command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHTTPManagedInternalRouteRetry$")
					command.Env = append(os.Environ(), helperEnv+"="+feature)
					if output, err := command.CombinedOutput(); err != nil {
						t.Fatalf("internal route retry subprocess failed: %v\n%s", err, output)
					}
				})
			})
		}
		return
	}
	gtest.C(t, func(t *gtest.T) {
		busy, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		defer func() { t.AssertNil(busy.Close()) }()
		var (
			server = newManagedHTTPTestServer()
			client = &http.Client{Timeout: time.Second}
		)
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.SetAddr(busy.Addr().String())
		// Configure each feature using the same public configuration accepted on startup.
		configure := func(path string) string {
			switch feature {
			case "openapi":
				server.SetOpenApiPath(path)
			case "pprof", "pprof-root":
				if feature == "pprof-root" {
					path = "/"
				}
				server.config.PProfEnabled = true
				server.config.PProfPattern = path
				path = strings.TrimRight(path, "/") + "/cmdline"
			case "swagger":
				server.SetOpenApiPath("/spec")
				server.SetSwaggerPath(path)
				path += "/"
			}
			return path
		}
		configure("/internal")
		if strings.HasPrefix(feature, "pprof") {
			plugin := &managedTestPlugin{installErr: gerror.New("plugin failed before listening")}
			server.Plugin(plugin)
			t.AssertNE(server.StartManaged(), nil)
			plugin.installErr = nil
		}
		t.AssertNE(server.StartManaged(), nil)
		server.SetAddr("127.0.0.1:0")
		defer func() { t.AssertNil(server.Close()) }()
		for _, prefix := range []string{"/internal", "/internal-next", "/internal"} {
			nextPath := configure(prefix)
			t.AssertNil(server.StartManaged())
			routeCount, staticCount := len(server.GetRoutes()), len(server.config.StaticPaths)
			t.AssertNE(server.StartManaged(), nil)
			t.Assert(len(server.GetRoutes()), routeCount)
			t.Assert(len(server.config.StaticPaths), staticCount)
			response, err := client.Get("http://" + server.GetListenedAddress() + nextPath)
			t.AssertNil(err)
			t.Assert(response.StatusCode, http.StatusOK)
			t.AssertNil(response.Body.Close())
			t.AssertNil(server.Shutdown())
		}
		if feature == "swagger" {
			t.Assert(len(server.config.StaticPaths), 2)
		}
	})
}

// TestHTTPManagedInternalRoutesReuse retains one registration per configured internal route.
func TestHTTPManagedInternalRoutesReuse(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := newManagedHTTPTestServer()
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.SetOpenApiPath("/spec")
		server.SetSwaggerPath("/docs")
		server.config.PProfEnabled = true
		defer func() { t.AssertNil(server.Close()) }()
		var routeCount int
		for range 3 {
			t.AssertNil(server.StartManaged())
			if routeCount == 0 {
				routeCount = len(server.GetRoutes())
			}
			t.Assert(len(server.GetRoutes()), routeCount)
			t.Assert(len(server.config.StaticPaths), 1)
			t.AssertNE(server.StartManaged(), nil)
			t.AssertNil(server.Shutdown())
		}
	})
}

// TestHTTPManagedSingletonRetry restores a stopped singleton without replacing a live owner.
func TestHTTPManagedSingletonRetry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		server := newManagedHTTPTestServer()
		server.SetLogger(glog.NewWithWriter(io.Discard))
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		replacement := GetServer(server.instance)
		t.Assert(replacement == server, false)
		t.AssertNil(server.StartManaged())
		t.Assert(GetServer(server.instance) == server, true)
		t.AssertNil(server.Close())
		replacement.SetAddr("127.0.0.1:0")
		replacement.SetRegistrar(nil)
		replacement.SetLogger(glog.NewWithWriter(io.Discard))
		replacement.SetDumpRouterMap(false)
		replacement.BindHandler("/", func(r *Request) { r.Response.Write("replacement") })
		t.AssertNil(replacement.StartManaged())
		defer func() { t.AssertNil(replacement.Close()) }()
		registrar := &managedTestRegistrar{}
		server.registrar = registrar
		server.config.Endpoints = []string{"localhost:0"}
		t.AssertNE(server.StartManaged(), nil)
		t.AssertNil(registrar.registered)
		t.Assert(GetServer(server.instance) == replacement, true)
		t.AssertNil(server.Close())
		t.Assert(GetServer(server.instance) == replacement, true)
	})
}

// TestHTTPManagedInternalRouteConflict preserves user duplicate route detection.
func TestHTTPManagedInternalRouteConflict(t *testing.T) {
	const helperEnv = "GF_TEST_HTTP_INTERNAL_ROUTE_CONFLICT"
	if os.Getenv(helperEnv) != "" {
		server := newManagedHTTPTestServer()
		server.BindHandler("/spec", func(r *Request) { r.Response.Write("user handler") })
		server.SetOpenApiPath("/spec")
		if err := server.StartManaged(); err != nil {
			t.Fatal(err)
		}
		t.Fatal("user route conflict was ignored")
	}
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHTTPManagedInternalRouteConflict$")
		command.Env = append(os.Environ(), helperEnv+"=1")
		output, err := command.CombinedOutput()
		t.AssertNE(err, nil)
		t.AssertNil(ctx.Err())
		t.Assert(strings.Contains(string(output), "duplicated route registry"), true)
	})
}

// TestHTTPManagedChildHandoff leaves process handoff to the managed lifecycle owner.
func TestHTTPManagedChildHandoff(t *testing.T) {
	const helperEnv = "GF_TEST_HTTP_CHILD_HANDOFF"
	mode := os.Getenv(helperEnv)
	if mode == "" {
		for _, mode := range []string{"failed-managed", "closed-managed", "standalone"} {
			t.Run(mode, func(t *testing.T) {
				gtest.C(t, func(t *gtest.T) {
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer cancel()
					command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHTTPManagedChildHandoff$")
					// An impossible PID keeps the legacy positive control away from a real parent.
					command.Env = append(os.Environ(), helperEnv+"="+mode, "GPROC_PPID=2147483647", "GF_DEBUG=1")
					output, err := command.CombinedOutput()
					if err != nil {
						t.Fatalf("child handoff subprocess failed: %v\n%s", err, output)
					}
					t.Assert(strings.Contains(string(output), "notice parent server graceful shuttingdown"), mode == "standalone")
				})
			})
		}
		return
	}
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gproc.IsChild(), true)
		server := newManagedHTTPTestServer()
		server.SetLogger(glog.NewWithWriter(os.Stdout))
		server.SetGraceful(false)
		server.SetGracefulTimeout(1)
		defer func() { t.AssertNil(server.Close()) }()
		switch mode {
		case "failed-managed":
			server.registrar = &managedTestRegistrar{registerErr: gerror.New("registration failed")}
			server.config.Endpoints = []string{"localhost:0"}
			t.AssertNE(server.StartManaged(), nil)
		case "closed-managed":
			t.AssertNil(server.StartManaged())
			t.AssertNil(server.Close())
		case "standalone":
			t.AssertNil(server.Start())
		}
		// The standalone positive control verifies that the one-second timer has fired.
		time.Sleep(1800 * time.Millisecond)
	})
}

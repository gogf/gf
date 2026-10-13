// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests managed HTTP startup failures and partial resource rollback.

package ghttp

import (
	"context"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// managedTestPlugin records successful installations and removals.
type managedTestPlugin struct {
	installErr error
	installed  int
	removed    int
}

// Name identifies the test plugin.
func (*managedTestPlugin) Name() string { return "managed-test" }

// Author identifies the test author.
func (*managedTestPlugin) Author() string { return "GoFrame" }

// Version identifies the test version.
func (*managedTestPlugin) Version() string { return "1" }

// Description describes the test plugin.
func (*managedTestPlugin) Description() string { return "managed startup test" }

// Install records installation and returns the configured failure.
func (p *managedTestPlugin) Install(*Server) error {
	if p.installErr == nil {
		p.installed++
	}
	return p.installErr
}

// Remove records completed cleanup.
func (p *managedTestPlugin) Remove() error {
	p.removed++
	return nil
}

// managedTestRegistrar can fail registration and records deregistration.
type managedTestRegistrar struct {
	registerErr  error
	deregistered int
}

// Register simulates a registry that returns no handle when registration fails.
func (r *managedTestRegistrar) Register(_ context.Context, service gsvc.Service) (gsvc.Service, error) {
	if r.registerErr != nil {
		return nil, r.registerErr
	}
	return service, nil
}

// Deregister records successful removal.
func (r *managedTestRegistrar) Deregister(context.Context, gsvc.Service) error {
	r.deregistered++
	return nil
}

// newManagedHTTPTestServer creates an isolated routed server.
func newManagedHTTPTestServer() *Server {
	s := GetServer(guid.S())
	s.registrar = nil
	s.SetAddr("127.0.0.1:0")
	s.SetDumpRouterMap(false)
	s.BindHandler("/", func(r *Request) { r.Response.Write("ok") })
	return s
}

// TestHTTPManagedImmediateRestart rejects a second start before Serve is scheduled.
func TestHTTPManagedImmediateRestart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		previous := runtime.GOMAXPROCS(1)
		defer runtime.GOMAXPROCS(previous)
		server := newManagedHTTPTestServer()
		server.SetLogger(glog.NewWithWriter(io.Discard))
		t.AssertNil(server.StartManaged())
		defer func() { t.AssertNil(server.Close()) }()
		t.AssertNE(server.StartManaged(), nil)
		t.Assert(len(server.servers), 1)
	})
}

// TestHTTPManagedPluginFailure rolls back only plugins whose installation succeeded.
func TestHTTPManagedPluginFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			server  = newManagedHTTPTestServer()
			good    = &managedTestPlugin{}
			failure = gerror.New("plugin installation failed")
			bad     = &managedTestPlugin{installErr: failure}
		)
		server.Plugin(good, bad)
		err := server.StartManaged()
		t.Assert(gerror.Is(err, failure), true)
		t.Assert(good.removed, 1)
		t.Assert(bad.removed, 0)
		t.AssertNil(server.Close())
		t.Assert(good.removed, 1)
		bad.installErr = nil
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		t.Assert(good.removed, 2)
		t.Assert(bad.removed, 1)
	})
}

// TestHTTPManagedPartialBindFailure closes an earlier listener when a later bind fails.
func TestHTTPManagedPartialBindFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		first, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		firstAddress := first.Addr().String()
		busy, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		defer func() { t.AssertNil(busy.Close()) }()
		server := newManagedHTTPTestServer()
		server.config.Listeners = []net.Listener{first}
		server.SetAddr(firstAddress + "," + busy.Addr().String())
		t.AssertNE(server.StartManaged(), nil)
		t.Assert(len(server.servers), 0)
		listener, err := net.Listen("tcp", firstAddress)
		t.AssertNil(err)
		t.AssertNil(listener.Close())
		server.config.Listeners = nil
		server.SetAddr("127.0.0.1:0")
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
	})
}

// TestHTTPManagedRegistrationFailure returns errors without deregistering failed definitions.
func TestHTTPManagedRegistrationFailure(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			server    = newManagedHTTPTestServer()
			failure   = gerror.New("HTTP registration failed")
			registrar = &managedTestRegistrar{registerErr: failure}
			plugin    = &managedTestPlugin{}
		)
		server.registrar = registrar
		server.config.Endpoints = []string{"localhost:0"}
		server.Plugin(plugin)
		err := server.StartManaged()
		t.Assert(gerror.Is(err, failure), true)
		t.Assert(len(server.servers), 0)
		t.Assert(plugin.removed, 1)
		t.AssertNil(server.Close())
		t.Assert(registrar.deregistered, 0)
		registrar.registerErr = nil
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		t.Assert(registrar.deregistered, 1)
		t.Assert(plugin.removed, 2)
	})
}

// TestHTTPLegacyBindFailureIsFatal verifies standalone startup retains process-exit behavior.
func TestHTTPLegacyBindFailureIsFatal(t *testing.T) {
	const helperEnv = "GF_TEST_HTTP_LEGACY_BIND_FAILURE"
	if os.Getenv(helperEnv) != "" {
		gtest.C(t, func(t *gtest.T) {
			busy, err := net.Listen("tcp", "127.0.0.1:0")
			t.AssertNil(err)
			defer func() { t.AssertNil(busy.Close()) }()
			server := newManagedHTTPTestServer()
			server.SetAddr(busy.Addr().String())
			if err := server.Start(); err != nil {
				t.Fatalf("standalone Start returned an error: %v", err)
			}
			t.Fatal("standalone Start returned instead of terminating")
		})
		return
	}
	gtest.C(t, func(t *gtest.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestHTTPLegacyBindFailureIsFatal$")
		cmd.Env = append(os.Environ(), helperEnv+"=1")
		output, err := cmd.CombinedOutput()
		t.AssertNE(err, nil)
		t.AssertNil(ctx.Err())
		t.Assert(strings.Contains(string(output), "[FATA]"), true)
		t.Assert(strings.Contains(string(output), "address already in use"), true)
	})
}

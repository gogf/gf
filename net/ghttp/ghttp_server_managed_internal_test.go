// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests managed HTTP startup failures and partial resource rollback.

package ghttp

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// TestHTTPManagedTLSCertificateRetry serves HTTPS after correcting a failed certificate load.
func TestHTTPManagedTLSCertificateRetry(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing")
	gtest.C(t, func(t *gtest.T) {
		var (
			server    = newManagedHTTPTestServer()
			config    = NewConfig()
			transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			client    = &http.Client{Timeout: time.Second, Transport: transport}
		)
		// The fixture certificate is self-signed; only the handshake and response are under test.
		defer transport.CloseIdleConnections()
		config.Address = "127.0.0.1:0"
		config.TLSConfig = &tls.Config{}
		config.HTTPSCertPath = missingPath + ".crt"
		config.HTTPSKeyPath = missingPath + ".key"
		config.Logger = glog.NewWithWriter(io.Discard)
		config.DumpRouterMap = false
		t.AssertNil(server.SetConfig(config))
		defer func() { t.AssertNil(server.Close()) }()
		t.AssertNE(server.StartManaged(), nil)
		server.EnableHTTPS(gtest.DataPath("https", "files", "server.crt"), gtest.DataPath("https", "files", "server.key"))
		t.AssertNil(server.StartManaged())
		response, err := client.Get(fmt.Sprintf("https://127.0.0.1:%d/", server.GetListenedHTTPSPort()))
		t.AssertNil(err)
		defer func() { t.AssertNil(response.Body.Close()) }()
		body, err := io.ReadAll(response.Body)
		t.AssertNil(err)
		t.Assert(response.StatusCode, http.StatusOK)
		t.Assert(string(body), "ok")
	})
}

// managedTestPlugin records successful installations and removals.
type managedTestPlugin struct {
	installErr     error
	installed      int
	removed        int
	removeFailures int
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
	if p.removeFailures > 0 {
		p.removeFailures--
		return gerror.New("temporary HTTP plugin removal failure")
	}
	p.removed++
	return nil
}

// TestHTTPManagedPluginRemovalRetry preserves failed removal before another installation.
func TestHTTPManagedPluginRemovalRetry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			server = newManagedHTTPTestServer()
			plugin = &managedTestPlugin{removeFailures: 3}
		)
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.Plugin(plugin)
		defer func() { t.AssertNil(server.Close()) }()
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		t.AssertNE(server.StartManaged(), nil)
		t.Assert(plugin.installed, 1)
		t.Assert(plugin.removed, 0)
		t.AssertNil(server.StartManaged())
		t.Assert(plugin.installed, 2)
		t.Assert(plugin.removed, 1)
		t.AssertNil(server.Close())
		t.AssertNil(server.Close())
		t.Assert(plugin.removed, 2)
	})
}

// TestHTTPGlobalShutdownRemovesManagedPlugins completes ownership cleanup before returning.
func TestHTTPGlobalShutdownRemovesManagedPlugins(t *testing.T) {
	for _, forceful := range []bool{false, true} {
		name := "graceful"
		if forceful {
			name = "forceful"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					server = newManagedHTTPTestServer()
					plugin = &managedTestPlugin{}
					status = serverProcessStatus.Val()
				)
				defer serverProcessStatus.Set(status)
				server.SetLogger(glog.NewWithWriter(io.Discard))
				server.Plugin(plugin)
				defer func() { t.AssertNil(server.Close()) }()
				t.AssertNil(server.StartManaged())
				if forceful {
					forceCloseWebServers(context.Background())
				} else {
					shutdownWebServersGracefully(context.Background(), nil)
				}
				server.managedServeWG.Wait()
				t.Assert(server.serverCount.Val(), 0)
				t.Assert(plugin.removed, 1)
				t.AssertNil(server.Close())
				t.Assert(plugin.removed, 1)
			})
		})
	}
}

// managedDrainPlugin observes whether global graceful shutdown completed the active handler.
type managedDrainPlugin struct {
	managedTestPlugin
	completed <-chan struct{}
	removed   chan bool
}

// Remove reports whether the handler has finished before plugin resources are released.
func (p *managedDrainPlugin) Remove() error {
	select {
	case <-p.completed:
		p.removed <- true
	default:
		p.removed <- false
	}
	return nil
}

// managedCloseListener signals when global shutdown stops accepting requests.
type managedCloseListener struct {
	net.Listener
	closed chan struct{}
	once   sync.Once
}

// Close publishes listener closure while preserving the underlying close result.
func (l *managedCloseListener) Close() error {
	err := l.Listener.Close()
	l.once.Do(func() { close(l.closed) })
	return err
}

// TestHTTPGlobalShutdownDrainsBeforePluginRemoval preserves the global graceful cleanup order.
func TestHTTPGlobalShutdownDrainsBeforePluginRemoval(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		t.AssertNil(err)
		var (
			observed  = &managedCloseListener{Listener: listener, closed: make(chan struct{})}
			entered   = make(chan struct{})
			release   = make(chan struct{})
			completed = make(chan struct{})
			plugin    = &managedDrainPlugin{completed: completed, removed: make(chan bool, 1)}
			server    = newManagedHTTPTestServer()
			status    = serverProcessStatus.Val()
		)
		defer serverProcessStatus.Set(status)
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.SetAddr(listener.Addr().String())
		server.config.Listeners = []net.Listener{observed}
		server.BindHandler("/active", func(r *Request) {
			close(entered)
			<-release
			r.Response.Write("done")
			close(completed)
		})
		server.Plugin(plugin)
		defer func() { t.AssertNil(server.Close()) }()
		t.AssertNil(server.StartManaged())
		responseDone := make(chan error, 1)
		go func() {
			client := &http.Client{Timeout: 5 * time.Second}
			response, err := client.Get("http://" + server.GetListenedAddress() + "/active")
			if err == nil {
				err = response.Body.Close()
			}
			responseDone <- err
		}()
		<-entered
		shutdownDone := make(chan struct{})
		go func() {
			shutdownWebServersGracefully(context.Background(), nil)
			close(shutdownDone)
		}()
		<-observed.closed
		close(release)
		t.AssertNil(<-responseDone)
		<-shutdownDone
		t.Assert(<-plugin.removed, true)
	})
}

// managedTestRegistrar can fail registration and records deregistration.
type managedTestRegistrar struct {
	registerErr  error
	registered   gsvc.Service
	deregistered int
}

// managedRetryRegistrar retains active registry handles across temporary cleanup failures.
type managedRetryRegistrar struct {
	active             map[gsvc.Service]bool
	registered         int
	deregisterFailures int
	deregistered       []gsvc.Service
}

// Register returns a distinct registration handle for each bound listener.
func (r *managedRetryRegistrar) Register(_ context.Context, service gsvc.Service) (gsvc.Service, error) {
	r.registered++
	r.active[service] = true
	return service, nil
}

// Deregister fails selected attempts while leaving the registry handle active.
func (r *managedRetryRegistrar) Deregister(_ context.Context, service gsvc.Service) error {
	r.deregistered = append(r.deregistered, service)
	if r.deregisterFailures > 0 {
		r.deregisterFailures--
		return gerror.New("temporary HTTP deregistration failure")
	}
	delete(r.active, service)
	return nil
}

// TestHTTPManagedDeregistrationRetry cleans pending handles before registering a replacement.
func TestHTTPManagedDeregistrationRetry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			server    = newManagedHTTPTestServer()
			registrar = &managedRetryRegistrar{active: make(map[gsvc.Service]bool), deregisterFailures: 2}
		)
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.registrar = registrar
		server.config.Endpoints = []string{"localhost:0"}
		defer func() { t.AssertNil(server.Close()) }()
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		t.Assert(len(registrar.active), 1)
		t.AssertNE(server.StartManaged(), nil)
		t.Assert(registrar.registered, 1)
		t.Assert(len(registrar.active), 0)
		t.Assert(len(server.servers), 0)
		t.AssertNil(server.StartManaged())
		t.Assert(registrar.registered, 2)
		t.AssertNil(server.Close())
		t.Assert(len(registrar.active), 0)
		t.Assert(len(registrar.deregistered), 4)
		t.Assert(registrar.deregistered[0] == registrar.deregistered[1], true)
		t.Assert(registrar.deregistered[1] == registrar.deregistered[2], true)
	})
}

// TestHTTPManagedRepeatedCloseRetriesDeregistration repeats only a failed registry cleanup.
func TestHTTPManagedRepeatedCloseRetriesDeregistration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			server    = newManagedHTTPTestServer()
			registrar = &managedRetryRegistrar{active: make(map[gsvc.Service]bool), deregisterFailures: 1}
		)
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.registrar = registrar
		server.config.Endpoints = []string{"localhost:0"}
		defer func() { t.AssertNil(server.Close()) }()
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		t.Assert(len(registrar.active), 1)
		t.AssertNil(server.Close())
		t.AssertNil(server.Close())
		t.Assert(len(registrar.active), 0)
		t.Assert(len(registrar.deregistered), 2)
		t.Assert(registrar.deregistered[0] == registrar.deregistered[1], true)
	})
}

// Register simulates a registry that returns no handle when registration fails.
func (r *managedTestRegistrar) Register(_ context.Context, service gsvc.Service) (gsvc.Service, error) {
	if r.registerErr != nil {
		return nil, r.registerErr
	}
	r.registered = service
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

// TestHTTPManagedCloseRetry replaces closed listeners and advertises only live endpoints.
func TestHTTPManagedCloseRetry(t *testing.T) {
	for _, forceful := range []bool{false, true} {
		name := "graceful"
		if forceful {
			name = "forceful"
		}
		t.Run(name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					server    = newManagedHTTPTestServer()
					registrar = &managedTestRegistrar{}
					client    = &http.Client{Timeout: time.Second}
				)
				server.SetLogger(glog.NewWithWriter(io.Discard))
				server.registrar = registrar
				server.config.Endpoints = []string{"localhost:0"}
				defer func() { t.AssertNil(server.Close()) }()
				for range 3 {
					t.AssertNil(server.StartManaged())
					t.Assert(len(server.servers), 1)
					t.Assert(len(registrar.registered.GetEndpoints()), 1)
					t.Assert(registrar.registered.GetEndpoints()[0].String(), fmt.Sprintf("localhost:%d", server.GetListenedPort()))
					response, err := client.Get("http://" + server.GetListenedAddress())
					t.AssertNil(err)
					t.Assert(response.StatusCode, http.StatusOK)
					t.AssertNil(response.Body.Close())
					if forceful {
						t.AssertNil(server.Close())
					} else {
						t.AssertNil(server.Shutdown())
					}
					t.Assert(server.serverCount.Val(), 0)
				}
			})
		})
	}
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

// TestHTTPManagedCloseFinishesServe permits immediate retry without standalone notification readers.
func TestHTTPManagedCloseFinishesServe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		previous := runtime.GOMAXPROCS(1)
		defer runtime.GOMAXPROCS(previous)
		server := newManagedHTTPTestServer()
		server.SetLogger(glog.NewWithWriter(io.Discard))
		server.closeChan = make(chan struct{}, 1)
		server.closeChan <- struct{}{}
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
		t.Assert(server.serverCount.Val(), 0)
		t.AssertNil(server.StartManaged())
		t.AssertNil(server.Close())
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

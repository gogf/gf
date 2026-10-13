// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Tests cover managed HTTP startup failure without terminating the application.

package gapp

import (
	"context"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
)

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

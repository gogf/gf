// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/otel"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

var testCtx = context.Background()

func Test_OTEL_RequestTracing_Disabled(t *testing.T) {
	s := g.Server(guid.S())
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"result": "ok"})
	})
	s.SetDumpRouterMap(false)
	
	// By default, request tracing should be disabled
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Post(ctx, "/test", g.Map{"param1": "value1"})
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
	})
}

func Test_OTEL_RequestTracing_Enabled(t *testing.T) {
	s := g.Server(guid.S())
	
	// Enable request tracing using SetConfigWithMap
	err := s.SetConfigWithMap(g.Map{
		"OtelTraceRequestEnabled": true,
	})
	gtest.AssertNil(err)
	
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"result": "ok"})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Post(ctx, "/test?query1=qvalue1", g.Map{"param1": "value1"})
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
		// Test passes if no errors occurred during tracing
	})
}

func Test_OTEL_ResponseTracing_Enabled(t *testing.T) {
	s := g.Server(guid.S())
	
	// Enable response tracing using SetConfigWithMap
	err := s.SetConfigWithMap(g.Map{
		"OtelTraceResponseEnabled": true,
	})
	gtest.AssertNil(err)
	
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"result": "success", "data": "test data"})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Get(ctx, "/test")
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
		// Test passes if no errors occurred during response tracing
	})
}

func Test_OTEL_BothTracingEnabled(t *testing.T) {
	s := g.Server(guid.S())
	
	// Enable both request and response tracing using SetConfigWithMap
	err := s.SetConfigWithMap(g.Map{
		"OtelTraceRequestEnabled":  true,
		"OtelTraceResponseEnabled": true,
	})
	gtest.AssertNil(err)
	
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"received_param": r.Get("param1"),
			"received_query": r.Get("query1"),
			"result": "success",
		})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Post(ctx, "/test?query1=testquery", g.Map{"param1": "testparam"})
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
		// Test passes if no errors occurred during both request and response tracing
	})
}

func Test_OTEL_NewConfiguration_RequestTracing(t *testing.T) {
	s := g.Server(guid.S())
	
	// Enable request tracing using new independent OTEL configuration
	config := ghttp.NewConfig()
	config.Otel.TraceRequestEnabled = true
	err := s.SetConfig(config)
	gtest.AssertNil(err)
	
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"result": "ok"})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Post(ctx, "/test?query1=qvalue1", g.Map{"param1": "value1"})
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
		// Test configuration helper methods
		t.Assert(s.GetConfig().IsOtelTraceRequestEnabled(), true)
		t.Assert(s.GetConfig().IsOtelTraceResponseEnabled(), false)
	})
}

func Test_OTEL_NewConfiguration_ResponseTracing(t *testing.T) {
	s := g.Server(guid.S())
	
	// Enable response tracing using new independent OTEL configuration
	config := ghttp.NewConfig()
	config.Otel.TraceResponseEnabled = true
	err := s.SetConfig(config)
	gtest.AssertNil(err)
	
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"result": "success", "data": "test data"})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Get(ctx, "/test")
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
		// Test configuration helper methods
		t.Assert(s.GetConfig().IsOtelTraceRequestEnabled(), false)
		t.Assert(s.GetConfig().IsOtelTraceResponseEnabled(), true)
	})
}

func Test_OTEL_BackwardCompatibility(t *testing.T) {
	s := g.Server(guid.S())
	
	// Test that legacy configuration still works alongside new configuration
	config := ghttp.NewConfig()
	config.OtelTraceRequestEnabled = true  // Legacy field
	config.Otel.TraceResponseEnabled = true // New field
	err := s.SetConfig(config)
	gtest.AssertNil(err)
	
	s.BindHandler("/test", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{"result": "backward_compatible"})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		prefix := fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort())
		client := g.Client()
		client.SetPrefix(prefix)
		res, err := client.Post(ctx, "/test?query=test", g.Map{"param": "test"})
		t.AssertNil(err)
		defer res.Close()
		
		t.Assert(res.StatusCode, 200)
		// Test that both legacy and new configuration work together
		t.Assert(s.GetConfig().IsOtelTraceRequestEnabled(), true)
		t.Assert(s.GetConfig().IsOtelTraceResponseEnabled(), true)
	})
}

func Test_OTEL_Configuration_Helpers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test new OTEL config helpers
		otelConfig := otel.NewConfig()
		t.Assert(otelConfig.IsTracingSQLEnabled(), false)
		t.Assert(otelConfig.IsTracingRequestEnabled(), false)
		t.Assert(otelConfig.IsTracingResponseEnabled(), false)
		
		otelConfig.TraceSQLEnabled = true
		otelConfig.TraceRequestEnabled = true
		t.Assert(otelConfig.IsTracingSQLEnabled(), true)
		t.Assert(otelConfig.IsTracingRequestEnabled(), true)
		t.Assert(otelConfig.IsTracingResponseEnabled(), false)
	})
}
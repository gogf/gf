// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/internal/tracing"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

type CustomProvider struct {
	*sdkTrace.TracerProvider
}

func NewCustomProvider() *CustomProvider {
	return &CustomProvider{
		TracerProvider: sdkTrace.NewTracerProvider(
			sdkTrace.WithIDGenerator(NewCustomIDGenerator()),
		),
	}
}

type CustomIDGenerator struct{}

func NewCustomIDGenerator() *CustomIDGenerator {
	return &CustomIDGenerator{}
}

func (id *CustomIDGenerator) NewIDs(ctx context.Context) (traceID trace.TraceID, spanID trace.SpanID) {
	return tracing.NewIDs()
}

func (id *CustomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) (spanID trace.SpanID) {
	return tracing.NewSpanID()
}

func TestClient_CustomProvider(t *testing.T) {
	provider := otel.GetTracerProvider()
	defer otel.SetTracerProvider(provider)

	otel.SetTracerProvider(NewCustomProvider())

	s := g.Server(guid.S())
	s.BindHandler("/hello", func(r *ghttp.Request) {
		r.Response.WriteHeader(200)
		r.Response.WriteJson(g.Map{"field": "test_for_response_body"})
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		url := fmt.Sprintf("127.0.0.1:%d/hello", s.GetListenedPort())
		resp, err := c.DoRequest(ctx, http.MethodGet, url)
		t.AssertNil(err)
		t.AssertNE(resp, nil)
		t.Assert(resp.ReadAllString(), "{\"field\":\"test_for_response_body\"}")
	})
}

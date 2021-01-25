// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command jaeger is an example program that creates spans
// and uploads to Jaeger.
package main

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/label"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "http-trace-demo",
			Tags: []label.KeyValue{
				label.String("exporter", "jaeger"),
				label.Float64("float", 312.23),
			},
		}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}
	return flush
}

func main() {
	ctx := context.Background()

	flush := initTracer()
	defer flush()

	ctx, span := otel.Tracer("test").Start(ctx, "test")
	defer span.End()

	for i := 0; i < 20; i++ {
		g.Client().Use(ghttp.MiddlewareClientTracing).Ctx(ctx).Header(g.MapStrStr{
			"test": "123",
			"john": "smith",
		}).Cookie(g.MapStrStr{
			"cookieKey": "cookieValue",
		}).GetContent(fmt.Sprintf("http://baidu.com/?q=test_%d", i))
	}
	foo(ctx)
}

func foo(ctx context.Context) {
	ctx, span := otel.Tracer("test").Start(ctx, "foo")
	defer span.End()
	span.AddEvent("Nice operation!", trace.WithAttributes(label.Int("bogons", 100)))
	span.SetAttributes(label.String("test2", "123"))
	time.Sleep(time.Second * 1)
	bar(ctx)
}

func bar(ctx context.Context) {
	_, span := otel.Tracer("test").Start(ctx, "bar")
	defer span.End()
	span.AddEvent("Nice operation!", trace.WithAttributes(label.Int("bogons", 100)))
	span.SetAttributes(label.String("test2", "123"))
	time.Sleep(time.Second * 1)
	// Do bar...
}

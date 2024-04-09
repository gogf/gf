// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package otlpgrpc

import (
	"github.com/gogf/gf/contrib/observability/obs-opentelemetry/provider/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	serviceName = "otlp-grpc-client"
	endpoint    = "localhost:4317"
	traceToken  = "******_******"
)

func main() {
	var ctx = gctx.New()

	p := provider.NewOpenTelemetryProvider(
		ctx,
		provider.WithServiceName(serviceName),
		// Support setting ExportEndpoint via environment variables: OTEL_EXPORTER_OTLP_ENDPOINT
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithHeaders(map[string]string{"Authentication": traceToken}),
		provider.WithInsecure(),
	)
	defer func() {
		if err := p.Shutdown(ctx); err != nil {
			g.Log().Fatal(ctx, err)
		}
	}()

	StartRequests()
}

// StartRequests starts requests.
func StartRequests() {
	ctx, span := gtrace.NewSpan(gctx.New(), "StartRequests")
	defer span.End()

	ctx = gtrace.SetBaggageValue(ctx, "name", "john")

	content := g.Client().GetContent(ctx, "http://127.0.0.1:8199/hello")
	g.Log().Print(ctx, content)
}

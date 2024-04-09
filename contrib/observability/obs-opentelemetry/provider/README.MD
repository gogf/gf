---
title: "OpenTelemetry"
date: 2024-04-09
weight: 5
keywords: ["OpenTelemetry"]
description: "GoFrame provides openTelemetry capabilities."
---

[OpenTelemetry](https://opentelemetry.io/) is an open source observability framework
from [CNCF](https://www.cncf.io/) that consist of a series of tools, APIs and SDKs, and it enables IT teams to detect, generate,
collect, and export remote monitoring data for analysis and understanding of software performance and behavior.

### Provider

- Out-of-the-box default opentelemetry provider
- Support setting via environment variables:
    - [Exporter](https://opentelemetry.io/docs/reference/specification/protocol/exporter/)
    - [SDK](https://opentelemetry.io/docs/reference/specification/sdk-environment-variables/#general-sdk-configuration)

Examples:

```go
package main

import (
	"context"
    
	"github.com/gogf/gf/contrib/observability/obs-opentelemetry/provider/v2"
)

func main() {
	serviceName := "GoFrame-echo"
    ctx := context.Background()
    
	p := provider.NewOpenTelemetryProvider(ctx,
		provider.WithServiceName(serviceName),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer func() {
        if err := p.Shutdown(ctx); err != nil {
            panic(err)
        }
    }()
   
    
}
```

[Code](https://github.com/gogf/gf/contrib/observability/obs-opentelemetry/provider/v2)

#### Provider Options

| Function                  | Description                                                              |
| ------------------------- | ------------------------------------------------------------------------ |
| WithServiceName           | Configure the resource properties of `service.name`                      |
| WithDeploymentEnvironment | Configure the resource properties of `deployment.environment`            |
| WithServiceNamespace      | Configure the resource properties of `service.namespace`                 |
| WithResourceAttributes    | Configure the resource properties                                        |
| WithResource              | Configure resources (`resource.Resource`)                                |
| WithEnableTracing         | Enable `tracing`                                                         |
| WithEnableMetrics         | Enable `metrics`                                                         |
| WithTextMapPropagator     | Configure `propagation.TextMapPropagator`                                |
| WithResourceDetector      | Configure `resource.Detector`                                            |
| WithHeaders               | Configure the gRPC request header for exporting telemetry data           |
| WithInsecure              | Configure whether to use secure authentication for exported gRPC clients |
| WithEnableCompression     | Configure whether to enable gzip transport compression                   |
| WithSampler               | Configure Trace sampler                                                  |
| WithSdkTracerProvider     | Configure `sdktrace.TracerProvider`                                      |
| WithMeterProvider         | Configure `metric.MeterProvider`                                         |

## Full Examples

For a full usage: [example](https://github.com/gogf/gf/tree/master/example/opentelemetry)
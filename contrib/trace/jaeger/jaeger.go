package jaeger

import (
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.7.0"
)

// Init initializes and registers jaeger to global TracerProvider.
//
// The output parameter `flush` is used for waiting exported trace spans to be uploaded,
// which is useful if your program is ending, and you do not want to lose recent spans.
func Init(serviceName, endpoint string) (*trace.TracerProvider, error) {
	// Create the Jaeger exporter
	var endpointOption jaeger.EndpointOption
	if strings.HasPrefix(endpoint, "http") {
		// HTTP.
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		// UDP.
		match, err := gregex.MatchString(`(.+):(\d+)`, endpoint)
		if err != nil {
			return nil, err
		}
		if len(match) < 3 {
			return nil, gerror.NewCodef(
				gcode.CodeInvalidParameter, `invalid endpoint "%s"`, endpoint,
			)
		}
		var (
			host = match[1]
			port = match[2]
		)
		endpointOption = jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(host), jaeger.WithAgentPort(port),
		)
	}
	exp, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(
		// Always be sure to batch in production.
		trace.WithBatcher(exp),
		// Record information about this application in a Resource.
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	return tp, nil
}

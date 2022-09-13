module github.com/gogf/gf/contrib/trace/jaeger/v2

go 1.15

require (
	github.com/gogf/gf/v2 v2.1.4
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/exporters/jaeger v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
)

replace github.com/gogf/gf/v2 => ../../../

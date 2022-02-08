module github.com/gogf/gf/contrib/trace/jaeger/v2

go 1.15

require (
	github.com/gogf/gf/v2 v2.0.0-rc2
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
)

replace github.com/gogf/gf/v2 => ../../../

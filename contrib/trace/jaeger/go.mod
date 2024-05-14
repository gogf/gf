module github.com/gogf/gf/contrib/trace/jaeger/v2

go 1.18

require (
	github.com/gogf/gf/v2 v2.7.1
	go.opentelemetry.io/otel v1.14.0
	go.opentelemetry.io/otel/exporters/jaeger v1.14.0
	go.opentelemetry.io/otel/sdk v1.14.0
)

require (
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
)

replace github.com/gogf/gf/v2 => ../../../

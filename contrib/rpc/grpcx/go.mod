module github.com/gogf/gf/contrib/rpc/grpcx/v2

go 1.15

require (
	github.com/gogf/gf/contrib/registry/file/v2 v2.0.0-20230223141509-94b2eae1bec0
	github.com/gogf/gf/v2 v2.0.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	golang.org/x/net v0.0.0-20220919232410-f2f64ebce3c1 // indirect
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

replace (
	github.com/gogf/gf/contrib/registry/file/v2 => ../../registry/file/
	github.com/gogf/gf/v2 => ../../../
)

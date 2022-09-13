module github.com/gogf/gf/example

go 1.15

require (
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.1.4
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.1.4
	github.com/gogf/gf/contrib/registry/polaris/v2 v2.1.4
	github.com/gogf/gf/contrib/trace/jaeger/v2 v2.1.4
	github.com/gogf/gf/v2 v2.1.4
	github.com/gogf/katyusha v0.4.1-0.20220620125113-f55d6f739773
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/polarismesh/polaris-go v1.2.0-beta.2
	google.golang.org/grpc v1.49.0
)

replace (
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/contrib/registry/polaris/v2 => ../contrib/registry/polaris/
	github.com/gogf/gf/contrib/trace/jaeger/v2 => ../contrib/trace/jaeger/
	github.com/gogf/gf/v2 => ../
)

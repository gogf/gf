module github.com/gogf/gf/example

go 1.15

require (
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.1.0-rc3
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.1.0-rc3.0.20220523034830-510fa3faf03f
	github.com/gogf/gf/contrib/registry/polaris/v2 v2.0.0-rc2
	github.com/gogf/gf/contrib/trace/jaeger/v2 v2.0.0-rc2
	github.com/gogf/gf/v2 v2.1.0-rc3.0.20220523034830-510fa3faf03f
	github.com/gogf/katyusha v0.4.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/polarismesh/polaris-go v1.2.0-beta.0.0.20220517041223-596a6a63b00f
	google.golang.org/grpc v1.46.2
)

replace (
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/contrib/registry/polaris/v2 => ../contrib/registry/polaris/
	github.com/gogf/gf/contrib/trace/jaeger/v2 => ../contrib/trace/jaeger/
	github.com/gogf/gf/v2 => ../
)

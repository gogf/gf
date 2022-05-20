module github.com/gogf/gf/example

go 1.15

require (
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.1.0-rc3
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.1.0-rc3.0.20220520082600-c90acf81d6a8
	github.com/gogf/gf/contrib/registry/polaris/v2 v2.0.0-rc2
	github.com/gogf/gf/contrib/trace/jaeger/v2 v2.0.0-rc2
	github.com/gogf/gf/v2 v2.1.0-rc3.0.20220520082600-c90acf81d6a8
	github.com/gogf/katyusha v0.3.1-0.20220520092852-cd19a309028c
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/polarismesh/polaris-go v1.1.0
	google.golang.org/grpc v1.46.2
)

replace (
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/contrib/registry/polaris/v2 => ../contrib/registry/polaris/
	github.com/gogf/gf/contrib/trace/jaeger/v2 => ../contrib/trace/jaeger/
	github.com/gogf/gf/v2 => ../
)

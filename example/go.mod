module github.com/gogf/gf/example

go 1.15

require (
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.0.0-rc2
	github.com/gogf/gf/contrib/resolver/v2 v2.0.0-rc2
	github.com/gogf/gf/v2 v2.0.0-rc2
	github.com/golang/protobuf v1.5.2
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/contrib/resolver/v2 => ../contrib/resolver/
	github.com/gogf/gf/v2 => ../
)

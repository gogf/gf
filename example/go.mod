module github.com/gogf/gf/example

go 1.15

require (
	github.com/gogf/gf/contrib/config/apollo/v2 v2.5.1
	github.com/gogf/gf/contrib/config/kubecm/v2 v2.5.1
	github.com/gogf/gf/contrib/config/nacos/v2 v2.5.1
	github.com/gogf/gf/contrib/config/polaris/v2 v2.5.1
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.5.1
	github.com/gogf/gf/contrib/nosql/redis/v2 v2.5.1
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.5.1
	github.com/gogf/gf/contrib/registry/file/v2 v2.5.1
	github.com/gogf/gf/contrib/registry/polaris/v2 v2.5.1
	github.com/gogf/gf/contrib/rpc/grpcx/v2 v2.5.1
	github.com/gogf/gf/contrib/trace/jaeger/v2 v2.5.1
	github.com/gogf/gf/contrib/trace/otlpgrpc/v2 v2.5.1
	github.com/gogf/gf/contrib/trace/otlphttp/v2 v2.5.1
	github.com/gogf/gf/v2 v2.5.1
	github.com/nacos-group/nacos-sdk-go v1.1.2
	github.com/polarismesh/polaris-go v1.5.1
	google.golang.org/grpc v1.55.0
	google.golang.org/protobuf v1.30.0
	k8s.io/client-go v0.25.2
)

replace (
	github.com/gogf/gf/contrib/config/apollo/v2 => ../contrib/config/apollo/
	github.com/gogf/gf/contrib/config/kubecm/v2 => ../contrib/config/kubecm/
	github.com/gogf/gf/contrib/config/nacos/v2 => ../contrib/config/nacos/
	github.com/gogf/gf/contrib/config/polaris/v2 => ../contrib/config/polaris/
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/nosql/redis/v2 => ../contrib/nosql/redis/
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/contrib/registry/file/v2 => ../contrib/registry/file/
	github.com/gogf/gf/contrib/registry/polaris/v2 => ../contrib/registry/polaris/
	github.com/gogf/gf/contrib/rpc/grpcx/v2 => ../contrib/rpc/grpcx/
	github.com/gogf/gf/contrib/trace/jaeger/v2 => ../contrib/trace/jaeger/
	github.com/gogf/gf/contrib/trace/otlpgrpc/v2 => ../contrib/trace/otlpgrpc
	github.com/gogf/gf/contrib/trace/otlphttp/v2 => ../contrib/trace/otlphttp
	github.com/gogf/gf/v2 => ../
)

replace (
	golang.org/x/net v0.2.0 => golang.org/x/net v0.0.0-20221019024206-cb67ada4b0ad
	golang.org/x/sys v0.2.0 => golang.org/x/sys v0.0.0-20220906165534-d0df966e6959
)

module github.com/gogf/gf/example

go 1.15

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/HdrHistogram/hdrhistogram-go v1.1.2 // indirect
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.248 // indirect
	github.com/apolloconfig/agollo/v4 v4.3.0 // indirect
	github.com/clbanning/mxj/v2 v2.5.7 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/dlclark/regexp2 v1.8.1 // indirect
	github.com/emicklei/go-restful/v3 v3.10.2 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gogf/gf/contrib/config/apollo/v2 v2.3.3
	github.com/gogf/gf/contrib/config/kubecm/v2 v2.3.3
	github.com/gogf/gf/contrib/config/nacos/v2 v2.3.3
	github.com/gogf/gf/contrib/config/polaris/v2 v2.3.3
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.3.3
	github.com/gogf/gf/contrib/nosql/redis/v2 v2.3.3
	github.com/gogf/gf/contrib/registry/etcd/v2 v2.3.3
	github.com/gogf/gf/contrib/registry/file/v2 v2.3.3
	github.com/gogf/gf/contrib/registry/polaris/v2 v2.3.3
	github.com/gogf/gf/contrib/rpc/grpcx/v2 v2.3.3
	github.com/gogf/gf/contrib/trace/jaeger/v2 v2.3.3
	github.com/gogf/gf/v2 v2.3.3
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/imdario/mergo v0.3.14 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/nacos-group/nacos-sdk-go v1.1.4
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/polarismesh/polaris-go v1.3.0
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/viper v1.15.0 // indirect
	go.etcd.io/etcd/client/v3 v3.5.7 // indirect
	go.opentelemetry.io/otel/exporters/jaeger v1.14.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20230320184635-7606e756e683 // indirect
	google.golang.org/grpc v1.54.0
	google.golang.org/protobuf v1.30.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	k8s.io/client-go v0.26.3
	k8s.io/klog/v2 v2.90.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230308215209-15aac26d736a // indirect
	k8s.io/utils v0.0.0-20230313181309-38a27ef9d749 // indirect
)

replace (
	github.com/gogf/gf/contrib/config/apollo/v2 => ../contrib/config/apollo/
	github.com/gogf/gf/contrib/config/kubecm/v2 => ../contrib/config/kubecm/
	github.com/gogf/gf/contrib/config/polaris/v2 => ../contrib/config/polaris/
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/nosql/redis/v2 => ../contrib/nosql/redis/
	github.com/gogf/gf/contrib/registry/etcd/v2 => ../contrib/registry/etcd/
	github.com/gogf/gf/contrib/registry/file/v2 => ../contrib/registry/file/
	github.com/gogf/gf/contrib/registry/polaris/v2 => ../contrib/registry/polaris/
	github.com/gogf/gf/contrib/rpc/grpcx/v2 => ../contrib/rpc/grpcx/
	github.com/gogf/gf/contrib/trace/jaeger/v2 => ../contrib/trace/jaeger/
	github.com/gogf/gf/v2 => ../
)

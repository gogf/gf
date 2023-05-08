package nacos

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"net"
	"strconv"
)

var _ gsvc.Registry = &Registry{}

type Registry struct {
	namingClient naming_client.INamingClient
	opts         *options
}

type Option func(o *options)

type options struct {
	namespaceId         string
	timeoutMs           uint64
	notLoadCacheAtStart bool
	logDir              string
	cacheDir            string
	logLevel            string
	contextPath         string
	weight              float64
	clusterName         string
	groupName           string
}

func New(address []string, opts ...Option) *Registry {
	// default options
	options := &options{
		namespaceId:         "",
		timeoutMs:           5000,
		notLoadCacheAtStart: true,
		logDir:              "/tmp/nacos/log",
		cacheDir:            "/tmp/nacos/cache",
		logLevel:            "debug",
		contextPath:         "/gogf",
		weight:              1,
		clusterName:         "service",
		groupName:           "default",
	}
	// apply options
	for _, o := range opts {
		o(options)
	}
	//create ServerConfig
	sc := make([]constant.ServerConfig, 0)
	for i := range address {
		host, portStr, err := net.SplitHostPort(address[i])
		if err != nil {
			panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid nacos address "%s"`, address[i]))
		}
		port, err := strconv.Atoi(portStr)
		if err != nil {
			panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid nacos address "%s"`, address[i]))
		}
		sc = append(sc, *constant.NewServerConfig(host, uint64(port), constant.WithContextPath(options.contextPath)))
	}

	//create ClientConfig
	cc := constant.ClientConfig{
		NamespaceId:         options.namespaceId,
		TimeoutMs:           options.timeoutMs,
		NotLoadCacheAtStart: options.notLoadCacheAtStart,
		LogDir:              options.logDir,
		CacheDir:            options.cacheDir,
		LogLevel:            options.logLevel,
		ContextPath:         options.contextPath,
	}

	// create naming client
	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(gerror.Wrap(err, `create nacos namingClient failed`))
	}
	return &Registry{
		namingClient: client,
		opts:         options,
	}
}

func getServiceFromInstances(key string, opts *options, client naming_client.INamingClient) ([]gsvc.Service, error) {
	name := gstr.Split(key, "/")[4]
	split := gstr.Split(key, "/")
	version := split[len(split)-1]
	instances, err := client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: key,
		GroupName:   opts.groupName,
		Clusters:    []string{opts.clusterName},
		HealthyOnly: true,
	})
	if err != nil {
		return nil, err
	}
	result := make([]gsvc.Service, 0)
	endpoints := gsvc.Endpoints{}
	for _, instance := range instances {
		endpoints = append(endpoints, gsvc.NewEndpoint(instance.Ip+":"+strconv.Itoa(int(instance.Port))))
	}
	for i := range instances {
		metadata := make(map[string]interface{}, 0)
		for k, v := range instances[i].Metadata {
			metadata[k] = v
		}
		result = append(result, &gsvc.LocalService{
			Name:       name,
			Version:    version,
			Endpoints:  endpoints,
			Head:       opts.clusterName,
			Deployment: opts.groupName,
			Namespace:  opts.namespaceId,
			Metadata:   metadata,
		})
	}
	return result, nil
}

func WithNameSpaceId(namespace string) Option {
	return func(o *options) {
		o.namespaceId = namespace
	}
}

func WithTimeoutMs(timeoutMs uint64) Option {
	return func(o *options) {
		o.timeoutMs = timeoutMs
	}
}

func WithNotLoadCacheAtStart(notLoadCacheAtStart bool) Option {
	return func(o *options) {
		o.notLoadCacheAtStart = notLoadCacheAtStart
	}
}

func WithLogDir(logDir string) Option {
	return func(o *options) {
		o.logDir = logDir
	}
}

func WithCacheDir(cacheDir string) Option {
	return func(o *options) {
		o.cacheDir = cacheDir
	}
}

func WithLogLevel(logLevel string) Option {
	return func(o *options) {
		o.logLevel = logLevel
	}
}

func WithContextPath(contextPath string) Option {
	return func(o *options) {
		o.contextPath = contextPath
	}
}

func WithWeight(weight float64) Option {
	return func(o *options) {
		o.weight = weight
	}
}

func WithClusterName(clusterName string) Option {
	return func(o *options) {
		o.clusterName = clusterName
	}
}

func WithGroupName(groupName string) Option {
	return func(o *options) {
		o.groupName = groupName
	}
}

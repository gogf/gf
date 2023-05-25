// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package nacos implements service Registry and Discovery using nacos.
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

const (
	DefaultNamespaceId = ""
	DefaultTimeoutMs   = 5 * 1000
	DefaultWeight      = 1
	DefaultContextPath = "/goframe"
)

// Registry is nacos registry.
type Registry struct {
	namingClient naming_client.INamingClient
	opts         *options
}

// Option is nacos registry option.
type Option func(o *options)

type options struct {
	namespaceId string
	timeoutMs   uint64
	weight      float64
	clusterName string
	groupName   string
	contextPath string
	version     string
}

// New create nacos registry.
func New(address []string, opts ...Option) *Registry {
	// default options
	options := &options{
		namespaceId: DefaultNamespaceId,
		timeoutMs:   DefaultTimeoutMs,
		weight:      DefaultWeight,
		clusterName: gsvc.DefaultHead,
		groupName:   gsvc.DefaultDeployment,
		contextPath: DefaultContextPath,
		version:     gsvc.DefaultVersion,
	}

	// apply options
	for _, o := range opts {
		o(options)
	}

	// check options
	options.groupName = gstr.Join(gstr.Split(options.groupName, "/"), "-")

	// create ServerConfig
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
		NamespaceId: options.namespaceId,
		TimeoutMs:   options.timeoutMs,
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
	var (
		versionLength = len(gstr.Split(opts.version, "/"))
		name          = gstr.Join(gstr.Split(key, gsvc.DefaultSeparator)[4:len(gstr.Split(key, gsvc.DefaultSeparator))-versionLength], "/")
		version       = opts.version
	)
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

// WithNameSpaceId set namespaceId
func WithNameSpaceId(namespace string) Option {
	return func(o *options) {
		o.namespaceId = namespace
	}
}

// WithTimeoutMs set timeoutMs
func WithTimeoutMs(timeoutMs uint64) Option {
	return func(o *options) {
		o.timeoutMs = timeoutMs
	}
}

// WithWeight set weight
func WithWeight(weight float64) Option {
	return func(o *options) {
		o.weight = weight
	}
}

// WithClusterName set clusterName
func WithClusterName(clusterName string) Option {
	return func(o *options) {
		o.clusterName = clusterName
	}
}

// WithGroupName set groupName
func WithGroupName(groupName string) Option {
	return func(o *options) {
		o.groupName = groupName
	}
}

// WithContextPath set contextPath
func WithContextPath(contextPath string) Option {
	return func(o *options) {
		o.contextPath = contextPath
	}
}

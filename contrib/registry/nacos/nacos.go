// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package nacos implements service Registry and Discovery using nacos.
package nacos

import (
	"context"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	cstServiceSeparator = "@@"
)

var (
	_ gsvc.Registry = &Registry{}
)

// Option is the option for creating Registry.
type Option func(*Registry)

// WithClusterName sets the cluster name for the registry.
// The default is 'DEFAULT'.
func WithClusterName(name string) Option {
	return func(r *Registry) {
		r.clusterName = name
	}
}

// WithGroupName sets the group name for the registry.
// The default is 'DEFAULT_GROUP'.
func WithGroupName(name string) Option {
	return func(r *Registry) {
		r.groupName = name
	}
}

// WithDefaultEndpoint sets the default endpoint for service registration.
// It overrides the service endpoints when registering if it's not empty.
func WithDefaultEndpoint(endpoint string) Option {
	return func(r *Registry) {
		r.defaultEndpoint = endpoint
	}
}

// WithDefaultMetadata sets the default metadata for service registration.
// It will be merged with service's original metadata when registering.
func WithDefaultMetadata(metadata map[string]string) Option {
	return func(r *Registry) {
		r.defaultMetadata = metadata
	}
}

// Registry is nacos registry.
type Registry struct {
	client          naming_client.INamingClient
	clusterName     string
	groupName       string
	defaultEndpoint string
	defaultMetadata map[string]string
}

// Config is the configuration object for nacos client.
type Config struct {
	ServerConfigs []constant.ServerConfig `v:"required"` // See constant.ServerConfig
	ClientConfig  *constant.ClientConfig  `v:"required"` // See constant.ClientConfig
}

// New new a registry with address and opts
func New(address string, opts ...constant.ClientOption) (reg *Registry) {
	endpoints := gstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid nacos address "%s"`, address))
	}

	clientConfig := constant.NewClientConfig(opts...)

	if len(clientConfig.NamespaceId) == 0 {
		clientConfig.NamespaceId = "public"
	}

	serverConfigs := make([]constant.ServerConfig, 0, len(endpoints))
	for _, endpoint := range endpoints {
		tmp := gstr.Split(endpoint, ":")
		ip := tmp[0]
		port := gconv.Uint64(tmp[1])
		if port == 0 {
			port = 8848
		}
		serverConfigs = append(serverConfigs, *constant.NewServerConfig(ip, port))
	}
	ctx := gctx.New()
	reg, err := NewWithConfig(ctx, Config{
		ServerConfigs: serverConfigs,
		ClientConfig:  clientConfig,
	})

	if err != nil {
		panic(gerror.Wrap(err, `create nacos client failed`))
	}
	return
}

// NewWithConfig creates and returns registry with Config.
func NewWithConfig(ctx context.Context, config Config, opts ...Option) (reg *Registry, err error) {
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}

	nameingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  config.ClientConfig,
		ServerConfigs: config.ServerConfigs,
	})
	if err != nil {
		return
	}
	return NewWithClient(nameingClient, opts...), nil
}

// NewWithClient new the instance with INamingClient
func NewWithClient(client naming_client.INamingClient, opts ...Option) *Registry {
	r := &Registry{
		client:          client,
		clusterName:     "DEFAULT",
		groupName:       "DEFAULT_GROUP",
		defaultMetadata: make(map[string]string),
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// SetClusterName can set the clusterName. The default is 'DEFAULT'
//
// Deprecated: Use WithClusterName option in NewWithClient/NewWithConfig instead.
func (reg *Registry) SetClusterName(clusterName string) *Registry {
	reg.clusterName = clusterName
	return reg
}

// SetGroupName can set the groupName. The default is 'DEFAULT_GROUP'
//
// Deprecated: Use WithGroupName option in NewWithClient/NewWithConfig instead.
func (reg *Registry) SetGroupName(groupName string) *Registry {
	reg.groupName = groupName
	return reg
}

// SetDefaultEndpoint sets the default endpoint for service registration.
// It overrides the service endpoints when registering if it's not empty.
//
// Deprecated: Use WithDefaultEndpoint option in NewWithClient/NewWithConfig instead.
func (reg *Registry) SetDefaultEndpoint(endpoint string) *Registry {
	reg.defaultEndpoint = endpoint
	return reg
}

// SetDefaultMetadata sets the default metadata for service registration.
// It will be merged with service's original metadata when registering.
//
// Deprecated: Use WithDefaultMetadata option in NewWithClient/NewWithConfig instead.
func (reg *Registry) SetDefaultMetadata(metadata map[string]string) *Registry {
	reg.defaultMetadata = metadata
	return reg
}
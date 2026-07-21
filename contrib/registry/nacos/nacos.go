// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package nacos implements service Registry and Discovery using nacos.
package nacos

import (
	"context"
	"sync/atomic"

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

// Registry is nacos registry.
type Registry struct {
	client          naming_client.INamingClient
	clusterName     string
	groupName       string
	defaultEndpoint string
	// defaultMetadata is stored as a snapshot map via atomic.Value to avoid
	// data races when SetDefaultMetadata races with Register (see #4649).
	defaultMetadata atomic.Value // map[string]string
}

// Option configures a Registry at construction time.
type Option func(r *Registry)

// WithClusterName sets the cluster name. Default is "DEFAULT".
func WithClusterName(clusterName string) Option {
	return func(r *Registry) {
		r.clusterName = clusterName
	}
}

// WithGroupName sets the group name. Default is "DEFAULT_GROUP".
func WithGroupName(groupName string) Option {
	return func(r *Registry) {
		r.groupName = groupName
	}
}

// WithDefaultEndpoint sets the default endpoint used on Register when non-empty.
func WithDefaultEndpoint(endpoint string) Option {
	return func(r *Registry) {
		r.defaultEndpoint = endpoint
	}
}

// WithDefaultMetadata sets default metadata merged into service metadata on Register.
// The map is copied; later mutations of the input map do not affect the registry.
func WithDefaultMetadata(metadata map[string]string) Option {
	return func(r *Registry) {
		r.storeDefaultMetadata(metadata)
	}
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

// NewWithClient new the instance with INamingClient.
// Prefer With* options for init-time configuration instead of Set* methods.
func NewWithClient(client naming_client.INamingClient, opts ...Option) *Registry {
	r := &Registry{
		client:      client,
		clusterName: "DEFAULT",
		groupName:   "DEFAULT_GROUP",
	}
	r.storeDefaultMetadata(nil)
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

// SetClusterName can set the clusterName. The default is 'DEFAULT'.
// Deprecated: use WithClusterName at construction time instead.
func (reg *Registry) SetClusterName(clusterName string) *Registry {
	reg.clusterName = clusterName
	return reg
}

// SetGroupName can set the groupName. The default is 'DEFAULT_GROUP'.
// Deprecated: use WithGroupName at construction time instead.
func (reg *Registry) SetGroupName(groupName string) *Registry {
	reg.groupName = groupName
	return reg
}

// SetDefaultEndpoint sets the default endpoint for service registration.
// It overrides the service endpoints when registering if it's not empty.
// Deprecated: use WithDefaultEndpoint at construction time instead.
func (reg *Registry) SetDefaultEndpoint(endpoint string) *Registry {
	reg.defaultEndpoint = endpoint
	return reg
}

// SetDefaultMetadata sets the default metadata for service registration.
// It will be merged with service's original metadata when registering.
// The map is copied so concurrent Register is race-free with later Sets.
// Deprecated: use WithDefaultMetadata at construction time instead.
func (reg *Registry) SetDefaultMetadata(metadata map[string]string) *Registry {
	reg.storeDefaultMetadata(metadata)
	return reg
}

func (reg *Registry) storeDefaultMetadata(metadata map[string]string) {
	cp := make(map[string]string, len(metadata))
	for k, v := range metadata {
		cp[k] = v
	}
	reg.defaultMetadata.Store(cp)
}

func (reg *Registry) loadDefaultMetadata() map[string]string {
	v := reg.defaultMetadata.Load()
	if v == nil {
		return nil
	}
	m, _ := v.(map[string]string)
	return m
}

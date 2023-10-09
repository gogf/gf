// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package nacos implements service Registry and Discovery using nacos.
package nacos

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/joy999/nacos-sdk-go/clients"
	"github.com/joy999/nacos-sdk-go/clients/naming_client"
	"github.com/joy999/nacos-sdk-go/common/constant"
	"github.com/joy999/nacos-sdk-go/vo"
)

// Registry is nacos registry.
type Registry struct {
	client      naming_client.INamingClient
	clusterName string
	groupName   string
}

// ClientOption is cname the constant.ClientOption
type ClientOption = constant.ClientOption

// ClientConfig is cname the constant.ClientConfig
type ClientConfig = constant.ClientConfig

// Config is the configuration object for nacos client.
type Config struct {
	ServerConfigs []constant.ServerConfig `v:"required"` // See constant.ServerConfig
	ClientConfig  constant.ClientConfig   `v:"required"` // See constant.ClientConfig
}

// New new a registry with address and opts
func New(address string, opts ...ClientOption) *Registry {
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

	nameingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		panic(gerror.Wrap(err, `create nacos client failed`))
	}
	return NewWithClient(nameingClient)
}

// New creates and returns registry with Config.
func NewWithConfig(ctx context.Context, config Config) (reg *Registry, err error) {
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}

	nameingClient, err := clients.NewNamingClient(vo.NacosClientParam{
		ClientConfig:  &config.ClientConfig,
		ServerConfigs: config.ServerConfigs,
	})
	if err != nil {
		return
	}
	return NewWithClient(nameingClient), nil
}

// NewWithClient new the instance with INamingClient
func NewWithClient(client naming_client.INamingClient) *Registry {
	r := &Registry{
		client:      client,
		clusterName: "DEFAULT",
		groupName:   "DEFAULT_GROUP",
	}
	return r
}

// SetClusterName can set the clusterName. The default is 'DEFAULT'
func (reg *Registry) SetClusterName(clusterName string) *Registry {
	reg.clusterName = clusterName
	return reg
}

// SetGroupName can set the groupName. The default is 'DEFAULT_GROUP'
func (reg *Registry) SetGroupName(groupName string) *Registry {
	reg.groupName = groupName
	return reg
}

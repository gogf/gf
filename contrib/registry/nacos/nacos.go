// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

type Registry struct {
	client      naming_client.INamingClient
	clusterName string
	groupName   string
}

type ClientOption = constant.ClientOption
type ClientConfig = constant.ClientConfig

// NewWithConfig new with the default config file.
func NewWithConfig(addrees string, opts ...ClientOption) *Registry {
	ctx := gctx.New()
	conf := g.Config()

	clusterName := conf.MustGet(ctx, "nacos.cluster_name", "DEFAULT").String()
	groupName := conf.MustGet(ctx, "nacos.group_name", "DEFAULT_GROUP").String()
	serviceName := conf.MustGet(ctx, "nacos.service_name").String()
	logDir := conf.MustGet(ctx, "nacos.log_dir").String()
	logDir = filepath.Join(logDir, serviceName)
	cacheDir := conf.MustGet(ctx, "nacos.cache_dir").String()
	cacheDir = filepath.Join(cacheDir, serviceName)

	return New(addrees, func(c *ClientConfig) {
		c.NamespaceId = conf.MustGet(ctx, "nacos.namespace_id", "").String()
		c.Endpoint = conf.MustGet(ctx, "nacos.endpoint", "").String()
		c.AppName = serviceName
		c.TimeoutMs = conf.MustGet(ctx, "nacos.timeout_ms", 5000).Uint64()
		c.CacheDir = cacheDir
		c.LogDir = logDir
		c.LogLevel = conf.MustGet(ctx, "nacos.log_level", "error").String()
	}).SetClusterName(clusterName).SetGroupName(groupName)
}

// New new a registry with address and opts
func New(address string, opts ...ClientOption) *Registry {
	endpoints := gstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid nacos address "%s"`, address))
	}

	clientConfig := &ClientConfig{
		TimeoutMs: 5000,
		LogLevel:  "error",
	}

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(clientConfig)
		}
	}

	serverConfigs := make([]constant.ServerConfig, 0, len(endpoints))
	for _, endpoint := range endpoints {
		tmp := gstr.Split(endpoint, ":")
		ip := tmp[0]
		port := gconv.Uint64(tmp[1])

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

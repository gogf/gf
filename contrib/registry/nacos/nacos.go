// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
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

func New(ctx context.Context, address string, opts ...constant.ClientOption) gsvc.Registry {
	endpoints := gstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid nacos address "%s"`, address))
	}

	conf := g.Config()

	clusterName := conf.MustGet(ctx, "nacos.cluster_name", "DEFAULT").String()
	groupName := conf.MustGet(ctx, "nacos.group_name", "DEFAULT_GROUP").String()
	serviceName := conf.MustGet(ctx, "nacos.service_name").String()
	logDir := conf.MustGet(ctx, "nacos.log_dir").String()
	logDir = filepath.Join(logDir, serviceName)
	cacheDir := conf.MustGet(ctx, "nacos.cache_dir").String()
	cacheDir = filepath.Join(cacheDir, serviceName)

	clientConfig := &constant.ClientConfig{
		NamespaceId: conf.MustGet(ctx, "nacos.namespace_id", "").String(),
		Endpoint:    conf.MustGet(ctx, "nacos.endpoint", "").String(),
		AppName:     serviceName,
		TimeoutMs:   conf.MustGet(ctx, "nacos.timeout_ms", 5000).Uint64(),
		CacheDir:    cacheDir,
		LogDir:      logDir,
		LogLevel:    conf.MustGet(ctx, "nacos.log_level", "error").String(),
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
	r := NewWithClient(ctx, nameingClient)
	r.clusterName = clusterName
	r.groupName = groupName

	return r
}

func NewWithClient(ctx context.Context, client naming_client.INamingClient) *Registry {
	r := &Registry{
		client: client,
	}
	return r
}

func (reg *Registry) SetClusterName(clusterName string) {
	reg.clusterName = clusterName
}

func (reg *Registry) SetGroupName(groupName string) {
	reg.groupName = groupName
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"context"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/resolver"
)

type etcdRegistryConfig struct {
	Endpoints []string
}

// autoLoadAndRegisterEtcdRegistry checks and registers ETCD service as default service registry
// if no registry is registered previously.
func autoLoadAndRegisterEtcdRegistry() {
	if gsvc.GetRegistry() != nil {
		return
	}
	var (
		config      *etcdRegistryConfig
		ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
	)
	defer cancel()

	if !g.Cfg().Available(ctx) {
		g.Log().Fatal(ctx, `no configuration available`)
		return
	}
	node, err := g.Cfg().Get(ctx, configNodeNameRegistry)
	if err != nil {
		g.Log().Fatalf(ctx, `configuration "%s" load failed`, configNodeNameRegistry)
		return
	}
	// If no configuration, nothing to do.
	if node == nil {
		return
	}
	if err = node.Scan(&config); err != nil || config == nil {
		g.Log().Fatalf(ctx, `configuration "%s" load failed`, configNodeNameRegistry)
		return
	}
	if len(config.Endpoints) == 0 {
		g.Log().Fatalf(ctx, `empty endpoints in configuration "%s"`, configNodeNameRegistry)
		return
	}
	g.Log().Debugf(ctx, `set default registry using etcd service, address: %s`, config.Endpoints[0])
	resolver.SetRegistry(etcd.New(config.Endpoints[0]))
}

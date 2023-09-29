// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (reg *Registry) Register(ctx context.Context, service gsvc.Service) (registered gsvc.Service, err error) {
	c := reg.client
	metadata := map[string]string{}
	for k, v := range service.GetMetadata() {
		metadata[k] = gconv.String(v)
	}

	for _, endpoint := range service.GetEndpoints() {
		if _, err = c.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          endpoint.Host(),
			Port:        uint64(endpoint.Port()),
			ServiceName: service.GetName(),
			Metadata:    metadata,
			Weight:      100,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			ClusterName: reg.clusterName,
			GroupName:   reg.groupName,
		}); err != nil {
			return
		}
	}

	registered = service

	return
}

// Deregister off-lines and removes `service` from the Registry.
func (reg *Registry) Deregister(ctx context.Context, service gsvc.Service) (err error) {
	c := reg.client

	for _, endpoint := range service.GetEndpoints() {
		if _, err = c.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          endpoint.Host(),
			Port:        uint64(endpoint.Port()),
			ServiceName: service.GetName(),
			Ephemeral:   true,
			Cluster:     reg.clusterName,
			GroupName:   reg.groupName,
		}); err != nil {
			return
		}
	}

	return
}

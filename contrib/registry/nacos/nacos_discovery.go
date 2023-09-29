// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// Search searches and returns services with specified condition.
func (reg *Registry) Search(ctx context.Context, in gsvc.SearchInput) (result []gsvc.Service, err error) {
	if in.Prefix == "" && in.Name != "" {
		in.Prefix = gsvc.NewServiceWithName(in.Name).GetPrefix()
	}

	c := reg.client

	serviceName := in.Name
	if serviceName == "" {
		info := gstr.SplitAndTrim(gstr.Trim(in.Prefix, "/"), "/")
		if len(info) >= 2 {
			serviceName = info[len(info)-2]
		}
	}
	param := vo.SelectInstancesParam{
		GroupName:   reg.groupName,
		Clusters:    []string{reg.clusterName},
		ServiceName: serviceName,
		HealthyOnly: true,
	}
	instances, err := c.SelectInstances(param)
	if err != nil {
		return
	}

	result = make([]gsvc.Service, 0, len(instances))
	for _, inst := range instances {
		result = append(result, NewServiceFromInstance(&inst))
	}

	return
}

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (reg *Registry) Watch(ctx context.Context, key string) (watcher gsvc.Watcher, err error) {
	c := reg.client

	w := newWather(ctx)

	fn := func(services []model.Instance, err error) {
		w.Push(services, err)
	}

	sArr := gstr.Split(key, "/")

	serviceName := sArr[4]

	param := &vo.SubscribeParam{
		ServiceName:       serviceName,
		GroupName:         reg.groupName,
		Clusters:          []string{reg.ClusterName},
		SubscribeCallback: fn,
	}

	w.SetCloseFunc(func() error {
		return c.Unsubscribe(param)
	})

	err = c.Subscribe(param)
	if err != nil {
		return
	}

	watcher = w
	return
}

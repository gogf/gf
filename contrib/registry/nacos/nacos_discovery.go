// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/joy999/nacos-sdk-go/model"
	"github.com/joy999/nacos-sdk-go/vo"
)

// Search searches and returns services with specified condition.
func (reg *Registry) Search(ctx context.Context, in gsvc.SearchInput) (result []gsvc.Service, err error) {
	if in.Prefix == "" && in.Name != "" {
		in.Prefix = gsvc.NewServiceWithName(in.Name).GetPrefix()
	}

	c := reg.client

	serviceName := in.Name
	if serviceName == "" {
		info := gstr.SplitAndTrim(gstr.Trim(in.Prefix, gsvc.DefaultSeparator), gsvc.DefaultSeparator)
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

	insts := make([]model.Instance, 0, len(instances))
inst_loop:
	for _, inst := range instances {
		if len(in.Metadata) > 0 {
			for k, v := range in.Metadata {
				if inst.Metadata[k] != v {
					continue inst_loop
				}
			}
		}
		insts = append(insts, inst)
	}

	result = NewServicesFromInstances(insts)
	return
}

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (reg *Registry) Watch(ctx context.Context, key string) (watcher gsvc.Watcher, err error) {
	c := reg.client

	w := newWatcher(ctx)

	fn := func(services []model.Instance, err error) {
		w.Push(services, err)
	}

	sArr := gstr.Split(key, gsvc.DefaultSeparator)
	if len(sArr) < 5 {
		err = gerror.NewCode(gcode.CodeInvalidParameter, "The 'key' is invalid")
		return
	}
	serviceName := sArr[4]

	param := &vo.SubscribeParam{
		ServiceName:       serviceName,
		GroupName:         reg.groupName,
		Clusters:          []string{reg.clusterName},
		SubscribeCallback: fn,
	}

	w.SetCloseFunc(func() error {
		return c.Unsubscribe(param)
	})

	if err = c.Subscribe(param); err != nil {
		return
	}

	watcher = w
	return
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"fmt"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/joy999/nacos-sdk-go/model"
)

// NewServiceFromInstance new one service from instance
func NewServiceFromInstance(instance []model.Instance) gsvc.Service {
	n := len(instance)
	if n == 0 {
		return nil
	}
	serviceName := instance[0].ServiceName
	endpoints := make(gsvc.Endpoints, 0, n)
	for i := 0; i < n; i++ {
		if instance[0].ServiceName != serviceName {
			return nil
		}
		endpoints = append(endpoints, gsvc.NewEndpoint(fmt.Sprintf("%s%s%d", instance[i].Ip, gsvc.EndpointHostPortDelimiter, int(instance[i].Port))))
	}
	if gstr.Contains(serviceName, cstServiceSeparator) {
		arr := gstr.SplitAndTrim(serviceName, cstServiceSeparator)
		serviceName = arr[1]
	}

	return &gsvc.LocalService{
		Endpoints: endpoints,
		Name:      serviceName,
		Metadata:  gmap.NewStrStrMapFrom(instance[0].Metadata).MapStrAny(),
		Version:   gsvc.DefaultVersion,
	}
}

// NewServicesFromInstances new some services from some instances
func NewServicesFromInstances(instances []model.Instance) []gsvc.Service {
	serviceMap := map[string][]model.Instance{}
	for _, inst := range instances {
		serviceMap[inst.ServiceName] = append(serviceMap[inst.ServiceName], inst)
	}

	services := make([]gsvc.Service, 0, len(serviceMap))
	for _, insts := range serviceMap {
		services = append(services, NewServiceFromInstance(insts))
	}

	return services
}

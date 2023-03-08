// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"
	"fmt"
	"strings"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Search returns the service instances in memory according to the service name.
func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) ([]gsvc.Service, error) {
	if in.Prefix == "" && in.Name != "" {
		service := &Service{
			Service: gsvc.NewServiceWithName(in.Name),
		}
		in.Prefix = service.GetPrefix()
	}
	// get all instances
	instancesResponse, err := r.consumer.GetAllInstances(&polaris.GetAllInstancesRequest{
		GetAllInstancesRequest: model.GetAllInstancesRequest{
			Service:    in.Prefix,
			Namespace:  r.opt.Namespace,
			Timeout:    &r.opt.Timeout,
			RetryCount: &r.opt.RetryCount,
		},
	})
	if err != nil {
		return nil, err
	}
	serviceInstances := instancesToServiceInstances(instancesResponse.GetInstances())
	// Service filter.
	filteredServices := make([]gsvc.Service, 0)
	for _, service := range serviceInstances {
		if in.Prefix != "" && !gstr.HasPrefix(service.GetKey(), in.Prefix) {
			continue
		}
		if in.Name != "" && service.GetName() != in.Name {
			continue
		}
		if in.Version != "" && service.GetVersion() != in.Version {
			continue
		}
		if len(in.Metadata) != 0 {
			m1 := gmap.NewStrAnyMapFrom(in.Metadata)
			m2 := gmap.NewStrAnyMapFrom(service.GetMetadata())
			if !m1.IsSubOf(m2) {
				continue
			}
		}
		resultItem := service
		filteredServices = append(filteredServices, resultItem)
	}
	return filteredServices, nil
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, serviceName string) (gsvc.Watcher, error) {
	return newWatcher(ctx, r.opt.Namespace, serviceName, r.consumer)
}

func instancesToServiceInstances(instances []model.Instance) []gsvc.Service {
	serviceInstances := make([]gsvc.Service, 0, len(instances))
	for _, instance := range instances {
		if instance.IsHealthy() {
			serviceInstances = append(serviceInstances, instanceToServiceInstance(instance))
		}
	}
	return serviceInstances
}

func instanceToServiceInstance(instance model.Instance) gsvc.Service {
	var (
		s         *gsvc.LocalService
		metadata  = instance.GetMetadata()
		names     = strings.Split(instance.GetService(), instanceIDSeparator)
		endpoints = gsvc.NewEndpoints(fmt.Sprintf("%s:%d", instance.GetHost(), instance.GetPort()))
	)
	if names != nil && len(names) > 4 {
		s = &gsvc.LocalService{
			Head:       names[0],
			Deployment: names[1],
			Namespace:  names[2],
			Name:       names[3],
			Version:    metadata[metadataKeyVersion],
			Metadata:   gconv.Map(metadata),
			Endpoints:  endpoints,
		}
	} else {
		s = &gsvc.LocalService{
			Name:      instance.GetService(),
			Namespace: instance.GetNamespace(),
			Version:   metadata[metadataKeyVersion],
			Metadata:  gconv.Map(metadata),
			Endpoints: endpoints,
		}
	}
	return &Service{
		Service: s,
	}
}

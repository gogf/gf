// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/net/gsvc"
)

// Search returns the service instances in memory according to the service name.
func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) ([]gsvc.Service, error) {
	in.Separator = instanceIDSeparator
	// get all instances
	instancesResponse, err := r.consumer.GetAllInstances(&polaris.GetAllInstancesRequest{
		GetAllInstancesRequest: model.GetAllInstancesRequest{
			Service:    in.Key(),
			Namespace:  r.opt.Namespace,
			Timeout:    &r.opt.Timeout,
			RetryCount: &r.opt.RetryCount,
		},
	})
	if err != nil {
		return nil, err
	}
	serviceInstances := instancesToServiceInstances(instancesResponse.GetInstances())
	return serviceInstances, nil
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, serviceName string) (gsvc.Watcher, error) {
	return newWatcher(ctx, r.opt.Namespace, serviceName, r.consumer)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package etcd

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	etcd3 "go.etcd.io/etcd/client/v3"
)

// Search is the etcd discovery search function.
func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) ([]gsvc.Service, error) {
	if in.Prefix == "" && in.Name != "" {
		in.Prefix = gsvc.NewServiceWithName(in.Name).GetPrefix()
	}

	res, err := r.kv.Get(ctx, in.Prefix, etcd3.WithPrefix())
	if err != nil {
		return nil, err
	}
	services, err := extractResponseToServices(res)
	if err != nil {
		return nil, err
	}
	// Service filter.
	filteredServices := make([]gsvc.Service, 0)
	for _, v := range services {
		if in.Name != "" && in.Name != v.GetName() {
			continue
		}
		if in.Version != "" && in.Version != v.GetVersion() {
			continue
		}
		service := v
		filteredServices = append(filteredServices, service)
	}
	return filteredServices, nil
}

// Watch is the etcd discovery watch function.
func (r *Registry) Watch(ctx context.Context, key string) (gsvc.Watcher, error) {
	return newWatcher(key, r.client)
}

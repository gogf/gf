// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package etcd

import (
	"context"

	etcd3 "go.etcd.io/etcd/client/v3"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// Search searches and returns services with specified condition.
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
	for _, service := range services {
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

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (r *Registry) Watch(ctx context.Context, key string) (gsvc.Watcher, error) {
	return newWatcher(key, r.client)
}

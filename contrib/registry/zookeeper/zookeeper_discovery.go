// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"context"
	"path"
	"strings"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// Search searches and returns services with specified condition.
func (r *Registry) Search(_ context.Context, in gsvc.SearchInput) ([]gsvc.Service, error) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(in.Prefix, "/", "-"), "-")
	instances, err, _ := r.group.Do(prefix, func() (interface{}, error) {
		serviceNamePath := path.Join(r.opts.namespace, prefix)
		servicesID, _, err := r.conn.Children(serviceNamePath)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				"Error with search the children node under %s",
				serviceNamePath,
			)
		}
		items := make([]gsvc.Service, 0, len(servicesID))
		for _, service := range servicesID {
			servicePath := path.Join(serviceNamePath, service)
			byteData, _, err := r.conn.Get(servicePath)
			if err != nil {
				return nil, gerror.Wrapf(
					err,
					"Error with node data which name is %s",
					servicePath,
				)
			}
			item, err := unmarshal(byteData)
			if err != nil {
				return nil, gerror.Wrapf(
					err,
					"Error with unmarshal node data to Content",
				)
			}
			svc, err := gsvc.NewServiceWithKV(item.Key, item.Value)
			if err != nil {
				return nil, gerror.Wrapf(
					err,
					"Error with new service with KV in Content",
				)
			}
			items = append(items, svc)
		}
		return items, nil
	})
	if err != nil {
		return nil, gerror.Wrapf(
			err,
			"Error with group do",
		)
	}
	// Service filter.
	filteredServices := make([]gsvc.Service, 0)
	for _, service := range instances.([]gsvc.Service) {
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
	return newWatcher(ctx, r.opts.namespace, key, r.conn)
}

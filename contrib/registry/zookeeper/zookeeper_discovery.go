// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"context"
	"github.com/gogf/gf/v2/net/gsvc"
	"path"
	"strings"
)

// Search is the etcd discovery search function.
func (r *Registry) Search(_ context.Context, in gsvc.SearchInput) ([]gsvc.Service, error) {
	prefix := strings.TrimPrefix(strings.ReplaceAll(in.Prefix, "/", "-"), "-")
	instances, err, _ := r.group.Do(prefix, func() (interface{}, error) {
		serviceNamePath := path.Join(r.opts.namespace, prefix)
		servicesID, _, err := r.conn.Children(serviceNamePath)
		if err != nil {
			return nil, err
		}
		items := make([]gsvc.Service, 0, len(servicesID))
		for _, service := range servicesID {
			servicePath := path.Join(serviceNamePath, service)
			byteData, _, err := r.conn.Get(servicePath)
			if err != nil {
				return nil, err
			}
			item, err := unmarshal(byteData)
			if err != nil {
				return nil, err
			}
			svc, err := gsvc.NewServiceWithKV(item.Key, item.Value)
			if err != nil {
				return nil, err
			}
			items = append(items, svc)
		}
		return items, nil
	})
	if err != nil {
		return nil, err
	}
	return instances.([]gsvc.Service), nil
}

// Watch is the etcd discovery watch function.
func (r *Registry) Watch(ctx context.Context, key string) (gsvc.Watcher, error) {
	return newWatcher(ctx, r.opts.namespace, key, r.conn)
}

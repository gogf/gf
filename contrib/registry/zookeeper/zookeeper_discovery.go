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

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
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
	return instances.([]gsvc.Service), nil
}

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (r *Registry) Watch(ctx context.Context, key string) (gsvc.Watcher, error) {
	return newWatcher(ctx, r.opts.namespace, key, r.conn)
}

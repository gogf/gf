// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"
	"github.com/gogf/gf/v2/net/gsvc"
)

//Search searches and returns services with specified condition.
func (r Registry) Search(ctx context.Context, in gsvc.SearchInput) (result []gsvc.Service, err error) {
	return getServiceFromInstances(in.Prefix, r.opts, r.namingClient)
}

// Watch watches specified condition changes. The `key` is the prefix of service key.
func (r Registry) Watch(ctx context.Context, key string) (watcher gsvc.Watcher, err error) {
	return newWatcher(key, r.namingClient, r.opts)
}

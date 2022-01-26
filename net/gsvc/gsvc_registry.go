// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
)

var defaultRegistry Registry

// SetRegistry sets the default Registry implements as your own implemented interface.
func SetRegistry(registry Registry) {
	if registry == nil {
		panic(gerror.New(`invalid Registry value "nil" given`))
	}
	defaultRegistry = registry
}

// Register registers `service` to default registry..
func Register(ctx context.Context, service *Service) error {
	return defaultRegistry.Register(ctx, service)
}

// Deregister removes `service` from default registry.
func Deregister(ctx context.Context, service *Service) error {
	return defaultRegistry.Deregister(ctx, service)
}

// Search searches and returns services with specified condition.
func Search(ctx context.Context, in SearchInput) ([]*Service, error) {
	return defaultRegistry.Search(ctx, in)
}

// Watch watches specified condition changes.
func Watch(ctx context.Context, key string) (Watcher, error) {
	return defaultRegistry.Watch(ctx, key)
}

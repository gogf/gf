// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsvc

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gutil"
)

// watchedMap stores discovery object and its watched service mapping.
var watchedMap = gmap.New(true)

// ServiceWatch is used to watch the service status.
type ServiceWatch func(service Service)

// Get retrieves and returns the service by service name.
func Get(ctx context.Context, name string) (service Service, err error) {
	return GetAndWatchWithDiscovery(ctx, defaultRegistry, name, nil)
}

// GetWithDiscovery retrieves and returns the service by service name in `discovery`.
func GetWithDiscovery(ctx context.Context, discovery Discovery, name string) (service Service, err error) {
	return GetAndWatchWithDiscovery(ctx, discovery, name, nil)
}

// GetAndWatch is used to getting the service with custom watch callback function.
func GetAndWatch(ctx context.Context, name string, watch ServiceWatch) (service Service, err error) {
	return GetAndWatchWithDiscovery(ctx, defaultRegistry, name, watch)
}

// GetAndWatchWithDiscovery is used to getting the service with custom watch callback function in `discovery`.
func GetAndWatchWithDiscovery(ctx context.Context, discovery Discovery, name string, watch ServiceWatch) (service Service, err error) {
	if discovery == nil {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `discovery cannot be nil`)
	}
	// Retrieve service map by discovery object.
	watchedServiceMap := watchedMap.GetOrSetFunc(discovery, func() interface{} {
		return gmap.NewStrAnyMap(true)
	}).(*gmap.StrAnyMap)
	// Retrieve service by name.
	storedService := watchedServiceMap.GetOrSetFuncLock(name, func() interface{} {
		var (
			services []Service
			watcher  Watcher
		)
		services, err = discovery.Search(ctx, SearchInput{
			Name: name,
		})
		if err != nil {
			return nil
		}
		if len(services) == 0 {
			err = gerror.NewCodef(gcode.CodeNotFound, `service not found with name "%s"`, name)
			return nil
		}

		// Just pick one if multiple.
		service = services[0]

		// Watch the service changes in goroutine.
		if watch != nil {
			if watcher, err = discovery.Watch(ctx, service.GetPrefix()); err != nil {
				return nil
			}
			go watchAndUpdateService(watchedServiceMap, watcher, service, watch)
		}
		return service
	})
	if storedService != nil {
		service = storedService.(Service)
	}
	return
}

// watchAndUpdateService watches and updates the service in memory if it is changed.
func watchAndUpdateService(watchedServiceMap *gmap.StrAnyMap, watcher Watcher, service Service, watchFunc ServiceWatch) {
	var (
		ctx      = context.Background()
		err      error
		services []Service
	)
	for {
		time.Sleep(time.Second)
		if services, err = watcher.Proceed(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
			continue
		}
		if len(services) > 0 {
			watchedServiceMap.Set(service.GetName(), services[0])
			if watchFunc != nil {
				gutil.TryCatch(ctx, func(ctx context.Context) {
					watchFunc(services[0])
				}, func(ctx context.Context, exception error) {
					intlog.Errorf(ctx, `%+v`, exception)
				})
			}
		}
	}
}

// Search searches and returns services with specified condition.
func Search(ctx context.Context, in SearchInput) ([]Service, error) {
	if defaultRegistry == nil {
		return nil, gerror.NewCodef(gcode.CodeNotImplemented, `no Registry is registered`)
	}
	ctx, _ = context.WithTimeout(ctx, defaultTimeout)
	return defaultRegistry.Search(ctx, in)
}

// Watch watches specified condition changes.
func Watch(ctx context.Context, key string) (Watcher, error) {
	if defaultRegistry == nil {
		return nil, gerror.NewCodef(gcode.CodeNotImplemented, `no Registry is registered`)
	}
	return defaultRegistry.Watch(ctx, key)
}

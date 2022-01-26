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
	"github.com/gogf/gf/v2/os/glog"
)

var (
	watchedServiceMap = gmap.New(true)
)

func Get(ctx context.Context, name string) (service *Service, err error) {
	v := watchedServiceMap.GetOrSetFuncLock(name, func() interface{} {
		var (
			s        = NewServiceWithName(name)
			services []*Service
		)
		services, err = Search(ctx, SearchInput{
			Prefix:     s.Prefix,
			Deployment: s.Deployment,
			Namespace:  s.Namespace,
			Name:       s.Name,
			Version:    s.Version,
		})
		if err != nil {
			return nil
		}
		if len(services) == 0 {
			err = gerror.NewCodef(gcode.CodeNotFound, `service not found with name "%s"`, name)
			return nil
		}
		service = services[0]
		// Watch the service changes in goroutine.
		go watchAndUpdateService(ctx, service)
		return service
	})
	if v != nil {
		service = v.(*Service)
	}
	return
}

func watchAndUpdateService(ctx context.Context, service *Service) {
	var (
		err      error
		watcher  Watcher
		services []*Service
	)
	for {
		time.Sleep(time.Second)
		watcher, err = Watch(ctx, service.KeyWithoutEndpoints())
		if err != nil {
			glog.Error(ctx, err)
			continue
		}
		services, err = watcher.Proceed()
		if err != nil {
			glog.Error(ctx, err)
			continue
		}
		if len(services) > 0 {
			watchedServiceMap.Set(service.Name, services[0])
		}
	}
}

// Search searches and returns services with specified condition.
func Search(ctx context.Context, in SearchInput) ([]*Service, error) {
	if defaultRegistry == nil {
		return nil, gerror.NewCodef(gcode.CodeNotImplemented, `no Registry is registered`)
	}
	return defaultRegistry.Search(ctx, in)
}

// Watch watches specified condition changes.
func Watch(ctx context.Context, key string) (Watcher, error) {
	if defaultRegistry == nil {
		return nil, gerror.NewCodef(gcode.CodeNotImplemented, `no Registry is registered`)
	}
	return defaultRegistry.Watch(ctx, key)
}

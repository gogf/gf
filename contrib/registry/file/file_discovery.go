// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package file

import (
	"context"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfsnotify"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
)

// Search searches and returns services with specified condition.
func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) (result []gsvc.Service, err error) {
	services, err := r.getServices(ctx)
	if err != nil {
		return nil, err
	}
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
		result = append(result, resultItem)
	}
	return
}

// Watch watches specified condition changes.
// The `key` is the prefix of service key.
func (r *Registry) Watch(ctx context.Context, key string) (watcher gsvc.Watcher, err error) {
	fileWatcher := &Watcher{
		prefix: r.getServiceKeyForFile(key),
		ch:     make(chan gsvc.Service, 100),
	}
	_, err = gfsnotify.Add(r.path, func(event *gfsnotify.Event) {
		if event.IsChmod() {
			return
		}
		if !gstr.HasPrefix(gfile.Basename(event.Path), fileWatcher.prefix) {
			return
		}
		service, err := r.getServiceByFilePath(event.Path)
		if err != nil {
			return
		}
		fileWatcher.ch <- service
	})
	return fileWatcher, err
}

func (r *Registry) getServices(ctx context.Context) (services []gsvc.Service, err error) {
	filePaths, err := gfile.ScanDirFile(r.path, "*", false)
	if err != nil {
		return nil, err
	}
	for _, filePath := range filePaths {
		s, e := r.getServiceByFilePath(filePath)
		if e != nil {
			return nil, e
		}
		// Check service TTL.
		var (
			updateAt    = s.GetMetadata().Get(updateAtKey).GTime()
			nowTime     = gtime.Now()
			subDuration = nowTime.Sub(updateAt)
		)
		if updateAt.IsZero() || subDuration > serviceTTL {
			g.Log().Debugf(
				ctx,
				`service "%s" is expired, update at: %s, current: %s, sub duration: %s`,
				s.GetKey(), updateAt.String(), nowTime.String(), subDuration.String(),
			)
			continue
		}
		services = append(services, s)
	}
	return
}

func (r *Registry) getServiceByFilePath(filePath string) (gsvc.Service, error) {
	var (
		fileName    = gfile.Basename(filePath)
		fileContent = gfile.GetContents(filePath)
		serviceKey  = gstr.Replace(fileName, defaultSeparator, gsvc.DefaultSeparator)
	)
	serviceKey = gsvc.DefaultSeparator + serviceKey
	return gsvc.NewServiceWithKV(serviceKey, fileContent)
}

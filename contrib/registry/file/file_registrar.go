// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package file

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gstr"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (r *Registry) Register(ctx context.Context, service gsvc.Service) (registered gsvc.Service, err error) {
	service.GetMetadata().Set(updateAtKey, gtime.Now())
	var (
		filePath    = r.getServiceFilePath(service)
		fileContent = service.GetValue()
	)
	err = gfile.PutContents(filePath, fileContent)
	if err == nil {
		gtimer.Add(ctx, serviceUpdateInterval, func(ctx context.Context) {
			if !gfile.Exists(filePath) {
				gtimer.Exit()
			}
			// Update TTL in timer.
			service, _ = r.getServiceByFilePath(filePath)
			if service != nil {
				service.GetMetadata().Set(updateAtKey, gtime.Now())
			}
			_ = gfile.PutContents(filePath, service.GetValue())
		})
	}
	return service, err
}

// Deregister off-lines and removes `service` from the Registry.
func (r *Registry) Deregister(ctx context.Context, service gsvc.Service) error {
	return gfile.Remove(r.getServiceFilePath(service))
}

func (r *Registry) getServiceFilePath(service gsvc.Service) string {
	return gfile.Join(r.path, r.getServiceFileName(service))
}

func (r *Registry) getServiceFileName(service gsvc.Service) string {
	return r.getServiceKeyForFile(service.GetKey())
}

func (r *Registry) getServiceKeyForFile(key string) string {
	key = gstr.Replace(key, gsvc.DefaultSeparator, defaultSeparator)
	key = gstr.Trim(key, defaultSeparator)
	return key
}

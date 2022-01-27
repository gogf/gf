// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package etcd implements service Registry and Discovery using etcd.
package etcd

import (
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
)

var (
	_ gsvc.Registry = &Registry{}
	_ gsvc.Watcher  = &watcher{}
)

type Registry struct {
	client       *etcd3.Client
	kv           etcd3.KV
	lease        etcd3.Lease
	keepaliveTTL time.Duration
	logger       *glog.Logger
}

type Option struct {
	Logger       *glog.Logger
	KeepaliveTTL time.Duration
}

const (
	DefaultKeepAliveTTL = 10 * time.Second
)

func New(address string, option ...Option) *Registry {
	endpoints := gstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid etcd address "%s"`, address))
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		panic(gerror.Wrap(err, `create etcd client failed`))
	}
	return NewWithClient(client, option...)
}

func NewWithClient(client *etcd3.Client, option ...Option) *Registry {
	r := &Registry{
		client: client,
		kv:     etcd3.NewKV(client),
	}
	if len(option) > 0 {
		r.logger = option[0].Logger
		r.keepaliveTTL = option[0].KeepaliveTTL
	}
	if r.logger == nil {
		r.logger = g.Log()
	}
	if r.keepaliveTTL == 0 {
		r.keepaliveTTL = DefaultKeepAliveTTL
	}
	return r
}

// extractResponseToServices extracts etcd watch response context to service list.
func extractResponseToServices(res *etcd3.GetResponse) ([]*gsvc.Service, error) {
	if res == nil || res.Kvs == nil {
		return nil, nil
	}
	var (
		services   []*gsvc.Service
		serviceKey string
		serviceMap = make(map[string]*gsvc.Service)
	)
	for _, kv := range res.Kvs {
		service, err := gsvc.NewServiceWithKV(kv.Key, kv.Value)
		if err != nil {
			return services, err
		}
		if service != nil {
			serviceKey = service.KeyWithoutEndpoints()
			if s, ok := serviceMap[serviceKey]; ok {
				s.Endpoints = append(s.Endpoints, service.Endpoints...)
			} else {
				serviceMap[serviceKey] = service
				services = append(services, service)
			}
		}
	}
	return services, nil
}

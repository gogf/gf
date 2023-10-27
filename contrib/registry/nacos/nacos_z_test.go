// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gogf/gf/contrib/registry/nacos/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/joy999/nacos-sdk-go/common/constant"
)

const (
	NACOS_ADDRESS   = `localhost:8848`
	NACOS_CACHE_DIR = `/tmp/nacos`
	NACOS_LOG_DIR   = `/tmp/nacos`
)

func TestRegistry(t *testing.T) {
	var (
		ctx      = gctx.GetInitCtx()
		registry = nacos.New(NACOS_ADDRESS, func(cc *constant.ClientConfig) {
			cc.CacheDir = NACOS_CACHE_DIR
			cc.LogDir = NACOS_LOG_DIR
		})
	)
	svc := &gsvc.LocalService{
		Name:      guid.S(),
		Endpoints: gsvc.NewEndpoints("127.0.0.1:8888"),
		Metadata: map[string]interface{}{
			"protocol": "https",
		},
	}
	gtest.C(t, func(t *gtest.T) {
		registered, err := registry.Register(ctx, svc)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc.GetName())
	})

	// Search by name.
	gtest.C(t, func(t *gtest.T) {
		result, err := registry.Search(ctx, gsvc.SearchInput{
			Name: svc.Name,
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.Name)
	})

	// Search by prefix.
	gtest.C(t, func(t *gtest.T) {
		result, err := registry.Search(ctx, gsvc.SearchInput{
			Prefix: svc.GetPrefix(),
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.Name)
	})

	// Search by metadata.
	gtest.C(t, func(t *gtest.T) {
		result, err := registry.Search(ctx, gsvc.SearchInput{
			Name: svc.GetName(),
			Metadata: map[string]interface{}{
				"protocol": "https",
			},
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.Name)
	})
	gtest.C(t, func(t *gtest.T) {
		result, err := registry.Search(ctx, gsvc.SearchInput{
			Name: svc.GetName(),
			Metadata: map[string]interface{}{
				"protocol": "grpc",
			},
		})
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})

	gtest.C(t, func(t *gtest.T) {
		err := registry.Deregister(ctx, svc)
		t.AssertNil(err)
	})
}

func TestWatch(t *testing.T) {
	var (
		ctx      = gctx.GetInitCtx()
		registry = nacos.New(NACOS_ADDRESS, func(cc *constant.ClientConfig) {
			cc.CacheDir = NACOS_CACHE_DIR
			cc.LogDir = NACOS_LOG_DIR
		})
		registry2 = nacos.New(NACOS_ADDRESS, func(cc *constant.ClientConfig) {
			cc.CacheDir = NACOS_CACHE_DIR
			cc.LogDir = NACOS_LOG_DIR
		})
	)

	svc1 := &gsvc.LocalService{
		Name:      guid.S(),
		Endpoints: gsvc.NewEndpoints("127.0.0.1:8888"),
		Metadata: map[string]interface{}{
			"protocol": "https",
		},
	}
	gtest.C(t, func(t *gtest.T) {
		registered, err := registry.Register(ctx, svc1)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc1.GetName())
	})

	gtest.C(t, func(t *gtest.T) {
		ctx := gctx.New()
		watcher, err := registry.Watch(ctx, svc1.GetPrefix())
		t.AssertNil(err)

		var latestProceedResult atomic.Value
		g.Go(ctx, func(ctx context.Context) {
			var (
				err error
				res []gsvc.Service
			)
			for err == nil {
				res, err = watcher.Proceed()
				t.AssertNil(err)
				latestProceedResult.Store(res)
			}
		}, func(ctx context.Context, exception error) {
			t.Fatal(exception)
		})

		// Register another service.
		svc2 := &gsvc.LocalService{
			Name:      svc1.Name,
			Endpoints: gsvc.NewEndpoints("127.0.0.2:9999"),
			Metadata: map[string]interface{}{
				"protocol": "https",
			},
		}
		registered, err := registry2.Register(ctx, svc2)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc2.GetName())

		time.Sleep(time.Second * 10)

		// Watch and retrieve the service changes:
		// svc1 and svc2 is the same service name, which has 2 endpoints.
		proceedResult, ok := latestProceedResult.Load().([]gsvc.Service)
		t.Assert(ok, true)
		t.Assert(len(proceedResult), 1)
		t.Assert(
			allEndpoints(proceedResult),
			gsvc.Endpoints{svc1.GetEndpoints()[0], svc2.GetEndpoints()[0]},
		)

		// Watch and retrieve the service changes:
		// left only svc1, which means this service has only 1 endpoint.
		err = registry2.Deregister(ctx, svc2)
		t.AssertNil(err)

		time.Sleep(time.Second * 10)
		proceedResult, ok = latestProceedResult.Load().([]gsvc.Service)
		t.Assert(ok, true)
		t.Assert(len(proceedResult), 1)
		t.Assert(
			allEndpoints(proceedResult),
			gsvc.Endpoints{svc1.GetEndpoints()[0]},
		)
		t.AssertNil(watcher.Close())
	})

	gtest.C(t, func(t *gtest.T) {
		err := registry.Deregister(ctx, svc1)
		t.AssertNil(err)
	})
}

func allEndpoints(services []gsvc.Service) gsvc.Endpoints {
	m := map[gsvc.Endpoint]struct{}{}
	for _, s := range services {
		for _, ep := range s.GetEndpoints() {
			m[ep] = struct{}{}
		}
	}
	var endpoints gsvc.Endpoints
	for ep := range m {
		endpoints = append(endpoints, ep)
	}
	return sortEndpoints(endpoints)
}

func sortEndpoints(in gsvc.Endpoints) gsvc.Endpoints {
	var endpoints gsvc.Endpoints
	endpoints = append(endpoints, in...)
	n := len(endpoints)
	for i := 0; i < n; i++ {
		for t := i; t < n; t++ {
			if endpoints[i].String() > endpoints[t].String() {
				endpoints[i], endpoints[t] = endpoints[t], endpoints[i]
			}
		}
	}
	return endpoints
}

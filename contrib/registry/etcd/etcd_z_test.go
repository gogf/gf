// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package etcd_test

import (
	"testing"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func TestRegistry(t *testing.T) {
	var (
		ctx      = gctx.GetInitCtx()
		registry = etcd.New(`127.0.0.1:2379`)
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
		registry = etcd.New(`127.0.0.1:2379`)
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
		watcher, err := registry.Watch(ctx, svc1.GetPrefix())
		t.AssertNil(err)

		// Register another service.
		svc2 := &gsvc.LocalService{
			Name:      svc1.Name,
			Endpoints: gsvc.NewEndpoints("127.0.0.1:9999"),
		}
		registered, err := registry.Register(ctx, svc2)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc2.GetName())

		// Watch and retrieve the service changes:
		// svc1 and svc2 is the same service name, which has 2 endpoints.
		proceedResult, err := watcher.Proceed()
		t.AssertNil(err)
		t.Assert(len(proceedResult), 1)
		t.Assert(
			proceedResult[0].GetEndpoints(),
			gsvc.Endpoints{svc1.GetEndpoints()[0], svc2.GetEndpoints()[0]},
		)

		// Watch and retrieve the service changes:
		// left only svc1, which means this service has only 1 endpoint.
		err = registry.Deregister(ctx, svc2)
		t.AssertNil(err)
		proceedResult, err = watcher.Proceed()
		t.AssertNil(err)
		t.Assert(
			proceedResult[0].GetEndpoints(),
			gsvc.Endpoints{svc1.GetEndpoints()[0]},
		)
		t.AssertNil(watcher.Close())
	})

	gtest.C(t, func(t *gtest.T) {
		err := registry.Deregister(ctx, svc1)
		t.AssertNil(err)
	})
}

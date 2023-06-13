// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
	"time"
)

// TestRegister TestRegisterService
func TestRegister(t *testing.T) {
	var (
		registry = New([]string{"127.0.0.1:8848"},
			//you can create a namespace in nacos, then copy the namespaceId to here, "" is public namespace
			WithNameSpaceId(""),
			WithClusterName("goframe"),
			WithGroupName("goframe/nacos"),
			WithContextPath("/nacos"))
		ctx = gctx.GetInitCtx()
	)

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-0-tcp",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:8080,127.0.0.1:8088"),
		Version:   "test",
	}

	gtest.C(t, func(t *gtest.T) {
		registered, err := registry.Register(ctx, svc)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc.GetName())
	})

	gtest.C(t, func(t *gtest.T) {
		err := registry.Deregister(ctx, svc)
		t.AssertNil(err)
	})
}

// TestSearch TestSearchService
func TestSearch(t *testing.T) {
	var (
		registry = New([]string{"127.0.0.1:8848"})
		ctx      = gctx.GetInitCtx()
	)

	svc := &gsvc.LocalService{
		Name:      "goframe/provider/0/tcp",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:8080"),
		Version:   "test/gf",
	}

	gtest.C(t, func(t *gtest.T) {
		registered, err := registry.Register(ctx, svc)
		t.AssertNil(err)
		t.Assert(registered.GetName(), svc.GetName())
	})

	time.Sleep(time.Second * 1)

	gtest.C(t, func(t *gtest.T) {
		result, err := registry.Search(ctx, gsvc.SearchInput{
			Prefix: svc.GetPrefix(),
		})
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0].GetName(), svc.GetName())
	})
}

// TestWatch TestWatchService
func TestWatch(t *testing.T) {
	var (
		r   = New([]string{"127.0.0.1:8848"})
		ctx = gctx.New()
	)

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-0-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}

	watch, err := r.Watch(ctx, svc.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}
	s1, err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	// watch svc
	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		g.Log().Info(ctx, "Register Proceed service: ", instance)
	}

	err = r.Deregister(ctx, s1)
	if err != nil {
		t.Log(err)
	}

	// svc deregister, DeleteEvent
	next, err = watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output nothing
		g.Log().Info(ctx, "Deregister Proceed service: ", instance)
	}

	err = watch.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = watch.Proceed()
	if err == nil {
		// if nil, stop failed
		t.Fatal()
	}
}

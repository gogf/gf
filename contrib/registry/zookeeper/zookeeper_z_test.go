// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

// TestRegistry TestRegistryManyService
func TestRegistry(t *testing.T) {
	r := New([]string{"127.0.0.1:2181"}, WithRootPath("/gogf"))
	ctx := context.Background()

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-0-tcp",
		Version:   "test",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	r := New([]string{"127.0.0.1:2181"}, WithRootPath("/gogf"))

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-1-tcp",
		Version:   "test",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}
	svc1 := &gsvc.LocalService{
		Name:      "goframe-provider-2-tcp",
		Version:   "test",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9001"),
	}
	svc2 := &gsvc.LocalService{
		Name:      "goframe-provider-3-tcp",
		Version:   "test",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9002"),
	}

	s0, err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	s2, err := r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), s0)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), s1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), s2)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetService Test GetService
func TestGetService(t *testing.T) {
	r := New([]string{"127.0.0.1:2181"}, WithRootPath("/gogf"))
	ctx := context.Background()

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}

	s, err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	serviceInstances, err := r.Search(ctx, gsvc.SearchInput{
		Prefix:   s.GetPrefix(),
		Name:     svc.Name,
		Version:  svc.Version,
		Metadata: svc.Metadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range serviceInstances {
		g.Log().Info(ctx, instance)
	}

	err = r.Deregister(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
}

// TestServiceInstanceLeaf verifies that serviceInstanceLeaf produces a valid,
// unique ZooKeeper node name from endpoint information.
func TestServiceInstanceLeaf(t *testing.T) {
	cases := []struct {
		endpoints string
		expected  string
	}{
		{"127.0.0.1:9000", "127.0.0.1-9000"},
		{"10.0.0.1:8080", "10.0.0.1-8080"},
	}
	for _, c := range cases {
		svc := &gsvc.LocalService{
			Name:      "leaf-test-svc",
			Version:   "v1.0.0",
			Endpoints: gsvc.NewEndpoints(c.endpoints),
		}
		got := serviceInstanceLeaf(svc)
		if got != c.expected {
			t.Fatalf("serviceInstanceLeaf(%q) = %q, want %q", c.endpoints, got, c.expected)
		}
	}
}

// TestRegistryMultiInstance verifies that multiple instances of the same service
// can be registered simultaneously without overwriting each other (issue #4149).
func TestRegistryMultiInstance(t *testing.T) {
	r := New([]string{"127.0.0.1:2181"}, WithRootPath("/gogf"))
	ctx := context.Background()

	// Two instances of the same service on different ports.
	svc1 := &gsvc.LocalService{
		Name:      "multi-instance-svc",
		Version:   "v1.0.0",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:19001"),
	}
	svc2 := &gsvc.LocalService{
		Name:      "multi-instance-svc",
		Version:   "v1.0.0",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:19002"),
	}

	s1, err := r.Register(ctx, svc1)
	if err != nil {
		t.Fatalf("Register svc1 failed: %v", err)
	}

	s2, err := r.Register(ctx, svc2)
	if err != nil {
		t.Fatalf("Register svc2 failed: %v", err)
	}

	// Allow ZK to settle.
	time.Sleep(200 * time.Millisecond)

	// Both instances should be discoverable.
	services, err := r.Search(ctx, gsvc.SearchInput{
		Prefix:  svc1.GetPrefix(),
		Name:    svc1.Name,
		Version: svc1.Version,
	})
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(services) != 2 {
		t.Fatalf("expected 2 instances after multi-instance registration, got %d", len(services))
	}
	g.Log().Infof(ctx, "Found %d service instances (expected 2)", len(services))

	if err = r.Deregister(ctx, s1); err != nil {
		t.Fatalf("Deregister svc1 failed: %v", err)
	}
	if err = r.Deregister(ctx, s2); err != nil {
		t.Fatalf("Deregister svc2 failed: %v", err)
	}
}

// TestWatch Test Watch
func TestWatch(t *testing.T) {
	r := New([]string{"127.0.0.1:2181"}, WithRootPath("/gogf"))

	ctx := gctx.New()

	svc := &gsvc.LocalService{
		Name:      "goframe-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]any{"app": "goframe", gsvc.MDProtocol: "tcp"},
		Endpoints: gsvc.NewEndpoints("127.0.0.1:9000"),
	}
	t.Log("watch")
	watch, err := r.Watch(context.Background(), svc.GetPrefix())
	if err != nil {
		t.Fatal(err)
	}

	s1, err := r.Register(context.Background(), svc)
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

	err = r.Deregister(context.Background(), s1)
	if err != nil {
		t.Fatal(err)
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

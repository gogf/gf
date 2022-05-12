// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"
	"testing"
	"time"

	"github.com/polarismesh/polaris-go/pkg/config"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

// TestRegistry TestRegistryManyService
func TestRegistry(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	ctx := context.Background()

	svc := &gsvc.Service{
		Name:      "goframe-provider-0-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		Separator: _instanceIDSeparator,
	}

	err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
}

// TestRegistryMany TestRegistryManyService
func TestRegistryMany(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	svc := &gsvc.Service{
		Name:      "goframe-provider-1-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		Separator: _instanceIDSeparator,
	}
	svc1 := &gsvc.Service{
		Name:      "goframe-provider-2-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe"},
		Endpoints: []string{"tcp://127.0.0.1:9001?isSecure=false"},
		Separator: _instanceIDSeparator,
	}
	svc2 := &gsvc.Service{
		Name:      "goframe-provider-3-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe"},
		Endpoints: []string{"tcp://127.0.0.1:9002?isSecure=false"},
		Separator: _instanceIDSeparator,
	}

	err := r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), svc1)
	if err != nil {
		t.Fatal(err)
	}

	err = r.Deregister(context.Background(), svc2)
	if err != nil {
		t.Fatal(err)
	}
}

// TestGetService Test GetService
func TestGetService(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	ctx := context.Background()

	svc := &gsvc.Service{
		Name:      "goframe-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		Separator: _instanceIDSeparator,
	}

	err := r.Register(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	serviceInstances, err := r.Search(ctx, gsvc.SearchInput{
		Prefix:     svc.Prefix,
		Deployment: svc.Deployment,
		Namespace:  svc.Namespace,
		Name:       svc.Name,
		Version:    svc.Version,
		Metadata:   svc.Metadata,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range serviceInstances {
		g.Log().Info(ctx, instance)
	}

	err = r.Deregister(ctx, svc)
	if err != nil {
		t.Fatal(err)
	}
}

// TestWatch Test Watch
func TestWatch(t *testing.T) {
	conf := config.NewDefaultConfiguration([]string{"127.0.0.1:8091"})

	r := NewRegistryWithConfig(
		conf,
		WithTimeout(time.Second*10),
		WithTTL(100),
	)

	ctx := gctx.New()

	svc := &gsvc.Service{
		Name:      "goframe-provider-4-tcp",
		Version:   "test",
		Metadata:  map[string]interface{}{"app": "goframe"},
		Endpoints: []string{"tcp://127.0.0.1:9000?isSecure=false"},
		Separator: _instanceIDSeparator,
	}

	watch, err := r.Watch(context.Background(), svc.KeyWithoutEndpoints())
	if err != nil {
		t.Fatal(err)
	}

	err = r.Register(context.Background(), svc)
	if err != nil {
		t.Fatal(err)
	}
	// watch svc
	time.Sleep(time.Second * 1)

	// svc register, AddEvent
	next, err := watch.Proceed()
	if err != nil {
		t.Fatal(err)
	}
	for _, instance := range next {
		// it will output one instance
		g.Log().Info(ctx, instance)
	}

	err = r.Deregister(context.Background(), svc)
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
		g.Log().Info(ctx, instance)
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

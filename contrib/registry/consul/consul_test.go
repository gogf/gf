// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consul

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	testServiceName    = "test-service"
	testServiceVersion = "1.0.0"
	testServiceAddress = "127.0.0.1"
	testServicePort    = 8000
)

func createTestService() gsvc.Service {
	return &gsvc.LocalService{
		Name:    testServiceName,
		Version: testServiceVersion,
		Metadata: map[string]interface{}{
			"region": "cn-east-1",
			"zone":   "a",
		},
		Endpoints: []gsvc.Endpoint{
			gsvc.NewEndpoint(fmt.Sprintf("%s:%d", testServiceAddress, testServicePort)),
		},
	}
}

func Test_Registry_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)
		t.AssertNE(registry, nil)

		// Create service
		service := createTestService()

		// Register service
		ctx := context.Background()
		registeredService, err := registry.Register(ctx, service)
		t.AssertNil(err)
		t.AssertNE(registeredService, nil)

		// Search service
		time.Sleep(time.Second) // Wait for service to be registered
		services, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: testServiceVersion,
		})
		t.AssertNil(err)
		t.Assert(len(services), 1)

		foundService := services[0]
		t.Assert(foundService.GetName(), testServiceName)
		t.Assert(foundService.GetVersion(), testServiceVersion)
		t.Assert(len(foundService.GetEndpoints()), 1)

		endpoint := foundService.GetEndpoints()[0]
		t.Assert(endpoint.Host(), testServiceAddress)
		t.Assert(endpoint.Port(), testServicePort)

		metadata := foundService.GetMetadata()
		t.AssertNE(metadata, nil)
		t.Assert(metadata["region"], "cn-east-1")
		t.Assert(metadata["zone"], "a")

		// Deregister service
		err = registry.Deregister(ctx, service)
		t.AssertNil(err)

		// Verify service is deregistered
		time.Sleep(time.Second) // Wait for service to be deregistered
		services, err = registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: testServiceVersion,
		})
		t.AssertNil(err)
		t.Assert(len(services), 0)
	})
}

func Test_Registry_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)

		// Create service
		service := createTestService()

		// Register service
		ctx := context.Background()
		_, err = registry.Register(ctx, service)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service)

		// Wait for service to be registered
		time.Sleep(time.Second)

		// Create watcher
		watcher, err := registry.Watch(ctx, testServiceName)
		t.AssertNil(err)
		t.AssertNE(watcher, nil)
		defer watcher.Close()

		// Watch for service changes
		watchDone := make(chan struct{})
		go func() {
			defer close(watchDone)
			services, watchErr := watcher.Proceed()
			if watchErr != nil {
				t.Error(watchErr)
				return
			}
			if len(services) != 1 {
				t.Error("expected 1 service")
				return
			}
			watchedService := services[0]
			if watchedService.GetName() != testServiceName {
				t.Error("unexpected service name")
				return
			}
			if watchedService.GetVersion() != testServiceVersion {
				t.Error("unexpected service version")
				return
			}
		}()

		// Wait for watch event
		select {
		case <-watchDone:
			// Success
		case <-time.After(5 * time.Second):
			t.Error("Watch timeout")
		}
	})
}

func Test_Registry_Options(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with custom address
		registry, err := New(WithAddress("localhost:8500"))
		t.AssertNil(err)
		t.Assert(registry.(*Registry).GetAddress(), "localhost:8500")

		// Test with custom token
		registry, err = New(WithToken("test-token"))
		t.AssertNil(err)
		t.Assert(registry.(*Registry).options["token"], "test-token")

		// Test with both options
		registry, err = New(
			WithAddress("localhost:8500"),
			WithToken("test-token"),
		)
		t.AssertNil(err)
		t.Assert(registry.(*Registry).GetAddress(), "localhost:8500")
		t.Assert(registry.(*Registry).options["token"], "test-token")
	})
}

func Test_Registry_MultipleServices(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)

		// Create services
		service1 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "1.0.0",
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8001"),
			},
		}
		service2 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "1.0.0",
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8002"),
			},
		}

		ctx := context.Background()

		// Register services
		_, err = registry.Register(ctx, service1)
		t.AssertNil(err)
		_, err = registry.Register(ctx, service2)
		t.AssertNil(err)

		// Search services
		time.Sleep(time.Second) // Wait for services to be registered
		services, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: "1.0.0",
		})
		t.AssertNil(err)
		t.Assert(len(services), 2)

		// Verify different endpoints
		endpoints := make(map[string]bool)
		for _, service := range services {
			for _, endpoint := range service.GetEndpoints() {
				endpoints[endpoint.String()] = true
			}
		}
		t.Assert(len(endpoints), 2)
		t.Assert(endpoints["127.0.0.1:8001"], true)
		t.Assert(endpoints["127.0.0.1:8002"], true)

		// Cleanup
		err = registry.Deregister(ctx, service1)
		t.AssertNil(err)
		err = registry.Deregister(ctx, service2)
		t.AssertNil(err)
	})
}

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
		t.Assert(registry != nil, true)

		// Test invalid service
		invalidService := &gsvc.LocalService{
			Name:    testServiceName,
			Version: testServiceVersion,
		}
		_, err = registry.Register(context.Background(), invalidService)
		t.AssertNE(err, nil) // Should fail due to no endpoints

		// Create service with invalid metadata
		serviceWithInvalidMeta := &gsvc.LocalService{
			Name:    testServiceName,
			Version: testServiceVersion,
			Metadata: map[string]interface{}{
				"invalid": make(chan int), // This will fail JSON marshaling
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint(fmt.Sprintf("%s:%d", testServiceAddress, testServicePort)),
			},
		}
		_, err = registry.Register(context.Background(), serviceWithInvalidMeta)
		t.AssertNE(err, nil) // Should fail due to invalid metadata

		// Create service
		service := createTestService()

		// Register service
		ctx := context.Background()
		registeredService, err := registry.Register(ctx, service)
		t.AssertNil(err)
		t.Assert(registeredService != nil, true)

		// Wait for service to be registered
		time.Sleep(2 * time.Second)

		// Search service
		services, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: testServiceVersion,
		})
		t.AssertNil(err)
		t.Assert(len(services), 1)

		// Test service properties
		foundService := services[0]
		t.Assert(foundService.GetName(), testServiceName)
		t.Assert(foundService.GetVersion(), testServiceVersion)
		t.Assert(len(foundService.GetEndpoints()), 1)

		endpoint := foundService.GetEndpoints()[0]
		t.Assert(endpoint.Host(), testServiceAddress)
		t.Assert(endpoint.Port(), testServicePort)

		metadata := foundService.GetMetadata()
		t.Assert(metadata != nil, true)
		t.Assert(metadata["region"], "cn-east-1")
		t.Assert(metadata["zone"], "a")

		// Search with invalid metadata
		servicesWithInvalidMeta, err := registry.Search(ctx, gsvc.SearchInput{
			Name:     testServiceName,
			Version:  testServiceVersion,
			Metadata: map[string]interface{}{"nonexistent": "value"},
		})
		t.AssertNil(err)
		t.Assert(len(servicesWithInvalidMeta), 0)

		// Test deregister with invalid service
		err = registry.Deregister(ctx, invalidService)
		t.AssertNE(err, nil) // Should fail due to no endpoints

		// Deregister service
		err = registry.Deregister(ctx, service)
		t.AssertNil(err)

		// Wait for service to be deregistered
		time.Sleep(2 * time.Second)

		// Verify service is deregistered
		deregisteredServices, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: testServiceVersion,
		})
		t.AssertNil(err)
		t.Assert(len(deregisteredServices), 0)
	})
}

func Test_Registry_Watch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)

		// Create service
		service := createTestService()

		// Register service first
		ctx := context.Background()
		_, err = registry.Register(ctx, service)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service)

		// Wait for service to be registered
		time.Sleep(time.Second)

		// Create watcher after service is registered
		watcher, err := registry.Watch(ctx, testServiceName)
		t.AssertNil(err)
		t.Assert(watcher != nil, true)
		defer watcher.Close()

		// Wait for initial service query
		time.Sleep(time.Second)

		// Should receive initial service list
		services, err := watcher.Proceed()
		t.AssertNil(err)
		t.Assert(len(services), 1)
		t.Assert(services[0].GetName(), testServiceName)
		t.Assert(services[0].GetVersion(), testServiceVersion)

		// Test closing watcher
		err = watcher.Close()
		t.AssertNil(err)

		// Test watch with invalid service name
		watcher, err = registry.Watch(ctx, "nonexistent-service")
		t.AssertNil(err)
		defer watcher.Close()

		// Wait for initial query
		time.Sleep(time.Second)

		// Should receive empty service list for non-existent service
		services, err = watcher.Proceed()
		t.AssertNil(err)
		t.Assert(len(services), 0)

		// Test watch after service deregistration
		watcher, err = registry.Watch(ctx, testServiceName)
		t.AssertNil(err)
		defer watcher.Close()

		// Wait for initial query
		time.Sleep(time.Second)

		err = registry.Deregister(ctx, service)
		t.AssertNil(err)

		// Wait for service to be deregistered
		time.Sleep(time.Second)

		// Should receive empty service list after deregistration
		services, err = watcher.Proceed()
		t.AssertNil(err)
		t.Assert(len(services), 0)
	})
}

func Test_Registry_MultipleServices(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)

		// Create multiple services
		service1 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "1.0.0",
			Metadata: map[string]interface{}{
				"region": "us-east-1",
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8001"),
			},
		}

		service2 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "2.0.0",
			Metadata: map[string]interface{}{
				"region": "us-west-1",
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8002"),
			},
		}

		// Register services
		ctx := context.Background()
		_, err = registry.Register(ctx, service1)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service1)

		_, err = registry.Register(ctx, service2)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service2)

		// Wait for services to be registered
		time.Sleep(2 * time.Second)

		// Search all services without version filter
		allServices, err := registry.Search(ctx, gsvc.SearchInput{
			Name: testServiceName,
		})
		t.AssertNil(err)
		t.Assert(len(allServices), 2)

		// Test search with different versions
		services1, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: "1.0.0",
		})
		t.AssertNil(err)
		t.Assert(len(services1), 1)
		t.Assert(services1[0].GetVersion(), "1.0.0")

		services2, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: "2.0.0",
		})
		t.AssertNil(err)
		t.Assert(len(services2), 1)
		t.Assert(services2[0].GetVersion(), "2.0.0")

		// Test search with metadata
		servicesEast, err := registry.Search(ctx, gsvc.SearchInput{
			Name: testServiceName,
			Metadata: map[string]interface{}{
				"region": "us-east-1",
			},
		})
		t.AssertNil(err)
		t.Assert(len(servicesEast), 1)
		t.Assert(servicesEast[0].GetMetadata()["region"], "us-east-1")

		// Watch both services
		watcher, err := registry.Watch(ctx, testServiceName)
		t.AssertNil(err)
		defer watcher.Close()

		// Wait for initial query
		time.Sleep(time.Second)

		// Should receive updates for both services
		services, err := watcher.Proceed()
		t.AssertNil(err)
		t.Assert(len(services), 2)

		// Verify services are sorted by version
		t.Assert(services[0].GetVersion() < services[1].GetVersion(), true)
	})
}

func Test_Registry_Options(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test with custom address
		registry1, err := New(WithAddress("localhost:8500"))
		t.AssertNil(err)
		t.Assert(registry1.(*Registry).GetAddress(), "localhost:8500")

		// Test with token
		registry2, err := New(WithAddress("localhost:8500"), WithToken("test-token"))
		t.AssertNil(err)
		t.Assert(registry2.(*Registry).options["token"], "test-token")

		// Test with invalid address (should still create registry but fail on operations)
		registry3, err := New(WithAddress("invalid:99999"))
		t.AssertNil(err)
		_, err = registry3.Register(context.Background(), createTestService())
		t.AssertNE(err, nil)
	})
}

func Test_Registry_MultipleServicesMetadataFiltering(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)

		// Create multiple services
		service1 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "1.0.0",
			Metadata: map[string]interface{}{
				"region": "us-east-1",
				"env":    "dev",
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8001"),
			},
		}

		service2 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "2.0.0",
			Metadata: map[string]interface{}{
				"region": "us-west-1",
				"env":    "prod",
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8002"),
			},
		}

		// Register services
		ctx := context.Background()
		_, err = registry.Register(ctx, service1)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service1)

		_, err = registry.Register(ctx, service2)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service2)

		time.Sleep(time.Second) // Wait for services to be registered

		// Test search with metadata filtering
		servicesDev, err := registry.Search(ctx, gsvc.SearchInput{
			Name: testServiceName,
			Metadata: map[string]interface{}{
				"env": "dev",
			},
		})
		t.AssertNil(err)
		t.Assert(len(servicesDev), 1)
		t.Assert(servicesDev[0].GetMetadata()["env"], "dev")

		servicesProd, err := registry.Search(ctx, gsvc.SearchInput{
			Name: testServiceName,
			Metadata: map[string]interface{}{
				"env": "prod",
			},
		})
		t.AssertNil(err)
		t.Assert(len(servicesProd), 1)
		t.Assert(servicesProd[0].GetMetadata()["env"], "prod")
	})
}

func Test_Registry_MultipleServicesVersionFiltering(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Create registry
		registry, err := New()
		t.AssertNil(err)

		// Create multiple services
		service1 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "1.0.0",
			Metadata: map[string]interface{}{
				"region": "us-east-1",
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8001"),
			},
		}

		service2 := &gsvc.LocalService{
			Name:    testServiceName,
			Version: "2.0.0",
			Metadata: map[string]interface{}{
				"region": "us-west-1",
			},
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint("127.0.0.1:8002"),
			},
		}

		// Register services
		ctx := context.Background()
		_, err = registry.Register(ctx, service1)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service1)

		_, err = registry.Register(ctx, service2)
		t.AssertNil(err)
		defer registry.Deregister(ctx, service2)

		time.Sleep(time.Second) // Wait for services to be registered

		// Test search with version filtering
		services, err := registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: "1.0.0",
		})
		t.AssertNil(err)
		t.Assert(len(services), 1)
		t.Assert(services[0].GetVersion(), "1.0.0")

		services, err = registry.Search(ctx, gsvc.SearchInput{
			Name:    testServiceName,
			Version: "2.0.0",
		})
		t.AssertNil(err)
		t.Assert(len(services), 1)
		t.Assert(services[0].GetVersion(), "2.0.0")
	})
}

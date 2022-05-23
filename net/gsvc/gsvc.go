// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsvc provides service registry and discovery definition.
package gsvc

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Registry interface for service.
type Registry interface {
	Registrar
	Discovery
}

// Registrar interface for service registrar.
type Registrar interface {
	// Register registers `service` to Registry.
	// Note that it returns a new Service if it changes the input Service with custom one.
	Register(ctx context.Context, service Service) (Service, error)

	// Deregister off-lines and removes `service` from the Registry.
	Deregister(ctx context.Context, service Service) error
}

// Discovery interface for service discovery.
type Discovery interface {
	// Search searches and returns services with specified condition.
	Search(ctx context.Context, in SearchInput) ([]Service, error)

	// Watch watches specified condition changes.
	Watch(ctx context.Context, key string) (Watcher, error)
}

// Watcher interface for service.
type Watcher interface {
	// Proceed proceeds watch in blocking way.
	Proceed() ([]Service, error)

	// Close closes the watcher.
	Close() error
}

// Service interface for service definition.
type Service interface {
	// GetName returns the name of the service.
	// The name is necessary for a service, and should be unique among services.
	GetName() string

	// GetVersion returns the version of the service.
	// It is suggested using GNU version naming like: v1.0.0, v2.0.1, v2.1.0-rc.
	// A service can have multiple versions deployed at once.
	// If no version set in service, the default version of service is "latest".
	GetVersion() string

	// GetKey formats and returns a unique key string for service.
	// The result key is commonly used for key-value registrar server.
	GetKey() string

	// GetValue formats and returns the value of the service.
	// The result value is commonly used for key-value registrar server.
	GetValue() string

	// GetPrefix formats and returns the key prefix string.
	// The result prefix string is commonly used in key-value registrar server
	// for service searching.
	//
	// Take etcd server for example, the prefix string is used like:
	// `etcdctl get /services/prod/hello.svc --prefix`
	GetPrefix() string

	// GetMetadata returns the Metadata map of service.
	// The Metadata is key-value pair map specifying extra attributes of a service.
	GetMetadata() Metadata

	// GetEndpoints returns the Endpoints of service.
	// The Endpoints contain multiple host/port information of service.
	GetEndpoints() Endpoints
}

// Endpoint interface for service.
type Endpoint interface {
	// Host returns the IPv4/IPv6 address of a service.
	Host() string

	// Port returns the port of a service.
	Port() int

	// String formats and returns the Endpoint as a string.
	String() string
}

// Endpoints are composed by multiple Endpoint.
type Endpoints []Endpoint

// Metadata stores custom key-value pairs.
type Metadata map[string]interface{}

// SearchInput is the input for service searching.
type SearchInput struct {
	Prefix   string   // Search by key prefix.
	Name     string   // Search by service name.
	Version  string   // Search by service version.
	Metadata Metadata // Filter by metadata if there are multiple result.
}

const (
	Schema                    = `services`
	DefaultHead               = `services`
	DefaultDeployment         = `default`
	DefaultNamespace          = `default`
	DefaultVersion            = `latest`
	EnvPrefix                 = `GF_GSVC_PREFIX`
	EnvDeployment             = `GF_GSVC_DEPLOYMENT`
	EnvNamespace              = `GF_GSVC_NAMESPACE`
	EnvName                   = `GF_GSVC_Name`
	EnvVersion                = `GF_GSVC_VERSION`
	MDProtocol                = `protocol`
	MDInsecure                = `insecure`
	MDWeight                  = `weight`
	DefaultProtocol           = `http`
	DefaultSeparator          = "/"
	defaultTimeout            = 5 * time.Second
	endpointHostPortDelimiter = ":"
	endpointsDelimiter        = ","
)

var defaultRegistry Registry

// SetRegistry sets the default Registry implements as your own implemented interface.
func SetRegistry(registry Registry) {
	if registry == nil {
		panic(gerror.New(`invalid Registry value "nil" given`))
	}
	defaultRegistry = registry
}

// GetRegistry returns the default Registry that is previously set.
// It returns nil if no Registry is set.
func GetRegistry() Registry {
	return defaultRegistry
}

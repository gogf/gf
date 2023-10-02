// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
)

// Service used to record and manage service information
type Service struct {
	// Ip          string                 `json:"ip"`
	// Port        int                    `json:"port"`
	Endpoints   gsvc.Endpoints         `json:"endpoints"`
	ServiceName string                 `json:"serviceName"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewService new one service.
func NewService() gsvc.Service {
	return &Service{
		Metadata: map[string]interface{}{},
	}
}

// NewServiceFromInstance new one service from instance
func NewServiceFromInstance(instance []model.Instance) gsvc.Service {
	n := len(instance)
	if n == 0 {
		return nil
	}
	serviceName := instance[0].ServiceName
	endpoints := make(gsvc.Endpoints, 0, n)
	for i := 0; i < n; i++ {
		if instance[0].ServiceName != serviceName {
			return nil
		}
		endpoints = append(endpoints, NewEndpoint(instance[i].Ip, int(instance[i].Port)))
	}
	if gstr.Contains(serviceName, "@@") {
		arr := gstr.SplitAndTrim(serviceName, "@@")
		serviceName = arr[1]
	}
	return &Service{
		Endpoints:   endpoints,
		ServiceName: serviceName,
		Metadata:    gmap.NewStrStrMapFrom(instance[0].Metadata).MapStrAny(),
	}
}

// NewServicesFromInstances new some services from some instances
func NewServicesFromInstances(instances []model.Instance) []gsvc.Service {
	serviceMap := map[string][]model.Instance{}
	for _, inst := range instances {
		serviceMap[inst.ServiceName] = append(serviceMap[inst.ServiceName], inst)
	}

	services := make([]gsvc.Service, 0, len(serviceMap))
	for _, insts := range serviceMap {
		services = append(services, NewServiceFromInstance(insts))
	}

	return services
}

// GetName returns the name of the service.
// The name is necessary for a service, and should be unique among services.
func (s *Service) GetName() string {
	return s.ServiceName
}

// GetVersion returns the version of the service.
// It is suggested using GNU version naming like: v1.0.0, v2.0.1, v2.1.0-rc.
// A service can have multiple versions deployed at once.
// If no version set in service, the default version of service is "latest".
func (s *Service) GetVersion() string {
	return "latest"
}

// GetKey formats and returns a unique key string for service.
// The result key is commonly used for key-value registrar server.
func (s *Service) GetKey() string {
	return ""
}

// GetValue formats and returns the value of the service.
// The result value is commonly used for key-value registrar server.
func (s *Service) GetValue() string {
	return ""
}

// GetPrefix formats and returns the key prefix string.
// The result prefix string is commonly used in key-value registrar server
// for service searching.
//
// Take etcd server for example, the prefix string is used like:
// `etcdctl get /services/prod/hello.svc --prefix`
func (s *Service) GetPrefix() string {
	return ""
}

// GetMetadata returns the Metadata map of service.
// The Metadata is key-value pair map specifying extra attributes of a service.
func (s *Service) GetMetadata() gsvc.Metadata {
	return s.Metadata
}

// GetEndpoints returns the Endpoints of service.
// The Endpoints contain multiple host/port information of service.
func (s *Service) GetEndpoints() gsvc.Endpoints {
	return s.Endpoints
}

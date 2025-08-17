// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package file

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/net/gsvc"
)

// Service wrapper.
type Service struct {
	gsvc.Service
	Endpoints gsvc.Endpoints
	Metadata  gsvc.Metadata
}

// NewService creates and returns local Service from gsvc.Service interface object.
func NewService(service gsvc.Service) *Service {
	s, ok := service.(*Service)
	if ok {
		if s.Endpoints == nil {
			s.Endpoints = make(gsvc.Endpoints, 0)
		}
		if s.Metadata == nil {
			s.Metadata = make(gsvc.Metadata)
		}
		return s
	}
	s = &Service{
		Service:   service,
		Endpoints: make(gsvc.Endpoints, 0),
		Metadata:  make(gsvc.Metadata),
	}
	if len(service.GetEndpoints()) > 0 {
		s.Endpoints = service.GetEndpoints()
	}
	if len(service.GetMetadata()) > 0 {
		s.Metadata = service.GetMetadata()
	}
	return s
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

// GetValue formats and returns the value of the service.
// The result value is commonly used for key-value registrar server.
func (s *Service) GetValue() string {
	b, _ := gjson.Marshal(s.Metadata)
	return string(b)
}

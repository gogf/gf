// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gsvc provides service registry and discovery definition.
package gsvc

import (
	"github.com/gogf/gf/v2/text/gstr"
)

// NewEndpoints creates and returns Endpoints from multiple addresses like:
// "192.168.1.100:80,192.168.1.101:80".
func NewEndpoints(addresses string) Endpoints {
	endpoints := make([]Endpoint, 0)
	for _, address := range gstr.SplitAndTrim(addresses, EndpointsDelimiter) {
		endpoints = append(endpoints, NewEndpoint(address))
	}
	return endpoints
}

// String formats and returns the Endpoints as a string like:
// "192.168.1.100:80,192.168.1.101:80"
func (es Endpoints) String() string {
	var s string
	for _, endpoint := range es {
		if s != "" {
			s += EndpointsDelimiter
		}
		s += endpoint.String()
	}
	return s
}

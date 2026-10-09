// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package zookeeper

import (
	"encoding/json"

	"github.com/gogf/gf/v2/net/gsvc"
)

func unmarshal(data []byte) (c *Content, err error) {
	err = json.Unmarshal(data, &c)
	return
}

func marshal(c *Content) ([]byte, error) {
	return json.Marshal(c)
}

// instanceNodeName returns the name of the node that stores a single instance of `service`.
// It uses the endpoints, so that instances of the same service do not overwrite each other.
func instanceNodeName(service gsvc.Service) string {
	if endpoints := service.GetEndpoints().String(); endpoints != "" {
		return endpoints
	}
	return service.GetName()
}

// mergeServices merges instances having the same service prefix into one service with all their endpoints.
func mergeServices(services []gsvc.Service) []gsvc.Service {
	var (
		merged   = make([]gsvc.Service, 0, len(services))
		prefixes = make(map[string]*gsvc.LocalService)
	)
	for _, service := range services {
		if s, ok := prefixes[service.GetPrefix()]; ok {
			s.Endpoints = append(s.Endpoints, service.GetEndpoints()...)
			continue
		}
		if s, ok := service.(*gsvc.LocalService); ok {
			prefixes[s.GetPrefix()] = s
		}
		merged = append(merged, service)
	}
	return merged
}

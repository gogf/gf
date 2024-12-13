// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consul

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
)

// Search searches and returns services with specified condition.
func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) ([]gsvc.Service, error) {
	// Get services from consul
	services, _, err := r.client.Health().Service(in.Name, "", true, &api.QueryOptions{
		WaitTime: time.Second * 3,
	})
	if err != nil {
		return nil, gerror.Wrap(err, "failed to get services from consul")
	}

	var result []gsvc.Service
	for _, service := range services {
		if service.Checks.AggregatedStatus() != api.HealthPassing {
			continue
		}

		// Parse metadata
		var metadata map[string]interface{}
		if metaStr, ok := service.Service.Meta["metadata"]; ok && metaStr != "" {
			if err = json.Unmarshal([]byte(metaStr), &metadata); err != nil {
				return nil, gerror.Wrap(err, "failed to unmarshal service metadata")
			}
		}

		// Skip if version doesn't match
		if in.Version != "" {
			if len(service.Service.Tags) == 0 || service.Service.Tags[0] != in.Version {
				continue
			}
		}

		// Skip if metadata doesn't match
		if len(in.Metadata) > 0 {
			if metadata == nil {
				continue
			}
			match := true
			for k, v := range in.Metadata {
				if mv, ok := metadata[k]; !ok || mv != v {
					match = false
					break
				}
			}
			if !match {
				continue
			}
		}

		// Get version from tags
		version := ""
		if len(service.Service.Tags) > 0 {
			version = service.Service.Tags[0]
		}

		// Create service instance
		localService := &gsvc.LocalService{
			Head:       "",
			Deployment: "",
			Namespace:  "",
			Name:       service.Service.Service,
			Version:    version,
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint(fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)),
			},
			Metadata: metadata,
		}
		result = append(result, localService)
	}

	return result, nil
}

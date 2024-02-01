// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

type discoveryNode struct {
	service gsvc.Service
	address string
}

// Service is the client discovery service.
func (n *discoveryNode) Service() gsvc.Service {
	return n.service
}

// Address returns the address of the node.
func (n *discoveryNode) Address() string {
	return n.address
}

// service prefix to its selector map cache.
var clientSelectorMap = gmap.New(true)

// internalMiddlewareDiscovery is a client middleware that enables service discovery feature for client.
func internalMiddlewareDiscovery(c *Client, r *http.Request) (response *Response, err error) {
	if c.discovery == nil && !isServiceName(r.URL.Host) {
		return c.Next(r)
	}
	var (
		ctx     = r.Context()
		service gsvc.Service
	)
	service, err = gsvc.GetAndWatchWithDiscovery(ctx, c.discovery, r.URL.Host, func(service gsvc.Service) {
		intlog.Printf(ctx, `http client watching service "%s" changed`, service.GetPrefix())
		if v := clientSelectorMap.Get(service.GetPrefix()); v != nil {
			if err = updateSelectorNodesByService(ctx, v.(gsel.Selector), service); err != nil {
				intlog.Errorf(ctx, `%+v`, err)
			}
		}
	})
	if err != nil {
		if gerror.Code(err) == gcode.CodeNotFound {
			intlog.Printf(
				ctx,
				`service discovery error with url "%s:%s":%s`,
				r.Method, r.URL.String(), err.Error(),
			)
			return c.Next(r)
		}
	}
	if service == nil {
		return c.Next(r)
	}
	// Balancer.
	var (
		selectorMapKey   = service.GetPrefix()
		selectorMapValue = clientSelectorMap.GetOrSetFuncLock(selectorMapKey, func() interface{} {
			intlog.Printf(ctx, `http client create selector for service "%s"`, selectorMapKey)
			selector := c.builder.Build()
			// Update selector nodes.
			if err = updateSelectorNodesByService(ctx, selector, service); err != nil {
				return nil
			}
			return selector
		})
	)
	if err != nil {
		return nil, err
	}
	selector := selectorMapValue.(gsel.Selector)
	// Pick one node from multiple addresses.
	node, done, err := selector.Pick(ctx)
	if err != nil {
		return nil, err
	}
	if done != nil {
		defer done(ctx, gsel.DoneInfo{})
	}
	r.Host = node.Address()
	r.URL.Host = node.Address()
	return c.Next(r)
}

func updateSelectorNodesByService(ctx context.Context, selector gsel.Selector, service gsvc.Service) error {
	nodes := make(gsel.Nodes, 0)
	for _, endpoint := range service.GetEndpoints() {
		nodes = append(nodes, &discoveryNode{
			service: service,
			address: endpoint.String(),
		})
	}
	return selector.Update(ctx, nodes)
}

// isServiceName checks and returns whether given input parameter is service name or not.
// It checks by whether the parameter is address by containing port delimiter character ':'.
//
// It does not contain any port number if using service discovery.
func isServiceName(serviceNameOrAddress string) bool {
	return !gstr.Contains(serviceNameOrAddress, gsvc.EndpointHostPortDelimiter)
}

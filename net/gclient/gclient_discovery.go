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
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	discoveryMiddlewareHandled gctx.StrKey = `MiddlewareClientDiscoveryHandled`
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

var clientSelectorMap = gmap.New(true)

// internalMiddlewareDiscovery is a client middleware that enables service discovery feature for client.
func internalMiddlewareDiscovery(c *Client, r *http.Request) (response *Response, err error) {
	ctx := r.Context()
	// Mark this request is handled by server tracing middleware,
	// to avoid repeated handling by the same middleware.
	if ctx.Value(discoveryMiddlewareHandled) != nil {
		return c.Next(r)
	}
	if gsvc.GetRegistry() == nil {
		return c.Next(r)
	}
	var service gsvc.Service
	service, err = gsvc.GetAndWatch(ctx, r.URL.Host, func(service gsvc.Service) {
		intlog.Printf(ctx, `http client watching service "%s" changed`, service.GetPrefix())
		if v := clientSelectorMap.Get(service.GetPrefix()); v != nil {
			if err = updateSelectorNodesByService(v.(gsel.Selector), service); err != nil {
				intlog.Errorf(context.Background(), `%+v`, err)
			}
		}
	})
	if err != nil {
		return nil, err
	}
	if service == nil {
		return c.Next(r)
	}
	// Balancer.
	selectorMapKey := service.GetPrefix()
	selector := clientSelectorMap.GetOrSetFuncLock(selectorMapKey, func() interface{} {
		intlog.Printf(ctx, `http client create selector for service "%s"`, selectorMapKey)
		return gsel.GetBuilder().Build()
	}).(gsel.Selector)
	// Update selector nodes.
	if err = updateSelectorNodesByService(selector, service); err != nil {
		return nil, err
	}
	// Pick one node from multiple addresses.
	node, done, err := selector.Pick(ctx)
	if err != nil {
		return nil, err
	}
	if done != nil {
		defer done(ctx, gsel.DoneInfo{})
	}
	r.URL.Host = node.Address()
	r.Host = node.Address()
	return c.Next(r)
}

func updateSelectorNodesByService(selector gsel.Selector, service gsvc.Service) error {
	nodes := make([]gsel.Node, 0)
	for _, endpoint := range service.GetEndpoints() {
		nodes = append(nodes, &discoveryNode{
			service: service,
			address: endpoint.String(),
		})
	}
	return selector.Update(nodes)
}

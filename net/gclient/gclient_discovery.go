// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"net/http"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
)

const (
	discoveryMiddlewareHandled gctx.StrKey = `MiddlewareClientDiscoveryHandled`
)

type discoveryNode struct {
	service *gsvc.Service
	address string
}

func (n *discoveryNode) Service() *gsvc.Service {
	return n.service
}

func (n *discoveryNode) Address() string {
	return n.address
}

var (
	clientSelectorMap = gmap.New(true)
)

// internalMiddlewareDiscovery is a client middleware that enables service discovery feature for client.
func internalMiddlewareDiscovery(c *Client, r *http.Request) (response *Response, err error) {
	var ctx = r.Context()
	// Mark this request is handled by server tracing middleware,
	// to avoid repeated handling by the same middleware.
	if ctx.Value(discoveryMiddlewareHandled) != nil {
		return c.Next(r)
	}
	if gsvc.GetRegistry() == nil {
		return c.Next(r)
	}
	var service *gsvc.Service
	service, err = gsvc.GetWithWatch(ctx, r.URL.Host, func(service *gsvc.Service) {
		intlog.Printf(ctx, `http client watching service "%s" changed`, service.KeyWithoutEndpoints())
		// If service changed, it removes it from map cache,
		// which makes it re-cache later.
		clientSelectorMap.Remove(service.KeyWithoutEndpoints())
	})
	if err != nil {
		return nil, err
	}
	if service == nil {
		return c.Next(r)
	}
	// Balancer.
	selector := clientSelectorMap.GetOrSetFuncLock(
		service.KeyWithoutEndpoints(),
		func() interface{} {
			intlog.Printf(ctx, `http client create selector for service "%s"`, service.KeyWithoutEndpoints())
			// Build selector and cache it in internal map.
			nodes := make([]gsel.Node, 0)
			for _, address := range service.Endpoints {
				nodes = append(nodes, &discoveryNode{
					service: service,
					address: address,
				})
			}
			selector := gsel.GetBuilder().Build()
			if err = selector.Update(nodes); err != nil {
				glog.Error(ctx, err)
			}
			return selector
		},
	).(gsel.Selector)
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

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gclient

import (
	"net/http"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	discoveryMiddlewareHandled gctx.StrKey = `MiddlewareClientDiscoveryHandled`
)

// internalMiddlewareDiscovery is a client middleware that enables service discovery feature for client.
func internalMiddlewareDiscovery(c *Client, r *http.Request) (response *Response, err error) {
	var ctx = r.Context()
	// Mark this request is handled by server tracing middleware,
	// to avoid repeated handling by the same middleware.
	if ctx.Value(discoveryMiddlewareHandled) != nil {
		return c.Next(r)
	}
	if gsvc.GetRegistry() != nil {
		service, err := gsvc.Get(ctx, r.URL.Host)
		if err != nil {
			return nil, err
		}
		if service != nil {
			r.URL.Host = service.Address()
			r.Host = service.Address()
		}
	}
	return c.Next(r)
}

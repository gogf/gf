// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gsvc"
)

type Endpoint struct {
	host string
	port int
}

// NewEndpoint new one Endpoint instance.
func NewEndpoint(host string, port int) gsvc.Endpoint {
	return &Endpoint{host: host, port: port}
}

// Host returns the IPv4/IPv6 address of a service.
func (ep *Endpoint) Host() string {
	return ep.host
}

// Port returns the port of a service.
func (ep *Endpoint) Port() int {
	return ep.port
}

// String formats and returns the Endpoint as a string.
func (ep *Endpoint) String() string {
	return fmt.Sprintf("%v:%v", ep.host, ep.port)
}

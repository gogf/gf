// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package balancer

import (
	"google.golang.org/grpc/balancer"

	"github.com/gogf/gf/v2/net/gsvc"
)

// Node is the node for the balancer.
type Node struct {
	service gsvc.Service
	conn    balancer.SubConn
}

// Service returns the service of the node.
func (n *Node) Service() gsvc.Service {
	return n.service
}

// Address returns the address of the node.
func (n *Node) Address() string {
	endpoints := n.service.GetEndpoints()
	if len(endpoints) == 0 {
		return ""
	}
	return endpoints[0].String()
}

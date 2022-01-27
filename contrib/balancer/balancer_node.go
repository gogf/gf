// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package balancer

import (
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/balancer"
)

type Node struct {
	service *gsvc.Service
	conn    balancer.SubConn
}

func (n *Node) Service() *gsvc.Service {
	return n.service
}

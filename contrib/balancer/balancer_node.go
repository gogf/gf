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

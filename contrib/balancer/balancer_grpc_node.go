package balancer

import (
	"google.golang.org/grpc/balancer"

	"github.com/gogf/gf/v2/net/gsvc"
)

type grpcNode struct {
	service *gsvc.Service
	subConn balancer.SubConn
}

func (n *grpcNode) Service() *gsvc.Service {
	return n.service
}

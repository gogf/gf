package grpc

import (
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/balancer"
)

type grpcNode struct {
	service *gsvc.Service
	subConn balancer.SubConn
}

func (n *grpcNode) Service() *gsvc.Service {
	return n.service
}

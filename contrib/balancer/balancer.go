package balancer

import (
	"github.com/gogf/gf/v2/net/gsel"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	RawSvcKeyInSubConnInfo = `RawService`
)

func Register(name string, selector gsel.Selector) {
	balancer.Register(base.NewBalancerBuilder(
		name,
		&Builder{selector: selector},
		base.Config{HealthCheck: true},
	))
}

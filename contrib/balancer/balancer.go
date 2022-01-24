package grpc

import (
	"sync"

	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

var (
	mu sync.Mutex
)

const (
	rawSvcKeyInSubConnInfo = `RawService`
)

// Builder is grpc balancer builder.
type Builder struct {
	selector gsel.Selector
}

// SetGlobalBalancer set grpc balancer with scheme.
func SetGlobalBalancer(scheme string) {
	mu.Lock()
	defer mu.Unlock()

	b := base.NewBalancerBuilder(
		scheme,
		&Builder{selector: nil},
		base.Config{HealthCheck: true},
	)
	balancer.Register(b)
}

// Build creates a grpc Picker.
func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	nodes := make([]gsel.Node, len(info.ReadySCs))
	for conn, subConnInfo := range info.ReadySCs {
		svc, _ := subConnInfo.Address.Attributes.Value(rawSvcKeyInSubConnInfo).(*gsvc.Service)
		nodes = append(nodes, &grpcNode{
			service: svc,
			subConn: conn,
		})
	}
	p := &Picker{
		selector: b.selector,
	}
	p.selector.Update(nodes)
	return p
}

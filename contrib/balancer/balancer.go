package balancer

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"

	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
)

const (
	rawSvcKeyInSubConnInfo = `RawService`
)

type Builder struct {
	selector gsel.Selector
}

func Register(name string, selector gsel.Selector) {
	balancer.Register(base.NewBalancerBuilder(
		name,
		&Builder{selector: selector},
		base.Config{HealthCheck: true},
	))
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

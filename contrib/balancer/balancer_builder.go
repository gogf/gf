package balancer

import (
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

// Builder implements grpc balancer base.PickerBuilder,
// which returns a picker that will be used by gRPC to pick a SubConn.
type Builder struct {
	selector gsel.Selector
}

// Build returns a picker that will be used by gRPC to pick a SubConn.
func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	nodes := make([]gsel.Node, len(info.ReadySCs))
	for conn, subConnInfo := range info.ReadySCs {
		svc, _ := subConnInfo.Address.Attributes.Value(RawSvcKeyInSubConnInfo).(*gsvc.Service)
		nodes = append(nodes, &Node{
			service: svc,
			conn:    conn,
		})
	}
	p := &Picker{
		selector: b.selector,
	}
	err := p.selector.Update(nodes)
	if err != nil {
		panic(err)
	}
	return p
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package balancer

import (
	"context"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
)

// Builder implements grpc balancer base.PickerBuilder,
// which returns a picker that will be used by gRPC to pick a SubConn.
type Builder struct {
	builder gsel.Builder
}

// Build returns a picker that will be used by gRPC to pick a SubConn.
func (b *Builder) Build(info base.PickerBuildInfo) balancer.Picker {
	if len(info.ReadySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}
	var (
		ctx   = context.Background()
		nodes = make([]gsel.Node, 0)
	)
	for conn, subConnInfo := range info.ReadySCs {
		svc, _ := subConnInfo.Address.Attributes.Value(rawSvcKeyInSubConnInfo).(gsvc.Service)
		if svc == nil {
			g.Log().Noticef(ctx, `empty service read from: %+v`, subConnInfo.Address)
			continue
		}
		nodes = append(nodes, &Node{
			service: svc,
			conn:    conn,
		})
	}
	p := &Picker{
		selector: b.builder.Build(),
	}
	if err := p.selector.Update(ctx, nodes); err != nil {
		panic(err)
	}
	return p
}

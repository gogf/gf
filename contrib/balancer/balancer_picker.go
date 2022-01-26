package balancer

import (
	"github.com/gogf/gf/v2/net/gsel"
	"google.golang.org/grpc/balancer"
)

type Picker struct {
	selector gsel.Selector
}

func (p *Picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	node, done, err := p.selector.Pick(info.Ctx)
	if err != nil {
		return balancer.PickResult{}, err
	}
	return balancer.PickResult{
		SubConn: node.(*grpcNode).subConn,
		Done: func(di balancer.DoneInfo) {
			done(info.Ctx, gsel.DoneInfo{
				Err:           di.Err,
				Trailer:       di.Trailer,
				BytesSent:     di.BytesSent,
				BytesReceived: di.BytesReceived,
				ServerLoad:    di.ServerLoad,
			})
		},
	}, nil
}

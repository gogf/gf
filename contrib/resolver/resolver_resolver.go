package resolver

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/contrib/balancer"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

type Resolver struct {
	watcher gsvc.Watcher
	cc      resolver.ClientConn
	ctx     context.Context
	cancel  context.CancelFunc
	logger  *glog.Logger
}

func (r *Resolver) watch() {
	for {
		select {
		case <-r.ctx.Done():
			return

		default:
			services, err := r.watcher.Proceed()
			if err != nil {
				r.logger.Warningf(r.ctx, `watcher.Proceed error: %+v`, err)
				time.Sleep(time.Second)
				continue
			}
			r.update(services)
		}
	}
}

func (r *Resolver) update(services []*gsvc.Service) {
	var (
		err       error
		addresses = make([]resolver.Address, 0)
	)
	for _, service := range services {
		for _, endpoint := range service.Endpoints {
			addr := resolver.Address{
				Addr:       endpoint,
				ServerName: service.Name,
				Attributes: newAttributesFromMetadata(service.Metadata),
			}
			addr.Attributes = addr.Attributes.WithValue(balancer.RawSvcKeyInSubConnInfo, service)
			addresses = append(addresses, addr)
		}
	}
	if len(addresses) == 0 {
		r.logger.Noticef(r.ctx, "empty addresses parsed from: %+v", services)
		return
	}
	if err = r.cc.UpdateState(resolver.State{Addresses: addresses}); err != nil {
		r.logger.Errorf(r.ctx, "UpdateState failed: %+v", err)
	}
}

func (r *Resolver) Close() {
	if err := r.watcher.Close(); err != nil {
		r.logger.Errorf(r.ctx, `%+v`, err)
	}
	r.cancel()
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func newAttributesFromMetadata(metadata map[string]interface{}) *attributes.Attributes {
	var a *attributes.Attributes
	for k, v := range metadata {
		if a == nil {
			a = attributes.New(k, v)
		} else {
			a = a.WithValue(k, v)
		}
	}
	return a
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package resolver

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Resolver implements grpc resolver.Resolver,
// which watches for the updates on the specified target.
// Updates include address updates and service config updates.
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
			if err != nil && !errors.Is(err, context.Canceled) {
				r.logger.Warningf(r.ctx, `watcher.Proceed error: %+v`, err)
				time.Sleep(time.Second)
				continue
			}
			if len(services) > 0 {
				r.update(services)
			}
		}
	}
}

func (r *Resolver) update(services []gsvc.Service) {
	var (
		err       error
		addresses = make([]resolver.Address, 0)
	)
	for _, service := range services {
		for _, endpoint := range service.GetEndpoints() {
			addr := resolver.Address{
				Addr:       endpoint.String(),
				ServerName: service.GetName(),
				Attributes: newAttributesFromMetadata(service.GetMetadata()),
			}
			addr.Attributes = addr.Attributes.WithValue(rawSvcKeyInSubConnInfo, service)
			addresses = append(addresses, addr)
		}
	}
	if len(addresses) == 0 {
		r.logger.Noticef(r.ctx, "empty addresses parsed from: %+v", services)
		return
	}
	r.logger.Debugf(r.ctx, "client conn updated with addresses %+v", addresses)
	if err = r.cc.UpdateState(resolver.State{Addresses: addresses}); err != nil {
		r.logger.Errorf(r.ctx, "UpdateState failed: %+v", err)
	}
}

// Close closes the resolver.
func (r *Resolver) Close() {
	r.logger.Debugf(r.ctx, `resolver closed`)
	if err := r.watcher.Close(); err != nil {
		r.logger.Errorf(r.ctx, `%+v`, err)
	}
	r.cancel()
}

// ResolveNow will be called by gRPC to try to resolve the target name
// again. It's just a hint, resolver can ignore this if it's not necessary.
//
// It could be called multiple times concurrently.
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

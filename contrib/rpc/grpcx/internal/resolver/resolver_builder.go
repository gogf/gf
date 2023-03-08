// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package resolver

import (
	"context"

	"google.golang.org/grpc/resolver"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

// Builder is the builder for the etcd discovery resolver.
type Builder struct {
	discovery gsvc.Discovery
}

// NewBuilder creates and returns a Builder.
func NewBuilder(discovery gsvc.Discovery) *Builder {
	return &Builder{
		discovery: discovery,
	}
}

// Build creates a new etcd discovery resolver.
func (b *Builder) Build(
	target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	var (
		err         error
		watcher     gsvc.Watcher
		watchKey    = target.URL.Path
		ctx, cancel = context.WithCancel(gctx.GetInitCtx())
	)
	g.Log().Debugf(ctx, `Watch key "%s"`, watchKey)
	if watcher, err = b.discovery.Watch(ctx, watchKey); err != nil {
		cancel()
		return nil, gerror.Wrap(err, `registry.Watch failed`)
	}
	r := &Resolver{
		discovery: b.discovery,
		watcher:   watcher,
		watchKey:  watchKey,
		cc:        cc,
		ctx:       ctx,
		cancel:    cancel,
		logger:    g.Log(),
	}
	go r.watch()
	return r, nil
}

// Scheme return scheme of discovery
func (*Builder) Scheme() string {
	return gsvc.Schema
}

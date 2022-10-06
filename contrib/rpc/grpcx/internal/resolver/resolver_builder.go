// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package resolver

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/resolver"
)

// Builder is the builder for the etcd discovery resolver.
type Builder struct{}

// Build creates a new etcd discovery resolver.
func (b *Builder) Build(
	target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	var (
		err         error
		watcher     gsvc.Watcher
		ctx, cancel = context.WithCancel(context.Background())
	)
	g.Log().Debugf(ctx, `etcd Watch key "%s"`, target.URL.Path)
	if watcher, err = gsvc.Watch(ctx, target.URL.Path); err != nil {
		cancel()
		return nil, gerror.Wrap(err, `registry.Watch failed`)
	}
	r := &Resolver{
		watcher: watcher,
		cc:      cc,
		ctx:     ctx,
		cancel:  cancel,
		logger:  g.Log(),
	}
	go r.watch()
	return r, nil
}

// Scheme return scheme of discovery
func (*Builder) Scheme() string {
	return gsvc.Schema
}

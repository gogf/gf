package resolver

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"google.golang.org/grpc/resolver"
)

type Builder struct{}

func (b *Builder) Build(
	target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	var (
		err         error
		watcher     gsvc.Watcher
		ctx, cancel = context.WithCancel(context.Background())
	)
	glog.Debugf(ctx, `etcd Watch key "%s"`, target.URL.Path)
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

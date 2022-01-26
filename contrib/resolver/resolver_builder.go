package resolver

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"google.golang.org/grpc/resolver"
)

const Name = "katyusha"

type Builder struct {
	registry gsvc.Registry
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(registry gsvc.Registry) resolver.Builder {
	b := &Builder{
		registry: registry,
	}
	return b
}

func (b *Builder) Build(
	target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	var (
		err         error
		watcher     gsvc.Watcher
		ctx, cancel = context.WithCancel(context.Background())
	)
	if watcher, err = b.registry.Watch(ctx, target.URL.Path); err != nil {
		cancel()
		return nil, gerror.Wrap(err, `registry.Watch failed`)
	}
	r := &Resolver{
		watcher: watcher,
		cc:      cc,
		ctx:     ctx,
		cancel:  cancel,
	}
	go r.watch()
	return r, nil
}

// Scheme return scheme of discovery
func (*Builder) Scheme() string {
	return Name
}

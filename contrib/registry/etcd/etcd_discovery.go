package etcd

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	etcd3 "go.etcd.io/etcd/client/v3"
)

func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) ([]*gsvc.Service, error) {
	res, err := r.kv.Get(ctx, in.Key(), etcd3.WithPrefix())
	if err != nil {
		return nil, err
	}
	services, err := extractResponseToServices(res)
	if err != nil {
		return nil, err
	}
	// Service filter.
	filteredServices := make([]*gsvc.Service, 0)
	for _, v := range services {
		if in.Deployment != "" && in.Deployment != v.Deployment {
			continue
		}
		if in.Namespace != "" && in.Namespace != v.Namespace {
			continue
		}
		if in.Name != "" && in.Name != v.Name {
			continue
		}
		if in.Version != "" && in.Version != v.Version {
			continue
		}
		service := v
		filteredServices = append(filteredServices, service)
	}
	return filteredServices, nil
}

func (r *Registry) Watch(ctx context.Context, key string) (gsvc.Watcher, error) {
	return newWatcher(ctx, key, r.client)
}

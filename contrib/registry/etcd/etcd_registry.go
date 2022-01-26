package etcd

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
)

func New(address string, option ...Option) (*Registry, error) {
	endpoints := gstr.SplitAndTrim(address, ",")
	if len(endpoints) == 0 {
		return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid etcd address "%s"`, address)
	}
	client, err := etcd3.New(etcd3.Config{
		Endpoints: endpoints,
	})
	if err != nil {
		return nil, gerror.Wrap(err, `create etcd client failed`)
	}
	return NewWithClient(client, option...), nil
}

func NewWithClient(client *etcd3.Client, option ...Option) *Registry {
	r := &Registry{
		client: client,
		kv:     etcd3.NewKV(client),
	}
	if len(option) > 0 {
		r.logger = option[0].Logger
		r.keepaliveTTL = option[0].KeepaliveTTL
	}
	if r.keepaliveTTL == 0 {
		r.keepaliveTTL = DefaultKeepAliveTTL
	}
	return r
}

func (r *Registry) Register(ctx context.Context, service *gsvc.Service) error {
	r.lease = etcd3.NewLease(r.client)
	grant, err := r.lease.Grant(ctx, int64(r.keepaliveTTL.Seconds()))
	if err != nil {
		return err
	}
	_, err = r.client.Put(ctx, service.Key(), service.Value(), etcd3.WithLease(grant.ID))
	if err != nil {
		return err
	}
	keepAliceCh, err := r.client.KeepAlive(ctx, grant.ID)
	if err != nil {
		return err
	}
	go r.doKeepAlive(ctx, grant.ID, keepAliceCh)
	return nil
}

func (r *Registry) Deregister(ctx context.Context, service *gsvc.Service) error {
	defer func() {
		if r.lease != nil {
			_ = r.lease.Close()
		}
	}()
	_, err := r.client.Delete(ctx, service.Key())
	return err
}

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

// doKeepAlive continuously keeps alive the lease from ETCD.
func (r *Registry) doKeepAlive(
	ctx context.Context, leaseID etcd3.LeaseID, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse,
) {
	for {
		select {
		case <-r.client.Ctx().Done():
			r.logger.Debugf(ctx, "keepalive done for lease id: %d", leaseID)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				r.logger.Debugf(ctx, `keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				r.logger.Debugf(ctx, `keepalive exit, lease id: %d`, leaseID)
				return
			}
		}
	}
}

// extractResponseToServices extracts etcd watch response context to service list.
func extractResponseToServices(res *etcd3.GetResponse) ([]*gsvc.Service, error) {
	var services []*gsvc.Service
	if res == nil || res.Kvs == nil {
		return services, nil
	}
	for _, kv := range res.Kvs {
		service, err := gsvc.NewServiceFromKV(kv.Key, kv.Value)
		if err != nil {
			return services, err
		}
		if service != nil {
			services = append(services, service)
		}
	}
	return services, nil
}

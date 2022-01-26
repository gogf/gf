package etcd

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	etcd3 "go.etcd.io/etcd/client/v3"
)

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

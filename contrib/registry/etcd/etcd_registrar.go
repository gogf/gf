// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package etcd

import (
	"context"
	"time"

	etcd3 "go.etcd.io/etcd/client/v3"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/grand"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (r *Registry) Register(ctx context.Context, service gsvc.Service) (gsvc.Service, error) {
	service = NewService(service)
	if err := r.doRegisterLease(ctx, service); err != nil {
		return nil, err
	}
	return service, nil
}

func (r *Registry) doRegisterLease(ctx context.Context, service gsvc.Service) error {
	r.lease = etcd3.NewLease(r.client)

	ctx, cancel := context.WithTimeout(context.Background(), r.etcdConfig.DialTimeout)
	defer cancel()

	grant, err := r.lease.Grant(ctx, int64(r.keepaliveTTL.Seconds()))
	if err != nil {
		return gerror.Wrapf(err, `etcd grant failed with keepalive ttl "%s"`, r.keepaliveTTL)
	}
	var (
		key   = service.GetKey()
		value = service.GetValue()
	)
	_, err = r.client.Put(ctx, key, value, etcd3.WithLease(grant.ID))
	if err != nil {
		return gerror.Wrapf(
			err,
			`etcd put failed with key "%s", value "%s", lease "%d"`,
			key, value, grant.ID,
		)
	}
	r.logger.Debugf(
		ctx,
		`etcd put success with key "%s", value "%s", lease "%d"`,
		key, value, grant.ID,
	)
	keepAliceCh, err := r.client.KeepAlive(context.Background(), grant.ID)
	if err != nil {
		return err
	}
	go r.doKeepAlive(service, grant.ID, keepAliceCh)
	return nil
}

// Deregister off-lines and removes `service` from the Registry.
func (r *Registry) Deregister(ctx context.Context, service gsvc.Service) error {
	_, err := r.client.Delete(ctx, service.GetKey())
	if r.lease != nil {
		_ = r.lease.Close()
	}
	return err
}

// doKeepAlive continuously keeps alive the lease from ETCD.
func (r *Registry) doKeepAlive(
	service gsvc.Service, leaseID etcd3.LeaseID, keepAliceCh <-chan *etcd3.LeaseKeepAliveResponse,
) {
	var ctx = context.Background()
	for {
		select {
		case <-r.client.Ctx().Done():
			r.logger.Infof(ctx, "keepalive done for lease id: %d", leaseID)
			return

		case res, ok := <-keepAliceCh:
			if res != nil {
				// r.logger.Debugf(ctx, `keepalive loop: %v, %s`, ok, res.String())
			}
			if !ok {
				r.logger.Warningf(ctx, `keepalive exit, lease id: %d, retry register`, leaseID)
				// Re-register the service.
				for {
					if err := r.doRegisterLease(ctx, service); err != nil {
						retryDuration := grand.D(time.Second, time.Second*3)
						r.logger.Errorf(
							ctx,
							`keepalive retry register failed, will retry in %s: %+v`,
							retryDuration, err,
						)
						time.Sleep(retryDuration)
						continue
					}
					break
				}
				return
			}
		}
	}
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"
	"time"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Register the registration.
func (r *Registry) Register(ctx context.Context, service gsvc.Service) error {
	// Replace input service to custom service type.
	service = &Service{
		Service: service,
	}
	// Register logic.
	var (
		ids            = make([]string, 0, len(service.GetEndpoints()))
		serviceVersion = service.GetVersion()
	)
	for _, endpoint := range service.GetEndpoints() {
		// medata
		var rmd map[string]interface{}
		if service.GetMetadata().IsEmpty() {
			rmd = map[string]interface{}{
				metadataKeyKind:    gsvc.DefaultProtocol,
				metadataKeyVersion: service.GetVersion(),
			}
		} else {
			rmd = make(map[string]interface{}, len(service.GetMetadata())+2)
			rmd[metadataKeyKind] = gsvc.DefaultProtocol
			if protocol, ok := service.GetMetadata()[gsvc.MDProtocol]; ok {
				rmd[metadataKeyKind] = gconv.String(protocol)
			}
			rmd[metadataKeyVersion] = serviceVersion
			for k, v := range service.GetMetadata() {
				rmd[k] = v
			}
		}
		// Register
		registeredService, err := r.provider.Register(
			&polaris.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      service.GetPrefix(),
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         endpoint.Host(),
					Port:         endpoint.Port(),
					Protocol:     r.opt.Protocol,
					Weight:       &r.opt.Weight,
					Priority:     &r.opt.Priority,
					Version:      &serviceVersion,
					Metadata:     gconv.MapStrStr(rmd),
					Healthy:      &r.opt.Healthy,
					Isolate:      &r.opt.Isolate,
					TTL:          &r.opt.TTL,
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			})
		if err != nil {
			return err
		}
		if r.opt.Heartbeat {
			r.doHeartBeat(ctx, registeredService.InstanceID, service, endpoint)
		}
		ids = append(ids, registeredService.InstanceID)
	}
	// need to set InstanceID for Deregister
	service.(*Service).ID = gstr.Join(ids, instanceIDSeparator)
	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, service gsvc.Service) error {
	r.c <- struct{}{}
	var (
		err   error
		split = gstr.Split(service.(*Service).ID, instanceIDSeparator)
	)
	for i, endpoint := range service.GetEndpoints() {
		// Deregister
		err = r.provider.Deregister(
			&polaris.InstanceDeRegisterRequest{
				InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
					Service:      service.GetPrefix(),
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					InstanceID:   split[i],
					Host:         endpoint.Host(),
					Port:         endpoint.Port(),
					Timeout:      &r.opt.Timeout,
					RetryCount:   &r.opt.RetryCount,
				},
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) doHeartBeat(ctx context.Context, instanceID string, service gsvc.Service, endpoint gsvc.Endpoint) {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(r.opt.TTL))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := r.provider.Heartbeat(&polaris.InstanceHeartbeatRequest{
					InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
						Service:      service.GetPrefix(),
						Namespace:    r.opt.Namespace,
						Host:         endpoint.Host(),
						Port:         endpoint.Port(),
						ServiceToken: r.opt.ServiceToken,
						InstanceID:   instanceID,
						Timeout:      &r.opt.Timeout,
						RetryCount:   &r.opt.RetryCount,
					},
				})
				if err != nil {
					g.Log().Error(ctx, err.Error())
					continue
				}
			case <-r.c:
				g.Log().Debug(ctx, "stop heartbeat")
				return
			}
		}
	}()
}

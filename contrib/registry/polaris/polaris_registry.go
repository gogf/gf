// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Register the registration.
func (r *Registry) Register(ctx context.Context, service gsvc.Service) (gsvc.Service, error) {
	// Replace input service to custom service types.
	service = &Service{
		Service: service,
	}
	// Register logic.
	var ids = make([]string, 0, len(service.GetEndpoints()))
	for _, endpoint := range service.GetEndpoints() {
		// medata
		var (
			rmd            map[string]any
			serviceName    = service.GetPrefix()
			serviceVersion = service.GetVersion()
		)
		if service.GetMetadata().IsEmpty() {
			rmd = map[string]any{
				metadataKeyKind:    gsvc.DefaultProtocol,
				metadataKeyVersion: serviceVersion,
			}
		} else {
			rmd = make(map[string]any, len(service.GetMetadata())+2)
			rmd[metadataKeyKind] = gsvc.DefaultProtocol
			if protocol, ok := service.GetMetadata()[gsvc.MDProtocol]; ok {
				rmd[metadataKeyKind] = gconv.String(protocol)
			}
			rmd[metadataKeyVersion] = serviceVersion
			for k, v := range service.GetMetadata() {
				rmd[k] = v
			}
		}
		// Register RegisterInstance Service registration is performed synchronously,
		// and heartbeat reporting is automatically performed
		registeredService, err := r.provider.RegisterInstance(
			&polaris.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      serviceName,
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
			return nil, err
		}
		ids = append(ids, registeredService.InstanceID)
	}
	// need to set InstanceID for Deregister
	service.(*Service).ID = gstr.Join(ids, instanceIDSeparator)
	return service, nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, service gsvc.Service) error {
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

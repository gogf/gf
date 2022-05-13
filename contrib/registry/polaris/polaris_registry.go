// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/gconv"
)

// Register the registration.
func (r *Registry) Register(ctx context.Context, serviceInstance *gsvc.Service) error {
	ids := make([]string, 0, len(serviceInstance.Endpoints))
	// set separator
	serviceInstance.Separator = instanceIDSeparator
	for _, endpoint := range serviceInstance.Endpoints {
		httpArr := strings.Split(endpoint, endpointDelimiter)
		if len(httpArr) < 2 {
			return errors.New("invalid endpoint")
		}
		host := httpArr[0]

		// port to int
		portNum, err := strconv.Atoi(httpArr[1])
		if err != nil {
			return err
		}

		// medata
		var rmd map[string]interface{}
		if serviceInstance.Metadata == nil {
			rmd = map[string]interface{}{
				"kind":    gsvc.DefaultProtocol,
				"version": serviceInstance.Version,
			}
		} else {
			rmd = make(map[string]interface{}, len(serviceInstance.Metadata)+2)
			if protocol, ok := serviceInstance.Metadata[gsvc.MDProtocol].(string); !ok {
				rmd["kind"] = gsvc.DefaultProtocol
			} else {
				rmd["kind"] = protocol
			}
			rmd["version"] = serviceInstance.Version
			for k, v := range serviceInstance.Metadata {
				rmd[k] = v
			}
		}
		// Register
		service, err := r.provider.Register(
			&api.InstanceRegisterRequest{
				InstanceRegisterRequest: model.InstanceRegisterRequest{
					Service:      serviceInstance.KeyWithoutEndpoints(),
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					Host:         host,
					Port:         portNum,
					Protocol:     r.opt.Protocol,
					Weight:       &r.opt.Weight,
					Priority:     &r.opt.Priority,
					Version:      &serviceInstance.Version,
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
		instanceID := service.InstanceID

		if r.opt.Heartbeat {
			// start heartbeat report
			go func() {
				ticker := time.NewTicker(time.Second * time.Duration(r.opt.TTL))
				defer ticker.Stop()

				for {
					<-ticker.C

					err = r.provider.Heartbeat(&api.InstanceHeartbeatRequest{
						InstanceHeartbeatRequest: model.InstanceHeartbeatRequest{
							Service:      serviceInstance.KeyWithoutEndpoints(),
							Namespace:    r.opt.Namespace,
							Host:         host,
							Port:         portNum,
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
				}
			}()
		}

		ids = append(ids, instanceID)
	}
	// need to set InstanceID for Deregister
	serviceInstance.ID = strings.Join(ids, instanceIDSeparator)

	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, serviceInstance *gsvc.Service) error {
	split := strings.Split(serviceInstance.ID, instanceIDSeparator)
	serviceInstance.Separator = instanceIDSeparator
	for i, endpoint := range serviceInstance.Endpoints {
		// get url
		u, err := url.Parse(endpoint)
		if err != nil {
			return err
		}

		// get host and port
		host, port, err := net.SplitHostPort(u.Host)
		if err != nil {
			return err
		}

		// port to int
		portNum, err := strconv.Atoi(port)
		if err != nil {
			return err
		}
		// Deregister
		err = r.provider.Deregister(
			&api.InstanceDeRegisterRequest{
				InstanceDeRegisterRequest: model.InstanceDeRegisterRequest{
					Service:      serviceInstance.KeyWithoutEndpoints(),
					ServiceToken: r.opt.ServiceToken,
					Namespace:    r.opt.Namespace,
					InstanceID:   split[i],
					Host:         host,
					Port:         portNum,
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

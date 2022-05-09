// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	_ gsvc.Registrar = &Registry{}
)

// _instanceIDSeparator Instance id Separator.
const _instanceIDSeparator = "-"

type options struct {
	// required, namespace in polaris
	Namespace string

	// required, service access token
	ServiceToken string

	// optional, protocol in polaris. Default value is nil, it means use protocol config in service
	Protocol *string

	// service weight in polaris. Default value is 100, 0 <= weight <= 10000
	Weight int

	// service priority. Default value is 0. The smaller the value, the lower the priority
	Priority int

	// To show service is healthy or not. Default value is True.
	Healthy bool

	// Heartbeat enable .Not in polaris . Default value is True.
	Heartbeat bool

	// To show service is isolate or not. Default value is False.
	Isolate bool

	// TTL timeout. if node needs to use heartbeat to report,required. If not set,server will throw ErrorCode-400141
	TTL int

	// optional, Timeout for single query. Default value is global config
	// Total is (1+RetryCount) * Timeout
	Timeout time.Duration

	// optional, retry count. Default value is global config
	RetryCount int
}

// Option The option is a polaris option.
type Option func(o *options)

// Registry is polaris registry.
type Registry struct {
	opt      options
	provider api.ProviderAPI
	consumer api.ConsumerAPI
}

// WithNamespace with the Namespace option.
func WithNamespace(namespace string) Option {
	return func(o *options) { o.Namespace = namespace }
}

// WithServiceToken with ServiceToken option.
func WithServiceToken(serviceToken string) Option {
	return func(o *options) { o.ServiceToken = serviceToken }
}

// WithProtocol with the Protocol option.
func WithProtocol(protocol string) Option {
	return func(o *options) { o.Protocol = &protocol }
}

// WithWeight with the Weight option.
func WithWeight(weight int) Option {
	return func(o *options) { o.Weight = weight }
}

// WithHealthy with the Healthy option.
func WithHealthy(healthy bool) Option {
	return func(o *options) { o.Healthy = healthy }
}

// WithIsolate with the Isolate option.
func WithIsolate(isolate bool) Option {
	return func(o *options) { o.Isolate = isolate }
}

// WithTTL with the TTL option.
func WithTTL(TTL int) Option {
	return func(o *options) { o.TTL = TTL }
}

// WithTimeout the Timeout option.
func WithTimeout(timeout time.Duration) Option {
	return func(o *options) { o.Timeout = timeout }
}

// WithRetryCount with RetryCount option.
func WithRetryCount(retryCount int) Option {
	return func(o *options) { o.RetryCount = retryCount }
}

// WithHeartbeat with the Heartbeat option.
func WithHeartbeat(heartbeat bool) Option {
	return func(o *options) { o.Heartbeat = heartbeat }
}

// NewRegistry 创建一个注册服务
func NewRegistry(provider api.ProviderAPI, consumer api.ConsumerAPI, opts ...Option) (r *Registry) {
	op := options{
		Namespace:    "default",
		ServiceToken: "",
		Protocol:     nil,
		Weight:       0,
		Priority:     0,
		Healthy:      true,
		Heartbeat:    true,
		Isolate:      false,
		TTL:          0,
		Timeout:      0,
		RetryCount:   0,
	}
	for _, option := range opts {
		option(&op)
	}
	return &Registry{
		opt:      op,
		provider: provider,
		consumer: consumer,
	}
}

// NewRegistryWithConfig new a registry with config.
func NewRegistryWithConfig(conf config.Configuration, opts ...Option) (r *Registry) {
	provider, err := api.NewProviderAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	consumer, err := api.NewConsumerAPIByConfig(conf)
	if err != nil {
		panic(err)
	}
	return NewRegistry(provider, consumer, opts...)
}

// Register the registration.
func (r *Registry) Register(ctx context.Context, serviceInstance *gsvc.Service) error {
	ids := make([]string, 0, len(serviceInstance.Endpoints))
	// set separator
	serviceInstance.Separator = _instanceIDSeparator
	for _, endpoint := range serviceInstance.Endpoints {
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

		// medata
		var rmd map[string]interface{}
		if serviceInstance.Metadata == nil {
			rmd = map[string]interface{}{
				"kind":    u.Scheme,
				"version": serviceInstance.Version,
			}
		} else {
			rmd = make(map[string]interface{}, len(serviceInstance.Metadata)+2)
			for k, v := range serviceInstance.Metadata {
				rmd[k] = v
			}
			rmd["kind"] = u.Scheme
			rmd["version"] = serviceInstance.Version
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
	serviceInstance.ID = strings.Join(ids, _instanceIDSeparator)

	return nil
}

// Deregister the registration.
func (r *Registry) Deregister(ctx context.Context, serviceInstance *gsvc.Service) error {
	split := strings.Split(serviceInstance.ID, _instanceIDSeparator)
	serviceInstance.Separator = _instanceIDSeparator
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

// Search returns the service instances in memory according to the service name.
func (r *Registry) Search(ctx context.Context, in gsvc.SearchInput) ([]*gsvc.Service, error) {
	in.Separator = _instanceIDSeparator
	// get all instances
	instancesResponse, err := r.consumer.GetAllInstances(&api.GetAllInstancesRequest{
		GetAllInstancesRequest: model.GetAllInstancesRequest{
			Service:    in.Key(),
			Namespace:  r.opt.Namespace,
			Timeout:    &r.opt.Timeout,
			RetryCount: &r.opt.RetryCount,
		},
	})
	if err != nil {
		return nil, err
	}

	serviceInstances := instancesToServiceInstances(instancesResponse.GetInstances())

	return serviceInstances, nil
}

// Watch creates a watcher according to the service name.
func (r *Registry) Watch(ctx context.Context, serviceName string) (gsvc.Watcher, error) {
	return newWatcher(ctx, r.opt.Namespace, serviceName, r.consumer)
}

// Watcher is a service watcher.
type Watcher struct {
	ServiceName      string
	Namespace        string
	Ctx              context.Context
	Cancel           context.CancelFunc
	Channel          <-chan model.SubScribeEvent
	ServiceInstances []*gsvc.Service
}

func newWatcher(ctx context.Context, namespace string, serviceName string, consumer api.ConsumerAPI) (*Watcher, error) {
	watchServiceResponse, err := consumer.WatchService(&api.WatchServiceRequest{
		WatchServiceRequest: model.WatchServiceRequest{
			Key: model.ServiceKey{
				Namespace: namespace,
				Service:   serviceName,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		Namespace:        namespace,
		ServiceName:      serviceName,
		Channel:          watchServiceResponse.EventChannel,
		ServiceInstances: instancesToServiceInstances(watchServiceResponse.GetAllInstancesResp.GetInstances()),
	}
	w.Ctx, w.Cancel = context.WithCancel(ctx)
	return w, nil
}

// Proceed returns services in the following two cases:
// 1.the first time to watch and the service instance list is not empty.
// 2.any service instance changes found.
// if the above two conditions are not met, it will block until the context deadline is exceeded or canceled
func (w *Watcher) Proceed() ([]*gsvc.Service, error) {
	select {
	case <-w.Ctx.Done():
		return nil, w.Ctx.Err()
	case event := <-w.Channel:
		if event.GetSubScribeEventType() == model.EventInstance {
			// these are always true, but we need to check it to make sure EventType not change
			if instanceEvent, ok := event.(*model.InstanceEvent); ok {
				// handle DeleteEvent
				if instanceEvent.DeleteEvent != nil {
					for _, instance := range instanceEvent.DeleteEvent.Instances {
						for i, serviceInstance := range w.ServiceInstances {
							if serviceInstance.ID == instance.GetId() {
								// remove equal
								if len(w.ServiceInstances) <= 1 {
									w.ServiceInstances = w.ServiceInstances[0:0]
									continue
								}
								w.ServiceInstances = append(w.ServiceInstances[:i], w.ServiceInstances[i+1:]...)
							}
						}
					}
				}
				// handle UpdateEvent
				if instanceEvent.UpdateEvent != nil {
					for i, serviceInstance := range w.ServiceInstances {
						for _, update := range instanceEvent.UpdateEvent.UpdateList {
							if serviceInstance.ID == update.Before.GetId() {
								w.ServiceInstances[i] = instanceToServiceInstance(update.After)
							}
						}
					}
				}
				// handle AddEvent
				if instanceEvent.AddEvent != nil {
					w.ServiceInstances = append(w.ServiceInstances, instancesToServiceInstances(instanceEvent.AddEvent.Instances)...)
				}
			}
			return w.ServiceInstances, nil
		}
	}
	return w.ServiceInstances, nil
}

// Close the watcher.
func (w *Watcher) Close() error {
	w.Cancel()
	return nil
}

func instancesToServiceInstances(instances []model.Instance) []*gsvc.Service {
	serviceInstances := make([]*gsvc.Service, 0, len(instances))
	for _, instance := range instances {
		if instance.IsHealthy() {
			serviceInstances = append(serviceInstances, instanceToServiceInstance(instance))
		}
	}
	return serviceInstances
}

func instanceToServiceInstance(instance model.Instance) *gsvc.Service {
	metadata := instance.GetMetadata()
	// Usually, it won't fail in goframe if register correctly
	kind := ""
	if k, ok := metadata["kind"]; ok {
		kind = k
	}
	return &gsvc.Service{
		Name:      instance.GetService(),
		Version:   metadata["version"],
		Metadata:  gconv.Map(metadata),
		Endpoints: []string{fmt.Sprintf("%s://%s:%d", kind, instance.GetHost(), instance.GetPort())},
		Separator: _instanceIDSeparator,
	}
}

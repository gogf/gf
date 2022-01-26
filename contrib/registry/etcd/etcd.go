package etcd

import (
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	etcd3 "go.etcd.io/etcd/client/v3"
)

var (
	_ gsvc.Registry = &Registry{}
	_ gsvc.Watcher  = &watcher{}
)

// Registry is etcd registry.
type Registry struct {
	client       *etcd3.Client
	kv           etcd3.KV
	lease        etcd3.Lease
	keepaliveTTL time.Duration
	logger       *glog.Logger
}

type Option struct {
	Logger       *glog.Logger
	KeepaliveTTL time.Duration
}

const (
	DefaultKeepAliveTTL = 10 * time.Second
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

// extractResponseToServices extracts etcd watch response context to service list.
func extractResponseToServices(res *etcd3.GetResponse) ([]*gsvc.Service, error) {
	if res == nil || res.Kvs == nil {
		return nil, nil
	}
	var (
		services   []*gsvc.Service
		serviceKey string
		serviceMap = make(map[string]*gsvc.Service)
	)
	for _, kv := range res.Kvs {
		service, err := gsvc.NewServiceWithKV(kv.Key, kv.Value)
		if err != nil {
			return services, err
		}
		if service != nil {
			serviceKey = service.KeyWithoutEndpoints()
			if s, ok := serviceMap[serviceKey]; ok {
				s.Endpoints = append(s.Endpoints, service.Endpoints...)
			} else {
				serviceMap[serviceKey] = service
				services = append(services, service)
			}
		}
	}
	return services, nil
}

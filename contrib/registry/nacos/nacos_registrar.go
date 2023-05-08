package nacos

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

// Register registers `service` to Registry.
// Note that it returns a new Service if it changes the input Service with custom one.
func (r Registry) Register(ctx context.Context, service gsvc.Service) (registered gsvc.Service, err error) {
	_ = service.GetKey()
	var version string
	if service.GetVersion() == "" {
		version = gsvc.DefaultVersion
	} else {
		version = service.GetVersion()
	}
	name := gstr.Join(gstr.Split(service.GetName(), "/"), "")
	s := &gsvc.LocalService{
		Name:       name,
		Version:    version,
		Head:       r.opts.clusterName,
		Deployment: r.opts.groupName,
		Namespace:  r.opts.namespaceId,
		Metadata:   service.GetMetadata(),
		Endpoints:  service.GetEndpoints(),
	}
	if err = r.registerByType(s.GetPrefix(), service); err != nil {
		return nil, err
	}
	return s, nil
}

// Deregister off-lines and removes `service` from the Registry.
func (r Registry) Deregister(ctx context.Context, service gsvc.Service) error {
	client := r.namingClient
	for i := range service.GetEndpoints() {
		ok, err := client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          service.GetEndpoints()[i].Host(),
			Port:        uint64(service.GetEndpoints()[i].Port()),
			ServiceName: service.GetPrefix(),
			Cluster:     r.opts.clusterName,
			GroupName:   r.opts.groupName,
			Ephemeral:   true,
		})
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("deregister instance failed")
		}
	}
	return nil
}

func (r Registry) registerByType(name string, service gsvc.Service) error {
	client := r.namingClient
	endpoints := service.GetEndpoints()
	for i := range endpoints {
		metadata := make(map[string]string)
		for k, v := range service.GetMetadata() {
			metadata[k] = fmt.Sprintf("%v", v)
		}
		_, err := client.RegisterInstance(vo.RegisterInstanceParam{
			Ip:          endpoints[i].Host(),
			Port:        uint64(endpoints[i].Port()),
			ServiceName: name,
			Weight:      r.opts.weight,
			Enable:      true,
			Healthy:     true,
			Ephemeral:   true,
			Metadata:    metadata,
			ClusterName: r.opts.clusterName,
			GroupName:   r.opts.groupName,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"bytes"
	"context"
	"fmt"

	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// Watcher is a service watcher.
type Watcher struct {
	ServiceName      string
	Namespace        string
	Ctx              context.Context
	Cancel           context.CancelFunc
	Channel          <-chan model.SubScribeEvent
	ServiceInstances []gsvc.Service
}

func newWatcher(ctx context.Context, namespace string, key string, consumer polaris.ConsumerAPI) (*Watcher, error) {
	watchServiceResponse, err := consumer.WatchService(&polaris.WatchServiceRequest{
		WatchServiceRequest: model.WatchServiceRequest{
			Key: model.ServiceKey{
				Namespace: namespace,
				Service:   key,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		Namespace:        namespace,
		ServiceName:      key,
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
func (w *Watcher) Proceed() ([]gsvc.Service, error) {
	select {
	case <-w.Ctx.Done():
		return nil, w.Ctx.Err()
	case event := <-w.Channel:
		if event.GetSubScribeEventType() == model.EventInstance {
			// these are always true, but we need to check it to make sure EventType not change
			instanceEvent, ok := event.(*model.InstanceEvent)
			if !ok {
				return w.ServiceInstances, nil
			}
			// handle DeleteEvent
			if instanceEvent.DeleteEvent != nil {
				var endpointStr bytes.Buffer
				for _, instance := range instanceEvent.DeleteEvent.Instances {
					// Iterate through existing service instances, deleting them if they exist
					for i, serviceInstance := range w.ServiceInstances {
						if serviceInstance.(*Service).ID == instance.GetId() {
							endpointStr.WriteString(fmt.Sprintf("%s:%d%s", instance.GetHost(), instance.GetPort(), gsvc.EndpointsDelimiter))
							if len(w.ServiceInstances) <= 1 {
								w.ServiceInstances = w.ServiceInstances[0:0]
								continue
							}
							w.ServiceInstances = append(w.ServiceInstances[:i], w.ServiceInstances[i+1:]...)
						}
					}
				}
				if endpointStr.Len() > 0 && len(w.ServiceInstances) > 0 {
					var (
						newEndpointStr     bytes.Buffer
						serviceEndpointStr = w.ServiceInstances[0].(*Service).GetEndpoints().String()
					)
					for _, address := range gstr.SplitAndTrim(serviceEndpointStr, gsvc.EndpointsDelimiter) {
						if !gstr.Contains(endpointStr.String(), address) {
							newEndpointStr.WriteString(fmt.Sprintf("%s%s", address, gsvc.EndpointsDelimiter))
						}
					}

					for i := 0; i < len(w.ServiceInstances); i++ {
						w.ServiceInstances[i] = instanceToServiceInstance(instanceEvent.DeleteEvent.Instances[0], gstr.TrimRight(newEndpointStr.String(), gsvc.EndpointsDelimiter), w.ServiceInstances[i].(*Service).ID)
					}
				}
			}
			// handle UpdateEvent
			if instanceEvent.UpdateEvent != nil {
				var (
					updateEndpointStr bytes.Buffer
					newEndpointStr    bytes.Buffer
				)
				for _, serviceInstance := range w.ServiceInstances {
					// update the current department or all instances
					for _, update := range instanceEvent.UpdateEvent.UpdateList {
						if serviceInstance.(*Service).ID == update.Before.GetId() {
							// update equal
							if update.After.IsHealthy() {
								newEndpointStr.WriteString(fmt.Sprintf("%s:%d%s", update.After.GetHost(), update.After.GetPort(), gsvc.EndpointsDelimiter))
							}
							updateEndpointStr.WriteString(fmt.Sprintf("%s:%d%s", update.Before.GetHost(), update.Before.GetPort(), gsvc.EndpointsDelimiter))
						}
					}
				}
				if len(w.ServiceInstances) > 0 {
					var serviceEndpointStr = w.ServiceInstances[0].(*Service).GetEndpoints().String()
					// old instance addresses are culled
					if updateEndpointStr.Len() > 0 {
						for _, address := range gstr.SplitAndTrim(serviceEndpointStr, gsvc.EndpointsDelimiter) {
							// If the historical instance is not in the change instance, it remains
							if !gstr.Contains(updateEndpointStr.String(), address) {
								newEndpointStr.WriteString(fmt.Sprintf("%s%s", address, gsvc.EndpointsDelimiter))
							}
						}
					}
					instance := instanceEvent.UpdateEvent.UpdateList[0].After
					for i := 0; i < len(w.ServiceInstances); i++ {
						w.ServiceInstances[i] = instanceToServiceInstance(instance, gstr.TrimRight(newEndpointStr.String(), gsvc.EndpointsDelimiter), w.ServiceInstances[i].(*Service).ID)
					}
				}
			}
			// handle AddEvent
			if instanceEvent.AddEvent != nil {
				var (
					newEndpointStr bytes.Buffer
					allEndpointStr string
				)
				if len(w.ServiceInstances) > 0 {
					allEndpointStr = w.ServiceInstances[0].(*Service).GetEndpoints().String()
				}
				for i := 0; i < len(instanceEvent.AddEvent.Instances); i++ {
					instance := instanceEvent.AddEvent.Instances[i]
					if instance.IsHealthy() {
						address := fmt.Sprintf("%s:%d", instance.GetHost(), instance.GetPort())
						if !gstr.Contains(allEndpointStr, address) {
							newEndpointStr.WriteString(fmt.Sprintf("%s%s", address, gsvc.EndpointsDelimiter))
						}
					}
				}
				if newEndpointStr.Len() > 0 {
					allEndpointStr = fmt.Sprintf("%s%s", newEndpointStr.String(), allEndpointStr)
				}
				for i := 0; i < len(w.ServiceInstances); i++ {
					w.ServiceInstances[i] = instanceToServiceInstance(instanceEvent.AddEvent.Instances[0], gstr.TrimRight(allEndpointStr, gsvc.EndpointsDelimiter), w.ServiceInstances[i].(*Service).ID)
				}

				for i := 0; i < len(instanceEvent.AddEvent.Instances); i++ {
					instance := instanceEvent.AddEvent.Instances[i]
					if instance.IsHealthy() {
						w.ServiceInstances = append(w.ServiceInstances, instanceToServiceInstance(instance, gstr.TrimRight(allEndpointStr, gsvc.EndpointsDelimiter), ""))
					}
				}
			}
		}
	}

	return w.ServiceInstances, nil
}

// Close the watcher.
func (w *Watcher) Close() error {
	w.Cancel()
	return nil
}

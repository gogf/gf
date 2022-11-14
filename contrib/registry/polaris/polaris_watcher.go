// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/pkg/model"
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

func newWatcher(ctx context.Context, namespace string, serviceName string, consumer polaris.ConsumerAPI) (*Watcher, error) {
	watchServiceResponse, err := consumer.WatchService(&polaris.WatchServiceRequest{
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
				for _, instance := range instanceEvent.DeleteEvent.Instances {
					for i, serviceInstance := range w.ServiceInstances {
						if serviceInstance.(*Service).ID == instance.GetId() {
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
						if serviceInstance.(*Service).ID == update.Before.GetId() {
							w.ServiceInstances[i] = instanceToServiceInstance(update.After)
						}
					}
				}
			}
			// handle AddEvent
			if instanceEvent.AddEvent != nil {
				w.ServiceInstances = append(
					w.ServiceInstances,
					instancesToServiceInstances(instanceEvent.AddEvent.Instances)...,
				)
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

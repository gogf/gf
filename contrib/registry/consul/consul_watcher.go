// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consul

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

// Watcher implements the gsvc.Watcher interface for consul.
type Watcher struct {
	client    *api.Client
	registry  *Registry
	key       string
	plan      *watch.Plan
	closeChan chan struct{}
	closeOnce sync.Once
	services  map[string]*api.ServiceEntry
	eventChan chan []gsvc.Service
	mu        sync.RWMutex
}

// newWatcher creates and returns a new Watcher.
func newWatcher(client *api.Client, registry *Registry, key string) *Watcher {
	w := &Watcher{
		client:    client,
		registry:  registry,
		key:       key,
		closeChan: make(chan struct{}),
		services:  make(map[string]*api.ServiceEntry),
		eventChan: make(chan []gsvc.Service, 1),
	}

	go w.watch()
	return w
}

// watch starts the watching process.
func (w *Watcher) watch() {
	// Create watch plan
	plan, err := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": w.key,
	})
	if err != nil {
		select {
		case w.eventChan <- nil:
		default:
		}
		return
	}
	w.plan = plan

	// Set handler
	plan.Handler = func(idx uint64, data interface{}) {
		if data == nil {
			return
		}
		entries, ok := data.([]*api.ServiceEntry)
		if !ok {
			return
		}

		w.mu.Lock()
		// Clear old services
		w.services = make(map[string]*api.ServiceEntry)
		// Add new services
		for _, entry := range entries {
			if entry.Checks.AggregatedStatus() != api.HealthPassing {
				continue
			}
			w.services[entry.Service.ID] = entry
		}
		w.mu.Unlock()

		services, _ := w.Services()
		select {
		case <-w.closeChan:
			return
		case w.eventChan <- services:
		}
	}

	// Run the plan
	go func() {
		// Initial service query
		services, _, err := w.client.Health().Service(w.key, "", true, &api.QueryOptions{
			WaitTime: time.Second * 3,
		})
		if err == nil && len(services) > 0 {
			w.mu.Lock()
			for _, entry := range services {
				if entry.Checks.AggregatedStatus() != api.HealthPassing {
					continue
				}
				w.services[entry.Service.ID] = entry
			}
			w.mu.Unlock()

			if initialServices, err := w.Services(); err == nil {
				select {
				case <-w.closeChan:
					return
				case w.eventChan <- initialServices:
				}
			}
		}

		// Start watching
		if err := plan.Run(w.registry.GetAddress()); err != nil {
			select {
			case w.eventChan <- nil:
			default:
			}
		}
	}()

	// Wait for close signal
	<-w.closeChan
	if w.plan != nil {
		w.plan.Stop()
	}
}

// Proceed returns a Service event. It blocks until the watcher receives a service change.
func (w *Watcher) Proceed() ([]gsvc.Service, error) {
	select {
	case <-w.closeChan:
		return nil, gerror.New("watcher closed")
	case services := <-w.eventChan:
		if services == nil {
			return nil, gerror.New("watch failed")
		}
		return services, nil
	}
}

// Close closes the watcher.
func (w *Watcher) Close() error {
	w.closeOnce.Do(func() {
		close(w.closeChan)
	})
	return nil
}

// Services returns the latest service list from watcher.
func (w *Watcher) Services() ([]gsvc.Service, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var result []gsvc.Service
	for _, entry := range w.services {
		// Parse metadata
		var metadata map[string]interface{}
		if metaStr, ok := entry.Service.Meta["metadata"]; ok {
			if err := json.Unmarshal([]byte(metaStr), &metadata); err != nil {
				return nil, gerror.Wrap(err, "failed to unmarshal service metadata")
			}
		}

		// Create service instance
		localService := &gsvc.LocalService{
			Head:       "",
			Deployment: "",
			Namespace:  "",
			Name:       entry.Service.Service,
			Version:    entry.Service.Tags[0],
			Endpoints: []gsvc.Endpoint{
				gsvc.NewEndpoint(fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)),
			},
			Metadata: metadata,
		}
		result = append(result, localService)
	}

	return result, nil
}

// Done returns the done channel for the watcher.
func (w *Watcher) Done() <-chan struct{} {
	return w.closeChan
}

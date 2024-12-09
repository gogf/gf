// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consul

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
)

// Watcher watches the service changes.
type Watcher struct {
	registry  *Registry      // The registry instance
	key       string         // The service name to watch
	closeChan chan struct{}  // Channel for closing
	eventChan chan struct{}  // Channel for notifying changes
	mu        sync.RWMutex   // Mutex for thread safety
	plan      *watch.Plan    // The watch plan
	services  []gsvc.Service // Current services
}

// New creates and returns a new watcher.
func newWatcher(registry *Registry, key string) (*Watcher, error) {
	w := &Watcher{
		registry:  registry,
		key:       key,
		closeChan: make(chan struct{}),
		eventChan: make(chan struct{}, 1),
	}

	// Start watching
	go w.watch()

	return w, nil
}

// watch starts the watching process.
func (w *Watcher) watch() {
	// Get initial service list
	initServices, err := w.Services()
	if err != nil {
		return
	}

	// Set initial services
	w.mu.Lock()
	w.services = initServices
	w.mu.Unlock()

	// Create watch plan
	plan, err := watch.Parse(map[string]interface{}{
		"type":    "service",
		"service": w.key,
	})
	if err != nil {
		return
	}

	w.mu.Lock()
	w.plan = plan
	w.mu.Unlock()

	// Set handler
	plan.Handler = func(idx uint64, data interface{}) {
		// Check if watcher is closed
		select {
		case <-w.closeChan:
			return
		default:
		}

		// Get current services
		services, _ := w.Services()

		// Update services
		w.mu.Lock()
		w.services = services
		w.mu.Unlock()

		// Notify changes
		select {
		case w.eventChan <- struct{}{}:
		default:
		}
	}

	// Start watching
	go func() {
		defer func() {
			w.mu.Lock()
			if w.plan != nil {
				w.plan.Stop()
				w.plan = nil
			}
			w.mu.Unlock()
		}()

		if err = plan.Run(w.registry.GetAddress()); err != nil {
			return
		}
	}()

	// Wait for close signal
	<-w.closeChan
}

// Proceed returns current services and waits for the next service change.
func (w *Watcher) Proceed() ([]gsvc.Service, error) {
	// Check if watcher is closed
	select {
	case <-w.closeChan:
		return nil, gerror.New("watcher closed")
	default:
	}

	w.mu.RLock()
	services := w.services
	w.mu.RUnlock()

	// Wait for changes
	select {
	case <-w.closeChan:
		return nil, gerror.New("watcher closed")
	case <-w.eventChan:
		w.mu.RLock()
		services = w.services
		w.mu.RUnlock()
		return services, nil
	}
}

// Close closes the watcher.
func (w *Watcher) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	select {
	case <-w.closeChan:
		return nil
	default:
		close(w.closeChan)
		if w.plan != nil {
			w.plan.Stop()
			w.plan = nil
		}
		return nil
	}
}

// Services returns current services from the watcher.
func (w *Watcher) Services() ([]gsvc.Service, error) {
	// Query services directly from Consul
	entries, _, err := w.registry.client.Health().Service(w.key, "", true, &api.QueryOptions{})
	if err != nil {
		return nil, err
	}
	// Convert entries to services
	var services []gsvc.Service
	for _, entry := range entries {
		if entry.Checks.AggregatedStatus() == api.HealthPassing {
			metadata := make(map[string]interface{})
			if entry.Service.Meta != nil {
				if metaStr, ok := entry.Service.Meta["metadata"]; ok {
					if err := json.Unmarshal([]byte(metaStr), &metadata); err != nil {
						return nil, gerror.Wrap(err, "failed to unmarshal metadata")
					}
				}
			}

			// Get version from metadata or tags
			version := ""
			if v, ok := entry.Service.Meta["version"]; ok {
				version = v
			} else if len(entry.Service.Tags) > 0 {
				version = entry.Service.Tags[0]
			}

			// Create service instance
			service := &gsvc.LocalService{
				Name:     entry.Service.Service,
				Version:  version,
				Metadata: metadata,
				Endpoints: []gsvc.Endpoint{
					gsvc.NewEndpoint(fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)),
				},
			}
			services = append(services, service)
		}
	}

	// Sort services by version
	if len(services) > 0 {
		sort.Slice(services, func(i, j int) bool {
			return services[i].GetVersion() < services[j].GetVersion()
		})
	}
	return services, nil
}

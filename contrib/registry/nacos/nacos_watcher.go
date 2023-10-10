// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/joy999/nacos-sdk-go/model"
)

// Watcher used to mange service event such as update.
type Watcher struct {
	ctx   context.Context
	event chan *watchEvent
	close func() error
}

// watchEvent
type watchEvent struct {
	Services []model.Instance
	Err      error
}

// newWatcher new a Watcher's instance
func newWatcher(ctx context.Context) *Watcher {
	w := &Watcher{
		ctx:   ctx,
		event: make(chan *watchEvent, 10),
	}
	return w
}

// Proceed proceeds watch in blocking way.
// It returns all complete services that watched by `key` if any change.
func (w *Watcher) Proceed() (services []gsvc.Service, err error) {
	n := len(w.event)
	servicesMap := map[string]gsvc.Service{}
	for i := 0; i < n; i++ {
		e := <-w.event
		if e.Err != nil {
			err = e.Err
			return
		}
		newServices := NewServicesFromInstances(e.Services)
		for _, s := range newServices {
			servicesMap[s.GetName()] = s
		}
	}
	services = make([]gsvc.Service, 0, len(servicesMap))
	for _, s := range servicesMap {
		services = append(services, s)
	}
	return
}

// Close closes the watcher.
func (w *Watcher) Close() (err error) {
	if w.close != nil {
		err = w.close()
	}
	return
}

// SetCloseFunc set the close callback function
func (w *Watcher) SetCloseFunc(close func() error) {
	w.close = close
}

// Push add the services watchevent to event queue
func (w *Watcher) Push(services []model.Instance, err error) {
	w.event <- &watchEvent{
		Services: services,
		Err:      err,
	}
}

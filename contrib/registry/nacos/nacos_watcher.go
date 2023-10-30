// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
	e, ok := <-w.event
	if !ok || e == nil {
		err = gerror.NewCode(gcode.CodeNil)
		return
	}
	if e.Err != nil {
		err = e.Err
		return
	}
	services = NewServicesFromInstances(e.Services)
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

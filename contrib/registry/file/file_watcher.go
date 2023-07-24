// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package file

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
)

// Watcher for file changes watch.
type Watcher struct {
	prefix    string            // Watched prefix key, not file name prefix.
	discovery gsvc.Discovery    // Service discovery.
	ch        chan gsvc.Service // Changes that caused by inotify.
}

// Proceed proceeds watch in blocking way.
// It returns all complete services that watched by `key` if any change.
func (w *Watcher) Proceed() (services []gsvc.Service, err error) {
	<-w.ch
	return w.discovery.Search(context.Background(), gsvc.SearchInput{
		Prefix: w.prefix,
	})
}

// Close closes the watcher.
func (w *Watcher) Close() error {
	close(w.ch)
	return nil
}

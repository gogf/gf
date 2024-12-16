// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package etcd

import (
	"context"
	"time"

	etcd3 "go.etcd.io/etcd/client/v3"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsvc"
)

var (
	_ gsvc.Watcher = &watcher{}
)

type watcher struct {
	key       string
	ctx       context.Context
	cancel    context.CancelFunc
	watchChan etcd3.WatchChan
	watcher   etcd3.Watcher
	kv        etcd3.KV
}

func newWatcher(key string, client *etcd3.Client, dialTimeout time.Duration) (*watcher, error) {
	w := &watcher{
		key:     key,
		watcher: etcd3.NewWatcher(client),
		kv:      etcd3.NewKV(client),
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), dialTimeout)
	defer cancel()

	// Test connection first.
	if _, err := client.Get(ctx, "ping"); err != nil {
		return nil, gerror.WrapCode(gcode.CodeOperationFailed, err, "failed to connect to etcd")
	}

	w.ctx, w.cancel = context.WithCancel(context.Background())
	w.watchChan = w.watcher.Watch(w.ctx, key, etcd3.WithPrefix(), etcd3.WithRev(0))

	if err := w.watcher.RequestProgress(context.Background()); err != nil {
		// Clean up
		w.cancel()
		return nil, gerror.WrapCode(gcode.CodeOperationFailed, err, "failed to establish watch connection")
	}

	return w, nil
}

// Proceed is used to watch the key.
func (w *watcher) Proceed() ([]gsvc.Service, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		// It retrieves, merges and returns all services by prefix if any changes.
		return w.getServicesByPrefix()
	}
}

// Close is used to close the watcher.
func (w *watcher) Close() error {
	w.cancel()
	return w.watcher.Close()
}

func (w *watcher) getServicesByPrefix() ([]gsvc.Service, error) {
	res, err := w.kv.Get(w.ctx, w.key, etcd3.WithPrefix())
	if err != nil {
		return nil, err
	}
	return extractResponseToServices(res)
}

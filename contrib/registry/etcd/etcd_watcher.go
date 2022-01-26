package etcd

import (
	"context"

	"github.com/gogf/gf/v2/net/gsvc"
	etcd3 "go.etcd.io/etcd/client/v3"
)

type watcher struct {
	key         string
	ctx         context.Context
	cancel      context.CancelFunc
	watchChan   etcd3.WatchChan
	watcher     etcd3.Watcher
	kv          etcd3.KV
	initialized bool
}

func newWatcher(ctx context.Context, key string, client *etcd3.Client) (*watcher, error) {
	w := &watcher{
		key:         key,
		watcher:     etcd3.NewWatcher(client),
		kv:          etcd3.NewKV(client),
		initialized: false,
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
	w.watchChan = w.watcher.Watch(w.ctx, key, etcd3.WithPrefix(), etcd3.WithRev(0))
	err := w.watcher.RequestProgress(context.Background())
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *watcher) Proceed() ([]*gsvc.Service, error) {
	if !w.initialized {
		w.initialized = true
		return w.getServicesByPrefix()
	}
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		return w.getServicesByPrefix()
	}
}

func (w *watcher) Close() error {
	w.cancel()
	return w.watcher.Close()
}

func (w *watcher) getServicesByPrefix() ([]*gsvc.Service, error) {
	res, err := w.kv.Get(w.ctx, w.key, etcd3.WithPrefix())
	if err != nil {
		return nil, err
	}
	return extractResponseToServices(res)
}

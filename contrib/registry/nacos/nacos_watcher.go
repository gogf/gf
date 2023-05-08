package nacos

import (
	"context"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

var _ gsvc.Watcher = &watcher{}

type watcher struct {
	key       string
	ctx       context.Context
	cancel    context.CancelFunc
	watchChan chan struct{}
	opts      *options
	client    naming_client.INamingClient
}

func newWatcher(key string, client naming_client.INamingClient, ops *options) (*watcher, error) {
	w := &watcher{
		key:       key,
		watchChan: make(chan struct{}),
		opts:      ops,
		client:    client,
	}
	w.ctx, w.cancel = context.WithCancel(context.Background())
	go w.watch(key)
	return w, nil
}

// Close is used to close the watcher.
func (w watcher) watch(key string) {
	for {
		//subscribe the key
		err := w.client.Subscribe(&vo.SubscribeParam{
			ServiceName: key,
			GroupName:   w.opts.groupName,
			Clusters:    []string{w.opts.clusterName},
			SubscribeCallback: func(services []model.Instance, err error) {
				//if service has changed,will send a channel to w.watchChan
				w.watchChan <- struct{}{}
			},
		})
		if err != nil {
			w.ctx.Done()
		}
	}

}

// Proceed is used to watch the key.
func (w watcher) Proceed() (services []gsvc.Service, err error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case <-w.watchChan:
		// It retrieves, merges and returns all services by prefix if any changes.
		instances, err := getServiceFromInstances(w.key, w.opts, w.client)
		if err != nil {
			return nil, err
		}
		return instances, nil
	}
}

func (w watcher) Close() error {
	w.cancel()
	return nil
}

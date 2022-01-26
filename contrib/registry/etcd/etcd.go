package etcd

import (
	"time"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	etcd3 "go.etcd.io/etcd/client/v3"
)

var (
	_ gsvc.Registry = &Registry{}
	_ gsvc.Watcher  = &watcher{}
)

// Registry is etcd registry.
type Registry struct {
	client       *etcd3.Client
	kv           etcd3.KV
	lease        etcd3.Lease
	keepaliveTTL time.Duration
	logger       *glog.Logger
}

type Option struct {
	Logger       *glog.Logger
	KeepaliveTTL time.Duration
}

const (
	DefaultKeepAliveTTL = 10 * time.Second
)

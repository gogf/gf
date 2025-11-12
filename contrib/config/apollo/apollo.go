// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package apollo implements gcfg.Adapter using apollo service.
package apollo

import (
	"context"
	"strings"

	"github.com/apolloconfig/agollo/v4"
	apolloConfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
)

var (
	// Compile-time checking for interface implementation.
	_ gcfg.Adapter        = (*Client)(nil)
	_ gcfg.WatcherAdapter = (*Client)(nil)
)

const (
	apolloNamespaceDelimiter = ","
)

// Config is the configuration object for apollo client.
type Config struct {
	AppID             string `v:"required"` // See apolloConfig.Config.
	IP                string `v:"required"` // See apolloConfig.Config.
	Cluster           string `v:"required"` // See apolloConfig.Config.
	NamespaceName     string // See apolloConfig.Config.
	IsBackupConfig    bool   // See apolloConfig.Config.
	BackupConfigPath  string // See apolloConfig.Config.
	Secret            string // See apolloConfig.Config.
	SyncServerTimeout int    // See apolloConfig.Config.
	MustStart         bool   // See apolloConfig.Config.
	Watch             bool   // Watch watches remote configuration updates, which updates local configuration in memory immediately when remote configuration changes.
}

// Client implements gcfg.Adapter implementing using apollo service.
type Client struct {
	config   Config                // Config object when created.
	client   agollo.Client         // Apollo client.
	value    *g.Var                // Configmap content cached. It is `*gjson.Json` value internally.
	watchers *gcfg.WatcherRegistry // Watchers for watching file changes.
}

// New creates and returns gcfg.Adapter implementing using apollo service.
func New(ctx context.Context, config Config) (adapter gcfg.Adapter, err error) {
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}
	if config.NamespaceName == "" {
		config.NamespaceName = storage.GetDefaultNamespace()
	}
	client := &Client{
		config:   config,
		value:    g.NewVar(nil, true),
		watchers: gcfg.NewWatcherRegistry(),
	}
	// Apollo client.
	client.client, err = agollo.StartWithConfig(func() (*apolloConfig.AppConfig, error) {
		return &apolloConfig.AppConfig{
			AppID:             config.AppID,
			Cluster:           config.Cluster,
			NamespaceName:     config.NamespaceName,
			IP:                config.IP,
			IsBackupConfig:    config.IsBackupConfig,
			BackupConfigPath:  config.BackupConfigPath,
			Secret:            config.Secret,
			SyncServerTimeout: config.SyncServerTimeout,
			MustStart:         config.MustStart,
		}, nil
	})
	if err != nil {
		return nil, gerror.Wrapf(err, `create apollo client failed with config: %+v`, config)
	}
	if config.Watch {
		client.client.AddChangeListener(client)
	}
	return client, nil
}

// Available checks and returns the backend configuration service is available.
// The optional parameter `resource` specifies certain configuration resource.
//
// Note that this function does not return error as it just does simply check for
// backend configuration service.
func (c *Client) Available(ctx context.Context, resource ...string) (ok bool) {
	if len(resource) == 0 && !c.value.IsNil() {
		return true
	}

	namespaces := strings.Split(c.config.NamespaceName, apolloNamespaceDelimiter)
	if len(resource) > 0 {
		namespaces = resource
	}

	for _, namespace := range namespaces {
		if c.client.GetConfig(namespace) == nil {
			return false
		}
	}

	return true
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (c *Client) Get(ctx context.Context, pattern string) (value any, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValue(ctx); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Get(pattern).Val(), nil
}

// Data retrieves and returns all configuration data in current resource as map.
// Note that this function may lead lots of memory usage if configuration data is too large,
// you can implement this function if necessary.
func (c *Client) Data(ctx context.Context) (data map[string]any, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValue(ctx); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Map(), nil
}

// OnChange is called when config changes.
func (c *Client) OnChange(event *storage.ChangeEvent) {
	_ = c.updateLocalValue(gctx.New())
}

// OnNewestChange is called when any config changes.
func (c *Client) OnNewestChange(event *storage.FullChangeEvent) {
	// Nothing to do.
}

func (c *Client) updateLocalValue(ctx context.Context) (err error) {
	j := gjson.New(nil)
	content := gjson.New(nil, true)

	for _, namespace := range strings.Split(c.config.NamespaceName, apolloNamespaceDelimiter) {
		cache := c.client.GetConfigCache(namespace)
		cache.Range(func(key, value any) bool {
			err = j.Set(gconv.String(key), value)
			if err != nil {
				return false
			}
			err = content.Set(gconv.String(key), value)
			if err != nil {
				return false
			}
			return true
		})

		if err != nil {
			break
		}
	}

	if err == nil {
		c.value.Set(j)
		adapterCtx := NewAdapterCtx(ctx).WithOperation(gcfg.OperationUpdate).WithNamespace(c.config.NamespaceName).
			WithAppId(c.config.AppID).WithCluster(c.config.Cluster).WithContent(content)
		c.notifyWatchers(adapterCtx.Ctx)
	}
	return
}

// AddWatcher adds a watcher for the specified configuration file.
func (c *Client) AddWatcher(name string, f func(ctx context.Context)) {
	c.watchers.Add(name, f)
}

// RemoveWatcher removes the watcher for the specified configuration file.
func (c *Client) RemoveWatcher(name string) {
	c.watchers.Remove(name)
}

// GetWatcherNames returns all watcher names.
func (c *Client) GetWatcherNames() []string {
	return c.watchers.GetNames()
}

// notifyWatchers notifies all watchers.
func (c *Client) notifyWatchers(ctx context.Context) {
	c.watchers.Notify(ctx)
}

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package consul implements gcfg.Adapter using consul service.
package consul

import (
	"context"
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/glog"
)

// Config is the configuration object for consul client.
type Config struct {
	// api.Config in consul package
	ConsulConfig api.Config `v:"required"`
	// As configuration file path key
	Path string `v:"required"`
	// Watch watches remote configuration updates, which updates local configuration in memory immediately when remote configuration changes.
	Watch bool
	// Logging interface, customized by user, default: glog.New()
	Logger glog.ILogger
}

// Client implements gcfg.Adapter implementing using consul service.
type Client struct {
	// Created config object
	config Config
	// Consul config client
	client *api.Client
	// Configmap content cached. It is `*gjson.Json` value internally.
	value *g.Var
}

// New creates and returns gcfg.Adapter implementing using consul service.
func New(ctx context.Context, config Config) (adapter gcfg.Adapter, err error) {
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}

	if config.Logger == nil {
		config.Logger = glog.New()
	}

	client := &Client{
		config: config,
		value:  g.NewVar(nil, true),
	}

	client.client, err = api.NewClient(&config.ConsulConfig)
	if err != nil {
		return nil, gerror.Wrapf(err, `create consul client failed with config: %+v`, config.ConsulConfig)
	}

	if err = client.addWatcher(); err != nil {
		return nil, gerror.Wrapf(err, `consul client add watcher failed with config: %+v`, config.ConsulConfig)
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

	_, _, err := c.client.KV().Get(c.config.Path, nil)

	return err == nil
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (c *Client) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValue(); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Get(pattern).Val(), nil
}

// Data retrieves and returns all configuration data in current resource as map.
// Note that this function may lead lots of memory usage if configuration data is too large,
// you can implement this function if necessary.
func (c *Client) Data(ctx context.Context) (data map[string]interface{}, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValue(); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Map(), nil
}

func (c *Client) updateLocalValue() (err error) {
	content, _, err := c.client.KV().Get(c.config.Path, nil)
	if err != nil {
		return gerror.Wrapf(err, `get config from consul path [%+v] failed`, c.config.Path)
	}
	if content == nil {
		return fmt.Errorf(`get config from consul path [%+v] value is nil`, c.config.Path)
	}
	return c.doUpdate(content.Value)
}

func (c *Client) doUpdate(content []byte) (err error) {
	var j *gjson.Json
	if j, err = gjson.LoadContent(content); err != nil {
		return gerror.Wrapf(err,
			`parse config map item from consul path [%+v] failed`, c.config.Path)
	}
	c.value.Set(j)
	return nil
}

func (c *Client) addWatcher() (err error) {
	if !c.config.Watch {
		return nil
	}

	plan, err := watch.Parse(map[string]interface{}{
		"type": "key",
		"key":  c.config.Path,
	})
	if err != nil {
		return gerror.Wrapf(err, `watch config from consul path %+v failed`, c.config.Path)
	}

	plan.Handler = func(idx uint64, raw interface{}) {
		var v *api.KVPair
		if raw == nil {
			// nil is a valid return value
			v = nil
			return
		}
		var ok bool
		if v, ok = raw.(*api.KVPair); !ok {
			return
		}

		if err = c.doUpdate(v.Value); err != nil {
			c.config.Logger.Errorf(
				context.Background(),
				"watch config from consul path %+v update failed: %s",
				c.config.Path, err,
			)
		}
	}

	plan.Datacenter = c.config.ConsulConfig.Datacenter
	plan.Token = c.config.ConsulConfig.Token

	go c.startAsynchronousWatch(plan)
	return nil
}

func (c *Client) startAsynchronousWatch(plan *watch.Plan) {
	if err := plan.Run(c.config.ConsulConfig.Address); err != nil {
		c.config.Logger.Errorf(
			context.Background(),
			"watch config from consul path %+v plan start failed: %s",
			c.config.Path, err,
		)
	}
}

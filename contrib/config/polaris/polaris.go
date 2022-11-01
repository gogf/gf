// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package polaris implements gcfg.Adapter using polaris service.
package polaris

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/polarismesh/polaris-go"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/model"
)

// LogDir sets the log directory for polaris.
func LogDir(dir string) error {
	return api.SetLoggersDir(dir)
}

// Config is the configuration for polaris.
type Config struct {
	// The namespace of the configuration.
	Namespace string `v:"required"`
	// The group of the configuration.
	FileGroup string `v:"required"`
	// The name of the configuration.
	FileName string `v:"required"`
	// The path of the polaris configuration file.
	Path string `v:"required"`
	// The log directory for polaris.
	LogDir string
	// Watch watches remote configuration updates, which updates local configuration in memory immediately when remote configuration changes.
	Watch bool
}

// Client implements gcfg.Adapter implementing using polaris service.
type Client struct {
	config Config
	client model.ConfigFile
	value  *g.Var
}

const defaultLogDir = "/tmp/polaris/log"

// New creates and returns gcfg.Adapter implementing using polaris service.
func New(ctx context.Context, config Config) (adapter gcfg.Adapter, err error) {
	if err = g.Validator().Data(config).Run(ctx); err != nil {
		err = gerror.Wrap(err, "invalid polaris config")
		return nil, err
	}
	var (
		client = &Client{
			config: config,
			value:  g.NewVar(nil, true),
		}
		configAPI polaris.ConfigAPI
	)

	if configAPI, err = polaris.NewConfigAPIByFile(config.Path); err != nil {
		err = gerror.Wrapf(err, "Polaris configuration initialization failed  with config: %+v", config)
		return
	}
	// set log dir
	if gstr.Trim(config.LogDir) == "" {
		config.LogDir = defaultLogDir
	}
	if err = LogDir(config.LogDir); err != nil {
		err = gerror.Wrap(err, "set polaris log dir failed")
		return
	}

	if client.client, err = configAPI.GetConfigFile(config.Namespace, config.FileGroup, config.FileName); err != nil {
		err = gerror.Wrapf(err, "failed to read data from Polaris configuration center  with config: %+v", config)
		return
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

	var namespace = c.config.Namespace
	if len(resource) > 0 {
		namespace = resource[0]
	}

	return c.client.GetNamespace() == namespace
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (c *Client) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValueAndWatch(ctx); err != nil {
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
		if err = c.updateLocalValueAndWatch(ctx); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Map(), nil
}

// init retrieves and caches the configmap content.
func (c *Client) updateLocalValueAndWatch(ctx context.Context) (err error) {
	if err = c.doUpdate(ctx); err != nil {
		err = gerror.Wrap(err, "failed to update local value")
		return err
	}
	if err = c.doWatch(ctx); err != nil {
		err = gerror.Wrap(err, "failed to watch configmap")
		return err
	}
	return nil
}

func (c *Client) doUpdate(ctx context.Context) (err error) {
	if !c.client.HasContent() {
		return gerror.New("config file is empty")
	}
	var j *gjson.Json
	if j, err = gjson.LoadContent(c.client.GetContent()); err != nil {
		return gerror.Wrap(err, `parse config map item from polaris failed`)
	}
	c.value.Set(j)
	return nil
}

func (c *Client) doWatch(ctx context.Context) (err error) {
	if !c.config.Watch {
		return nil
	}
	var changeChan = make(chan model.ConfigFileChangeEvent)
	c.client.AddChangeListenerWithChannel(changeChan)
	go func() {
		for {
			select {
			case <-changeChan:
				_ = c.doUpdate(ctx)
			}
		}
	}()
	return nil
}

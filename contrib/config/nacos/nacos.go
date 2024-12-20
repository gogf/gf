// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package nacos implements gcfg.Adapter using nacos service.
package nacos

import (
	"context"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

const (
	componentNameNacos = "gf.component.nacos"
)

var (
	instance = gmap.StrAnyMap{}
)

// Config is the configuration object for nacos client.
type Config struct {
	ServerConfigs  []constant.ServerConfig                     `v:"required"` // See constant.ServerConfig
	ClientConfig   constant.ClientConfig                       `v:"required"` // See constant.ClientConfig
	ConfigParam    vo.ConfigParam                              `v:"required"` // See vo.ConfigParam
	Watch          bool                                        // Watch watches remote configuration updates, which updates local configuration in memory immediately when remote configuration changes.
	OnConfigChange func(namespace, group, dataId, data string) // Configure change callback function
}

// Client implements gcfg.Adapter implementing using nacos service.
type Client struct {
	config Config                      // Config object when created.
	client config_client.IConfigClient // Nacos config client.
	value  *g.Var                      // Configmap content cached. It is `*gjson.Json` value internally.
}

// New creates and returns gcfg.Adapter implementing using nacos service.
func New(ctx context.Context, config Config) (adapter gcfg.Adapter, err error) {
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}

	client := &Client{
		config: config,
		value:  g.NewVar(nil, true),
	}
	client.initConfigClient()

	err = client.addWatcher()
	if err != nil {
		return nil, err
	}

	return client, nil
}

// initConfigClient create a singleton config_client.IConfigClient
func (c *Client) initConfigClient() {
	key := componentNameNacos
	for _, v := range c.config.ServerConfigs {
		key = fmt.Sprintf("%s:%s:%d", key, v.IpAddr, v.Port)
	}
	if c.config.ClientConfig.NamespaceId != "" {
		key += ":" + c.config.ClientConfig.NamespaceId
	}
	if c.config.ClientConfig.AccessKey != "" {
		key += ":" + c.config.ClientConfig.AccessKey
	}
	if c.config.ClientConfig.Username != "" {
		key += ":" + c.config.ClientConfig.Username
	}
	res := instance.GetOrSetFuncLock(key, func() interface{} {
		ins, err := clients.CreateConfigClient(map[string]interface{}{
			"serverConfigs": c.config.ServerConfigs,
			"clientConfig":  c.config.ClientConfig,
		})
		if err != nil {
			panic("create nacos client failed with config: " + err.Error())
		}
		return ins
	})
	c.client = res.(config_client.IConfigClient)
	return
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
	_, err := c.client.GetConfig(c.config.ConfigParam)
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
	content, err := c.client.GetConfig(c.config.ConfigParam)
	if err != nil {
		return gerror.Wrap(err, `retrieve config from nacos failed`)
	}

	return c.doUpdate(content)
}

func (c *Client) doUpdate(content string) (err error) {
	var j *gjson.Json
	if j, err = gjson.LoadContent([]byte(content)); err != nil {
		return gerror.Wrap(err, `parse config map item from nacos failed`)
	}
	c.value.Set(j)
	return nil
}

func (c *Client) addWatcher() error {
	if !c.config.Watch {
		return nil
	}
	c.config.ConfigParam.OnChange = func(namespace, group, dataId, data string) {
		c.doUpdate(data)
		if c.config.OnConfigChange != nil {
			go c.config.OnConfigChange(namespace, group, dataId, data)
		}
	}

	if err := c.client.ListenConfig(c.config.ConfigParam); err != nil {
		return gerror.Wrap(err, `watch config from namespace failed`)
	}

	return nil
}

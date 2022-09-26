// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package kubecm implements gcfg.Adapter using kubernetes configmap.
package kubecm

import (
	"context"

	kubeMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/util/gutil"
)

// Client implements gcfg.Adapter.
type Client struct {
	Config        // Config object when created.
	value  *g.Var // Configmap content cached. It is `*gjson.Json` value internally.
}

// Config for Client.
type Config struct {
	ConfigMap  string                `v:"required"` // ConfigMap name.
	DataItem   string                `v:"required"` // DataItem is the key item in Configmap data.
	Namespace  string                // (Optional) Specify the namespace for configmap.
	RestConfig *rest.Config          // (Optional) Custom rest config for kube client.
	KubeClient *kubernetes.Clientset // (Optional) Custom kube client.
}

// New creates and returns gcfg.Adapter implementing using kubernetes configmap.
func New(ctx context.Context, config Config) (adapter gcfg.Adapter, err error) {
	// Data validation.
	err = g.Validator().Data(config).Run(ctx)
	if err != nil {
		return nil, err
	}
	// Kubernetes client creating.
	if config.KubeClient == nil {
		if config.RestConfig == nil {
			config.RestConfig, err = NewDefaultKubeConfig(ctx)
			if err != nil {
				return nil, gerror.Wrapf(err, `create kube config failed`)
			}
		}
		config.KubeClient, err = kubernetes.NewForConfig(config.RestConfig)
		if err != nil {
			return nil, gerror.Wrapf(err, `create kube client failed`)
		}
	}
	adapter = &Client{
		Config: config,
		value:  g.NewVar(nil, true),
	}
	return
}

// Available checks and returns the backend configuration service is available.
// The optional parameter `resource` specifies certain configuration resource.
//
// Note that this function does not return error as it just does simply check for
// backend configuration service.
func (c *Client) Available(ctx context.Context, configMap ...string) (ok bool) {
	err := c.init(ctx, configMap...)
	return err == nil
}

// Get retrieves and returns value by specified `pattern` in current resource.
// Pattern like:
// "x.y.z" for map item.
// "x.0.y" for slice item.
func (c *Client) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if c.value.IsNil() {
		if err = c.init(ctx); err != nil {
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
		if err = c.init(ctx); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Map(), nil
}

// init retrieves and caches the configmap content.
func (c *Client) init(ctx context.Context, configMap ...string) (err error) {
	var (
		namespace     = gutil.GetOrDefaultStr(Namespace(), c.Namespace)
		configMapName = gutil.GetOrDefaultStr(c.ConfigMap, configMap...)
	)
	cm, err := c.KubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, kubeMetaV1.GetOptions{})
	if err != nil {
		return gerror.Wrapf(
			err,
			`retrieve configmap "%s" from namespace "%s" failed`,
			configMapName, namespace,
		)
	}
	if c.value.IsNil() {
		var j *gjson.Json
		if c.DataItem != "" {
			j, err = gjson.LoadContent(cm.Data[c.DataItem])
			if err != nil {
				return gerror.Wrapf(
					err,
					`parse config map item from %s[%s] failed`, configMapName, c.DataItem,
				)
			}
			c.value.Set(j)
		} else {
			j = gjson.New(cm.Data)
			c.value.Set(j)
		}
	}
	return nil
}

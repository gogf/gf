// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package kubecm implements gcfg.Adapter using kubernetes configmap.
package kubecm

import (
	"context"
	"fmt"

	kubeMetaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
	config Config                // Config object when created.
	client *kubernetes.Clientset // Kubernetes client.
	value  *g.Var                // Configmap content cached. It is `*gjson.Json` value internally.
}

// Config for Client.
type Config struct {
	ConfigMap  string                `v:"required"` // ConfigMap name.
	DataItem   string                `v:"required"` // DataItem is the key item in Configmap data.
	Namespace  string                // Specify the namespace for configmap.
	RestConfig *rest.Config          // Custom rest config for kube client.
	KubeClient *kubernetes.Clientset // Custom kube client.
	Watch      bool                  // Watch updates, which updates configuration when configmap changes.
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
		config: config,
		client: config.KubeClient,
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
	if len(configMap) == 0 && !c.value.IsNil() {
		return true
	}

	var (
		namespace     = gutil.GetOrDefaultStr(Namespace(), c.config.Namespace)
		configMapName = gutil.GetOrDefaultStr(c.config.ConfigMap, configMap...)
	)
	_, err := c.config.KubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, configMapName, kubeMetaV1.GetOptions{})
	if err != nil {
		return false
	}
	return true
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
	var namespace = gutil.GetOrDefaultStr(Namespace(), c.config.Namespace)
	err = c.doUpdate(ctx, namespace)
	if err != nil {
		return err
	}
	err = c.doWatch(ctx, namespace)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) doUpdate(ctx context.Context, namespace string) (err error) {
	cm, err := c.client.CoreV1().ConfigMaps(namespace).Get(ctx, c.config.ConfigMap, kubeMetaV1.GetOptions{})
	if err != nil {
		return gerror.Wrapf(
			err,
			`retrieve configmap "%s" from namespace "%s" failed`,
			c.config.ConfigMap, namespace,
		)
	}
	var j *gjson.Json
	if j, err = gjson.LoadContent(cm.Data[c.config.DataItem]); err != nil {
		return gerror.Wrapf(
			err,
			`parse config map item from %s[%s] failed`, c.config.ConfigMap, c.config.DataItem,
		)
	}
	c.value.Set(j)
	return nil
}

func (c *Client) doWatch(ctx context.Context, namespace string) (err error) {
	if !c.config.Watch {
		return nil
	}
	var watchHandler watch.Interface
	watchHandler, err = c.client.CoreV1().ConfigMaps(namespace).Watch(ctx, kubeMetaV1.ListOptions{
		FieldSelector: fmt.Sprintf(`metadata.name=%s`, c.config.ConfigMap),
		Watch:         true,
	})
	if err != nil {
		return gerror.Wrapf(
			err,
			`watch configmap "%s" from namespace "%s" failed`,
			c.config.ConfigMap, namespace,
		)
	}
	go func() {
		for {
			event := <-watchHandler.ResultChan()
			switch event.Type {
			case watch.Modified:
				_ = c.doUpdate(ctx, namespace)
			}
		}
	}()
	return nil
}

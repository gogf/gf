package etcd

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"time"
)

type Config struct {
	Endpoints []string `v:"required"` // See etcd Config.
	DialTimeout time.Duration
	Username string
	Password string
	TLS *tls.Config
	DialOptions []grpc.DialOption
	Watch          bool   // Watch watches remote configuration updates
	ConfigKey     string `v:"required"`
}

// Client implements gcfg.Adapter implementing using apollo service.
type Client struct {
	config Config        // Config object when created.
	client *clientv3.Client // Etcd client.
	value  *g.Var        // Configmap content cached. It is `*gjson.Json` value internally.
}

// New creates and returns gcfg.Adapter implementing using etcd service.
func New(ctx context.Context, config Config) (adapter gcfg.Adapter, err error) {
	err = g.Validator().Data(config).Run(ctx)
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	client := &Client{
		config: config,
		value:  g.NewVar(nil, true),
	}
	// etcd client.
	client.client, err = clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: config.DialTimeout,
		Username: config.Username,
		Password: config.Password,
		TLS: config.TLS,
		DialOptions: config.DialOptions,
		Context : ctx,
	})
	if err != nil {
		return nil, err
	}
	if config.Watch {
		go client.addWatcher(ctx)
	}

	return client, nil
}

func (c *Client) Available(ctx context.Context, resource ...string) (ok bool) {
	if len(resource) == 0 && !c.value.IsNil() {
		return true
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	_, err := c.client.Status(ctx, c.config.Endpoints[0])
	cancel()
	return err == nil
}

func (c *Client) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValue(ctx); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Get(pattern).Val(), nil
}

func (c *Client) Data(ctx context.Context) (data map[string]interface{}, err error) {
	if c.value.IsNil() {
		if err = c.updateLocalValue(ctx); err != nil {
			return nil, err
		}
	}
	return c.value.Val().(*gjson.Json).Map(), nil
}

func (c *Client) updateLocalValue(ctx context.Context) (err error) {
	if !c.config.Watch {
		defer c.client.Close()
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	resp, err := c.client.Get(ctx, c.config.ConfigKey)
	cancel()
	if err != nil {
		return gerror.Wrap(err, `retrieve config from etcd failed`)
	}
	if len(resp.Kvs) > 0 {
		content := string(resp.Kvs[0].Value)
		if err := c.doUpdate(content); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) doUpdate(content string) (err error) {
	var j *gjson.Json
	if j, err = gjson.LoadContent(content); err != nil {
		return gerror.Wrap(err, `parse config from etcd failed`)
	}
	c.value.Set(j)
	return nil
}

func (c *Client) addWatcher(ctx context.Context) {
	defer c.client.Close()
	watchChan := c.client.Watch(ctx, c.config.ConfigKey)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			_ = c.doUpdate(string(event.Kv.Value))
		}
	}
}
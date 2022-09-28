package apollo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/spf13/viper"
	"strings"
)

var (
	appId   = "SampleApp"
	cluster = "default"
	ip      = "http://localhost:8080"
)

func init() {
	ctx := gctx.New()
	adapterApollo, err := NewApollo(appId, cluster, ip).Run()
	if err != nil {
		g.Log().Errorf(ctx, "apollo init fail, err=%v", err)
		return
	}
	g.Cfg().SetAdapter(adapterApollo)

}

type ApolloConfig struct {
	*config.AppConfig
}

func NewApollo(appId, cluster, ip string) *ApolloConfig {
	return &ApolloConfig{&config.AppConfig{
		AppID:         appId,
		IP:            ip,
		Cluster:       cluster,
		NamespaceName: storage.GetDefaultNamespace(),
	}}
}

type ConfigCenter interface {
	Run() error
}

func (apolloConfig *ApolloConfig) Run() (*AdapterApollo, error) {
	client, _ := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return apolloConfig.AppConfig, nil
	})
	client.AddChangeListener(&ConfigChangeListener{})

	cfgMap := make(map[string]interface{})
	cache := client.GetConfigCache(storage.GetDefaultNamespace())
	cache.Range(func(key, value interface{}) bool {
		cfgMap[key.(string)] = value
		return true
	})

	err := refreshViperConfig(cfgMap)
	if err == nil {
		apolloConfig := &AdapterApollo{
			client: client,
		}
		return apolloConfig, nil
	} else {
		return nil, nil
	}
}

func refreshViperConfig(cfgMap interface{}) error {
	cfgBytes, err := json.Marshal(cfgMap)
	if err != nil {
		return err
	}

	viper.SetConfigType("json")
	if err := viper.ReadConfig(bytes.NewBuffer(cfgBytes)); err != nil {
		return fmt.Errorf("viper.ReadConfig Err: %v", err)
	}

	return nil
}

type AdapterApollo struct {
	client agollo.Client
}

func (adapterApollo AdapterApollo) Available(ctx context.Context, resource ...string) (ok bool) {
	return adapterApollo.client != nil
}

func (adapterApollo AdapterApollo) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	if strings.Contains(pattern, ".") {
		return viper.GetString(pattern), nil
	} else {
		return viper.AllSettings()[pattern], nil
	}
}

func (adapterApollo AdapterApollo) Data(ctx context.Context) (data map[string]interface{}, err error) {
	return viper.AllSettings(), nil
}

type ConfigChangeListener struct {
}

// hot update
func (c *ConfigChangeListener) OnChange(event *storage.ChangeEvent) {
	cfgBytes, _ := json.Marshal(event)
	g.Log().Infof(gctx.New(), "apollo change, %s", string(cfgBytes))

	cfgMap := make(map[string]interface{})
	for key, value := range event.Changes {
		cfgMap[key] = value.NewValue
	}
	_ = refreshViperConfig(cfgMap)
}

func (c *ConfigChangeListener) OnNewestChange(event *storage.FullChangeEvent) {

}

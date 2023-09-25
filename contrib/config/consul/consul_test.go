// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package consul_test

import (
	"testing"
	"time"

	consul "github.com/gogf/gf/contrib/config/consul/v2"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"

	"github.com/gogf/gf/v2/frame/g"
)

func TestConsul(t *testing.T) {
	ctx := gctx.GetInitCtx()
	gtest.C(t, func(t *gtest.T) {
		configuration := consul.Config{
			ConsulConfig: api.Config{
				Address:    "127.0.0.1:8500",
				Scheme:     "http",
				Datacenter: "dc1",
				Transport:  cleanhttp.DefaultPooledTransport(),
				Token:      "3f8aeba2-f1f7-42d0-b912-fcb041d4546d",
			},
			Path:  "server/message",
			Watch: true,
		}

		var configValue string

		configValue = `redis:
  addr: 127.0.0.1:6379`

		// Write test configuration
		consulClient, err := api.NewClient(&configuration.ConsulConfig)
		t.AssertNil(err)
		kv := consulClient.KV()
		_, err = kv.Put(&api.KVPair{Key: configuration.Path, Value: []byte(configValue)}, nil)
		t.AssertNil(err)

		// Create gcfg.Adapter
		adapter, err := consul.New(ctx, configuration)
		t.AssertNil(err)
		conf := g.Cfg(guid.S())
		conf.SetAdapter(adapter)

		t.Assert(conf.Available(ctx), true)

		v, err := conf.Get(ctx, "redis.addr")
		t.AssertNil(err)
		t.Assert(v.String(), "127.0.0.1:6379")

		m, err := conf.Data(ctx)
		t.AssertNil(err)
		t.AssertGT(len(m), 0)
		// g.Dump(m)

		// Test changes after modifying configuration
		configValue = `redis:
  addr: localhost:6379`
		_, err = kv.Put(&api.KVPair{Key: configuration.Path, Value: []byte(configValue)}, nil)
		t.AssertNil(err)

		time.Sleep(time.Second)

		v, err = conf.Get(ctx, "redis.addr")
		t.AssertNil(err)
		t.Assert(v.String(), "localhost:6379")

		m, err = conf.Data(ctx)
		t.AssertNil(err)
		g.Dump(m)
	})
}

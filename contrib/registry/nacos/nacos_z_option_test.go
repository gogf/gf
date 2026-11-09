// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos_test

import (
	"sync"
	"testing"

	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/registry/nacos/v2"
)

// fakeNamingClient records register params without requiring a nacos server.
type fakeNamingClient struct {
	naming_client.INamingClient
	mu     sync.Mutex
	params []vo.BatchRegisterInstanceParam
}

// BatchRegisterInstance records the given param.
func (c *fakeNamingClient) BatchRegisterInstance(param vo.BatchRegisterInstanceParam) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.params = append(c.params, param)
	return true, nil
}

func TestRegistry_Options(t *testing.T) {
	var (
		ctx      = gctx.GetInitCtx()
		client   = &fakeNamingClient{}
		metadata = map[string]string{"env": "test"}
		reg      = nacos.NewWithClient(
			client,
			nacos.WithClusterName("cluster"),
			nacos.WithGroupName("group"),
			nacos.WithDefaultEndpoint("10.0.0.1:9000"),
			nacos.WithDefaultMetadata(metadata),
		)
		svc = &gsvc.LocalService{
			Name:      "svc",
			Endpoints: gsvc.NewEndpoints("127.0.0.1:8888"),
			Metadata:  map[string]any{"protocol": "https"},
		}
		wg sync.WaitGroup
	)
	gtest.C(t, func(t *gtest.T) {
		// Changing the caller's map after creation must not affect or race with Register.
		wg.Add(1)
		go func() {
			defer wg.Done()
			metadata["env"] = "changed"
		}()
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := reg.Register(ctx, svc)
				t.AssertNil(err)
			}()
		}
		wg.Wait()

		t.Assert(len(client.params), 10)
		p := client.params[0]
		t.Assert(p.GroupName, "group")
		t.Assert(len(p.Instances), 1)
		t.Assert(p.Instances[0].Ip, "10.0.0.1")
		t.Assert(p.Instances[0].Port, 9000)
		t.Assert(p.Instances[0].ClusterName, "cluster")
		t.Assert(p.Instances[0].GroupName, "group")
		t.Assert(p.Instances[0].Metadata, map[string]string{"env": "test", "protocol": "https"})
	})
}

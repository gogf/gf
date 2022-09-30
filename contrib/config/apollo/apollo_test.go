// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package apollo_test

import (
	"testing"

	"github.com/gogf/gf/contrib/config/apollo/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

var (
	ctx     = gctx.GetInitCtx()
	appId   = "SampleApp"
	cluster = "default"
	ip      = "http://localhost:8080"
)

func TestApollo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		adapter, err := apollo.New(ctx, apollo.Config{
			AppID:   appId,
			IP:      ip,
			Cluster: cluster,
		})
		t.AssertNil(err)
		config := g.Cfg(guid.S())
		config.SetAdapter(adapter)

		t.Assert(config.Available(ctx), true)
		t.Assert(config.Available(ctx, "non-exist"), false)

		v, err := config.Get(ctx, `server.address`)
		t.AssertNil(err)
		t.Assert(v.String(), ":8000")

		m, err := config.Data(ctx)
		t.AssertNil(err)
		t.AssertGT(len(m), 0)
	})
}

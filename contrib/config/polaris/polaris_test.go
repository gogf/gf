// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package polaris_test

// import (
// 	"fmt"
// 	"testing"
//
// 	"github.com/gogf/gf/contrib/config/polaris/v2"
// 	"github.com/gogf/gf/v2/frame/g"
// 	"github.com/gogf/gf/v2/os/gctx"
// 	"github.com/gogf/gf/v2/test/gtest"
// 	"github.com/gogf/gf/v2/util/guid"
// )
//
// var (
// 	ctx       = gctx.GetInitCtx()
// 	namespace = "default"
// 	fileGroup = "goframe"
// 	fileName  = "config.yaml"
// 	path      = "testdata/polaris.yaml"
// 	logDir    = "/tmp/polaris/log"
// )
//
// func TestPolaris(t *testing.T) {
// 	gtest.C(t, func(t *gtest.T) {
// 		adapter, err := polaris.New(ctx, polaris.Config{
// 			Namespace: namespace,
// 			FileGroup: fileGroup,
// 			FileName:  fileName,
// 			Path:      path,
// 			LogDir:    logDir,
// 			Watch:     true,
// 		})
// 		t.AssertNil(err)
// 		config := g.Cfg(guid.S())
// 		config.SetAdapter(adapter)
//
// 		fmt.Println(adapter.Get(ctx, "server.address"))
//
// 		t.Assert(config.Available(ctx), true)
// 		t.Assert(config.Available(ctx, "non-exist"), false)
//
// 		v, err := config.Get(ctx, `server.address`)
// 		t.AssertNil(err)
// 		t.Assert(v.String(), ":8000")
//
// 		m, err := config.Data(ctx)
// 		t.AssertNil(err)
// 		t.AssertGT(len(m), 0)
// 	})
// }

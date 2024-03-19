// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package boot

import (
	"github.com/wangyougui/gf/contrib/config/kubecm/v2"
	"github.com/wangyougui/gf/v2/frame/g"
	"github.com/wangyougui/gf/v2/os/gctx"
)

const (
	configmapName       = "test-configmap"
	dataItemInConfigmap = "config.yaml"
)

func init() {
	var (
		err error
		ctx = gctx.GetInitCtx()
	)
	// Create kubecm Client that implements gcfg.Adapter.
	adapter, err := kubecm.New(gctx.GetInitCtx(), kubecm.Config{
		ConfigMap: configmapName,
		DataItem:  dataItemInConfigmap,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}

	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)
}

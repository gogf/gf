// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package boot

import (
	"github.com/gogf/gf/contrib/config/polaris/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func init() {
	var (
		ctx       = gctx.GetInitCtx()
		namespace = "default"
		fileGroup = "TestGroup"
		fileName  = "config.yaml"
		path      = "manifest/config/polaris.yaml"
		logDir    = "/tmp/polaris/log"
	)
	// Create polaris Client that implements gcfg.Adapter.
	adapter, err := polaris.New(ctx, polaris.Config{
		Namespace: namespace,
		FileGroup: fileGroup,
		FileName:  fileName,
		Path:      path,
		LogDir:    logDir,
		Watch:     true,
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	// Change the adapter of default configuration instance.
	g.Cfg().SetAdapter(adapter)
}

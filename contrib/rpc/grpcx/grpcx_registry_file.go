// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gogf/gf/contrib/registry/file/v2"
)

// autoLoadAndRegisterFileRegistry checks and registers ETCD service as default service registry
// if no registry is registered previously.
func autoLoadAndRegisterFileRegistry() {
	// It ignores etcd registry if any registry already registered.
	if gsvc.GetRegistry() != nil {
		return
	}
	var (
		ctx           = gctx.GetInitCtx()
		directoryPath = gfile.Temp("gsvc")
		fileRegistry  = file.New(directoryPath)
	)

	g.Log().Debugf(
		ctx,
		`set default registry using file registry as no custom registry set, path: %s`,
		directoryPath,
	)
	Resolver.Register(fileRegistry)
}

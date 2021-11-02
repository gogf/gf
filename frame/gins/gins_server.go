// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	frameCoreComponentNameServer = "gf.core.component.server"
	configNodeNameServer         = "server"
)

// Server returns an instance of http server with specified name.
// Note that it panics if any error occurs duration instance creating.
func Server(name ...interface{}) *ghttp.Server {
	var (
		ctx          = context.Background()
		instanceName = ghttp.DefaultServerName
		instanceKey  = fmt.Sprintf("%s.%v", frameCoreComponentNameServer, name)
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = gconv.String(name[0])
	}
	return localInstances.GetOrSetFuncLock(instanceKey, func() interface{} {
		s := ghttp.GetServer(instanceName)
		// It ignores returned error to avoid file no found error while it's not necessary.
		var (
			serverConfigMap       map[string]interface{}
			serverLoggerConfigMap map[string]interface{}
			configNodeName        = configNodeNameServer
		)
		if configData, _ := Config().Data(ctx); len(configData) > 0 {
			if v, _ := gutil.MapPossibleItemByKey(configData, configNodeNameServer); v != "" {
				configNodeName = v
			}
		}
		// Server configuration.
		certainConfigNodeName := fmt.Sprintf(`%s.%s`, configNodeName, s.GetName())
		if v, _ := Config().Get(ctx, certainConfigNodeName); !v.IsEmpty() {
			serverConfigMap = v.Map()
		}
		if len(serverConfigMap) == 0 {
			if v, _ := Config().Get(ctx, configNodeName); !v.IsEmpty() {
				serverConfigMap = v.Map()
			}
		}
		if len(serverConfigMap) > 0 {
			if err := s.SetConfigWithMap(serverConfigMap); err != nil {
				panic(err)
			}
		} else {
			// The configuration is not necessary, so it just prints internal logs.
			intlog.Printf(ctx, `missing configuration for HTTP server "%s"`, instanceName)
		}

		// Server logger configuration checks.
		serverLoggerNodeName := fmt.Sprintf(`%s.%s.%s`, configNodeName, s.GetName(), configNodeNameLogger)
		if v, _ := Config().Get(ctx, serverLoggerNodeName); !v.IsEmpty() {
			serverLoggerConfigMap = v.Map()
		}
		if len(serverLoggerConfigMap) > 0 {
			if err := s.Logger().SetConfigWithMap(serverLoggerConfigMap); err != nil {
				panic(err)
			}
		}
		// As it might use template feature,
		// it initializes the view instance as well.
		_ = getViewInstance()
		return s
	}).(*ghttp.Server)
}

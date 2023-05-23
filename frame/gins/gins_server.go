// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/internal/consts"
	"github.com/gogf/gf/v2/internal/instance"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// Server returns an instance of http server with specified name.
// Note that it panics if any error occurs duration instance creating.
func Server(name ...interface{}) *ghttp.Server {
	var (
		err          error
		ctx          = context.Background()
		instanceName = ghttp.DefaultServerName
		instanceKey  = fmt.Sprintf("%s.%v", frameCoreComponentNameServer, name)
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = gconv.String(name[0])
	}
	return instance.GetOrSetFuncLock(instanceKey, func() interface{} {
		server := ghttp.GetServer(instanceName)
		if Config().Available(ctx) {
			// Server initialization from configuration.
			var (
				configMap             map[string]interface{}
				serverConfigMap       map[string]interface{}
				serverLoggerConfigMap map[string]interface{}
				configNodeName        string
			)
			if configMap, err = Config().Data(ctx); err != nil {
				intlog.Errorf(ctx, `retrieve config data map failed: %+v`, err)
			}
			// Find possible server configuration item by possible names.
			if len(configMap) > 0 {
				if v, _ := gutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameServer); v != "" {
					configNodeName = v
				}
				if configNodeName == "" {
					if v, _ := gutil.MapPossibleItemByKey(configMap, consts.ConfigNodeNameServerSecondary); v != "" {
						configNodeName = v
					}
				}
			}
			// Automatically retrieve configuration by instance name.
			serverConfigMap = Config().MustGet(
				ctx,
				fmt.Sprintf(`%s.%s`, configNodeName, instanceName),
			).Map()
			if len(serverConfigMap) == 0 {
				serverConfigMap = Config().MustGet(ctx, configNodeName).Map()
			}
			if len(serverConfigMap) > 0 {
				if err = server.SetConfigWithMap(serverConfigMap); err != nil {
					panic(err)
				}
			} else {
				// The configuration is not necessary, so it just prints internal logs.
				intlog.Printf(
					ctx,
					`missing configuration from configuration component for HTTP server "%s"`,
					instanceName,
				)
			}
			// Server logger configuration checks.
			serverLoggerConfigMap = Config().MustGet(
				ctx,
				fmt.Sprintf(`%s.%s.%s`, configNodeName, instanceName, consts.ConfigNodeNameLogger),
			).Map()
			if len(serverLoggerConfigMap) == 0 && len(serverConfigMap) > 0 {
				serverLoggerConfigMap = gconv.Map(serverConfigMap[consts.ConfigNodeNameLogger])
			}
			if len(serverLoggerConfigMap) > 0 {
				if err = server.Logger().SetConfigWithMap(serverLoggerConfigMap); err != nil {
					panic(err)
				}
			}
		}
		// The server name is necessary. It sets a default server name is it is not configured.
		if server.GetName() == "" || server.GetName() == ghttp.DefaultServerName {
			server.SetName(instanceName)
		}
		// As it might use template feature,
		// it initializes the view instance as well.
		_ = getViewInstance()
		return server
	}).(*ghttp.Server)
}

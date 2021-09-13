// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gutil"
)

const (
	frameCoreComponentNameServer = "gf.core.component.server"
	configNodeNameServer         = "server"
)

// Server returns an instance of http server with specified name.
func Server(name ...interface{}) *ghttp.Server {
	instanceKey := fmt.Sprintf("%s.%v", frameCoreComponentNameServer, name)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		s := ghttp.GetServer(name...)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var (
				serverConfigMap       map[string]interface{}
				serverLoggerConfigMap map[string]interface{}
			)
			nodeKey, _ := gutil.MapPossibleItemByKey(Config().GetMap("."), configNodeNameServer)
			if nodeKey == "" {
				nodeKey = configNodeNameServer
			}
			// Server configuration.
			serverConfigMap = Config().GetMap(fmt.Sprintf(`%s.%s`, nodeKey, s.GetName()))
			if len(serverConfigMap) == 0 {
				serverConfigMap = Config().GetMap(nodeKey)
			}
			if len(serverConfigMap) > 0 {
				if err := s.SetConfigWithMap(serverConfigMap); err != nil {
					panic(err)
				}
			}
			// Server logger configuration.
			serverLoggerConfigMap = Config().GetMap(
				fmt.Sprintf(`%s.%s.%s`, nodeKey, s.GetName(), configNodeNameLogger),
			)
			if len(serverLoggerConfigMap) > 0 {
				if err := s.Logger().SetConfigWithMap(serverLoggerConfigMap); err != nil {
					panic(err)
				}
			}
			// As it might use template feature,
			// it initialize the view instance as well.
			_ = getViewInstance()
		}
		return s
	}).(*ghttp.Server)
}

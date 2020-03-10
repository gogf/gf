// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
)

const (
	gFRAME_CORE_COMPONENT_NAME_SERVER = "gf.core.component.server"
)

// Server returns an instance of http server with specified name.
func Server(name ...interface{}) *ghttp.Server {
	instanceKey := fmt.Sprintf("%s.%v", gFRAME_CORE_COMPONENT_NAME_SERVER, name)
	return instances.GetOrSetFunc(instanceKey, func() interface{} {
		s := ghttp.GetServer(name...)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			// It firstly searches the configuration of the instance name.
			if m = Config().GetMap(fmt.Sprintf(`server.%s`, s.GetName())); m == nil {
				// If the configuration for the instance does not exist,
				// it uses the default server configuration.
				m = Config().GetMap("server")
			}
			if m != nil {
				if err := s.SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
			// As it might use template feature,
			// it initialize the view instance as well.
			View()
		}
		return s
	}).(*ghttp.Server)
}

// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

const (
	gFRAME_CORE_COMPONENT_NAME_SERVER = "gf.core.component.server"
	gSERVER_NODE_NAME                 = "server"
)

// Server returns an instance of http server with specified name.
func Server(name ...interface{}) *ghttp.Server {
	instanceKey := fmt.Sprintf("%s.%v", gFRAME_CORE_COMPONENT_NAME_SERVER, name)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		s := ghttp.GetServer(name...)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			// It firstly searches the configuration of the instance name.
			if _, v := gutil.MapPossibleItemByKey(
				Config().GetMap("."),
				fmt.Sprintf(`%s.%s`, gSERVER_NODE_NAME, s.GetName()),
			); v != nil {
				m = gconv.Map(v)
			} else {
				// If the configuration for the instance does not exist,
				// it uses the default server configuration.
				if _, v := gutil.MapPossibleItemByKey(
					Config().GetMap("."),
					gSERVER_NODE_NAME,
				); v != nil {
					m = gconv.Map(v)
				}
			}
			if len(m) > 0 {
				if err := s.SetConfigWithMap(m); err != nil {
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

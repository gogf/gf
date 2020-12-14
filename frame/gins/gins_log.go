// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gutil"
)

const (
	gFRAME_CORE_COMPONENT_NAME_LOGGER = "gf.core.component.logger"
	gLOGGER_NODE_NAME                 = "logger"
)

// Log returns an instance of glog.Logger.
// The parameter <name> is the name for the instance.
func Log(name ...string) *glog.Logger {
	instanceName := glog.DefaultName
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_LOGGER, instanceName)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		logger := glog.Instance(instanceName)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			nodeKey, _ := gutil.MapPossibleItemByKey(Config().GetMap("."), gLOGGER_NODE_NAME)
			if nodeKey == "" {
				nodeKey = gLOGGER_NODE_NAME
			}
			m = Config().GetMap(fmt.Sprintf(`%s.%s`, nodeKey, instanceName))
			if len(m) == 0 {
				m = Config().GetMap(nodeKey)
			}
			if len(m) > 0 {
				if err := logger.SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
		}
		return logger
	}).(*glog.Logger)
}

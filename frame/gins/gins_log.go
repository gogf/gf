// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	frameCoreComponentNameLogger = "gf.core.component.logger"
	configNodeNameLogger         = "logger"
)

// Log returns an instance of glog.Logger.
// The parameter `name` is the name for the instance.
// Note that it panics if any error occurs duration instance creating.
func Log(name ...string) *glog.Logger {
	var (
		ctx          = context.Background()
		instanceName = glog.DefaultName
	)
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameLogger, instanceName)
	return localInstances.GetOrSetFuncLock(instanceKey, func() interface{} {
		logger := glog.Instance(instanceName)
		// To avoid file no found error while it's not necessary.
		var (
			configMap      map[string]interface{}
			loggerNodeName = configNodeNameLogger
		)
		// Try to find possible `loggerNodeName` in case-insensitive way.
		if configData, _ := Config().Data(ctx); len(configData) > 0 {
			if v, _ := gutil.MapPossibleItemByKey(configData, configNodeNameLogger); v != "" {
				loggerNodeName = v
			}
		}
		// Retrieve certain logger configuration by logger name.
		certainLoggerNodeName := fmt.Sprintf(`%s.%s`, loggerNodeName, instanceName)
		if v, _ := Config().Get(ctx, certainLoggerNodeName); !v.IsEmpty() {
			configMap = v.Map()
		}
		// Retrieve global logger configuration if configuration for certain logger name does not exist.
		if len(configMap) == 0 {
			if v, _ := Config().Get(ctx, loggerNodeName); !v.IsEmpty() {
				configMap = v.Map()
			}
		}
		// Set logger config if config map is not empty.
		if len(configMap) > 0 {
			if err := logger.SetConfigWithMap(configMap); err != nil {
				panic(err)
			}
		}
		return logger
	}).(*glog.Logger)
}

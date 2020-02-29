// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/util/gconv"
)

const (
	gFRAME_CORE_COMPONENT_NAME_REDIS = "gf.core.component.redis"
)

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gredis.Redis {
	config := Config()
	group := gredis.DEFAULT_GROUP_NAME
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_REDIS, group)
	result := instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		// If already configured, it returns the redis instance.
		if _, ok := gredis.GetConfig(group); ok {
			return gredis.Instance(group)
		}
		// Or else, it parses the default configuration file and returns a new redis instance.
		if m := config.GetMap("redis"); m != nil {
			if v, ok := m[group]; ok {
				redisConfig, err := gredis.ConfigFromStr(gconv.String(v))
				if err != nil {
					panic(err)
				}
				addConfigMonitor(instanceKey, config)
				return gredis.New(redisConfig)
			} else {
				panic(fmt.Sprintf(`configuration for redis not found for group "%s"`, group))
			}
		} else {
			panic(fmt.Sprintf(`incomplete configuration for redis: "redis" node not found in config file "%s"`, config.FilePath()))
		}
		return nil
	})
	if result != nil {
		return result.(*gredis.Redis)
	}
	return nil
}

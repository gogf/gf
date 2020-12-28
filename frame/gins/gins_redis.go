// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

const (
	frameCoreComponentNameRedis = "gf.core.component.redis"
	configNodeNameRedis         = "redis"
)

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gredis.Redis {
	config := Config()
	group := gredis.DefaultGroupName
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameRedis, group)
	result := instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		// If already configured, it returns the redis instance.
		if _, ok := gredis.GetConfig(group); ok {
			return gredis.Instance(group)
		}
		// Or else, it parses the default configuration file and returns a new redis instance.
		var m map[string]interface{}
		if _, v := gutil.MapPossibleItemByKey(Config().GetMap("."), configNodeNameRedis); v != nil {
			m = gconv.Map(v)
		}
		if len(m) > 0 {
			if v, ok := m[group]; ok {
				redisConfig, err := gredis.ConfigFromStr(gconv.String(v))
				if err != nil {
					panic(err)
				}
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

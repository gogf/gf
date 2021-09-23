// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
)

const (
	frameCoreComponentNameRedis = "gf.core.component.redis"
	configNodeNameRedis         = "redis"
)

// Redis returns an instance of redis client with specified configuration group name.
// Note that it panics if any error occurs duration instance creating.
func Redis(name ...string) *gredis.Redis {
	var (
		ctx   = context.Background()
		group = gredis.DefaultGroupName
	)
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameRedis, group)
	result := localInstances.GetOrSetFuncLock(instanceKey, func() interface{} {
		// If already configured, it returns the redis instance.
		if _, ok := gredis.GetConfig(group); ok {
			return gredis.Instance(group)
		}
		// Or else, it parses the default configuration file and returns a new redis instance.
		var (
			configMap map[string]interface{}
		)

		if configData, err := Config().Data(ctx); err != nil {
			panic(gerror.WrapCode(gcode.CodeOperationFailed, err, `retrieving redis configuration failed`))
		} else {
			if _, v := gutil.MapPossibleItemByKey(configData, configNodeNameRedis); v != nil {
				configMap = gconv.Map(v)
			}
		}

		if len(configMap) > 0 {
			if v, ok := configMap[group]; ok {
				redisConfig, err := gredis.ConfigFromStr(gconv.String(v))
				if err != nil {
					panic(err)
				}
				redisClient, err := gredis.New(redisConfig)
				if err != nil {
					panic(err)
				}
				return redisClient
			} else {
				panic(fmt.Sprintf(`missing configuration for redis group "%s"`, group))
			}
		} else {
			panic(`missing configuration for redis: "redis" node not found`)
		}
		return nil
	})
	if result != nil {
		return result.(*gredis.Redis)
	}
	return nil
}

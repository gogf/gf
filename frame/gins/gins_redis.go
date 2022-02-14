// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	frameCoreComponentNameRedis = "gf.core.component.redis"
	configNodeNameRedis         = "redis"
)

// Redis returns an instance of redis client with specified configuration group name.
// Note that it panics if any error occurs duration instance creating.
func Redis(name ...string) *gredis.Redis {
	var (
		err   error
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
		if Config().Available(ctx) {
			var (
				configMap   map[string]interface{}
				redisConfig *gredis.Config
				redisClient *gredis.Redis
			)
			if configMap, err = Config().Data(ctx); err != nil {
				intlog.Errorf(ctx, `retrieve config data map failed: %+v`, err)
			}
			if _, v := gutil.MapPossibleItemByKey(configMap, configNodeNameRedis); v != nil {
				configMap = gconv.Map(v)
			}
			if len(configMap) > 0 {
				if v, ok := configMap[group]; ok {
					if redisConfig, err = gredis.ConfigFromMap(gconv.Map(v)); err != nil {
						panic(err)
					}
				}
				intlog.Printf(ctx, `missing configuration for redis group "%s"`, group)
			} else {
				intlog.Printf(ctx, `missing configuration for redis: "redis" node not found`)
			}
			if redisClient, err = gredis.New(redisConfig); err != nil {
				panic(err)
			}
			return redisClient
		}
		return nil
	})
	if result != nil {
		return result.(*gredis.Redis)
	}
	return nil
}

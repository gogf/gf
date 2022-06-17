// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/consts"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	frameCoreComponentNameDatabase = "gf.core.component.database"
)

// Database returns an instance of database ORM object with specified configuration group name.
// Note that it panics if any error occurs duration instance creating.
func Database(name ...string) gdb.DB {
	var (
		ctx   = context.Background()
		group = gdb.DefaultGroupName
	)

	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", frameCoreComponentNameDatabase, group)
	db := localInstances.GetOrSetFuncLock(instanceKey, func() interface{} {
		// It ignores returned error to avoid file no found error while it's not necessary.
		var (
			configMap     map[string]interface{}
			configNodeKey = consts.ConfigNodeNameDatabase
		)
		// It firstly searches the configuration of the instance name.
		if configData, _ := Config().Data(ctx); len(configData) > 0 {
			if v, _ := gutil.MapPossibleItemByKey(configData, consts.ConfigNodeNameDatabase); v != "" {
				configNodeKey = v
			}
		}
		if v, _ := Config().Get(ctx, configNodeKey); !v.IsEmpty() {
			configMap = v.Map()
		}
		// No configuration found, it formats and panics error.
		if len(configMap) == 0 && !gdb.IsConfigured() {
			// File configuration object checks.
			var err error
			if fileConfig, ok := Config().GetAdapter().(*gcfg.AdapterFile); ok {
				if _, err = fileConfig.GetFilePath(); err != nil {
					panic(gerror.WrapCode(gcode.CodeMissingConfiguration, err,
						`configuration not found, did you miss the configuration file or the misspell the configuration file name`,
					))
				}
			}
			// Panic if nothing found in Config object or in gdb configuration.
			if len(configMap) == 0 && !gdb.IsConfigured() {
				panic(gerror.NewCodef(
					gcode.CodeMissingConfiguration,
					`database initialization failed: configuration missing for database node "%s"`,
					consts.ConfigNodeNameDatabase,
				))
			}
		}

		if len(configMap) == 0 {
			configMap = make(map[string]interface{})
		}
		// Parse `m` as map-slice and adds it to global configurations for package gdb.
		for g, groupConfig := range configMap {
			cg := gdb.ConfigGroup{}
			switch value := groupConfig.(type) {
			case []interface{}:
				for _, v := range value {
					if node := parseDBConfigNode(v); node != nil {
						cg = append(cg, *node)
					}
				}
			case map[string]interface{}:
				if node := parseDBConfigNode(value); node != nil {
					cg = append(cg, *node)
				}
			}
			if len(cg) > 0 {
				if gdb.GetConfig(group) == nil {
					intlog.Printf(ctx, "add configuration for group: %s, %#v", g, cg)
					gdb.SetConfigGroup(g, cg)
				} else {
					intlog.Printf(ctx, "ignore configuration as it already exists for group: %s, %#v", g, cg)
					intlog.Printf(ctx, "%s, %#v", g, cg)
				}
			}
		}
		// Parse `m` as a single node configuration,
		// which is the default group configuration.
		if node := parseDBConfigNode(configMap); node != nil {
			cg := gdb.ConfigGroup{}
			if node.Link != "" || node.Host != "" {
				cg = append(cg, *node)
			}
			if len(cg) > 0 {
				if gdb.GetConfig(group) == nil {
					intlog.Printf(ctx, "add configuration for group: %s, %#v", gdb.DefaultGroupName, cg)
					gdb.SetConfigGroup(gdb.DefaultGroupName, cg)
				} else {
					intlog.Printf(
						ctx,
						"ignore configuration as it already exists for group: %s, %#v",
						gdb.DefaultGroupName, cg,
					)
					intlog.Printf(ctx, "%s, %#v", gdb.DefaultGroupName, cg)
				}
			}
		}

		// Create a new ORM object with given configurations.
		if db, err := gdb.NewByGroup(name...); err == nil {
			// Initialize logger for ORM.
			var (
				loggerConfigMap map[string]interface{}
				loggerNodeName  = fmt.Sprintf("%s.%s", configNodeKey, consts.ConfigNodeNameLogger)
			)
			if v, _ := Config().Get(ctx, loggerNodeName); !v.IsEmpty() {
				loggerConfigMap = v.Map()
			}
			if len(loggerConfigMap) == 0 {
				if v, _ := Config().Get(ctx, configNodeKey); !v.IsEmpty() {
					loggerConfigMap = v.Map()
				}
			}
			if len(loggerConfigMap) > 0 {
				if err = db.GetLogger().SetConfigWithMap(loggerConfigMap); err != nil {
					panic(err)
				}
			}
			return db
		} else {
			// If panics, often because it does not find its configuration for given group.
			panic(err)
		}
		return nil
	})
	if db != nil {
		return db.(gdb.DB)
	}
	return nil
}

func parseDBConfigNode(value interface{}) *gdb.ConfigNode {
	nodeMap, ok := value.(map[string]interface{})
	if !ok {
		return nil
	}
	var (
		node = &gdb.ConfigNode{}
		err  = gconv.Struct(nodeMap, node)
	)
	if err != nil {
		panic(err)
	}
	// Find possible `Link` configuration content.
	if _, v := gutil.MapPossibleItemByKey(nodeMap, "Link"); v != nil {
		node.Link = gconv.String(v)
	}
	return node
}

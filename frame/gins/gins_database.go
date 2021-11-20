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
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	frameCoreComponentNameDatabase = "gf.core.component.database"
	configNodeNameDatabase         = "database"
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
			configNodeKey = configNodeNameDatabase
		)
		// It firstly searches the configuration of the instance name.
		if configData, _ := Config().Data(ctx); len(configData) > 0 {
			if v, _ := gutil.MapPossibleItemByKey(configData, configNodeNameDatabase); v != "" {
				configNodeKey = v
			}
		}
		if v, _ := Config().Get(ctx, configNodeKey); !v.IsEmpty() {
			configMap = v.Map()
		}
		if len(configMap) == 0 && !gdb.IsConfigured() {
			// File configuration object checks.
			var (
				err            error
				configFilePath string
			)
			if fileConfig, ok := Config().GetAdapter().(*gcfg.AdapterFile); ok {
				if configFilePath, err = fileConfig.GetFilePath(); configFilePath == "" {
					exampleFileName := "config.example.toml"
					if exampleConfigFilePath, _ := fileConfig.GetFilePath(exampleFileName); exampleConfigFilePath != "" {
						err = gerror.WrapCodef(
							gcode.CodeMissingConfiguration,
							err,
							`configuration file "%s" not found, but found "%s", did you miss renaming the example configuration file?`,
							fileConfig.GetFileName(),
							exampleFileName,
						)
					} else {
						err = gerror.WrapCodef(
							gcode.CodeMissingConfiguration,
							err,
							`configuration file "%s" not found, did you miss the configuration file or the misspell the configuration file name?`,
							fileConfig.GetFileName(),
						)
					}
					if err != nil {
						panic(err)
					}
				}
			}
			// Panic if nothing found in Config object or in gdb configuration.
			if len(configMap) == 0 && !gdb.IsConfigured() {
				err = gerror.WrapCodef(
					gcode.CodeMissingConfiguration,
					err,
					`database initialization failed: "%s" node not found, is configuration file or configuration node missing?`,
					configNodeNameDatabase,
				)
				panic(err)
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
					intlog.Printf(ctx, "ignore configuration as it already exists for group: %s, %#v", gdb.DefaultGroupName, cg)
					intlog.Printf(ctx, "%s, %#v", gdb.DefaultGroupName, cg)
				}
			}
		}

		// Create a new ORM object with given configurations.
		if db, err := gdb.New(name...); err == nil {
			// Initialize logger for ORM.
			var (
				loggerConfigMap map[string]interface{}
				loggerNodeName  = fmt.Sprintf("%s.%s", configNodeKey, configNodeNameLogger)
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
	node := &gdb.ConfigNode{}
	err := gconv.Struct(nodeMap, node)
	if err != nil {
		panic(err)
	}
	// Be compatible with old version.
	if _, v := gutil.MapPossibleItemByKey(nodeMap, "LinkInfo"); v != nil {
		node.Link = gconv.String(v)
	}
	if _, v := gutil.MapPossibleItemByKey(nodeMap, "Link"); v != nil {
		node.Link = gconv.String(v)
	}
	// Parse link syntax.
	if node.Link != "" && node.Type == "" {
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, node.Link)
		if len(match) == 3 {
			node.Type = gstr.Trim(match[1])
			node.Link = gstr.Trim(match[2])
		}
	}
	return node
}

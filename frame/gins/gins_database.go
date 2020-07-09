// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gins

import (
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gutil"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

const (
	gFRAME_CORE_COMPONENT_NAME_DATABASE = "gf.core.component.database"
	gDATABASE_NODE_NAME                 = "database"
)

// Database returns an instance of database ORM object
// with specified configuration group name.
func Database(name ...string) gdb.DB {
	group := gdb.DEFAULT_GROUP_NAME
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_DATABASE, group)
	db := instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		// Configuration already exists.
		if gdb.GetConfig(group) != nil {
			db, err := gdb.Instance(group)
			if err != nil {
				panic(err)
			}
			return db
		}
		var m map[string]interface{}
		// It firstly searches the configuration of the instance name.
		nodeKey, _ := gutil.MapPossibleItemByKey(Config().GetMap("."), gVIEWER_NODE_NAME)
		if nodeKey == "" {
			nodeKey = gDATABASE_NODE_NAME
		}
		if m = Config().GetMap(nodeKey); len(m) == 0 {
			panic(fmt.Sprintf(`database init failed: "%s" node not found, is config file or configuration missing?`, gDATABASE_NODE_NAME))
		}
		// Parse <m> as map-slice and adds it to gdb's global configurations.
		for group, groupConfig := range m {
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
				intlog.Printf("%s, %#v", group, cg)
				gdb.SetConfigGroup(group, cg)
			}
		}
		// Parse <m> as a single node configuration,
		// which is the default group configuration.
		if node := parseDBConfigNode(m); node != nil {
			cg := gdb.ConfigGroup{}
			if node.LinkInfo != "" || node.Host != "" {
				cg = append(cg, *node)
			}
			if len(cg) > 0 {
				intlog.Printf("%s, %#v", gdb.DEFAULT_GROUP_NAME, cg)
				gdb.SetConfigGroup(gdb.DEFAULT_GROUP_NAME, cg)
			}
		}

		if db, err := gdb.New(name...); err == nil {
			// Initialize logger for ORM.
			var m map[string]interface{}
			m = Config().GetMap(fmt.Sprintf("%s.%s", nodeKey, gLOGGER_NODE_NAME))
			if len(m) == 0 {
				m = Config().GetMap(nodeKey)
			}
			if len(m) > 0 {
				if err := db.GetLogger().SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
			return db
		} else {
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
	if _, v := gutil.MapPossibleItemByKey(nodeMap, "link"); v != nil {
		node.LinkInfo = gconv.String(v)
	}
	// Parse link syntax.
	if node.LinkInfo != "" && node.Type == "" {
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, node.LinkInfo)
		if len(match) == 3 {
			node.Type = gstr.Trim(match[1])
			node.LinkInfo = gstr.Trim(match[2])
		}
	}
	return node
}

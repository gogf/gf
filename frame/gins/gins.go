// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gins provides instances management and core components management.
package gins

import (
	"fmt"

	"github.com/gogf/gf/os/gfile"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/i18n/gi18n"
	"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gfsnotify"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gres"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

const (
	gFRAME_CORE_COMPONENT_NAME_REDIS    = "gf.core.component.redis"
	gFRAME_CORE_COMPONENT_NAME_DATABASE = "gf.core.component.database"
)

// 单例对象存储器
var instances = gmap.NewStrAnyMap(true)

// 获取单例对象
func Get(key string) interface{} {
	return instances.Get(key)
}

// 设置单例对象
func Set(key string, value interface{}) {
	instances.Set(key, value)
}

// 当键名存在时返回其键值，否则写入指定的键值
func GetOrSet(key string, value interface{}) interface{} {
	return instances.GetOrSet(key, value)
}

// 当键名存在时返回其键值，否则写入指定的键值，键值由指定的函数生成
func GetOrSetFunc(key string, f func() interface{}) interface{} {
	return instances.GetOrSetFunc(key, f)
}

// 与GetOrSetFunc不同的是，f是在写锁机制内执行
func GetOrSetFuncLock(key string, f func() interface{}) interface{} {
	return instances.GetOrSetFuncLock(key, f)
}

// 当键名不存在时写入，并返回true；否则返回false。
func SetIfNotExist(key string, value interface{}) bool {
	return instances.SetIfNotExist(key, value)
}

// View returns an instance of View with default settings.
// The parameter <name> is the name for the instance.
func View(name ...string) *gview.View {
	return gview.Instance(name...)
}

// Config returns an instance of View with default settings.
// The parameter <name> is the name for the instance.
func Config(name ...string) *gcfg.Config {
	return gcfg.Instance(name...)
}

// Resource returns an instance of Resource.
// The parameter <name> is the name for the instance.
func Resource(name ...string) *gres.Resource {
	return gres.Instance(name...)
}

// I18n returns an instance of gi18n.Manager.
// The parameter <name> is the name for the instance.
func I18n(name ...string) *gi18n.Manager {
	return gi18n.Instance(name...)
}

// Database returns an instance of database ORM object
// with specified configuration group name.
func Database(name ...string) gdb.DB {
	config := Config()
	group := gdb.DEFAULT_GROUP_NAME
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_DATABASE, group)
	db := instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		if gdb.GetConfig(group) != nil {
			db, err := gdb.Instance(group)
			if err != nil {
				glog.Error(err)
			}
			return db
		}
		m := config.GetMap("database")
		if m == nil {
			glog.Error(`database init failed: "database" node not found, is config file or configuration missing?`)
			return nil
		}
		// Parse <m> as map-slice.
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
				gdb.SetConfigGroup(group, cg)
			}
		}
		// Parse <m> as a single node configuration.
		if node := parseDBConfigNode(m); node != nil {
			cg := gdb.ConfigGroup{}
			if node.LinkInfo != "" || node.Host != "" {
				cg = append(cg, *node)
			}
			if len(cg) > 0 {
				gdb.SetConfigGroup(group, cg)
			}
		}
		addConfigMonitor(instanceKey, config)

		if db, err := gdb.New(name...); err == nil {
			return db
		} else {
			glog.Error(err)
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
		glog.Error(err)
	}
	if value, ok := nodeMap["link"]; ok {
		node.LinkInfo = gconv.String(value)
	}
	// Parse link syntax.
	if node.LinkInfo != "" && node.Type == "" {
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, node.LinkInfo)
		if len(match) == 3 {
			node.Type = match[1]
			node.LinkInfo = match[2]
		}
	}
	return node
}

// Redis returns an instance of redis client with specified configuration group name.
func Redis(name ...string) *gredis.Redis {
	config := Config()
	group := "default"
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
					glog.Error(err)
					return nil
				}
				addConfigMonitor(instanceKey, config)
				return gredis.New(redisConfig)
			} else {
				glog.Errorf(`configuration for redis not found for group "%s"`, group)
			}
		} else {
			glog.Errorf(`incomplete configuration for redis: "redis" node not found in config file "%s"`, config.FilePath())
		}
		return nil
	})
	if result != nil {
		return result.(*gredis.Redis)
	}
	return nil
}

func addConfigMonitor(key string, config *gcfg.Config) {
	if path := config.FilePath(); path != "" && gfile.Exists(path) {
		gfsnotify.Add(path, func(event *gfsnotify.Event) {
			instances.Remove(key)
		})
	}
}

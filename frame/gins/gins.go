// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gins provides instances and core components management.
//
// Note that it should not used glog.Panic* functions for panics if you do not want
// to log all the panics.
package gins

import (
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gutil"

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
	gFRAME_CORE_COMPONENT_NAME_LOGGER   = "gf.core.component.logger"
	gFRAME_CORE_COMPONENT_NAME_SERVER   = "gf.core.component.server"
	gFRAME_CORE_COMPONENT_NAME_VIEWER   = "gf.core.component.viewer"
	gFRAME_CORE_COMPONENT_NAME_DATABASE = "gf.core.component.database"
	gLOGGER_NODE_NAME                   = "logger"
	gVIEWER_NODE_NAME                   = "viewer"
)

var (
	// instances is the instance map for common used components.
	instances = gmap.NewStrAnyMap(true)
)

// Get returns the instance by given name.
func Get(name string) interface{} {
	return instances.Get(name)
}

// Set sets a instance object to the instance manager with given name.
func Set(name string, instance interface{}) {
	instances.Set(name, instance)
}

// GetOrSet returns the instance by name,
// or set instance to the instance manager if it does not exist and returns this instance.
func GetOrSet(name string, instance interface{}) interface{} {
	return instances.GetOrSet(name, instance)
}

// GetOrSetFunc returns the instance by name,
// or sets instance with returned value of callback function <f> if it does not exist
// and then returns this instance.
func GetOrSetFunc(name string, f func() interface{}) interface{} {
	return instances.GetOrSetFunc(name, f)
}

// GetOrSetFuncLock returns the instance by name,
// or sets instance with returned value of callback function <f> if it does not exist
// and then returns this instance.
//
// GetOrSetFuncLock differs with GetOrSetFunc function is that it executes function <f>
// with mutex.Lock of the hash map.
func GetOrSetFuncLock(name string, f func() interface{}) interface{} {
	return instances.GetOrSetFuncLock(name, f)
}

// SetIfNotExist sets <instance> to the map if the <name> does not exist, then returns true.
// It returns false if <name> exists, and <instance> would be ignored.
func SetIfNotExist(name string, instance interface{}) bool {
	return instances.SetIfNotExist(name, instance)
}

// View returns an instance of View with default settings.
// The parameter <name> is the name for the instance.
func View(name ...string) *gview.View {
	instanceName := gview.DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_VIEWER, instanceName)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		view := gview.Instance(instanceName)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			// It firstly searches the configuration of the instance name.
			if m = Config().GetMap(fmt.Sprintf(`%s.%s`, gVIEWER_NODE_NAME, instanceName)); m == nil {
				// If the configuration for the instance does not exist,
				// it uses the default view configuration.
				m = Config().GetMap(gVIEWER_NODE_NAME)
			}
			if m != nil {
				if err := view.SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
		}
		return view
	}).(*gview.View)
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

// Log returns an instance of glog.Logger.
// The parameter <name> is the name for the instance.
func Log(name ...string) *glog.Logger {
	instanceName := glog.DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		instanceName = name[0]
	}
	instanceKey := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_LOGGER, instanceName)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		logger := glog.Instance(instanceName)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			// It firstly searches the configuration of the instance name.
			if m = Config().GetMap(fmt.Sprintf(`%s.%s`, gLOGGER_NODE_NAME, instanceName)); m == nil {
				// If the configuration for the instance does not exist,
				// it uses the default logging configuration.
				m = Config().GetMap(gLOGGER_NODE_NAME)
			}
			if m != nil {
				if err := logger.SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
		}
		return logger
	}).(*glog.Logger)
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
		// Configuration already exists.
		if gdb.GetConfig(group) != nil {
			db, err := gdb.Instance(group)
			if err != nil {
				panic(err)
			}
			return db
		}
		m := config.GetMap("database")
		if m == nil {
			panic(`database init failed: "database" node not found, is config file or configuration missing?`)
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
				gdb.SetConfigGroup(gdb.DEFAULT_GROUP_NAME, cg)
			}
		}
		addConfigMonitor(instanceKey, config)

		if db, err := gdb.New(name...); err == nil {
			// Initialize logger for ORM.
			m := config.GetMap(fmt.Sprintf("database.%s", gLOGGER_NODE_NAME))
			if m == nil {
				m = config.GetMap(gLOGGER_NODE_NAME)
			}
			if m != nil {
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
			node.Type = match[1]
			node.LinkInfo = match[2]
		}
	}
	return node
}

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

// Server returns an instance of http server with specified name.
func Server(name ...interface{}) *ghttp.Server {
	instanceKey := fmt.Sprintf("%s.%v", gFRAME_CORE_COMPONENT_NAME_SERVER, name)
	return instances.GetOrSetFuncLock(instanceKey, func() interface{} {
		s := ghttp.GetServer(name...)
		// To avoid file no found error while it's not necessary.
		if Config().Available() {
			var m map[string]interface{}
			// It firstly searches the configuration of the instance name.
			if m = Config().GetMap(fmt.Sprintf(`server.%s`, s.GetName())); m == nil {
				// If the configuration for the instance does not exist,
				// it uses the default server configuration.
				m = Config().GetMap("server")
			}
			if m != nil {
				if err := s.SetConfigWithMap(m); err != nil {
					panic(err)
				}
			}
		}
		return s
	}).(*ghttp.Server)
}

func addConfigMonitor(key string, config *gcfg.Config) {
	if path := config.FilePath(); path != "" && gfile.Exists(path) {
		_, err := gfsnotify.Add(path, func(event *gfsnotify.Event) {
			instances.Remove(key)
		})
		if err != nil {
			intlog.Error(err)
		}
	}
}

// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gins provides instances and core components management.
package gins

import (
	"github.com/gogf/gf/container/gmap"

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

<<<<<<< HEAD
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
	//============================== If you have a cluster configuration, optimize the use of clustering
	if config.GetString("rediscluster."+group+".host") != "" && gredis.FlagBanCluster == false {
		clusters := RedisCluster(config, group)
		if clusters != nil {
			return clusters
		}
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

func RedisCluster(config *gcfg.Config, group string) *gredis.Redis {
	if gredis.FlagBanCluster {
		return nil
	}
	key := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_REDIS, group)
	result := instances.GetOrSetFuncLock(key, func() interface{} {
		if m := config.GetMap("rediscluster"); m != nil {
			// host1:port1,host2:port2
			if v, ok := m[group]; ok {
				lines := gconv.Map(v)
				hosts := strings.Split(gconv.String(lines["host"]), ",")
				return gredis.NewClusterClient(&gredis.ClusterOption{
					Nodes: hosts,
					Pwd:   gconv.String(lines["pwd"]),
				})

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
=======
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
>>>>>>> bd3e25adea5d01b7371b5122e67751262da52ff6
}

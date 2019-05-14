// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gins provides instances management and core components management.
package gins

import (
	"fmt"
	"github.com/gogf/gf/g/container/gmap"
	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/database/gredis"
	"github.com/gogf/gf/g/os/gcfg"
	"github.com/gogf/gf/g/os/gfsnotify"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/os/gview"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/text/gstr"
	"github.com/gogf/gf/g/util/gconv"
	"time"
)

const (
    gFRAME_CORE_COMPONENT_NAME_REDIS      = "gf.core.component.redis"
    gFRAME_CORE_COMPONENT_NAME_DATABASE   = "gf.core.component.database"
)

// 单例对象存储器
var instances = gmap.NewStrAnyMap()

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
// The param <name> is the name for the instance.
func View(name ...string) *gview.View {
	return gview.Instance(name ...)
}

// Config returns an instance of View with default settings.
// The param <name> is the name for the instance.
func Config(name ...string) *gcfg.Config {
    return gcfg.Instance(name ...)
}

// 数据库操作对象，使用了连接池
func Database(name ...string) gdb.DB {
    config := Config()
    group  := gdb.DEFAULT_GROUP_NAME
    if len(name) > 0 {
        group = name[0]
    }
    key := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_DATABASE, group)
    db  := instances.GetOrSetFuncLock(key, func() interface{} {
        if gdb.GetConfig(group) == nil {
            m := config.GetMap("database")
            if m == nil {
                glog.Error(`database init failed: "database" node not found, is config file or configuration missing?`)
                return nil
            }
            for group, v := range m {
                cg := gdb.ConfigGroup{}
                if list, ok := v.([]interface{}); ok {
                    for _, nodeValue := range list {
                        node    := gdb.ConfigNode{}
                        nodeMap := nodeValue.(map[string]interface{})
                        if value, ok := nodeMap["host"]; ok {
                            node.Host = gconv.String(value)
                        }
                        if value, ok := nodeMap["port"]; ok {
                            node.Port = gconv.String(value)
                        }
                        if value, ok := nodeMap["user"]; ok {
                            node.User = gconv.String(value)
                        }
                        if value, ok := nodeMap["pass"]; ok {
                            node.Pass = gconv.String(value)
                        }
                        if value, ok := nodeMap["name"]; ok {
                            node.Name = gconv.String(value)
                        }
                        if value, ok := nodeMap["type"]; ok {
                            node.Type = gconv.String(value)
                        }
                        if value, ok := nodeMap["role"]; ok {
                            node.Role = gconv.String(value)
                        }
                        if value, ok := nodeMap["charset"]; ok {
                            node.Charset = gconv.String(value)
                        }
                        if value, ok := nodeMap["priority"]; ok {
                            node.Priority = gconv.Int(value)
                        }
                        // Deprecated
                        if value, ok := nodeMap["linkinfo"]; ok {
                            node.LinkInfo = gconv.String(value)
                        }
                        // Deprecated
                        if value, ok := nodeMap["link-info"]; ok {
                            node.LinkInfo = gconv.String(value)
                        }
                        if value, ok := nodeMap["linkInfo"]; ok {
                            node.LinkInfo = gconv.String(value)
                        }
                        // Deprecated
                        if value, ok := nodeMap["max-idle"]; ok {
                            node.MaxIdleConnCount = gconv.Int(value)
                        }
                        if value, ok := nodeMap["maxIdle"]; ok {
                            node.MaxIdleConnCount = gconv.Int(value)
                        }
                        // Deprecated
                        if value, ok := nodeMap["max-open"]; ok {
                            node.MaxOpenConnCount = gconv.Int(value)
                        }
                        if value, ok := nodeMap["maxOpen"]; ok {
                            node.MaxOpenConnCount = gconv.Int(value)
                        }
                        // Deprecated
                        if value, ok := nodeMap["max-lifetime"]; ok {
                            node.MaxConnLifetime = gconv.Int(value)
                        }
                        if value, ok := nodeMap["maxLifetime"]; ok {
                            node.MaxConnLifetime = gconv.Int(value)
                        }
                        cg = append(cg, node)
                    }
                }
                gdb.AddConfigGroup(group, cg)
            }
            addConfigMonitor(key, config)
        }
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

// Redis操作对象，使用了连接池
func Redis(name...string) *gredis.Redis {
    config := Config()
    group  := "default"
    if len(name) > 0 {
        group = name[0]
    }
    key    := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_REDIS, group)
    result := instances.GetOrSetFuncLock(key, func() interface{} {
        if m := config.GetMap("redis"); m != nil {
            // host:port[,db,pass?maxIdle=x&maxActive=x&idleTimeout=x&maxConnLifetime=x]
            if v, ok := m[group]; ok {
                line := gconv.String(v)
                array, _ := gregex.MatchString(`(.+):(\d+),{0,1}(\d*),{0,1}(.*)\?(.+)`, line)
                if len(array) == 6 {
                    parse, _    := gstr.Parse(array[5])
                    redisConfig := gredis.Config{
                        Host : array[1],
                        Port : gconv.Int(array[2]),
                        Db   : gconv.Int(array[3]),
                        Pass : array[4],
                    }
                    if v, ok := parse["maxIdle"]; ok {
                        redisConfig.MaxIdle = gconv.Int(v)
                    }
                    if v, ok := parse["maxActive"]; ok {
                        redisConfig.MaxActive = gconv.Int(v)
                    }
                    if v, ok := parse["idleTimeout"]; ok {
                        redisConfig.IdleTimeout = gconv.Duration(v)*time.Second
                    }
                    if v, ok := parse["maxConnLifetime"]; ok {
                        redisConfig.MaxConnLifetime = gconv.Duration(v)*time.Second
                    }
                    addConfigMonitor(key, config)
                    return gredis.New(redisConfig)
                }
                array, _ = gregex.MatchString(`(.+):(\d+),{0,1}(\d*),{0,1}(.*)`, line)
                if len(array) == 5 {
                    addConfigMonitor(key, config)
                    return gredis.New(gredis.Config{
                        Host : array[1],
                        Port : gconv.Int(array[2]),
                        Db   : gconv.Int(array[3]),
                        Pass : array[4],
                    })
                } else {
                    glog.Errorfln(`invalid redis node configuration: "%s"`, line)
                }
            } else {
                glog.Errorfln(`configuration for redis not found for group "%s"`, group)
            }
        } else {
            glog.Errorfln(`incomplete configuration for redis: "redis" node not found in config file "%s"`, config.FilePath())
        }
        return nil
    })
    if result != nil {
        return result.(*gredis.Redis)
    }
    return nil
}

// 添加对单例对象的配置文件inotify监控
func addConfigMonitor(key string, config *gcfg.Config) {
    // 使用gfsnotify进行文件监控，当配置文件有任何变化时，清空对象单例缓存
    if path := config.FilePath(); path != "" {
        gfsnotify.Add(path, func(event *gfsnotify.Event) {
            instances.Remove(key)
        })
    }
}

// 模板内置方法：config
func funcConfig(pattern string, file...interface{}) string {
    return Config().GetString(pattern, file...)
}


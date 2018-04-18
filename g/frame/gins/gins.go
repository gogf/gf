// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 单例对象管理.
// 框架内置了一些核心对象获取方法，并且可以通过Set和Get方法实现IoC以及对内置核心对象的自定义替换
package gins

import (
    "strconv"
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/os/gcmd"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/database/gredis"
    "strings"
    "gitee.com/johng/gf/g/util/gconv"
)

const (
    gFRAME_CORE_COMPONENT_NAME_VIEW       = "gf.core.component.view"
    gFRAME_CORE_COMPONENT_NAME_CONFIG     = "gf.core.component.config"
    gFRAME_CORE_COMPONENT_NAME_DATABASE   = "gf.core.component.database"
    gFRAME_CORE_COMPONENT_NAME_REDIS      = "gf.core.component.redis"
)

// 单例对象存储器
var instances = gmap.NewStringInterfaceMap()

// 获取单例对象
func Get(k string) interface{} {
    return instances.Get(k)
}

// 设置单例对象
func Set(k string, v interface{}) {
    instances.Set(k, v)
}

// 自定义框架核心组件：View
func SetView(v *gview.View) {
    instances.Set(gFRAME_CORE_COMPONENT_NAME_VIEW, v)
}

// 自定义框架核心组件：Config
func SetConfig(v *gcfg.Config) {
    instances.Set(gFRAME_CORE_COMPONENT_NAME_CONFIG, v)
}

// 自定义框架核心组件：Database
func SetDatabase(v *gdb.Db, name...string) {
    dbCacheKey := gFRAME_CORE_COMPONENT_NAME_DATABASE
    if len(name) > 0 {
        dbCacheKey += name[0]
    }
    instances.Set(dbCacheKey, v)
}

// 核心对象：View
func View() *gview.View {
    result := Get(gFRAME_CORE_COMPONENT_NAME_VIEW)
    if result != nil {
        return result.(*gview.View)
    } else {
        path := gcmd.Option.Get("viewpath")
        if path == "" {
            path = genv.Get("gf.viewpath")
            if path == "" {
                path = gfile.SelfDir()
            }
        }
        view := gview.Get(path)
        Set(gFRAME_CORE_COMPONENT_NAME_VIEW, view)
        return view
    }
    return nil
}

// 核心对象：Config
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config() *gcfg.Config {
    result := Get(gFRAME_CORE_COMPONENT_NAME_CONFIG)
    if result != nil {
        return result.(*gcfg.Config)
    } else {
        path := gcmd.Option.Get("cfgpath")
        if path == "" {
            path = genv.Get("gf.cfgpath")
            if path == "" {
                path = gfile.SelfDir()
            }
        }
        config := gcfg.New(path)
        Set(gFRAME_CORE_COMPONENT_NAME_CONFIG, config)
        return config
    }
    return nil
}

// 核心对象：Database
func Database(name...string) *gdb.Db {
    dbCacheKey := gFRAME_CORE_COMPONENT_NAME_DATABASE
    if len(name) > 0 {
        dbCacheKey += name[0]
    }
    result := Get(dbCacheKey)
    if result != nil {
        return result.(*gdb.Db)
    } else {
        config := Config()
        if config == nil {
            return nil
        }
        if m := config.GetMap("database"); m != nil {
            for group, v := range m {
                if list, ok := v.([]interface{}); ok {
                    for _, nodev := range list {
                        node  := gdb.ConfigNode{}
                        nodem := nodev.(map[string]interface{})
                        if value, ok := nodem["host"]; ok {
                            node.Host = value.(string)
                        }
                        if value, ok := nodem["port"]; ok {
                            node.Port = value.(string)
                        }
                        if value, ok := nodem["user"]; ok {
                            node.User = value.(string)
                        }
                        if value, ok := nodem["pass"]; ok {
                            node.Pass = value.(string)
                        }
                        if value, ok := nodem["name"]; ok {
                            node.Name = value.(string)
                        }
                        if value, ok := nodem["type"]; ok {
                            node.Type = value.(string)
                        }
                        if value, ok := nodem["role"]; ok {
                            node.Role = value.(string)
                        }
                        if value, ok := nodem["charset"]; ok {
                            node.Charset = value.(string)
                        }
                        if value, ok := nodem["priority"]; ok {
                            node.Priority, _ = strconv.Atoi(value.(string))
                        }
                        gdb.AddConfigNode(group, node)
                    }
                }
            }

            if db, err := gdb.Instance(name...); err == nil {
                Set(dbCacheKey, db)
                return db
            } else {
                return nil
            }
        }
    }
    return nil
}

// 核心对象：Redis
func Redis(name...string) *gredis.Redis {
    group         := "default"
    redisCacheKey := gFRAME_CORE_COMPONENT_NAME_REDIS
    if len(name) > 0 {
        group          = name[0]
        redisCacheKey += name[0]
    }
    result := Get(redisCacheKey)
    if result != nil {
        return result.(*gredis.Redis)
    } else {
        config := Config()
        if config == nil {
            return nil
        }
        if m := config.GetMap("redis"); m != nil {
            if v, ok := m[group]; ok {
                array := strings.Split(gconv.String(v), ",")
                if len(array) > 1 {
                    redis := gredis.New(array[0], array[1])
                    Set(redisCacheKey, redis)
                    return redis
                }
            }
        }
    }
    return nil
}
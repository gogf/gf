// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gins provides instances management and some core components.
// 
// 单例对象管理.
// 框架内置了一些核心对象获取方法，并且可以通过Set和Get方法实现IoC以及对内置核心对象的自定义替换
package gins

import (
    "fmt"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/database/gdb"
    "github.com/gogf/gf/g/database/gredis"
    "github.com/gogf/gf/g/internal/cmdenv"
    "github.com/gogf/gf/g/os/gcfg"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/gfsnotify"
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gview"
    "github.com/gogf/gf/g/text/gstr"
    "github.com/gogf/gf/g/util/gconv"
    "github.com/gogf/gf/g/text/gregex"
    "time"
)

const (
    gFRAME_CORE_COMPONENT_NAME_VIEW       = "gf.core.component.view"
    gFRAME_CORE_COMPONENT_NAME_CONFIG     = "gf.core.component.config"
    gFRAME_CORE_COMPONENT_NAME_REDIS      = "gf.core.component.redis"
    gFRAME_CORE_COMPONENT_NAME_DATABASE   = "gf.core.component.database"
)

// 单例对象存储器
var instances = gmap.NewStringInterfaceMap()

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

// 核心对象：View
func View(name...string) *gview.View {
    group := "default"
    if len(name) > 0 {
        group = name[0]
    }
    key := fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_VIEW, group)
    return instances.GetOrSetFuncLock(key, func() interface{} {
        path := cmdenv.Get("gf.gview.path", gfile.Pwd()).String()
        view := gview.New(path)
        // 添加基于源码的搜索目录检索地址，常用于开发环境调试，只添加入口文件目录
        if p := gfile.MainPkgPath(); p != "" && gfile.Exists(p) {
            view.AddPath(p)
        }
        // 框架内置函数
        view.BindFunc("config", funcConfig)
        return view
    }).(*gview.View)
}

// 核心对象：Config
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config(file...string) *gcfg.Config {
    configFile := gcfg.DEFAULT_CONFIG_FILE
    if len(file) > 0 {
        configFile = file[0]
    }
    return instances.GetOrSetFuncLock(fmt.Sprintf("%s.%s", gFRAME_CORE_COMPONENT_NAME_CONFIG, configFile),
        func() interface{} {
            pwdPath  := gfile.Pwd()
            envPath  := cmdenv.Get("gf.gcfg.path").String()
            selfPath := gfile.SelfDir()
            mainPath := gfile.MainPkgPath()
            config   := gcfg.New(pwdPath, configFile)
            // 自定义的环境变量/启动参数路径，优先级最高，覆盖默认的工作目录
            if envPath != "" && gfile.Exists(envPath) {
                config.SetPath(envPath)
            }
            // 二进制文件执行目录
            if selfPath != "" && gfile.Exists(selfPath) {
                config.AddPath(selfPath)
            }
            // 开发环境源码main包目录
            if mainPath != "" && gfile.Exists(mainPath) {
                config.AddPath(mainPath)
            }
            return config
    }).(*gcfg.Config)
}

// 数据库操作对象，使用了连接池
func Database(name...string) gdb.DB {
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
                    for _, nodev := range list {
                        node  := gdb.ConfigNode{}
                        nodem := nodev.(map[string]interface{})
                        if value, ok := nodem["host"]; ok {
                            node.Host = gconv.String(value)
                        }
                        if value, ok := nodem["port"]; ok {
                            node.Port = gconv.String(value)
                        }
                        if value, ok := nodem["user"]; ok {
                            node.User = gconv.String(value)
                        }
                        if value, ok := nodem["pass"]; ok {
                            node.Pass = gconv.String(value)
                        }
                        if value, ok := nodem["name"]; ok {
                            node.Name = gconv.String(value)
                        }
                        if value, ok := nodem["type"]; ok {
                            node.Type = gconv.String(value)
                        }
                        if value, ok := nodem["role"]; ok {
                            node.Role = gconv.String(value)
                        }
                        if value, ok := nodem["charset"]; ok {
                            node.Charset = gconv.String(value)
                        }
                        if value, ok := nodem["priority"]; ok {
                            node.Priority = gconv.Int(value)
                        }
                        // Deprecated
                        if value, ok := nodem["linkinfo"]; ok {
                            node.Linkinfo = gconv.String(value)
                        }
                        if value, ok := nodem["linkInfo"]; ok {
                            node.Linkinfo = gconv.String(value)
                        }
                        // Deprecated
                        if value, ok := nodem["max-idle"]; ok {
                            node.MaxIdleConnCount = gconv.Int(value)
                        }
                        if value, ok := nodem["maxIdle"]; ok {
                            node.MaxIdleConnCount = gconv.Int(value)
                        }
                        // Deprecated
                        if value, ok := nodem["max-open"]; ok {
                            node.MaxOpenConnCount = gconv.Int(value)
                        }
                        if value, ok := nodem["maxOpen"]; ok {
                            node.MaxOpenConnCount = gconv.Int(value)
                        }
                        // Deprecated
                        if value, ok := nodem["max-lifetime"]; ok {
                            node.MaxConnLifetime = gconv.Int(value)
                        }
                        if value, ok := nodem["maxLifetime"]; ok {
                            node.MaxConnLifetime = gconv.Int(value)
                        }
                        cg = append(cg, node)
                    }
                }
                gdb.AddConfigGroup(group, cg)
            }
            // 使用gfsnotify进行文件监控，当配置文件有任何变化时，清空数据库配置缓存
            gfsnotify.Add(config.GetFilePath(), func(event *gfsnotify.Event) {
                instances.Remove(key)
            })
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
                    parse, _ := gstr.Parse(array[5])
                    config   := gredis.Config{
                        Host : array[1],
                        Port : gconv.Int(array[2]),
                        Db   : gconv.Int(array[3]),
                        Pass : array[4],
                    }
                    if v, ok := parse["maxIdle"]; ok {
                        config.MaxIdle = gconv.Int(v)
                    }
                    if v, ok := parse["maxActive"]; ok {
                        config.MaxActive = gconv.Int(v)
                    }
                    if v, ok := parse["idleTimeout"]; ok {
                        config.IdleTimeout = gconv.TimeDuration(v)*time.Second
                    }
                    if v, ok := parse["maxConnLifetime"]; ok {
                        config.MaxConnLifetime = gconv.TimeDuration(v)*time.Second
                    }
                    return gredis.New(config)
                }
                array, _ = gregex.MatchString(`(.+):(\d+),{0,1}(\d*),{0,1}(.*)`, line)
                if len(array) == 5 {
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
            glog.Errorfln(`incomplete configuration for redis: "redis" node not found in config file "%s"`, config.GetFilePath())
        }
        return nil
    })
    if result != nil {
        return result.(*gredis.Redis)
    }
    return nil
}

// 模板内置方法：config
func funcConfig(pattern string, file...string) string {
    return Config().GetString(pattern, file...)
}


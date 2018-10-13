// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 单例对象管理.
// 框架内置了一些核心对象获取方法，并且可以通过Set和Get方法实现IoC以及对内置核心对象的自定义替换
package gins

import (
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/os/gcmd"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "fmt"
)

const (
    gFRAME_CORE_COMPONENT_NAME_VIEW       = "gf.core.component.view"
    gFRAME_CORE_COMPONENT_NAME_CONFIG     = "gf.core.component.config"
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
    instances.Set(key, key)
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
func View() *gview.View {
    return instances.GetOrSetFuncLock(gFRAME_CORE_COMPONENT_NAME_VIEW, func() interface{} {
        path := gcmd.Option.Get("gf.viewpath")
        if path == "" {
            path = genv.Get("gf.viewpath")
            if path == "" {
                path = gfile.SelfDir()
            }
        }
        view := gview.Get(path)
        // 添加基于源码的搜索目录检索地址，常用于开发环境调试，只添加入口文件目录
        if p := gfile.MainPkgPath(); gfile.Exists(p) {
            view.AddPath(p)
        }
        // 框架内置函数
        view.BindFunc("config", funcConfig)
        return view
    }).(*gview.View)
}

// 核心对象：Config
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config() *gcfg.Config {
    return instances.GetOrSetFuncLock(gFRAME_CORE_COMPONENT_NAME_CONFIG, func() interface{} {
        path := gcmd.Option.Get("gf.cfgpath")
        if path == "" {
            path = genv.Get("gf.cfgpath")
            if path == "" {
                path = gfile.SelfDir()
            }
        }
        config := gcfg.New(path)
        // 添加基于源码的搜索目录检索地址，常用于开发环境调试，只添加入口文件目录
        if p := gfile.MainPkgPath(); gfile.Exists(p) {
            config.AddPath(p)
        }
        return config
    }).(*gcfg.Config)
}

// 数据库操作对象，使用了连接池
func Database(name...string) *gdb.Db {
    config := Config()
    db := instances.GetOrSetFuncLock(gFRAME_CORE_COMPONENT_NAME_DATABASE, func() interface{} {
        m := config.GetMap("database")
        if m == nil {
            panic(fmt.Sprintf(`incomplete configuration for database: "database" node not found in config file "%s"`, config.GetFilePath()))
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
                    if value, ok := nodem["max-idle"]; ok {
                        node.MaxIdleConnCount = gconv.Int(value)
                    }
                    if value, ok := nodem["max-open"]; ok {
                        node.MaxOpenConnCount = gconv.Int(value)
                    }
                    if value, ok := nodem["max-lifetime"]; ok {
                        node.MaxConnLifetime = gconv.Int(value)
                    }
                    cg = append(cg, node)
                }
            }
            gdb.AddConfigGroup(group, cg)
        }
        // 使用gfsnotify进行文件监控，当配置文件有任何变化时，清空数据库配置缓存
        gfsnotify.Add(config.GetFilePath(), func(event *gfsnotify.Event) {
            instances.Remove(gFRAME_CORE_COMPONENT_NAME_DATABASE)
        })
        if db, err := gdb.New(name...); err == nil {
            return db
        } else {
            panic(err)
        }
        return nil
    })
    if db != nil {
        return db.(*gdb.Db)
    }
    return nil
}

// 模板内置方法：config
func funcConfig(pattern string, file...string) gview.HTML {
    return gview.HTML(Config().GetString(pattern, file...))
}


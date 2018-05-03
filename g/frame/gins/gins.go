// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 单例对象管理.
// 框架内置了一些核心对象获取方法，并且可以通过Set和Get方法实现IoC以及对内置核心对象的自定义替换
package gins

import (
    "runtime"
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/os/gcmd"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gregx"
    "gitee.com/johng/gf/g/container/gmap"
)

const (
    gFRAME_CORE_COMPONENT_NAME_VIEW       = "gf.core.component.view"
    gFRAME_CORE_COMPONENT_NAME_CONFIG     = "gf.core.component.config"
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
        // 添加代码级的搜索目录检索地址，常用于开发环境调试，只添加入口文件目录
        path = ""
        for i := 1; i < 10000; i++ {
            if _, file, _, ok := runtime.Caller(i); ok {
                // 不包含go源码路径
                if !gregx.IsMatchString("^" + runtime.GOROOT(), file) {
                    path = file
                }
            } else {
                break
            }
        }
        if path != "" {
            config.AddPath(gfile.Dir(path))
        }
        // 单例对象缓存控制
        Set(gFRAME_CORE_COMPONENT_NAME_CONFIG, config)
        return config
    }
    return nil
}

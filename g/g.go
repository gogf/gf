// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 常用数据类型以及对象封装

package g

import (
    "strings"
    "strconv"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/database/gredis"
)

// 常用map数据结构
type Map map[string]interface{}

// 常用list数据结构
type List []Map

// 核心对象：View
func View() *gview.View {
    return gins.View()
}

// Config配置管理对象
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config() *gcfg.Config {
    return gins.Config()
}


// 数据库操作对象，使用了连接池
func Database(name...string) *gdb.Db {
    config := gins.Config()
    if config == nil {
        return nil
    }
    if m := config.GetMap("database"); m != nil {
        c := gdb.Config{}
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
                        node.Priority, _ = strconv.Atoi(gconv.String(value))
                    }
                    cg = append(cg, node)
                }
            }
            c[group] = cg
        }
        gdb.SetConfig(c)

        if db, err := gdb.New(name...); err == nil {
            return db
        } else {
            return nil
        }
    }

    return nil
}

// Redis操作对象，使用了连接池
func Redis(name...string) *gredis.Redis {
    group := "default"
    if len(name) > 0 {
        group = name[0]
    }
    config := gins.Config()
    if config == nil {
        return nil
    }
    if m := config.GetMap("redis"); m != nil {
        if v, ok := m[group]; ok {
            array := strings.Split(gconv.String(v), ",")
            if len(array) > 1 {
                return gredis.New(array[0], array[1])
            }
        }
    }
    return nil
}
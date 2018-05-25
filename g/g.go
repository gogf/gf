// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 常用数据类型以及对象封装

package g

import (
    "strings"
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/database/gredis"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/net/gudp"
)

const (
    gIS_DATABASE_CONFIG_CACHED = "gf.core.component.database.cached"
)

// 常用map数据结构(使用别名)
type Map  = map[string]interface{}

// 常用list数据结构(使用别名)
type List = []Map


// 阻塞等待HTTPServer执行完成(同一进程多HTTPServer情况下)
func Wait() {
    ghttp.Wait()
}

// HTTPServer单例对象
func Server(name...interface{}) *ghttp.Server {
    return ghttp.GetServer(name...)
}

// TCPServer单例对象
func TcpServer(name...interface{}) *gtcp.Server {
    return gtcp.GetServer(name...)
}

// UDPServer单例对象
func UdpServer(name...interface{}) *gudp.Server {
    return gudp.GetServer(name...)
}

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
    // 数据库配置是否已经设置
    if gcache.Get(gIS_DATABASE_CONFIG_CACHED) == nil {
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
                            node.Priority = gconv.Int(value)
                        }
                        cg = append(cg, node)
                    }
                }
                c[group] = cg
            }
            gdb.SetConfig(c)
            gcache.Set(gIS_DATABASE_CONFIG_CACHED, struct{}{}, 0)
            // 使用gfsnotify进行文件监控，当配置文件有任何变化时，清空数据库配置缓存
            gfsnotify.Add(Config().GetFilePath(), func(event *gfsnotify.Event) {
                gcache.Remove(gIS_DATABASE_CONFIG_CACHED)
            })
        }
    }
    if db, err := gdb.New(name...); err == nil {
        return db
    } else {
        return nil
    }
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



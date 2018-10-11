// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package g

import (
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/os/gfsnotify"
    "gitee.com/johng/gf/g/database/gredis"
    "gitee.com/johng/gf/g/frame/gins"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/net/gudp"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gcfg"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/os/glog"
    "fmt"
)
const (
    gIS_DATABASE_CONFIG_CACHED = "gf.core.component.database.cached"
)


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
        glog.Error("Config component init failed")
        return nil
    }
    // 数据库配置是否已经设置到gdb模块中
    gcache.GetOrSetFuncLock(gIS_DATABASE_CONFIG_CACHED, func() interface{} {
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
                c[group] = cg
            }
            gdb.SetConfig(c)
            // 使用gfsnotify进行文件监控，当配置文件有任何变化时，清空数据库配置缓存
            gfsnotify.Add(Config().GetFilePath(), func(event *gfsnotify.Event) {
                gcache.Remove(gIS_DATABASE_CONFIG_CACHED)
            })
            return struct{}{}
        } else {
            glog.Error(fmt.Sprintf(`incomplete configuration for database: "database" node not found in config file "%s"`, config.GetFilePath()))
        }
        return nil
    }, 0)
    if db, err := gdb.New(name...); err == nil {
        return db
    } else {
        glog.Error(err)
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
        // host:port[,db[,pass]]
        if v, ok := m[group]; ok {
            array, err := gregex.MatchString(`(.+):(\d+),{0,1}(\d*),{0,1}(.*)`, gconv.String(v))
            if err == nil {
                return gredis.New(gredis.Config{
                    Host : array[1],
                    Port : gconv.Int(array[2]),
                    Db : gconv.Int(array[3]),
                    Pass : array[4],
                })
            }
        }
    }
    return nil
}
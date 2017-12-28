// 单例对象管理工具
// 框架内置了一些核心对象，并且可以通过Set和Get方法实现IoC以及对内置核心对象的自定义替换
package ginstance

import (
    "strconv"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/genv"
    "gitee.com/johng/gf/g/os/gview"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gconsole"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/frame/gconfig"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/net/ghttp"
)

const (
    FRAME_CORE_COMPONENT_NAME_VIEW     = "gf.component.view"
    FRAME_CORE_COMPONENT_NAME_CONFIG   = "gf.component.config"
    FRAME_CORE_COMPONENT_NAME_DATABASE = "gf.component.database"
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

// 核心对象：Server
// 框架支持多服务器对象，通过传入不同的name进行区分
func Server(names...string) *ghttp.Server {
    name := "default"
    if len(names) > 0 {
        name = names[0]
    }
    return ghttp.GetServer(name)
}

// 核心对象：View
func View() *gview.View {
    result := Get(FRAME_CORE_COMPONENT_NAME_VIEW)
    if result != nil {
        return result.(*gview.View)
    } else {
        path := gconsole.Option.Get("viewpath")
        if path == "" {
            path = genv.Get("viewpath")
            if path == "" {
                path = gfile.SelfDir()
            }
        }
        view := gview.GetView(path)
        Set(FRAME_CORE_COMPONENT_NAME_VIEW, view)
        return view
    }
    return nil
}

// 核心对象：Config
// 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
func Config() *gconfig.Config {
    result := Get(FRAME_CORE_COMPONENT_NAME_CONFIG)
    if result != nil {
        return result.(*gconfig.Config)
    } else {
        path := gconsole.Option.Get("cfgpath")
        if path == "" {
            path = genv.Get("cfgpath")
            if path == "" {
                path = gfile.SelfDir()
            }
        }
        config := gconfig.New(path)
        Set(FRAME_CORE_COMPONENT_NAME_CONFIG, config)
        return config
    }
    return nil
}

// 核心对象：Database
func Database(names...string) gdb.Link {
    result := Get(FRAME_CORE_COMPONENT_NAME_DATABASE)
    if result != nil {
        return result.(gdb.Link)
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
            var db gdb.Link = nil
            if len(names) == 0 {
                if link, err := gdb.Instance(); err == nil {
                    db = link
                } else {
                    glog.Error(err)
                }
            } else {
                if link, err := gdb.InstanceByGroup(names[0]); err == nil {
                    db = link
                } else {
                    glog.Error(err)
                }
            }
            if db != nil {
                Set(FRAME_CORE_COMPONENT_NAME_DATABASE, db)
                return db
            }
        }
    }
    return nil
}
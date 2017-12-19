// 基类
package gbase

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gconsole"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/frame/gconfig"
    "gitee.com/johng/gf/g/frame/ginstance"
)

// 框架基类，所有的基于gf框架的类对象都继承于此，以便使用框架的一些封装的核心组件
type Base struct {
    Db       gdb.Link        // 数据库操作对象
    Config   *gconfig.Config // 配置管理对象
}

// 基类初始化，如若需要自定义初始化内置核心对象组件，可在继承子类中覆盖此方法
func (b *Base) Init() {
    // 配置文件目录查找依次为：启动参数cfgpath、当前程序运行目录
    if b.Config == nil {
        path := gconsole.Option.Get("cfgpath")
        if path == "" {
            path = gfile.SelfDir()
        }
        ckey   := "gf_config_with_path_" + path
        result := ginstance.Get(ckey)
        if result != nil {
            b.Config = result.(*gconfig.Config)
        } else {
            b.Config = gconfig.New(path)
            ginstance.Set(ckey, b.Config)
        }
    }
    // 数据库操作对象初始化
    // 全局只有一个数据库单例对象，可以配置不同分组的配置进行使用
    if b.Db == nil {
        ckey   := "gf_database"
        result := ginstance.Get(ckey)
        if result != nil {
            b.Db = result.(gdb.Link)
        } else {
            if m := b.Config.GetMap("database"); m != nil {
                for group, v := range m {
                    if list, ok := v.([]interface{}); ok {
                        for _, nodei := range list {
                            gdb.AddConfigNodeByString(group, nodei.(string))
                        }
                    }
                }
                if link, err := gdb.Instance(); err != nil {
                    b.Db = link
                    ginstance.Set(ckey, b.Db)
                }
            }
        }
    }
}
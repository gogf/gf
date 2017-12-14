// 基类
package gbase

import (
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/frame/gconfig"
    "gitee.com/johng/gf/g/frame/ginstance"
)

const (
    gDEFAULT_CONFIG_FILE = "config.json" // 默认读取的配置文件名称
)

// 框架基类
type Base struct {
    Db     gdb.Link
    Config *gconfig.Config
}

// 基类初始化，如若需要自定义初始化内置核心对象组件，可在继承子类中覆盖此方法
func (b *Base) Init() {
    // 默认配置目录为当前程序运行目录
    if b.Config == nil {
        path   := gfile.SelfDir()
        ckey   := "gf_config_with_path_" + path
        result := ginstance.Get(ckey)
        if result != nil {
            b.Config = result.(*gconfig.Config)
        } else {
            b.Config = gconfig.New(path)
            b.Config.Add(gDEFAULT_CONFIG_FILE)
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
            }
            if link, err := gdb.Instance(); err != nil {
                b.Db = link
                ginstance.Set(ckey, b.Db)
            }
        }
    }
}
package gdb

import (
    "database/sql"
    "g"
    "errors"
    "g/core/gfunc"
    "fmt"
    "log"
)

const (
    OPTION_INSERT  = 0
    OPTION_REPLACE = 1
    OPTION_SAVE    = 2
    OPTION_IGNORE  = 3
)

const (
    gMAX_LOADBALANCE_CHECK_SIZE = 1000000
)

// 数据库配置包内对象
var config struct {
    c  Config
    lb gLoadBalance
}

// 数据库事务操作对象
type gTrasaction struct {
    db *sql.DB
    tx *sql.Tx
}

// 数据库链接对象
type gLink struct {
    Transaction gTrasaction
    master *sql.DB
    slave  *sql.DB
}

// 应用层的数据库负载均衡
type gLoadBalance struct {
    master []int
    slave  []int
}

// 数据库配置
type Config      map[string]ConfigGroup

// 数据库集群配置
type ConfigGroup []ConfigItem

// 数据库单项配置
type ConfigItem  struct {
    Host     string // 地址
    Port     int    // 端口
    User     string // 账号
    Pass     string // 密码
    Name     string // 数据库名称
    Type     string // 数据库类型：mysql, sqlite, mssql, postgresql, oracle(目前仅支持mysql)
    Role     string // 数据库的角色，用于主从操作分离，至少需要有一个master，参数值：master, slave
    Charset  string // 编码，默认为 utf-8
}

// 记录关联数组
type DataMap  map[string]string

// 记录关联数组列表(索引从0开始的数组)
type DataList []DataMap

// 数据库集群配置示例，支持主从处理，多数据库集群支持
/*
var Database = Config {
    // 数据库集群配置名称
    "default" : ConfigGroup {
        {
            Host    : "127.0.0.1",
            Port    : 3306,
            User    : "root",
            Pass    : "123456",
            Name    : "test",
            Type    : "mysql",
            Role    : "master",
            Charset : "utf-8",
        },
    },
}
*/

// 设置数据库配置信息
func SetConfig (c Config) {
    config.c = c
}

// 使用默认选项进行连接
func New() (*gLink, error) {
    return NewByConfig("default")
}

// 根据数据库配置项创建一个数据库操作对象
// @todo 需要检测当master和slave都是使用同一个配置时进行链接，是否能够io复用
// @todo 增加负载均衡的处理
// @todo 增加同一配置的单例实现
func NewByConfig(dbConfigItemKeyName string) (*gLink, error) {
    if len(config.c) < 1 {
        return nil, errors.New("empty database configuration")
    }
    if list, ok := config.c[dbConfigItemKeyName]; ok {
        masterList  := make([]ConfigItem, 0)
        slaveList   := make([]ConfigItem, 0)
        for i := 0; i < len(list); i++ {
            if list[i].Role == "master" {
                masterList = append(masterList, list[i])
            }
            if list[i].Role == "slave" {
                slaveList = append(slaveList, list[i])
            }
        }
        if len(masterList) < 1 {
            return nil, errors.New("empty master node of database configuration")
        }
        link       := gLink{}
        masterItem := masterList[gfunc.Rand(0, len(masterList))]
        slaveItem  := masterItem
        if len(slaveList) > 0 {
            slaveItem = slaveList[gfunc.Rand(0, len(slaveList))]
        }
        master, err := sql.Open(masterItem.Type, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", masterItem.User, masterItem.Pass, masterItem.Host, masterItem.Port, masterItem.Name))
        if err != nil {
            log.Fatal(err)
        }
        slave,  err := sql.Open(slaveItem.Type,  fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", slaveItem.User, slaveItem.Pass, slaveItem.Host, slaveItem.Port, slaveItem.Name))
        if err != nil {
            log.Fatal(err)
        }
        link.master         = master
        link.slave          = slave
        link.Transaction.db = master
        return &link, nil
    } else {
        return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", dbConfigItemKeyName))
    }

}


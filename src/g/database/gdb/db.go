package gdb

import (
    "database/sql"
    "errors"
    "g/core/gfunc"
    "fmt"
    "log"
    "g/core/ginstance"
)

const (
    OPTION_INSERT  = 0
    OPTION_REPLACE = 1
    OPTION_SAVE    = 2
    OPTION_IGNORE  = 3
)

const (
    // 用于负载均衡判断的最大的记录链接次数，超过该次数则将每个数据库资源链接的记录次数重置为0
    gMAX_LOADBALANCE_CHECK_SIZE = 1000000
)

// 数据库配置包内对象
var config struct {
    c  Config
    d  string
    lb gLoadBalance
}

// 数据库事务操作对象
type gTrasaction struct {
    db *sql.DB
    tx *sql.Tx
}

// 数据库链接对象
type Link struct {
    Transaction gTrasaction
    master *sql.DB
    slave  *sql.DB
}

// 应用层的数据库负载均衡(平均负载)
type gLoadBalance struct {
    master []int
    slave  []int
}

// 数据库配置
type Config      map[string]ConfigGroup

// 数据库集群配置
type ConfigGroup []ConfigNode

// 数据库单项配置
type ConfigNode  struct {
    Host     string // 地址
    Port     string // 端口
    User     string // 账号
    Pass     string // 密码
    Name     string // 数据库名称
    Type     string // 数据库类型：mysql, sqlite, mssql, postgresql, oracle(目前仅支持mysql)
    Role     string // (可选)数据库的角色，用于主从操作分离，至少需要有一个master，参数值：master, slave
    Charset  string // (可选)编码，默认为 utf-8
}

// 记录关联数组
type DataMap  map[string]string

// 记录关联数组列表(索引从0开始的数组)
type DataList []DataMap

// 数据库集群配置示例，支持主从处理，多数据库集群支持
/*
var DatabaseConfiguration = Config {
    // 数据库集群配置名称
    "default" : ConfigGroup {
        {
            Host    : "127.0.0.1",
            Port    : "3306",
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

// 包初始化
func init() {
    config.d = "default"
}

// 设置当前应用的数据库配置信息
// 支持三种数据类型的输入参数：Config, ConfigGroup, ConfigNode
func SetConfig (c interface{}) {
    switch c.(type) {
        case Config:
            config.c = c.(Config)

        case ConfigGroup:
            config.c = Config {"default" : c.(ConfigGroup)}

        case ConfigNode:
            config.c = Config {"default" : ConfigGroup { c.(ConfigNode) }}

        default:
            panic("invalid config type, valid types are: Config, ConfigGroup, ConfigNode")
    }
}

// 设置默认链接的数据库链接配置项(默认是 default)
func SetDefaultGroup (groupName string) {
    config.d = groupName
}

// 根据配置项获取一个数据库操作对象单例
func instance (groupName string) (*Link, error) {
    instanceName := "gdb_instance_" + groupName
    result       := ginstance.Get(instanceName)
    if result == nil {
        link, err := NewByGroup(groupName)
        if err == nil {
            ginstance.Set(instanceName, link)
            return link, nil
        } else {
            return nil, err
        }
    } else {
        return result.(*Link), nil
    }
}

// 获得默认的数据库操作对象单例
func Instance () (*Link, error) {
    return instance(config.d)
}

// 获得指定配置项的数据库草最对象单例
func InstanceByGroup(groupName string) (*Link, error) {
    return instance(groupName)
}

// 使用默认选项进行连接，数据库集群配置项：default
func New() (*Link, error) {
    return NewByGroup(config.d)
}

// 根据数据库配置项创建一个数据库操作对象
// @todo 增加负载均衡的处理
func NewByGroup(groupName string) (*Link, error) {
    if len(config.c) < 1 {
        return nil, errors.New("empty database configuration")
    }
    if list, ok := config.c[groupName]; ok {
        // 将master, slave集群列表拆分出来
        masterList  := make(ConfigGroup, 0)
        slaveList   := make(ConfigGroup, 0)
        for i := 0; i < len(list); i++ {
            if list[i].Role == "slave" {
                slaveList  = append(slaveList, list[i])
            } else {
                // 默认配置项的角色为master
                masterList = append(masterList, list[i])
            }
        }
        if len(masterList) < 1 {
            return nil, errors.New("at least one master node configuration's need to make sense")
        }
        link        := Link{}
        // master链接对象
        masterItem  := masterList[gfunc.Rand(0, len(masterList))]
        master, err := openSql(&masterItem)
        if err != nil {
            log.Fatal(err)
        }
        // slave链接对象
        // 如果整个配置中仅有一个master配置项，那么slave和master共用一个链接对象
        slave := master
        if len(slaveList) > 0 {
            slaveItem  := slaveList[gfunc.Rand(0, len(slaveList))]
            slave,  err = openSql(&slaveItem)
            if err != nil {
                log.Fatal(err)
            }
        }
        link.master         = master
        link.slave          = slave
        link.Transaction.db = master
        return &link, nil
    } else {
        return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", groupName))
    }
}

//// 按照负载均衡算法从数据库集群中选择一个配置节点出来使用
//func getConfigNodeFromGroup (cg *ConfigGroup) *ConfigNode {
//
//}

// 创建SQL操作对象，内部采用了lazy link处理
func openSql (c *ConfigNode) (*sql.DB, error) {
    db,  err := sql.Open(c.Type,  fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pass, c.Host, c.Port, c.Name))
    if err != nil {
        log.Fatal(err)
    }
    return db, err
}


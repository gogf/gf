// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//
// 对常用关系数据库的封装管理包

package gdb

import (
    "fmt"
    "errors"
    "database/sql"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/grand"
    _ "github.com/lib/pq"
    _ "github.com/go-sql-driver/mysql"
)

const (
    OPTION_INSERT  = 0
    OPTION_REPLACE = 1
    OPTION_SAVE    = 2
    OPTION_IGNORE  = 3
)

// 数据库操作接口
type Link interface {
    Open (c *ConfigNode) (*sql.DB, error)
    Close() error
    Query(q string, args ...interface{}) (*sql.Rows, error)
    Exec(q string, args ...interface{}) (sql.Result, error)
    Prepare(q string) (*sql.Stmt, error)

    GetAll(q string, args ...interface{}) (List, error)
    GetOne(q string, args ...interface{}) (Map, error)
    GetValue(q string, args ...interface{}) (interface{}, error)

    PingMaster() error
    PingSlave() error

    SetMaxIdleConns(n int)
    SetMaxOpenConns(n int)

    setMaster(master *sql.DB)
    setSlave(slave *sql.DB)
    setQuoteChar(left string, right string)
    setLink(link Link)
    getQuoteCharLeft () string
    getQuoteCharRight () string
    handleSqlBeforeExec(q *string) *string

    Begin() (*sql.Tx, error)
    Commit() error
    Rollback() error

    insert(table string, data Map, option uint8) (sql.Result, error)
    Insert(table string, data Map) (sql.Result, error)
    Replace(table string, data Map) (sql.Result, error)
    Save(table string, data Map) (sql.Result, error)

    batchInsert(table string, list List, batch int, option uint8) error
    BatchInsert(table string, list List, batch int) error
    BatchReplace(table string, list List, batch int) error
    BatchSave(table string, list List, batch int) error

    Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error)
    Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error)

    Table(tables string) (*gLinkOp)
}

// 数据库链接对象
type dbLink struct {
    link        Link
    transaction *sql.Tx
    master      *sql.DB
    slave       *sql.DB
    charl        string
    charr        string
}

// 关联数组，绑定一条数据表记录
type Map  map[string]interface{}

// 关联数组列表(索引从0开始的数组)，绑定多条记录
type List []Map

// 获得默认的数据库操作对象单例
func Instance () (Link, error) {
    return instance(config.d)
}

// 获得指定配置项的数据库草最对象单例
func InstanceByGroup(groupName string) (Link, error) {
    return instance(groupName)
}

// 根据配置项获取一个数据库操作对象单例
func instance (groupName string) (Link, error) {
    instanceName := "gdb_instance_" + groupName
    result       := gcache.Get(instanceName)
    if result == nil {
        link, err := NewByGroup(groupName)
        if err == nil {
            gcache.Set(instanceName, link, 0)
            return link, nil
        } else {
            return nil, err
        }
    } else {
        return result.(Link), nil
    }
}

// 使用默认选项进行连接，数据库集群配置项：default
func New() (Link, error) {
    return NewByGroup(config.d)
}

// 根据数据库配置项创建一个数据库操作对象
func NewByGroup(groupName string) (Link, error) {
    config.RLock()
    defer config.RUnlock()

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
        masterNode := getConfigNodeByPriority(&masterList)
        var slaveNode *ConfigNode
        if len(slaveList) > 0 {
            slaveNode = getConfigNodeByPriority(&slaveList)
        }
        return newLink(masterNode, slaveNode)
    } else {
        return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", groupName))
    }
}

// 根据单点数据库配置获得一个数据库草最对象
func NewByConfigNode(node ConfigNode) (Link, error) {
    return newLink (&node, nil)
}

// 按照负载均衡算法(优先级配置)从数据库集群中选择一个配置节点出来使用
func getConfigNodeByPriority (cg *ConfigGroup) *ConfigNode {
    if len(*cg) < 2 {
        return &(*cg)[0]
    }
    var total int
    for i := 0; i < len(*cg); i++ {
        total += (*cg)[i].Priority * 100
    }
    r   := grand.Rand(0, total)
    min := 0
    max := 0
    for i := 0; i < len(*cg); i++ {
        max = min + (*cg)[i].Priority * 100
        //fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
        if r >= min && r < max {
            return &(*cg)[i]
        } else {
            min = max
        }
    }
    return nil
}

// 创建数据库链接对象
func newLink (masterNode *ConfigNode, slaveNode *ConfigNode) (Link, error) {
    var link Link
    switch masterNode.Type {
        case "mysql":
            link = Link(&mysqlLink{})

        case "pgsql":
            link = Link(&pgsqlLink{})

        default:
            return nil, errors.New(fmt.Sprintf("unsupported db type '%s'", masterNode.Type))
    }
    master, err := link.Open(masterNode)
    if err != nil {
        glog.Fatal(err)
    }
    slave := master
    if slaveNode != nil {
        slave,  err = link.Open(slaveNode)
        if err != nil {
            glog.Fatal(err)
        }
    }
    link.setLink(link)
    link.setMaster(master)
    link.setSlave(slave)
    link.setQuoteChar(link.getQuoteCharLeft(), link.getQuoteCharRight())
    return link, nil
}

// 设置master链接对象
func (l *dbLink) setMaster(master *sql.DB) {
    l.master = master
}

// 设置slave链接对象
func (l *dbLink) setSlave(slave *sql.DB) {
    l.slave = slave
}

// 设置当前数据库类型引用字符
func (l *dbLink) setQuoteChar(left string, right string) {
    l.charl = left
    l.charr = right
}

// 设置挡脸操作的link接口
func (l *dbLink) setLink(link Link) {
    l.link = link
}


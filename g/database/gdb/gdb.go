// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 数据库ORM.
package gdb

import (
    "fmt"
    "errors"
    "database/sql"
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
    // 打开数据库连接，建立数据库操作对象
    Open (c *ConfigNode) (*sql.DB, error)

    // SQL操作方法
    Query(q string, args ...interface{}) (*sql.Rows, error)
    Exec(q string, args ...interface{}) (sql.Result, error)
    Prepare(q string) (*sql.Stmt, error)

    // 数据库查询
    GetAll(q string, args ...interface{}) (List, error)
    GetOne(q string, args ...interface{}) (Map, error)
    GetValue(q string, args ...interface{}) (interface{}, error)

    // Ping
    PingMaster() error
    PingSlave() error

    // 连接属性设置
    SetMaxIdleConns(n int)
    SetMaxOpenConns(n int)

    // 开启事务操作
    Begin() (*Tx, error)

    // 数据表插入/更新/保存操作
    Insert(table string, data Map) (sql.Result, error)
    Replace(table string, data Map) (sql.Result, error)
    Save(table string, data Map) (sql.Result, error)

    // 数据表插入/更新/保存操作(批量)
    BatchInsert(table string, list List, batch int) (sql.Result, error)
    BatchReplace(table string, list List, batch int) (sql.Result, error)
    BatchSave(table string, list List, batch int) (sql.Result, error)

    // 数据修改/删除
    Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error)
    Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error)

    // 创建链式操作对象(Table为From的别名)
    Table(tables string) (*DbOp)
    From(tables string)  (*DbOp)

    // 关闭数据库操作对象
    Close() error

    // 内部方法
    insert(table string, data Map, option uint8) (sql.Result, error)
    batchInsert(table string, list List, batch int, option uint8) (sql.Result, error)

    getQuoteCharLeft () string
    getQuoteCharRight () string
    handleSqlBeforeExec(q *string) *string
}

// 数据库链接对象
type Db struct {
    link   Link
    master *sql.DB
    slave  *sql.DB
    charl  string
    charr  string
}

// 关联数组，绑定一条数据表记录
type Map  map[string]interface{}

// 关联数组列表(索引从0开始的数组)，绑定多条记录
type List []Map

// 获得默认/指定分组名称的数据库操作对象单例
func Instance (groupName...string) (*Db, error) {
    name := config.d
    if len(groupName) > 0 {
        name = groupName[0]
    }
    return instance(name)
}

// 根据配置项获取一个数据库操作对象单例
func instance (groupName string) (*Db, error) {
    instanceName := "gdb_instance_" + groupName
    result       := gcache.Get(instanceName)
    if result == nil {
        db, err := New(groupName)
        if err == nil {
            gcache.Set(instanceName, db, 0)
            return db, nil
        } else {
            return nil, err
        }
    } else {
        return result.(*Db), nil
    }
}

// 使用默认/指定分组配置进行连接，数据库集群配置项：default
func New(groupName...string) (*Db, error) {
    name := config.d
    if len(groupName) > 0 {
        name = groupName[0]
    }
    config.RLock()
    defer config.RUnlock()

    if len(config.c) < 1 {
        return nil, errors.New("empty database configuration")
    }
    if list, ok := config.c[name]; ok {
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
        return newDb(masterNode, slaveNode)
    } else {
        return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", name))
    }
}

// 根据单点数据库配置获得一个数据库草最对象
func NewByNode(node ConfigNode) (*Db, error) {
    return newDb (&node, nil)
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
func newDb (masterNode *ConfigNode, slaveNode *ConfigNode) (*Db, error) {
    var link Link
    switch masterNode.Type {
        case "mysql":
            link = Link(&dbmysql{})

        case "pgsql":
            link = Link(&dbpgsql{})

        default:
            return nil, errors.New(fmt.Sprintf("unsupported db type '%s'", masterNode.Type))
    }
    master, err := link.Open(masterNode)
    if err != nil {
        return nil, err
    }
    slave := master
    if slaveNode != nil {
        slave,  err = link.Open(slaveNode)
        if err != nil {
            return nil, err
        }
    }
    //link.setLink(link)
    //link.setMaster(master)
    //link.setSlave(slave)
    //link.setQuoteChar(link.getQuoteCharLeft(), link.getQuoteCharRight())
    return &Db {
        link   : link,
        master : master,
        slave  : slave,
        charl  : link.getQuoteCharLeft(),
        charr  : link.getQuoteCharRight(),
    }, nil
}


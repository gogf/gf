// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// 对常用关系数据库的封装管理包

package gdb

import (
    "fmt"
    "gitee.com/johng/gf/g/container/gring"
    "sync"
)

const (
    DEFAULT_GROUP_NAME = "default" // 默认配置名称
)

// 数据库配置包内对象
var config struct {
    sync.RWMutex
    c  Config // 数据库配置
    d  string // 默认数据库分组名称
}

// 数据库配置
type Config      map[string]ConfigGroup

// 数据库集群配置
type ConfigGroup []ConfigNode

// 数据库单项配置
type ConfigNode  struct {
    Host             string   // 地址
    Port             string   // 端口
    User             string   // 账号
    Pass             string   // 密码
    Name             string   // 数据库名称
    Type             string   // 数据库类型：mysql, sqlite, mssql, pgsql, oracle(目前仅支持mysql)
    Role             string   // (可选，默认为master)数据库的角色，用于主从操作分离，至少需要有一个master，参数值：master, slave
    Charset          string   // (可选，默认为 utf8)编码，默认为 utf8
    Priority         int      // (可选)用于负载均衡的权重计算，当集群中只有一个节点时，权重没有任何意义
    Linkinfo         string   // (可选)自定义链接信息，当该字段被设置值时，以上链接字段(Host,Port,User,Pass,Name)将失效(该字段是一个扩展功能)
    MaxIdleConnCount int      // (可选)连接池最大限制的连接数
    MaxOpenConnCount int      // (可选)连接池最大打开的连接数
    MaxConnLifetime  int      // (可选，单位秒)连接对象可重复使用的时间长度
}

// 数据库集群配置示例，支持主从处理，多数据库集群支持
/*
var DatabaseConfiguration = Config {
    // 数据库集群配置名称
    "default" : ConfigGroup {
        {
            Host     : "192.168.1.100",
            Port     : "3306",
            User     : "root",
            Pass     : "123456",
            Name     : "test",
            Type     : "mysql",
            Role     : "master",
            Charset  : "utf8",
            Priority : 100,
        },
        {
            Host     : "192.168.1.101",
            Port     : "3306",
            User     : "root",
            Pass     : "123456",
            Name     : "test",
            Type     : "mysql",
            Role     : "slave",
            Charset  : "utf8",
            Priority : 100,
        },
    },
}
*/

// 包初始化
func init() {
    config.c = make(Config)
    config.d = DEFAULT_GROUP_NAME
}

// 设置当前应用的数据库配置信息，进行全局数据库配置覆盖操作
func SetConfig (c Config) {
    config.Lock()
    defer config.Unlock()
    config.c = c
}

// 添加数据库服务器集群配置
func AddConfigGroup (group string, nodes ConfigGroup) {
    config.Lock()
    config.c[group] = nodes
    config.Unlock()
}

// 添加一台数据库服务器配置
func AddConfigNode (group string, node ConfigNode) {
    config.Lock()
    config.c[group] = append(config.c[group], node)
    config.Unlock()
}

// 添加默认链接的一台数据库服务器配置
func AddDefaultConfigNode (node ConfigNode) {
    AddConfigNode(DEFAULT_GROUP_NAME, node)
}

// 添加默认链接的数据库服务器集群配置
func AddDefaultConfigGroup (nodes ConfigGroup) {
    AddConfigGroup(DEFAULT_GROUP_NAME, nodes)
}

// 添加一台数据库服务器配置
func GetConfig (group string) ConfigGroup {
    config.RLock()
    defer config.RUnlock()
    return config.c[group]
}

// 设置默认链接的数据库链接配置项(默认是 default)
func SetDefaultGroup (groupName string) {
    config.Lock()
    config.d = groupName
    config.Unlock()
}

// 设置数据库连接池中空闲链接的大小
func (bs *dbBase) SetMaxIdleConns(n int) {
    bs.maxIdleConnCount.Set(n)
}

// 设置数据库连接池最大打开的链接数量
func (bs *dbBase) SetMaxOpenConns(n int) {
    bs.maxOpenConnCount.Set(n)
}

// 设置数据库连接可重复利用的时间，超过该时间则被关闭废弃
// 如果 d <= 0 表示该链接会一直重复利用
func (bs *dbBase) SetConnMaxLifetime(n int) {
    bs.maxConnLifetime.Set(n)
}

// 节点配置转换为字符串
func (node *ConfigNode) String() string {
    if node.Linkinfo != "" {
        return node.Linkinfo
    }
    return fmt.Sprintf(`%s@%s:%s,%s,%s,%s,%s,%d-%d-%d`, node.User, node.Host, node.Port,
        node.Name, node.Type, node.Role, node.Charset,
        node.MaxIdleConnCount, node.MaxOpenConnCount, node.MaxConnLifetime,
    )
}

// 是否开启调试服务
func (bs *dbBase) SetDebug(debug bool) {
    bs.debug.Set(debug)
    if debug && bs.sqls == nil {
        bs.sqls = gring.New(gDEFAULT_DEBUG_SQL_LENGTH)
    }
}

// 获取是否开启调试服务
func (bs *dbBase) getDebug() bool {
    return bs.debug.Val()
}
// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"sync"

	"github.com/gogf/gf/g/container/gring"
)

const (
	DEFAULT_GROUP_NAME = "default" // 默认配置名称
)

// 数据库分组配置
type Config map[string]ConfigGroup

// 数据库集群配置
type ConfigGroup []ConfigNode

// 数据库单项配置
type ConfigNode struct {
	Host             string // 地址
	Port             string // 端口
	User             string // 账号
	Pass             string // 密码
	Name             string // 数据库名称
	Type             string // 数据库类型：mysql, sqlite, mssql, pgsql, oracle
	Role             string // (可选，默认为master)数据库的角色，用于主从操作分离，至少需要有一个master，参数值：master, slave
	Debug            bool   // (可选)开启调试模式
	Weight           int    // (可选)用于负载均衡的权重计算，当集群中只有一个节点时，权重没有任何意义
	Charset          string // (可选，默认为 utf8)编码，默认为 utf8
	LinkInfo         string // (可选)自定义链接信息，当该字段被设置值时，以上链接字段(Host,Port,User,Pass,Name)将失效(该字段是一个扩展功能)
	MaxIdleConnCount int    // (可选)连接池最大限制的连接数
	MaxOpenConnCount int    // (可选)连接池最大打开的连接数
	MaxConnLifetime  int    // (可选，单位秒)连接对象可重复使用的时间长度
}

// 数据库配置包内对象
var configs struct {
	sync.RWMutex        // 并发安全互斥锁
	config       Config // 数据库分组配置
	defaultGroup string // 默认数据库分组名称
}

// 包初始化
func init() {
	configs.config = make(Config)
	configs.defaultGroup = DEFAULT_GROUP_NAME
}

// 设置当前应用的数据库配置信息，进行全局数据库配置覆盖操作
func SetConfig(config Config) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config = config
}

// 添加数据库服务器集群配置
func AddConfigGroup(group string, nodes ConfigGroup) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = nodes
}

// 添加一台数据库服务器配置
func AddConfigNode(group string, node ConfigNode) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.config[group] = append(configs.config[group], node)
}

// 添加默认链接的一台数据库服务器配置
func AddDefaultConfigNode(node ConfigNode) {
	AddConfigNode(DEFAULT_GROUP_NAME, node)
}

// 添加默认链接的数据库服务器集群配置
func AddDefaultConfigGroup(nodes ConfigGroup) {
	AddConfigGroup(DEFAULT_GROUP_NAME, nodes)
}

// 添加一台数据库服务器配置
func GetConfig(group string) ConfigGroup {
	configs.RLock()
	defer configs.RUnlock()
	return configs.config[group]
}

// 设置默认链接的数据库链接配置项(默认是 default)
func SetDefaultGroup(name string) {
	defer instances.Clear()
	configs.Lock()
	defer configs.Unlock()
	configs.defaultGroup = name
}

// 获取默认链接的数据库链接配置项(默认是 default)
func GetDefaultGroup() string {
	defer instances.Clear()
	configs.RLock()
	defer configs.RUnlock()
	return configs.defaultGroup
}

// 设置数据库连接池中空闲链接的大小
func (bs *dbBase) SetMaxIdleConnCount(n int) {
	bs.maxIdleConnCount.Set(n)
}

// 设置数据库连接池最大打开的链接数量
func (bs *dbBase) SetMaxOpenConnCount(n int) {
	bs.maxOpenConnCount.Set(n)
}

// 设置数据库连接可重复利用的时间，超过该时间则被关闭废弃
// 如果 d <= 0 表示该链接会一直重复利用
func (bs *dbBase) SetMaxConnLifetime(n int) {
	bs.maxConnLifetime.Set(n)
}

// 节点配置转换为字符串
func (node *ConfigNode) String() string {
	if node.LinkInfo != "" {
		return node.LinkInfo
	}
	return fmt.Sprintf(`%s@%s:%s,%s,%s,%s,%s,%v,%d-%d-%d`, node.User, node.Host, node.Port,
		node.Name, node.Type, node.Role, node.Charset, node.Debug,
		node.MaxIdleConnCount, node.MaxOpenConnCount, node.MaxConnLifetime,
	)
}

// 是否开启调试服务
func (bs *dbBase) SetDebug(debug bool) {
	if bs.debug.Val() == debug {
		return
	}
	bs.debug.Set(debug)
	if debug && bs.sqls == nil {
		bs.sqls = gring.New(gDEFAULT_DEBUG_SQL_LENGTH)
	}
}

// 获取是否开启调试服务
func (bs *dbBase) getDebug() bool {
	return bs.debug.Val()
}

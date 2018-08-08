// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 数据库ORM.
// 默认内置支持MySQL, 其他数据库需要手动import对应的数据库引擎第三方包.
package gdb

import (
	"fmt"
	"time"
    "errors"
    "database/sql"
	"gitee.com/johng/gf/g/container/gmap"
	"gitee.com/johng/gf/g/container/gring"
	"gitee.com/johng/gf/g/container/gtype"
	"gitee.com/johng/gf/g/os/gcache"
	"gitee.com/johng/gf/g/util/grand"
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
	Open(c *ConfigNode) (*sql.DB, error)

	// SQL操作方法
	Query(q string, args ...interface{}) (*sql.Rows, error)
	Exec(q string, args ...interface{}) (sql.Result, error)
	Prepare(q string) (*sql.Stmt, error)

	// 数据库查询
	GetAll(q string, args ...interface{}) (Result, error)
	GetOne(q string, args ...interface{}) (Record, error)
	GetValue(q string, args ...interface{}) (Value, error)

	// Ping
	PingMaster() error
	PingSlave() error

	// 连接属性设置
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	SetConnMaxLifetime(d time.Duration)

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
	Table(tables string) *Model
	From(tables string) *Model

	// 关闭数据库操作对象
	Close() error

	// 内部方法
	insert(table string, data Map, option uint8) (sql.Result, error)
	batchInsert(table string, list List, batch int, option uint8) (sql.Result, error)

	getQuoteCharLeft() string
	getQuoteCharRight() string
	handleSqlBeforeExec(q *string) *string
}

// 数据库链接对象
type Db struct {
	link   Link          // 底层数据库类型管理对象
	master *sql.DB       // 实例化数据库链接(master)
	slave  *sql.DB       // 实例化数据库链接(slave，可能会与master相同)
	charl  string        // SQL安全符号(左)
	charr  string        // SQL安全符号(右)
	debug  *gtype.Bool   // (默认关闭)是否开启调试模式，当开启时会启用一些调试特性
	sqls   *gring.Ring   // (debug=true时有效)已执行的SQL列表
	cache  *gcache.Cache // 查询缓存，需要注意的是，事务查询不支持缓存
}

// 执行的SQL对象
type Sql struct {
	Sql   string        // SQL语句(可能带有预处理占位符)
	Args  []interface{} // 预处理参数值列表
	Error error         // 执行结果(nil为成功)
	Cost  int64         // 执行时间消耗(毫秒)
	Func  string        // 执行方法名称
}

// 返回数据表记录值
type Value []byte

// 返回数据表记录Map
type Record map[string]Value

// 返回数据表记录List
type Result []Record

// 关联数组，绑定一条数据表记录(使用别名)
type Map = map[string]interface{}

// 关联数组列表(索引从0开始的数组)，绑定多条记录(使用别名)
type List = []Map

// MySQL接口对象
var linkMysql = &dbmysql{}

// PostgreSQL接口对象
var linkPgsql = &dbpgsql{}

// Sqlite接口对象
// @author wxkj<wxscz@qq.com>
var linkSqlite = &dbsqlite{}

// 数据库查询缓存对象map，使用数据库连接名称作为键名，键值为查询缓存对象
var dbCaches = gmap.NewStringInterfaceMap()

// 使用默认/指定分组配置进行连接，数据库集群配置项：default
func New(groupName ...string) (*Db, error) {
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
		masterList := make(ConfigGroup, 0)
		slaveList := make(ConfigGroup, 0)
		for i := 0; i < len(list); i++ {
			if list[i].Role == "slave" {
				slaveList = append(slaveList, list[i])
			} else {
				// 默认配置项的角色为master
				masterList = append(masterList, list[i])
			}
		}
		if len(masterList) < 1 {
			return nil, errors.New("at least one master node configuration's need to make sense")
		}
		masterNode := getConfigNodeByPriority(masterList)
		var slaveNode *ConfigNode
		if len(slaveList) > 0 {
			slaveNode = getConfigNodeByPriority(slaveList)
		}
		return newDb(masterNode, slaveNode, name)
	} else {
		return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", name))
	}
}

// 按照负载均衡算法(优先级配置)从数据库集群中选择一个配置节点出来使用
// 算法说明举例，
// 1、假如2个节点的priority都是1，那么随机大小范围为[0, 199]；
// 2、那么节点1的权重范围为[0, 99]，节点2的权重范围为[100, 199]，比例为1:1；
// 3、假如计算出的随机数为99;
// 4、那么选择的配置为节点1;
func getConfigNodeByPriority(cg ConfigGroup) *ConfigNode {
	if len(cg) < 2 {
		return &cg[0]
	}
	var total int
	for i := 0; i < len(cg); i++ {
		total += cg[i].Priority * 100
	}
	// 不能取到末尾的边界点
	r := grand.Rand(0, total)
	if r > 0 {
		r -= 1
	}
	min := 0
	max := 0
	for i := 0; i < len(cg); i++ {
		max = min + cg[i].Priority*100
		//fmt.Printf("r: %d, min: %d, max: %d\n", r, min, max)
		if r >= min && r < max {
			return &cg[i]
		} else {
			min = max
		}
	}
	return nil
}

// 创建数据库链接对象
func newDb(masterNode *ConfigNode, slaveNode *ConfigNode, groupName string) (*Db, error) {
	var link Link
	switch masterNode.Type {
		case "mysql":
			link = linkMysql
		case "pgsql":
			link = linkPgsql
		case "sqlite":
			link = linkSqlite
		default:
			return nil, errors.New(fmt.Sprintf("unsupported db type '%s'", masterNode.Type))
	}
	master, err := link.Open(masterNode)
	if err != nil {
		return nil, err
	}
	slave := master
	if slaveNode != nil {
		slave, err = link.Open(slaveNode)
		if err != nil {
			return nil, err
		}
	}
	db := &Db{
		link:   link,
		master: master,
		slave:  slave,
		charl:  link.getQuoteCharLeft(),
		charr:  link.getQuoteCharRight(),
		debug:  gtype.NewBool(),
	}
	// 设置连接属性，master和slave必须是一致的，所以这里使用的是master的属性设置
	if masterNode.MaxIdleConnCount > 0 {
		db.SetMaxIdleConns(masterNode.MaxIdleConnCount)
	}
	if masterNode.MaxOpenConnCount > 0 {
		db.SetMaxOpenConns(masterNode.MaxOpenConnCount)
	}
	if masterNode.MaxConnLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(masterNode.MaxConnLifetime) * time.Second)
	}
	if v := dbCaches.Get(groupName); v == nil {
		dbCaches.LockFunc(func(m map[string]interface{}) {
			if v, ok := m[groupName]; !ok {
				db.cache = gcache.New()
				m[groupName] = db.cache
			} else {
				db.cache = v.(*gcache.Cache)
			}
		})
	} else {
		db.cache = v.(*gcache.Cache)
	}

	return db, nil
}

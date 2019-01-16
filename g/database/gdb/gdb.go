// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gdb provides ORM features for popular relationship databases.
//
// 数据库ORM,
// 默认内置支持MySQL, 其他数据库需要手动import对应的数据库引擎第三方包.
package gdb

import (
    "database/sql"
    "errors"
    "fmt"
    "gitee.com/johng/gf/g/container/gring"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/container/gvar"
    "gitee.com/johng/gf/g/os/gcache"
    "gitee.com/johng/gf/g/util/grand"
    _ "gitee.com/johng/gf/third/github.com/go-sql-driver/mysql"
    "time"
)

// 数据库操作接口
type DB interface {
    // 建立数据库连接方法(开发者一般不需要直接调用)
    Open(config *ConfigNode) (*sql.DB, error)

	// SQL操作方法 API
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	Prepare(sql string, execOnMaster...bool) (*sql.Stmt, error)

    // 内部实现API的方法(不同数据库可覆盖这些方法实现自定义的操作)
    doQuery(link dbLink, query string, args ...interface{}) (rows *sql.Rows, err error)
    doExec(link dbLink, query string, args ...interface{}) (result sql.Result, err error)
    doPrepare(link dbLink, query string) (*sql.Stmt, error)
    doInsert(link dbLink, table string, data Map, option int) (result sql.Result, err error)
    doBatchInsert(link dbLink, table string, list List, batch int, option int) (result sql.Result, err error)
    doUpdate(link dbLink, table string, data interface{}, condition interface{}, args ...interface{}) (result sql.Result, err error)
    doDelete(link dbLink, table string, condition interface{}, args ...interface{}) (result sql.Result, err error)

	// 数据库查询
	GetAll(query string, args ...interface{}) (Result, error)
	GetOne(query string, args ...interface{}) (Record, error)
	GetValue(query string, args ...interface{}) (Value, error)
    GetCount(query string, args ...interface{}) (int, error)
    GetStruct(obj interface{}, query string, args ...interface{}) error

    // 创建底层数据库master/slave链接对象
    Master() (*sql.DB, error)
    Slave() (*sql.DB, error)

    // Ping
	PingMaster() error
	PingSlave() error

	// 开启事务操作
	Begin() (*TX, error)

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

	// 设置管理
    SetDebug(debug bool)
    SetSchema(schema string)
    GetQueriedSqls() []*Sql
    PrintQueriedSqls()
    SetMaxIdleConns(n int)
    SetMaxOpenConns(n int)
    SetConnMaxLifetime(n int)

	// 内部方法接口
	getCache() (*gcache.Cache)
	getChars() (charLeft string, charRight string)
	getDebug() bool
    filterFields(table string, data map[string]interface{}) map[string]interface{}
    convertValue(fieldValue interface{}, fieldType string) interface{}
    getTableFields(table string) (map[string]string, error)
    rowsToResult(rows *sql.Rows) (Result, error)
    handleSqlBeforeExec(sql string) string
}

// 执行底层数据库操作的核心接口
type dbLink interface {
    Query(query string, args ...interface{}) (*sql.Rows, error)
    Exec(sql string, args ...interface{}) (sql.Result, error)
    Prepare(sql string) (*sql.Stmt, error)
}

// 数据库链接对象
type dbBase struct {
	db               DB                           // 数据库对象
	group            string                       // 配置分组名称
	debug            *gtype.Bool                  // (默认关闭)是否开启调试模式，当开启时会启用一些调试特性
	sqls             *gring.Ring                  // (debug=true时有效)已执行的SQL列表
	cache            *gcache.Cache                // 数据库缓存，包括底层连接池对象缓存及查询缓存；需要注意的是，事务查询不支持查询缓存
    schema           *gtype.String                // 手动切换的数据库名称
    tables           map[string]map[string]string // 数据库表结构
	maxIdleConnCount *gtype.Int                   // 连接池最大限制的连接数
    maxOpenConnCount *gtype.Int                   // 连接池最大打开的连接数
    maxConnLifetime  *gtype.Int                   // (单位秒)连接对象可重复使用的时间长度
}

// 执行的SQL对象
type Sql struct {
	Sql   string        // SQL语句(可能带有预处理占位符)
	Args  []interface{} // 预处理参数值列表
	Error error         // 执行结果(nil为成功)
	Start int64         // 执行开始时间(毫秒)
	End   int64         // 执行结束时间(毫秒)
	Func  string        // 执行方法
}

// 返回数据表记录值
type Value = *gvar.Var

// 返回数据表记录Map
type Record map[string]Value

// 返回数据表记录List
type Result []Record

// 关联数组，绑定一条数据表记录(使用别名)
type Map  = map[string]interface{}

// 关联数组列表(索引从0开始的数组)，绑定多条记录(使用别名)
type List = []Map

const (
    OPTION_INSERT  = 0
    OPTION_REPLACE = 1
    OPTION_SAVE    = 2
    OPTION_IGNORE  = 3
    // 默认的连接池连接存活时间(秒)
    gDEFAULT_CONN_MAX_LIFE_TIME = 30
)

// 使用默认/指定分组配置进行连接，数据库集群配置项：default
func New(groupName ...string) (db DB, err error) {
	group := config.d
	if len(groupName) > 0 {
        group = groupName[0]
	}
	config.RLock()
	defer config.RUnlock()

	if len(config.c) < 1 {
		return nil, errors.New("empty database configuration")
	}
	if _, ok := config.c[group]; ok {
	    if node, err := getConfigNodeByGroup(group, true); err == nil {
	        base := &dbBase {
                group            : group,
                debug            : gtype.NewBool(),
                cache            : gcache.New(),
                schema           : gtype.NewString(),
                maxIdleConnCount : gtype.NewInt(),
                maxOpenConnCount : gtype.NewInt(),
                maxConnLifetime  : gtype.NewInt(gDEFAULT_CONN_MAX_LIFE_TIME),
            }
            switch node.Type {
                case "mysql":
                    base.db = &dbMysql{dbBase  : base}
                case "pgsql":
                    base.db = &dbPgsql{dbBase  : base}
                case "mssql":
                    base.db = &dbMssql{dbBase  : base}
                case "sqlite":
                    base.db = &dbSqlite{dbBase : base}
                case "oracle":
                    base.db = &dbOracle{dbBase : base}
                default:
                    return nil, errors.New(fmt.Sprintf(`unsupported database type "%s"`, node.Type))
            }
            return base.db, nil
        } else {
            return nil, err
        }
	} else {
		return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", group))
	}
}

// 获取指定数据库角色的一个配置项，内部根据权重计算负载均衡
func getConfigNodeByGroup(group string, master bool) (*ConfigNode, error) {
    if list, ok := config.c[group]; ok {
        // 将master, slave集群列表拆分出来
        masterList := make(ConfigGroup, 0)
        slaveList  := make(ConfigGroup, 0)
        for i := 0; i < len(list); i++ {
            if list[i].Role == "slave" {
                slaveList = append(slaveList, list[i])
            } else {
                masterList = append(masterList, list[i])
            }
        }
        if len(masterList) < 1 {
            return nil, errors.New("at least one master node configuration's need to make sense")
        }
        if len(slaveList) < 1 {
            slaveList = masterList
        }
        if master {
            return getConfigNodeByPriority(masterList), nil
        } else {
            return getConfigNodeByPriority(slaveList), nil
        }
    } else {
        return nil, errors.New(fmt.Sprintf("empty database configuration for item name '%s'", group))
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
	// 如果total为0表示所有连接都没有配置priority属性，那么默认都是1
	if total == 0 {
        for i := 0; i < len(cg); i++ {
            cg[i].Priority = 1
            total         += cg[i].Priority * 100
        }
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

// 获得底层数据库链接对象
func (bs *dbBase) getSqlDb(master bool) (sqlDb *sql.DB, err error) {
    // 负载均衡
    node, err := getConfigNodeByGroup(bs.group, master)
    if err != nil {
        return nil, err
    }
    // 默认值设定
    if node.Charset == "" {
        node.Charset = "utf8"
    }
    v := bs.cache.GetOrSetFuncLock(node.String(), func() interface{} {
        sqlDb, err = bs.db.Open(node)
        if err != nil {
            return nil
        }

        if n := bs.maxIdleConnCount.Val(); n > 0 {
            sqlDb.SetMaxIdleConns(n)
        } else if node.MaxIdleConnCount > 0 {
            sqlDb.SetMaxIdleConns(node.MaxIdleConnCount)
        }

        if n := bs.maxOpenConnCount.Val(); n > 0 {
            sqlDb.SetMaxOpenConns(n)
        } else if node.MaxOpenConnCount > 0 {
            sqlDb.SetMaxOpenConns(node.MaxOpenConnCount)
        }

        if n := bs.maxConnLifetime.Val(); n > 0 {
            sqlDb.SetConnMaxLifetime(time.Duration(n) * time.Second)
        } else if node.MaxConnLifetime > 0 {
            sqlDb.SetConnMaxLifetime(time.Duration(node.MaxConnLifetime) * time.Second)
        }
        return sqlDb
    }, 0)
    if v != nil && sqlDb == nil {
        sqlDb = v.(*sql.DB)
    }
    // 是否手动选择数据库
    if v := bs.schema.Val(); v != "" {
        sqlDb.Exec("USE " + v)
    }
    return
}

// 切换操作的数据库(注意该切换是全局的)
func (bs *dbBase) SetSchema(schema string) {
    bs.schema.Set(schema)
}

// 创建底层数据库master链接对象
func (bs *dbBase) Master() (*sql.DB, error) {
	return bs.getSqlDb(true)
}

// 创建底层数据库slave链接对象
func (bs *dbBase) Slave() (*sql.DB, error) {
    return bs.getSqlDb(false)
}

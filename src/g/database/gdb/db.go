package gdb

import (
    "database/sql"
    "errors"
    "fmt"
    "g/core/ginstance"
    "log"
    "g/util/grand"
    "sync"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/lib/pq"
)

const (
    OPTION_INSERT  = 0
    OPTION_REPLACE = 1
    OPTION_SAVE    = 2
    OPTION_IGNORE  = 3
)

// 数据库配置包内对象
var config struct {
    sync.RWMutex
    c  Config // 数据库配置
    d  string // 默认数据库分组名称
}

// 数据库操作接口
type Link interface {
    Open (c *ConfigNode) (*sql.DB, error)
    Close() error
    Query(q string, args ...interface{}) (*sql.Rows, error)
    Exec(q string, args ...interface{}) (sql.Result, error)
    Prepare(q string) (*sql.Stmt, error)

    GetAll(q string, args ...interface{}) (*List, error)
    GetOne(q string, args ...interface{}) (*Map, error)
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

    insert(table string, data *Map, option uint8) (sql.Result, error)
    Insert(table string, data *Map) (sql.Result, error)
    Replace(table string, data *Map) (sql.Result, error)
    Save(table string, data *Map) (sql.Result, error)

    batchInsert(table string, list *List, batch int, option uint8) error
    BatchInsert(table string, list *List, batch int) error
    BatchReplace(table string, list *List, batch int) error
    BatchSave(table string, list *List, batch int) error

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

// 数据库配置
type Config      map[string]ConfigGroup

// 数据库集群配置
type ConfigGroup []ConfigNode

// 数据库单项配置
type ConfigNode  struct {
    Host     string // 地址
    Port        int // 端口
    User     string // 账号
    Pass     string // 密码
    Name     string // 数据库名称
    Type     string // 数据库类型：mysql, sqlite, mssql, pgsql, oracle(目前仅支持mysql)
    Role     string // (可选)数据库的角色，用于主从操作分离，至少需要有一个master，参数值：master, slave
    Charset  string // (可选)编码，默认为 utf-8
    Priority    int // (可选)用于负载均衡的权重计算，当集群中只有一个节点时，权重没有任何意义
    Linkinfo string // (可选)自定义链接信息，当该字段被设置值时，以上链接字段(Host,Port,User,Pass,Name)将失效(该字段是一个扩展功能)
}

// 关联数组，绑定一条数据表记录
type Map  map[string]interface{}

// 关联数组列表(索引从0开始的数组)，绑定多条记录
type List []Map

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
            Charset  : "utf-8",
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
            Charset  : "utf-8",
            Priority : 100,
        },
    },
}
*/

// 包初始化
func init() {
    config.c = make(Config)
    config.d = "default"
}

// 设置当前应用的数据库配置信息，进行全局数据库配置覆盖操作
// 支持三种数据类型的输入参数：Config, ConfigGroup, ConfigNode
func SetConfig (c interface{}) error {
    config.Lock()
    defer config.Unlock()

    switch c.(type) {
        case Config:
            config.c = c.(Config)

        case ConfigGroup:
            config.c = Config {"default" : c.(ConfigGroup)}

        case ConfigNode:
            config.c = Config {"default" : ConfigGroup { c.(ConfigNode) }}

        default:
            return errors.New("invalid config type, types should be in: Config, ConfigGroup, ConfigNode")
    }
    return nil
}

// 添加一台数据库服务器配置
func AddConfigNode (group string, node ConfigNode) {
    config.Lock()
    config.c[group] = append(config.c[group], node)
    config.Unlock()
}

// 添加数据库服务器集群配置
func AddConfigGroup (group string, nodes ConfigGroup) {
    config.Lock()
    config.c[group] = nodes
    config.Unlock()
}

// 添加默认链接的一台数据库服务器配置
func AddDefaultConfigNode (node ConfigNode) {
    AddConfigNode("default", node)
}

// 添加默认链接的数据库服务器集群配置
func AddDefaultConfigGroup (nodes ConfigGroup) {
    AddConfigGroup("default", nodes)
}

// 设置默认链接的数据库链接配置项(默认是 default)
func SetDefaultGroup (groupName string) {
    config.Lock()
    config.d = groupName
    config.Unlock()
}

// 根据配置项获取一个数据库操作对象单例
func instance (groupName string) (Link, error) {
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
        return result.(Link), nil
    }
}

// 获得默认的数据库操作对象单例
func Instance () (Link, error) {
    return instance(config.d)
}

// 获得指定配置项的数据库草最对象单例
func InstanceByGroup(groupName string) (Link, error) {
    return instance(groupName)
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
        log.Fatal(err)
    }
    slave := master
    if slaveNode != nil {
        slave,  err = link.Open(slaveNode)
        if err != nil {
            log.Fatal(err)
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


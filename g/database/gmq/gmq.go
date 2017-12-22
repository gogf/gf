// Go语言实现的高性能消息队列。
// 每个xxx.mq消息文件存储10万条消息，顶部为消息索引域，底部为数据存储域。
// 顶部索引域：消息状态(1byte) 数据开始位置(40bit,1TB) 数据长度(24bit, 16MB)
// 底部数据域：[压缩消息数据](变长)

package gmq

import (
    "errors"
    "gitee.com/johng/gf/g/container/gmap"
    "gitee.com/johng/gf/g/os/gfile"
    "sync"
)

const (
    gDEAFULT_MQGROUP_NAME   = "default" // 默认的队列分组名称
    gMQFILE_MAX_COUNT       = 100000    // 每个队列文件存储的消息条数上限(不能随便该，和数据结构设计有关系)
    gMQFILE_INDEX_ITEM_SIZE = 9         // 每个队列文件的索引域大小
    gMQFILEPOOL_TIMEOUT     = 60        // 消息队列文件指针池过期时间
    gMQFILE_INDEX_LENGTH    = gMQFILE_MAX_COUNT*gMQFILE_INDEX_ITEM_SIZE // 消息队列文件索引域大小
)

// 消息队列管理对象
type MQ struct {
    path   string                   // 消息队列数据文件存放目录(绝对路径)
    groups *gmap.StringInterfaceMap // 所有的消息队列分类
}

// 消息队列分类管理对象
type MQGroup struct {
    mu    sync.RWMutex
    minid uint64 // 队列当前最小id
    maxid uint64 // 队列当前最大id
    path  string // 该分类的消息队列数据文件存放目录(绝对路径)
    name  string // 分组名称(命名规则和文件名一致，因为要生成对应的目录)
}

// 消息队列遍历器
type MQIterator struct {
    group *MQGroup // 遍历的分类
    id     uint64  // 当前遍历的消息id
}

// 创建消息队列管理对象
func New(path string) (*MQ, error) {
    if !gfile.Exists(path) {
        if err := gfile.Mkdir(path); err != nil {
            return nil, errors.New("creating mq folder failed: " + err.Error())
        }
    }
    if !gfile.IsWritable(path) || !gfile.IsReadable(path) {
        return nil, errors.New("permission denied to mq folder: " + path)
    }
    mq := &MQ {
        path   : path,
        groups : gmap.NewStringInterfaceMap(),
    }
    return mq, nil
}

// 获取或者创建一个消息队列分类
func (mq *MQ) Group(name string) *MQGroup {
    if result := mq.groups.Get(name); result != nil {
        return result.(*MQGroup)
    }
    path  := mq.path + gfile.Separator + name
    mqg   := &MQGroup {path : path, name : name}
    mqg.init()
    mq.groups.Set(name, mqg)
    return mqg
}

// 向默认分类队列写入消息
func (mq *MQ) Push(msg []byte) (uint64, error) {
    return mq.Group(gDEAFULT_MQGROUP_NAME).Push(msg)
}

// 向默认分类队列头获取消息
func (mq *MQ) Pop() []byte {
    return mq.Group(gDEAFULT_MQGROUP_NAME).Pop()
}
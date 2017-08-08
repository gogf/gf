/*
    使用raft算法处理集群的一致性
    @todo 解决split brains造成的数据一致性问题
 */

package graft

import (
    "log"
    "os"
    "g/core/types/gmap"
    "sync"
    "g/core/types/glist"
)

const (
    gVERSION                    = "0.6"   // 当前版本
    // 集群端口定义
    gPORT_RAFT                  = 4166    // 集群协议通信接口
    gPORT_REPL                  = 4167    // 集群数据同步接口
    gPORT_API                   = 4168    // 服务器对外API接口

    // 节点状态
    gSTATUS_ALIVE               = 1
    gSTATUS_DEAD                = 0

    // raft 角色
    gROLE_FOLLOWER              = 0
    gROLE_CANDIDATE             = 1
    gROLE_LEADER                = 2

    // 超时时间设置
    gTCP_RETRY_COUNT            = 3
    gTCP_READ_TIMEOUT           = 3000    // 毫秒
    gTCP_WRITE_TIMEOUT          = 3000    // 毫秒
    gELECTION_TIMEOUT_MIN       = 1000    // 毫秒
    gELECTION_TIMEOUT_MAX       = 3000    // 毫秒
    gELECTION_TIMEOUT_HEARTBEAT = 500     // 毫秒
    gLOG_REPL_TIMEOUT_HEARTBEAT = 1000    // 毫秒
    gLOG_REPL_AUTOSAVE_INTERVAL = 5000    // 毫秒

    // 选举操作
    gMSG_HEAD_HI                = iota
    gMSG_HEAD_HI2
    gMSG_HEAD_HEARTBEAT
    gMSG_HEAD_I_AM_LEADER
    gMSG_HEAD_RAFT_RESPONSE
    gMSG_HEAD_SCORE_REQUEST
    gMSG_HEAD_SCORE_COMPARE_REQUEST
    gMSG_HEAD_SCORE_COMPARE_FAILURE
    gMSG_HEAD_SCORE_COMPARE_SUCCESS
    gMSG_HEAD_SPLIT_BRAINS_OCCURRED

    // 数据同步操作
    gMSG_HEAD_SET
    gMSG_HEAD_REMOVE
    gMSG_HEAD_UPDATE
    gMSG_HEAD_REPL_HEARTBEAT
    gMSG_HEAD_REPL_RESPONSE
    gMSG_HEAD_REPL_NEED_UPDATE_LEADER
    gMSG_HEAD_REPL_NEED_UPDATE_FOLLOWER

    // API相关
    gMSG_HEAD_PEERS_INFO
    gMSG_HEAD_PEERS_ADD
    gMSG_HEAD_PEERS_REMOVE
)

// 消息
type Msg struct {
    Head int
    Body string
    Info NodeInfo
}

// 服务器节点信息
type Node struct {
    mutex            sync.RWMutex

    Name             string                   // 节点名称
    Ip               string                   // 主机节点的局域网ip
    Peers            *gmap.StringInterfaceMap // 集群所有的节点信息(ip->节点信息)，不包含自身
    Role             int                      // raft角色
    Leader           string                   // Leader节点ip
    Score            int64                    // 选举比分
    ScoreCount       int                      // 选举比分的节点数
    ElectionDeadline int64                    // 选举超时时间点

    LastLogId        int64                    // 最后一次未保存log的id，用以数据同步识别
    LastSavedLogId   int64                    // 最后一次物理化log的id，用以物理化保存识别
    LogList          *glist.SafeList          // leader日志列表，用以数据同步
    LogCount         int64                    // 日志的总数，用以核对一致性
    DataPath         string                   // 物理存储的本地数据目录绝对路径
    KVMap            *gmap.StringStringMap    // 存储的K-V哈希表
}

// 节点信息
type NodeInfo struct {
    Name             string
    Ip               string
    Status           int
    Role             int
    Score            int64
    ScoreCount       int
    LastLogId        int64
    LogCount         int64
    LastHeartbeat    int64  // 上一次心跳检查的毫秒数
    Version          string
}

// 数据保存结构体
type SaveInfo struct {
    LastLogId        int64
    LogCount         int64
    Peers            map[string]interface{}
    DataMap          map[string]string
}

// 日志记录项
type LogEntry struct {
    Id               int64                  // 唯一ID
    Act              int
    Items            interface{}            // map[string]string或[]string
}

// 绑定本地IP并创建一个服务节点
func NewServerByIp(ip string) *Node {
    hostname, err := os.Hostname()
    if err != nil {
        log.Fatalln("getting local hostname failed")
        return nil
    }
    node := Node {
        Name         : hostname,
        Ip           : ip,
        Role         : gROLE_FOLLOWER,
        Peers        : gmap.NewStringInterfaceMap(),
        DataPath     : os.TempDir(),
        LogList      : glist.NewSafeList(),
        KVMap        : gmap.NewStringStringMap(),
    }
    return &node
}

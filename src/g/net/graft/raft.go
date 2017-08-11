/*
    使用raft算法处理集群的一致性
    已解决split brains造成的数据一致性问题：
        经典分区多节点脑裂问题：出现两个及以上的分区网络，网络之间无法相互连接
        复杂非分区三点脑裂问题：A-B-C，AC之间无法相互连接（B是双网卡），这样会造成A、C为leader，B为follower
        以上需要解决的是数据一致性问题，解决方案：检测集群必需为环形网络，剔除掉非环形网络的节点
 */

package graft

import (
    "log"
    "os"
    "g/core/types/gmap"
    "sync"
    "g/core/types/glist"
    "g/net/ghttp"
)

const (
    gVERSION                        = "0.6"   // 当前版本
    // 集群端口定义
    gPORT_RAFT                      = 4166    // 集群协议通信接口
    gPORT_REPL                      = 4167    // 集群数据同步接口
    gPORT_API                       = 4168    // 服务器对外API接口
    gPORT_MONITOR                   = 4169    // 监控服务接口
    gPORT_WEB                       = 4170    // WEB管理界面

    // 节点状态
    gSTATUS_ALIVE                   = 1
    gSTATUS_DEAD                    = 0

    // raft 角色
    gROLE_FOLLOWER                  = 0
    gROLE_CANDIDATE                 = 1
    gROLE_LEADER                    = 2

    // 超时时间设置
    gTCP_RETRY_COUNT                = 3       // TCP请求失败时的重试次数
    gTCP_READ_TIMEOUT               = 3000    // (毫秒)TCP链接读取超时
    gTCP_WRITE_TIMEOUT              = 3000    // (毫秒)TCP链接写入超时
    gELECTION_TIMEOUT_MIN           = 1000    // (毫秒)RAFT选举超时最小值
    gELECTION_TIMEOUT_MAX           = 3000    // (毫秒)RAFT选举超时最大值
    gELECTION_TIMEOUT_HEARTBEAT     = 500     // (毫秒)RAFT Leader统治维持心跳间隔
    gLOG_REPL_TIMEOUT_HEARTBEAT     = 1000    // (毫秒)数据同步检测心跳间隔(数据包括kv数据及service数据)
    gLOG_REPL_AUTOSAVE_INTERVAL     = 5000    // (毫秒)数据自动物理化保存的间隔
    gLOG_REPL_LOGCLEAN_INTERVAL     = 5000    // (毫秒)数据同步时的日志清理间隔
    gSERVICE_HEALTH_CHECK_INTERVAL  = 5000    // (毫秒)健康检查默认间隔

    // RAFT操作
    gMSG_RAFT_HI                    = 110
    gMSG_RAFT_HI2                   = 120
    gMSG_RAFT_RESPONSE              = 130
    gMSG_RAFT_HEARTBEAT             = 140
    gMSG_RAFT_I_AM_LEADER           = 150
    gMSG_RAFT_SPLIT_BRAINS_CHECK    = 160
    gMSG_RAFT_SPLIT_BRAINS_UNSET    = 170
    gMSG_RAFT_SCORE_REQUEST         = 180
    gMSG_RAFT_SCORE_COMPARE_REQUEST = 190
    gMSG_RAFT_SCORE_COMPARE_FAILURE = 200
    gMSG_RAFT_SCORE_COMPARE_SUCCESS = 210

    // 数据同步操作
    gMSG_REPL_SET                   = 310
    gMSG_REPL_REMOVE                = 320
    gMSG_REPL_INCREMENTAL_UPDATE    = 330
    gMSG_REPL_COMPLETELY_UPDATE     = 340
    gMSG_REPL_HEARTBEAT             = 350
    gMSG_REPL_RESPONSE              = 360
    gMSG_REPL_NEED_UPDATE_LEADER    = 370
    gMSG_REPL_NEED_UPDATE_FOLLOWER  = 380

    // API相关
    gMSG_API_PEERS_INFO             = 510
    gMSG_API_PEERS_ADD              = 520
    gMSG_API_PEERS_REMOVE           = 530
    gMSG_API_SERVICE_SET            = 540
    gMSG_API_SERVICE_REMOVE         = 550
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
    Monitor          string                   // Monitor节点ip
    Score            int64                    // 选举比分
    ScoreCount       int                      // 选举比分的节点数
    ElectionDeadline int64                    // 选举超时时间点

    LastLogId        int64                    // 最后一次保存log的id，用以数据同步识别
    LastSavedLogId   int64                    // 最后一次物理化log的id，用以物理化保存识别
    LastServiceLogId int64                    // 最后一次保存的service id号，用以识别service信息同步
    LogChan          chan LogEntry            // 用于数据同步的通道
    LogList          *glist.SafeList          // leader日志列表，用以数据同步
    SavePath         string                   // 物理存储的本地数据目录绝对路径
    FileName         string                   // 数据文件名称(包含后缀)
    Service          *gmap.StringInterfaceMap // 存储的服务配置表
    KVMap            *gmap.StringStringMap    // 存储的K-V哈希表
}

// 服务配置对象
type Service struct {
    Name string                   `json:"name"`
    Type string                   `json:"type"`
    List []map[string]interface{} `json:"list"`
}

// 用于KV API接口的对象
type NodeApiKv struct {
    ghttp.ControllerBase
    node *Node
}

// 用于Node API接口的对象
type NodeApiNode struct {
    ghttp.ControllerBase
    node *Node
}

// 用于Service API接口的对象
type NodeApiService struct {
    ghttp.ControllerBase
    node *Node
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
    LastServiceLogId int64
    LastHeartbeat    int64  // 上一次心跳检查的毫秒数
    Version          string // 节点的版本
}

// 数据保存结构体
type SaveInfo struct {
    LastLogId        int64
    LastServiceLogId int64
    Service          map[string]interface{}
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
        SavePath     : os.TempDir(),
        FileName     : "graft.db",
        LogChan      : make(chan LogEntry, 1024),
        LogList      : glist.NewSafeList(),
        Service      : gmap.NewStringInterfaceMap(),
        KVMap        : gmap.NewStringStringMap(),
    }
    return &node
}

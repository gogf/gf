/*
    使用raft算法处理集群的一致性
    @todo 当集群节点 < 3时的leader选取问题
 */

package graft

import (
    "log"
    "os"
    "g/core/types/gmap"
    "sync"
    "g/core/types/glist"
    "g/os/gfile"
)

const (
    // 集群端口定义
    gCLUSTER_PORT_RAFT    = 4166 // 集群协议通信接口
    gCLUSTER_PORT_REPLI   = 4167 // 集群数据同步接口
    gCLUSTER_PORT_API     = 4168 // 服务器对外API接口

    // 集群角色
    gCLUSTER_ROLE_SERVER  = 0
    gCLUSTER_ROLE_CLIENT  = 1
    gCLUSTER_ROLE_MONITOR = 2

    // raft 角色
    gRAFT_ROLE_FOLLOWER   = 0
    gRAFT_ROLE_CANDIDATE  = 1
    gRAFT_ROLE_LEADER     = 2

    // 超时时间设置
    gRAFT_ELECTION_TIMEOUT_MIN = 150    // 毫秒， 官方推荐 150ms - 300ms
    gRAFT_ELECTION_TIMEOUT_MAX = 300   // 毫秒， 官方推荐 150ms - 300ms
    gRAFT_HEARTBEAT_TIMEOUT    = 100    // 毫秒

    gRAFT_MSG_HEAD_HI          = iota
    gRAFT_MSG_HEAD_HI2
    gRAFT_MSG_HEAD_HEARTBEAT
    gRAFT_MSG_HEAD_KEEPALIVED
    gRAFT_MSG_HEAD_I_AM_LEADER
    gRAFT_MSG_HEAD_VOTE_REQUEST
    gRAFT_MSG_HEAD_VOTE_YES
    gRAFT_MSG_HEAD_VOTE_NO

    gREPLI_MSG_HEAD_SET    = 100
    gREPLI_MSG_HEAD_REMOVE = 101

)

// 消息
type Msg struct {
    Head int
    Body interface{}
    From MsgFrom
}

// 消息来源信息
type MsgFrom struct {
    Name string
    Role int
    RaftInfo struct{
        Role int
        Term int
    }
}

// 服务器节点信息
type Node struct {
    mutex     sync.RWMutex

    Name      string                // 节点名称
    Ip        string                // 主机节点的局域网ip
    Role      int                   // 集群角色
    Peers     *gmap.StringIntMap    // 集群所有的节点(ip->raft角色)，不包含自身
    RaftInfo  RaftInfo
    LastLogId int64                 // 最后一次物理化log的id，用以数据同步识别
    LogList   *glist.SafeList       // 未提交的日志列表
    LogChan   chan struct{}         // 用于数据同步的事件通知
    LogCount  uint64                // 日志的总数，用以核对一致性
    DataPath  string                // 物理存储的本地数据文件绝对路径
    KVMap     *gmap.StringStringMap // 存储的K-V哈希表
}

// raft信息结构体
type RaftInfo struct {
    Role             int          // raft角色
    Term             int          // 时间阶段，改进：只在数据冲突时作为处理冲突的判断条件之一
    Leader           string       // Leader节点ip
    VoteFor          string       // 当前node投票的节点
    VoteCount        int          // 获得的选票数量
    ElectionDeadline int64        // 毫秒
}

// 日志记录项
type LogEntry struct {
    Id               int64        // 唯一ID
    Act              int
    Key              string
    Value            string
}

// 日志请求，用以向leader发送日志请求，不带ID，ID由leader统一生成
type LogRequest struct {
    Key              string
    Value            string
}

// 创建一个服务节点
func NewServer(ip string) *Node {
    hostname, err := os.Hostname()
    if err != nil {
        log.Fatalln("getting local hostname failed")
        return nil
    }
    node := Node {
        Name     : hostname,
        Ip       : ip,
        Role     : gCLUSTER_ROLE_SERVER,
        Peers    : gmap.NewStringIntMap(),
        DataPath : os.TempDir() + gfile.Separator + "graft.db",
        LogList  : glist.NewSafeList(),
        LogChan  : make(chan struct{}, 1024),
        KVMap    : gmap.NewStringStringMap(),
    }
    return &node
}

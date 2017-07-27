package graft

import (
    "time"
    "log"
    "os"
)

const (
    // 集群端口定义
    gCLUSTER_PORT_RAFT    = 4166 // 集群协议通信接口
    gCLUSTER_PORT_REPLI   = 4167 // 集群数据同步接口
    // 集群角色
    gCLUSTER_ROLE_SERVER  = 0
    gCLUSTER_ROLE_CLIENT  = 1
    gCLUSTER_ROLE_MONITOR = 2
    // raft 角色
    gRAFT_ROLE_FOLLOWER   = 0
    gRAFT_ROLE_CANDIDATE  = 1
    gRAFT_ROLE_LEADER     = 2
    // 超时时间设置
    gELECTION_TIMEOUT_MIN = 500  * time.Millisecond // 官方推荐 150ms - 300ms
    gELECTION_TIMEOUT_MAX = 1000 * time.Millisecond // 官方推荐 150ms - 300ms
)

// 通信消息结构体
type Msg struct {
    Act  string
    Data interface{}
    From MsgFrom
}

// 消息来源节点信息
type MsgFrom struct {
    Name string
    Role int
    RaftInfo struct{
        Role int
        Term int
    }
}

// 节点结构体
type Node struct {
    Name     string            // 节点名称
    Ip       string            // 主机节点的局域网ip
    Role     int               // 集群角色
    Peers    map[string]int    // 集群所有的节点(ip->raft角色)，不包含自身
    RaftInfo RaftInfo
}

// raft信息结构体
type RaftInfo struct {
    Role       int             // raft角色
    Term       int             // 时间阶段
    Vote       map[string]bool // 当前node投票的节点
    Leader     string          // Leader节点ip
    VoteCount  int             // 获得的选票数量
    TotalCount int             // 总共节点数
}


// 创建一个服务节点
func NewServer(ip string) *Node {
    hostname, err := os.Hostname()
    if err != nil {
        log.Fatalln("getting local hostname failed")
        return nil
    }
    node := Node{
        Name  : hostname,
        Ip    : ip,
        Role  : gCLUSTER_ROLE_SERVER,
        Peers : make(map[string]int),
    }
    return &node
}

// 创建一个客户端节点
func NewClient(ip string) *Node {
    hostname, err := os.Hostname()
    if err != nil {
        log.Fatalln("getting local hostname failed")
        return nil
    }
    node := Node{
        Name  : hostname,
        Ip    : ip,
        Role  : gCLUSTER_ROLE_CLIENT,
        Peers : make(map[string]int),
    }
    return &node
}

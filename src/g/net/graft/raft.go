package graft

import "time"

const (
    // raft 角色
    ROLE_FOLLOWER  = 0
    ROLE_CANDIDATE = 1
    ROLE_LEADER    = 2
    // 超时时间设置
    ELECTION_TIMEOUT_MIN = 500  * time.Millisecond // 官方推荐 150ms - 300ms
    ELECTION_TIMEOUT_MAX = 1000 * time.Millisecond // 官方推荐 150ms - 300ms
)

// 节点结构体
type Node struct {
    Name       string          // 节点名称
    Role       int             // 节点角色
    Term       int             // 时间阶段
    Vote       Node            // 投票的节点
    Leader     Node            // Leader节点
    VoteCount  int             // 获得的选票数量
    TotalCount int             // 总共节点数
    Peers      map[string]Node // 集群所有的节点(ip->信息)
}
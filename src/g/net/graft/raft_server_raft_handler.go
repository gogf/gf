package graft

import (
    "g/net/gip"
    "net"
    "sync"
    "time"
    "g/util/gtime"
    "log"
    "g/util/grand"
    "g/core/types/gset"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg       := n.receiveMsg(conn)
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    if msg == nil {
        conn.Close()
        return
    }
    // 保存peers
    n.Peers.Set(fromip, 0)
    // 消息处理
    switch msg.Head {
        // 上线通知
        case gRAFT_MSG_HEAD_HI:
            n.Peers.Set(fromip, msg.From.RaftInfo.Role)
            n.sendMsg(conn, gRAFT_MSG_HEAD_HI2, nil)
            //log.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())

        // 节点存活性保持
        case gRAFT_MSG_HEAD_KEEPALIVED:
            //log.Println("keepalived from", fromip)
            n.sendMsg(conn, gRAFT_MSG_HEAD_KEEPALIVED, nil)

        // 心跳保持
        case gRAFT_MSG_HEAD_HEARTBEAT:
            n.updateRaftElectionDeadline()
            //log.Println("heartbeat from", fromip)
            if n.getRaftRole() == gRAFT_ROLE_LEADER && msg.From.RaftInfo.Role == gRAFT_ROLE_LEADER {
                // 脑裂问题判断，term最大的leader节点认定为集群的leader
                if n.getRaftTerm() > msg.From.RaftInfo.Term {
                    n.sendMsg(conn, gRAFT_MSG_HEAD_I_AM_LEADER, nil)
                    goto continueHandleConn
                } else {
                    log.Println("two leader occured, set ", fromip, "as my leader")
                    n.setRaftRole(gRAFT_ROLE_FOLLOWER)
                    n.setRaftLeader(fromip)
                }
            } else if n.getRaftLeader() != fromip {
                n.setRaftLeader(fromip)
            } else if n.getRaftTerm() != msg.From.RaftInfo.Term{
                n.setRaftTerm(msg.From.RaftInfo.Term)
            }
            n.sendMsg(conn, gRAFT_MSG_HEAD_HEARTBEAT, nil)


        // raft选举，在多节点选举中，比较关键的一个操作
        case gRAFT_MSG_HEAD_VOTE_REQUEST:
            var result int
            if n.getRaftRole() == gRAFT_ROLE_LEADER && n.getRaftTerm() >= msg.From.RaftInfo.Term {
                result = gRAFT_MSG_HEAD_I_AM_LEADER
                log.Println("vote i am leader for", fromip)
            } else if n.getRaftVoteFor() == "" {
                n.setRaftVoteFor(fromip)
                result = gRAFT_MSG_HEAD_VOTE_YES
                log.Println("vote yes for", fromip)
            } else {
                result = gRAFT_MSG_HEAD_VOTE_NO
                log.Println("vote no for", fromip)
            }

            n.sendMsg(conn, result,  nil)
    }

    continueHandleConn:
        n.raftTcpHandler(conn)
}


// 通过心跳维持集群统治，如果心跳不及时，那么选民会重新选举进入选举流程
// 每一台服务器节点保持一个tcp链接，异步goroutine保持心跳
// 每个节点中都会有一个线程处理该回调函数，但是只有leader节点才会激活
// 改进：
// 不仅是通过心跳维持集群统治，并且可以保持与其他节点的链接
func (n *Node) heartbeatHandler() {
    conns := gset.NewStringSet()
    for {
        ips := n.Peers.Keys()
        for _, ip := range ips {
            if conns.Contains(ip) {
                continue
            }
            c := n.getConn(ip, gCLUSTER_PORT_RAFT)
            if c == nil {
                n.Peers.Remove(ip)
                conns.Remove(ip)
                continue
            }
            conns.Add(ip)
            go func(conn net.Conn) {
                for {
                    msgstr := gRAFT_MSG_HEAD_KEEPALIVED
                    if n.getRaftRole() == gRAFT_ROLE_LEADER {
                        msgstr = gRAFT_MSG_HEAD_HEARTBEAT
                    }
                    ip, _ := gip.ParseAddress(conn.RemoteAddr().String())
                    n.sendMsg(conn, msgstr, nil)
                    conn.SetReadDeadline(time.Now().Add(3 * gRAFT_HEARTBEAT_TIMEOUT * time.Millisecond))
                    msg := n.receiveMsg(conn)
                    if msg == nil {
                        log.Println(ip, "no response, removing this peer")
                        n.Peers.Remove(ip)
                        conns.Remove(ip)
                        conn.Close()
                        return
                    } else {
                        switch msg.Head {
                            case gRAFT_MSG_HEAD_I_AM_LEADER:
                                //log.Println("two leader occured, set ", ip, "as my leader, done heartbeating")
                                n.setRaftRole(gRAFT_ROLE_FOLLOWER)
                                n.setRaftLeader(ip)
                        }
                    }
                    time.Sleep(gRAFT_HEARTBEAT_TIMEOUT * time.Millisecond)
                }
            }(c)
        }
        time.Sleep(gRAFT_HEARTBEAT_TIMEOUT * time.Millisecond)
    }
}

// 服务器节点选举
// 改进：
// 3个节点以内的集群也可以完成leader选举
func (n *Node) electionHandler() {
    n.updateRaftElectionDeadline()
    for {
        if n.getRaftRole() != gRAFT_ROLE_LEADER && gtime.Millisecond() >= n.getRaftElectionDeadline() {
            // 重新进入选举流程时，需要清空已有的信息
            if n.getRaftLeader() != "" {
                n.Peers.Remove(n.getRaftLeader())
            }
            if n.Peers.Size() > 0 {
                // 集群是2个节点及以上
                n.resetAsCandidate()
                n.addRaftTerm()
                n.beginVote()
            } else {
                // 集群目前仅有1个节点
                log.Println("only one node in this cluster, so i'll be the leader")
                n.setRaftRole(gRAFT_ROLE_LEADER)
                n.setRaftLeader(n.Ip)
            }
            n.updateRaftElectionDeadline()
            // 改进：
            // 如果本时间段候选人选举失败，在下一段选举时期内，角色重置为选民，并清空所有选票信息，
            // 采用 随机超时+避让策略 让集群更容易达成绝大多数的选举，以便快速选举
            if n.getRaftRole() != gRAFT_ROLE_LEADER {
                n.resetAsFollower()
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 一轮选举投票
func (n *Node) beginVote() {
    var wg sync.WaitGroup
    log.Println("begin new voting")
    ips := n.Peers.Keys()
    for _, ipstr := range ips {
        wg.Add(1)
        go func(ip string) {
            time.Sleep(5*time.Second)
            if n.getRaftLeader() != "" {
                wg.Done()
                return
            }
            conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
            if conn != nil {
                log.Println("request vote to", ip)
                n.sendMsg(conn, gRAFT_MSG_HEAD_VOTE_REQUEST, nil)
                msg := n.receiveMsg(conn)
                if msg != nil {
                    switch msg.Head {
                        // 同意投票
                        case gRAFT_MSG_HEAD_VOTE_YES:
                            log.Println("receive yes from", ip)
                            n.addRaftVoteCount()
                            if  n.getRaftLeader() == "" &&
                                n.getRaftRole() == gRAFT_ROLE_CANDIDATE &&
                                n.getRaftVoteCount() >= (int(n.Peers.Size()/2) + 1) {
                                log.Println("i am voted to be leader now, vote count:", n.getRaftVoteCount())
                                n.setRaftRole(gRAFT_ROLE_LEADER)
                                n.setRaftLeader(n.Ip)
                            }

                        // 对方是leader，并且term不小于自身，那么转变自身角色
                        case gRAFT_MSG_HEAD_I_AM_LEADER:
                            log.Println("receive iamleader from", ip)
                            if  msg.From.RaftInfo.Role == gRAFT_ROLE_LEADER &&
                                msg.From.RaftInfo.Term >= n.getRaftTerm() {
                                log.Println("set him as my leader, done voting")
                                n.setRaftLeader(ip)
                                n.setRaftRole(gRAFT_ROLE_FOLLOWER)
                            }
                    }
                } else {
                    log.Println("receive nil from", ip)
                    n.Peers.Remove(ip)
                }
                conn.Close()
            } else {
                log.Println("could not connect to", ip)
                n.Peers.Remove(ip)
            }
            wg.Done()
        }(ipstr)
    }
    wg.Wait()
}

func (n *Node) getRaftLeader() string {
    n.mutex.RLock()
    r := n.RaftInfo.Leader
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRaftRole() int {
    n.mutex.RLock()
    r := n.RaftInfo.Role
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRaftTerm() int {
    n.mutex.RLock()
    r := n.RaftInfo.Term
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRaftElectionDeadline() int64 {
    n.mutex.RLock()
    r := n.RaftInfo.ElectionDeadline
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRaftVoteFor() string {
    n.mutex.RLock()
    r := n.RaftInfo.VoteFor
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRaftVoteCount() int {
    n.mutex.RLock()
    r := n.RaftInfo.VoteCount
    n.mutex.RUnlock()
    return r
}

// 添加投票轮数
func (n *Node) addRaftTerm() {
    n.mutex.Lock()
    n.RaftInfo.Term++
    n.mutex.Unlock()
}

// 添加投票次数
func (n *Node) addRaftVoteCount() {
    n.mutex.Lock()
    n.RaftInfo.VoteCount++
    n.mutex.Unlock()
}

// 重置为候选者，并初始化投票给自己
func (n *Node) resetAsCandidate() {
    n.mutex.Lock()
    n.RaftInfo.Role      = gRAFT_ROLE_CANDIDATE
    n.RaftInfo.Leader    = ""
    n.RaftInfo.VoteFor   = n.Ip
    n.RaftInfo.VoteCount = 1
    n.mutex.Unlock()
}

// 重置为选民，并清空选票信息
func (n *Node) resetAsFollower() {
    n.mutex.Lock()
    n.RaftInfo.Role      = gRAFT_ROLE_FOLLOWER
    n.RaftInfo.Leader    = ""
    n.RaftInfo.VoteFor   = ""
    n.RaftInfo.VoteCount = 0
    n.mutex.Unlock()
}

func (n *Node) setRaftRole(role int) {
    n.mutex.Lock()
    n.RaftInfo.Role      = role
    n.mutex.Unlock()
}

func (n *Node) setRaftLeader(ip string) {
    n.mutex.Lock()
    n.RaftInfo.Leader    = ip
    n.mutex.Unlock()
}

func (n *Node) setRaftTerm(term int) {
    n.mutex.Lock()
    n.RaftInfo.Term = term
    n.mutex.Unlock()
}

func (n *Node) setRaftVoteFor(ip string) {
    n.mutex.Lock()
    n.RaftInfo.VoteFor = ip
    n.mutex.Unlock()
}

func (n *Node) setRaftVoteCount(count int) {
    n.mutex.Lock()
    n.RaftInfo.VoteCount = count
    n.mutex.Unlock()
}

// 更新选举截止时间
func (n *Node) updateRaftElectionDeadline() {
    n.mutex.Lock()
    n.RaftInfo.ElectionDeadline = gtime.Millisecond() + int64(grand.Rand(gRAFT_ELECTION_TIMEOUT_MIN, gRAFT_ELECTION_TIMEOUT_MAX))
    n.mutex.Unlock()
}
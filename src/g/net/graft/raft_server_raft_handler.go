package graft

import (
    "g/net/gip"
    "net"
    "sync"
    "time"
    "g/util/gtime"
    "log"
    "g/util/grand"
    "g/core/types/gmap"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg := n.recieveMsg(conn)
    if msg == nil {
        return
    }
    // 任何raft通信都会更新选举的超时时间
    n.updateRaftElectionDeadline()
    // 保存peers
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    n.Peers.Set(fromip, 0)
    // 消息处理
    switch msg.Head {
        // 上线通知
        case "hi":
            n.Peers.Set(fromip, msg.From.RaftInfo.Role)
            n.sendMsg(conn, "hi2", nil)
            //log.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())

        // 心跳保持
        case "heartbeat":
            log.Println("heartbeat from", fromip)
            n.sendMsg(conn, "heartbeat", nil)
            n.raftTcpHandler(conn)

        // raft选举
        case "voteme":
            n.mutex.Lock()
            result := ""
            if n.RaftInfo.Role == gRAFT_ROLE_LEADER {
                log.Println("i am leader vs", fromip)
                result = "imleader"
            } else {
                if n.RaftInfo.Vote == "" {
                    n.RaftInfo.Vote = fromip
                    result          = "yes"
                } else {
                    result = "no"
                }
            }
            n.mutex.Unlock()

            log.Println("vote", result, "for", fromip)
            n.sendMsg(conn, result,  nil)

        // raft选举完成
        case "electiondone":
            log.Println("electiondone from", fromip)
            n.checkAndsetRaftLeader(fromip, msg.From.RaftInfo.Term)

    }
}

// 服务器节点选举
func (n *Node) electionHandler() {
    n.updateRaftElectionDeadline()
    for {
        if n.getRaftRole() != gRAFT_ROLE_LEADER && n.getRaftLeader() == "" && gtime.Millisecond() >= n.getRaftElectionDeadline() {
            n.beginVote()
            n.updateRaftElectionDeadline()
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// leader heartbeat
// 每一台服务器节点保持一个tcp链接，异步goroutine保持心跳
func (n *Node) heartbeatHandler() {
    conns := gmap.NewStringInterfaceMap()
    for n.getRaftRole() == gRAFT_ROLE_LEADER {
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
            conns.Set(ip, c)
            // 异步心跳保持
            go func(conn net.Conn) {
                for {
                    if n.getRaftRole() != gRAFT_ROLE_LEADER {
                        log.Println("i am not leader, quit sending heartbeat")
                        conn.Close()
                        return
                    }
                    ip, _ := gip.ParseAddress(conn.RemoteAddr().String())
                    log.Println("sending heartbeat to", ip)
                    n.sendMsg(conn, "heartbeat", nil)
                    // 设置read的时间期限，防止无限期阻塞
                    conn.SetReadDeadline(time.Now().Add(gHEARTBEAT_TIMEOUT * time.Millisecond))
                    if n.recieveMsg(conn) == nil {
                        log.Println(ip, "no response, removing this peer")
                        n.Peers.Remove(ip)
                        conns.Remove(ip)
                        return
                    }
                    time.Sleep(gHEARTBEAT_TIMEOUT * time.Millisecond)
                }
            }(c)
        }
        time.Sleep(gHEARTBEAT_TIMEOUT * time.Millisecond)
    }
}

// 一轮选举投票
func (n *Node) beginVote() {
    var wg sync.WaitGroup
    if n.Peers.Size() < 2 {
        return
    }
    log.Println("begin new voting")
    n.addRaftTerm()
    n.resetAsCandidate()
    ips := n.Peers.Keys()
    for _, ipstr := range ips {
        go func(ip string) {
            wg.Add(1)
            if n.getRaftLeader() != "" {
                return
            }
            conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
            if conn != nil {
                log.Println("request vote to", ip)
                n.sendMsg(conn, "voteme", nil)
                msg := n.recieveMsg(conn)
                if msg != nil {
                    log.Println("recieve", msg.Head, "from", ip)
                    switch msg.Head {
                        // 同意投票
                        case "yes":
                            n.addRaftVoteCount()
                            if n.getRaftLeader() == "" && n.canBeLeader() {
                                log.Println("i am voted to be leader now, vote count:", n.getRaftVoteCount())
                                n.setRaftRole(gRAFT_ROLE_LEADER)
                                n.setRaftLeader(n.Ip)
                                n.electionDone()
                            }

                        // 对方是leader，那么转变自身角色
                        case "imleader":
                            if n.getRaftLeader() == "" {
                                log.Println("he is leader, so i should have a check myself")
                                if (!n.canBeLeader()) {
                                    log.Println("set him as my leader, done voting")
                                    n.setRaftLeader(ip)
                                }
                            }
                    }
                } else {
                    log.Println("recieve nil from", ip)
                }
                conn.Close()
            }
            wg.Done()
        }(ipstr)
    }
    wg.Wait()
    // 选举为leader失败，那么重新初始化为选民，并根据情况决定是否清空选票信息
    if n.getRaftRole() != gRAFT_ROLE_LEADER {
        if n.getRaftLeader() == "" {
            // 如果本轮没有任何节点被选举为leader，那么重置角色信息，并清空选票，以便进行下一轮选举
            n.resetAsFollower()
        } else {
            // 如果其他节点被选举为leader，那么保留自身的选票信息
            n.setRaftRole(gRAFT_ROLE_FOLLOWER)
            n.setRaftVoteCount(0)
        }
    }
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

// 检查如果当前服务器节点没有设置leader，那么设置ip为leader
func (n *Node) checkAndsetRaftLeader(ip string, term int) {
    if n.getRaftLeader() == "" {
        n.setRaftLeader(ip)
        n.setRaftTerm(term)
    }
}

// 重置为候选者，并初始化投票给自己
func (n *Node) resetAsCandidate() {
    n.mutex.Lock()
    n.RaftInfo.Role      = gRAFT_ROLE_CANDIDATE
    n.RaftInfo.Vote      = n.Ip
    n.RaftInfo.Leader    = ""
    n.RaftInfo.VoteCount = 1
    n.mutex.Unlock()
}

// 重置为选民，并清空选票信息
func (n *Node) resetAsFollower() {
    n.mutex.Lock()
    n.RaftInfo.Role      = gRAFT_ROLE_FOLLOWER
    n.RaftInfo.Vote      = ""
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

func (n *Node) setRaftVoteCount(count int) {
    n.mutex.Lock()
    n.RaftInfo.VoteCount = count
    n.mutex.Unlock()
}

// 更新选举截止时间
func (n *Node) updateRaftElectionDeadline() {
    n.mutex.Lock()
    n.RaftInfo.ElectionDeadline = gtime.Millisecond() + int64(grand.Rand(gELECTION_TIMEOUT_MIN, gELECTION_TIMEOUT_MAX))
    n.mutex.Unlock()
}

// 当前服务器节点石佛营满足leader的要求
func (n *Node) canBeLeader() bool {
    n.mutex.RLock()
    result := n.RaftInfo.VoteCount >= (int(n.Peers.Size()/2) + 1)
    n.mutex.RUnlock()
    return result
}

// 异步通知选举完成
func (n *Node) electionDone() {
    ips := n.Peers.Keys()
    for _, ipstr := range ips {
        go func(ip string) {
            conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
            if conn != nil {
                n.sendMsg(conn, "electiondone", nil)
                conn.Close()
            }
        }(ipstr)
    }
    // 开始心跳保持
    go n.heartbeatHandler()
}

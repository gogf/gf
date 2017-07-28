package graft

import (
    "g/net/gip"
    "net"
    "sync"
    "time"
    "fmt"
    "g/util/gtime"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg := n.recieveMsg(conn)
    if msg == nil {
        return
    }
    // 任何raft通信都会更新选举的超时时间
    n.updateElectionDeadline()

    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    switch msg.Head {
        // 上线通知
        case "hi":
            n.Peers.Set(fromip, msg.From.RaftInfo.Role)
            n.sendMsg(conn, "hi2", nil)
            //fmt.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())

        // 心跳保持
        case "heartbeat":
            fmt.Println("heartbeat from", fromip)
            //n.checkAndSetLeader(fromip)


        // raft选举
        case "voteme":
            result := ""
            n.mutex.Lock()
            if n.RaftInfo.Role == gRAFT_ROLE_LEADER {
                result = "imleader"
            } else {
                if n.RaftInfo.Vote == "" {
                    n.RaftInfo.Vote = fromip
                    result = "yes"
                } else {
                    result = "no"
                }
            }
            n.mutex.Unlock()

            fmt.Println("vote", result, "for", fromip)
            n.sendMsg(conn, result,  nil)

        // raft选举完成
        case "electiondone":
            fmt.Println("electiondone from", fromip)
            n.checkAndSetLeader(fromip)

    }
}

// 选举超时检查
func (n *Node) checkElectionTimeout() {
    n.updateElectionDeadline()
    for {
        if n.RaftInfo.Role != gRAFT_ROLE_LEADER &&
            n.RaftInfo.Leader == "" &&
            gtime.Millisecond() >= n.RaftInfo.ElectionDeadline {
            //fmt.Println(gtime.Millisecond(), ">", n.RaftInfo.ElectionDeadline)
            n.beginVote()
            n.updateElectionDeadline()
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// leader heartbeat
// 每一台服务器节点保持一个tcp链接，异步goroutine保持心跳
func (n *Node) checkHeartbeat() {
    for ip, _ := range n.Peers.M {
        c := n.getConn(ip, gCLUSTER_PORT_RAFT)
        if c != nil {
            go func(conn net.Conn) {
                for {
                    n.mutex.RLock()
                    role := n.RaftInfo.Role
                    n.mutex.RUnlock()
                    if role != gRAFT_ROLE_LEADER {
                        fmt.Println("i am not leader, quit sending heartbeat")
                        conn.Close()
                        return
                    }
                    fmt.Println("sending heartbeat to", conn.RemoteAddr())
                    n.sendMsg(conn, "heartbeat", nil)
                    time.Sleep(gHEARTBEAT_TIMEOUT * time.Millisecond)
                }
            }(c)
        } else {
            n.Peers.Remove(ip)
        }
    }
}

// 一轮选举投票
func (n *Node) beginVote() {
    var wg sync.WaitGroup
    if n.Peers.Size() < 2 {
        return
    }
    fmt.Println("begin new voting")
    n.RaftInfo.Term ++
    n.RaftInfo.Role      = gRAFT_ROLE_CANDIDATE
    n.RaftInfo.Vote      = ""
    n.RaftInfo.VoteCount = 1
    for ip, _ := range n.Peers.M {
        go func(ip string) {
            wg.Add(1)
            if n.RaftInfo.Leader != "" {
                return
            }
            conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
            if conn != nil {
                fmt.Println("request vote from", ip)
                n.sendMsg(conn, "voteme", nil)
                msg := n.recieveMsg(conn)
                fmt.Println("recieve msg", msg, "from", ip)
                if msg != nil {
                    switch msg.Head {
                        // 同意投票
                        case "yes":
                            n.mutex.Lock()
                            n.RaftInfo.VoteCount ++
                            n.mutex.Unlock()
                            if n.RaftInfo.Leader == "" && n.canBeLeader() {
                                n.mutex.Lock()
                                n.RaftInfo.Role   = gRAFT_ROLE_LEADER
                                n.RaftInfo.Leader = n.Ip
                                n.mutex.Unlock()
                                n.electionDone()
                            } else {

                            }

                        // 对方是leader，那么转变自身角色
                        case "imleader":
                            if n.RaftInfo.Leader == "" {
                                fmt.Println("he is leader, so i should have a check myself")
                                if (!n.canBeLeader()) {
                                    fmt.Println("set him as my leader, done voting")
                                    n.setLeader(ip)
                                }
                            }
                    }
                }
                conn.Close()
            }
            wg.Done()
        }(ip)
    }
    wg.Wait()
}

// 检查如果当前服务器节点没有设置leader，那么设置ip为leader
func (n *Node) checkAndSetLeader(ip string) {
    n.mutex.RLock()
    leader := n.RaftInfo.Leader
    n.mutex.RUnlock()
    if leader == "" {
        n.setLeader(ip)
    }
}


// 设置leader
func (n *Node) setLeader(ip string) {
    n.mutex.Lock()
    n.RaftInfo.Role      = gRAFT_ROLE_FOLLOWER
    n.RaftInfo.Leader    = ip
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
    for ip, _ := range n.Peers.M {
        go func() {
            conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
            if conn != nil {
                n.sendMsg(conn, "electiondone", nil)
                conn.Close()
            }
        }()
    }
    // 开始心跳保持
    go n.checkHeartbeat()
}

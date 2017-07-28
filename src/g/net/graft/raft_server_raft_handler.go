package graft

import (
    "g/net/gip"
    "net"
    "sync"
    "time"
    "fmt"
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

        // raft选举
        case "voteme":
            n.mutex.RLock()
            role := n.RaftInfo.Role
            n.mutex.RUnlock()
            if role != gRAFT_ROLE_LEADER {
                if n.RaftInfo.Vote == "" {
                    n.mutex.Lock()
                    n.RaftInfo.Vote = fromip
                    n.mutex.Unlock()
                    n.sendMsg(conn, "yes", nil)
                } else {
                    n.sendMsg(conn, "no",  nil)
                }
            } else {
                n.sendMsg(conn, "imleader",  nil)
            }

        // raft选举完成
        case "electiondone":
            fmt.Println("electiondone")
            n.mutex.Lock()
            n.RaftInfo.Role      = gRAFT_ROLE_FOLLOWER
            n.RaftInfo.Leader    = fromip
            //n.RaftInfo.Term      = msg.From.RaftInfo.Term
            n.RaftInfo.VoteCount = 0
            //n.RaftInfo.Vote      = ""
            n.mutex.Unlock()

    }
}

// 选举超时检查
func (n *Node) checkElectionTimeout() {
    for {
        if n.RaftInfo.Role != gRAFT_ROLE_LEADER && time.Now().UnixNano() >= n.RaftInfo.ElectionDeadline {
            n.beginVote()
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// leader heartbeat
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
                        conn.Close()
                        return
                    }
                    n.sendMsg(conn, "heartbeat", nil)
                    time.Sleep(gHEARTBEAT_TIMEOUT * time.Millisecond)
                }
            }(c)
        } else {
            n.Peers.Remove(ip)
        }
    }
}

// 选举投票
func (n *Node) beginVote() {
    var wg sync.WaitGroup
    n.mutex.Lock()
    n.RaftInfo.Term = 0
    n.mutex.Unlock()

    for {
        // 直到选取到leader才会退出
        if n.RaftInfo.Leader == "" {
            if n.Peers.Size() < 2 {
                time.Sleep(1 * time.Second)
                continue
            }

            n.RaftInfo.Term ++
            n.RaftInfo.Role      = gRAFT_ROLE_CANDIDATE
            n.RaftInfo.Vote      = n.Ip
            n.RaftInfo.VoteCount = 1
            for ip, _ := range n.Peers.M {
                if n.RaftInfo.Leader != "" {
                    break;
                }
                go func() {
                    wg.Add(1)
                    n.mutex.RLock()
                    leader := n.RaftInfo.Leader
                    n.mutex.RUnlock()
                    if leader != "" {
                        return
                    }

                    conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
                    if conn != nil {
                        n.sendMsg(conn, "voteme", nil)
                        msg := n.recieveMsg(conn)
                        if msg != nil {
                            switch msg.Head {
                                case "yes":
                                    n.mutex.Lock()
                                    n.RaftInfo.VoteCount ++
                                    if n.RaftInfo.VoteCount >= (int(n.Peers.Size()/2) + 1) {
                                        n.RaftInfo.Role   = gRAFT_ROLE_LEADER
                                        n.RaftInfo.Leader = n.Ip
                                    }
                                    n.mutex.Unlock()
                                    n.electionDone()
                                
                                case "imleader":
                                    n.mutex.Lock()
                                    n.RaftInfo.Leader = ip
                                    n.RaftInfo.Role   = gRAFT_ROLE_FOLLOWER
                                    n.mutex.Unlock()
                            }

                        }
                        conn.Close()
                    }
                    wg.Done()
                }()
            }
            wg.Wait()
        } else {
            break;
        }
    }
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
    n.checkHeartbeat()
}

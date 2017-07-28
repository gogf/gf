package graft

import (
    "g/net/gip"
    "net"
    "sync"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg := n.recieveMsg(conn)
    if msg == nil {
        return
    }
    fromip, _ := gip.ParseAddress(conn.RemoteAddr().String())
    switch msg.Head {
        // 上线通知
        case "hi":
            n.Peers.Set(fromip, msg.From.RaftInfo.Role)
            n.RaftInfo.TotalCount = n.Peers.Size()
            n.sendMsg(conn, "hi2", nil)
            //fmt.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())
        // raft选举
        case "voteme":
            if n.RaftInfo.Role != gRAFT_ROLE_LEADER {
                n.mutex.Lock()
                if n.RaftInfo.Vote == "" {
                    n.RaftInfo.Vote = fromip
                    n.mutex.Unlock()
                    n.sendMsg(conn, "yes", nil)
                } else {
                    n.mutex.Unlock()
                    n.sendMsg(conn, "no",  nil)
                }
            }
        // raft选举完成
        case "electiondone":
            n.mutex.Lock()
            n.RaftInfo.Role      = gRAFT_ROLE_FOLLOWER
            n.RaftInfo.Leader    = fromip
            n.RaftInfo.Term      = msg.From.RaftInfo.Term
            n.RaftInfo.VoteCount = 0
            n.RaftInfo.Vote      = ""
            n.mutex.Unlock()

    }
}

// 选举投票
func (n *Node) beginVote() {
    var wg sync.WaitGroup
    for {
        if n.RaftInfo.Leader == "" {
            n.RaftInfo.Term ++
            n.RaftInfo.Role = gRAFT_ROLE_CANDIDATE
            for ip, _ := range n.Peers.M {
                go func() {
                    wg.Add(1)
                    n.mutex.RLock()
                    if n.RaftInfo.Leader != "" {
                        n.mutex.RUnlock()
                        return
                    }
                    n.mutex.RUnlock()

                    conn := n.getConn(ip, gCLUSTER_PORT_RAFT)
                    if conn != nil {
                        n.sendMsg(conn, "voteme", nil)
                        msg := n.recieveMsg(conn)
                        if msg.Head == "yes" {
                            n.mutex.Lock()
                            n.RaftInfo.VoteCount ++
                            if n.RaftInfo.VoteCount >= (int(n.RaftInfo.TotalCount/2) + 1) {
                                n.RaftInfo.Role   = gRAFT_ROLE_LEADER
                                n.RaftInfo.Leader = n.Ip
                            }
                            n.mutex.Unlock()
                            // 执行选举成功通知
                        }
                        conn.Close()
                    }
                    wg.Done()
                }()
            }
            wg.Wait()
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
}

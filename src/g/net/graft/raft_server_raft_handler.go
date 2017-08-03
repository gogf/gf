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
    "g/core/types/gmap"
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
    n.Peers.Set(fromip, gSTATUS_ALIVE)

    // 消息处理
    switch msg.Head {
        // 上线通知
        case gMSG_HEAD_HI:
            n.sendMsg(conn, gMSG_HEAD_HI2, nil)
            //log.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())

        // 心跳保持
        case gMSG_HEAD_HEARTBEAT:
            n.updateElectionDeadline()
            result := gMSG_HEAD_HEARTBEAT
            if n.getRole() == gROLE_LEADER {
                // 脑裂问题处理
                if n.getScoreCount() > msg.From.ScoreCount {
                    result = gMSG_HEAD_I_AM_LEADER
                } else if n.getScoreCount() == msg.From.ScoreCount {
                    if n.getScore() > msg.From.Score {
                        result = gMSG_HEAD_I_AM_LEADER
                    } else if n.getScore() == msg.From.Score {
                        // 极少数情况会出现两个节点ScoreCount和Score都相等的情况, 这个时候采用随机策略
                        if grand.Rand(0, 1) == 1 {
                            result = gMSG_HEAD_I_AM_LEADER
                        }
                    }
                }
                if result == gMSG_HEAD_HEARTBEAT {
                    n.setLeader(fromip)
                    n.setRole(gROLE_FOLLOWER)
                }
            } else {
                if n.getLeader() != fromip {
                    n.setLeader(fromip)
                    n.setRole(gROLE_FOLLOWER)
                }
            }
            n.sendMsg(conn, result, nil)

        // 选举比分获取
        case gMSG_HEAD_SCORE_REQUEST:
            if n.getRole() == gROLE_LEADER {
                n.sendMsg(conn, gMSG_HEAD_I_AM_LEADER,  nil)
            } else {
                n.sendMsg(conn, gMSG_HEAD_SCORE_RESPONSE,  nil)
            }

        // 选举比分对比
        case gMSG_HEAD_SCORE_COMPARE_REQUEST:
            result := gMSG_HEAD_SCORE_COMPARE_SUCCESS
            if n.getRole() == gROLE_LEADER {
                result = gMSG_HEAD_I_AM_LEADER
            } else {
                if n.getScoreCount() > msg.From.ScoreCount {
                    result = gMSG_HEAD_SCORE_COMPARE_FAILURE
                } else if n.getScoreCount() == msg.From.ScoreCount {
                    if n.getScore() > msg.From.Score {
                        result = gMSG_HEAD_SCORE_COMPARE_FAILURE
                    } else if n.getScore() == msg.From.Score {
                        // 极少数情况会出现两个节点ScoreCount和Score都相等的情况, 这个时候采用随机策略
                        if grand.Rand(0, 1) == 1 {
                            result = gMSG_HEAD_SCORE_COMPARE_FAILURE
                        }
                    }
                } else {
                    result = gMSG_HEAD_SCORE_COMPARE_SUCCESS
                }
            }
            if result == gMSG_HEAD_SCORE_COMPARE_SUCCESS {
                n.setLeader(fromip)
                n.setRole(gROLE_FOLLOWER)
            }
            n.sendMsg(conn, result,  nil)

    }

    n.raftTcpHandler(conn)
}


// 通过心跳维持集群统治，如果心跳不及时，那么选民会重新进入选举流程
func (n *Node) heartbeatHandler() {
    conns := gset.NewStringSet()
    for {
        if n.getRole() == gROLE_LEADER {
            ips := n.Peers.Keys()
            for _, ip := range ips {
                if conns.Contains(ip) {
                    continue
                }
                conn := n.getConn(ip, gPORT_RAFT)
                if conn == nil {
                    n.Peers.Set(ip, gSTATUS_DEAD)
                    conns.Remove(ip)
                    continue
                }
                conns.Add(ip)
                go func(ip string, conn net.Conn) {
                    for {
                        if n.getRole() != gROLE_LEADER {
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        if err := n.sendMsg(conn, gMSG_HEAD_HEARTBEAT, nil); err != nil {
                            log.Println(err)
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            log.Println(ip, "no response, removing this peer")
                            n.Peers.Set(ip, gSTATUS_DEAD)
                            conns.Remove(ip)
                            conn.Close()
                            return
                        } else {
                            switch msg.Head {
                                case gMSG_HEAD_I_AM_LEADER:
                                    log.Println("two leader occured, set", ip, "as my leader, done heartbeating")
                                    n.setRole(gROLE_FOLLOWER)
                                    n.setLeader(ip)

                                default:
                                    time.Sleep(gELECTION_TIMEOUT_HEARTBEAT * time.Millisecond)
                            }
                        }
                    }
                }(ip, conn)
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 服务器节点选举
// 改进：
// 3个节点以内的集群也可以完成leader选举
func (n *Node) electionHandler() {
    for {
        if n.getRole() != gROLE_LEADER && gtime.Millisecond() >= n.getElectionDeadline() {
            // 重新进入选举流程时，需要清空已有的信息
            if n.getLeader() != "" {
                n.Peers.Set(n.getLeader(), gSTATUS_DEAD)
            }
            if n.Peers.Size() > 0 {
                // 集群是2个节点及以上
                n.resetAsCandidate()
                n.beginScore()
            } else {
                // 集群目前仅有1个节点
                log.Println("only one node in this cluster, so i'll be the leader")
                n.setRole(gROLE_LEADER)
                n.setLeader(n.Ip)
            }
            n.updateElectionDeadline()
            // 改进：采用 随机超时+避让策略 让集群更容易达成绝大多数的选举，以便快速选举
            if n.getRole() != gROLE_LEADER {
                n.resetAsFollower()
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 一轮选举比分
func (n *Node) beginScore() {
    var wg sync.WaitGroup
    log.Println("begin new score")
    conns := gmap.NewStringInterfaceMap()
    // 请求比分，获取比分数据
    for ip, status := range n.Peers.M {
        if status != gSTATUS_ALIVE {
            continue
        }
        wg.Add(1)
        go func(ip string) {
            if n.getLeader() != "" || n.getRole() != gROLE_CANDIDATE {
                wg.Done()
                return
            }
            stime := time.Now().UnixNano()
            conn  := n.getConnFromPool(ip, gPORT_RAFT, conns)
            if conn == nil {
                n.Peers.Set(ip, gSTATUS_DEAD)
                return
            }
            if err := n.sendMsg(conn, gMSG_HEAD_SCORE_REQUEST, nil); err != nil {
                log.Println(err)
                conn.Close()
                return
            }
            msg := n.receiveMsg(conn)
            if msg != nil {
                switch msg.Head {
                    case gMSG_HEAD_I_AM_LEADER:
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)

                    case gMSG_HEAD_SCORE_RESPONSE:
                        etime := time.Now().UnixNano()
                        score := etime - stime
                        n.addScore(score)
                        n.addScoreCount()
                }
            } else {
                n.Peers.Set(ip, gSTATUS_DEAD)
            }
            wg.Done()
        }(ip)
    }
    wg.Wait()

    // 执行比分，对比比分数据，选举出leader
    for ip, status := range n.Peers.M {
        if status != gSTATUS_ALIVE {
            continue
        }
        conn := n.getConnFromPool(ip, gPORT_RAFT, conns)
        if conn == nil {
            n.Peers.Set(ip, gSTATUS_DEAD)
            continue
        }
        wg.Add(1)
        go func(ip string, conn net.Conn) {
            defer func() {
                conn.Close()
                wg.Done()
            }()
            if n.getLeader() != "" || n.getRole() != gROLE_CANDIDATE {
                return
            }

            if err := n.sendMsg(conn, gMSG_HEAD_SCORE_COMPARE_REQUEST, nil); err != nil {
                log.Println(err)
                return
            }
            msg := n.receiveMsg(conn)
            if msg != nil {
                switch msg.Head {
                    case gMSG_HEAD_I_AM_LEADER:
                        log.Println("score comparison: get leader from", ip)
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)

                    case gMSG_HEAD_SCORE_COMPARE_FAILURE:
                        log.Println("score comparison: get failure from", ip)
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)

                    case gMSG_HEAD_SCORE_COMPARE_SUCCESS:
                        log.Println("score comparison: get success from", ip)
                }
            }
        }(ip, conn)
    }
    wg.Wait()

    // 如果其他节点均没有条件满足leader，那么选举自身为leader
    if n.getRole() != gROLE_FOLLOWER {
        log.Println("I've won this score comparison")
        n.setRole(gROLE_LEADER)
        n.setLeader(n.Ip)
    }
}

func (n *Node) getLeader() string {
    n.mutex.RLock()
    r := n.Leader
    n.mutex.RUnlock()
    return r
}

func (n *Node) getRole() int {
    n.mutex.RLock()
    r := n.Role
    n.mutex.RUnlock()
    return r
}

func (n *Node) getScore() int64 {
    n.mutex.RLock()
    r := n.Score
    n.mutex.RUnlock()
    return r
}

func (n *Node) getScoreCount() int {
    n.mutex.RLock()
    r := n.ScoreCount
    n.mutex.RUnlock()
    return r
}

func (n *Node) getElectionDeadline() int64 {
    n.mutex.RLock()
    r := n.ElectionDeadline
    n.mutex.RUnlock()
    return r
}

// 添加比分节
func (n *Node) addScore(s int64) {
    n.mutex.Lock()
    n.Score += s
    n.mutex.Unlock()
}

// 添加比分节点数
func (n *Node) addScoreCount() {
    n.mutex.Lock()
    n.ScoreCount++
    n.mutex.Unlock()
}

// 重置为候选者，并初始化投票给自己
func (n *Node) resetAsCandidate() {
    n.mutex.Lock()
    n.Role       = gROLE_CANDIDATE
    n.Leader     = ""
    n.Score      = 0
    n.ScoreCount = 0
    n.mutex.Unlock()
}

// 重置为选民，并清空选票信息
func (n *Node) resetAsFollower() {
    n.mutex.Lock()
    n.Role      = gROLE_FOLLOWER
    n.Leader    = ""
    n.Score      = 0
    n.ScoreCount = 0
    n.mutex.Unlock()
}

func (n *Node) setRole(role int) {
    n.mutex.Lock()
    n.Role = role
    n.mutex.Unlock()
}

func (n *Node) setLeader(ip string) {
    n.mutex.Lock()
    n.Leader    = ip
    n.mutex.Unlock()
}

// 更新选举截止时间
func (n *Node) updateElectionDeadline() {
    n.mutex.Lock()
    n.ElectionDeadline = gtime.Millisecond() + int64(grand.Rand(gELECTION_TIMEOUT_MIN, gELECTION_TIMEOUT_MAX))
    n.mutex.Unlock()
}
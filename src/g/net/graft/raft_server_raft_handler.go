package graft

import (
    "net"
    "sync"
    "time"
    "g/util/gtime"
    "log"
    "g/util/grand"
    "g/core/types/gset"
    "g/core/types/gmap"
    "g/encoding/gjson"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil {
        conn.Close()
        return
    }
    // 保存peers
    if msg.Info.Ip != n.Ip {
        if n.Peers.Contains(msg.Info.Ip) {
            n.updatePeerStatus(msg.Info.Ip, gSTATUS_ALIVE)
        } else {
            n.updatePeerInfo(msg.Info.Ip, msg.Info)
        }
    }

    // 消息处理
    switch msg.Head {
        // 上线通知
        case gMSG_HEAD_HI:
            n.sendMsg(conn, gMSG_HEAD_HI2, "")
            //log.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())

        // 心跳保持
        case gMSG_HEAD_HEARTBEAT:
            n.updateElectionDeadline()
            result := gMSG_HEAD_HEARTBEAT
            if n.getRole() == gROLE_LEADER {
                if n.getScoreCount() > msg.Info.ScoreCount {
                    result = gMSG_HEAD_I_AM_LEADER
                } else if n.getScoreCount() == msg.Info.ScoreCount {
                    if n.getScore() > msg.Info.Score {
                        result = gMSG_HEAD_I_AM_LEADER
                    } else if n.getScore() == msg.Info.Score {
                        // 极少数情况会出现两个节点ScoreCount和Score都相等的情况, 这个时候采用随机策略
                        if grand.Rand(0, 1) == 1 {
                            result = gMSG_HEAD_I_AM_LEADER
                        }
                    }
                }
                if result == gMSG_HEAD_HEARTBEAT {
                    n.setLeader(msg.Info.Ip)
                    n.setRole(gROLE_FOLLOWER)
                }
            } else if n.getLeader() == "" {
                n.setLeader(msg.Info.Ip)
                n.setRole(gROLE_FOLLOWER)
            } else {
                // 脑裂问题，一个节点处于两个网路中，并且两个网络的leader无法相互通信，会引起数据一致性问题
                if n.getLeader() != msg.Info.Ip {
                    if  n.getScoreCount() < msg.Info.ScoreCount ||
                        (n.getScoreCount() == msg.Info.ScoreCount && n.getLastLogId() < msg.Info.LastLogId) {
                        n.setLeader(msg.Info.Ip)
                        n.setRole(gROLE_FOLLOWER)
                    }
                }
            }
            n.sendMsg(conn, result, "")
            if result == gMSG_HEAD_HEARTBEAT {
                n.raftTcpHandler(conn)
            }

        // 选举比分获取
        case gMSG_HEAD_SCORE_REQUEST:
            if n.getRole() == gROLE_LEADER {
                n.sendMsg(conn, gMSG_HEAD_I_AM_LEADER, "")
            } else {
                n.sendMsg(conn, gMSG_HEAD_RAFT_RESPONSE, "")
            }

        // 选举比分对比
        case gMSG_HEAD_SCORE_COMPARE_REQUEST:
            result := gMSG_HEAD_SCORE_COMPARE_SUCCESS
            if n.getRole() == gROLE_LEADER {
                result = gMSG_HEAD_I_AM_LEADER
            } else {
                if n.getScoreCount() > msg.Info.ScoreCount {
                    result = gMSG_HEAD_SCORE_COMPARE_FAILURE
                } else if n.getScoreCount() == msg.Info.ScoreCount {
                    if n.getScore() > msg.Info.Score {
                        result = gMSG_HEAD_SCORE_COMPARE_FAILURE
                    } else if n.getScore() == msg.Info.Score {
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
                n.setLeader(msg.Info.Ip)
                n.setRole(gROLE_FOLLOWER)
            }
            n.sendMsg(conn, result, "")

        // 节点信息查询
        case gMSG_HEAD_PEERS_INFO:
            list := make([]NodeInfo, 0)
            list  = append(list, *n.getNodeInfo())
            for _, v := range n.Peers.Values() {
                list = append(list, v.(NodeInfo))
            }
            n.sendMsg(conn, gMSG_HEAD_PEERS_INFO, *gjson.Encode(list))

        // 新增节点
        case gMSG_HEAD_PEERS_ADD:
            list := make([]string, 0)
            gjson.DecodeTo(&msg.Body, &list)
            if list != nil && len(list) > 0 {
                for _, ip := range list {
                    if n.Peers.Contains(ip) {
                        continue
                    }
                    // log.Println("adding peer:", ip)
                    go func(ip string) {
                        conn := n.getConn(ip, gPORT_RAFT)
                        if conn != nil {
                            n.sendMsg(conn, gMSG_HEAD_HI, "")
                            msg := n.receiveMsg(conn)
                            if msg != nil && msg.Head == gMSG_HEAD_HI2{
                                n.updatePeerInfo(ip, msg.Info)
                            }
                        }
                        // 判断是否添加成功，如果没有，那么添加一个默认的信息
                        if !n.Peers.Contains(ip) {
                            info       := NodeInfo{}
                            info.Status = gSTATUS_DEAD
                            n.updatePeerInfo(ip, msg.Info)
                        }
                    }(ip)
                }
            }
            n.sendMsg(conn, gMSG_HEAD_RAFT_RESPONSE, "")

        // 删除节点
        case gMSG_HEAD_PEERS_REMOVE:
            list := make([]string, 0)
            gjson.DecodeTo(&msg.Body, &list)
            if list != nil && len(list) > 0 {
                for _, ip := range list {
                    // log.Println("removing peer:", ip)
                    n.Peers.Remove(ip)
                }
            }
            n.sendMsg(conn, gMSG_HEAD_RAFT_RESPONSE, "")
    }

    conn.Close()
}


// 通过心跳维持集群统治，如果心跳不及时，那么选民会重新进入选举流程
func (n *Node) heartbeatHandler() {
    conns := gset.NewStringSet()
    for {
        if n.getRole() == gROLE_LEADER {
            for _, v := range n.Peers.Values() {
                info := v.(NodeInfo)
                if conns.Contains(info.Ip) {
                    continue
                }
                conn := n.getConn(info.Ip, gPORT_RAFT)
                if conn == nil {
                    n.updatePeerStatus(info.Ip, gSTATUS_DEAD)
                    conns.Remove(info.Ip)
                    // 如果失联超过3天，那么将该节点移除
                    if gtime.Millisecond() - info.LastHeartbeat > 3 * 86400 * 1000 {
                        log.Println(info.Ip, "was dead over 3 days, removing from peers")
                        n.Peers.Remove(info.Ip)
                    }
                    continue
                }
                conns.Add(info.Ip)
                go func(ip string, conn net.Conn) {
                    for {
                        // 如果当前节点不再是leader，或者节点表中已经删除该节点信息
                        if n.getRole() != gROLE_LEADER || !n.Peers.Contains(ip){
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        if err := n.sendMsg(conn, gMSG_HEAD_HEARTBEAT, ""); err != nil {
                            log.Println(err)
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            log.Println(ip, "was dead")
                            n.updatePeerStatus(ip, gSTATUS_DEAD)
                            conns.Remove(ip)
                            conn.Close()
                            return
                        } else {
                            // 更新节点信息
                            n.updatePeerInfo(ip, msg.Info)
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
                }(info.Ip, conn)
            }
        }
        time.Sleep(gELECTION_TIMEOUT_HEARTBEAT * time.Millisecond)
    }
}

// 服务器节点选举
// 改进：
// 3个节点以内的集群也可以完成leader选举
func (n *Node) electionHandler() {
    n.updateElectionDeadline()
    for {
        if n.getRole() != gROLE_LEADER && gtime.Millisecond() >= n.getElectionDeadline() {
            // 重新进入选举流程时，需要清空已有的信息
            if n.getLeader() != "" {
                n.updatePeerStatus(n.getLeader(), gSTATUS_DEAD)
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
    for _, v := range n.Peers.Values() {
        info := v.(NodeInfo)
        if info.Status != gSTATUS_ALIVE {
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
                n.updatePeerStatus(ip, gSTATUS_DEAD)
                return
            }
            if err := n.sendMsg(conn, gMSG_HEAD_SCORE_REQUEST, ""); err != nil {
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

                    case gMSG_HEAD_RAFT_RESPONSE:
                        etime := time.Now().UnixNano()
                        score := etime - stime
                        n.addScore(score)
                        n.addScoreCount()
                }
            } else {
                n.updatePeerStatus(ip, gSTATUS_DEAD)
            }
            wg.Done()
        }(info.Ip)
    }
    wg.Wait()

    // 执行比分，对比比分数据，选举出leader
    for _, v := range n.Peers.Values() {
        info := v.(NodeInfo)
        if info.Status != gSTATUS_ALIVE {
            continue
        }
        conn := n.getConnFromPool(info.Ip, gPORT_RAFT, conns)
        if conn == nil {
            n.updatePeerStatus(info.Ip, gSTATUS_DEAD)
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

            if err := n.sendMsg(conn, gMSG_HEAD_SCORE_COMPARE_REQUEST, ""); err != nil {
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
        }(info.Ip, conn)
    }
    wg.Wait()

    // 如果其他节点均没有条件满足leader，那么选举自身为leader
    if n.getRole() != gROLE_FOLLOWER {
        log.Println("I've won this score comparison")
        n.setRole(gROLE_LEADER)
        n.setLeader(n.Ip)
    }
}

// 获取当前节点的信息
func (n *Node) getNodeInfo() *NodeInfo {
    return &NodeInfo {
        Name          : n.Name,
        Ip            : n.Ip,
        Status        : gSTATUS_ALIVE,
        Role          : n.getRole(),
        Score         : n.getScore(),
        ScoreCount    : n.getScoreCount(),
        LastLogId     : n.getLastLogId(),
        LogCount      : n.getLogCount(),
        LastHeartbeat : gtime.Millisecond(),
        Version       : gVERSION,
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

func (n *Node) updatePeerStatus(ip string, status int) {
    r := n.Peers.Get(ip)
    if r != nil {
        info       := r.(NodeInfo)
        info.Status = status
        if status == gSTATUS_ALIVE {
            info.LastHeartbeat = gtime.Millisecond()
        }
        n.Peers.Set(ip, info)
    }
}

// 更新节点信息
func (n *Node) updatePeerInfo(ip string, info NodeInfo) {
    n.Peers.Set(ip, info)
}

// 更新选举截止时间
func (n *Node) updateElectionDeadline() {
    n.mutex.Lock()
    n.ElectionDeadline = gtime.Millisecond() + int64(grand.Rand(gELECTION_TIMEOUT_MIN, gELECTION_TIMEOUT_MAX))
    n.mutex.Unlock()
}
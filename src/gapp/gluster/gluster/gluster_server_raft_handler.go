package gluster

import (
    "net"
    "g/util/grand"
    "g/encoding/gjson"
    "g/os/glog"
)

// 集群协议通信接口回调函数
func (n *Node) raftTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil || msg.Info.Group != n.Group {
        conn.Close()
        return
    }
    // 保存peers
    if msg.Info.Id != n.Id {
        if n.Peers.Contains(msg.Info.Id) {
            n.updatePeerStatus(msg.Info.Id, gSTATUS_ALIVE)
        } else {
            n.updatePeerInfo(msg.Info)
        }
    }
    // 消息处理
    switch msg.Head {
        case gMSG_RAFT_HI:                      n.onMsgRaftHi(conn, msg)
        case gMSG_RAFT_HEARTBEAT:               n.onMsgRaftHeartbeat(conn, msg)
        case gMSG_RAFT_SCORE_REQUEST:           n.onMsgRaftScoreRequest(conn, msg)
        case gMSG_RAFT_SCORE_COMPARE_REQUEST:   n.onMsgRaftScoreCompareRequest(conn, msg)
        case gMSG_RAFT_SPLIT_BRAINS_CHECK:      n.onMsgRaftSplitBrainsCheck(conn, msg)
        case gMSG_RAFT_SPLIT_BRAINS_UNSET:      n.onMsgRaftSplitBrainsUnset(conn, msg)
        case gMSG_API_PEERS_ADD:                n.onMsgApiPeersAdd(conn, msg)
        case gMSG_API_PEERS_REMOVE:             n.onMsgApiPeersRemove(conn, msg)
    }
    //这里不用自动关闭链接，由于链接有读取超时，当一段时间没有数据时会自动关闭
    n.raftTcpHandler(conn)
}

// 检测split brains问题，检查两个leader的连通性
// 如果不连通，那么follower保持当前leader不变
// 如果能够连通，那么需要在两个leader中确定一个
func (n *Node) onMsgRaftSplitBrainsCheck(conn net.Conn, msg *Msg) {
    checkip := msg.Body
    result  := gMSG_RAFT_SPLIT_BRAINS_UNSET
    tconn   := n.getConn(checkip, gPORT_RAFT)
    if tconn != nil {
        if n.sendMsg(tconn, gMSG_RAFT_HI, "") == nil {
            rmsg := n.receiveMsg(tconn)
            if rmsg != nil {
                n.updatePeerInfo(rmsg.Info)
                if n.getLogCount() < msg.Info.LogCount && n.getLastLogId() < msg.Info.LastLogId {
                    n.setLeader(&rmsg.Info)
                    n.setRaftRole(gROLE_RAFT_FOLLOWER)
                    result = gMSG_RAFT_RESPONSE
                }
            }
        }
        tconn.Close()
    }
    n.sendMsg(conn, result, "")
}

// 处理split brains问题
func (n *Node) onMsgRaftSplitBrainsUnset(conn net.Conn, msg *Msg) {
    glog.Println("split brains occurred, remove node:", msg.Info.Name)
    n.Peers.Remove(msg.Info.Id)
}

// 上线通知
func (n *Node) onMsgRaftHi(conn net.Conn, msg *Msg) {
    n.sendMsg(conn, gMSG_RAFT_HI2, "")
}

// 心跳保持
func (n *Node) onMsgRaftHeartbeat(conn net.Conn, msg *Msg) {
    n.updateElectionDeadline()
    if n.checkConnInLocalNode(conn) {
        n.Peers.Remove(msg.Info.Id)
        conn.Close()
        return
    }
    result := gMSG_RAFT_HEARTBEAT

    if n.getRaftRole() == gROLE_RAFT_LEADER {
        // 如果是两个leader相互心跳，表示两个leader是连通的，这时根据算法算出一个leader即可
        // 需要同时对比日志信息及选举比分
        if n.getLogCount() > msg.Info.LogCount && n.getLastLogId() > msg.Info.LastLogId {
            result = gMSG_RAFT_I_AM_LEADER
        } else if n.getLogCount() == msg.Info.LogCount && n.getLastLogId() == msg.Info.LastLogId {
            if n.getScoreCount() > msg.Info.ScoreCount {
                result = gMSG_RAFT_I_AM_LEADER
            } else if n.getScoreCount() == msg.Info.ScoreCount {
                if n.getScore() > msg.Info.Score {
                    result = gMSG_RAFT_I_AM_LEADER
                } else if n.getScore() == msg.Info.Score {
                    // 极少数情况, 这时采用随机策略
                    if grand.Rand(0, 1) == 1 {
                        result = gMSG_RAFT_I_AM_LEADER
                    }
                }
            }
        }
        if result == gMSG_RAFT_HEARTBEAT {
            n.setLeader(&msg.Info)
            n.setRaftRole(gROLE_RAFT_FOLLOWER)
        }
    } else if n.getLeader() == nil {
        // 如果没有leader，那么设置leader
        n.setLeader(&msg.Info)
        n.setRaftRole(gROLE_RAFT_FOLLOWER)
    } else {
        // 脑裂问题，集群节点规划或者网络异常造成
        // 1、两个leader无法相互通信，那么两个leader处于不同的两个网络，因此需要将其中一个网络中的该follower剔除掉，只保留其在一个网络中
        // 2、两个leader可以相互通信，那么两个leader处于相同的网络，于是将两个leader相互比较，最终留下一个作为leader，另外一个作为follower
        if n.getLeader().Id != msg.Info.Id {
            glog.Println("split brains occurred:", n.getLeader().Name, "and", msg.Info.Name)
            leaderConn := n.getConn(n.getLeader().Ip, gPORT_RAFT)
            if leaderConn != nil {
                if n.sendMsg(leaderConn, gMSG_RAFT_SPLIT_BRAINS_CHECK, msg.Info.Ip) == nil {
                    rmsg := n.receiveMsg(leaderConn)
                    if rmsg != nil {
                        if msg.Head == gMSG_RAFT_SPLIT_BRAINS_UNSET {
                            result = gMSG_RAFT_SPLIT_BRAINS_UNSET
                            n.updatePeerStatus(msg.Info.Id, gSTATUS_DEAD)
                        } else {
                            n.setLeader(&msg.Info)
                        }
                    }
                }
                leaderConn.Close()
            }
        }
    }
    n.sendMsg(conn, result, "")
}

// 选举比分获取
func (n *Node) onMsgRaftScoreRequest(conn net.Conn, msg *Msg) {
    if n.getRaftRole() == gROLE_RAFT_LEADER {
        n.sendMsg(conn, gMSG_RAFT_I_AM_LEADER, "")
    } else {
        n.sendMsg(conn, gMSG_RAFT_RESPONSE, "")
    }
}

// 选举比分对比
// 注意：这里除了比分选举，还需要判断数据一致性的对比
func (n *Node) onMsgRaftScoreCompareRequest(conn net.Conn, msg *Msg) {
    result := gMSG_RAFT_SCORE_COMPARE_SUCCESS
    if n.getRaftRole() == gROLE_RAFT_LEADER {
        result = gMSG_RAFT_I_AM_LEADER
    } else {
        // 需要同时对比日志信息和比分
        if n.getLogCount() > msg.Info.LogCount && n.getLastLogId() > msg.Info.LastLogId {
            result = gMSG_RAFT_SCORE_COMPARE_FAILURE
        } else if n.getLogCount() == msg.Info.LogCount && n.getLastLogId() == msg.Info.LastLogId {
            if n.getScoreCount() > msg.Info.ScoreCount {
                result = gMSG_RAFT_SCORE_COMPARE_FAILURE
            } else if n.getScoreCount() == msg.Info.ScoreCount {
                if n.getScore() > msg.Info.Score {
                    result = gMSG_RAFT_SCORE_COMPARE_FAILURE
                } else if n.getScore() == msg.Info.Score {
                    // 极少数情况, 这时采用随机策略
                    if grand.Rand(0, 1) == 1 {
                        result = gMSG_RAFT_SCORE_COMPARE_FAILURE
                    }
                }
            }
        }
    }
    if result == gMSG_RAFT_SCORE_COMPARE_SUCCESS {
        n.setLeader(&msg.Info)
        n.setRaftRole(gROLE_RAFT_FOLLOWER)
    }
    n.sendMsg(conn, result, "")
}

// 新增节点,通过IP添加
func (n *Node) onMsgApiPeersAdd(conn net.Conn, msg *Msg) {
    list := make([]string, 0)
    gjson.DecodeTo(msg.Body, &list)
    if list != nil && len(list) > 0 {
        for _, ip := range list {
            if n.Peers.Contains(ip) {
                continue
            }
            n.updatePeerInfo(NodeInfo{Id: ip, Ip: ip})
        }
    }
    n.sendMsg(conn, gMSG_RAFT_RESPONSE, "")
}

// 删除节点，目前通过IP删除，效率较低
func (n *Node) onMsgApiPeersRemove(conn net.Conn, msg *Msg) {
    list := make([]string, 0)
    gjson.DecodeTo(msg.Body, &list)
    if list != nil && len(list) > 0 {
        peers := n.Peers.Values()
        for _, ip := range list {
            // glog.Println("removing peer:", ip)
            for _, v := range peers {
                info := v.(NodeInfo)
                if ip == info.Ip {
                    n.Peers.Remove(ip)
                    break;
                }
            }
        }
    }
    n.sendMsg(conn, gMSG_RAFT_RESPONSE, "")
}

package gluster

import (
    "net"
    "g/util/grand"
    "g/encoding/gjson"
    "g/os/glog"
    "g/util/gtime"
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
        case gMSG_API_PEERS_INFO:               n.onMsgApiPeersInfo(conn, msg)
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
    ip      := n.Ip
    checkip := msg.Body
    if !n.Peers.Contains(checkip) {
        tconn := n.getConn(checkip, gPORT_RAFT)
        if tconn != nil {
            if n.sendMsg(tconn, gMSG_RAFT_HI, "") == nil {
                rmsg := n.receiveMsg(tconn)
                if rmsg != nil {
                    n.updatePeerInfo(rmsg.Info)
                    if n.getLastLogId() < msg.Info.LastLogId {
                        ip = checkip
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)
                    }
                }
            }
            tconn.Close()
        } else {
            ip = ""
        }
    }
    n.sendMsg(conn, gMSG_RAFT_RESPONSE, ip)
}

// 处理split brains问题
func (n *Node) onMsgRaftSplitBrainsUnset(conn net.Conn, msg *Msg) {
    n.Peers.Remove(msg.Info.Ip)
}

// 上线通知
func (n *Node) onMsgRaftHi(conn net.Conn, msg *Msg) {
    n.sendMsg(conn, gMSG_RAFT_HI2, "")
    //glog.Println("add peer:", fromip, "to", n.Ip, ", remote", conn.RemoteAddr(), ", local", conn.LocalAddr())
}

// 心跳保持
func (n *Node) onMsgRaftHeartbeat(conn net.Conn, msg *Msg) {
    n.updateElectionDeadline()
    result := gMSG_RAFT_HEARTBEAT
    if n.getRole() == gROLE_LEADER {
        // 如果是两个leader相互心跳，表示两个leader是连通的，这时根据算法算出一个leader即可
        if n.getScoreCount() > msg.Info.ScoreCount {
            result = gMSG_RAFT_I_AM_LEADER
        } else if n.getScoreCount() == msg.Info.ScoreCount {
            if n.getScore() > msg.Info.Score {
                result = gMSG_RAFT_I_AM_LEADER
            } else if n.getScore() == msg.Info.Score {
                // 极少数情况会出现两个节点ScoreCount和Score都相等的情况, 这个时候采用随机策略
                if grand.Rand(0, 1) == 1 {
                    result = gMSG_RAFT_I_AM_LEADER
                }
            }
        }
        if result == gMSG_RAFT_HEARTBEAT {
            n.setLeader(msg.Info.Ip)
            n.setRole(gROLE_FOLLOWER)
        }
    } else if n.getLeader() == "" {
        // 如果没有leader，那么设置leader
        n.setLeader(msg.Info.Ip)
        n.setRole(gROLE_FOLLOWER)
    } else {
        // 脑裂问题，一个节点处于两个网路中，并且两个网络的leader无法相互通信，会引起数据一致性问题
        if n.getLeader() != msg.Info.Ip {
            glog.Println("split brains occurred:", n.getLeader(), "and", msg.Info.Ip)
            leaderConn := n.getConn(n.getLeader(), gPORT_RAFT)
            if leaderConn != nil {
                if n.sendMsg(leaderConn, gMSG_RAFT_SPLIT_BRAINS_CHECK, msg.Info.Ip) == nil {
                    rmsg := n.receiveMsg(leaderConn)
                    if rmsg != nil {
                        if msg.Body == "" {
                            // 请求返回空表示该节点不与对方leader为一个集群，或者网络分区，关闭联系方式
                            result = gMSG_RAFT_SPLIT_BRAINS_UNSET
                            n.updatePeerStatus(msg.Info.Ip, gSTATUS_DEAD)
                        } else if n.getLeader() != msg.Body {
                            n.setLeader(msg.Body)
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
    if n.getRole() == gROLE_LEADER {
        n.sendMsg(conn, gMSG_RAFT_I_AM_LEADER, "")
    } else {
        n.sendMsg(conn, gMSG_RAFT_RESPONSE, "")
    }
}

// 选举比分对比
func (n *Node) onMsgRaftScoreCompareRequest(conn net.Conn, msg *Msg) {
    result := gMSG_RAFT_SCORE_COMPARE_SUCCESS
    if n.getRole() == gROLE_LEADER {
        result = gMSG_RAFT_I_AM_LEADER
    } else {
        if n.getScoreCount() > msg.Info.ScoreCount {
            result = gMSG_RAFT_SCORE_COMPARE_FAILURE
        } else if n.getScoreCount() == msg.Info.ScoreCount {
            if n.getScore() > msg.Info.Score {
                result = gMSG_RAFT_SCORE_COMPARE_FAILURE
            } else if n.getScore() == msg.Info.Score {
                // 极少数情况会出现两个节点ScoreCount和Score都相等的情况, 这个时候采用随机策略
                if grand.Rand(0, 1) == 1 {
                    result = gMSG_RAFT_SCORE_COMPARE_FAILURE
                }
            }
        } else {
            result = gMSG_RAFT_SCORE_COMPARE_SUCCESS
        }
    }
    if result == gMSG_RAFT_SCORE_COMPARE_SUCCESS {
        n.setLeader(msg.Info.Ip)
        n.setRole(gROLE_FOLLOWER)
    }
    n.sendMsg(conn, result, "")
}

// 节点信息查询
func (n *Node) onMsgApiPeersInfo(conn net.Conn, msg *Msg) {
    list := make([]NodeInfo, 0)
    list  = append(list, *n.getNodeInfo())
    for _, v := range n.Peers.Values() {
        list = append(list, v.(NodeInfo))
    }
    n.sendMsg(conn, gMSG_API_PEERS_INFO, *gjson.Encode(list))
}

// 新增节点
func (n *Node) onMsgApiPeersAdd(conn net.Conn, msg *Msg) {
    list := make([]string, 0)
    gjson.DecodeTo(&(msg.Body), &list)
    if list != nil && len(list) > 0 {
        for _, ip := range list {
            if n.Peers.Contains(ip) {
                continue
            }
            // glog.Println("adding peer:", ip)
            go func(ip string) {
                conn := n.getConn(ip, gPORT_RAFT)
                if conn != nil {
                    n.sendMsg(conn, gMSG_RAFT_HI, "")
                    rmsg := n.receiveMsg(conn)
                    if rmsg != nil && rmsg.Head == gMSG_RAFT_HI2 {
                        n.updatePeerInfo(rmsg.Info)
                    }
                }
                // 判断是否添加成功，如果没有，那么添加一个默认的信息
                if !n.Peers.Contains(ip) {
                    info       := NodeInfo{}
                    info.Ip     = ip
                    info.Status = gSTATUS_DEAD
                    info.LastActiveTime = gtime.Millisecond()
                    n.updatePeerInfo(info)
                }
            }(ip)
        }
    }
    n.sendMsg(conn, gMSG_RAFT_RESPONSE, "")
}

// 删除节点
func (n *Node) onMsgApiPeersRemove(conn net.Conn, msg *Msg) {
    list := make([]string, 0)
    gjson.DecodeTo(&(msg.Body), &list)
    if list != nil && len(list) > 0 {
        for _, ip := range list {
            // glog.Println("removing peer:", ip)
            n.Peers.Remove(ip)
        }
    }
    n.sendMsg(conn, gMSG_RAFT_RESPONSE, "")
}

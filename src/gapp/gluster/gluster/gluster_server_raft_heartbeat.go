package gluster

import (
    "net"
    "time"
    "g/core/types/gset"
    "g/os/glog"
)

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
                    // 注意一个节点添加的时候会给定一个初始化的活跃时间(添加时间)，因此该字段不会为0
                    //if info.LastActiveTime != 0 && gtime.Millisecond() - info.LastActiveTime > 3 * 86400 * 1000 {
                    //    glog.Println(info.Ip, "was dead over 3 days, removing from peers")
                    //    n.Peers.Remove(info.Ip)
                    //}
                    continue
                }
                conns.Add(info.Ip)
                go func(ip string, conn net.Conn) {
                    defer func() {
                        conn.Close()
                        conns.Remove(ip)
                    }()
                    for {
                        // 如果当前节点不再是leader，或者节点表中已经删除该节点信息
                        if n.getRole() != gROLE_LEADER || !n.Peers.Contains(ip){
                            return
                        }
                        if n.sendMsg(conn, gMSG_RAFT_HEARTBEAT, "") != nil {
                            n.updatePeerStatus(ip, gSTATUS_DEAD)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            n.updatePeerStatus(ip, gSTATUS_DEAD)
                            return
                        } else {
                            //glog.Println("receive heartbeat back from", ip)
                            // 更新节点信息
                            n.updatePeerInfo(msg.Info)
                            switch msg.Head {
                                case gMSG_RAFT_I_AM_LEADER:
                                    glog.Println("two leader occured, set", ip, "as my leader, done heartbeating")
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

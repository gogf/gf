package gluster

import (
    "time"
    "g/core/types/gset"
    "g/os/glog"
)

// 通过心跳维持集群统治，如果心跳不及时，那么选民会重新进入选举流程
func (n *Node) heartbeatHandler() {
    // 存储已经保持心跳的节点
    conns := gset.NewStringSet()
    for {
        if n.getRaftRole() == gROLE_RAFT_LEADER {
            for _, v := range n.Peers.Values() {
                info := v.(NodeInfo)
                if conns.Contains(info.Id) {
                    continue
                }
                go func(info *NodeInfo) {
                    conns.Add(info.Id)
                    defer conns.Remove(info.Id)
                    conn := n.getConn(info.Ip, gPORT_RAFT)
                    if conn == nil {
                        n.updatePeerStatus(info.Id, gSTATUS_DEAD)
                        return
                    }
                    defer conn.Close()
                    // 如果是本地同一节点通信，那么移除掉
                    if n.checkConnInLocalNode(conn) {
                        n.Peers.Remove(info.Id)
                        return
                    }
                    for {
                        // 如果当前节点不再是leader，或者节点表中已经删除该节点信息
                        if n.getRaftRole() != gROLE_RAFT_LEADER || !n.Peers.Contains(info.Id){
                            return
                        }
                        if n.sendMsg(conn, gMSG_RAFT_HEARTBEAT, "") != nil {
                            n.updatePeerStatus(info.Id, gSTATUS_DEAD)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            n.updatePeerStatus(info.Id, gSTATUS_DEAD)
                            return
                        } else {
                            //glog.Println("receive heartbeat back from", ip)
                            // 更新节点信息
                            n.updatePeerInfo(msg.Info)
                            switch msg.Head {
                                case gMSG_RAFT_I_AM_LEADER:
                                    glog.Println("two leader occured, set", msg.Info.Name, "as my leader, done heartbeating")
                                    n.setLeader(&(msg.Info))
                                    n.setRaftRole(gROLE_RAFT_FOLLOWER)

                                default:
                                    time.Sleep(gELECTION_TIMEOUT_HEARTBEAT * time.Millisecond)
                            }
                        }
                    }
                }(&info)
            }
        }
        time.Sleep(gELECTION_TIMEOUT_HEARTBEAT * time.Millisecond)
    }
}

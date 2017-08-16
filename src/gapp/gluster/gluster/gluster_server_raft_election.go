package gluster

import (
    "net"
    "sync"
    "time"
    "g/util/gtime"
    "g/core/types/gmap"
    "g/os/glog"
)

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
                glog.Println("only one node in this cluster, so i'll be the leader")
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
    glog.Println(n.Ip + ":", "begin new election")
    conns := gmap.NewStringInterfaceMap()
    // 请求比分，获取比分数据
    for _, v := range n.Peers.Values() {
        info := v.(NodeInfo)
        if info.Status != gSTATUS_ALIVE {
            continue
        }
        wg.Add(1)
        go func(ip string) {
            defer wg.Done()
            if n.getLeader() != "" || n.getRole() != gROLE_CANDIDATE {
                return
            }
            stime := time.Now().UnixNano()
            conn  := n.getConnFromPool(ip, gPORT_RAFT, conns)
            if conn == nil {
                n.updatePeerStatus(ip, gSTATUS_DEAD)
                return
            }
            defer conn.Close()
            if err := n.sendMsg(conn, gMSG_RAFT_SCORE_REQUEST, ""); err != nil {
                glog.Println(err)
                return
            }
            msg := n.receiveMsg(conn)
            if msg != nil {
                if n.getLeader() != "" || n.getRole() != gROLE_CANDIDATE {
                    return
                }
                switch msg.Head {
                    case gMSG_RAFT_I_AM_LEADER:
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)

                    case gMSG_RAFT_RESPONSE:
                        etime := time.Now().UnixNano()
                        score := etime - stime
                        n.addScore(score)
                        n.addScoreCount()
                }
            } else {
                n.updatePeerStatus(ip, gSTATUS_DEAD)
            }
        }(info.Ip)
    }
    wg.Wait()

    // 如果在计算比分的过程中发现了leader，那么不再继续比分，退出选举
    if n.getLeader() != "" {
        return;
    }

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
            if err := n.sendMsg(conn, gMSG_RAFT_SCORE_COMPARE_REQUEST, ""); err != nil {
                glog.Println(err)
                return
            }
            msg := n.receiveMsg(conn)
            if msg != nil {
                if n.getLeader() != "" || n.getRole() != gROLE_CANDIDATE {
                    return
                }
                switch msg.Head {
                    case gMSG_RAFT_I_AM_LEADER:
                        glog.Println("score comparison: get leader from", ip)
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)

                    case gMSG_RAFT_SCORE_COMPARE_FAILURE:
                        glog.Println("score comparison: get failure from", ip)
                        n.setLeader(ip)
                        n.setRole(gROLE_FOLLOWER)

                    case gMSG_RAFT_SCORE_COMPARE_SUCCESS:
                        glog.Println("score comparison: get success from", ip)
                }
            }
        }(info.Ip, conn)
    }
    wg.Wait()

    // 如果peers中的节点均没有条件满足leader，那么选举自身为leader
    if n.getRole() != gROLE_FOLLOWER {
        glog.Println(n.Ip + ":", "I've won this score comparison")
        n.setRole(gROLE_LEADER)
        n.setLeader(n.Ip)
    }
}


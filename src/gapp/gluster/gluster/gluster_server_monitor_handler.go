package gluster

import (
    "g/core/types/gmap"
    "net"
)

// monitor通信接口回调函数
func (n *Node) monitorTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil || msg.Info.Group != n.Group {
        conn.Close()
        return
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

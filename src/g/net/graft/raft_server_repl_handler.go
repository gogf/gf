// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package graft

import (
    "net"
    "g/encoding/gjson"
    "time"
    "log"
)

// 集群数据同步接口回调函数
func (n *Node) replTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil {
        //log.Println("receive nil")
        conn.Close()
        return
    }
    switch msg.Head {
        case gMSG_REPL_SET:         n.onMsgReplSet(conn, msg)
        case gMSG_REPL_REMOVE:      n.onMsgReplRemove(conn, msg)
        case gMSG_REPL_HEARTBEAT:   n.onMsgReplHeartbeat(conn, msg)
        case gMSG_REPL_UPDATE:      n.onMsgReplUpdate(conn, msg)
    }
    conn.Close()
}

// kv删除
func (n *Node) onMsgReplRemove(conn net.Conn, msg *Msg) {
    n.onMsgReplSet(conn, msg)
}

// kv设置
func (n *Node) onMsgReplSet(conn net.Conn, msg *Msg) {
    n.setStatusInReplication(true)
    if n.getRole() == gROLE_LEADER {
        var items interface{}
        if gjson.DecodeTo(&msg.Body, &items) == nil {
            var entry = LogEntry {
                Id    : time.Now().UnixNano(),
                Act   : msg.Head,
                Items : items,
            }
            n.LogList.PushFront(entry)
            n.LogChan <- entry
        }
    } else {
        var entry LogEntry
        gjson.DecodeTo(&msg.Body, &entry)
        n.saveLogEntry(entry)
    }
    n.setStatusInReplication(false)
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// 心跳
func (n *Node) onMsgReplHeartbeat(conn net.Conn, msg *Msg) {
    result := gMSG_REPL_HEARTBEAT
    //log.Println("heartbeat:", n.getLastLogId(), msg.Info.LastLogId, n.getStatusInReplication())
    if n.getLastLogId() < msg.Info.LastLogId {
        if !n.getStatusInReplication() {
            result = gMSG_REPL_NEED_UPDATE_FOLLOWER
        }
    } else if n.getLastLogId() > msg.Info.LastLogId {
        if !n.getStatusInReplication() {
            result = gMSG_REPL_NEED_UPDATE_LEADER
        }
    }
    if result == gMSG_REPL_NEED_UPDATE_LEADER {
        n.sendMsg(conn, result, *gjson.Encode(n.getLogEntriesByLastLogId(msg.Info.LastLogId)))
    } else {
        n.sendMsg(conn, result, "")
    }
    n.replTcpHandler(conn)
}

// 数据同步，更新本地数据
func (n *Node) onMsgReplUpdate(conn net.Conn, msg *Msg) {
    log.Println("receive data replication update")
    if n.getLastLogId() < msg.Info.LastLogId {
        if !n.getStatusInReplication() {
            n.updateFromLogEntriesJson(&msg.Body)
        }
    }
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

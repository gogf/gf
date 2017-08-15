// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package gluster

import (
    "net"
    "g/encoding/gjson"
    "g/core/types/gmap"
    "g/util/gtime"
    "g/os/glog"
)

// 集群数据同步接口回调函数
func (n *Node) replTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil {
        //glog.Println("receive nil")
        conn.Close()
        return
    }
    switch msg.Head {
        case gMSG_REPL_SET:                 n.onMsgReplSet(conn, msg)
        case gMSG_REPL_REMOVE:              n.onMsgReplRemove(conn, msg)
        case gMSG_REPL_HEARTBEAT:           n.onMsgReplHeartbeat(conn, msg)
        case gMSG_REPL_COMPLETELY_UPDATE:   n.onMsgReplUpdate(conn, msg)
        case gMSG_REPL_INCREMENTAL_UPDATE:  n.onMsgReplUpdate(conn, msg)
        case gMSG_API_SERVICE_SET:          n.onMsgServiceSet(conn, msg)
        case gMSG_API_SERVICE_REMOVE:       n.onMsgServiceRemove(conn, msg)
    }
    //这里不用自动关闭链接，由于链接有读取超时，当一段时间没有数据时会自动关闭
    n.replTcpHandler(conn)
}

// service删除
func (n *Node) onMsgServiceRemove(conn net.Conn, msg *Msg) {
    list := make([]interface{}, 0)
    if gjson.DecodeTo(&msg.Body, &list) == nil {
        for _, name := range list {
            n.Service.Remove(name.(string))
            n.setLastServiceLogId(gtime.Microsecond())
        }
    }
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// service设置
func (n *Node) onMsgServiceSet(conn net.Conn, msg *Msg) {
    var service Service
    if gjson.DecodeTo(&msg.Body, &service) == nil {
        n.Service.Set(service.Name, service)
        n.setLastServiceLogId(gtime.Microsecond())
    }
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
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
                Id    : gtime.Microsecond(),
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

// 心跳响应
func (n *Node) onMsgReplHeartbeat(conn net.Conn, msg *Msg) {
    result := gMSG_REPL_HEARTBEAT
    //glog.Println("heartbeat:", n.getLastLogId(), msg.Info.LastLogId, n.getStatusInReplication())
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
        n.updateDataToRemoteNode(conn, msg)
    } else {
        n.sendMsg(conn, result, "")
    }
}

// 数据同步，更新本地数据
func (n *Node) onMsgReplUpdate(conn net.Conn, msg *Msg) {
    glog.Println("receive data replication update")
    n.updateDataFromRemoteNode(conn, msg)
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// 从目标节点同步数据，采用增量+全量模式
func (n *Node) updateDataFromRemoteNode(conn net.Conn, msg *Msg) {
    if msg.Head == gMSG_REPL_INCREMENTAL_UPDATE {
        // 增量同步，LastLogId会根据保存的LogEntry自动更新
        if n.getLastLogId() < msg.Info.LastLogId {
            if !n.getStatusInReplication() {
                n.updateFromLogEntriesJson(&msg.Body)
            }
        }
    } else {
        // 全量同步，完整的kv数据覆盖，因此需要设置LastLogId
        m   := make(map[string]string)
        err := gjson.DecodeTo(&(msg.Body), &m)
        if err == nil {
            newm := gmap.NewStringStringMap()
            newm.BatchSet(m)
            n.setKVMap(newm)
            n.setLastLogId(msg.Info.LastLogId)
        } else {
            glog.Println(err)
        }
    }
}

// 同步数据到目标节点，采用增量+全量模式
func (n *Node) updateDataToRemoteNode(conn net.Conn, msg *Msg) {
    n.setStatusInReplication(true)
    glog.Println("send data replication update from", n.Ip, "to", msg.Info.Ip)
    // 首先进行增量同步
    if err := n.sendMsg(conn, gMSG_REPL_INCREMENTAL_UPDATE, *gjson.Encode(n.getLogEntriesByLastLogId(msg.Info.LastLogId))); err != nil {
        glog.Println(err)
        return
    }
    rmsg := n.receiveMsg(conn)
    if rmsg != nil {
        // 如果增量同步失败则采用全量同步
        if n.getLastLogId() > rmsg.Info.LastLogId {
            glog.Println(rmsg.Info.Ip + ":", "incremental update failed, now try completely update")
            if err := n.sendMsg(conn, gMSG_REPL_COMPLETELY_UPDATE, *gjson.Encode(*n.KVMap.Clone())); err != nil {
                glog.Println(err)
                return
            }
            n.receiveMsg(conn)
        }
        glog.Println(rmsg.Info.Ip + ":", "data replication update done")
    }
    n.setStatusInReplication(false)
}


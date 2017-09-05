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
    "sync"
)

// 集群数据同步接口回调函数
func (n *Node) replTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil || msg.Info.Group != n.Group {
        //glog.Println("receive nil")
        conn.Close()
        return
    }
    switch msg.Head {
        case gMSG_REPL_DATA_SET:                    n.onMsgReplDataSet(conn, msg)
        case gMSG_REPL_DATA_REMOVE:                 n.onMsgReplDataRemove(conn, msg)
        case gMSG_REPL_HEARTBEAT:                   n.onMsgReplHeartbeat(conn, msg)
        case gMSG_REPL_PEERS_UPDATE:                n.onMsgPeersUpdate(conn, msg)
        case gMSG_REPL_INCREMENTAL_UPDATE:          n.onMsgReplUpdate(conn, msg)
        case gMSG_REPL_COMPLETELY_UPDATE:           n.onMsgReplUpdate(conn, msg)
        case gMSG_REPL_CONFIG_FROM_FOLLOWER     :   n.onMsgConfigFromFollower(conn, msg)
        case gMSG_REPL_SERVICE_COMPLETELY_UPDATE:   n.onMsgServiceCompletelyUpdate(conn, msg)
        case gMSG_API_SERVICE_SET:                  n.onMsgServiceSet(conn, msg)
        case gMSG_API_SERVICE_REMOVE:               n.onMsgServiceRemove(conn, msg)
    }
    //这里不用自动关闭链接，由于链接有读取超时，当一段时间没有数据时会自动关闭
    n.replTcpHandler(conn)
}

// Follower->Leader的配置同步
func (n *Node) onMsgConfigFromFollower(conn net.Conn, msg *Msg) {
    //glog.Println("config replication from", msg.Info.Name)
    j := gjson.DecodeToJson(msg.Body)
    if j != nil {
        // 初始化节点列表，包含自定义的所需添加的服务器IP或者域名列表
        peers := j.GetArray("Peers")
        if peers != nil {
            for _, v := range peers {
                ip := v.(string)
                if ip == n.Ip || n.Peers.Contains(ip){
                    continue
                }
                go func(ip string) {
                    if !n.sayHi(ip) {
                        n.updatePeerInfo(NodeInfo{Id: ip, Ip: ip})
                    }
                }(ip)
            }
        }
    }
    conn.Close()
    //glog.Println("config replication from", msg.Info.Name, "done")
}

// Peers信息更新
func (n *Node) onMsgPeersUpdate(conn net.Conn, msg *Msg) {
    //glog.Println("receive peers update", msg.Body)
    m := make([]NodeInfo, 0)
    if gjson.DecodeTo(msg.Body, &m) == nil {
        for _, v := range m {
            if v.Id != n.Id {
                n.updatePeerInfo(v)
            } else {
                n.setIp(v.Ip)
            }
        }
    }
    conn.Close()
}

// 心跳消息提交的完整更新消息
func (n *Node) onMsgServiceCompletelyUpdate(conn net.Conn, msg *Msg) {
    n.updateServiceFromRemoteNode(conn, msg)
}

// Service删除
func (n *Node) onMsgServiceRemove(conn net.Conn, msg *Msg) {
    list := make([]interface{}, 0)
    if gjson.DecodeTo(msg.Body, &list) == nil {
        updated := false
        for _, name := range list {
            if n.Service.Contains(name.(string)) {
                n.Service.Remove(name.(string))
                updated = true
            }
        }
        if updated {
            n.setLastServiceLogId(gtime.Microsecond())
        }
    }
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// Service设置
func (n *Node) onMsgServiceSet(conn net.Conn, msg *Msg) {
    var st ServiceStruct
    if gjson.DecodeTo(msg.Body, &st) == nil {
        n.Service.Set(st.Name, *n.serviceSructToService(&st))
        n.ServiceForApi.Set(st.Name, st)
        n.setLastServiceLogId(gtime.Microsecond())
    }
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// kv删除
func (n *Node) onMsgReplDataRemove(conn net.Conn, msg *Msg) {
    n.onMsgReplDataSet(conn, msg)
}

// kv设置，最终一致性
func (n *Node) onMsgReplDataSet(conn net.Conn, msg *Msg) {
    n.setStatusInReplication(true)
    if n.getRaftRole() == gROLE_RAFT_LEADER {
        var items interface{}
        if gjson.DecodeTo(msg.Body, &items) == nil {
            var entry = LogEntry {
                Id    : gtime.Microsecond(),
                Act   : msg.Head,
                Items : items,
            }
            n.LogList.PushFront(entry)
            n.saveLogEntry(entry)
            // 这里不做主动通知数据同步，而是依靠心跳检测时的单线程数据同步
            // 并且为保证客户端能够及时相应（例如在写入请求的下一次获取请求将一定能够获取到最新的数据），
            // 因此，请求端应当在leader返回成功后，同时将该数据写入到本地
            // n.sendLogEntryToPeers(entry)
        }
    } else {
        var entry LogEntry
        gjson.DecodeTo(msg.Body, &entry)
        n.saveLogEntry(entry)
    }
    n.setStatusInReplication(false)
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// 发送数据操作到其他节点,为保证数据的强一致性，所有节点返回结果后，才算成功
// 只要数据请求完整流程执行完毕，即使其中几个节点失败也不影响，因为有另外的数据同步方式进行进一步的数据一致性保证
func (n *Node) sendLogEntryToPeers(entry LogEntry) {
    var wg sync.WaitGroup
    n.setStatusInReplication(true)
    // 异步并发发送数据操作请求到其他节点
    //glog.Println("sending log entry", entry)
    for _, v := range n.Peers.Values() {
        info := v.(NodeInfo)
        if info.Status != gSTATUS_ALIVE {
            continue
        }
        wg.Add(1)
        go func(info *NodeInfo, entry LogEntry) {
            defer wg.Done()
            conn := n.getConn(info.Ip, gPORT_REPL)
            if conn == nil {
                return
            }
            defer conn.Close()
            if n.sendMsg(conn, entry.Act, gjson.Encode(entry)) == nil {
                n.receiveMsg(conn)
            }
        }(&info, entry)
    }
    wg.Wait()
    n.setStatusInReplication(false)
}

// 心跳响应
func (n *Node) onMsgReplHeartbeat(conn net.Conn, msg *Msg) {
    result := gMSG_REPL_HEARTBEAT
    //glog.Println("heartbeat:", n.getLastLogId(), msg.Info.LastLogId, n.getStatusInReplication())
    // 日志检测同步
    lastLogId := n.getLastLogId()
    if lastLogId < msg.Info.LastLogId {
        if !n.getStatusInReplication() {
            result = gMSG_REPL_NEED_UPDATE_FOLLOWER
        }
    } else if lastLogId > msg.Info.LastLogId {
        if !n.getStatusInReplication() {
            result = gMSG_REPL_NEED_UPDATE_LEADER
        }
    } else {
        // service同步检测
        lastServiceLogId := n.getLastServiceLogId()
        if lastServiceLogId < msg.Info.LastServiceLogId {
            result = gMSG_REPL_SERVICE_NEED_UPDATE_FOLLOWER
        } else if lastServiceLogId > msg.Info.LastServiceLogId {
            result = gMSG_REPL_SERVICE_NEED_UPDATE_LEADER
        }
    }
    switch result {
        case gMSG_REPL_NEED_UPDATE_LEADER:          n.updateDataToRemoteNode(conn, msg)
        case gMSG_REPL_SERVICE_NEED_UPDATE_LEADER:  n.updateServiceToRemoteNode(conn, msg)
        default:
            n.sendMsg(conn, result, "")
    }
}

// 数据同步，更新本地数据
func (n *Node) onMsgReplUpdate(conn net.Conn, msg *Msg) {
    //glog.Println("receive data replication update from", msg.Info.Name)
    n.updateDataFromRemoteNode(conn, msg)
    n.sendMsg(conn, gMSG_REPL_RESPONSE, "")
}

// 保存日志数据
func (n *Node) saveLogEntry(entry LogEntry) {
    lastLogId := n.getLastLogId()
    if entry.Id < lastLogId {
        glog.Printf("expired log entry, received:%v, current:%v\n", entry.Id, lastLogId)
        return
    }
    switch entry.Act {
        case gMSG_REPL_DATA_SET:
            //glog.Println("setting log entry", entry)
            for k, v := range entry.Items.(map[string]interface{}) {
                n.DataMap.Set(k, v.(string))
            }

        case gMSG_REPL_DATA_REMOVE:
            //glog.Println("removing log entry", entry)
            for _, v := range entry.Items.([]interface{}) {
                n.DataMap.Remove(v.(string))
            }

    }
    n.addLogCount()
    n.setLastLogId(entry.Id)
}

// 从目标节点同步数据，采用增量+全量模式
func (n *Node) updateDataFromRemoteNode(conn net.Conn, msg *Msg) {
    if msg.Head == gMSG_REPL_INCREMENTAL_UPDATE {
        // 增量同步，LogCount和LastLogId会根据保存的LogEntry自动更新
        if n.getLastLogId() < msg.Info.LastLogId {
            if !n.getStatusInReplication() {
                n.updateFromLogEntriesJson(msg.Body)
            }
        }
    } else {
        // 全量同步，完整的kv数据覆盖
        m   := make(map[string]string)
        err := gjson.DecodeTo(msg.Body, &m)
        if err == nil {
            newm := gmap.NewStringStringMap()
            newm.BatchSet(m)
            n.setDataMap(newm)
            n.setLogCount(msg.Info.LogCount)
            n.setLastLogId(msg.Info.LastLogId)
        } else {
            glog.Error(err)
        }
    }
}

// 同步数据到目标节点，采用增量+全量模式
func (n *Node) updateDataToRemoteNode(conn net.Conn, msg *Msg) {
    n.setStatusInReplication(true)
    defer n.setStatusInReplication(false)

    //glog.Println("send data replication update to", msg.Info.Name)
    // 首先进行增量同步
    updated := true
    list    := n.getLogEntriesByLastLogId(msg.Info.LastLogId)
    length  := len(list)
    if length > 0 && list[length - 1].Id == n.getLastLogId() && (msg.Info.LogCount + length) == n.getLogCount() {
        if err := n.sendMsg(conn, gMSG_REPL_INCREMENTAL_UPDATE, gjson.Encode(list)); err != nil {
            glog.Error(err)
            return
        }
        rmsg := n.receiveMsg(conn)
        if rmsg != nil {
            if n.getLastLogId() > rmsg.Info.LastLogId {
                //glog.Error(rmsg.Info.Name + ":", "incremental update failed, now try completely update")
                updated = false
            }
        }
    } else {
        updated = false
    }
    if !updated {
        // 如果增量同步失败，或者判断需要完整同步，则采用全量同步
        if err := n.sendMsg(conn, gMSG_REPL_COMPLETELY_UPDATE, gjson.Encode(*n.DataMap.Clone())); err != nil {
            glog.Error(err)
            return
        }
        n.receiveMsg(conn)
    }
}

// 从目标节点同步Service数据
func (n *Node) updateServiceFromRemoteNode(conn net.Conn, msg *Msg) {
    //glog.Println("receive service replication update from", msg.Info.Name)
    m   := make(map[string]ServiceStruct)
    err := gjson.DecodeTo(msg.Body, &m)
    if err == nil {
        newmForService    := gmap.NewStringInterfaceMap()
        newmForServiceApi := gmap.NewStringInterfaceMap()
        for k, v := range m {
            newmForService.Set(k, *n.serviceSructToService(&v))
            newmForServiceApi.Set(k, v)
        }
        n.setService(newmForService)
        n.setServiceForApi(newmForServiceApi)
        n.setLastServiceLogId(msg.Info.LastServiceLogId)
    } else {
        glog.Error(err)
    }
}

// 同步Service到目标节点
func (n *Node) updateServiceToRemoteNode(conn net.Conn, msg *Msg) {
    //glog.Println("send service replication update to", msg.Info.Name)
    if err := n.sendMsg(conn, gMSG_REPL_SERVICE_COMPLETELY_UPDATE, gjson.Encode(*n.ServiceForApi.Clone())); err != nil {
        glog.Error(err)
        return
    }
}


// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package graft

import (
    "net"
    "g/encoding/gjson"
    "time"
    "log"
    "g/core/types/gmap"
    "g/os/gfile"
    "g/core/types/gset"
    "g/util/gtime"
    "sync"
)

// 用以识别节点当前是否正在数据同步中
var isInReplication bool

// leader到其他节点的数据同步监听
func (n *Node) logAutoReplicationHandler() {
    var wg sync.WaitGroup
    // 初始化数据同步心跳检测
    go n.logAutoReplicationCheckHandler()
    // 进入循环监听日志事件
    for {
        select {
            case entry := <- n.LogChan:
                n.setStatusInReplication(true)
                log.Println("sending log entry", entry)
                for _, v := range n.Peers.Values() {
                    info := v.(NodeInfo)
                    if info.Status != gSTATUS_ALIVE {
                        continue
                    }
                    conn := n.getConn(info.Ip, gPORT_REPL)
                    if conn != nil {
                        wg.Add(1)
                        go func(conn net.Conn, entry LogEntry) {
                            if err := n.sendMsg(conn, entry.Act, *gjson.Encode(entry)); err != nil {
                                log.Println(err)
                                conn.Close()
                                wg.Done()
                                return
                            }
                            n.receiveMsg(conn)
                            wg.Done()
                        }(conn, entry)
                    }
                }
                wg.Wait()
                // 当所有节点的请求处理后，再保存数据到自身
                // 以便leader与follower之间的数据同步判断
                n.saveLogEntry(entry)
                n.setStatusInReplication(false)
        }
    }
}

// 日志自动同步检查，类似心跳
func (n *Node) logAutoReplicationCheckHandler() {
    conns := gset.NewStringSet()
    for {
        if n.getRole() == gROLE_LEADER {
            ips := n.Peers.Keys()
            for _, ip := range ips {
                if conns.Contains(ip) {
                    continue
                }
                conn := n.getConn(ip, gPORT_REPL)
                if conn == nil {
                    conns.Remove(ip)
                    continue
                }
                conns.Add(ip)
                go func(ip string, conn net.Conn) {
                    defer func() {
                        conn.Close()
                        conns.Remove(ip)
                    }()
                    for {
                        // 如果当前正在数据同步操作中，那么等待
                        for n.getStatusInReplication() {
                            time.Sleep(100 * time.Millisecond)
                        }
                        if n.getRole() != gROLE_LEADER {
                            return
                        }
                        //log.Println("sending replication heartbeat to", ip)
                        if err := n.sendMsg(conn, gMSG_REPL_HEARTBEAT, ""); err != nil {
                            log.Println(err)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            n.updatePeerStatus(ip, gSTATUS_DEAD)
                            return
                        } else {
                            switch msg.Head {
                                case gMSG_REPL_NEED_UPDATE_FOLLOWER:
                                    n.setStatusInReplication(true)
                                    log.Println("request data replication update to", ip)
                                    array := n.getLogEntriesByLastLogId(msg.Info.LastLogId)
                                    if err := n.sendMsg(conn, gMSG_REPL_UPDATE, *gjson.Encode(array)); err != nil {
                                        log.Println(err)
                                        return
                                    }
                                    msg := n.receiveMsg(conn)
                                    if msg != nil {
                                        log.Println("follower data replication update done")
                                    }
                                    n.setStatusInReplication(false)

                                case gMSG_REPL_NEED_UPDATE_LEADER:
                                    n.setStatusInReplication(true)
                                    log.Println("request data replication update from", ip)
                                    if n.updateFromLogEntriesJson(&msg.Body) == nil {
                                        log.Println("leader data replication update done")
                                    } else {
                                        log.Println("leader data replication update failed")
                                    }
                                    n.setStatusInReplication(false)

                                default:
                                    time.Sleep(gLOG_REPL_TIMEOUT_HEARTBEAT * time.Millisecond)
                            }
                        }
                    }
                }(ip, conn)
            }
        }
        time.Sleep(100 * time.Millisecond)
    }
}

// 保存日志数据
func (n *Node) saveLogEntry(entry LogEntry) {
    switch entry.Act {
        case gMSG_REPL_SET:
            log.Println("setting log entry", entry)
            for k, v := range entry.Items.(map[string]interface{}) {
                n.KVMap.Set(k, v.(string))
            }


        case gMSG_REPL_REMOVE:
            log.Println("removing log entry", entry)
            for _, v := range entry.Items.([]interface{}) {
                n.KVMap.Remove(v.(string))
            }

    }
    n.setLastLogId(entry.Id)
}

// 日志自动保存处理
func (n *Node) logAutoSavingHandler() {
    t := gtime.Millisecond()
    for {
        // 当日志列表的最新ID与保存的ID不相等，或者超过超时时间
        if n.getLastLogId() != n.getLastSavedLogId() || gtime.Millisecond() - t > gLOG_REPL_AUTOSAVE_INTERVAL {
            //log.Println("saving data to file")
            n.saveData()
            t = gtime.Millisecond()
        } else {
            time.Sleep(100 * time.Millisecond)
        }
    }
}

// 保存数据到磁盘
func (n *Node) saveData() {
    data := SaveInfo {
        LastLogId : n.getLastLogId(),
        Peers     : *n.Peers.Clone(),
        DataMap   : *n.KVMap.Clone(),
    }
    content := gjson.Encode(&data)
    gfile.PutContents(n.getDataFilePath(), *content)
    n.setLastSavedLogId(n.getLastLogId())
}

// 从物理化文件中恢复变量
func (n *Node) restoreData() {
    path := n.getDataFilePath()
    if gfile.Exists(path) {
        content := gfile.GetContents(path)
        if content != nil {
            //log.Println("initializing kvmap from data file")
            var data = SaveInfo {
                Peers   : make(map[string]interface{}),
                DataMap : make(map[string]string),
            }
            content := string(content)
            if gjson.DecodeTo(&content, &data) == nil {
                dataMap := gmap.NewStringStringMap()
                peerMap := gmap.NewStringInterfaceMap()
                infoMap := make(map[string]NodeInfo)
                gjson.DecodeTo(gjson.Encode(data.Peers), &infoMap)
                dataMap.BatchSet(data.DataMap)
                for k, v := range infoMap {
                    peerMap.Set(k, v)
                }
                n.setLastLogId(data.LastLogId)
                n.setLastSavedLogId(data.LastLogId)
                n.setPeers(peerMap)
                n.setKVMap(dataMap)
            }
        }
    } else {
        //log.Println("no data file found at", path)
    }
}

// 使用logentry数组更新本地的日志列表
func (n *Node) updateFromLogEntriesJson(jsonContent *string) error {
    array := make([]LogEntry, 0)
    err   := gjson.DecodeTo(jsonContent, &array)
    if err != nil {
        log.Println(err)
        return err
    }
    if array != nil && len(array) > 0 {
        for _, v := range array {
            if v.Id > n.getLastLogId() {
                n.saveLogEntry(v)
            }
        }
    }
    return nil
}

// 根据logid获取还未更新的日志列表
func (n *Node) getLogEntriesByLastLogId(id int64) []LogEntry {
    if n.getLastLogId() > id {
        array := make([]LogEntry, 0)
        n.LogList.RLock()
        l := n.LogList.L.Back()
        for l != nil {
            r := l.Value.(LogEntry)
            if r.Id > id {
                array = append(array, r)
            }
            l = l.Prev()
        }
        n.LogList.RUnlock()
        return array
    }
    return nil
}


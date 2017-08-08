// 数据同步需要注意的是：
// leader只有在通知完所有follower更新完数据之后，自身才会进行数据更新
// 因此leader
package graft

import (
    "net"
    "g/encoding/gjson"
    "time"
    "log"
    "sync"
    "g/core/types/gmap"
    "g/os/gfile"
    "g/core/types/gset"
    "g/util/gtime"
)

// 用以识别节点当前是否正在数据同步中
var isInReplication bool

// 集群数据同步接口回调函数
func (n *Node) replTcpHandler(conn net.Conn) {
    msg := n.receiveMsg(conn)
    if msg == nil {
        //log.Println("receive nil")
        conn.Close()
        return
    }
    switch msg.Head {
        case gMSG_HEAD_SET:
            fallthrough
        case gMSG_HEAD_REMOVE:
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
                    n.LogChan <- struct{}{}
                }
            } else {
                var entry LogEntry
                gjson.DecodeTo(&msg.Body, &entry)
                n.saveLogEntry(entry)
            }
            n.setStatusInReplication(false)
            n.sendMsg(conn, gMSG_HEAD_REPL_RESPONSE, "")

        // 数据同步自动检测
        case gMSG_HEAD_REPL_HEARTBEAT:
            result := gMSG_HEAD_REPL_HEARTBEAT
            //log.Println("heartbeat:", n.getLastLogId(), msg.Info.LastLogId, n.getStatusInReplication())
            if n.getLogCount() < msg.Info.LogCount {
                if !n.getStatusInReplication() {
                    result = gMSG_HEAD_REPL_NEED_UPDATE_FOLLOWER
                }
            } else if n.getLogCount() > msg.Info.LogCount {
                if !n.getStatusInReplication() {
                    result = gMSG_HEAD_REPL_NEED_UPDATE_LEADER
                }
            } else {
                if n.getLastLogId() < msg.Info.LastLogId {
                    if !n.getStatusInReplication() {
                        result = gMSG_HEAD_REPL_NEED_UPDATE_FOLLOWER
                    }
                } else if n.getLastLogId() > msg.Info.LastLogId {
                    if !n.getStatusInReplication() {
                        result = gMSG_HEAD_REPL_NEED_UPDATE_LEADER
                    }
                }
            }
            if result == gMSG_HEAD_REPL_NEED_UPDATE_LEADER {
                n.sendMsg(conn, result, *gjson.Encode(*n.KVMap.Clone()))
            } else {
                n.sendMsg(conn, result, "")
                n.replTcpHandler(conn)
            }

        // 数据完整同步更新
        case gMSG_HEAD_UPDATE:
            log.Println("receive data replication update")
            if n.getLastLogId() < msg.Info.LastLogId || n.getLogCount() < msg.Info.LogCount {
                if !n.getStatusInReplication() {
                    if n.updateFromDataMapJson(&msg.Body) == nil {
                        n.setLastLogId(msg.Info.LastLogId)
                        n.setLogCount(msg.Info.LogCount)
                    }
                }
            }
            n.sendMsg(conn, gMSG_HEAD_REPL_RESPONSE, "")
    }

    conn.Close()
}

// leader到其他节点的数据同步监听
func (n *Node) logAutoReplicationHandler() {
    var wg sync.WaitGroup
    // 初始化数据同步心跳检测
    go n.logAutoReplicationCheckHandler()

    for {
        select {
            case <- n.LogChan:
                n.setStatusInReplication(true)
                conns := gmap.NewStringInterfaceMap()
                for {
                    item := n.LogList.PopBack()
                    if item == nil {
                        break;
                    }
                    entry := item.(LogEntry)
                    log.Println("sending log entry", entry)
                    for _, v := range n.Peers.Values() {
                        info := v.(NodeInfo)
                        if info.Status != gSTATUS_ALIVE {
                            continue
                        }
                        conn := n.getConnFromPool(info.Ip, gPORT_REPL, conns)
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
                }
                n.setStatusInReplication(false)
                for _, c := range conns.M {
                    c.(net.Conn).Close()
                }
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
                    for {
                        // 如果当前正在数据同步操作中，那么等待
                        for n.getStatusInReplication() {
                            time.Sleep(100 * time.Millisecond)
                        }
                        if n.getRole() != gROLE_LEADER {
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        //log.Println("sending replication heartbeat to", ip)
                        if err := n.sendMsg(conn, gMSG_HEAD_REPL_HEARTBEAT, ""); err != nil {
                            log.Println(err)
                            conn.Close()
                            conns.Remove(ip)
                            return
                        }
                        msg := n.receiveMsg(conn)
                        if msg == nil {
                            n.updatePeerStatus(ip, gSTATUS_DEAD)
                            conns.Remove(ip)
                            conn.Close()
                            return
                        } else {
                            switch msg.Head {
                                case gMSG_HEAD_REPL_NEED_UPDATE_FOLLOWER:
                                    log.Println("request data replication update to", ip)
                                    if err := n.sendMsg(conn, gMSG_HEAD_UPDATE, *gjson.Encode(*n.KVMap.Clone())); err != nil {
                                        log.Println(err)
                                        conn.Close()
                                        conns.Remove(ip)
                                        return
                                    }
                                    msg := n.receiveMsg(conn)
                                    if msg != nil {
                                        log.Println("follower data replication update done")
                                    }

                                case gMSG_HEAD_REPL_NEED_UPDATE_LEADER:
                                    log.Println("request data replication update from", ip)
                                    if n.updateFromDataMapJson(&msg.Body) == nil {
                                        n.setLastLogId(msg.Info.LastLogId)
                                        log.Println("leader data replication update done")
                                    } else {
                                        log.Println("leader data replication update failed")
                                    }

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
        case gMSG_HEAD_SET:
            log.Println("setting log entry", entry)
            for k, v := range entry.Items.(map[string]interface{}) {
                n.KVMap.Set(k, v.(string))
            }


        case gMSG_HEAD_REMOVE:
            log.Println("removing log entry", entry)
            for _, v := range entry.Items.([]interface{}) {
                n.KVMap.Remove(v.(string))
            }

    }
    n.setLastLogId(entry.Id)
    n.addLogCount()
}

// 日志自动保存处理
func (n *Node) logAutoSavingHandler() {
    t := gtime.Millisecond()
    for {
        if n.getLastLogId() != n.getLastSavedLogId() || gtime.Millisecond() - t > gLOG_REPL_AUTOSAVE_INTERVAL {
            //log.Println("saving data to file")
            n.saveData()
            t = gtime.Millisecond()
        } else {
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func (n *Node) saveData() {
    var data SaveInfo
    data.LastLogId = n.getLastLogId()
    data.LogCount  = n.getLogCount()
    data.Peers     = *n.Peers.Clone()
    data.DataMap   = *n.KVMap.Clone()
    content       := gjson.Encode(&data)
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
                n.setLogCount(data.LogCount)
                n.setPeers(peerMap)
                n.setKVMap(dataMap)
            }
        }
    } else {
        //log.Println("no data file found at", path)
    }
}

func (n *Node) updateFromDataMapJson(jsonContent *string) error {
    m   := make(map[string]string)
    err := gjson.DecodeTo(jsonContent, &m)
    if err != nil {
        log.Println(err)
        return err
    }
    kvmap := gmap.NewStringStringMap()
    kvmap.BatchSet(m)
    n.setKVMap(kvmap)
    return nil
}

func (n *Node) getLastLogId() int64 {
    n.mutex.RLock()
    r := n.LastLogId
    n.mutex.RUnlock()
    return r
}

func (n *Node) getLogCount() int64 {
    n.mutex.RLock()
    r := n.LogCount
    n.mutex.RUnlock()
    return r
}

func (n *Node) getLastSavedLogId() int64 {
    n.mutex.Lock()
    r := n.LastSavedLogId
    n.mutex.Unlock()
    return r
}

func (n *Node) getStatusInReplication() bool {
    n.mutex.RLock()
    r := isInReplication
    n.mutex.RUnlock()
    return r
}

func (n *Node) setLastLogId(id int64) {
    n.mutex.Lock()
    n.LastLogId = id
    n.mutex.Unlock()
}

func (n *Node) setLastSavedLogId(id int64) {
    n.mutex.Lock()
    n.LastSavedLogId = id
    n.mutex.Unlock()
}

func (n *Node) setLogCount(c int64) {
    n.mutex.Lock()
    n.LogCount = c
    n.mutex.Unlock()
}

func (n *Node) setStatusInReplication(status bool ) {
    n.mutex.Lock()
    isInReplication = status
    n.mutex.Unlock()
}

func (n *Node) setPeers(m *gmap.StringInterfaceMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.Peers = m
    n.mutex.Unlock()
}

func (n *Node) setKVMap(m *gmap.StringStringMap) {
    if m == nil {
        return
    }
    n.mutex.Lock()
    n.KVMap = m
    n.mutex.Unlock()
}

func (n *Node) addLogCount() {
    n.mutex.Lock()
    n.LogCount++
    n.mutex.Unlock()
}
